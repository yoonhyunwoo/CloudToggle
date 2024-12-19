package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/scheduler"
)

type ScheduleRequest struct {
	StartTime string `json:"start_time"`
	StopTime  string `json:"stop_time"`
}

func ScheduleGroupHandler(scheduler *scheduler.Scheduler, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID := vars["group_id"]

		var req ScheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// ✅ HH:mm -> cron 포맷으로 변환
		startCron, err := convertToCron(req.StartTime)
		if err != nil {
			log.Printf("Invalid start time format: %v", err)
			http.Error(w, "Invalid start time format", http.StatusBadRequest)
			return
		}

		stopCron, err := convertToCron(req.StopTime)
		if err != nil {
			log.Printf("Invalid stop time format: %v", err)
			http.Error(w, "Invalid stop time format", http.StatusBadRequest)
			return
		}

		// 스케줄 추가 로직
		err = scheduler.ScheduleGroup(groupID, startCron, stopCron)
		if err != nil {
			log.Printf("Failed to schedule group %s: %v", groupID, err)
			http.Error(w, "Failed to create schedule", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"status":  "success",
			"message": "Schedule successfully created for group " + groupID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// convertToCron은 HH:mm 형식을 cron 포맷으로 변환합니다.
func convertToCron(hhmm string) (string, error) {
	parts := strings.Split(hhmm, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid time format: %s", hhmm)
	}

	hour := parts[0]
	minute := parts[1]

	// cron 포맷: 초 분 시 일 월 요일 (예: "0 30 2 * * ?")
	return fmt.Sprintf("0 %s %s * * ?", minute, hour), nil
}
