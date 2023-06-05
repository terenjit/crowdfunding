package domain

import (
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/token"
	"time"
)

type Campaign struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	Name             string           `json:"name"`
	ShortDescription string           `json:"short_description"`
	Description      string           `json:"description" `
	User             models.User      `json:"user,omitempty"`
	Perks            string           `json:"perks"`
	BackerCount      int64            `json:"backer_count"`
	GoalAmount       int64            `json:"goal_amount"`
	CurrentAmount    int              `json:"current_amount"`
	Slug             string           `json:"slug"`
	Images           []CampaignImages `json:"images" gorm:"-"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type CampaignsFormat struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description" `
	Perks            string    `json:"perks"`
	BackerCount      int64     `json:"backer_count"`
	GoalAmount       int64     `json:"goal_amount"`
	CurrentAmount    int       `json:"current_amount"`
	Slug             string    `json:"slug"`
	ImagesURL        string    `json:"images_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CampaignFormatter struct {
	ID               string                    `json:"id"`
	UserID           string                    `json:"user_id"`
	Name             string                    `json:"name"`
	ShortDescription string                    `json:"short_description"`
	Description      string                    `json:"description" `
	ImagesURL        string                    `json:"images_url"`
	User             UserFormat                `json:"user"`
	Perks            []string                  `json:"perks"`
	BackerCount      int64                     `json:"backer_count"`
	GoalAmount       int64                     `json:"goal_amount"`
	CurrentAmount    int64                     `json:"current_amount"`
	Slug             string                    `json:"slug"`
	Images           []CampaignImagesFormatter `json:"images" gorm:"-"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

type CampaignImages struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	FileName   string    `json:"name"`
	IsPrimary  int       `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UploadImageRequest struct {
	CampaignID string      `json:"campaign_id" form:"campaign_id"`
	IsPrimary  bool        `json:"is_primary" form:"is_primary"`
	Opts       token.Claim `json:"opts,omitempty"`
}

type CampaignImagesFormatter struct {
	FileName  string `json:"name"`
	IsPrimary bool   `json:"is_primary"`
}

type CampainGetList struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"`
	Search   string      `json:"search"`
	Page     int         `json:"page"`
	Quantity int         `json:"quantity"`
	Query    string      `query:"query"`
	SortBy   string      `json:"sort_by"`
	UserID   string      `json:"user_id"`
	Opts     token.Claim `json:"opts"`
}

type CampaignGetDetail struct {
	ID string `json:"id"`
}

type UserFormat struct {
	Name           string `json:"name"`
	AvatarFileName string `json:"avatar"`
}

type UpdateCampaign struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	Name             string           `json:"name"`
	ShortDescription string           `json:"short_description"`
	Description      string           `json:"description" `
	Perks            string           `json:"perks"`
	BackerCount      int64            `json:"backer_count"`
	GoalAmount       int64            `json:"goal_amount"`
	CurrentAmount    int64            `json:"current_amount"`
	Slug             string           `json:"slug"`
	Images           []CampaignImages `json:"images" gorm:"-"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	UpdatedBy        string           `json:"update_by"`
	Opts             token.Claim      `json:"opts"`
}

type CreateRequest struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	ShortDescription string      `json:"short_description"`
	Description      string      `json:"description"`
	GoalAmount       int64       `json:"goal_amount"`
	Slug             string      `json:"slug"`
	Perks            string      `json:"perks"`
	Opts             token.Claim `json:"opts,omitempty"`
}
