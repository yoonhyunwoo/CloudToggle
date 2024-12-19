package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWTSecret = "supersecretkey" // 환경변수로 관리하는 것이 좋습니다.

// Middleware는 API 요청에 대한 인증을 처리하는 미들웨어입니다.
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Authorization 헤더 확인
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Bearer 토큰 추출
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. JWT 검증
		claims, err := validateJWT(token)
		if err != nil {
			http.Error(w, "401 Unauthorized  Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. 사용자 정보 컨텍스트에 추가
		r = r.WithContext(WithUserContext(r.Context(), claims))

		// 5. 다음 핸들러 호출
		next.ServeHTTP(w, r)
	}
}

// validateJWT는 JWT 토큰의 유효성을 검증하고 클레임(Claims)을 반환합니다.
func validateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	// JWT 파서 초기화
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HMAC 서명 방법 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Claims를 파싱하여 반환
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, http.ErrAbortHandler
	}
}

// GenerateJWT는 사용자 정보를 바탕으로 JWT를 생성합니다.
func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24시간 후 만료
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// WithUserContext는 사용자 정보를 컨텍스트에 추가합니다.
func WithUserContext(ctx context.Context, claims *jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, "user", claims)
}

// GetUserFromContext는 컨텍스트에서 사용자 정보를 가져옵니다.
func GetUserFromContext(ctx context.Context) *jwt.RegisteredClaims {
	if user, ok := ctx.Value("user").(*jwt.RegisteredClaims); ok {
		return user
	}
	return nil
}
