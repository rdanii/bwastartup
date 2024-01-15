package campaign

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserId(userID int) ([]Campaign, error)
	FindByID(id int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	CreateImage(createImage CampaignImage) (CampaignImage, error)
	IsPrimaryImage(id int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaign []Campaign
	err := r.db.Debug().Preload("CampaignImages", "campaign_images.is_primary=1").Find(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) FindByUserId(userID int) ([]Campaign, error) {
	var campaign []Campaign
	err := r.db.Debug().Preload("CampaignImages", "campaign_images.is_primary=1").Find(&campaign, "user_id = ?", userID).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) FindByID(id int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Debug().Preload("CampaignImages").Preload("User").Find(&campaign, "id = ?", id).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Debug().Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Debug().Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) CreateImage(createImage CampaignImage) (CampaignImage, error) {
	err := r.db.Debug().Create(&createImage).Error
	if err != nil {
		return createImage, err
	}

	return createImage, nil
}

func (r *repository) IsPrimaryImage(id int) (bool, error) {
	err := r.db.Debug().Model(&CampaignImage{}).Where("campaign_id = ?", id).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
