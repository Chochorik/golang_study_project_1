package api

import (
	"net/http"
	"time"
)

func NewServer(addr string, store TaskStore) *http.Server {
	h := &handler{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("POST /tasks", h.create)
	mux.HandleFunc("GET /tasks", h.list)
	mux.HandleFunc("GET /tasks/{id}", h.getById)
	mux.HandleFunc("PATCH /tasks/{id}", h.update)
	mux.HandleFunc("DELETE /tasks/{id}", h.delete)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
