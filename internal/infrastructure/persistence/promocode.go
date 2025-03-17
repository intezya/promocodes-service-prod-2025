package persistence

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"solution/internal/domain/errors"
	"solution/internal/domain/promocode"
	"solution/internal/domain/user"
	"strings"
	"time"
)

type PromoCodeRepository struct {
	db *gorm.DB
}

func (r *PromoCodeRepository) Create(p *promocode.PromoCode) error {
	r.db.Create(p)
	return nil
}

func (r *PromoCodeRepository) Get(id uuid.UUID) (*promocode.PromoCode, *customerrors.RepositoryError) {
	found := &promocode.PromoCode{}

	result := r.db.Model(&promocode.PromoCode{}).First(found, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerrors.NotFoundInRepository()
		}
		return nil, customerrors.UnknownErrorInRepository(result.Error.Error())
	}
	return found, nil
}

func (r *PromoCodeRepository) GetByCompanyIDAsCompanyList(
	id uuid.UUID,
	params *promocode.GetAsCompanyListParams,
) ([]*promocode.PromoCode, int) {
	var count int64
	query := r.db.Model(&promocode.PromoCode{}).Where("company_id = ?", id)
	if params.CountryCode != nil && len(*params.CountryCode) > 0 {
		query = query.Where("target_country_lower in ? OR target_country_lower is NULL", *params.CountryCode)
	}
	query.Count(&count)
	if params.SortBy == "" {
		query = query.Order("created_at DESC")
	} else {
		if params.SortBy == "active_from" {
			query = query.Order("active_from DESC")
		} else if params.SortBy == "active_until" {
			query = query.Order("active_until ASC")
		}
	}
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}
	var promos []*promocode.PromoCode
	query.Find(&promos)
	return promos, int(count)
}

func (r *PromoCodeRepository) GetAsUserFeed(params *promocode.GetAsUserFeedParams) ([]*promocode.PromoCode, int) {
	query := r.db.Model(&promocode.PromoCode{})
	var count int64
	if params.Category != "" {
		query = query.Where("? = ANY(target_categories_lower)", strings.ToLower(params.Category))
	}
	query.Count(&count)
	zap.S().Debugw("after categories", "count", count)

	query = query.Where(
		"(target_age_from IS NULL OR ? >= target_age_from) AND (target_age_until IS NULL OR ? <= target_age_until)",
		params.Age, params.Age,
	)
	query.Count(&count)
	zap.S().Debugw("after age", "count", count)

	if params.Country != "" {
		query = query.Where("target_country_lower = ? OR target_country_lower IS NULL", strings.ToLower(params.Country))
	}
	query.Count(&count)
	zap.S().Debugw("after country", "count", count)

	if params.Active != nil {
		now := time.Now()
		if *params.Active {
			query = query.Where(
				"(mode = ? AND max_count > used_count) OR (mode = ? AND array_length(available_promo, 1) > 0)",
				"COMMON", "UNIQUE",
			).Where(
				`
(active_from <= ? AND active_until >= ?) OR
(active_from <= ? AND active_until = '0001-01-01') OR
(active_until = '0001-01-01' AND active_from <= ?)
`, now, now, now, now,
			)
		} else {
			query = query.Where(
				`
(mode = ? AND max_count <= used_count) OR
(mode = ? AND array_length(available_promo, 1) = 0) OR
active_from > ? OR
(active_until < ? AND active_until != '0001-01-01')
`,
				"COMMON", "UNIQUE", now, now,
			)
		}
	}
	query.Count(&count)
	zap.S().Debugw("after active", "count", count)
	query.Count(&count)
	query.Order("created_at DESC")
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	query = query.Offset(params.Offset)

	var promoCodes []*promocode.PromoCode
	query.Find(&promoCodes)
	return promoCodes, int(count)
}

func (r *PromoCodeRepository) Delete(id uuid.UUID) *customerrors.RepositoryError {
	err := r.db.Model(&promocode.PromoCode{}).Where("id = ?", id).Delete(&promocode.PromoCode{}).Error
	if err != nil {
		return customerrors.NotFoundInRepository()
	}
	return nil
}

func (r *PromoCodeRepository) GetUsesCount(promoCodeID uuid.UUID) int {
	var count int64
	r.db.Model(&promocode.Use{}).Where("promo_code_id = ?", promoCodeID).Count(&count)
	return int(count)
}

func NewPromoCodeRepository(db *gorm.DB) *PromoCodeRepository {
	return &PromoCodeRepository{db: db}
}

func (r *PromoCodeRepository) GetLikesCount(promoCodeID uuid.UUID) int {
	var count int64
	r.db.Model(&promocode.Like{}).Where("promo_code_id = ?", promoCodeID).Count(&count)
	return int(count)
}

func (r *PromoCodeRepository) GetCommentsCount(promoCodeID uuid.UUID) int {
	var count int64
	r.db.Model(&promocode.Comment{}).Where("promo_code_id = ?", promoCodeID).Count(&count)
	zap.S().Debugw("comments count", "count", count)
	return int(count)
}

func (r *PromoCodeRepository) Save(p *promocode.PromoCode) {
	r.db.Save(p)
}

func (r *PromoCodeRepository) GetUsageStatistics(promoCodeID uuid.UUID) map[string]interface{} {
	var uses []promocode.Use
	var stats = make(map[string]interface{})
	var countryCounts = make(map[string]int)

	result := r.db.Where("promo_code_id = ?", promoCodeID).Find(&uses)
	totalActivations := int(result.RowsAffected)
	stats["activations_count"] = totalActivations

	for _, use := range uses {
		countryCounts[use.CountryLower]++
	}
	var countries []map[string]interface{}
	for country, count := range countryCounts {
		countries = append(
			countries, map[string]interface{}{
				"country":           country,
				"activations_count": count,
			},
		)
	}
	stats["countries"] = countries
	return stats
}

func (r *PromoCodeRepository) IsLiked(promoCodeID uuid.UUID, userID uuid.UUID) bool {
	var count int64
	r.db.Model(&promocode.Like{}).Where("promo_code_id = ?", promoCodeID).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

func (r *PromoCodeRepository) IsActivated(promoCodeID uuid.UUID, userID uuid.UUID) bool {
	var count int64
	r.db.Model(&promocode.Use{}).Where("promo_code_id = ?", promoCodeID).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

func (r *PromoCodeRepository) Like(id uuid.UUID, sub uuid.UUID) *customerrors.RepositoryError {
	result := r.db.Create(&promocode.Like{PromoCodeID: id, UserID: sub})
	if result.Error != nil {
		return &customerrors.RepositoryError{
			Code:        409,
			Message:     "already exists",
			DebugDetail: "",
		}
	}
	return nil
}
func (r *PromoCodeRepository) Unlike(id uuid.UUID, sub uuid.UUID) *customerrors.RepositoryError {
	result := r.db.Where("promo_code_id = ?", id).Where("user_id = ?", sub).Delete(&promocode.Like{})
	if result.Error != nil {
		return customerrors.NotFoundInRepository()
	}
	return nil
}

func (r *PromoCodeRepository) Comment(id uuid.UUID, sub uuid.UUID, comment string) *promocode.Comment {
	c := &promocode.Comment{
		ID:          uuid.New(),
		PromoCodeID: id,
		UserID:      sub,
		Content:     comment,
	}
	r.db.Model(&promocode.Comment{}).Create(c)
	return c
}

func (r *PromoCodeRepository) EditComment(
	comment uuid.UUID,
	promo uuid.UUID,
	commentText string,
) *customerrors.RepositoryError {
	err := r.db.Model(&promocode.Comment{}).Where("id = ? AND promo_code_id = ?", comment, promo).Update(
		"content",
		commentText,
	).Error
	if err != nil {
		return customerrors.NotFoundInRepository()
	}
	return nil
}

func (r *PromoCodeRepository) DeleteComment(comment uuid.UUID, promo uuid.UUID) *customerrors.RepositoryError {
	err := r.db.Model(&promocode.Comment{}).Where(
		"id = ? AND promo_code_id = ?",
		comment,
		promo,
	).Delete(&promocode.Comment{}).Error
	if err != nil {
		return customerrors.NotFoundInRepository()
	}
	return nil
}

func (r *PromoCodeRepository) GetComment(commentID uuid.UUID, promo uuid.UUID) (
	*promocode.CommentView,
	*customerrors.RepositoryError,
) {
	var comment promocode.Comment
	var u user.User

	result := r.db.First(&comment, "id = ? AND promo_code_id = ?", commentID, promo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerrors.NotFoundInRepository()
		}
		return nil, customerrors.UnknownErrorInRepository(result.Error.Error())
	}

	userResult := r.db.First(&u, comment.UserID)
	if userResult.Error != nil {
		return nil, customerrors.UnknownErrorInRepository(userResult.Error.Error())
	}

	return &promocode.CommentView{
		Id:   comment.ID,
		Text: comment.Content,
		Date: comment.CreatedAt.Format(time.RFC3339),
		Author: struct {
			Id        uuid.UUID `json:"id"`
			Name      string    `json:"name"`
			Surname   string    `json:"surname"`
			AvatarUrl *string   `json:"avatar_url"`
		}{
			Id:        u.ID,
			Name:      u.Name,
			Surname:   u.Surname,
			AvatarUrl: u.AvatarURL,
		},
	}, nil
}

func (r *PromoCodeRepository) GetComments(promoID uuid.UUID) []*promocode.CommentView {
	var comments []promocode.Comment
	var commentViews []*promocode.CommentView

	r.db.Where("promo_code_id = ?", promoID).Order("created_at DESC").Find(&comments)

	for _, comment := range comments {
		var u user.User
		userResult := r.db.First(&u, comment.UserID)
		if userResult.Error != nil {
			continue
		}

		commentViews = append(
			commentViews, &promocode.CommentView{
				Id:   comment.ID,
				Text: comment.Content,
				Date: comment.CreatedAt.Format(time.RFC3339),
				Author: struct {
					Id        uuid.UUID `json:"id"`
					Name      string    `json:"name"`
					Surname   string    `json:"surname"`
					AvatarUrl *string   `json:"avatar_url"`
				}{
					Id:        u.ID,
					Name:      u.Name,
					Surname:   u.Surname,
					AvatarUrl: u.AvatarURL,
				},
			},
		)
	}

	return commentViews
}

func (r *PromoCodeRepository) AddUse(u *promocode.Use) *customerrors.RepositoryError {
	r.db.Create(u)
	return nil
}

func (r *PromoCodeRepository) UseHistory(id uuid.UUID) []*promocode.Use {
	var uses []*promocode.Use
	r.db.Where("user_id = ?", id).Order("created_at DESC").Find(&uses)
	return uses
}
