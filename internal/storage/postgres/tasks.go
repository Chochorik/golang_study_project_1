package postgres

import (
	"context"
	"errors"
	"fmt"
	"todo-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const taskColumns = `id, title, description, is_done, done_at, notify_status, notify_attempts, notified_at, created_at, updated_at`

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func ScanTask(row pgx.Row, t *models.Task) error {
	return row.Scan(
		&t.Id,
		&t.Title,
		&t.Description,
		&t.IsDone,
		&t.DoneAt,
		&t.NotifyStatus,
		&t.NotifyAttempts,
		&t.NotifiedAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
}

func (r *TaskRepository) Create(ctx context.Context, t *models.Task) error {
	const q = `INSERT INTO tasks (title, description, done_at)
		VALUES ($1, $2, $3)
		RETURNING ` + taskColumns

	row := r.pool.QueryRow(ctx, q, t.Title, t.Description, t.DoneAt)

	if err := ScanTask(row, t); err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetById(ctx context.Context, id int64) (*models.Task, error) {
	const q = `SELECT ` + taskColumns + ` FROM tasks WHERE id = $1`

	var t models.Task

	row := r.pool.QueryRow(ctx, q, id)

	if err := ScanTask(row, &t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrTaskNotFound
		}

		return nil, fmt.Errorf("get task %d: %w", id, err)
	}

	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context, isDone *bool) ([]models.Task, error) {
	const q = `SELECT ` + taskColumns + ` FROM tasks
		WHERE ($1::boolean IS NULL OR is_done = $1)
		ORDER BY id`

	rows, err := r.pool.Query(ctx, q, isDone)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		var t models.Task

		if err := ScanTask(rows, &t); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, upd models.TaskUpdate) (*models.Task, error) {
	const q = `UPDATE tasks SET
		title = COALESCE($2, title),
		description = COALESCE($3, description),
		is_done = COALESCE($4, is_done),
		done_at = COALESCE($5, done_at),
		updated_at = now()
		WHERE id = $1
		RETURNING ` + taskColumns

	var t models.Task

	row := r.pool.QueryRow(ctx, q, id, upd.Title, upd.Description, upd.IsDone, upd.DoneAt)

	if err := ScanTask(row, &t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrTaskNotFound
		}

		return nil, fmt.Errorf("update task %d: %w", id, err)
	}

	return &t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM tasks WHERE id = $1`

	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete task %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}
