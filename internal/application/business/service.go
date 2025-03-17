package business

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"solution/config"
	"solution/internal/domain/business"
	customerrors "solution/internal/domain/errors"
	"solution/internal/domain/promocode"
	"strings"
	"time"
)

type ApplicationService struct {
	ds      *business.DomainService
	promoDS *promocode.DomainService
	cfg     *config.Config
}

func NewApplicationService(
	ds *business.DomainService,
	promoDS *promocode.DomainService,
	cfg *config.Config,
) *ApplicationService {
	return &ApplicationService{ds: ds, promoDS: promoDS, cfg: cfg}
}

func (s *ApplicationService) SignUp(
	request *CreateBusinessRequest,
) (
	*CreateBusinessResponse,
	*customerrors.DomainError,
) {
	id, token, err := s.ds.Create(request.Name, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &CreateBusinessResponse{
		ID:    id,
		Token: token,
	}, nil
}

func (s *ApplicationService) SignIn(
	request *LoginBusinessRequest,
) (
	*LoginBusinessResponse,
	*customerrors.DomainError,
) {
	token, err := s.ds.Authorize(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &LoginBusinessResponse{
		Token: token,
	}, nil
}

func (s *ApplicationService) CreatePromoCode(
	sub uuid.UUID,
	request *CreatePromoCodeRequest,
) (
	*CreatePromoCodeResponse,
	*customerrors.DomainError,
) {
	var activeFrom, activeUntil time.Time
	activeFrom, err := request.ActiveFrom.ToDate()
	if err != nil {
		return nil, customerrors.BadRequest("active from")
	}
	activeUntil, err = request.ActiveUntil.ToDate()
	if err != nil {
		return nil, customerrors.BadRequest("active until")
	}
	company, er := s.ds.GetByID(sub)
	if er != nil {
		return nil, customerrors.NotFound(er.Message)
	}
	var promoID uuid.UUID
	var countryLower *string
	if request.Target.Country != nil {
		res := strings.ToLower(*request.Target.Country)
		countryLower = &res
	}
	var categoriesLower []string
	if request.Target.Categories != nil {
		for i := 0; i < len(*request.Target.Categories); i++ {
			res := strings.ToLower((*request.Target.Categories)[i])
			categoriesLower = append(categoriesLower, res)
		}
	}
	if request.Mode == promocode.COMMON {
		promo := &promocode.PromoCode{
			ID:                    uuid.New(),
			Description:           request.Description,
			CompanyID:             company.ID,
			CompanyName:           company.CompanyName,
			MaxCount:              *request.MaxCount,
			TargetAgeFrom:         request.Target.AgeFrom,
			TargetAgeUntil:        request.Target.AgeUntil,
			TargetCountry:         request.Target.Country,
			TargetCountryLower:    countryLower,
			TargetCategories:      (*pq.StringArray)(request.Target.Categories),
			Promo:                 pq.StringArray{*request.PromoCommon},
			ImageURL:              request.ImageURL,
			ActiveFrom:            &activeFrom,
			ActiveUntil:           &activeUntil,
			Mode:                  promocode.COMMON,
			TargetCategoriesLower: (*pq.StringArray)(&categoriesLower),
		}
		_ = s.promoDS.Create(promo)
		promoID = promo.ID
	} else if request.Mode == promocode.UNIQUE {
		promo := &promocode.PromoCode{
			ID:                    uuid.New(),
			Description:           request.Description,
			CompanyID:             company.ID,
			CompanyName:           company.CompanyName,
			TargetAgeFrom:         request.Target.AgeFrom,
			TargetAgeUntil:        request.Target.AgeUntil,
			TargetCountry:         request.Target.Country,
			TargetCountryLower:    countryLower,
			TargetCategories:      (*pq.StringArray)(request.Target.Categories),
			Promo:                 *request.PromoUnique,
			AvailablePromo:        *request.PromoUnique,
			ImageURL:              request.ImageURL,
			ActiveFrom:            &activeFrom,
			ActiveUntil:           &activeUntil,
			Mode:                  promocode.UNIQUE,
			TargetCategoriesLower: (*pq.StringArray)(&categoriesLower),
		}
		_ = s.promoDS.Create(promo)
		promoID = promo.ID
	}
	return &CreatePromoCodeResponse{
		ID: promoID,
	}, nil
}

func (s *ApplicationService) EditPromoCode(
	sub uuid.UUID,
	promoID uuid.UUID,
	request *EditPromoCodeRequest,
) (
	map[string]interface{},
	*customerrors.DomainError,
) {
	promo, err := s.promoDS.Get(promoID)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	if promo.CompanyID != sub {
		return nil, customerrors.Forbidden()
	}
	if promo.Mode == promocode.UNIQUE && request.MaxCount != nil {
		return nil, customerrors.BadRequest("max count")
	}
	zap.S().Infow("EditPromoCode", "request", request)
	p, d, er := s.promoDS.Update(promo.ID, (*promocode.UpdatePromoCode)(request))
	if er != nil {
		return nil, er
	}
	if p.Mode == promocode.UNIQUE {
		return p.ToOwnerViewUNIQUE(d.Active, d.Likes, d.Uses), nil
	} else if p.Mode == promocode.COMMON {
		return p.ToOwnerViewCOMMON(d.Active, d.Likes, d.Uses), nil
	}
	return nil, customerrors.NotFound()
}

func (s *ApplicationService) GetAllPromoCodes(
	sub uuid.UUID,
	p *GetPromoCodesQueryParams,
) ([]map[string]interface{}, int, *customerrors.DomainError) {
	company, err := s.ds.GetByID(sub)
	if err != nil {
		return nil, 0, customerrors.NotFound()
	}
	zap.S().Debugf("GetAllPromoCodes: %+v", p)
	res, c := s.promoDS.GetByCompanyID(company.ID, p.Limit, p.Offset, p.SortBy, p.Countries)
	return res, c, nil
}

func (s *ApplicationService) GetPromoCode(
	sub uuid.UUID,
	promoID uuid.UUID,
) (map[string]interface{}, *customerrors.DomainError) {
	company, err := s.ds.GetByID(sub)
	if err != nil {
		zap.S().Info(err.Message)
		return nil, customerrors.NotFound()
	}
	p, d, er := s.promoDS.Get2(promoID)
	if er != nil {
		zap.S().Info(er.Message)
		return nil, customerrors.NotFound()
	}
	if p.CompanyID != company.ID {
		return nil, customerrors.Forbidden()
	}
	if p.Mode == promocode.UNIQUE {
		return p.ToOwnerViewUNIQUE(d.Active, d.Likes, d.Uses), nil
	} else if p.Mode == promocode.COMMON {
		return p.ToOwnerViewCOMMON(d.Active, d.Likes, d.Uses), nil
	}
	zap.S().Info(p.Mode)
	return nil, customerrors.NotFound()
}

func (s *ApplicationService) GetUsageStatistic(sub uuid.UUID, promo uuid.UUID) (
	map[string]interface{},
	*customerrors.DomainError,
) {
	p, _, err := s.promoDS.Get2(promo)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	if p.CompanyID != sub {
		return nil, customerrors.Forbidden()
	}
	return s.promoDS.UsageStatistic(promo), nil
}
