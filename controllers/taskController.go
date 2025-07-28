package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/manishmaang/TODO-APPLICATION/config"
	"github.com/manishmaang/TODO-APPLICATION/models"
	"net/http"
	"time"
)

func CreateTask(ctx *gin.Context) {
	var payload models.TodoSchema

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload.CreatedAt = time.Now()

	// Validate struct
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		// Convert error into readable form
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+": "+err.ActualTag())
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	// Set up context with timeout
	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	query := `INSERT INTO tasks (task_name, description, created_at, expiry_time, username) VALUES 
	($1, $2, $3, $4, $5) RETURNING id`
	err := config.DB.QueryRow(ct, query,
		payload.TaskName,
		payload.Description,
		payload.CreatedAt,
		payload.ExpiryTime,
		payload.Username,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"id":      id,
	})

}

func DeleteTask(ctx *gin.Context) {
	type TaskRequest struct {
		TaskName string `json:"task_name"`
	}

	var req TaskRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM tasks WHERE task_name = $1`
	_, err := config.DB.Exec(ct, query, req.TaskName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
