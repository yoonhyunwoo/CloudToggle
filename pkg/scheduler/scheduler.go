package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/aws"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
)

// Scheduler는 작업을 관리하고 AWS 리소스를 제어하기 위한 구조체입니다.
type Scheduler struct {
	cron       *cron.Cron              // cron 스케줄러 인스턴스
	jobMutex   sync.Mutex              // 작업 등록/삭제 보호를 위한 Mutex
	jobEntries map[string]cron.EntryID // 작업 ID를 저장하는 맵
	AWSClient  *aws.AWSClient          // AWS 리소스 매니저 클라이언트
	DB         *database.DB            // 데이터베이스 클라이언트
	Context    context.Context         // 작업 실행 시 사용할 기본 Context
}

// NewScheduler는 새로운 Scheduler 인스턴스를 생성합니다.
func NewScheduler(db *database.DB, awsClient *aws.AWSClient) *Scheduler {
	return &Scheduler{
		cron:       cron.New(cron.WithSeconds()), // 초 단위 스케줄링을 지원
		jobEntries: make(map[string]cron.EntryID),
		AWSClient:  awsClient,
		DB:         db,
		Context:    context.TODO(), // 기본 컨텍스트 생성
	}
}

// Start는 스케줄러를 시작하여 등록된 작업을 실행합니다.
func (s *Scheduler) Start() {
	log.Println("[Scheduler] Starting scheduler...")
	s.cron.Start()
}

// Stop은 스케줄러를 중지하여 모든 작업 실행을 멈춥니다.
func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Stopping scheduler...")
	s.cron.Stop()
}

// ScheduleGroup은 특정 리소스 그룹에 대한 시작 및 중지 작업을 스케줄에 등록합니다.
// 그룹 ID, 시작 시간, 중지 시간을 입력받아 작업을 등록합니다.
func (s *Scheduler) ScheduleGroup(groupID, startTime, stopTime string) error {
	// 시작 작업 등록
	startJobID, err := s.addJob(startTime, func() {
		log.Printf("[Scheduler] Automatically starting group %s", groupID)
		_, err := s.StartGroup(groupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to start group %s: %v", groupID, err)
		}
	})
	if err != nil {
		return err
	}

	// 중지 작업 등록
	stopJobID, err := s.addJob(stopTime, func() {
		log.Printf("[Scheduler] Automatically stopping group %s", groupID)
		_, err := s.StopGroup(groupID)
		if err != nil {
			log.Printf("[Scheduler] Failed to stop group %s: %v", groupID, err)
		}
	})
	if err != nil {
		return err
	}

	log.Printf("[Scheduler] Scheduled group %s (Start: %s, Stop: %s) with job IDs (%d, %d)", groupID, startTime, stopTime, startJobID, stopJobID)
	return nil
}

// addJob은 스케줄러에 특정 시간에 실행할 작업을 등록합니다.
func (s *Scheduler) addJob(schedule string, task func()) (cron.EntryID, error) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	entryID, err := s.cron.AddFunc(schedule, task)
	if err != nil {
		return 0, err
	}

	s.jobEntries[fmt.Sprintf("%d", entryID)] = entryID
	return entryID, nil
}
