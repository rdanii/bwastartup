package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByCampaignID(id int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(id int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Debug().Preload("User").Find(&transaction, "campaign_id = ?", id).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Debug().Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Find(&transactions, "user_id = ?", userID).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
