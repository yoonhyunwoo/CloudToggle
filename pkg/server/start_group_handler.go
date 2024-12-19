package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/scheduler"
)

func StartGroupHandler(scheduler *scheduler.Scheduler, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID := vars["group_id"]

		actionID, err := scheduler.StartGroup(groupID)
		if err != nil {
			log.Printf("Scheduler error: %v", err)
			http.Error(w, "Failed to start group", http.StatusInternalServerError)
			return
		}

		db.RecordAction(actionID, groupID, "start")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "success",
			"message":   "Group is starting",
			"action_id": actionID,
		})
	}
}
