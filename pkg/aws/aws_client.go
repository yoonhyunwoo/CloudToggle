package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type AWSClient struct {
	EC2Client *ec2.Client
}

// NewAWSClient는 AWS SDK 클라이언트를 생성합니다.
func NewAWSClient() *AWSClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK configuration: %v", err)
	}

	return &AWSClient{
		EC2Client: ec2.NewFromConfig(cfg),
	}
}

// StartInstances는 EC2 인스턴스를 시작합니다.
func (c *AWSClient) StartInstances(instanceIDs []string) error {
	log.Printf("Starting EC2 instances: %v", instanceIDs)

	input := &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := c.EC2Client.StartInstances(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Printf("Successfully started EC2 instances: %v", instanceIDs)
	return nil
}

// StopInstances는 EC2 인스턴스를 중지합니다.
func (c *AWSClient) StopInstances(instanceIDs []string) error {
	log.Printf("Stopping EC2 instances: %v", instanceIDs)

	input := &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := c.EC2Client.StopInstances(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Printf("Successfully stopped EC2 instances: %v", instanceIDs)
	return nil
}

// GetEC2InstancesByTags는 태그에 해당하는 EC2 인스턴스를 조회합니다.
func (c *AWSClient) GetEC2InstancesByTags(resourceTags []models.ResourceTag) ([]string, error) {
	var instanceIDs []string

	// GetEC2InstancesByTags 함수 내
	var filters []types.Filter
	for _, resourceTag := range resourceTags {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:" + resourceTag.Key),
			Values: []string{resourceTag.Value},
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	output, err := c.EC2Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	// 인스턴스 ID 수집
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	return instanceIDs, nil
}
