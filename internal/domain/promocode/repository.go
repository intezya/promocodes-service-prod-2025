package promocode

import (
	"github.com/google/uuid"
	"solution/internal/domain/errors"
	"time"
)

type GetAsCompanyListParams struct {
	Limit       *int
	Offset      int
	SortBy      string
	CountryCode *[]string
}

type GetAsUserFeedParams struct {
	Limit    *int
	Offset   int
	Category string
	Active   *bool
	Age      int
	Country  string
}

type CreatedAt = time.Time

type Repository interface {
	Create(p *PromoCode) error
	Get(id uuid.UUID) (*PromoCode, *customerrors.RepositoryError)
	GetByCompanyIDAsCompanyList(id uuid.UUID, params *GetAsCompanyListParams) ([]*PromoCode, int)
	GetAsUserFeed(params *GetAsUserFeedParams) ([]*PromoCode, int)
	GetCommentsCount(promoCodeID uuid.UUID) int
	GetLikesCount(promoCodeID uuid.UUID) int
	GetUsesCount(promoCodeID uuid.UUID) int
	Delete(id uuid.UUID) *customerrors.RepositoryError
	GetUsageStatistics(promoCodeID uuid.UUID) map[string]interface{}
	Save(p *PromoCode)

	IsLiked(promoCodeID uuid.UUID, userID uuid.UUID) bool
	IsActivated(promoCodeID uuid.UUID, userID uuid.UUID) bool

	Like(id uuid.UUID, sub uuid.UUID) *customerrors.RepositoryError
	Unlike(id uuid.UUID, sub uuid.UUID) *customerrors.RepositoryError

	Comment(id uuid.UUID, sub uuid.UUID, comment string) *Comment
	GetComment(comment uuid.UUID, promo uuid.UUID) (*CommentView, *customerrors.RepositoryError)
	EditComment(comment uuid.UUID, promo uuid.UUID, commentText string) *customerrors.RepositoryError
	DeleteComment(comment uuid.UUID, promo uuid.UUID) *customerrors.RepositoryError
	GetComments(promoid uuid.UUID) []*CommentView

	AddUse(u *Use) *customerrors.RepositoryError
	UseHistory(id uuid.UUID) []*Use
}
