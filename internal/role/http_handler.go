package role

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/middleware/validation"
	"simple-management-employee/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HttpRoleHandler struct {
	roleService domain.RoleService
}

func NewHttpRoleHandler(r fiber.Router, roleService domain.RoleService) {
	h := &HttpRoleHandler{
		roleService: roleService,
	}
	
	r.Get("/", h.GetAllRoles)
	r.Post("/", validation.New[domain.CreateRoleRequest](), h.CreateRole)
	r.Get("/:id", h.GetRoleByID)
	r.Put("/:id", validation.New[domain.UpdateRoleRequest](), h.UpdateRole)
	r.Delete("/:id", h.DeleteRole)
	r.Get("/name/:name", h.GetRoleByName)
}

func (h *HttpRoleHandler) GetAllRoles(c *fiber.Ctx) error {
	page, size, query := c.QueryInt("page", 1), c.QueryInt("size", 10), c.Query("q")

	if page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "page must be a positive integer",
		})
	}
	if size <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "size must be a positive integer",
		})
	}

	filter := &domain.Role{
		Name: query,
	}

	roles, nextPage, err := h.roleService.FindAll(uint(page), uint(size), filter)
	if err != nil {
		return err
	}

	if roles == nil{
		c.JSON([]domain.Role{})
	}

	totalItem, err := h.roleService.Count(filter)
	if err != nil {
		return err
	}

	maxPage := int(totalItem)/size


	if nextPage > 0 && nextPage <= uint(maxPage) {
		c.Set("X-Cursor", strconv.Itoa(int(nextPage)))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))

	return c.JSON(roles)
}

func (h *HttpRoleHandler) CreateRole(c *fiber.Ctx) error {
	roleReq := utilities.ExtractStructFromValidator[domain.CreateRoleRequest](c)

	role := &domain.Role{
		Name: roleReq.Name,
	}

	if err := h.roleService.Create(roleReq); err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(role)
}

func (h *HttpRoleHandler) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	role, err := h.roleService.FindByID(id)
	if err != nil {
		return err
	}

	return c.JSON(role)
}

func (h *HttpRoleHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	roleReq := utilities.ExtractStructFromValidator[domain.UpdateRoleRequest](c)
	
	role := &domain.Role{
		ID:   roleReq.ID,
		Name: roleReq.Name,
	}

	if err := h.roleService.Update(role); err != nil {
		return err
	}

	return c.JSON(role)
}

func (h *HttpRoleHandler) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	if err := h.roleService.Delete(id); err != nil {
		return err
	}

	return c.JSON(domain.Message{
		Code:   fiber.StatusOK,
		Message: "success delete role",
	})
}

func (h *HttpRoleHandler) GetRoleByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "name is required",
		})
	}

	role, err := h.roleService.FindByName(name)
	if err != nil {
		return err
	}

	return c.JSON(role)
}