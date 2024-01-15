package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {

	transactionsformatter := []CampaignTransactionFormatter{}

	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionsformatter = append(transactionsformatter, formatter)
	}

	return transactionsformatter
}

type UserTransactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign: CampaignFormatter{
			Name: transaction.Campaign.Name,
		},
	}

	if len(transaction.Campaign.CampaignImages) > 0 {
		formatter.Campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	return formatter
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	transactionsFormatter := []UserTransactionFormatter{}

	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}
