# **CloudToggle API**

CloudToggle은 AWS의 EC2, RDS 등 주요 리소스를 **그룹으로 관리**하고, **지정된 스케줄에 따라 자동으로 시작/중지**할 수 있는 **셀프 서비스 리소스 제어 시스템**입니다. **태그 기반의 필터링 기능**을 통해 특정 리소스를 그룹에 자동으로 추가할 수 있습니다.

---

## **📁 폴더 구조**
```
cloudtoggle/
├── cmd/                   # CLI 명령어 정의 (프로그램 진입점)
│   └── main.go            # 프로그램의 진입점
├── config/                # 설정 파일 (YAML, JSON 등)
├── internal/              # 외부에 노출되지 않는 내부 모듈
│   ├── auth/              # JWT 및 인증 관련 모듈
│   └── validator/         # 요청 유효성 검증 모듈
├── pkg/                   # 애플리케이션의 핵심 로직
│   ├── aws/               # AWS SDK 클라이언트 및 EC2 제어 로직
│   ├── models/            # 데이터베이스 및 요청/응답에 필요한 구조체 정의
│   ├── scheduler/         # 스케줄러 관련 코드 (스케줄 추가/수정/삭제/조회 기능)
│   ├── server/            # API 서버 핸들러 (라우터 및 엔드포인트 정의)
│   └── database/          # 데이터베이스 연결 및 쿼리 실행
├── migrations/            # DB 마이그레이션 SQL 파일 (초기 스키마, 데이터 삽입 등)
├── go.mod                 # Go 모듈 설정 파일
├── .env                   # 환경변수 파일 (AWS 키, DB URL 등)
└── Dockerfile             # Docker 이미지 빌드 파일
```

---

## **📋 API 엔드포인트 및 요청 예시**

| **HTTP 메서드** | **URL**                     | **설명**                          |
|-----------------|-----------------------------|-----------------------------------|
| **POST**        | `/api/v1/login`             | 관리자 로그인 (JWT 발급)          |
| **POST**        | `/api/v1/resource-groups`   | 리소스 그룹 추가                  |
| **DELETE**      | `/api/v1/resource-groups/{group_id}` | 리소스 그룹 삭제           |
| **POST**        | `/api/v1/groups/{group_id}/schedule` | 리소스 그룹의 스케줄 추가  |
| **POST**        | `/api/v1/groups/{group_id}/start` | 특정 리소스 그룹 시작          |
| **POST**        | `/api/v1/groups/{group_id}/stop`  | 특정 리소스 그룹 중지          |

---

## **📋 요청 예시 및 응답**

### **1️⃣ 관리자 로그인**
- **URL**: `POST /api/v1/login`
- **Request Body**:
  ```json
  {
      "username": "admin",
      "password": "<서버 시작 시 생성된 랜덤 비밀번호>"
  }
  ```
- **Response**:
  ```json
  {
      "token": "<JWT Token>"
  }
  ```

---

### **2️⃣ 리소스 그룹 추가**
- **URL**: `POST /api/v1/resource-groups`
- **Request Body**:
  ```json
  {
      "name": "Development Group",
      "status": "stopped",
      "resources": [
          {
              "type": "EC2",
              "tags": [
                  { "key": "Environment", "value": "Development" }
              ]
          },
          {
              "type": "RDS",
              "tags": [
                  { "key": "Team", "value": "Engineering" }
              ]
          }
      ]
  }
  ```
- **Response**:
  ```json
  {
      "id": 1,
      "message": "Resource group created successfully"
  }
  ```

---

### **3️⃣ 스케줄 추가**
- **URL**: `POST /api/v1/groups/{group_id}/schedule`
- **Request Body**:
  ```json
  {
      "start_time": "10:00",
      "stop_time": "19:00"
  }
  ```
- **Response**:
  ```json
  {
      "status": "success",
      "message": "Schedule successfully created for group"
  }
  ```

---

### **4️⃣ EC2 인스턴스 시작/중지**
AWS SDK를 사용하여 **EC2 인스턴스를 시작 및 중지**할 수 있습니다.

- **StartEC2Instances**: 특정 EC2 인스턴스를 시작합니다.
- **StopEC2Instances**: 특정 EC2 인스턴스를 중지합니다.

---

## **💡 추가 정보**

### **초기 관리자 비밀번호**
- 서버가 시작될 때 랜덤한 **관리자 비밀번호**가 생성됩니다.
- 이 비밀번호는 서버의 **로그에 출력**되며, 운영 환경에서는 보안에 유의해야 합니다.

### **환경 변수**
- **JWT_SECRET**: JWT 서명에 사용하는 시크릿 키.
- **DB_URL**: PostgreSQL 연결 URL.
- **AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY**: AWS SDK에서 사용하는 자격 증명.

### **AWS 권한 요구 사항**
AWS SDK에서 **EC2 인스턴스 시작/중지**를 수행하기 위해 다음 IAM 권한이 필요합니다.
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:StartInstances",
                "ec2:StopInstances"
            ],
            "Resource": "*"
        }
    ]
}
```

---

## **🛠️ 빌드 및 실행**

1. **환경 설정**
```
cp .env.example .env
```

2. **모듈 다운로드**
```
go mod tidy
```

3. **빌드 및 실행**
```
make dev  # PostgreSQL 데이터베이스 실행
make run  # CloudToggle 서버 실행
```

4. **Postman 요청 테스트**
- **로그인** → **Bearer Token 추가** → **API 요청 전송**

---

이 문서는 **CloudToggle API의 종합 문서**로, 개발자와 운영자가 빠르게 프로젝트의 구조와 기능을 이해할 수 있도록 작성되었습니다. 추가적인 수정이 필요하거나 새로운 기능이 추가되면 README.md도 함께 업데이트되어야 합니다.
