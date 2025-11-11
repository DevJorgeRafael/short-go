package repository

import "go-task-easy-list/internal/auth/domain/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Update(user *model.User) error
}