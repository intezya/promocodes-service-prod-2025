package promocode

import (
	"solution/internal/domain/types"
	"solution/pkg"
	"time"
)

func (p *PromoCode) ToOwnerViewUNIQUE(
	active bool,
	likes,
	uses int,
) map[string]interface{} {
	var af, au string
	if p.ActiveFrom != nil {
		af = p.ActiveFrom.Format(time.DateOnly)
	}
	if p.ActiveUntil != nil {
		au = p.ActiveUntil.Format(time.DateOnly)
	}
	t := map[string]interface{}{
		"age_from":   p.TargetAgeFrom,
		"age_until":  p.TargetAgeUntil,
		"country":    p.TargetCountry,
		"categories": p.TargetCategories,
	}
	r := map[string]interface{}{
		"active":       active,
		"company_id":   p.CompanyID,
		"company_name": p.CompanyName,
		"description":  p.Description,
		"like_count":   likes,
		"mode":         p.Mode,
		"promo_id":     p.ID,
		"target":       t,
		"used_count":   uses,
		"promo_unique": p.Promo,
		"image_url":    p.ImageURL,
		"active_from":  af,
		"active_until": au,
		"max_count":    1,
	}
	pkg.RecursiveRemoveNulls(r)
	if r["target"] == nil {
		r["target"] = map[string]interface{}{}
	}
	if r["active_from"] == "0001-01-01T00:00:00Z" {
		delete(r, "active_from")
	}
	if r["active_until"] == "0001-01-01T00:00:00Z" {
		delete(r, "active_until")
	}
	return r
}
func (p *PromoCode) ToOwnerViewCOMMON(
	active bool,
	likes,
	uses int,
) map[string]interface{} {
	var af, au string
	if p.ActiveFrom != nil {
		af = p.ActiveFrom.Format(time.DateOnly)
	}
	if p.ActiveUntil != nil {
		au = p.ActiveUntil.Format(time.DateOnly)
	}
	t := map[string]interface{}{
		"age_from":   p.TargetAgeFrom,
		"age_until":  p.TargetAgeUntil,
		"country":    p.TargetCountry,
		"categories": p.TargetCategories,
	}
	r := map[string]interface{}{
		"active":       active,
		"company_id":   p.CompanyID,
		"company_name": p.CompanyName,
		"description":  p.Description,
		"like_count":   likes,
		"max_count":    p.MaxCount,
		"mode":         p.Mode,
		"promo_id":     p.ID,
		"target":       t,
		"used_count":   uses,
		"promo_common": p.Promo[0],
		"image_url":    p.ImageURL,
		"active_from":  af,
		"active_until": au,
	}
	pkg.RecursiveRemoveNulls(r)
	if r["target"] == nil {
		r["target"] = map[string]interface{}{}
	}
	if r["active_from"] == "0001-01-01" {
		delete(r, "active_from")
	}
	if r["active_until"] == "0001-01-01" {
		delete(r, "active_until")
	}
	return r
}

func (p *PromoCode) ToUserView(
	active bool, activated, liked bool,
	likes, comments int,
) map[string]interface{} {
	r := map[string]interface{}{
		"promo_id":             p.ID,
		"company_id":           p.CompanyID,
		"company_name":         p.CompanyName,
		"description":          p.Description,
		"active":               active,
		"is_activated_by_user": activated,
		"like_count":           likes,
		"is_liked_by_user":     liked,
		"comment_count":        comments,
		"image_url":            p.ImageURL,
	}
	if r["active"] == nil {
		delete(r, "active")
	}
	return r
}

type UpdatePromoCode struct {
	Description *string
	ImageURL    *string
	Target      *struct {
		AgeFrom    *int
		AgeUntil   *int
		Country    *string
		Categories *[]string
	} `json:"target"`
	MaxCount    *int
	ActiveFrom  *types.SolutionDate
	ActiveUntil *types.SolutionDate
}

type PromoSimpleData struct {
	Active   bool
	Likes    int
	Uses     int
	Comments int
}
