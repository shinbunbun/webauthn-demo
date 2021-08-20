package main

import (
	"fmt"
	"net/http"
	"time"

	"crypto/rand"
	"math/big"

	"encoding/base64"

	"crypto/sha256"

	"encoding/binary"

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
	ClientDataJSONString string `json:"clientDataJSONString"`
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

type AuthenticatorData struct {
	rpIdHash []byte
	flags    struct {
		ed byte
		at byte
		uv byte
		up byte
	}
	signCount              []byte
	attestedCredentialData struct {
		aaguid              []byte
		credentialIdLength  uint16
		credentialId        []byte
		credentialPublicKey []byte
	}
}

func getSHA256Binary(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

func Bytes2uint(bytes ...byte) uint16 {
	padding := make([]byte, 8-len(bytes))
	i := binary.BigEndian.Uint16(append(padding, bytes...))
	return i
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
		var registerJson RegisterJson
		if err := ctx.ShouldBindJSON(&registerJson); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		/* fmt.Println(registerJson.ClientDataJSONString)
		fmt.Println("まああああああああああにいいいいいいいいいいいいいい！！！！！！！！") */

		db, dbErr := gorm.Open("sqlite3", "store.sqlite")
		if dbErr != nil {
			println(dbErr)
		}

		var isValid int = 1

		// typeを検証
		if registerJson.Response.ClientDataJSON.Type != "webauthn.create" {
			isValid = 0
		}

		//originを検証
		if registerJson.Response.ClientDataJSON.Origin != "http://localhost:8080" {
			isValid = 0
		}

		// Authenticator Dataを取り出す
		authData, _ := base64.StdEncoding.DecodeString(registerJson.Response.AttestationObject.AuthData)
		var authenticatorData AuthenticatorData
		authenticatorData.rpIdHash = authData[0:32]
		var flags = authData[32]
		authenticatorData.flags.ed = (flags >> 7) & 1
		authenticatorData.flags.at = (flags >> 6) & 1
		authenticatorData.flags.uv = (flags >> 2) & 1
		authenticatorData.flags.up = flags & 1
		authenticatorData.signCount = authData[33:37]
		authenticatorData.attestedCredentialData.aaguid = authData[37:53]
		authenticatorData.attestedCredentialData.credentialIdLength = Bytes2uint(authData[53:56]...)
		authenticatorData.attestedCredentialData.credentialId = authData[53 : 53+authenticatorData.attestedCredentialData.credentialIdLength+1]
		authenticatorData.attestedCredentialData.credentialPublicKey = authData[53+authenticatorData.attestedCredentialData.credentialIdLength+1 : len(authData)]

		/* fmt.Printf("%#v\n", authenticatorData) */

		// rpidのhashを検証
		/* var clientDataJSON, clientDataJSONErr = json.Marshal(registerJson.Response.ClientDataJSON)
		if clientDataJSONErr != nil {
			println(clientDataJSONErr)
		} */
		fmt.Println(getSHA256Binary("localhost"))
		fmt.Println(authenticatorData.rpIdHash)
		if string(getSHA256Binary("localhost")) != string(authenticatorData.rpIdHash) {
			isValid = 0
		}

		// flagsを検証
		if authenticatorData.flags.uv != 1 && authenticatorData.flags.up != 1 {
			isValid = 0
		}

		/* 		//よくわからんエンコードだから、JSで書いて
		   		//うんちぶりぶり
		   		var credentialPublicKey []byte
		   		err := cbor.Unmarshal(authenticatorData.attestedCredentialData.credentialPublicKey, &credentialPublicKey)
		   		if err != nil {
		   			fmt.Println("ERROR:", err)
		   		}

		   		//map[interface{}]interface{}
		   		// 0xなんとか
		   		fmt.Println(credentialPublicKey) */

		println(isValid)

		var session Session

		db.Where("challenge = ?", registerJson.Response.ClientDataJSON.Challenge).First(&session)

		var userData = User{Id: registerJson.ID, Email: session.Email, DisplayName: session.DisplayName}
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
