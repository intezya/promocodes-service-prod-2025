package business

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"solution/internal/domain/promocode"
	"solution/internal/domain/types"
	"solution/pkg"
	"strings"
)

type CreateBusinessRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (r *CreateBusinessRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type LoginBusinessRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (r *LoginBusinessRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type CreatePromoCodeRequest struct {
	Description string         `json:"description" validate:"required,description"`
	Mode        promocode.Mode `json:"mode" validate:"required,oneof=UNIQUE COMMON" mode_logic:"PromoCommon,PromoUnique,MaxCount"`
	MaxCount    *int           `json:"max_count" validate:"required,gte=0,lte=100000000"`
	Target      *struct {
		AgeFrom    *int      `json:"age_from" validate:"omitempty,gte=0,lte=100"`
		AgeUntil   *int      `json:"age_until" validate:"omitempty,gte=0,lte=100"`
		Country    *string   `json:"country" validate:"omitempty,country"`
		Categories *[]string `json:"categories" validate:"omitempty,dive,min=2,max=20"`
	} `json:"target" validate:"required"`
	PromoCommon *string             `json:"promo_common" validate:"required_if=Mode COMMON"`
	PromoUnique *[]string           `json:"promo_unique" validate:"required_if=Mode UNIQUE"`
	ImageURL    *string             `json:"image_url" validate:"omitempty,url"`
	ActiveFrom  *types.SolutionDate `json:"active_from"`
	ActiveUntil *types.SolutionDate `json:"active_until"`
}

func (r *CreatePromoCodeRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type EditPromoCodeRequest struct {
	Description *string `json:"description" validate:"omitempty,description"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
	Target      *struct {
		AgeFrom    *int      `json:"age_from" validate:"omitempty,gte=0,lte=100"`
		AgeUntil   *int      `json:"age_until" validate:"omitempty,gte=0,lte=100"`
		Country    *string   `json:"country" validate:"omitempty,country"`
		Categories *[]string `json:"categories" validate:"omitempty,dive,min=2,max=20"`
	} `json:"target" validate:"omitempty"`
	MaxCount    *int                `json:"max_count" validate:"omitempty,gte=0,lte=100000000"`
	ActiveFrom  *types.SolutionDate `json:"active_from"`
	ActiveUntil *types.SolutionDate `json:"active_until"`
}

func (r *EditPromoCodeRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type GetPromoCodesQueryParams struct {
	Limit     *int      `query:"limit"`
	Offset    int       `query:"offset"`
	SortBy    string    `query:"sort_by"`
	Countries *[]string `query:"country"`
}

func (r *GetPromoCodesQueryParams) Bind(c *fiber.Ctx) error {
	if err := c.QueryParser(r); err != nil {
		return err
	}
	if r.Countries == nil {
		return nil
	}
	var countries []string
	for _, v := range *r.Countries {
		countries = append(
			countries,
			pkg.Map(
				func(v2 string) string {
					return strings.ToLower(v2)
				}, strings.Split(v, ","),
			)...,
		)
	}
	r.Countries = &countries
	return nil
}
