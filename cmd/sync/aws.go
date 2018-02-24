package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// AWSClient allows you to get the list of IP addresses of instanes of an Auto Scaling group
type AWSClient struct {
	svcEC2         ec2iface.EC2API
}

// NewAWSClient creates an AWSClient
func NewAWSClient(svcEC2 ec2iface.EC2API) *AWSClient {
	return &AWSClient{svcEC2,}
}

// GetPrivateIPsOfInstancesOfAutoscalingGroup returns the list of IP addresses of instanes of the Auto Scaling group
func (client *AWSClient) GetPrivateIPsOfInstancesOfAutoscalingGroup(name string) ([]string, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:aws:autoscaling:groupName"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	resverations, err := client.svcEC2.DescribeInstances(params)

	if err != nil {
		return nil, err
	}
	if len(resverations.Reservations) == 0 {
		return nil, fmt.Errorf("autoscaling group %v doesn't exists", name)
	}

	instances, err := client.getInstancesOfReservations(resverations)

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


func (client *AWSClient) getInstancesOfReservations(group *ec2.DescribeInstancesOutput) ([]*ec2.Instance, error) {
	var result []*ec2.Instance

	if len(group.Reservations) == 0 {
		return result, nil
	}

	for _, res := range group.Reservations {
		for _, ins := range res.Instances {
			result = append(result, ins)
		}
	}

	return result, nil
}
