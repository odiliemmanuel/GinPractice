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
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}


type CreateTaskInput struct {
	Title       string     `json:"title"   binding:"required"`
	Description string     `json:"description"`
}



type UpdateTaskInput struct {

	Title       string     `json:"title"   binding:"required"`
	Description string     `json:"description"`
	Status      Status     `json:"status" `	

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
	Status      *Status     `json:"status" `

}
