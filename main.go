package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manishmaang/TODO-APPLICATION/routes" 
	_ "github.com/manishmaang/TODO-APPLICATION/config" // run init() only
)

func main() {
	router := gin.Default()
	// Hello world program !!
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "Hello world !!",
	// 	})
	// });

	routes.SetupRoutes(router) // Register all routesroutes.SetupRoutes(router);

	router.Run(":3000")
}
