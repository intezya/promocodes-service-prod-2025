package business

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	customerrors "solution/internal/domain/errors"
	"time"
)

type TokenClaims struct {
	Sub       uuid.UUID `json:"business_id"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Sub       uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	RevokedAt *time.Time
}

type TokenManager interface {
	GenerateToken(businessID uuid.UUID, email string) string
	ValidateToken(tokenString string) (*TokenClaims, *customerrors.TokenError)
	RevokeToken(tokenString string)
}

type Repository interface {
	Create(b *Business) *customerrors.RepositoryError
	GetByEmail(email string) (*Business, *customerrors.RepositoryError)
	Get(id uuid.UUID) (*Business, *customerrors.RepositoryError)
}
