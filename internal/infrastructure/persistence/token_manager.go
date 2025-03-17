package persistence

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"solution/internal/domain/business"
	"solution/internal/domain/errors"
	"time"
)

type TokenManagerRepository struct {
	db        *gorm.DB
	secretKey []byte
}

func NewTokenManagerRepository(db *gorm.DB, secretKey []byte) *TokenManagerRepository {
	return &TokenManagerRepository{db: db, secretKey: secretKey}
}

func (r *TokenManagerRepository) GenerateToken(businessID uuid.UUID, email string) string {
	sessionID := uuid.New()
	created := time.Now()
	session := business.Session{
		ID:        sessionID,
		Sub:       businessID,
		CreatedAt: created,
	}

	r.db.Table("sessions").Where("sub = ? AND revoked_at IS NULL", businessID).
		Update("revoked_at", time.Now())
	r.db.Create(&session)

	claims := business.TokenClaims{
		Sub:       businessID,
		Email:     email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(created.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(created),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(r.secretKey)
	return signed
}

func (r *TokenManagerRepository) ValidateToken(tokenString string) (
	*business.TokenClaims,
	*customerrors.TokenError,
) {
	token, err := jwt.ParseWithClaims(
		tokenString, &business.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return r.secretKey, nil
		},
	)

	if err != nil {
		log.Error(err)
		return nil, &customerrors.TokenError{
			Message: "invalid token",
		}
	}

	if claims, ok := token.Claims.(*business.TokenClaims); ok && token.Valid {
		var session business.Session
		if err := r.db.First(&session, claims.SessionID).Error; err != nil {
			return nil, &customerrors.TokenError{Message: "invalid session"}
		}

		if session.RevokedAt != nil {
			return nil, &customerrors.TokenError{Message: "session revoked"}
		}

		return claims, nil
	}
	log.Info(":(")
	return nil, &customerrors.TokenError{Message: "invalid token"}
}

func (r *TokenManagerRepository) RevokeToken(tokenString string) {
	claims, err := r.ValidateToken(tokenString)
	if err != nil {
		return
	}
	r.db.Model(&business.Session{}).
		Where("id = ?", claims.SessionID).
		Update("revoked_at", time.Now())
}
