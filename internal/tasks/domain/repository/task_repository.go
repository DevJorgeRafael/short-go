package repository

import "go-task-easy-list/internal/tasks/domain/model"

type TaskRepository interface {
	Create(task *model.Task) error
	FindByUserID(userID string) ([]*model.Task, error)
	FindByID(id string) (*model.Task, error)
	Update(task *model.Task) error
	Delete(id string) error
}