package domain

import (
	"crowdfunding/pkg/token"
	"time"
)

type Transaction struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	FileName   string    `json:"file_name"`
	IsPrimary  int       `json:"is_primary"`
	PaymentURL string    `json:"payment_url"`
	Amount     int64     `json:"amount"`
	Status     string    `json:"status"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TransactionList struct {
	Status     string      `json:"status"`
	Search     interface{} `json:"search"`
	Page       int         `json:"page"`
	Quantity   int         `json:"quantity"`
	Query      string      `query:"query"`
	SortBy     string      `json:"sort_by"`
	CampaignID string      `json:"campaign_id" param:"campaign_id"`
	IsPrimary  int         `json:"is_primary"`
	UserID     string      `json:"user_id"`
	Opts       token.Claim `json:"opts"`
}

type CreateRequest struct {
	ID         string      `json:"id"`
	CampaignID string      `json:"campaign_id"`
	Amount     int64       `json:"amount"`
	Opts       token.Claim `json:"opts,omitempty" gorm:"-"`
}

type TransactionModel struct {
	ID         string      `json:"id,omitempty"`
	CampaignID string      `json:"campaign_id,omitempty"`
	UserID     string      `json:"user_id,omitempty"`
	Amount     int64       `json:"amount,omitempty"`
	Status     string      `json:"status,omitempty"`
	Code       string      `json:"code,omitempty"`
	PaymentURL string      `json:"payment_url"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Opts       token.Claim `json:"opts" gorm:"-"`
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

type CampaignModel struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description" `
	Perks            string    `json:"perks"`
	BackerCount      int64     `json:"backer_count"`
	GoalAmount       int64     `json:"goal_amount"`
	CurrentAmount    int64     `json:"current_amount"`
	Slug             string    `json:"slug"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
