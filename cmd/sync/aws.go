package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	yaml "gopkg.in/yaml.v2"
)

// AWSClient allows you to get the list of IP addresses of instanes of an Auto Scaling group. It implements the CloudProvider interface
type AWSClient struct {
	svcEC2 ec2iface.EC2API
	config *awsConfig
}

// NewAWSClient creates and configures an AWSClient
func NewAWSClient(data []byte) (*AWSClient, error) {
	awsClient := &AWSClient{}
	cfg, err := parseAWSConfig(data)
	if err != nil {
		return nil, fmt.Errorf("error validating config: %v", err)
	}

	if cfg.Region == "self" {
		httpClient := &http.Client{Timeout: connTimeoutInSecs * time.Second}
		params := &aws.Config{HTTPClient: httpClient}

		metaSession, err := session.NewSession(params)
		if err != nil {
			return nil, err
		}

		metaClient := ec2metadata.New(metaSession)
		if !metaClient.Available() {
			return nil, fmt.Errorf("ec2metadata service is unavailable")
		}

		region, err := metaClient.Region()
		if err != nil {
			return nil, fmt.Errorf("unable to retreive region from ec2metadata: %v", err)
		}
		cfg.Region = region
	}

	awsClient.config = cfg

	err = awsClient.configure()
	if err != nil {
		return nil, fmt.Errorf("error configuring AWS Client: %v", err)
	}

	return awsClient, nil
}

// GetUpstreams returns the Upstreams list
func (client *AWSClient) GetUpstreams() []Upstream {
	var upstreams []Upstream
	for i := 0; i < len(client.config.Upstreams); i++ {
		u := Upstream{
			Name:         client.config.Upstreams[i].Name,
			Port:         client.config.Upstreams[i].Port,
			Kind:         client.config.Upstreams[i].Kind,
			ScalingGroup: client.config.Upstreams[i].AutoscalingGroup,
			MaxConns:     &client.config.Upstreams[i].MaxConns,
			MaxFails:     &client.config.Upstreams[i].MaxFails,
			FailTimeout:  client.config.Upstreams[i].FailTimeout,
			SlowStart:    client.config.Upstreams[i].SlowStart,
		}
		upstreams = append(upstreams, u)
	}
	return upstreams
}

// configure configures the AWSClient with necessary parameters
func (client *AWSClient) configure() error {
	httpClient := &http.Client{Timeout: connTimeoutInSecs * time.Second}
	cfg := &aws.Config{Region: aws.String(client.config.Region), HTTPClient: httpClient}

	session, err := session.NewSession(cfg)
	if err != nil {
		return err
	}

	svcEC2 := ec2.New(session)
	client.svcEC2 = svcEC2
	return nil
}

// parseAWSConfig parses and validates AWSClient config
func parseAWSConfig(data []byte) (*awsConfig, error) {
	cfg := &awsConfig{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	err = validateAWSConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// CheckIfScalingGroupExists checks if the Auto Scaling group exists
func (client *AWSClient) CheckIfScalingGroupExists(name string) (bool, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:aws:autoscaling:groupName"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}

	response, err := client.svcEC2.DescribeInstances(params)
	if err != nil {
		return false, fmt.Errorf("couldn't check if an AutoScaling group exists: %v", err)
	}

	return len(response.Reservations) > 0, nil
}

// GetPrivateIPsForScalingGroup returns the list of IP addresses of instances of the Auto Scaling group
func (client *AWSClient) GetPrivateIPsForScalingGroup(name string) ([]string, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:aws:autoscaling:groupName"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}

	response, err := client.svcEC2.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	if len(response.Reservations) == 0 {
		return nil, fmt.Errorf("autoscaling group %v doesn't exist", name)
	}

	var result []string
	for _, res := range response.Reservations {
		for _, ins := range res.Instances {
			if len(ins.NetworkInterfaces) > 0 && ins.NetworkInterfaces[0].PrivateIpAddress != nil {
				result = append(result, *ins.NetworkInterfaces[0].PrivateIpAddress)
			}
		}
	}

	return result, nil
}

// Configuration for AWS Cloud Provider

type awsConfig struct {
	Region    string
	Upstreams []*awsUpstream
}

type awsUpstream struct {
	Name             string
	AutoscalingGroup string `yaml:"autoscaling_group"`
	Port             int
	Kind             string
	MaxConns         int    `yaml:"max_conns"`
	MaxFails         int    `yaml:"max_fails"`
	FailTimeout      string `yaml:"fail_timeout"`
	SlowStart        string `yaml:"slow_start"`
}

func validateAWSConfig(cfg *awsConfig) error {
	if cfg.Region == "" {
		return fmt.Errorf(errorMsgFormat, "region")
	}

	if len(cfg.Upstreams) == 0 {
		return fmt.Errorf("There are no upstreams found in the config file")
	}

	for _, ups := range cfg.Upstreams {
		if ups.Name == "" {
			return fmt.Errorf(upstreamNameErrorMsg)
		}
		if ups.AutoscalingGroup == "" {
			return fmt.Errorf(upstreamErrorMsgFormat, "autoscaling_group", ups.Name)
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
