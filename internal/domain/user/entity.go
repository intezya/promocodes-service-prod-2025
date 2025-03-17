package user

import (
	"github.com/google/uuid"
	"github.com/intezya/pkglib"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name    string `gorm:"type:varchar(255);not null"`
	Surname string `gorm:"type:varchar(255);not null"`

	Email string `gorm:"type:varchar(255);not null;unique"`

	Age     int    `gorm:"column:age"`
	Country string `gorm:"type:varchar(255);column:country"`

	PasswordHash string `gorm:"type:TEXT;not null"`

	AvatarURL *string `gorm:"type:TEXT"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) Token(secret string) string {
	return pkglib.JWT.GenerateToken(
		map[string]interface{}{
			"sub": u.ID.String(),
		}, secret,
	)
}

type Profile struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Other   struct {
		Age     int    `json:"age"`
		Country string `json:"country"`
	} `json:"other"`
	AvatarUrl *string `json:"avatar_url"`
}
