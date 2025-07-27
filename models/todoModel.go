package models

import(
	"time"
)

type TodoSchema struct {
	TaskName    string    `json:"task_name" validate:"required,min=5"`
	Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiryTime  time.Time `json:"expiry_time" validate:"required"`
	Username string `json:"username" validate:"required"`
}

// The struct tag validate:"required" is just a string — Go does not automatically use it unless you explicitly run validation using the validator library.
// Only when you call validator.New().Struct(data), the library reads those tags and checks them.

// Therefore I don't need to import any package right now.
