package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/scheduler"
)

func StopGroupHandler(scheduler *scheduler.Scheduler, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID := vars["group_id"]

		actionID, err := scheduler.StopGroup(groupID)
		if err != nil {
			log.Printf("Scheduler error: %v", err)
			http.Error(w, "Failed to stop group", http.StatusInternalServerError)
			return
		}

		db.RecordAction(actionID, groupID, "stop")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "success",
			"message":   "Group is stopping",
			"action_id": actionID,
		})
	}
}
