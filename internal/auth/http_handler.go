package auth

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/middleware/authentication"
	"simple-management-employee/internal/middleware/validation"
	"simple-management-employee/internal/utilities"
	"simple-management-employee/pkg/xlogger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpAuthHandler struct {
	authService domain.AuthService
	jwtSecret string
}

func NewHttpAuthHandler(r fiber.Router, authService domain.AuthService, jwtSecret string) {
	handler := &HttpAuthHandler{
		authService: authService,
		jwtSecret: jwtSecret,
	}

	r.Post("/login",validation.New[domain.LoginRequest](),handler.Login)
	adminGroup := r.Group("/register")
	adminGroup.Use(authentication.JwtProtect(handler.jwtSecret))
	adminGroup.Post("/admin",  validation.New[domain.RegisterAdminRequest](),handler.RegisterAdmin)
	adminGroup.Post("/employee", validation.New[domain.RegisterEmployeeRequest](),handler.RegisterEmployee)
}

func (h *HttpAuthHandler) Login(c *fiber.Ctx) error {
	req := utilities.ExtractStructFromValidator[domain.LoginRequest](c)
	authReq := &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
xlogger.Logger.Info().Msgf("Login request: %v", authReq)
	loginResponse, err := h.authService.Login(authReq)
	if err != nil {
		xlogger.Logger.Error().Err(err).Msgf("Login failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid credentials",
		})
	}

	return c.JSON(loginResponse)
}

func (h *HttpAuthHandler) RegisterAdmin(c *fiber.Ctx) error {
	req := utilities.ExtractStructFromValidator[domain.RegisterAdminRequest](c)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	authReq := &domain.User{
		ID: 	 uuid.NewString(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		RoleID:    req.RoleID,
	}

	if err := h.authService.RegisterAdmin(authReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to register admin",
		})
	}

	return c.JSON(domain.Message{
		Code:    fiber.StatusCreated,
		Message: "Admin registered successfully",
	})
}

func (h *HttpAuthHandler) RegisterEmployee(c *fiber.Ctx) error {
	req := utilities.ExtractStructFromValidator[domain.RegisterEmployeeRequest](c)

	authReq := &domain.User{
		ID: 	 uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		Address:     &req.Address,
		PhoneNumber: &req.PhoneNumber,
		Gender: &req.Gender,
		RoleID:      req.RoleID,
	}
	if err := h.authService.RegisterEmployee(authReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to register employee",
		})
	}

	return c.JSON(domain.Message{
		Code:    fiber.StatusCreated,
		Message: "Employee registered successfully",
	})
}