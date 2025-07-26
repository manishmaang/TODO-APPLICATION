package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/manishmaang/TODO-APPLICATION/config"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LogInUsers(ctx *gin.Context) {

}

func RegisterUsers(ctx *gin.Context) {

	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("user name is : ", req.Username)
	fmt.Println("password is : ", req.Password)

	UserExists, err := IsUserExists(req.Username)
	if err != nil {
		fmt.Println("error is : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	} else if UserExists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("User already exists with username: %s", req.Username),
		})
		return
	}

	fmt.Println("Everything is fine will hash the password");

	hash, err := HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error" : "Error while hashing the password",
		})
		return;
	}

	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO users (username, password) VALUES ($1, $2)";

	_, err = config.DB.Exec(ct, query, req.Username, hash);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert user into database",
		});
		return;
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User has been registered successfully",
	})
}

// helper functions
func IsUserExists(username string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"
	var exists bool

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := config.DB.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error while Hashing the password, error is : ", err)
		return "", err
	}

	return string(hashed), nil
}
