package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"

	"log"
)

// TODO : 현재 Start 작동 안함, GetbyTags시 리전 및 Account ID 종속되지 않게, 혹은 자동으로 받아오게 해야 함

type RDSManager struct {
	client *rds.Client
}

func NewRDSManager(client *rds.Client) *RDSManager {
	return &RDSManager{
		client: client,
	}
}

func (r *RDSManager) Start(ctx context.Context, dbInstanceIdentifiers []string) error {
	for _, dbInstance := range dbInstanceIdentifiers {
		log.Printf("Starting RDS instance: %s", dbInstance)

		_, err := r.client.StartDBInstance(ctx, &rds.StartDBInstanceInput{
			DBInstanceIdentifier: aws.String(dbInstance),
		})
		if err != nil {
			return fmt.Errorf("failed to start RDS instance %s: %w", dbInstance, err)
		}

		log.Printf("Successfully started RDS instance: %s", dbInstance)
	}
	return nil
}

func (r *RDSManager) Stop(ctx context.Context, dbInstanceIdentifiers []string) error {
	for _, dbInstance := range dbInstanceIdentifiers {
		log.Printf("Stopping RDS instance: %s", dbInstance)

		_, err := r.client.StopDBInstance(ctx, &rds.StopDBInstanceInput{
			DBInstanceIdentifier: aws.String(dbInstance),
		})
		if err != nil {
			return fmt.Errorf("failed to stop RDS instance %s: %w", dbInstance, err)
		}

		log.Printf("Successfully stopped RDS instance: %s", dbInstance)
	}
	return nil
}

func (r *RDSManager) GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error) {
	log.Printf("Retrieving RDS instances with tags: %v", resourceTags)

	// List all RDS instances
	var allInstances []string
	input := &rds.DescribeDBInstancesInput{}
	for {
		output, err := r.client.DescribeDBInstances(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to describe RDS instances: %w", err)
		}

		for _, dbInstance := range output.DBInstances {
			allInstances = append(allInstances, *dbInstance.DBInstanceIdentifier)
		}

		if output.Marker == nil {
			break
		}
		input.Marker = output.Marker
	}

	log.Printf("Retrieved all RDS instances: %v", allInstances)

	var matchingInstances []string
	for _, instanceIdentifier := range allInstances {
		// Retrieve tags for the current RDS instance
		tagInput := &rds.ListTagsForResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:rds:ap-northeast-2:593634833876:db:%s", instanceIdentifier)),
		}
		tagOutput, err := r.client.ListTagsForResource(ctx, tagInput)
		if err != nil {
			log.Printf("failed to get tags for RDS instance %s: %v", instanceIdentifier, err)
			continue
		}

		// Check if the instance's tags match the provided tags
		if r.matchTags(resourceTags, tagOutput.TagList) {
			matchingInstances = append(matchingInstances, instanceIdentifier)
		}
	}

	log.Printf("Matching RDS instances: %v", matchingInstances)
	return matchingInstances, nil
}

// Helper function to check if tags match
func (r *RDSManager) matchTags(resourceTags []models.ResourceTag, awsTags []types.Tag) bool {
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
