package service

import (
	"context"
	"notification-system/internal/deliver"
	"notification-system/internal/model"
	"notification-system/internal/repository"
	"notification-system/internal/semantic"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo        repository.Repository
	deliverer   *deliver.Deliverer
	semanticMgr *semantic.SemanticManager
}

func NewService(repo repository.Repository, deliverer *deliver.Deliverer, semanticMgr *semantic.SemanticManager) *Service {
	return &Service{
		repo:        repo,
		deliverer:   deliverer,
		semanticMgr: semanticMgr,
	}
}

func (s *Service) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (string, error) {
	maxRetries := 3
	if req.MaxRetries != nil {
		maxRetries = *req.MaxRetries
	}

	task := &model.NotificationTask{
		ID:            uuid.New().String(),
		TargetAddress: req.TargetAddress,
		Content:       req.Content,
		Semantic:      req.Semantic,
		MaxRetries:    maxRetries,
		Deadline:      req.Deadline,
		CallbackURL:   req.CallbackURL,
		CallbackMQ:    req.CallbackMQ,
		Status:        model.TaskStatusPending,
		RetryCount:    0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return "", err
	}

	return task.ID, nil
}

func (s *Service) GetTask(ctx context.Context, taskID string) (*model.TaskQueryResponse, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return &model.TaskQueryResponse{
		ID:            task.ID,
		TargetAddress: task.TargetAddress,
		Content:       task.Content,
		Status:        task.Status,
		RetryCount:    task.RetryCount,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}, nil
}

func (s *Service) CancelTask(ctx context.Context, taskID string) error {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.Status == model.TaskStatusSuccess {
		return nil
	}

	return s.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusCancelled)
}

func (s *Service) Start(ctx context.Context) {
	s.deliverer.Start(ctx)
}
