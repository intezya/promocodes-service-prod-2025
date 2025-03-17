package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	customerrors "solution/internal/domain/errors"
	"solution/internal/domain/user"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(u *user.User) *customerrors.RepositoryError {
	err := r.db.Create(u)
	if err.Error != nil {
		return customerrors.UnknownErrorInRepository(err.Error.Error())
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, *customerrors.RepositoryError) {
	var u user.User
	err := r.db.Model(&user.User{}).Where("email = ?", email).First(&u)
	if err.Error != nil {
		return nil, customerrors.UnknownErrorInRepository(err.Error.Error())
	}
	return &u, nil
}

func (r *UserRepository) Get(id uuid.UUID) (*user.User, *customerrors.RepositoryError) {
	var u user.User
	err := r.db.Model(&user.User{}).Where("id = ?", id).First(&u)
	if err.Error != nil {
		return nil, customerrors.UnknownErrorInRepository(err.Error.Error())
	}
	return &u, nil
}

func (r *UserRepository) Save(u *user.User) {
	r.db.Save(u)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
