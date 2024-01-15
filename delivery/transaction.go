package delivery

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionDelivery interface {
	GetCampaignTransactions(c *gin.Context)
	GetTransactionsByUserID(c *gin.Context)
}

type transactionDelivery struct {
	usecase transaction.Usecase
}

func NewTransactionDelivery(usecase transaction.Usecase) TransactionDelivery {
	return &transactionDelivery{usecase}
}

func (d *transactionDelivery) GetCampaignTransactions(c *gin.Context) {

	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMsg := gin.H{"errors": errors}
		response := helper.ResponseContainer("Failed to get campaign's transactions", http.StatusBadRequest, "error", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	result, errResult := d.usecase.GetTransactionsByCampaignID(input)
	if errResult != nil {
		response := helper.ResponseContainer("Failed to get campaign's transactions", http.StatusBadRequest, "error", errResult.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatCampaignTransactions(result)
	response := helper.ResponseContainer("Campaign's transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (d *transactionDelivery) GetTransactionsByUserID(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	result, err := d.usecase.GetTransactionsByUserID(currentUser.ID)
	if err != nil {
		response := helper.ResponseContainer("Failed to get users's transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatUserTransactions(result)
	response := helper.ResponseContainer("Users's transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
