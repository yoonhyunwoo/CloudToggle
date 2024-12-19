package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
)

func GetGroupsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := db.GetAllGroups()
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to get groups", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
