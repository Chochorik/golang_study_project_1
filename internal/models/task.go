package models

import (
	"errors"
	"strings"
	"time"
)

type NotifyStatus string

const (
	NotifyStatusPending NotifyStatus = "pending"
	NotifyStatusQueued  NotifyStatus = "queued"
	NotifyStatusSent    NotifyStatus = "sent"
	NotifyStatusFailed  NotifyStatus = "failed"
)

type Task struct {
	Id             int          `json:"id"              gorm:"column:id;primaryKey"`
	Title          string       `json:"title"           gorm:"column:title"`
	Description    string       `json:"description"     gorm:"column:description"`
	IsDone         bool         `json:"is_done"         gorm:"column:is_done"`
	DoneAt         *time.Time   `json:"done_at"         gorm:"column:done_at"`
	NotifyStatus   NotifyStatus `json:"notify_status"   gorm:"column:notify_status"`
	NotifyAttempts int          `json:"notify_attempts" gorm:"column:notify_attempts"`
	NotifiedAt     *time.Time   `json:"notified_at"     gorm:"column:notified_at"`
	CreatedAt      time.Time    `json:"created_at"      gorm:"column:created_at"`
	UpdatedAt      time.Time    `json:"updated_at"      gorm:"column:updated_at"`
}

func (Task) TableName() string { return "tasks" }

var ErrTaskTitleRequired = errors.New("title is required")
var ErrTaskNotFound = errors.New("task not found")

type TaskUpdate struct {
	Title       *string
	Description *string
	IsDone      *bool
	DoneAt      *time.Time
}

func NewTask(title, description string, doneAt *time.Time) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTaskTitleRequired
	}

	return &Task{
		Title:        title,
		Description:  description,
		DoneAt:       doneAt,
		NotifyStatus: NotifyStatusPending,
	}, nil
}
