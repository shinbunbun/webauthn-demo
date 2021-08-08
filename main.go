package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")
	router.Static("/script", "./script")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/register", func(ctx *gin.Context) {
		var json JsonRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"email": json.Email, "display_name": json.DisplayName})
	})

	router.Run(":8080")
}
