package database

import (
	"database/sql"
	"fmt"
	"github.com/yoonhyunwoo/cloudtoggle/pkg/models"
	"log"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL 드라이버
)

// DB는 데이터베이스 연결을 나타내는 구조체입니다.
type DB struct {
	Conn *sql.DB
}

// InitDB는 데이터베이스 연결을 초기화합니다.
func InitDB(dbURL string) (*DB, error) {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 연결이 정상적인지 확인
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Connected to the database successfully")
	return &DB{Conn: conn}, nil
}

// GetAllGroups는 모든 리소스 그룹을 반환합니다.
func (db *DB) GetAllGroups() ([]map[string]interface{}, error) {
	rows, err := db.Conn.Query("SELECT id, name, status FROM resource_groups")
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %v", err)
	}
	defer rows.Close()

	var groups []map[string]interface{}
	for rows.Next() {
		var id, name, status string
		if err := rows.Scan(&id, &name, &status); err != nil {
			return nil, fmt.Errorf("failed to scan group: %v", err)
		}
		groups = append(groups, map[string]interface{}{
			"id":     id,
			"name":   name,
			"status": status,
		})
	}

	return groups, nil
}

func (db *DB) GetGroupByID(groupID string) (map[string]interface{}, error) {
	gid, err := strconv.Atoi(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid groupID: %v", err)
	}

	// LEFT JOIN을 사용하여 해당 그룹에 리소스가 없더라도 그룹 정보는 조회할 수 있도록 함
	rows, err := db.Conn.Query(`
        SELECT rg.id, rg.name, rg.status, rgr.resource_type, rgr.tag_key, rgr.tag_value
        FROM resource_groups rg
        LEFT JOIN resource_group_resources rgr ON rg.id = rgr.group_id
        WHERE rg.id = $1
    `, gid)
	if err != nil {
		return nil, fmt.Errorf("failed to query group and resources: %v", err)
	}
	defer rows.Close()

	var (
		groupIDVal  int
		groupName   string
		groupStatus string
		resources   []map[string]interface{}
		foundGroup  bool
	)

	for rows.Next() {
		var (
			resourceType sql.NullString
			tagKey       sql.NullString
			tagValue     sql.NullString
		)

		if err := rows.Scan(&groupIDVal, &groupName, &groupStatus, &resourceType, &tagKey, &tagValue); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// 그룹 정보는 첫 행에서 가져오고 이후 행에서도 동일하므로 별도 처리 필요 없음
		// 단지 첫 루프 진입 시 그룹 존재 여부를 확인
		foundGroup = true

		// 리소스 정보가 있는 경우에만 리소스 리스트에 추가
		if resourceType.Valid && tagKey.Valid && tagValue.Valid {
			resources = append(resources, map[string]interface{}{
				"resource_type": resourceType.String,
				"tag_key":       tagKey.String,
				"tag_value":     tagValue.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	// 그룹이 없는 경우
	if !foundGroup {
		return nil, fmt.Errorf("group not found")
	}

	// 최종 결과 반환
	return map[string]interface{}{
		"id":        groupIDVal,
		"name":      groupName,
		"status":    groupStatus,
		"resources": resources,
	}, nil
}

// RecordAction은 리소스 그룹의 작업 기록을 데이터베이스에 저장합니다.
func (db *DB) RecordAction(actionID, groupID, actionType string) error {
	query := "INSERT INTO action_logs (action_id, group_id, action_type) VALUES ($1, $2, $3)"
	_, err := db.Conn.Exec(query, actionID, groupID, actionType)
	if err != nil {
		return fmt.Errorf("failed to record action: %v", err)
	}
	return nil
}

// GroupExists는 특정 그룹이 데이터베이스에 존재하는지 확인합니다.
func (db *DB) GroupExists(groupID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM resource_groups WHERE id = $1)"
	err := db.Conn.QueryRow(query, groupID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if group exists: %v", err)
	}
	return exists, nil
}

// GetActionStatus는 특정 작업의 상태를 반환합니다.
func (db *DB) GetActionStatus(actionID string) (map[string]interface{}, error) {
	query := "SELECT action_id, group_id, action_type, created_at FROM action_logs WHERE action_id = $1"
	row := db.Conn.QueryRow(query, actionID)

	var actionIDRes, groupID, actionType, createdAt string
	err := row.Scan(&actionIDRes, &groupID, &actionType, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("action not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan action status: %v", err)
	}

	return map[string]interface{}{
		"action_id":  actionIDRes,
		"group_id":   groupID,
		"action":     actionType,
		"created_at": createdAt,
	}, nil
}

// Close는 데이터베이스 연결을 닫습니다.
func (db *DB) Close() error {
	if db.Conn != nil {
		log.Println("Closing database connection...")
		err := db.Conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close the database connection: %v", err)
		}
		log.Println("Database connection closed successfully")
		return nil
	}
	return nil
}

// AddResourceGroup는 새로운 리소스 그룹을 데이터베이스에 추가하고 생성된 ID를 반환합니다.
func (db *DB) AddResourceGroup(name, status string, resources []models.AWSResource) (int, error) {
	tx, err := db.Conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}

	var groupID int
	query := "INSERT INTO resource_groups (name, status) VALUES ($1, $2) RETURNING id"
	err = tx.QueryRow(query, name, status).Scan(&groupID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to add resource group: %v", err)
	}

	for _, resource := range resources {
		for _, tag := range resource.Tags {
			query := `
				INSERT INTO resource_group_resources (group_id, resource_type, tag_key, tag_value)
				VALUES ($1, $2, $3, $4)
			`
			_, err = tx.Exec(query, groupID, resource.Type, tag.Key, tag.Value)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to add resources to group: %v", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return groupID, nil
}

func (db *DB) DeleteResourceGroup(groupID string) error {
	query := "DELETE FROM resource_groups WHERE id = $1"
	_, err := db.Conn.Exec(query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete resource group: %v", err)
	}
	return nil
}
