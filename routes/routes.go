package routes

import(
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router * gin.Engine) {
	UserRoutes(router)
	TodoRoutes(router)
}