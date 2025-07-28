package rest

import (
	"task-management/internal/handlers"

	"github.com/gorilla/mux"
)

// NewRouter returns a new router instance with configured routes
func NewRouter(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	// Task endpoints
	router.HandleFunc("/api/tasks", h.V1.GetAllTasks).Methods("GET")
	router.HandleFunc("/api/tasks", h.V1.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", h.V1.GetTaskByID).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", h.V1.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}", h.V1.DeleteTask).Methods("DELETE")

	return router
}
