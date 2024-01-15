package delivery

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignDelivery interface {
	GetCampaigns(c *gin.Context)
	GetCampaignDetail(c *gin.Context)
	CreateCampaign(c *gin.Context)
	UpdateCampaign(c *gin.Context)
	UploadCampaignImage(c *gin.Context)
}

type campaignDelivery struct {
	campaignUsecase campaign.Usecase
}

func NewCampaignDelivery(campaignUsecase campaign.Usecase) CampaignDelivery {
	return &campaignDelivery{campaignUsecase}
}

func (d *campaignDelivery) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	result, errResult := d.campaignUsecase.GetCampaigns(userID)
	if errResult != nil {
		response := helper.ResponseContainer("error to get campaigns", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaigns(result)
	response := helper.ResponseContainer("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (d *campaignDelivery) GetCampaignDetail(c *gin.Context) {
	var input campaign.GetCampaignIDInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ResponseContainer("failed to get detail of campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, errResult := d.campaignUsecase.GetCampaignByID(input)
	if errResult != nil {
		response := helper.ResponseContainer("failed to get detail of campaign", http.StatusNotFound, "error", errResult.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	formatter := campaign.FormatDetailCampaign(result)
	response := helper.ResponseContainer("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (d *campaignDelivery) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("failed to create campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	result, errResult := d.campaignUsecase.CreateCampaign(input)
	if errResult != nil {
		response := helper.ResponseContainer("failed to create campaign", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaign(result)
	response := helper.ResponseContainer("success to create campaign", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (d *campaignDelivery) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignIDInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("failed to update campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("failed to update campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	result, errResult := d.campaignUsecase.UpdateCampaign(inputID, inputData)
	if errResult != nil {
		response := helper.ResponseContainer("failed to update campaign", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaign(result)
	response := helper.ResponseContainer("success to update campaign", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (d *campaignDelivery) UploadCampaignImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("failed to upload campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, errFile := c.FormFile("file")
	if errFile != nil {
		errMsg := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("failed to upload campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errMsg := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("failed to upload campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, errResult := d.campaignUsecase.SaveCampaignImage(input, path)
	if errResult != nil {
		errMsg := gin.H{"is_uploaded": false}
		response := helper.ResponseContainer("failed to upload campaign", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	errMsg := gin.H{"is_uploaded": true}
	response := helper.ResponseContainer("success to update campaign", http.StatusOK, "success", errMsg)
	c.JSON(http.StatusOK, response)
}
