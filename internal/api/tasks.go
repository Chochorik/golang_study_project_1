package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo-service/internal/models"
)

type TaskStore interface {
	Create(ctx context.Context, t *models.Task) error
	GetById(ctx context.Context, id int64) (*models.Task, error)
	List(ctx context.Context, isDone *bool) ([]models.Task, error)
	Update(ctx context.Context, id int64, upd models.TaskUpdate) (*models.Task, error)
	Delete(ctx context.Context, id int64) error
}

type handler struct {
	store TaskStore
}

type createTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
}

type updateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	IsDone      *bool      `json:"is_done,omitempty"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
}

func parseID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func (h *handler) handleStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrTaskNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	slog.Error("store error", "err", err)
	writeError(w, http.StatusInternalServerError, "internal error")
}

func (h *handler) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	task, err := models.NewTask(req.Title, req.Description, req.DoneAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.Create(r.Context(), task); err != nil {
		h.handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	var isDone *bool
	if v := r.URL.Query().Get("is_done"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid argument 'is_done'")
			return
		}
		isDone = &b
	}

	tasks, err := h.store.List(r.Context(), isDone)
	if err != nil {
		h.handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *handler) getById(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.store.GetById(r.Context(), id)
	if err != nil {
		h.handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateTaskRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			writeError(w, http.StatusBadRequest, models.ErrTaskTitleRequired.Error())
			return
		}
		req.Title = &trimmed
	}

	if req.IsDone != nil && *req.IsDone && req.DoneAt != nil {
		now := time.Now()
		req.DoneAt = &now
	}

	task, err := h.store.Update(r.Context(), id, models.TaskUpdate{
		Title:       req.Title,
		Description: req.Description,
		DoneAt:      req.DoneAt,
		IsDone:      req.IsDone,
	})
	if err != nil {
		h.handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		h.handleStoreError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
