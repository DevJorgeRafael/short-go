package gorm

import (
	"short-go/internal/tasks/domain/model"
	derefUtils "short-go/internal/shared/http/utils"

	"gorm.io/gorm"
)

type TaskRepositoryGorm struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepositoryGorm {
	return &TaskRepositoryGorm{db: db}
}

func (r *TaskRepositoryGorm) Create(task *model.Task) error {
	taskModel := &TaskModel{
		ID: task.ID,
		UserID: task.UserID,
		Title: task.Title,
		Description: task.Description,
		StatusID: task.StatusID,
		PriorityID: task.PriorityID,
		StartsAt: &task.StartsAt,
		DueDate: &task.DueDate,
		CompletedAt: &task.CompletedAt,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	if err := r.db.Create(taskModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *TaskRepositoryGorm) FindByUserID(userID string) ([]*model.Task, error) {
	var taskModels []TaskModel
	if err := r.db.Where("user_id = ?", userID).Find(&taskModels).Error; err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for _, tm := range taskModels {
		task := &model.Task{
			ID:        tm.ID,
			UserID:    tm.UserID,
			Title:     tm.Title,
			StatusID:  tm.StatusID,
			PriorityID: tm.PriorityID,
			Description: tm.Description,
			StartsAt:  derefUtils.DerefTime(tm.StartsAt),
			DueDate:   derefUtils.DerefTime(tm.DueDate),
			CompletedAt: derefUtils.DerefTime(tm.CompletedAt),
			CreatedAt: tm.CreatedAt,
			UpdatedAt: tm.UpdatedAt,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepositoryGorm) FindByID(id string) (*model.Task, error) {
	var taskModel TaskModel
	if err := r.db.First(&taskModel, "id = ?", id).Error; err != nil {
		return nil, err
	}

	task := &model.Task{
		ID:        taskModel.ID,
		UserID:    taskModel.UserID,
		Title:     taskModel.Title,
		Description: taskModel.Description,
		StartsAt:  derefUtils.DerefTime(taskModel.StartsAt),
		DueDate:   derefUtils.DerefTime(taskModel.DueDate),
		CompletedAt: derefUtils.DerefTime(taskModel.CompletedAt),
		CreatedAt: taskModel.CreatedAt,
		UpdatedAt: taskModel.UpdatedAt,
	}

	return task, nil
}

func (r *TaskRepositoryGorm) Update(task *model.Task) error {
	taskModel := &TaskModel{
		ID: task.ID,
		UserID: task.UserID,
		Title: task.Title,
		Description: task.Description,
		StatusID: task.StatusID,
		PriorityID: task.PriorityID,
		StartsAt: &task.StartsAt,
		DueDate: &task.DueDate,
		CompletedAt: &task.CompletedAt,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	if err := r.db.Save(taskModel).Error; err != nil {
		return err
	}	
	return nil
}

func (r *TaskRepositoryGorm) Delete(id string) error {
	if err := r.db.Delete(&TaskModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryGorm) ChangeStatus(taskID string, statusID int) error {
	if err := r.db.Model(&TaskModel{}).Where("id = ?", taskID).Update("status_id", statusID).Error; err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryGorm) ChangePriority(taskID string, priorityID int) error {
	if err := r.db.Model(&TaskModel{}).Where("id = ?", taskID).Update("priority_id", priorityID).Error; err != nil {
		return err
	}
	return nil
}
