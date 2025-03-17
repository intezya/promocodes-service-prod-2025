package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
	"solution/internal/application/user"
	customerrors "solution/internal/domain/errors"
	"solution/pkg"
	"strconv"
)

type UserAPI struct {
	userAS *user.ApplicationService
}

func NewUserAPI(userAS *user.ApplicationService) *UserAPI {
	return &UserAPI{
		userAS: userAS,
	}
}

func (u *UserAPI) SignUp(c *fiber.Ctx) error {
	request := &user.CreateUserRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, err := u.userAS.SignUp(request)
	if err != nil {
		return err.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) SignIn(c *fiber.Ctx) error {
	request := &user.LoginUserRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, err := u.userAS.SignIn(request)
	if err != nil {
		return err.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) GetProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := u.userAS.GetProfile(userID)
	if er != nil {
		return er.ToFiber(c)
	}
	if response.AvatarUrl == nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"name":    response.Name,
				"surname": response.Surname,
				"email":   response.Email,
				"other": fiber.Map{
					"age":     response.Other.Age,
					"country": response.Other.Country,
				},
			},
		)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) EditProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	request := &user.EditProfileRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, er := u.userAS.EditProfile(userID, request)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(response)
	if response.AvatarUrl == nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"name":    response.Name,
				"surname": response.Surname,
				"email":   response.Email,
				"other": fiber.Map{
					"age":     response.Other.Age,
					"country": response.Other.Country,
				},
			},
		)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) GetFeed(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	params := &user.GetPromoFeedQueryParams{}
	if err := params.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, count, er := u.userAS.GetFeed(userID, params)
	if er != nil {
		return er.ToFiber(c)
	}
	log.Info(len(response))
	pkg.RecursiveRemoveNulls(response)
	log.Info(len(response))
	c.Set("X-Total-Count", strconv.Itoa(count))
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) GetPromoCode(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	sub, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := u.userAS.GetPromoCode(sub, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(response)
	c.Set("X-Total-Count", strconv.Itoa(len(response)))

	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) Like(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	er := u.userAS.LikePromoCode(userID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status": "ok",
		},
	)
}

func (u *UserAPI) Unlike(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	er := u.userAS.UnlikePromoCode(userID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status": "ok",
		},
	)
}

func (u *UserAPI) CommentPromoCode(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	request := &user.CommentPromoRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, er := u.userAS.CommentPromoCode(userID, promoID, request.Content)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(response)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (u *UserAPI) GetPromoCodeComment(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("comment_id"))
	if err != nil {
		return customerrors.BadRequest("comment_id" + err.Error())
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	response, er := u.userAS.GetComment(commentID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	c.Set("X-Total-Count", "1")
	pkg.RecursiveRemoveNulls(response)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) EditPromoCodeComment(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("comment_id"))
	if err != nil {
		return customerrors.BadRequest("comment_id" + err.Error())
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	request := &user.CommentPromoRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	sub, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	p, er := u.userAS.EditComment(sub, commentID, request.Content, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(p)
	return c.Status(fiber.StatusOK).JSON(p)
}

func (u *UserAPI) DeletePromoCodeComment(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("comment_id"))
	if err != nil {
		return customerrors.BadRequest("comment_id" + err.Error())
	}
	sub, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	er := u.userAS.DeleteComment(sub, commentID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status": "ok",
		},
	)
}

func (u *UserAPI) GetPromoCodeComments(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	response, er := u.userAS.GetPromoCodeComments(promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(response)))
	pkg.RecursiveRemoveNulls(response)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (u *UserAPI) ActivatePromoCode(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	sub, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := u.userAS.ActivatePromoCode(sub, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

type GetUseHistoryQueryParams struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}

func (p *GetUseHistoryQueryParams) Bind(c *fiber.Ctx) error {
	if err := c.QueryParser(p); err != nil {
		return err
	}
	return nil
}

func (u *UserAPI) GetUseHistory(c *fiber.Ctx) error {
	sub, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	p := GetUseHistoryQueryParams{}
	if err := p.Bind(c); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	zap.S().Debugw("p", "p", p)
	response, er := u.userAS.UseHistory(sub)
	if er != nil {
		return er.ToFiber(c)
	}
	zap.S().Debugw("response", "response", response)
	var offset, limit int
	if p.Limit != nil {
		limit = *p.Limit
	} else {
		limit = math.MaxInt
	}
	if p.Offset != nil {
		offset = *p.Offset
	}
	start := offset
	if start > len(response) {
		start = len(response)
	}
	end := start + limit
	if end > len(response) {
		end = len(response)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(response)))
	limited := response[start:end]
	pkg.RecursiveRemoveNulls(limited)
	return c.Status(fiber.StatusOK).JSON(limited)
}
