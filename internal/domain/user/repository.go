package user

import (
	"github.com/google/uuid"
	"solution/internal/domain/business"
	customerrors "solution/internal/domain/errors"
	"time"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Sub       uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	RevokedAt *time.Time
}

type TokenManager interface {
	GenerateToken(businessID uuid.UUID, email string) string
	ValidateToken(tokenString string) (*business.TokenClaims, *customerrors.TokenError)
	RevokeToken(tokenString string)
}

type Repository interface {
	Create(u *User) *customerrors.RepositoryError
	GetByEmail(email string) (*User, *customerrors.RepositoryError)
	Get(id uuid.UUID) (*User, *customerrors.RepositoryError)
	Save(u *User)
}
