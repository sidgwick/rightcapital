package semantic

import (
	"context"
	"errors"
	"notification-system/internal/model"
	"notification-system/internal/repository"
)

type SemanticHandler interface {
	Handle(ctx context.Context, taskID string, deliverFunc func() error) error
}

var ErrNoRetry = errors.New("delivery failed, no retry")

type AtLeastOnceHandler struct{}

func NewAtLeastOnceHandler() *AtLeastOnceHandler {
	return &AtLeastOnceHandler{}
}

func (h *AtLeastOnceHandler) Handle(ctx context.Context, taskID string, deliverFunc func() error) error {
	err := deliverFunc()
	if err != nil {
		return errors.New("delivery failed, retry needed")
	}
	return nil
}

type AtMostOnceHandler struct{}

func NewAtMostOnceHandler() *AtMostOnceHandler {
	return &AtMostOnceHandler{}
}

func (h *AtMostOnceHandler) Handle(ctx context.Context, taskID string, deliverFunc func() error) error {
	err := deliverFunc()
	if err != nil {
		return ErrNoRetry
	}
	return nil
}

type ExactlyOnceHandler struct {
	repo repository.Repository
}

func NewExactlyOnceHandler(repo repository.Repository) *ExactlyOnceHandler {
	return &ExactlyOnceHandler{
		repo: repo,
	}
}

func (h *ExactlyOnceHandler) Handle(ctx context.Context, taskID string, deliverFunc func() error) error {
	marked, err := h.repo.MarkTaskAsSuccessIfNotDelivered(ctx, taskID)
	if err != nil {
		return err
	}

	if !marked {
		return errors.New("task already delivered")
	}

	deliverErr := deliverFunc()
	if deliverErr != nil {
		h.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusPending)
		return errors.New("delivery failed, task status reverted")
	}

	return nil
}

type SemanticManager struct {
	handlers map[string]SemanticHandler
}

func NewSemanticManager() *SemanticManager {
	return &SemanticManager{
		handlers: make(map[string]SemanticHandler),
	}
}

func (sm *SemanticManager) RegisterHandler(semantic string, handler SemanticHandler) {
	sm.handlers[semantic] = handler
}

func (sm *SemanticManager) GetHandler(semantic string) (SemanticHandler, bool) {
	handler, ok := sm.handlers[semantic]
	return handler, ok
}
