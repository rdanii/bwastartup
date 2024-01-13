package main

import (
	"bwastartup/delivery"
	"bwastartup/user"
	"log"
	"net/http"

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

	r := user.NewRepository(db)
	u := user.NewUsecase(r)
	d := delivery.NewDelivery(u)

	router := gin.Default()

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "route not found",
		})
	})

	v1 := router.Group("/api/v1")
	v1.POST("/users", d.Register)

	router.Run()

}
