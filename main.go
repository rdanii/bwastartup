package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/delivery"
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	r := user.NewRepository(db)
	u := user.NewUsecase(r)
	d := delivery.NewDelivery(u, authService)

	campaignRepo := campaign.NewRepository(db)
	campaignUsecase := campaign.NewUsecase(campaignRepo)
	campaignDelivery := delivery.NewCampaignDelivery(campaignUsecase)

	transactionRepo := transaction.NewRepository(db)
	transactionUsecase := transaction.NewUsecase(transactionRepo, campaignRepo)
	transactionDelivery := delivery.NewTransactionDelivery(transactionUsecase)

	router := gin.Default()

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

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		strToken := strings.Split(authHeader, " ")
		if len(strToken) == 2 {
			tokenString = strToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.ResponseContainer("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
