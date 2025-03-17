package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Other struct {
		Age     int    `json:"age" validate:"required,gte=0,lte=100"`
		Country string `json:"country" validate:"required,country"`
	} `json:"other" validate:"required"`
	Password  string  `json:"password" validate:"required,password"`
	Surname   string  `json:"surname" validate:"required,min=1,max=120"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
}

func (r *CreateUserRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (r *LoginUserRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type EditProfileRequest struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=100"`
	Surname   *string `json:"surname" validate:"omitempty,min=1,max=120"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
	Password  *string `json:"password" validate:"omitempty,password"`
}

func (r *EditProfileRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type GetPromoFeedQueryParams struct {
	Limit    *int   `json:"limit" validate:"omitempty,gte=0"`
	Offset   int    `json:"offset" validate:"omitempty,gte=0"`
	Category string `json:"category"`
	Active   *bool  `json:"active"`
}

func (r *GetPromoFeedQueryParams) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.QueryParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}

type CommentPromoRequest struct {
	Content string `json:"text" validate:"required,min=10,max=1000"`
}

func (r *CommentPromoRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	return v.Struct(r)
}
