package main

import (
	"bwastartup/auth"
	"bwastartup/delivery"
	"bwastartup/helper"
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

	router := gin.Default()

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "route not found",
		})
	})

	v1 := router.Group("/api/v1")
	v1.POST("/users", d.Register)
	v1.POST("/sessions", d.Login)
	v1.POST("/email_checkers", d.CheckEmailIsAvailable)
	v1.POST("/avatars", authMiddleware(authService, u), d.UploadAvatar)

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
