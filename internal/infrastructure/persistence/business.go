package persistence

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"solution/internal/domain/business"
	"solution/internal/domain/errors"
)

type BusinessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) *BusinessRepository {
	return &BusinessRepository{db: db}
}

func (r *BusinessRepository) Create(b *business.Business) *customerrors.RepositoryError {
	if err := r.db.Create(b).Error; err != nil {
		return &customerrors.RepositoryError{
			Code:    409,
			Message: "failed to create business",
		}
	}
	return nil
}

func (r *BusinessRepository) GetByEmail(email string) (*business.Business, *customerrors.RepositoryError) {
	var b business.Business
	if err := r.db.Where("email = ?", email).First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customerrors.RepositoryError{
				Code:    404,
				Message: "business not found",
			}
		}
		return nil, &customerrors.RepositoryError{
			Code:    500,
			Message: "failed to retrieve business",
		}
	}
	return &b, nil
}

func (r *BusinessRepository) Get(id uuid.UUID) (*business.Business, *customerrors.RepositoryError) {
	var b business.Business
	if err := r.db.First(&b, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customerrors.RepositoryError{
				Code:    404,
				Message: "business not found",
			}
		}
		return nil, &customerrors.RepositoryError{
			Code:    500,
			Message: "failed to retrieve business",
		}
	}
	return &b, nil
}
