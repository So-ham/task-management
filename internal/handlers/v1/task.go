package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-management/internal/entities"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get all tasks with pagination and optional status filter
// @Tags tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param status query string false "Task status filter" Enums(Pending,InProgress,Completed,Cancelled)
// @Success 200 {array} entities.TaskResponse
// @Failure 500 {object} map[string]string
// @Router /api/tasks [get]
func (h *handlerV1) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	pageParam := r.URL.Query().Get("page")
	if pageParam != "" {
		pageInt, err := strconv.Atoi(pageParam)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	pageSizeParam := r.URL.Query().Get("pageSize")
	if pageSizeParam != "" {
		pageSizeInt, err := strconv.Atoi(pageSizeParam)
		if err == nil && pageSizeInt > 0 {
			pageSize = pageSizeInt
		}
	}

	var status *entities.TaskStatus
	statusParam := r.URL.Query().Get("status")
	if statusParam != "" {
		taskStatus := entities.TaskStatus(statusParam)
		if taskStatus == entities.StatusPending ||
			taskStatus == entities.StatusInProgress ||
			taskStatus == entities.StatusCompleted ||
			taskStatus == entities.StatusCancelled {
			status = &taskStatus
		}
	}

	tasks, err := h.Service.GetAllTasks(r.Context(), page, pageSize, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID godoc
// @Summary Get a task by ID
// @Description Get a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} entities.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tasks/{id} [get]
func (h *handlerV1) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.Service.GetTaskByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body entities.TaskRequest true "Task request body"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks [post]
func (h *handlerV1) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req entities.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateTask(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateTask godoc
// @Summary Update an existing task
// @Description Update a task by its ID with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body entities.TaskRequest true "Task request body"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [put]
func (h *handlerV1) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req entities.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateTask(r.Context(), id, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [delete]
func (h *handlerV1) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteTask(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
