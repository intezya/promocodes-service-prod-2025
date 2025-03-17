package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"solution/internal/application/business"
	customerrors "solution/internal/domain/errors"
	"solution/internal/domain/promocode"
	"solution/pkg"
	"strconv"
)

type BusinessAPI struct {
	as *business.ApplicationService
}

func NewBusinessAPI(as *business.ApplicationService) *BusinessAPI {
	return &BusinessAPI{as: as}
}

func (b *BusinessAPI) SignUp(c *fiber.Ctx) error {
	request := &business.CreateBusinessRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, err := b.as.SignUp(request)
	if err != nil {
		return err.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (b *BusinessAPI) SignIn(c *fiber.Ctx) error {
	request := &business.LoginBusinessRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	response, err := b.as.SignIn(request)
	if err != nil {
		return err.ToFiber(c)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (b *BusinessAPI) CreatePromoCode(c *fiber.Ctx) error {
	request := &business.CreatePromoCodeRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	if request.Mode == promocode.UNIQUE && request.MaxCount != nil && *request.MaxCount != 1 {
		return customerrors.BadRequest("max_count must be 1 for unique mode").ToFiber(c)
	}
	if request.Target == nil {
		return customerrors.BadRequest("target is required").ToFiber(c)
	}
	if request.Target.AgeFrom != nil && request.Target.AgeUntil != nil && *request.Target.AgeFrom > *request.Target.AgeUntil {
		return customerrors.BadRequest("age_from must be less than age_until").ToFiber(c)
	}
	companyID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := b.as.CreatePromoCode(companyID, request)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(response)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (b *BusinessAPI) GetPromoCodes(c *fiber.Ctx) error {
	params := &business.GetPromoCodesQueryParams{}
	if err := params.Bind(c); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	companyID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, count, er := b.as.GetAllPromoCodes(companyID, params)
	if er != nil {
		return er.ToFiber(c)
	}
	zap.S().Debugw("get promo codes response", "c", count, "offset", params.Offset)
	c.Set("X-Total-Count", strconv.Itoa(count))
	return c.Status(fiber.StatusOK).JSON(response)
}
func (b *BusinessAPI) GetPromoCode(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	companyID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := b.as.GetPromoCode(companyID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	c.Set("X-Total-Count", "1")
	return c.Status(fiber.StatusOK).JSON(response)
}

func (b *BusinessAPI) EditPromoCode(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	companyID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	request := &business.EditPromoCodeRequest{}
	if err := request.Bind(c, v); err != nil {
		return customerrors.BadRequest("req " + err.Error()).ToFiber(c)
	}
	if request.Target != nil &&
		request.Target.AgeFrom != nil &&
		request.Target.AgeUntil != nil &&
		*request.Target.AgeFrom > *request.Target.AgeUntil {
		return customerrors.BadRequest("age_from must be less than age_until").ToFiber(c)
	}
	zap.S().Debug(request)
	zap.S().Debugw("edit promo code", "companyID", companyID, "promoID", promoID, "request", request)
	response, er := b.as.EditPromoCode(companyID, promoID, request)
	if er != nil {
		zap.S().Info(er.Message)
		return er.ToFiber(c)
	}
	zap.S().Debugw("edit promo code", "response", response)
	pkg.RecursiveRemoveNulls(response)
	if response["target"] == nil {
		response["target"] = fiber.Map{}
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (b *BusinessAPI) UsageStatistic(c *fiber.Ctx) error {
	promoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return customerrors.BadRequest("promo_id" + err.Error())
	}
	companyID, err := uuid.Parse(c.Locals("sub").(string))
	if err != nil {
		return customerrors.BadRequest("sub " + err.Error()).ToFiber(c)
	}
	response, er := b.as.GetUsageStatistic(companyID, promoID)
	if er != nil {
		return er.ToFiber(c)
	}
	pkg.RecursiveRemoveNulls(response)
	return c.Status(fiber.StatusOK).JSON(response)
}
