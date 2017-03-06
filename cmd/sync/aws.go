package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// AWSClient allows you to get the list of IP addresses of instanes of an Auto Scaling group
type AWSClient struct {
	svcEC2         ec2iface.EC2API
	svcAutoscaling autoscalingiface.AutoScalingAPI
}

// NewAWSClient creates an AWSClient
func NewAWSClient(svcEC2 ec2iface.EC2API, svcAutoscaling autoscalingiface.AutoScalingAPI) *AWSClient {
	return &AWSClient{svcEC2, svcAutoscaling}
}

// CheckIfAutoscalingGroupExists checks if the Auto Scaling group exists
func (client *AWSClient) CheckIfAutoscalingGroupExists(name string) (bool, error) {
	_, exists, err := client.getAutoscalingGroup(name)
	return exists, err
}

// GetPrivateIPsOfInstancesOfAutoscalingGroup returns the list of IP addresses of instanes of the Auto Scaling group
func (client *AWSClient) GetPrivateIPsOfInstancesOfAutoscalingGroup(name string) ([]string, error) {
	group, exists, err := client.getAutoscalingGroup(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("autoscaling group %v doesn't exists", name)
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
		for _, ins := range res.Instances {
			result = append(result, ins)
		}
	}

	return result, nil
}
