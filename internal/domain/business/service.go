package business

import (
	"github.com/google/uuid"
	"github.com/intezya/pkglib"
	"solution/internal/domain/errors"
)

type DomainService struct {
	repo Repository
	tm   TokenManager
}

func NewDomainService(repo Repository, tm TokenManager) *DomainService {
	return &DomainService{repo: repo, tm: tm}
}

func (s *DomainService) Create(name, email, password string) (
	id uuid.UUID,
	token string,
	err *customerrors.DomainError,
) {
	salt := pkglib.Crypto.Salt(32)
	business := &Business{
		ID:           uuid.New(),
		CompanyName:  name,
		Email:        email,
		PasswordHash: pkglib.Crypto.EncodeBase64(pkglib.Crypto.HashArgon2(password, salt)),
	}
	if err := s.repo.Create(business); err != nil {
		return uuid.Nil, "", err.ToDomain()
	}
	return business.ID, s.tm.GenerateToken(business.ID, email), nil
}

func (s *DomainService) Authorize(email, password string) (token string, err *customerrors.DomainError) {
	if business, err := s.repo.GetByEmail(email); err != nil {
		return "", customerrors.Unauthorized("business not found")
	} else {
		pw := pkglib.Crypto.DecodeBase64(business.PasswordHash)
		ok := pkglib.Crypto.VerifyArgon2(password, pw)
		if !ok {
			return "", customerrors.Unauthorized("wrong password")
		}
		return s.tm.GenerateToken(business.ID, email), nil
	}
}

func (s *DomainService) GetByID(id uuid.UUID) (*Business, *customerrors.DomainError) {
	b, err := s.repo.Get(id)
	if err != nil {
		return nil, err.ToDomain()
	}
	return b, nil
}
