package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/aws"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/database"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/scheduler"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/server"
)

func main() {
	// 1. 환경 변수 로드
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. 데이터베이스 연결 초기화
	db, err := database.InitDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. 스케줄러 초기화
	awsClient := aws.NewAWSClient()
	mainScheduler := scheduler.NewScheduler(db, awsClient)
	mainScheduler.Start()

	// 4. 서버 실행 (서버는 스케줄러와 데이터베이스를 의존성으로 가짐)
	server.StartServer(mainScheduler, db)
}
