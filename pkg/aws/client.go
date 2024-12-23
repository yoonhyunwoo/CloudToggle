package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// AWSClient는 모든 AWS 리소스 매니저를 포함하는 클라이언트입니다.
type AWSClient struct {
	EC2Manager *EC2Manager // EC2 리소스 매니저
	ECSManager *ECSManager // ECS 리소스 매니저
}

// NewAWSClient는 모든 AWS 리소스 매니저를 초기화하여 AWSClient를 반환합니다.
func NewAWSClient() *AWSClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK configuration: %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	ecsClient := ecs.NewFromConfig(cfg)

	return &AWSClient{
		EC2Manager: NewEC2Manager(ec2Client), // EC2 매니저 초기화
		ECSManager: NewECSManager(ecsClient), // ECS 매니저 초기화
	}
}
