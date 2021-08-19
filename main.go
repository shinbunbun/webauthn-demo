package main

import (
	"net/http"
	"time"

	"crypto/rand"
	"math/big"

	"encoding/base64"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type JsonRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type RegisterJson struct {
	ID       string `json:"id"`
	Response struct {
		ClientDataJSON struct {
			Type        string `json:"type"`
			Challenge   string `json:"challenge"`
			Origin      string `json:"origin"`
			CrossOrigin bool   `json:"crossOrigin"`
		} `json:"clientDataJSON"`
		AttestationObject struct {
			Fmt     string `json:"fmt"`
			AttStmt struct {
				Alg int    `json:"alg"`
				Sig string `json:"sig"`
			} `json:"attStmt"`
			AuthData string `json:"authData"`
		} `json:"attestationObject"`
	} `json:"response"`
}

type LoginRequestParam struct {
}

type Session struct {
	Id          string `gorm:"column:id"`
	Email       string `gorm:"column:email"`
	Challenge   string `gorm:"column:challenge"`
	DisplayName string `gorm:"column:displayname"`
}

type User struct {
	Id          string `gorm:"column:id"`
	Email       string `gorm:"column:email"`
	DisplayName string `gorm:"column:displayname"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")
	router.Static("/script", "./script")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/success-sign-in", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "success-sign-in.html", gin.H{})
	})

	router.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost:8080",
		},
		// アクセスを許可したいHTTPメソッド(以下の例だとPUTやDELETEはアクセスできません)
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: false,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	router.POST("/register-request", func(ctx *gin.Context) {
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
		var id string = base64.StdEncoding.EncodeToString([]byte(uuidObj.String()))
		var challengeStr string = base64.StdEncoding.EncodeToString([]byte(challenge.String()))
		var sessionData = Session{Id: id, Email: json.Email, Challenge: challengeStr, DisplayName: json.DisplayName}
		db.Create(&sessionData)
		ctx.JSON(http.StatusOK, gin.H{"id": id, "challenge": challengeStr, "rp": "bunbun-test-rp"})
	})

	router.POST("/register", func(ctx *gin.Context) {
		var json RegisterJson
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db, dbErr := gorm.Open("sqlite3", "store.sqlite")
		if dbErr != nil {
			println(dbErr)
		}

		var session Session

		db.Where("challenge = ?", json.Response.ClientDataJSON.Challenge).First(&session)

		var userData = User{Id: json.ID, Email: session.Email, DisplayName: session.DisplayName}
		db.Create(&userData)

		ctx.JSON(http.StatusOK, gin.H{"verificationStatus": "succeeded"})
	})

	router.POST("/login-request", func(ctx *gin.Context) {

		var json JsonRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db, dbErr := gorm.Open("sqlite3", "store.sqlite")
		if dbErr != nil {
			println(dbErr)
		}

		var userData User
		db.Where("Email = ?", json.Email).First(&userData)

		challenge, challengeErr := rand.Int(rand.Reader, big.NewInt(999999999999999999))
		if challengeErr != nil {
			println(challengeErr)
		}

		var challengeStr string = base64.StdEncoding.EncodeToString([]byte(challenge.String()))

		ctx.JSON(http.StatusOK, gin.H{"challenge": challengeStr, "id": userData.Id})
	})

	router.Run(":8080")
}
