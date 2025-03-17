package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"solution/config"
	business2 "solution/internal/application/business"
	user2 "solution/internal/application/user"
	"solution/internal/domain/business"
	"solution/internal/domain/promocode"
	"solution/internal/domain/user"
	"solution/internal/infrastructure/persistence"
	"solution/internal/interfaces/http"
	"solution/internal/interfaces/middleware"
)

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	server := fiber.New()
	cfg := config.New()

	api := server.Group("/api")

	api.Get(
		"/ping", func(c *fiber.Ctx) error { // 01
			return c.Status(fiber.StatusOK).SendString("PROOOOOOOOOOOOOOOOOD")
		},
	)

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: cfg.PostgresConn,
			},
		), &gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(
		&business.Business{},
		&business.Session{},
		&promocode.PromoCode{},
		&promocode.Like{},
		&promocode.Comment{},
		&promocode.Use{},
		&user.User{},
	)
	promocodeRepository := persistence.NewPromoCodeRepository(db)
	tokenManager := persistence.NewTokenManagerRepository(db, []byte(cfg.RandomSecret))
	businessRepository := persistence.NewBusinessRepository(db)
	userRepository := persistence.NewUserRepository(db)

	authMiddleware := middleware.TokenAuth(tokenManager)

	businessDS := business.NewDomainService(businessRepository, tokenManager)
	promoDS := promocode.NewDomainService(promocodeRepository)
	userDS := user.NewDomainService(userRepository, tokenManager)

	businessAS := business2.NewApplicationService(businessDS, promoDS, cfg)
	userAS := user2.NewApplicationService(userDS, promoDS)

	businessAPI := http.NewBusinessAPI(businessAS)
	userAPI := http.NewUserAPI(userAS)

	api.Post("/business/auth/sign-up", businessAPI.SignUp) // 02
	api.Post("/business/auth/sign-in", businessAPI.SignIn) // 03

	api.Post("/business/promo/", authMiddleware, businessAPI.CreatePromoCode)   // 04
	api.Get("/business/promo", authMiddleware, businessAPI.GetPromoCodes)       // 05
	api.Get("/business/promo/:id", authMiddleware, businessAPI.GetPromoCode)    // 06
	api.Patch("/business/promo/:id", authMiddleware, businessAPI.EditPromoCode) // 06

	api.Post("/user/auth/sign-up", userAPI.SignUp)                  // 07
	api.Post("/user/auth/sign-in", userAPI.SignIn)                  // 08
	api.Get("/user/profile", authMiddleware, userAPI.GetProfile)    // 09
	api.Patch("/user/profile", authMiddleware, userAPI.EditProfile) // 09
	api.Get("/user/feed", authMiddleware, userAPI.GetFeed)          // 10

	api.Get("/user/promo/history", authMiddleware, userAPI.GetUseHistory)

	api.Get("/user/promo/:id", authMiddleware, userAPI.GetPromoCode) //10

	api.Post("/user/promo/:id/like", authMiddleware, userAPI.Like)     // 11
	api.Delete("/user/promo/:id/like", authMiddleware, userAPI.Unlike) // 11

	api.Get("/user/promo/:id/comments/:comment_id", authMiddleware, userAPI.GetPromoCodeComment)       // 12
	api.Put("/user/promo/:id/comments/:comment_id", authMiddleware, userAPI.EditPromoCodeComment)      // 12
	api.Delete("/user/promo/:id/comments/:comment_id", authMiddleware, userAPI.DeletePromoCodeComment) // 12
	api.Post("/user/promo/:id/comments", authMiddleware, userAPI.CommentPromoCode)                     // 12
	api.Get("/user/promo/:id/comments", authMiddleware, userAPI.GetPromoCodeComments)                  // 12

	antifraud := middleware.AntiFraud(cfg.RedisHost+":"+cfg.RedisPort, cfg.AntifraudAddress)

	api.Post("/user/promo/:id/activate", authMiddleware, antifraud, userAPI.ActivatePromoCode)
	// 13 POST user/promo/{id}/activate
	// 13 GET /user/promo/history

	api.Get("/business/promo/:id/stat", authMiddleware, businessAPI.UsageStatistic) // 14
	log.Info(server.Listen(":" + cfg.ServerPort))
}
