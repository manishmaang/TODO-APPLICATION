package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/manishmaang/TODO-APPLICATION/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

// | Function         | Purpose                                                                               |
// | ---------------- | ------------------------------------------------------------------------------------- |
// | **`Exec()`**     | Used for statements that **don’t return rows**, like `INSERT`, `UPDATE`, or `DELETE`. |
// | **`QueryRow()`** | Used when you expect **a single row result** from the database (e.g., `SELECT`).      |
// | **`Query()`**    | Used when you expect **multiple rows** from a `SELECT`.                               |

func LogInUsers(ctx *gin.Context) {
	type LoginInRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req_body LoginInRequest
	if err := ctx.BindJSON(&req_body); err != nil {
		fmt.Println("Error while binding the req body, error is : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// different hash is generated every time even for the same password there use comparer and hash password function
	// hash, err := HashPassword(req_body.Password)
	// if err != nil {
	// 	fmt.Println("Error while hasing the password error is : ", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Internal Server Error",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	query := "SELECT username, password FROM users WHERE username = $1"
	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var username, hashedPassword string
	err := config.DB.QueryRow(ct, query, req_body.Username).Scan(&username, &hashedPassword)
	if err != nil {
		fmt.Println("Error while searching for the doc, error is : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found or DB error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req_body.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	accessToken, err := GenerateToken(req_body.Username, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken, err := GenerateToken(req_body.Username, true) // only refreshToken is new variable err is used from above declaration
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie(
		"refresh_token", // Cookie name
		refreshToken,    // Value
		60*60*24*7,      // MaxAge in seconds (7 days)
		"/",             // Path
		"",              // Domain ("" means current)
		false,           // Secure (use true in production with HTTPS)
		true,            // HTTPOnly (prevents JS access — good!)
	)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func RegisterUsers(ctx *gin.Context) { // ctx represents everything about the HTTP request and response.

	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil { // shouldBinJson copy request body value into our go struct
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 	Go's request body is a raw stream (io.ReadCloser), not a ready-made JSON object.
	// So you can't directly access properties like body.username.
	// You must read the stream and parse it into a struct.

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

	fmt.Println("Everything is fine will hash the password")

	hash, err := HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while hashing the password",
		})
		return
	}

	ct, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"

	_, err = config.DB.Exec(ct, query, req.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert user into database",
		})
		return
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

func GenerateToken(username string, isRefresh bool) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file, error is ", err)
		return "", err
	}

	secret := os.Getenv("JWT_SECRET")
	expiryTime := time.Hour * 24
	if isRefresh {
		secret = os.Getenv("REFRESH_SECRET")
		expiryTime = time.Hour * 24 * 7
	}

	claims := jwt.MapClaims{
		"user":   username,
		"expiry": time.Now().Add(expiryTime).Unix(),
		"type": func() string {
			if isRefresh {
				return "refresh"
			} else {
				return "access"
			}
		}(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	} else {
		return signedToken, nil
	}
}
