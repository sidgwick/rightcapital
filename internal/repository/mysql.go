package repository

import (
	"context"
	"encoding/json"
	"notification-system/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NotificationTaskDB struct {
	ID            string                 `gorm:"primaryKey;size:36"`
	TargetAddress string                 `gorm:"size:255;not null"`
	Content       string                 `gorm:"type:text"`
	Semantic      model.DeliverySemantic `gorm:"size:50;not null"`
	MaxRetries    int                    `gorm:"not null;default:3"`
	Deadline      *time.Time             `gorm:"type:datetime"`
	CallbackURL   string                 `gorm:"size:500"`
	CallbackMQ    string                 `gorm:"size:255"`
	Status        model.TaskStatus       `gorm:"size:50;not null;default:'pending'"`
	RetryCount    int                    `gorm:"not null;default:0"`
	CreatedAt     time.Time              `gorm:"autoCreateTime"`
	UpdatedAt     time.Time              `gorm:"autoUpdateTime"`
}

func (NotificationTaskDB) TableName() string {
	return "notification_tasks"
}

func (db *NotificationTaskDB) ToModel() *model.NotificationTask {
	return &model.NotificationTask{
		ID:            db.ID,
		TargetAddress: db.TargetAddress,
		Content:       db.Content,
		Semantic:      db.Semantic,
		MaxRetries:    db.MaxRetries,
		Deadline:      db.Deadline,
		CallbackURL:   db.CallbackURL,
		CallbackMQ:    db.CallbackMQ,
		Status:        db.Status,
		RetryCount:    db.RetryCount,
		CreatedAt:     db.CreatedAt,
		UpdatedAt:     db.UpdatedAt,
	}
}

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(dsn string) (*MySQLRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	repo := &MySQLRepository{db: db}

	if err := repo.AutoMigrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *MySQLRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&NotificationTaskDB{})
}

func (r *MySQLRepository) CreateTask(ctx context.Context, task *model.NotificationTask) error {
	contentJSON, err := json.Marshal(task.Content)
	if err != nil {
		return err
	}

	dbTask := &NotificationTaskDB{
		ID:            task.ID,
		TargetAddress: task.TargetAddress,
		Content:       string(contentJSON),
		Semantic:      task.Semantic,
		MaxRetries:    task.MaxRetries,
		Deadline:      task.Deadline,
		CallbackURL:   task.CallbackURL,
		CallbackMQ:    task.CallbackMQ,
		Status:        task.Status,
		RetryCount:    task.RetryCount,
	}

	return r.db.WithContext(ctx).Create(dbTask).Error
}

func (r *MySQLRepository) GetTaskByID(ctx context.Context, id string) (*model.NotificationTask, error) {
	var dbTask NotificationTaskDB
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbTask).Error
	if err != nil {
		return nil, err
	}

	return dbTask.ToModel(), nil
}

func (r *MySQLRepository) UpdateTaskStatus(ctx context.Context, id string, status model.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&NotificationTaskDB{}).
		Where("id = ?", id).
		Update("status", status).
		Update("updated_at", time.Now()).
		Error
}

func (r *MySQLRepository) UpdateTaskRetryCount(ctx context.Context, id string, retryCount int) error {
	return r.db.WithContext(ctx).Model(&NotificationTaskDB{}).
		Where("id = ?", id).
		Update("retry_count", retryCount).
		Update("updated_at", time.Now()).
		Error
}

func (r *MySQLRepository) GetPendingTasks(ctx context.Context) ([]*model.NotificationTask, error) {
	var dbTasks []NotificationTaskDB
	err := r.db.WithContext(ctx).
		Where("status = ?", model.TaskStatusPending).
		Find(&dbTasks).
		Error

	if err != nil {
		return nil, err
	}

	tasks := make([]*model.NotificationTask, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = dbTask.ToModel()
	}

	return tasks, nil
}

func (r *MySQLRepository) DeleteTask(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&NotificationTaskDB{}, "id = ?", id).Error
}

func (r *MySQLRepository) MarkTaskAsSuccessIfNotDelivered(ctx context.Context, id string) (bool, error) {
	var dbTask NotificationTaskDB
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Raw("SELECT * FROM notification_tasks WHERE id = ? FOR UPDATE", id).Scan(&dbTask).Error
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	if dbTask.Status == model.TaskStatusSuccess {
		tx.Rollback()
		return false, nil
	}

	result := tx.Model(&NotificationTaskDB{}).
		Where("id = ?", id).
		Update("status", model.TaskStatusSuccess).
		Update("updated_at", time.Now())

	if result.Error != nil {
		tx.Rollback()
		return false, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
}
