package routes

import(
	"github.com/gin-gonic/gin"
	"github.com/manishmaang/TODO-APPLICATION/controllers"

)

func UserRoutes(router * gin.Engine){
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":"Hello World !!",
		})
	});

	router.POST("/register", controllers.RegisterUsers);
}