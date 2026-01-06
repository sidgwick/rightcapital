package repository

import (
	"context"
	"notification-system/internal/model"
	"time"
)

type Repository interface {
	CreateTask(ctx context.Context, task *model.NotificationTask) error
	GetTaskByID(ctx context.Context, id string) (*model.NotificationTask, error)
	UpdateTaskStatus(ctx context.Context, id string, status model.TaskStatus) error
	UpdateTaskRetryCount(ctx context.Context, id string, retryCount int) error
	GetPendingTasks(ctx context.Context) ([]*model.NotificationTask, error)
	DeleteTask(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	tasks map[string]*model.NotificationTask
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		tasks: make(map[string]*model.NotificationTask),
	}
}

func (r *InMemoryRepository) CreateTask(ctx context.Context, task *model.NotificationTask) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryRepository) GetTaskByID(ctx context.Context, id string) (*model.NotificationTask, error) {
	return r.tasks[id], nil
}

func (r *InMemoryRepository) UpdateTaskStatus(ctx context.Context, id string, status model.TaskStatus) error {
	if task, ok := r.tasks[id]; ok {
		task.Status = status
		task.UpdatedAt = time.Now()
	}
	return nil
}

func (r *InMemoryRepository) UpdateTaskRetryCount(ctx context.Context, id string, retryCount int) error {
	if task, ok := r.tasks[id]; ok {
		task.RetryCount = retryCount
		task.UpdatedAt = time.Now()
	}
	return nil
}

func (r *InMemoryRepository) GetPendingTasks(ctx context.Context) ([]*model.NotificationTask, error) {
	var tasks []*model.NotificationTask
	for _, task := range r.tasks {
		if task.Status == model.TaskStatusPending {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *InMemoryRepository) DeleteTask(ctx context.Context, id string) error {
	delete(r.tasks, id)
	return nil
}
