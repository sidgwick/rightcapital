package deliver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		platform, ok := d.platformMgr.GetPlatform(task.Platform)
		if !ok {
			return fmt.Errorf("platform not found")
		}
		return platform.Deliver(ctx, task.TargetAddress, task.Content)
	})

	if err != nil {
		if err == semantic.ErrNoRetry {
			d.repo.UpdateTaskStatus(ctx, task.ID, model.TaskStatusFailed)
			d.enqueueCallback(task.ID, false, err.Error(), nil)
			return
		}

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(callback)
	if err != nil {
		log.Printf("[HTTP Callback] Failed to marshal callback for task %s: %v", callback.TaskID, err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[HTTP Callback] Failed to create request for task %s: %v", callback.TaskID, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[HTTP Callback] Failed to send callback for task %s to %s: %v", callback.TaskID, url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[HTTP Callback] Received non-2xx status %d for task %s", resp.StatusCode, callback.TaskID)
		return
	}

	log.Printf("[HTTP Callback] Successfully sent callback for task %s to %s", callback.TaskID, url)
}

func (d *Deliverer) sendMQCallback(mq string, callback *model.CallbackRequest) {
	jsonData, err := json.Marshal(callback)
	if err != nil {
		log.Printf("[MQ Callback] Failed to marshal callback for task %s: %v", callback.TaskID, err)
		return
	}

	log.Printf("[MQ Callback] Sending message to %s for task %s: %s", mq, callback.TaskID, string(jsonData))
}

func (d *Deliverer) GetCallbackQueue() <-chan *model.CallbackRequest {
	return d.callbackQueue
}
