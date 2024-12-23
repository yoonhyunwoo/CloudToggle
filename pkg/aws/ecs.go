package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type ECSManager struct {
	client *ecs.Client
}

func NewECSManager(client *ecs.Client) *ECSManager {
	return &ECSManager{client: client}
}

func (e *ECSManager) Start(ctx context.Context, serviceNames []string) error {
	log.Printf("Starting ECS services: %v", serviceNames)
	for _, serviceName := range serviceNames {
		_, err := e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
			Cluster:      aws.String("default"), // Replace with your cluster
			Service:      aws.String(serviceName),
			DesiredCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		log.Printf("Successfully started ECS service: %s", serviceName)
	}
	return nil
}

func (e *ECSManager) Stop(ctx context.Context, serviceNames []string) error {
	log.Printf("Stopping ECS services: %v", serviceNames)
	for _, serviceName := range serviceNames {
		_, err := e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
			Cluster:      aws.String("default"), // Replace with your cluster
			Service:      aws.String(serviceName),
			DesiredCount: aws.Int32(0),
		})
		if err != nil {
			return err
		}
		log.Printf("Successfully stopped ECS service: %s", serviceName)
	}
	return nil
}

func (e *ECSManager) GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error) {
	log.Printf("Retrieving ECS services with tags: %v", resourceTags)
	cluster := "default" // Replace with your cluster name

	output, err := e.client.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}

	var serviceNames []string
	serviceNames = append(serviceNames, output.ServiceArns...)
	log.Printf("Found ECS services: %v", serviceNames)
	return serviceNames, nil
}
