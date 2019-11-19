package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	yaml "gopkg.in/yaml.v2"
)

// AzureClient allows you to get the list of IP addresses of VirtualMachines of a VirtualMachine Scale Set. It implements the CloudProvider interface
type AzureClient struct {
	config      *azureConfig
	vMSSClient  compute.VirtualMachineScaleSetsClient
	iFaceClient network.InterfacesClient
}

// NewAzureClient creates an AzureClient
func NewAzureClient(data []byte) (*AzureClient, error) {
	azureClient := &AzureClient{}
	cfg, err := parseAzureConfig(data)
	if err != nil {
		return nil, fmt.Errorf("error validating config: %v", err)
	}

	azureClient.config = cfg

	err = azureClient.configure()
	if err != nil {
		return nil, fmt.Errorf("error configuring Azure Client: %v", err)
	}

	return azureClient, nil
}

// parseAzureConfig parses and validates AzureClient config
func parseAzureConfig(data []byte) (*azureConfig, error) {
	cfg := &azureConfig{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	err = validateAzureConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetPrivateIPsForScalingGroup returns the list of IP addresses of instances of the Virtual Machine Scale Set
func (client *AzureClient) GetPrivateIPsForScalingGroup(name string) ([]string, error) {
	var ips []string

	ctx := context.TODO()

	for iFaces, err := client.iFaceClient.ListVirtualMachineScaleSetNetworkInterfaces(ctx, client.config.ResourceGroupName, name); iFaces.NotDone() || err != nil; err = iFaces.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}

		for _, iFace := range iFaces.Values() {
			if iFace.VirtualMachine != nil && iFace.VirtualMachine.ID != nil && iFace.IPConfigurations != nil {
				for _, n := range *iFace.IPConfigurations {
					ip := getPrimaryIPFromInterfaceIPConfiguration(n)
					if ip != "" {
						ips = append(ips, *n.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
						break
					}
				}
			}
		}
	}
	return ips, nil
}

func getPrimaryIPFromInterfaceIPConfiguration(ipConfig network.InterfaceIPConfiguration) string {
	if ipConfig == (network.InterfaceIPConfiguration{}) {
		return ""
	}

	if ipConfig.Primary == nil {
		return ""
	}

	if !*ipConfig.Primary {
		return ""
	}

	if ipConfig.InterfaceIPConfigurationPropertiesFormat == nil {
		return ""
	}

	if ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress == nil {
		return ""
	}

	return *ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
}

// CheckIfScalingGroupExists checks if the Virtual Machine Scale Set exists
func (client *AzureClient) CheckIfScalingGroupExists(name string) (bool, error) {
	ctx := context.TODO()
	vmss, err := client.vMSSClient.Get(ctx, client.config.ResourceGroupName, name)
	if err != nil {
		return false, fmt.Errorf("couldn't check if a Virtual Machine Scale Set exists: %v", err)
	}

	return vmss.ID != nil, nil
}

func (client *AzureClient) configure() error {
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return err
	}

	client.vMSSClient = compute.NewVirtualMachineScaleSetsClient(client.config.SubscriptionID)
	client.vMSSClient.Authorizer = authorizer

	client.iFaceClient = network.NewInterfacesClient(client.config.SubscriptionID)
	client.iFaceClient.Authorizer = authorizer
	return nil
}

// GetUpstreams returns the Upstreams list
func (client *AzureClient) GetUpstreams() []Upstream {
	var upstreams []Upstream
	for i := 0; i < len(client.config.Upstreams); i++ {
		u := Upstream{
			Name:         client.config.Upstreams[i].Name,
			Port:         client.config.Upstreams[i].Port,
			Kind:         client.config.Upstreams[i].Kind,
			ScalingGroup: client.config.Upstreams[i].VMScaleSet,
			MaxConns:     &client.config.Upstreams[i].MaxConns,
			MaxFails:     &client.config.Upstreams[i].MaxFails,
			FailTimeout:  client.config.Upstreams[i].FailTimeout,
			SlowStart:    client.config.Upstreams[i].SlowStart,
		}
		upstreams = append(upstreams, u)
	}
	return upstreams
}

type azureConfig struct {
	SubscriptionID    string `yaml:"subscription_id"`
	ResourceGroupName string `yaml:"resource_group_name"`
	Upstreams         []*azureUpstream
}

type azureUpstream struct {
	Name        string
	VMScaleSet  string `yaml:"virtual_machine_scale_set"`
	Port        int
	Kind        string
	MaxConns    int    `yaml:"max_conns"`
	MaxFails    int    `yaml:"max_fails"`
	FailTimeout string `yaml:"fail_timeout"`
	SlowStart   string `yaml:"slow_start"`
}

func validateAzureConfig(cfg *azureConfig) error {
	if cfg.SubscriptionID == "" {
		return fmt.Errorf(errorMsgFormat, "subscription_id")
	}

	if cfg.ResourceGroupName == "" {
		return fmt.Errorf(errorMsgFormat, "resource_group_name")
	}

	if len(cfg.Upstreams) == 0 {
		return fmt.Errorf("There are no upstreams found in the config file")
	}

	for _, ups := range cfg.Upstreams {
		if ups.Name == "" {
			return fmt.Errorf(upstreamNameErrorMsg)
		}
		if ups.VMScaleSet == "" {
			return fmt.Errorf(upstreamErrorMsgFormat, "virtual_machine_scale_set", ups.Name)
		}
		if ups.Port == 0 {
			return fmt.Errorf(upstreamPortErrorMsgFormat, ups.Name)
		}
		if ups.Kind == "" || !(ups.Kind == "http" || ups.Kind == "stream") {
			return fmt.Errorf(upstreamKindErrorMsgFormat, ups.Name)
		}
		if ups.MaxConns < 0 {
			return fmt.Errorf(upstreamMaxConnsErrorMsgFmt, ups.MaxConns)
		}
		if ups.MaxFails < 0 {
			return fmt.Errorf(upstreamMaxFailsErrorMsgFmt, ups.MaxFails)
		}
		if !isValidTime(ups.FailTimeout) {
			return fmt.Errorf(upstreamFailTimeoutErrorMsgFmt, ups.FailTimeout)
		}
		if !isValidTime(ups.SlowStart) {
			return fmt.Errorf(upstreamSlowStartErrorMsgFmt, ups.SlowStart)
		}
	}
	return nil
}
