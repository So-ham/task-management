package v1

import (
	"net/http"

	"task-management/internal/services"

	"github.com/go-playground/validator"
)

type handlerV1 struct {
	Service  services.Service
	Validate *validator.Validate
}

type HandlerV1 interface {

	// Task handlers
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

func New(s services.Service, v *validator.Validate) HandlerV1 {
	return &handlerV1{Service: s, Validate: v}
}
