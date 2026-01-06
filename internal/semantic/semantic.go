package semantic

import (
	"context"
	"errors"
)

type SemanticHandler interface {
	Handle(ctx context.Context, taskID string, deliverFunc func() error) error
}

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
		return errors.New("delivery failed, no retry for at_most_once semantic")
	}
	return nil
}

type ExactlyOnceHandler struct {
	delivered map[string]bool
}

func NewExactlyOnceHandler() *ExactlyOnceHandler {
	return &ExactlyOnceHandler{
		delivered: make(map[string]bool),
	}
}

func (h *ExactlyOnceHandler) Handle(ctx context.Context, taskID string, deliverFunc func() error) error {
	if h.delivered[taskID] {
		return errors.New("task already delivered")
	}

	err := deliverFunc()
	if err != nil {
		return errors.New("delivery failed")
	}

	h.delivered[taskID] = true
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
