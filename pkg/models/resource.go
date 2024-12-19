package models

// AWS 리소스를 정의하는 구조체
type AWSResource struct {
	Type string        `json:"type"` // EC2, RDS, S3 등 AWS 리소스 유형
	Tags []ResourceTag `json:"tags"` // 리소스를 필터링할 태그 목록
}

// 리소스 태그를 정의하는 구조체
type ResourceTag struct {
	Key   string `json:"key"`   // 태그 키 (예: Environment)
	Value string `json:"value"` // 태그 값 (예: Development)
}

// 리소스 그룹 구조체 정의
type ResourceGroup struct {
	Name      string        `json:"name"`      // 리소스 그룹의 이름
	Resources []AWSResource `json:"resources"` // 리소스 목록 (EC2, RDS 등)
}
