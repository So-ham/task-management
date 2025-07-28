package handlers

import (
	v1 "task-management/internal/handlers/v1"
	"task-management/internal/services"

	"github.com/go-playground/validator"
)

type Handler struct {
	V1 v1.HandlerV1
}

func New(service services.Service, validate *validator.Validate) *Handler {
	return &Handler{
		V1: v1.New(service, validate),
	}

}
