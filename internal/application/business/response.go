package business

import (
	"github.com/google/uuid"
	"solution/internal/domain/promocode"
	"solution/internal/domain/types"
)

type CreateBusinessResponse struct {
	ID    uuid.UUID `json:"id"`
	Token string    `json:"token"`
}

type LoginBusinessResponse struct {
	Token string `json:"token"`
}

type CreatePromoCodeResponse struct {
	ID uuid.UUID `json:"id"`
}

type EditPromoCodeResponse struct {
	Active      bool           `json:"active"`
	CompanyID   uuid.UUID      `json:"company_id"`
	CompanyName string         `json:"company_name"`
	Description string         `json:"description"`
	LikeCount   int            `json:"like_count"`
	MaxCount    int            `json:"max_count"`
	Mode        promocode.Mode `json:"mode"`
	PromoID     uuid.UUID      `json:"promo_id"`
	Target      struct {
		AgeFrom    int      `json:"age_from"`
		AgeUntil   int      `json:"age_until"`
		Country    string   `json:"country"`
		Categories []string `json:"categories"`
	}
	UsedCount   int                `json:"used_count"`
	PromoCommon *string            `json:"promo_common"`
	PromoUnique *[]string          `json:"promo_unique"`
	ImageURL    string             `json:"image_url"`
	ActiveFrom  types.SolutionDate `json:"active_from"`
	ActiveUntil types.SolutionDate `json:"active_until"`
}
