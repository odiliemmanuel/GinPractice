package models

import (
	"time"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)



type Task struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	Title       string    
	Description string    
	Status      Status 
	DueDate     *time.Time   
	CreatedAt   time.Time 
}


type CreateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate	 *time.Time    `json:"dueDate"`
}



type UpdateTaskInput struct {

	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status" `	
	DueDate	 *time.Time    `json:"dueDate"`

}


type DeleteTaskInput struct {
	ID    uint   `json:"id" binding:"required"`
}

type FilterTaskByStatusInput struct {
	Status    Status   `json:"status" binding:"required"`
}


type PatchTaskInput struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Status      *Status     `json:"status"`
	DueDate	 *time.Time    `json:"dueDate"`

}
