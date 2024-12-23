package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type EC2Manager struct {
	client *ec2.Client
}

func NewEC2Manager(client *ec2.Client) *EC2Manager {
	return &EC2Manager{client: client}
}

func (e *EC2Manager) Start(ctx context.Context, instanceIDs []string) error {
	log.Printf("Starting EC2 instances: %v", instanceIDs)
	_, err := e.client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return err
	}
	log.Printf("Successfully started EC2 instances: %v", instanceIDs)
	return nil
}

func (e *EC2Manager) Stop(ctx context.Context, instanceIDs []string) error {
	log.Printf("Stopping EC2 instances: %v", instanceIDs)
	_, err := e.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return err
	}
	log.Printf("Successfully stopped EC2 instances: %v", instanceIDs)
	return nil
}

func (e *EC2Manager) GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error) {
	var instanceIDs []string
	filters := BuildTagFilters(resourceTags)

	output, err := e.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}
	return instanceIDs, nil
}
