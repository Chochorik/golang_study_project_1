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
	Id             int
	Title          string
	Description    string
	IsDone         bool
	DoneAt         *time.Time
	NotifyStatus   NotifyStatus
	NotifyAttempts int
	NotifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewTask(title, description string, doneAt *time.Time) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title is required")
	}

	return &Task{
		Title:        title,
		Description:  description,
		DoneAt:       doneAt,
		NotifyStatus: NotifyStatusPending,
	}, nil
}
