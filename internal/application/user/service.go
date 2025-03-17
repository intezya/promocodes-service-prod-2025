package user

import (
	"github.com/google/uuid"
	"github.com/intezya/pkglib"
	"go.uber.org/zap"
	customerrors "solution/internal/domain/errors"
	"solution/internal/domain/promocode"
	"solution/internal/domain/user"
)

type ApplicationService struct {
	ds      *user.DomainService
	promoDS *promocode.DomainService
}

func NewApplicationService(ds *user.DomainService, promoDS *promocode.DomainService) *ApplicationService {
	return &ApplicationService{ds: ds, promoDS: promoDS}
}

func (s *ApplicationService) SignUp(request *CreateUserRequest) (*CreateUserResponse, *customerrors.DomainError) {
	_, token, err := s.ds.Create(
		request.Name,
		request.Surname,
		request.Email,
		request.Other.Country,
		request.Password,
		request.AvatarURL,
		request.Other.Age,
	)
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{
		Token: token,
	}, nil
}

func (s *ApplicationService) SignIn(request *LoginUserRequest) (*LoginUserResponse, *customerrors.DomainError) {
	token, err := s.ds.Authorize(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &LoginUserResponse{
		Token: token,
	}, nil
}

func (s *ApplicationService) GetProfile(sub uuid.UUID) (*user.Profile, *customerrors.DomainError) {
	u, err := s.ds.GetByID(sub)
	if err != nil {
		return nil, err
	}
	return &user.Profile{
		Name:    u.Name,
		Surname: u.Surname,
		Email:   u.Email,
		Other: struct {
			Age     int    `json:"age"`
			Country string `json:"country"`
		}{
			Age:     u.Age,
			Country: u.Country,
		},
		AvatarUrl: u.AvatarURL,
	}, nil
}

func (s *ApplicationService) EditProfile(sub uuid.UUID, request *EditProfileRequest) (
	*user.Profile,
	*customerrors.DomainError,
) {
	u, err := s.ds.GetByID(sub)
	if err != nil {
		return nil, err
	}
	if request.Name != nil {
		u.Name = *request.Name
	}
	if request.Surname != nil {
		u.Surname = *request.Surname
	}
	if request.AvatarURL != nil {
		u.AvatarURL = request.AvatarURL
	}
	if request.Password != nil {
		u.PasswordHash = pkglib.Crypto.EncodeBase64(pkglib.Crypto.HashArgon2(*request.Password, pkglib.Crypto.Salt(32)))
	}
	s.ds.Save(u)
	return &user.Profile{
		Name:    u.Name,
		Surname: u.Surname,
		Email:   u.Email,
		Other: struct {
			Age     int    `json:"age"`
			Country string `json:"country"`
		}{
			Age:     u.Age,
			Country: u.Country,
		},
		AvatarUrl: u.AvatarURL,
	}, nil
}

func (s *ApplicationService) GetFeed(sub uuid.UUID, params *GetPromoFeedQueryParams) (
	[]map[string]interface{},
	int,
	*customerrors.DomainError,
) {
	u, err := s.ds.GetByID(sub)
	if err != nil {
		return nil, 0, err
	}
	res, c := s.promoDS.GetFeed(
		sub,
		params.Limit,
		params.Offset,
		params.Category,
		params.Active,
		u.Age,
		u.Country,
	)
	return res, c, nil
}

func (s *ApplicationService) GetPromoCode(sub uuid.UUID, promo uuid.UUID) (
	map[string]interface{},
	*customerrors.DomainError,
) {
	p, d, err := s.promoDS.Get2(promo)
	if err != nil {
		return nil, customerrors.NotFound()
	}

	return p.ToUserView(
		d.Active,
		s.promoDS.Activated(p.ID, sub),
		s.promoDS.Liked(p.ID, sub),
		d.Likes,
		d.Comments,
	), nil
}

func (s *ApplicationService) LikePromoCode(sub uuid.UUID, promo uuid.UUID) *customerrors.DomainError {
	if _, err := s.promoDS.Get(promo); err != nil {
		return customerrors.NotFound()
	}
	_ = s.promoDS.Like(promo, sub)
	return nil
}

func (s *ApplicationService) UnlikePromoCode(sub uuid.UUID, promo uuid.UUID) *customerrors.DomainError {
	return s.promoDS.Unlike(promo, sub)
}

func (s *ApplicationService) CommentPromoCode(
	sub uuid.UUID,
	promo uuid.UUID,
	comment string,
) (*promocode.CommentView, *customerrors.DomainError) {
	if _, err := s.promoDS.Get(promo); err != nil {
		zap.S().Error(err)
		return nil, customerrors.NotFound()
	}
	return s.promoDS.Comment(promo, sub, comment)
}

func (s *ApplicationService) GetPromoCodeComments(promo uuid.UUID) (
	[]*promocode.CommentView,
	*customerrors.DomainError,
) {
	_, err := s.promoDS.Get(promo)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	return s.promoDS.GetComments(promo), nil
}

func (s *ApplicationService) EditComment(sub uuid.UUID, commentID uuid.UUID, comment string, promo uuid.UUID) (
	*promocode.CommentView,
	*customerrors.DomainError,
) {
	c, err := s.promoDS.GetComment(commentID, promo)
	if err != nil {
		return nil, err
	}
	if c.Author.Id != sub {
		return nil, customerrors.Forbidden()
	}
	v, _ := s.promoDS.EditComment(commentID, promo, comment)
	return v, nil
}

func (s *ApplicationService) DeleteComment(
	sub uuid.UUID,
	comment uuid.UUID,
	promo uuid.UUID,
) *customerrors.DomainError {
	c, err := s.promoDS.GetComment(comment, promo)
	if err != nil {
		return err
	}
	if c.Author.Id != sub {
		return customerrors.Forbidden()
	}
	return s.promoDS.DeleteComment(comment, promo)
}

func (s *ApplicationService) GetComment(comment uuid.UUID, promo uuid.UUID) (
	*promocode.CommentView,
	*customerrors.DomainError,
) {
	return s.promoDS.GetComment(comment, promo)
}

func (s *ApplicationService) ActivatePromoCode(sub uuid.UUID, promo uuid.UUID) (
	*ActivatePromoResponse,
	*customerrors.DomainError,
) {
	u, _ := s.ds.GetByID(sub)
	p, err := s.promoDS.Get(promo)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	if !s.promoDS.IsActive(p) {
		return nil, customerrors.Forbidden()
	}
	if p.TargetAgeFrom != nil && u.Age < *p.TargetAgeFrom {
		return nil, customerrors.Forbidden()
	}
	if p.TargetAgeUntil != nil && u.Age > *p.TargetAgeUntil {
		return nil, customerrors.Forbidden()
	}

	if p.TargetCountry != nil && u.Country != *p.TargetCountry {
		return nil, customerrors.Forbidden()
	}

	var usedValue string
	if p.Mode == promocode.UNIQUE {
		if len(p.AvailablePromo) == 0 {
			return nil, customerrors.Forbidden()
		}
		usedValue = p.AvailablePromo[0]
		if len(p.AvailablePromo) > 1 {
			p.AvailablePromo = p.AvailablePromo[1:]
		} else {
			p.AvailablePromo = nil
		}
		s.promoDS.SavePromo(p)
	} else {
		if p.UsedCount >= p.MaxCount {
			return nil, customerrors.Forbidden()
		}
		usedValue = p.Promo[0]
		p.UsedCount += 1
	}

	er := s.promoDS.UsePromo(p.ID, sub, usedValue)
	if er != nil {
		return nil, er
	}
	s.promoDS.SavePromo(p)

	return &ActivatePromoResponse{
		Promo: usedValue,
	}, nil
}

func (s *ApplicationService) UseHistory(sub uuid.UUID) (
	[]map[string]interface{},
	*customerrors.DomainError,
) {
	return s.promoDS.UseHistory(sub), nil
}
