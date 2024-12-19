package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type AddResourceGroupRequest struct {
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	Resources []models.AWSResource `json:"resources"` // AWS 리소스 타입 참조
}

type AddResourceGroupResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func AddResourceGroupHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AddResourceGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "Group name is required", http.StatusBadRequest)
			return
		}

		groupID, err := db.AddResourceGroup(req.Name, req.Status, req.Resources)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to create resource group", http.StatusInternalServerError)
			return
		}

		response := AddResourceGroupResponse{
			ID:      groupID,
			Message: "Resource group created successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
