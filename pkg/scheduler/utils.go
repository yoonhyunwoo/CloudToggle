package scheduler

import (
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

// extractResources는 그룹 데이터에서 AWS 리소스 정보를 추출합니다.
func extractResources(groupData map[string]interface{}) []models.AWSResource {
	var resources []models.AWSResource

	rawResources, ok := groupData["resources"].([]map[string]interface{})
	if !ok {
		return resources
	}

	for _, r := range rawResources {
		rType, _ := r["resource_type"].(string)
		tagKey, _ := r["tag_key"].(string)
		tagValue, _ := r["tag_value"].(string)

		resources = append(resources, models.AWSResource{
			Type: rType,
			Tags: []models.ResourceTag{
				{Key: tagKey, Value: tagValue},
			},
		})
	}
	return resources
}
