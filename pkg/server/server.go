package server

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/scheduler"

	"github.com/gorilla/mux"
	"github.com/yoonhyunwoo/cloudtoggle/internal/auth"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
)

func StartServer(scheduler *scheduler.Scheduler, db *database.DB) {

	InitializeAdminPassword()

	router := mux.NewRouter()

	// API 경로 및 핸들러 연결
	router.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/resource-groups", auth.Middleware(AddResourceGroupHandler(db))).Methods("POST")
	router.HandleFunc("/api/v1/resource-groups/{group_id}", auth.Middleware(DeleteResourceGroupHandler(db))).Methods("DELETE")
	router.HandleFunc("/api/v1/groups", auth.Middleware(GetGroupsHandler(db))).Methods("GET")
	router.HandleFunc("/api/v1/groups/{group_id}", auth.Middleware(GetGroupHandler(db))).Methods("GET")
	router.HandleFunc("/api/v1/groups/{group_id}/start", auth.Middleware(StartGroupHandler(scheduler, db))).Methods("POST")
	router.HandleFunc("/api/v1/groups/{group_id}/stop", auth.Middleware(StopGroupHandler(scheduler, db))).Methods("POST")
	router.HandleFunc("/api/v1/groups/{group_id}/schedule", auth.Middleware(ScheduleGroupHandler(scheduler, db))).Methods("POST")
	router.HandleFunc("/api/v1/actions/{action_id}", auth.Middleware(GetActionStatusHandler(db))).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),                   // 허용할 클라이언트 URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // 허용할 HTTP 메서드
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // 허용할 헤더
	)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler(router))
}
