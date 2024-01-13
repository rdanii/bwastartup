package delivery

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Delivery interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	CheckEmailIsAvailable(c *gin.Context)
	UploadAvatar(c *gin.Context)
}

type delivery struct {
	usecase     user.Usecase
	authService auth.Service
}

func NewDelivery(usecase user.Usecase, authService auth.Service) Delivery {
	return &delivery{usecase, authService}
}

func (d *delivery) Register(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMessage := gin.H{"errors": errors}
		response := helper.ResponseContainer("register account failed", http.StatusBadRequest, "error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, errResult := d.usecase.RegisterInput(input)
	if errResult != nil {
		if strings.Contains(errResult.Error(), "email already exists") {
			response := helper.ResponseContainer("the email you entered is already registered", http.StatusNotFound, "error", nil)
			c.JSON(http.StatusNotFound, response)
			return
		}
		response := helper.ResponseContainer("register account failed", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, errToken := d.authService.GenerateToken(result.ID)
	if errToken != nil {
		response := helper.ResponseContainer("register account failed", http.StatusBadRequest, "error", errToken.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(result, token)

	response := helper.ResponseContainer("account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (d *delivery) Login(c *gin.Context) {
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMessage := gin.H{"errors": errors}
		response := helper.ResponseContainer("login failed", http.StatusBadRequest, "error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, errResult := d.usecase.Login(input)
	if errResult != nil {
		if strings.Contains(errResult.Error(), "user not found") {
			response := helper.ResponseContainer("no user found on that email", http.StatusNotFound, "error", nil)
			c.JSON(http.StatusNotFound, response)
			return
		}
		response := helper.ResponseContainer("login failed", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, errToken := d.authService.GenerateToken(result.ID)
	if errToken != nil {
		response := helper.ResponseContainer("register account failed", http.StatusBadRequest, "error", errToken.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(result, token)
	response := helper.ResponseContainer("login success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (d *delivery) CheckEmailIsAvailable(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("email checking failed", http.StatusUnprocessableEntity, "error", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isAvailable, errCheckEmail := d.usecase.IsAvailableEmail(input)
	if errCheckEmail != nil {
		response := helper.ResponseContainer("email checking failed", http.StatusBadRequest, "error", errCheckEmail.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"isAvailable": isAvailable,
	}

	msg := "email has been registered"

	if isAvailable {
		msg = "email is available"
	}

	response := helper.ResponseContainer(msg, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (d *delivery) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User).ID

	path := fmt.Sprintf("images/%d-%s", currentUser, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = d.usecase.SaveAvatar(currentUser, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ResponseContainer("Avatar successfuly uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
