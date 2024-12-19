package scheduler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/aws"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// JobInfo는 작업의 세부 정보를 담고 있습니다.
type JobInfo struct {
	ID        cron.EntryID
	Schedule  string
	NextRun   time.Time
	IsRunning bool
}

// Scheduler는 작업을 관리하는 스케줄러 구조체입니다.
type Scheduler struct {
	cron       *cron.Cron
	jobMutex   sync.Mutex
	jobEntries map[string]cron.EntryID
	AWSClient  *aws.AWSClient
	DB         *database.DB
}

// NewScheduler는 새로운 스케줄러 인스턴스를 생성합니다.
func NewScheduler(db *database.DB) *Scheduler {
	return &Scheduler{
		cron:       cron.New(cron.WithSeconds()), // 초 단위 스케줄을 지원
		jobEntries: make(map[string]cron.EntryID),
		AWSClient:  aws.NewAWSClient(),
		DB:         db,
	}
}

// ScheduleGroup은 리소스 그룹의 시작 및 중지 스케줄을 등록합니다.
func (s *Scheduler) ScheduleGroup(groupID, startTime, stopTime string) error {
	// 스케줄 시작 작업 추가
	startJobID, err := s.AddJob(startTime, func() {
		log.Printf("[Scheduler] Automatically starting group %s\n", groupID)

		_, err := s.StartGroup(groupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to start group %s: %v\n", groupID, err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule start job for group %s: %v", groupID, err)
	}

	// 스케줄 중지 작업 추가
	stopJobID, err := s.AddJob(stopTime, func() {
		log.Printf("[Scheduler] Automatically stopping group %s\n", groupID)
		_, err := s.StopGroup(groupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to stop group %s: %v\n", groupID, err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule stop job for group %s: %v", groupID, err)
	}

	log.Printf("[Scheduler] Scheduled group %s (Start: %s, Stop: %s) with job IDs (%d, %d)\n", groupID, startTime, stopTime, startJobID, stopJobID)
	return nil
}

// AddJob은 특정 시간에 작업을 실행하도록 작업을 추가합니다.
func (s *Scheduler) AddJob(schedule string, task func()) (cron.EntryID, error) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	entryID, err := s.cron.AddFunc(schedule, task)
	if err != nil {
		return 0, fmt.Errorf("failed to add job: %v", err)
	}

	s.jobEntries[fmt.Sprintf("%d", entryID)] = entryID
	return entryID, nil
}

// Start는 스케줄러를 시작하여 모든 작업을 활성화합니다.
func (s *Scheduler) Start() {
	log.Println("[Scheduler] Starting scheduler...")
	s.cron.Start()
}

// Stop은 스케줄러를 중지합니다.
func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Stopping scheduler...")
	s.cron.Stop()
}

// ListJobs는 현재 스케줄러에 등록된 모든 작업의 목록을 반환합니다.
func (s *Scheduler) ListJobs() []JobInfo {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	var jobs []JobInfo
	entries := s.cron.Entries()
	for _, entry := range entries {
		jobs = append(jobs, JobInfo{
			ID:        entry.ID,
			Schedule:  "cron style",
			NextRun:   entry.Next,
			IsRunning: entry.Valid(),
		})
	}

	return jobs
}

// StartEC2Instances는 특정 인스턴스를 시작합니다.
func (s *Scheduler) StartEC2Instances(instanceIDs []string) error {
	log.Printf("Scheduler: Starting EC2 instances: %v", instanceIDs)
	return s.AWSClient.StartInstances(instanceIDs)
}

// StopEC2Instances는 특정 인스턴스를 중지합니다.
func (s *Scheduler) StopEC2Instances(instanceIDs []string) error {
	log.Printf("Scheduler: Stopping EC2 instances: %v", instanceIDs)
	return s.AWSClient.StopInstances(instanceIDs)
}

// StartGroup은 특정 리소스 그룹의 EC2 인스턴스를 시작합니다.
func (s *Scheduler) StartGroup(resourceGroupID string) (string, error) {
	actionID := uuid.New().String()

	go func() {
		log.Printf("[Scheduler] Starting resources for group: %s at %s\n", resourceGroupID, time.Now().Format(time.RFC3339))
		resourceGroups, err := s.DB.GetGroupByID(resourceGroupID)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		// "resources" 키에서 자원 리스트를 가져오기
		rawResources, ok := resourceGroups["resources"].([]map[string]interface{})
		if !ok {
			log.Println("Invalid resources structure")
			return
		}

		// resource_type을 키로 태그 리스트를 모아둠
		resourceMap := make(map[string][]models.ResourceTag)
		for _, r := range rawResources {
			rType, _ := r["resource_type"].(string)
			tagKey, _ := r["tag_key"].(string)
			tagValue, _ := r["tag_value"].(string)

			resourceMap[rType] = append(resourceMap[rType], models.ResourceTag{
				Key:   tagKey,
				Value: tagValue,
			})
		}

		// 2. resourceMap을 바탕으로 resource 구조 리스트를 구성
		var resources []models.AWSResource
		for rType, tags := range resourceMap {
			resources = append(resources, models.AWSResource{
				Type: rType,
				Tags: tags,
			})
		}

		// 3. EC2 타입 리소스의 태그만 추출
		ec2TagList := []models.ResourceTag{}
		for _, resource := range resources {
			if resource.Type == "EC2" {
				for _, tag := range resource.Tags {
					ec2TagList = append(ec2TagList, models.ResourceTag{
						Key:   tag.Key,
						Value: tag.Value,
					})
				}
			}
		}

		// EC2 Type
		instanceIDs, err := s.AWSClient.GetEC2InstancesByTags(ec2TagList)
		if err != nil {
			log.Printf("[Scheduler] Failed to get EC2 instances for group %s: %v", resourceGroupID, err)
			return
		}

		if len(instanceIDs) > 0 {
			err = s.AWSClient.StartInstances(instanceIDs)
			if err != nil {
				log.Printf("[Scheduler] Failed to start EC2 instances: %v", err)
			} else {
				log.Printf("[Scheduler] Successfully started EC2 instances: %v", instanceIDs)
			}
		} else {
			log.Printf("[Scheduler] No matching EC2 instances found for group %s", resourceGroupID)
		}
	}()
	return actionID, nil
}

// StopGroup은 특정 리소스 그룹의 EC2 인스턴스를 중지합니다.
func (s *Scheduler) StopGroup(resourceGroupID string) (string, error) {
	actionID := uuid.New().String()

	go func() {
		log.Printf("[Scheduler] Stopping resources for group: %s at %s\n", resourceGroupID, time.Now().Format(time.RFC3339))

		resourceGroups, err := s.DB.GetGroupByID(resourceGroupID)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		// "resources" 키에서 자원 리스트를 가져오기
		rawResources, ok := resourceGroups["resources"].([]map[string]interface{})
		if !ok {
			log.Println("Invalid resources structure")
			return
		}

		// resource_type을 키로 태그 리스트를 모아둠
		resourceMap := make(map[string][]models.ResourceTag)
		for _, r := range rawResources {
			rType, _ := r["resource_type"].(string)
			tagKey, _ := r["tag_key"].(string)
			tagValue, _ := r["tag_value"].(string)

			resourceMap[rType] = append(resourceMap[rType], models.ResourceTag{
				Key:   tagKey,
				Value: tagValue,
			})
		}

		// 2. resourceMap을 바탕으로 resource 구조 리스트를 구성
		var resources []models.AWSResource
		for rType, tags := range resourceMap {
			resources = append(resources, models.AWSResource{
				Type: rType,
				Tags: tags,
			})
		}

		// 3. EC2 타입 리소스의 태그만 추출
		ec2TagList := []models.ResourceTag{}
		for _, resource := range resources {
			if resource.Type == "EC2" {
				for _, tag := range resource.Tags {
					ec2TagList = append(ec2TagList, models.ResourceTag{
						Key:   tag.Key,
						Value: tag.Value,
					})
				}
			}
		}

		// EC2 인스턴스 ID 조회
		instanceIDs, err := s.AWSClient.GetEC2InstancesByTags(ec2TagList)
		if err != nil {
			log.Printf("[Scheduler] Failed to get EC2 instances for group %s: %v", resourceGroupID, err)
			return
		}

		if len(instanceIDs) > 0 {
			err = s.AWSClient.StopInstances(instanceIDs)
			if err != nil {
				log.Printf("[Scheduler] Failed to stop EC2 instances: %v", err)
			} else {
				log.Printf("[Scheduler] Successfully stopped EC2 instances: %v", instanceIDs)
			}
		} else {
			log.Printf("[Scheduler] No matching EC2 instances found for group %s", resourceGroupID)
		}
	}()
	return actionID, nil
}
