package promocode

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	customerrors "solution/internal/domain/errors"
	"strings"
	"time"
)

type DomainService struct {
	repository Repository
}

func NewDomainService(repository Repository) *DomainService {
	return &DomainService{
		repository: repository,
	}
}

func (d *DomainService) Create(p *PromoCode) error {
	return d.repository.Create(p)
}

func (d *DomainService) Get(id uuid.UUID) (*PromoCode, error) {
	p, err := d.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *DomainService) Get2(id uuid.UUID) (p *PromoCode, dt *PromoSimpleData, err *customerrors.DomainError) {
	p, er := d.repository.Get(id)
	if er != nil {
		return nil, nil, er.ToDomain()
	}
	return p, &PromoSimpleData{
		Active:   d.IsActive(p),
		Likes:    d.repository.GetLikesCount(p.ID),
		Uses:     d.repository.GetUsesCount(p.ID),
		Comments: d.repository.GetCommentsCount(p.ID),
	}, nil
}

func (d *DomainService) Update(id uuid.UUID, u *UpdatePromoCode) (
	p *PromoCode,
	data *PromoSimpleData,
	err *customerrors.DomainError,
) {
	p, errr := d.repository.Get(id)
	if errr != nil {
		return nil, nil, errr.ToDomain()
	}
	au, er := u.ActiveUntil.ToDate()
	if er != nil {
		return nil, nil, customerrors.BadRequest("active until")
	}
	af, er := u.ActiveFrom.ToDate()
	if er != nil {
		return nil, nil, customerrors.BadRequest("active from")
	}
	if u.Description != nil {
		p.Description = *u.Description
	}
	if u.ImageURL != nil {
		p.ImageURL = u.ImageURL
	}
	if u.Target != nil {
		if u.Target.AgeFrom != nil {
			p.TargetAgeFrom = u.Target.AgeFrom
		}
		if u.Target.AgeUntil != nil {

			p.TargetAgeUntil = u.Target.AgeUntil
		}
		if u.Target.Country != nil {
			p.TargetCountry = u.Target.Country
			sd := strings.ToLower(*u.Target.Country)
			p.TargetCountryLower = &sd
		}
		zap.S().Info(u.Target.Categories)
		if u.Target.Categories != nil {
			p.TargetCategories = (*pq.StringArray)(u.Target.Categories)
			var l []string
			for _, v := range *p.TargetCategories {
				l = append(l, strings.ToLower(v))
			}
			p.TargetCategoriesLower = (*pq.StringArray)(&l)
		}
	}
	if u.ActiveFrom != nil {
		p.ActiveFrom = &af
	}
	if u.ActiveUntil != nil {
		p.ActiveUntil = &au
	}
	if p.Mode == COMMON {
		if u.MaxCount != nil {
			p.MaxCount = *u.MaxCount
		}
	}
	lc := d.repository.GetLikesCount(p.ID)
	uc := d.repository.GetUsesCount(p.ID)
	d.repository.Save(p)
	return p, &PromoSimpleData{
		Active: d.IsActive(p),
		Likes:  lc,
		Uses:   uc,
	}, nil
}

func (d *DomainService) IsActive(p *PromoCode) bool {
	if p.Mode == UNIQUE {
		if len(p.AvailablePromo) == 0 {
			return false
		}
	}
	if p.Mode == COMMON {
		if p.UsedCount >= p.MaxCount {
			return false
		}
	}
	if p.ActiveFrom.Before(time.Now()) && p.ActiveUntil.Format(time.DateOnly) == "0001-01-01" {
		return true
	}
	if p.ActiveUntil.After(time.Now()) && p.ActiveFrom.Before(time.Now()) {
		return true
	}
	if p.ActiveUntil.After(time.Now()) && p.ActiveFrom.Format(time.DateOnly) == "0001-01-01" {
		return true
	}
	return false
}

func (d *DomainService) GetByCompanyID(
	id uuid.UUID,
	limit *int,
	offset int,
	sort string,
	countryCode *[]string,
) ([]map[string]interface{}, int) {
	promos, count := d.repository.GetByCompanyIDAsCompanyList(
		id, &GetAsCompanyListParams{
			Limit:       limit,
			Offset:      offset,
			SortBy:      sort,
			CountryCode: countryCode,
		},
	)
	var result []map[string]interface{}

	for _, p := range promos {
		isactive := d.IsActive(p)
		if p.Mode == COMMON {
			result = append(
				result,
				p.ToOwnerViewCOMMON(
					isactive,
					d.repository.GetLikesCount(p.ID),
					d.repository.GetUsesCount(p.ID),
				),
			)
		} else {
			result = append(
				result,
				p.ToOwnerViewUNIQUE(
					isactive,
					d.repository.GetLikesCount(p.ID),
					d.repository.GetUsesCount(p.ID),
				),
			)
		}

	}
	zap.S().Debugw("Get by company", "result", result)
	return result, count
}

func (d *DomainService) UsageStatistic(id uuid.UUID) map[string]interface{} {
	return d.repository.GetUsageStatistics(id)
}

func (d *DomainService) GetFeed(
	sub uuid.UUID,
	limit *int,
	offset int,
	category string,
	active *bool,
	age int,
	country string,
) ([]map[string]interface{}, int) {
	result, count := d.repository.GetAsUserFeed(
		&GetAsUserFeedParams{
			Limit:    limit,
			Offset:   offset,
			Category: category,
			Active:   active,
			Age:      age,
			Country:  country,
		},
	)
	var r []map[string]interface{}
	for _, p := range result {
		r = append(
			r, p.ToUserView(
				d.IsActive(p),
				d.repository.IsActivated(p.ID, sub),
				d.repository.IsLiked(p.ID, sub),
				d.repository.GetLikesCount(p.ID),
				d.repository.GetCommentsCount(p.ID),
			),
		)
	}
	return r, count
}

func (d *DomainService) Activated(id uuid.UUID, sub uuid.UUID) bool {
	return d.repository.IsActivated(id, sub)
}

func (d *DomainService) Liked(id uuid.UUID, sub uuid.UUID) bool {
	return d.repository.IsLiked(id, sub)
}

func (d *DomainService) Like(id uuid.UUID, sub uuid.UUID) *customerrors.DomainError {
	err := d.repository.Like(id, sub)
	if err != nil {
		return err.ToDomain()
	}
	return nil
}

func (d *DomainService) Unlike(id uuid.UUID, sub uuid.UUID) *customerrors.DomainError {
	err := d.repository.Unlike(id, sub)
	if err != nil {
		return customerrors.NotFound()
	}
	return nil
}

func (d *DomainService) Comment(id uuid.UUID, sub uuid.UUID, comment string) (*CommentView, *customerrors.DomainError) {
	if p, er := d.repository.Get(id); p == nil || er != nil {
		return nil, customerrors.NotFound()
	}
	c := d.repository.Comment(id, sub, comment)
	return d.GetComment(c.ID, c.PromoCodeID)
}

func (d *DomainService) GetComments(id uuid.UUID) []*CommentView {
	return d.repository.GetComments(id)
}

func (d *DomainService) GetComment(comment uuid.UUID, promo uuid.UUID) (*CommentView, *customerrors.DomainError) {
	result, err := d.repository.GetComment(comment, promo)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	return result, nil
}

func (d *DomainService) EditComment(comment uuid.UUID, promo uuid.UUID, commentText string) (
	*CommentView,
	*customerrors.DomainError,
) {
	err := d.repository.EditComment(comment, promo, commentText)
	if err != nil {
		return nil, customerrors.NotFound()
	}
	res, _ := d.repository.GetComment(comment, promo)
	return res, nil
}

func (d *DomainService) DeleteComment(comment uuid.UUID, promo uuid.UUID) *customerrors.DomainError {
	err := d.repository.DeleteComment(comment, promo)
	if err != nil {
		return customerrors.NotFound()
	}
	return nil
}

func (d *DomainService) SavePromo(p *PromoCode) {
	d.repository.Save(p)
}
func (d *DomainService) UsePromo(p uuid.UUID, u uuid.UUID, c string) *customerrors.DomainError {
	err := d.repository.AddUse(
		&Use{
			ID:           uuid.New(),
			PromoCodeID:  p,
			UserID:       u,
			Country:      c,
			CountryLower: strings.ToLower(c),
			CreatedAt:    time.Now(),
		},
	)
	if err != nil {
		return err.ToDomain()
	}
	return nil
}

func (d *DomainService) UseHistory(id uuid.UUID) []map[string]interface{} {
	uses := d.repository.UseHistory(id)
	var result []map[string]interface{}
	for _, u := range uses {
		p, pd, _ := d.Get2(u.PromoCodeID)

		result = append(
			result,
			p.ToUserView(
				pd.Active,
				d.repository.IsActivated(p.ID, u.UserID),
				d.repository.IsLiked(p.ID, u.UserID),
				pd.Likes,
				pd.Comments,
			),
		)
	}
	return result
}
