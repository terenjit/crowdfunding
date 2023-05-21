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
	UserID     string      `json:"user_id"`
	Opts       token.Claim `json:"opts"`
}
