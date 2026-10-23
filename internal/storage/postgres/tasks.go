package postgres

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"todo-service/internal/models"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *models.Task) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetById(ctx context.Context, id int64) (*models.Task, error) {
	var t models.Task

	err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrTaskNotFound
		}

		return nil, fmt.Errorf("get task %d: %w", id, err)
	}

	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context, isDone *bool) ([]models.Task, error) {
	tasks := make([]models.Task, 0)

	q := r.db.WithContext(ctx).Order("id")
	if isDone != nil {
		q = q.Where("is_done = ?", *isDone)
	}

	if err := q.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, upd models.TaskUpdate) (*models.Task, error) {
	fields := map[string]any{}

	if upd.Title != nil {
		fields["title"] = *upd.Title
	}
	if upd.Description != nil {
		fields["description"] = *upd.Description
	}
	if upd.IsDone != nil {
		fields["is_done"] = *upd.IsDone
		if *upd.IsDone {
			fields["done_at"] = gorm.Expr("COALESCE(done_at, now())")
		} else {
			fields["done_at"] = nil
		}
	}

	if len(fields) == 0 {
		return r.GetById(ctx, id)
	}

	var t models.Task

	res := r.db.WithContext(ctx).
		Model(&t).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(fields)

	if res.Error != nil {
		return nil, fmt.Errorf("update task %d: %w", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, models.ErrTaskNotFound
	}

	return &t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id)

	if res.Error != nil {
		return fmt.Errorf("delete task %d: %w", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}
