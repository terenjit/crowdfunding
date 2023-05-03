package domain

import (
	"crowdfunding/pkg/token"
	"time"
)

type Campaign struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description" `
	Perks            string    `json:"perks"`
	BackerCount      int       `json:"backer_count"`
	GoalAmount       int       `json:"goal_amount"`
	CurrentAmount    int       `json:"current_amount"`
	Slug             int       `json:"slug"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CampaignImage struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	FileName   string    `json:"name"`
	IsPrimary  int       `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CampainGetList struct {
	Status   string      `json:"status"`
	Search   interface{} `json:"search"`
	Page     int         `json:"page"`
	Quantity int         `json:"quantity"`
	Query    string      `query:"query"`
	SortBy   string      `json:"sort_by"`
	UserID   string      `json:"user_id"`
	Opts     token.Claim `json:"opts"`
}
