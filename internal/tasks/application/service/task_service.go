package service

import (
	"errors"
	"go-task-easy-list/internal/tasks/domain/model"
	"go-task-easy-list/internal/tasks/domain/repository"
	"time"

	"github.com/google/uuid"
)

// Errores del dominio
var (
	ErrInvalidTitle   = errors.New("el título no puede estar vacío")
	ErrTaskNotFound   = errors.New("tarea no encontrada")
	ErrUnauthorized   = errors.New("no autorizado para esta tarea")
	ErrInvalidDueDate = errors.New("fecha de vencimiento debe ser en el futuro")
	ErrInvalidDates   = errors.New("fecha de inicio no puede ser posterior a la fecha de vencimiento")
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(title, description string, statusID, priorityID int, startsAt, dueDate time.Time, userID string) (*model.Task, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}

	if !dueDate.IsZero() && dueDate.Before(time.Now()) {
		return nil, ErrInvalidDueDate
	}

	if !startsAt.IsZero() && !dueDate.IsZero() && startsAt.After(dueDate) {
		return nil, ErrInvalidDates
	}

	newTask := &model.Task{
		ID: uuid.New().String(),
		UserID: userID,
		Title: title,
		Description: description,
		StatusID: statusID,
		PriorityID: priorityID,
		StartsAt: startsAt,
		DueDate: dueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.taskRepo.Create(newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

func (s *TaskService) GetTasksByUserID(userID string) ([]*model.Task, error) {
	return s.taskRepo.FindByUserID(userID)
}

func (s *TaskService) GetTaskByID(id, userID string) (*model.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, ErrUnauthorized
	}

	return task, nil
}

func (s *TaskService) UpdateTask(updatedTask *model.Task, userID string) (*model.Task, error) {
	existingTask, err := s.taskRepo.FindByID(updatedTask.ID)
	if err != nil || existingTask == nil {
		return nil, ErrTaskNotFound
	}

	if existingTask.UserID != userID {
		return nil, ErrUnauthorized
	}

	if updatedTask.Title == "" {
		return nil, ErrInvalidTitle
	}

	if !updatedTask.DueDate.IsZero() && updatedTask.DueDate.Before(existingTask.CreatedAt) {
		return nil, ErrInvalidDueDate
	}

	taskResponse := &model.Task{
		ID:        existingTask.ID,
		UserID:    existingTask.UserID,
		Title:     updatedTask.Title,
		Description: updatedTask.Description,
		StatusID:  updatedTask.StatusID,
		PriorityID: updatedTask.PriorityID,
		StartsAt:  updatedTask.StartsAt,
		DueDate:   updatedTask.DueDate,
		CompletedAt: updatedTask.CompletedAt,
		CreatedAt: existingTask.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.taskRepo.Update(taskResponse); err != nil {
		return nil, err
	}

	return taskResponse, nil
}

func (s *TaskService) DeleteTask(id, userID string) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil || task == nil {
		return ErrTaskNotFound
	}
	if task.UserID != userID {
		return ErrUnauthorized
	}

	return s.taskRepo.Delete(id)
}

func (s *TaskService) ChangeStatus(taskID, userID string, statusID int) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil || task == nil {
		return ErrTaskNotFound
	}

	if task.UserID != userID {
		return ErrUnauthorized
	}

	task.StatusID = statusID
	task.UpdatedAt = time.Now()

	return s.taskRepo.Update(task)
}

func (s *TaskService) ChangePriority(taskID, userID string, priorityID int) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil || task == nil {
		return ErrTaskNotFound
	}

	if task.UserID != userID {
		return ErrUnauthorized
	}

	task.PriorityID = priorityID
	task.UpdatedAt = time.Now()

	return s.taskRepo.Update(task)
}