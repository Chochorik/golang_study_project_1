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
	Id             int          `json:"id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	IsDone         bool         `json:"is_done"`
	DoneAt         *time.Time   `json:"done_at"`
	NotifyStatus   NotifyStatus `json:"notify_status"`
	NotifyAttempts int          `json:"notify_attempts"`
	NotifiedAt     *time.Time   `json:"notified_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

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
