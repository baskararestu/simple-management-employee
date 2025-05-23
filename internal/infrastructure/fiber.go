package infrastructure

import (
	"fmt"
	annualleave "simple-management-employee/internal/annual_leave"
	"simple-management-employee/internal/auth"
	"simple-management-employee/internal/docs"
	"simple-management-employee/internal/role"
	"simple-management-employee/internal/user"
	"simple-management-employee/pkg/xlogger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Run() {
	logger := xlogger.Logger

	app := fiber.New(fiber.Config{
		ProxyHeader:           cfg.ProxyHeader,
		DisableStartupMessage: true,
		ErrorHandler:          defaultErrorHandler,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))
	app.Use(recover2.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	api := app.Group("/api")
	docs.NewHttpHandler(api.Group("/docs"))
	// article.NewHttpHandler(api.Group("/articles"), articleService)
	user.NewHttpUserHandler(api.Group("/users"), userSvc,roleSvc, cfg.JwtSecret)
	auth.NewHttpAuthHandler(api.Group("/auth"), authSvc, cfg.JwtSecret)
	role.NewHttpRoleHandler(api.Group("/roles"), roleSvc)
	annualleave.NewHttpAnnualLeaveHandler(api.Group("/annual-leaves"), annualLeaveSvc, cfg.JwtSecret)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
