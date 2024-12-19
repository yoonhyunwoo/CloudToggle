package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
)

func GetActionStatusHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		actionID := vars["action_id"]

		status, err := db.GetActionStatus(actionID)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Action not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
