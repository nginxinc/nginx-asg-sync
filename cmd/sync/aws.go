package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	yaml "gopkg.in/yaml.v2"
)

const upstreamNameErrorMsg = "The mandatory field name is either empty or missing for an upstream in the config file"
const upstreamErrorMsgFormat = "The mandatory field %v is either empty or missing for the upstream %v in the config file"
const upstreamPortErrorMsgFormat = "The mandatory field port is either zero or missing for the upstream %v in the config file"
const upstreamKindErrorMsgFormat = "The mandatory field kind is either not equal to http or tcp or missing for the upstream %v in the config file"

// AWSClient allows you to get the list of IP addresses of instanes of an Auto Scaling group. It implements the CloudProvider interface
type AWSClient struct {
	svcEC2         ec2iface.EC2API
	svcAutoscaling autoscalingiface.AutoScalingAPI
	config         *awsConfig
}

// NewAWSClient creates and configures an AWSClient
func NewAWSClient(data []byte) (*AWSClient, error) {
	awsClient := &AWSClient{}
	cfg, err := parseAWSConfig(data)
	if err != nil {
		return nil, fmt.Errorf("error validating config: %v", err)
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
	for _, awsU := range client.config.Upstreams {
		u := Upstream{
			Name:         awsU.Name,
			Port:         awsU.Port,
			Kind:         awsU.Kind,
			ScalingGroup: awsU.AutoscalingGroup,
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

	svcAutoscaling := autoscaling.New(session)
	svcEC2 := ec2.New(session)
	client.svcEC2 = svcEC2
	client.svcAutoscaling = svcAutoscaling
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
	_, exists, err := client.getAutoscalingGroup(name)
	if err != nil {
		return exists, fmt.Errorf("couldn't check if an AutoScaling group exists: %v", err)
	}
	return exists, nil
}

// GetPrivateIPsForScalingGroup returns the list of IP addresses of instanes of the Auto Scaling group
func (client *AWSClient) GetPrivateIPsForScalingGroup(name string) ([]string, error) {
	group, exists, err := client.getAutoscalingGroup(name)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("autoscaling group %v doesn't exist", name)
	}

	instances, err := client.getInstancesOfAutoscalingGroup(group)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, ins := range instances {
		if len(ins.NetworkInterfaces) > 0 && ins.NetworkInterfaces[0].PrivateIpAddress != nil {
			result = append(result, *ins.NetworkInterfaces[0].PrivateIpAddress)
		}
	}

	return result, nil
}

func (client *AWSClient) getAutoscalingGroup(name string) (*autoscaling.Group, bool, error) {
	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(name),
		},
	}

	resp, err := client.svcAutoscaling.DescribeAutoScalingGroups(params)
	if err != nil {
		return nil, false, err
	}

	if len(resp.AutoScalingGroups) != 1 {
		return nil, false, nil
	}

	return resp.AutoScalingGroups[0], true, nil
}

func (client *AWSClient) getInstancesOfAutoscalingGroup(group *autoscaling.Group) ([]*ec2.Instance, error) {
	var result []*ec2.Instance

	if len(group.Instances) == 0 {
		return result, nil
	}

	var ids []*string
	for _, ins := range group.Instances {
		ids = append(ids, ins.InstanceId)
	}
	params := &ec2.DescribeInstancesInput{
		InstanceIds: ids,
	}

	resp, err := client.svcEC2.DescribeInstances(params)
	if err != nil {
		return result, err
	}
	for _, res := range resp.Reservations {
		result = append(result, res.Instances...)
	}

	return result, nil
}

// Configuration for AWS Cloud Provider

type awsConfig struct {
	Region    string
	Upstreams []awsUpstream
}

type awsUpstream struct {
	Name             string
	AutoscalingGroup string `yaml:"autoscaling_group"`
	Port             int
	Kind             string
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
	}

	return nil
}
