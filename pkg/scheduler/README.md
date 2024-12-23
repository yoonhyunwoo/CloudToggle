# **스케줄러: 중지 가능한 리소스 추가 방법**

이 문서는 스케줄러에서 새로운 AWS 리소스 타입을 지원하도록 추가하는 방법을 설명합니다. 스케줄러는 리소스 타입(예: EC2, ECS)에 따라 리소스를 시작 및 중지하기 위해 리소스 매니저를 사용합니다. 아래 단계를 따라 새로운 리소스 타입을 통합하세요.

---

## **1. 리소스 매니저 정의하기**

스케줄러는 리소스와 상호작용하기 위해 `AWSResourceManager` 인터페이스를 사용합니다. 새로운 리소스 매니저는 이 인터페이스를 구현해야 합니다.

### **`AWSResourceManager` 인터페이스**
```go
type AWSResourceManager interface {
    Start(ctx context.Context, resourceIDs []string) error
    Stop(ctx context.Context, resourceIDs []string) error
    GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error)
}
```

예시: RDS 리소스 매니저 추가하기

1.	pkg/aws 디렉토리에 새로운 파일을 생성합니다(예: rds.go).
2.	아래와 같이 AWSResourceManager 인터페이스를 구현합니다.

```go
package aws

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/rds"
    "log"
    "github.com/yoonhyunwoo/cloudtoggle/pkg/models"
)

type RDSManager struct {
    client *rds.Client
}

func NewRDSManager(client *rds.Client) *RDSManager {
    return &RDSManager{client: client}
}

func (r *RDSManager) Start(ctx context.Context, resourceIDs []string) error {
    log.Printf("Starting RDS instances: %v", resourceIDs)
    for _, id := range resourceIDs {
        _, err := r.client.StartDBInstance(ctx, &rds.StartDBInstanceInput{
            DBInstanceIdentifier: &id,
        })
        if err != nil {
            return err
        }
        log.Printf("Successfully started RDS instance: %s", id)
    }
    return nil
}

func (r *RDSManager) Stop(ctx context.Context, resourceIDs []string) error {
    log.Printf("Stopping RDS instances: %v", resourceIDs)
    for _, id := range resourceIDs {
        _, err := r.client.StopDBInstance(ctx, &rds.StopDBInstanceInput{
            DBInstanceIdentifier: &id,
        })
        if err != nil {
            return err
        }
        log.Printf("Successfully stopped RDS instance: %s", id)
    }
    return nil
}

func (r *RDSManager) GetByTags(ctx context.Context, resourceTags []models.ResourceTag) ([]string, error) {
    // 태그를 기반으로 RDS 인스턴스를 가져오는 로직 추가
    log.Println("Fetching RDS instances by tags is not implemented yet")
    return nil, nil
}
```
2. 리소스 매니저 등록하기
```go
pkg/aws/client.go 파일을 수정하여 새로운 리소스 매니저를 포함합니다.

package aws

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/rds"
    "log"
)

type AWSClient struct {
    EC2Manager *EC2Manager
    ECSManager *ECSManager
    RDSManager *RDSManager // 새로운 RDSManager 추가
}

func NewAWSClient() *AWSClient {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatalf("Unable to load AWS SDK configuration: %v", err)
    }

    return &AWSClient{
        EC2Manager: NewEC2Manager(ec2.NewFromConfig(cfg)),
        ECSManager: NewECSManager(ecs.NewFromConfig(cfg)),
        RDSManager: NewRDSManager(rds.NewFromConfig(cfg)), // RDSManager 초기화
    }
}

3. 스케줄러에 통합하기

pkg/scheduler/group_operations.go 파일에서 getResourceManager 함수를 수정하여 새로운 리소스 타입을 인식하도록 추가합니다.

func (s *Scheduler) getResourceManager(resourceType string) aws.AWSResourceManager {
    switch resourceType {
    case "EC2":
        return s.AWSClient.EC2Manager
    case "ECS":
        return s.AWSClient.ECSManager
    case "RDS": // 새로운 RDS 리소스 타입 추가
        return s.AWSClient.RDSManager
    default:
        return nil
    }
}
```
4. 데이터베이스에 리소스 타입 추가하기
```go 
데이터베이스 스키마가 새로운 리소스 타입(RDS 등)을 지원하도록 업데이트해야 합니다. 데이터베이스 또는 샘플 데이터를 수정하세요.

예시:

{
    "group_id": "example-group",
    "resources": [
        {
            "resource_type": "RDS",
            "tag_key": "Environment",
            "tag_value": "Production"
        }
    ]
}
```
5. 통합 테스트
	1.	스케줄러를 실행하여 새로운 리소스 타입이 포함된 그룹 작업을 스케줄링합니다.
	2.	로그나 AWS 콘솔을 통해 리소스가 올바르게 시작 및 중지되는지 확인합니다.


다른 리소스 타입을 추가하려면 다음 단계를 따르세요:
1.	AWSResourceManager 인터페이스 구현: 새로운 리소스 매니저를 작성합니다.
2.	pkg/aws/client.go에 등록: 클라이언트에 새로운 매니저를 추가합니다.
3.	getResourceManager 함수 수정: 스케줄러에서 리소스 타입을 인식하도록 업데이트합니다.
4.	데이터베이스 수정: 새로운 리소스 타입에 대한 데이터를 지원하도록 수정합니다.

이 가이드를 통해 스케줄러에 새로운 AWS 리소스 타입을 쉽게 추가할 수 있습니다. 🎉

