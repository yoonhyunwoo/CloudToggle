package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type AWSResourceManager interface {
	Start(ctx context.Context, resourceIDs []string) error
	Stop(ctx context.Context, resourceIDs []string) error
	GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error)
}

// 공통 태그 필터링 함수
func BuildTagFilters(resourceTags []models.ResourceTag) []types.Filter {
	var filters []types.Filter
	for _, tag := range resourceTags {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:" + tag.Key), // "tag:Key" 형식의 필터 이름
			Values: []string{tag.Value},          // 필터 값 리스트
		})
	}
	return filters
}
