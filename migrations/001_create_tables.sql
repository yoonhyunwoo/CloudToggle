-- SQL Script for Initial Database Schema

-- 리소스 그룹 테이블
CREATE TABLE IF NOT EXISTS resource_groups (
                                               id SERIAL PRIMARY KEY,
                                               name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'stopped',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 리소스 그룹과 AWS 리소스 간의 관계 테이블
CREATE TABLE IF NOT EXISTS resource_group_resources (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES resource_groups(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL, -- EC2, RDS, S3, DynamoDB 등
    tag_key VARCHAR(50) NOT NULL,
    tag_value VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 작업 로그 테이블
CREATE TABLE IF NOT EXISTS action_logs (
                                           action_id UUID PRIMARY KEY,
                                           group_id INT REFERENCES resource_groups(id) ON DELETE CASCADE,
    action_type VARCHAR(20) NOT NULL, -- start, stop 등 작업의 유형
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 작업 상태 테이블
CREATE TABLE IF NOT EXISTS job_status (
                                          id SERIAL PRIMARY KEY,
                                          action_id UUID REFERENCES action_logs(action_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL, -- in_progress, completed, failed
    message TEXT, -- 작업 관련 메시지 (예: 에러 원인)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 인덱스 추가
CREATE INDEX IF NOT EXISTS idx_resource_groups_name ON resource_groups (name);
CREATE INDEX IF NOT EXISTS idx_action_logs_group_id ON action_logs (group_id);
CREATE INDEX IF NOT EXISTS idx_job_status_action_id ON job_status (action_id);
