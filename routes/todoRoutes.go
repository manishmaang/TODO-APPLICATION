package routes

import(
	"github.com/gin-gonic/gin"
)

func TodoRoutes(router * gin.Engine){
	router.GET("/hello", func (ctx * gin.Context)  {
		ctx.JSON(200, gin.H{
			"message" : "Hello from todo routes",
		})
	})
}