package delivery

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignDelivery interface {
	GetCampaigns(c *gin.Context)
	GetCampaignDetail(c *gin.Context)
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
	var input campaign.GetCampaignDetailInput

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

	formatter := campaign.FormatCampaignDetail(result)
	response := helper.ResponseContainer("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
