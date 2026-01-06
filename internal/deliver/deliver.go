package deliver

import (
	"context"
	"fmt"
	"notification-system/internal/model"
	"notification-system/internal/platform"
	"notification-system/internal/repository"
	"notification-system/internal/semantic"
	"time"
)

type Deliverer struct {
	repo          repository.Repository
	platformMgr   *platform.PlatformManager
	semanticMgr   *semantic.SemanticManager
	callbackQueue chan *model.CallbackRequest
}

func NewDeliverer(repo repository.Repository, platformMgr *platform.PlatformManager, semanticMgr *semantic.SemanticManager) *Deliverer {
	return &Deliverer{
		repo:          repo,
		platformMgr:   platformMgr,
		semanticMgr:   semanticMgr,
		callbackQueue: make(chan *model.CallbackRequest, 1000),
	}
}

func (d *Deliverer) Start(ctx context.Context) {
	go d.processDeliveries(ctx)
	go d.processCallbacks(ctx)
}

func (d *Deliverer) processDeliveries(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.deliverPendingTasks(ctx)
		}
	}
}

func (d *Deliverer) deliverPendingTasks(ctx context.Context) {
	tasks, err := d.repo.GetPendingTasks(ctx)
	if err != nil {
		return
	}

	for _, task := range tasks {
		d.deliverTask(ctx, task)
	}
}

func (d *Deliverer) deliverTask(ctx context.Context, task *model.NotificationTask) {
	d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusDelivering)

	handler, ok := d.semanticMgr.GetHandler(string(task.Semantic))
	if !ok {
		d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusFailed)
		d.enqueueCallback(task.ID, false, "unsupported semantic", nil)
		return
	}

	err := handler.Handle(ctx, task.ID, func() error {
		platform, ok := d.platformMgr.GetPlatform("default")
		if !ok {
			return fmt.Errorf("platform not found")
		}
		return platform.Deliver(ctx, task.TargetAddress, task.Content)
	})

	if err != nil {
		if task.RetryCount >= task.MaxRetries {
			d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusFailed)
			d.enqueueCallback(task.ID, false, "max retries exceeded", nil)
			return
		}

		task.RetryCount++
		d.repo.UpdateTaskRetryCount(ctx, task.ID, task.RetryCount)
		d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusPending)
		return
	}

	d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusSuccess)
	d.enqueueCallback(task.ID, true, "delivered successfully", nil)
}

func (d *Deliverer) enqueueCallback(taskID string, success bool, message string, metadata interface{}) {
	d.callbackQueue <- &model.CallbackRequest{
		TaskID:   taskID,
		Success:  success,
		Message:  message,
		Metadata: metadata,
	}
}

func (d *Deliverer) processCallbacks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case callback := <-d.callbackQueue:
			d.sendCallback(callback)
		}
	}
}

func (d *Deliverer) sendCallback(callback *model.CallbackRequest) {
	task, err := d.repo.GetTaskByID(context.Background(), callback.TaskID)
	if err != nil {
		return
	}

	if task.CallbackURL != "" {
		d.sendHTTPCallback(task.CallbackURL, callback)
	} else if task.CallbackMQ != "" {
		d.sendMQCallback(task.CallbackMQ, callback)
	}
}

func (d *Deliverer) sendHTTPCallback(url string, callback *model.CallbackRequest) {
}

func (d *Deliverer) sendMQCallback(mq string, callback *model.CallbackRequest) {
}

func (d *Deliverer) GetCallbackQueue() <-chan *model.CallbackRequest {
	return d.callbackQueue
}
