package customerrors

import "github.com/gofiber/fiber/v2"

type DomainError struct {
	Code        int
	Message     string
	DebugDetail string
}

func Unauthorized(detail ...string) *DomainError {
	if len(detail) > 0 {
		return &DomainError{
			Code:        401,
			Message:     "unauthorized",
			DebugDetail: detail[0],
		}
	}
	return &DomainError{
		Code:    401,
		Message: "unauthorized",
	}
}

func BadRequest(detail ...string) *DomainError {
	if len(detail) > 0 {
		return &DomainError{
			Code:        400,
			Message:     "bad request",
			DebugDetail: detail[0],
		}
	}
	return &DomainError{
		Code:    400,
		Message: "bad request",
	}
}

func (e *DomainError) Error() string {
	return e.Message
}

func NotFound(details ...string) *DomainError {
	if len(details) > 0 {
		return &DomainError{
			Code:        404,
			Message:     "not found",
			DebugDetail: details[0],
		}
	}
	return &DomainError{
		Code:    404,
		Message: "not found",
	}
}

func Forbidden() *DomainError {
	return &DomainError{
		Code:    403,
		Message: "forbidden",
	}
}

func (e *DomainError) ToFiber(c *fiber.Ctx) error {
	return c.Status(e.Code).JSON(
		fiber.Map{
			"message": e.Message,
			"detail":  e.DebugDetail,
		},
	)
}
