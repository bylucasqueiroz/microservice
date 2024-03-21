package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.GetString("Hello, World!")
	})
	router.Run("localhost:5050")
}
