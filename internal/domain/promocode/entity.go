package promocode

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Target struct {
	AgeFrom    int
	AgeUntil   int
	Country    string
	Categories []string
}

type PromoCode struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description string    `gorm:"type:text"`

	CompanyID   uuid.UUID `gorm:"type:uuid"`
	CompanyName string    `gorm:"type:varchar(255)"`

	MaxCount  int `gorm:"column:max_count"`
	UsedCount int `gorm:"column:used_count"`

	TargetAgeFrom         *int            `gorm:"column:target_age_from"`
	TargetAgeUntil        *int            `gorm:"column:target_age_until"`
	TargetCountry         *string         `gorm:"column:target_country"`
	TargetCountryLower    *string         `gorm:"column:target_country_lower"`
	TargetCategories      *pq.StringArray `gorm:"type:text[];column:target_categories"`
	TargetCategoriesLower *pq.StringArray `gorm:"type:text[];column:target_categories_lower"`

	Mode Mode `gorm:"column:mode"`

	Promo          pq.StringArray `gorm:"type:text[]"`
	AvailablePromo pq.StringArray `gorm:"type:text[]"`

	ImageURL    *string    `gorm:"type:text"`
	ActiveFrom  *time.Time `gorm:"column:active_from;type:date"`
	ActiveUntil *time.Time `gorm:"column:active_until;type:date"`
}

type Like struct {
	PromoCodeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time
}

type Comment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PromoCodeID uuid.UUID `gorm:"type:uuid"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	Content     string    `gorm:"type:text"`
	CreatedAt   time.Time
}

type CommentView struct {
	Id     uuid.UUID `json:"id"`
	Text   string    `json:"text"`
	Date   string    `json:"date"`
	Author struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Surname   string    `json:"surname"`
		AvatarUrl *string   `json:"avatar_url"`
	} `json:"author"`
}

type Use struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	PromoCodeID  uuid.UUID `gorm:"type:uuid"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	Country      string    `gorm:"type:varchar(255)"`
	CountryLower string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
}
