package scheduler

import (
	"log"

	"github.com/google/uuid"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/aws"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

// StartGroup은 특정 리소스 그룹의 인스턴스를 시작합니다.
// 그룹 ID를 사용해 리소스를 조회하고, 리소스 타입별로 시작 작업을 실행합니다.
func (s *Scheduler) StartGroup(resourceGroupID string) (string, error) {
	actionID := uuid.New().String()

	go func() {
		log.Printf("[Scheduler] Starting resources for group: %s", resourceGroupID)

		// 그룹 데이터를 가져와 리소스를 처리
		resources, err := s.getResourcesForGroup(resourceGroupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to get resources for group %s: %v", resourceGroupID, err)
			return
		}

		// 각 리소스에 대해 시작 작업 실행
		for _, resource := range resources {
			manager := s.getResourceManager(resource.Type)
			if manager == nil {
				log.Printf("[Scheduler] No manager found for resource type: %s", resource.Type)
				continue
			}

			resourceIDs, err := manager.GetByTags(s.Context, resource.Tags)
			if err != nil {
				log.Printf("[Scheduler] Failed to get resources for resource type %s: %v", resource.Type, err)
				continue
			}

			if len(resourceIDs) > 0 {
				err = manager.Start(s.Context, resourceIDs)
				if err != nil {
					log.Printf("[Scheduler] Failed to start resources for resource type %s: %v", resource.Type, err)
				} else {
					log.Printf("[Scheduler] Successfully started resources: %v", resourceIDs)
				}
			} else {
				log.Printf("[Scheduler] No matching resources found for resource type: %s", resource.Type)
			}
		}
	}()

	return actionID, nil
}

// StopGroup은 특정 리소스 그룹의 인스턴스를 중지합니다.
// 그룹 ID를 사용해 리소스를 조회하고, 리소스 타입별로 중지 작업을 실행합니다.
func (s *Scheduler) StopGroup(resourceGroupID string) (string, error) {
	actionID := uuid.New().String()

	go func() {
		log.Printf("[Scheduler] Stopping resources for group: %s", resourceGroupID)

		// 그룹 데이터를 가져와 리소스를 처리
		resources, err := s.getResourcesForGroup(resourceGroupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to get resources for group %s: %v", resourceGroupID, err)
			return
		}

		// 각 리소스에 대해 중지 작업 실행
		for _, resource := range resources {
			manager := s.getResourceManager(resource.Type)
			if manager == nil {
				log.Printf("[Scheduler] No manager found for resource type: %s", resource.Type)
				continue
			}

			resourceIDs, err := manager.GetByTags(s.Context, resource.Tags)
			if err != nil {
				log.Printf("[Scheduler] Failed to get resources for resource type %s: %v", resource.Type, err)
				continue
			}

			if len(resourceIDs) > 0 {
				err = manager.Stop(s.Context, resourceIDs)
				if err != nil {
					log.Printf("[Scheduler] Failed to stop resources for resource type %s: %v", resource.Type, err)
				} else {
					log.Printf("[Scheduler] Successfully stopped resources: %v", resourceIDs)
				}
			} else {
				log.Printf("[Scheduler] No matching resources found for resource type: %s", resource.Type)
			}
		}
	}()

	return actionID, nil
}

// getResourcesForGroup은 그룹 ID를 사용해 리소스 데이터를 가져옵니다.
func (s *Scheduler) getResourcesForGroup(groupID string) ([]models.AWSResource, error) {
	groupData, err := s.DB.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	return extractResources(groupData), nil
}

// getResourceManager는 리소스 유형에 맞는 매니저를 반환합니다.
// 현재 EC2와 ECS를 지원하며, 추가 리소스 매니저를 확장할 수 있습니다.
func (s *Scheduler) getResourceManager(resourceType string) aws.AWSResourceManager {
	switch resourceType {
	case "EC2":
		return s.AWSClient.EC2Manager
	case "ECS":
		return s.AWSClient.ECSManager
	default:
		return nil
	}
}
