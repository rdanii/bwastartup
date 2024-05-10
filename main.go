package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/delivery"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/bwa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	authService := auth.NewService()
	paymentService := payment.NewService()

	r := user.NewRepository(db)
	u := user.NewUsecase(r)
	d := delivery.NewDelivery(u, authService)

	campaignRepo := campaign.NewRepository(db)
	campaignUsecase := campaign.NewUsecase(campaignRepo)
	campaignDelivery := delivery.NewCampaignDelivery(campaignUsecase)

	transactionRepo := transaction.NewRepository(db)
	transactionUsecase := transaction.NewUsecase(transactionRepo, campaignRepo, paymentService)
	transactionDelivery := delivery.NewTransactionDelivery(transactionUsecase)

	router := gin.Default()

	// Menggunakan middleware cors.Default() untuk menangani CORS.
	router.Use(cors.Default())

	// Membuat instance baru dari cookie store dengan menggunakan SECRET_KEY untuk enkripsi.
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))

	// Menggunakan middleware sessions.Sessions() untuk manajemen sesi.
	router.Use(sessions.Sessions("bwastartup", cookieStore))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "route not found",
		})
	})

	api := router.Group("/api/v1")

	router.Static("/images", "./images")

	// user API
	api.POST("/users", d.Register)
	api.POST("/sessions", d.Login)
	api.POST("/email_checkers", d.CheckEmailIsAvailable)
	api.POST("/avatars", authMiddleware(authService, u), d.UploadAvatar)

	// campaigns API
	api.GET("/campaigns", campaignDelivery.GetCampaigns)
	api.GET("/campaigns/:id", campaignDelivery.GetCampaignDetail)
	api.POST("/campaigns", authMiddleware(authService, u), campaignDelivery.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, u), campaignDelivery.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, u), campaignDelivery.UploadCampaignImage)
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, u), transactionDelivery.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, u), transactionDelivery.GetTransactionsByUserID)
	api.POST("/transactions", authMiddleware(authService, u), transactionDelivery.CreateTransaction)

	router.Run()
}

// Middleware Gin yang digunakan untuk melakukan otentikasi
func authMiddleware(authService auth.Service, userService user.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengambil nilai Authorization header dari request.
		authHeader := c.GetHeader("Authorization")

		// Memeriksa apakah Authorization header mengandung token JWT dengan skema "Bearer".
		if !strings.Contains(authHeader, "Bearer") {
			// Jika tidak, kembalikan respons Unauthorized.
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Mendapatkan token JWT dari Authorization header.
		var tokenString string
		strToken := strings.Split(authHeader, " ")
		if len(strToken) == 2 {
			tokenString = strToken[1]
		}

		// Memvalidasi token JWT menggunakan AuthService.
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Jika token tidak valid, kembalikan respons Unauthorized.
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Mendapatkan klaim (claims) dari token JWT.
		claim, ok := token.Claims.(jwt.MapClaims)

		// Memeriksa apakah klaim valid.
		if !ok || !token.Valid {
			// Jika klaim tidak valid, kembalikan respons Unauthorized.
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Mendapatkan user ID dari klaim.
		userID := int(claim["user_id"].(float64))

		// Mendapatkan informasi user berdasarkan user ID menggunakan UserService.
		user, err := userService.GetUserByID(userID)
		if err != nil {
			// Jika tidak dapat menemukan user, kembalikan respons Unauthorized.
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Menetapkan informasi user ke dalam konteks (context) Gin.
		c.Set("currentUser", user)
	}
}
