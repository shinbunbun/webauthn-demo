package main

import (
	"net/http"

	"crypto/rand"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type JsonRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type Session struct {
	Id        string `gorm:"column:id"`
	Email     string `gorm:"column:email"`
	Challenge string `gorm:"column:challenge"`
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

		db, dbErr := gorm.Open("sqlite3", "store.sqlite")
		if dbErr != nil {
			println(dbErr)
		}

		challenge, challengeErr := rand.Int(rand.Reader, big.NewInt(999999999999999999))
		if challengeErr != nil {
			println(challengeErr)
		}

		uuidObj, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			println(challengeErr)
		}

		var sessionData = Session{Id: uuidObj.String(), Email: json.Email, Challenge: challenge.String()}
		db.Create(&sessionData)
		ctx.JSON(http.StatusOK, gin.H{"id": sessionData.Id, "challenge": sessionData.Challenge})
	})

	router.Run(":8080")
}
