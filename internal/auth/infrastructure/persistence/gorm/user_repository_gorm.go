package gorm

import (
	"go-task-easy-list/internal/auth/domain/model"
	"go-task-easy-list/internal/auth/domain/repository"

	"gorm.io/gorm"
)

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) Create(user *model.User) error {
	// Convert domain.User -> gorm.UserModel
	userModel := &UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		IsActive: user.IsActive,
	}

	// db.Create
	if err := r.db.Create(userModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryGorm) FindByEmail(email string) (*model.User, error) {
	// db.Where("email = ?", email).First(&userModel)
	userModel := &UserModel{}
	if err := r.db.Where("email = ?", email).First(userModel).Error; err != nil {
		return nil, err
	}

	// Convert gorm.UserModel -> domain.User
	user := &model.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Password: userModel.Password,
		Name:     userModel.Name,
		IsActive: userModel.IsActive,
	}
	return user, nil
}

func (r *UserRepositoryGorm) FindByID(id string) (*model.User, error) {
	// db.First(&userModel, "id = ?", id)
	userModel := &UserModel{}
	if err := r.db.First(userModel, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Convert gorm.UserModel -> domain.User
	user := &model.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Password: userModel.Password,
		Name:     userModel.Name,
		IsActive: userModel.IsActive,
	}

	return user, nil
}

func (r *UserRepositoryGorm) Update(user *model.User) error {
	// Convert domain.User -> gorm.UserModel
	userModel := &UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		IsActive: user.IsActive,
	}

	// db.Save(&userModel)
	if err := r.db.Save(userModel).Error; err != nil {
		return err
	}

	return nil
}