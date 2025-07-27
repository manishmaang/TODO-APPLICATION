package models

import(
	"time"
)

type TodoSchema struct {
	TaskName    string    `json:"task_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiryTime  time.Time `json:"expiry_time"`
}
