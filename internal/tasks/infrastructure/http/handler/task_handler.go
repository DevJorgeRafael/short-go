package handler

import (
	"encoding/json"
	sharedContext "go-task-easy-list/internal/shared/context"
	sharedhttp "go-task-easy-list/internal/shared/http"
	format "go-task-easy-list/internal/shared/http/utils"
	sharedValidation "go-task-easy-list/internal/shared/validation"
	"go-task-easy-list/internal/tasks/application/service"
	"go-task-easy-list/internal/tasks/domain/model"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	taskService *service.TaskService
	validator   *validator.Validate
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		validator:   sharedValidation.NewValidator(),
	}
}

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	StatusId    int    `json:"statusId" validate:"required,min=1,max=3"`
	PriorityId  int    `json:"priorityId" validate:"required,min=1,max=3"`
	StartsAt    string `json:"startsAt"`
	DueDate     string `json:"dueDate"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusId    int    `json:"statusId"`
	PriorityId  int    `json:"priorityId"`
	StartsAt    string `json:"startsAt"`
	DueDate     string `json:"dueDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateTask - POST /api/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "StartsAt inválido")
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "DueDate inválido")
		return
	}

	task, err := h.taskService.CreateTask(
		req.Title, 
		req.Description, 
		req.StatusId, 
		req.PriorityId, 
		startsAt, 
		dueDate, 
		userID,
	)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		StatusId:    task.StatusID,
		PriorityId:  task.PriorityID,
		StartsAt:    formatTime(task.StartsAt),
		DueDate:     formatTime(task.DueDate),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}

	sharedhttp.SuccessResponse(w, http.StatusCreated, resp)
}

// GetTasks - GET /api/tasks
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())

	tasks, err := h.taskService.GetTasksByUserID(userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusInternalServerError, "Error al obtener tareas")
		return
	}

	var resp []TaskResponse
	for _, task := range tasks {
		resp = append(resp, TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			StatusId:    task.StatusID,
			PriorityId:  task.PriorityID,
			StartsAt:    task.StartsAt.Format(time.RFC3339),
			DueDate:     task.DueDate.Format(time.RFC3339),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		})
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, resp)
}

// GetTask - GET /api/tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())
	taskID := chi.URLParam(r, "id")

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	resp := TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		StatusId:    task.StatusID,
		PriorityId:  task.PriorityID,
		StartsAt:    task.StartsAt.Format(time.RFC3339),
		DueDate:     task.DueDate.Format(time.RFC3339),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, resp)
}

// PUT /api/tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())
	taskID := chi.URLParam(r, "id")

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, format.FormatValidationError(err))
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "StartsAt inválido")
		return
	}
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, "DueDate inválido")
		return
	}

	updatedTask, err := h.taskService.UpdateTask(&model.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		StatusID:    req.StatusId,
		PriorityID:  req.PriorityId,
		StartsAt:    startsAt,
		DueDate:     dueDate,
	}, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return			
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, updatedTask)
}

// DELETE /api/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := sharedContext.GetUserID(r.Context())

	taskID := chi.URLParam(r, "id")
	err := h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedhttp.SuccessResponse(w, http.StatusNoContent, nil)
}

// ------------------------- HELPERS ------------------------- //
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}