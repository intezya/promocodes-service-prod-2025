package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	customerrors "solution/internal/domain/errors"
	"strings"
	"time"
)

type RequestForAntiFraud struct {
	Email   string    `json:"user_email" validate:"required,email"`
	PromoID uuid.UUID `json:"promo_id" validate:"required"`
}

func AntiFraudForbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(
		fiber.Map{
			"message": "forbidden",
			"debug":   "forbidden by antifraud service",
		},
	)
}

type AntiFraudServiceResponse struct {
	Ok         bool   `json:"ok"`
	CacheUntil string `json:"cache_until,omitempty"`
}
type AntiFraudServiceRequest struct {
	Email   string    `json:"user_email" validate:"required,email"`
	PromoID uuid.UUID `json:"promo_id" validate:"required"`
}

type antiFraudService struct {
	url    string
	client *http.Client
}

func (a *antiFraudService) sendToService(user string, promo uuid.UUID) (*AntiFraudServiceResponse, error) {
	body := AntiFraudServiceRequest{
		user,
		promo,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", a.url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	response, err := a.client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	if response.StatusCode != 200 {
		response, _ = a.client.Do(req)
	}
	defer response.Body.Close()
	zap.S().Debugw(
		"response",
		"body", response.Body,
		"status", response.StatusCode,
		"request", response.Request.Body,
		"remoteAddr", response.Request.RequestURI,
	)
	antiFraudResp := &AntiFraudServiceResponse{}
	_ = json.NewDecoder(response.Body).Decode(antiFraudResp)
	return antiFraudResp, nil
}

func AntiFraud(redisURL string, serviceURL string) fiber.Handler {
	if !strings.Contains(serviceURL, "://") {
		serviceURL = "http://" + serviceURL
	}
	service := &antiFraudService{
		url:    serviceURL + "/api/validate",
		client: &http.Client{},
	}
	client := redis.NewClient(&redis.Options{Addr: redisURL})
	for {
		if client.Ping(context.Background()).Err() == nil {
			break
		}
		log.Warn("redis is not ready")
		time.Sleep(time.Second)
	}
	return func(c *fiber.Ctx) error {
		pID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return customerrors.BadRequest("promo_id" + err.Error())
		}
		request := &RequestForAntiFraud{
			Email:   c.Locals("email").(string),
			PromoID: pID,
		}
		err = client.Get(c.Context(), request.Email).Err()
		if err == nil {
			return c.Next()
		}
		res, err := service.sendToService(request.Email, request.PromoID)
		if err != nil {
			zap.S().Info(err)
			return AntiFraudForbidden(c)
		}
		if !res.Ok { // forbidden
			zap.S().Info(err)
			return AntiFraudForbidden(c)
		}
		if res.CacheUntil == "" {
			return c.Next()
		}
		t, _ := time.Parse(time.RFC3339, res.CacheUntil)
		client.Set(c.Context(), request.Email, request.PromoID, time.Until(t))
		return c.Next()
	}
}
