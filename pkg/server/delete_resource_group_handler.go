package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
)

type DeleteResourceGroupResponse struct {
	Message string `json:"message"`
}

// DeleteResourceGroupHandler는 특정 리소스 그룹을 삭제하는 핸들러입니다.
func DeleteResourceGroupHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 요청 메서드 확인
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// URL 경로에서 group_id 추출
		vars := mux.Vars(r)
		groupID := vars["group_id"]
		if groupID == "" {
			http.Error(w, "group_id is required", http.StatusBadRequest)
			return
		}

		// 데이터베이스에서 리소스 그룹 삭제
		err := db.DeleteResourceGroup(groupID)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to delete resource group", http.StatusInternalServerError)
			return
		}

		// 응답 생성
		response := DeleteResourceGroupResponse{
			Message: "Resource group deleted successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
