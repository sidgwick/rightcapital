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
	ID            string
	TargetAddress string
	Content       interface{}
	Semantic      DeliverySemantic
	MaxRetries    int
	Deadline      *time.Time
	CallbackURL   string
	CallbackMQ    string
	Status        TaskStatus
	RetryCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateTaskRequest struct {
	TargetAddress string
	Content       interface{}
	Semantic      DeliverySemantic
	MaxRetries    *int
	Deadline      *time.Time
	CallbackURL   string
	CallbackMQ    string
}

type TaskQueryResponse struct {
	ID            string
	TargetAddress string
	Content       interface{}
	Status        TaskStatus
	RetryCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CallbackRequest struct {
	TaskID   string
	Success  bool
	Message  string
	Metadata interface{}
}
