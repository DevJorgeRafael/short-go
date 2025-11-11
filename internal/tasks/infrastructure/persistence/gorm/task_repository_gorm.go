package gorm

import (
	"go-task-easy-list/internal/tasks/domain/model"
	"time"

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
			StartsAt:  derefTime(tm.StartsAt),
			DueDate:   derefTime(tm.DueDate),
			CompletedAt: derefTime(tm.CompletedAt),
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
		StartsAt:  derefTime(taskModel.StartsAt),
		DueDate:   derefTime(taskModel.DueDate),
		CompletedAt: derefTime(taskModel.CompletedAt),
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

// ------------------- Helper ---------------------
func derefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}