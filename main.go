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

// -----------------------------------------------------------
//  Context Summary Table (Go/Gin)
// -----------------------------------------------------------
// context.Background() 
// → base context, never times out or cancels—used as a starting point
//
// context.WithTimeout(ctx, duration) 
// → wraps a context with a deadline—automatically cancels after given time
//
// cancel() 
// → function returned by WithTimeout that manually stops the context & releases resources
//
// defer cancel() 
// → schedules cancel() to run after current function finishes,
//    ensures cleanup happens even on early return or error
//
// Best Practice: always use 'defer cancel()' right after creating
// a timeout context to prevent resource leaks and dangling operations
// -----------------------------------------------------------
