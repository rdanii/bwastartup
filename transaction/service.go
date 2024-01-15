package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type Usecase interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
}

type usecase struct {
	repository   Repository
	repoCampaign campaign.Repository
}

func NewUsecase(repository Repository, repoCampaign campaign.Repository) Usecase {
	return &usecase{repository, repoCampaign}
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
