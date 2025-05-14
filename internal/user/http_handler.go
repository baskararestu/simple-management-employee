package user

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/middleware/validation"
	"simple-management-employee/internal/utilities"
	"simple-management-employee/pkg/xlogger"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	userSvc domain.UserService
	roleSvc domain.RoleService
	jwtSecret string
}

func NewHttpUserHandler(r fiber.Router, userSvc domain.UserService, roleSvc domain.RoleService, jwtSecret string){
	h:=&HttpUserHandler{
		userSvc: userSvc,
		roleSvc: roleSvc,
		jwtSecret: jwtSecret,
	}

	r.Get("/", h.GetAllUsers)
	r.Get("/:id", h.GetUserById)
	r.Put("/:id", validation.New[domain.UpdateUserRequest](),h.UpdateUser)
	r.Delete("/:id", h.DeleteUser)
	r.Get("/role/name/:roleName", h.FindByRoleName)
	r.Get("/role/id/:roleId", h.FindByRoleID)
}

func (h *HttpUserHandler) GetAllUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	sizeStr := c.Query("size", "10")

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "page must be a positive number",
		})
	}

	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil || sizeInt <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "size must be a positive number",
		})
	}

	page := uint(pageInt)
	size := uint(sizeInt)
	
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")
	email := c.Query("email")

	filter := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	users, err := h.userSvc.FindAll(page, size, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}

func (h *HttpUserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	user, err := h.userSvc.FindByID(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *HttpUserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	// Extract token
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Authorization header missing",
		})
	}

	var token string
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token format",
		})
	}

	claims, err := utilities.ExtractClaimsFromToken(token, h.jwtSecret) // use env/config
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token: " + err.Error(),
		})
	}

	// Get role name
	xlogger.Logger.Info().Msgf("Claims: %v", claims)
	role, err := h.roleSvc.FindByID(claims.RoleID)
	if err != nil {
		xlogger.Logger.Error().Err(err).Msgf("error finding role: %+v", err)
		return err
	}

	// Authorization logic
	roleName := strings.ToLower(role.Name)
	if roleName != "admin" && claims.UserID != id {
		return c.Status(fiber.StatusForbidden).JSON(domain.Error{
			Code:    fiber.StatusForbidden,
			Message: "You don't have permission to update this user",
		})
	}

	// Parse request
	userReq := utilities.ExtractStructFromValidator[domain.UpdateUserRequest](c)

	checkUser, err := h.userSvc.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	user := &domain.User{
		ID:          id,
		FirstName:   userReq.FirstName,
		LastName:    userReq.LastName,
		Email:       userReq.Email,
		Password:    checkUser.Password,
		Address:     userReq.Address,
		PhoneNumber: userReq.PhoneNumber,
		Gender:      userReq.Gender,
	}
	
	if roleName == "admin" && claims.UserID != id {
		user.RoleID = checkUser.RoleID
	} else {
		user.RoleID = claims.RoleID
	}

	if err := h.userSvc.Update(user); err != nil {
		return err
	}

	return c.JSON(domain.Message{
		Code:    fiber.StatusOK,
		Message: "success update user",
	})
}


func (h *HttpUserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	xlogger.Logger.Info().Msgf("Deleting user with ID: %s", id)

	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Authorization header missing",
		})
	}

	var token string
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token format",
		})
	}

	claims, err := utilities.ExtractClaimsFromToken(token, h.jwtSecret) // use env/config
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token: " + err.Error(),
		})
	}

	// Get role name from role ID
	role, err := h.roleSvc.FindByID(claims.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch role: " + err.Error(),
		})
	}

	roleName := strings.ToLower(role.Name)
	if roleName != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(domain.Error{
			Code:    fiber.StatusForbidden,
			Message: "Only admin can delete users",
		})
	}

	// Proceed with deletion
	if err := h.userSvc.Delete(id); err != nil {
		xlogger.Logger.Error().Err(err).Msgf("error deleting user: %+v", err)
		return err
	}

	return c.JSON(domain.Message{
		Code:    fiber.StatusOK,
		Message: "success delete user",
	})
}


func (h *HttpUserHandler) FindByRoleName(c *fiber.Ctx) error {
	roleName := c.Params("roleName")
	if roleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "role name is required",
		})
	}

	users, err := h.userSvc.FindByRoleName(roleName)
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *HttpUserHandler) FindByRoleID(c *fiber.Ctx) error {
	roleID := c.Params("roleId")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "role id is required",
		})
	}

	users, err := h.userSvc.FindByRoleID(roleID)
	if err != nil {
		return err
	}

	return c.JSON(users)
}

