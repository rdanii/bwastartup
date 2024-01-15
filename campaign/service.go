package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Usecase interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignIDInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignIDInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
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

func (u *usecase) GetCampaignByID(input GetCampaignIDInput) (Campaign, error) {
	campaign, err := u.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		return campaign, errors.New("no user found on that id")
	}

	return campaign, nil
}

func (u *usecase) CreateCampaign(input CreateCampaignInput) (Campaign, error) {

	campaign := Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Perks:            input.Perks,
		GoalAmount:       input.GoalAmount,
		User:             input.User,
	}

	slugName := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugName)

	newCampaign, err := u.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (u *usecase) UpdateCampaign(inputID GetCampaignIDInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := u.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	dataUpdated, errUpdated := u.repository.Update(campaign)
	if errUpdated != nil {
		return dataUpdated, errUpdated
	}

	return dataUpdated, nil
}

func (u *usecase) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {

	campaign, err := u.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("not an owner of the campaign image")
	}

	var isPrimary int

	if input.IsPrimary {
		isPrimary = 1
		_, errIsPrimary := u.repository.IsPrimaryImage(input.CampaignID)
		if errIsPrimary != nil {
			return CampaignImage{}, errIsPrimary
		}
	}

	campaignImage := CampaignImage{
		CampaignID: input.CampaignID,
		IsPrimary:  isPrimary,
		FileName:   fileLocation,
	}

	data, errData := u.repository.CreateImage(campaignImage)
	if errData != nil {
		return data, errData
	}

	return data, nil
}
