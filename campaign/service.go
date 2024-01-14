package campaign

import "errors"

type Usecase interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
}

type usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{repository}
}

func (u *usecase) GetCampaigns(userID int) ([]Campaign, error) {
	if userID > 0 {
		campaign, err := u.repository.FindByUserId(userID)
		if err != nil {
			return nil, err
		}

		return campaign, err
	}

	campaigns, err := u.repository.FindAll()
	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (u *usecase) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := u.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		return campaign, errors.New("no user found on that id")
	}

	return campaign, nil
}
