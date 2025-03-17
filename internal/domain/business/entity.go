package business

import (
	"github.com/google/uuid"
	"github.com/intezya/pkglib"
)

type Business struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CompanyName string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);not null;unique"`

	PasswordHash string `gorm:"type:TEXT;not null"`
}

func (*Business) TableName() string {
	return "businesses"
}

func (b *Business) Token(secret string) string {
	return pkglib.JWT.GenerateToken(
		map[string]interface{}{
			"sub": b.ID.String(),
		}, secret,
	)
}
