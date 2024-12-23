package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
	"log"
)

type ECSManager struct {
	client *ecs.Client
	// TODO : NoSQL DB로 변경하기
	serviceTaskCounts map[string]int32
}

func NewECSManager(client *ecs.Client) *ECSManager {
	return &ECSManager{
		client:            client,
		serviceTaskCounts: make(map[string]int32),
	}
}

func (e *ECSManager) Start(ctx context.Context, clusterNames []string) error {
	for _, clusterName := range clusterNames {
		log.Printf("Starting ECS services in cluster: %s", clusterName)

		// 서비스 목록 가져오기
		serviceNames, err := listServices(e.client, clusterName)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		for _, serviceName := range serviceNames {
			taskCount := e.serviceTaskCounts[serviceName]
			if taskCount == 0 {
				taskCount = 1 // 기본적으로 1로 설정
			}

			_, err := e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
				Cluster:      aws.String(clusterName),
				Service:      aws.String(serviceName),
				DesiredCount: aws.Int32(taskCount),
			})
			if err != nil {
				return fmt.Errorf("failed to start service %s: %w", serviceName, err)
			}

			log.Printf("Successfully started ECS service: %s with task count: %d", serviceName, taskCount)
		}
	}
	return nil
}

func (e *ECSManager) Stop(ctx context.Context, clusterNames []string) error {
	for _, clusterName := range clusterNames {

		log.Printf("Stopping ECS services in cluster: %s", clusterName)

		// 서비스 목록 가져오기
		serviceNames, err := listServices(e.client, clusterName)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		for _, serviceName := range serviceNames {
			// 서비스의 현재 태스크 수 가져오기
			serviceDesc, err := e.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
				Cluster:  aws.String(clusterName),
				Services: []string{serviceName},
			})
			if err != nil {
				return fmt.Errorf("failed to describe service %s: %w", serviceName, err)
			}

			if len(serviceDesc.Services) > 0 {
				e.serviceTaskCounts[serviceName] = serviceDesc.Services[0].DesiredCount
			}

			// 태스크 수를 0으로 설정
			_, err = e.client.UpdateService(ctx, &ecs.UpdateServiceInput{
				Cluster:      aws.String(clusterName),
				Service:      aws.String(serviceName),
				DesiredCount: aws.Int32(0),
			})
			if err != nil {
				return fmt.Errorf("failed to stop service %s: %w", serviceName, err)
			}

			log.Printf("Successfully stopped ECS service: %s", serviceName)
		}
	}
	return nil
}

func (e *ECSManager) GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error) {
	log.Printf("Retrieving ECS clusters with tags: %v", resourceTags)

	// List all clusters with pagination
	var allClusters []string
	input := &ecs.ListClustersInput{}
	for {
		output, err := e.client.ListClusters(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list ECS clusters: %w", err)
		}
		allClusters = append(allClusters, output.ClusterArns...)

		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}

	log.Printf("Retrieved all clusters: %v", allClusters)

	var matchingClusters []string
	for _, clusterArn := range allClusters {
		// Retrieve tags for the current cluster
		tagInput := &ecs.ListTagsForResourceInput{
			ResourceArn: &clusterArn,
		}
		tagOutput, err := e.client.ListTagsForResource(ctx, tagInput)
		if err != nil {
			log.Printf("failed to get tags for cluster %s: %v", clusterArn, err)
			continue
		}

		// Check if the cluster's tags match the provided tags
		if matchTags(resourceTags, tagOutput.Tags) {
			matchingClusters = append(matchingClusters, clusterArn)
		}
	}

	log.Printf("Matching ECS clusters: %v", matchingClusters)
	return matchingClusters, nil
}

func listServices(client *ecs.Client, clusterName string) ([]string, error) {
	var services []string
	var nextToken *string

	for {
		output, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
			Cluster:    aws.String(clusterName),
			NextToken:  nextToken,
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return nil, err
		}

		services = append(services, output.ServiceArns...)
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return services, nil
}

// Helper function to check if tags match
func matchTags(resourceTags []models.ResourceTag, awsTags []types.Tag) bool {
	tagMap := make(map[string]string)
	for _, tag := range awsTags {
		tagMap[*tag.Key] = *tag.Value
	}

	for _, resourceTag := range resourceTags {
		if value, exists := tagMap[resourceTag.Key]; !exists || value != resourceTag.Value {
			return false
		}
	}
	return true
}
