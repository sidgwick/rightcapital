package repository

import (
	"context"
	"notification-system/internal/model"
)

type Repository interface {
	CreateTask(ctx context.Context, task *model.NotificationTask) error
	GetTaskByID(ctx context.Context, id string) (*model.NotificationTask, error)
	UpdateTaskStatus(ctx context.Context, id string, status model.TaskStatus) error
	UpdateTaskRetryCount(ctx context.Context, id string, retryCount int) error
	GetPendingTasks(ctx context.Context) ([]*model.NotificationTask, error)
	DeleteTask(ctx context.Context, id string) error
	MarkTaskAsSuccessIfNotDelivered(ctx context.Context, id string) (bool, error)
}
