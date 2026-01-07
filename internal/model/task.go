package model

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusDelivering TaskStatus = "delivering"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type DeliverySemantic string

const (
	DeliverySemanticAtLeastOnce DeliverySemantic = "at_least_once"
	DeliverySemanticAtMostOnce  DeliverySemantic = "at_most_once"
	DeliverySemanticExactlyOnce DeliverySemantic = "exactly_once"
)

type NotificationTask struct {
	ID            string           `json:"id" gorm:"primaryKey;size:36"`
	TargetAddress string           `json:"target_address" gorm:"size:255;not null"`
	Content       string           `json:"content" gorm:"type:text"`
	Semantic      DeliverySemantic `json:"semantic" gorm:"size:50;not null"`
	MaxRetries    int              `json:"max_retries" gorm:"not null;default:3"`
	Deadline      *time.Time       `json:"deadline" gorm:"type:datetime"`
	CallbackURL   string           `json:"callback_url" gorm:"size:500"`
	CallbackMQ    string           `json:"callback_mq" gorm:"size:255"`
	Status        TaskStatus       `json:"status" gorm:"size:50;not null;default:'pending'"`
	RetryCount    int              `json:"retry_count" gorm:"not null;default:0"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateTaskRequest struct {
	TargetAddress string           `json:"target_address"`
	Content       interface{}      `json:"content"`
	Semantic      DeliverySemantic `json:"semantic"`
	MaxRetries    *int             `json:"max_retries"`
	Deadline      *time.Time       `json:"deadline"`
	CallbackURL   string           `json:"callback_url"`
	CallbackMQ    string           `json:"callback_mq"`
}

type TaskQueryResponse struct {
	ID            string      `json:"id"`
	TargetAddress string      `json:"target_address"`
	Content       interface{} `json:"content"`
	Status        TaskStatus  `json:"status"`
	RetryCount    int         `json:"retry_count"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type CallbackRequest struct {
	TaskID   string      `json:"task_id"`
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata"`
}
