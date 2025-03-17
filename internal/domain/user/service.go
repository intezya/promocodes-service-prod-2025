package user

import (
	"github.com/google/uuid"
	"github.com/intezya/pkglib"
	customerrors "solution/internal/domain/errors"
)

type DomainService struct {
	repo Repository
	tm   TokenManager
}

func NewDomainService(repo Repository, tm TokenManager) *DomainService {
	return &DomainService{repo: repo, tm: tm}
}

func (s *DomainService) Create(name, surname, email, country, password string, avatarURL *string, age int) (
	id uuid.UUID,
	token string,
	err *customerrors.DomainError,
) {
	salt := pkglib.Crypto.Salt(32)
	user := &User{
		ID:           uuid.New(),
		Name:         name,
		Surname:      surname,
		Email:        email,
		Age:          age,
		Country:      country,
		PasswordHash: pkglib.Crypto.EncodeBase64(pkglib.Crypto.HashArgon2(password, salt)),
		AvatarURL:    avatarURL,
	}
	if err := s.repo.Create(user); err != nil {
		return uuid.Nil, "", &customerrors.DomainError{
			Code:        409,
			Message:     "conflict",
			DebugDetail: "email conflict",
		}
	}
	return user.ID, s.tm.GenerateToken(user.ID, user.Email), nil
}

func (s *DomainService) Authorize(email, password string) (token string, err *customerrors.DomainError) {
	if user, err := s.repo.GetByEmail(email); err != nil {
		return "", customerrors.Unauthorized("business not found")
	} else {
		pw := pkglib.Crypto.DecodeBase64(user.PasswordHash)
		ok := pkglib.Crypto.VerifyArgon2(password, pw)
		if !ok {
			return "", customerrors.Unauthorized("wrong password")
		}
		return s.tm.GenerateToken(user.ID, user.Email), nil
	}
}

func (s *DomainService) GetByID(id uuid.UUID) (*User, *customerrors.DomainError) {
	u, err := s.repo.Get(id)
	if err != nil {
		return nil, err.ToDomain()
	}
	return u, nil
}

func (s *DomainService) Save(u *User) {
	s.repo.Save(u)
}
