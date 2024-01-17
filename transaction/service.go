package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
)

type Usecase interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

type usecase struct {
	repository     Repository
	repoCampaign   campaign.Repository
	paymentService payment.Service
}

func NewUsecase(repository Repository, repoCampaign campaign.Repository, paymentService payment.Service) Usecase {
	return &usecase{repository, repoCampaign, paymentService}
}

func (u *usecase) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaign, err := u.repoCampaign.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not an owner of the campaigns")
	}

	transactions, errTransaction := u.repository.GetByCampaignID(input.ID)
	if errTransaction != nil {
		return transactions, errTransaction
	}

	return transactions, nil
}

func (u *usecase) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := u.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (u *usecase) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{
		CampaignID: input.CampaignID,
		UserID:     input.User.ID,
		Amount:     input.Amount,
		Status:     "pending",
	}

	createTransaction, err := u.repository.Save(transaction)
	if err != nil {
		return createTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     createTransaction.ID,
		Amount: createTransaction.Amount,
	}

	paymentURL, errPyamentURL := u.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if errPyamentURL != nil {
		return createTransaction, errPyamentURL
	}

	createTransaction.PaymentURL = paymentURL

	createTransaction, err = u.repository.Update(createTransaction)
	if err != nil {
		return createTransaction, err
	}

	return createTransaction, nil
}
