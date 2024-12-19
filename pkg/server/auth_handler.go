package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey") // 보안상의 이유로 환경 변수로 관리하는 것이 좋습니다.
var adminPassword string                 // 서버 시작 시 생성되는 랜덤 비밀번호

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// generateRandomToken는 32자리의 랜덤 문자열을 생성합니다.
func generateRandomToken() (string, error) {
	bytes := make([]byte, 16) // 16 바이트 -> 32자리 Hex 문자열
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

// InitializeAdminPassword는 서버 시작 시 호출되어 랜덤한 admin 비밀번호를 생성합니다.
func InitializeAdminPassword() {
	var err error
	adminPassword, err = generateRandomToken()
	if err != nil {
		log.Fatalf("Failed to generate admin password: %v", err)
	}
	log.Println("==============================")
	log.Printf("[INFO] Admin password: %s", adminPassword) // 서버 시작 시 비밀번호 출력
	log.Println("==============================")
	// 이 정보를 보안 로그로 남기지 않도록 주의해야 합니다.
}

// LoginHandler는 로그인 요청을 처리하고 JWT 토큰을 발급합니다.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest

	// 요청 본문 파싱
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 로그인 검증
	if loginRequest.Username != "admin" || loginRequest.Password != adminPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// JWT 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24시간 후 만료
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
