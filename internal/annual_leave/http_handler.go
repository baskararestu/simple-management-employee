package annualleave

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/middleware/authentication"
	"simple-management-employee/internal/middleware/validation"
	"simple-management-employee/internal/utilities"
	"simple-management-employee/pkg/xlogger"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpAnnualLeaveHandler struct {
	annualLeaveService domain.AnnualLeaveService
	jwtSecret        string
}

func NewHttpAnnualLeaveHandler(r fiber.Router, annualLeaveService domain.AnnualLeaveService,jwtSecret string) {
	h := &HttpAnnualLeaveHandler{
		annualLeaveService: annualLeaveService,
		jwtSecret: jwtSecret,
	}
	annualR := r.Group("/")
	annualR.Use(authentication.JwtProtect(h.jwtSecret))
	annualR.Get("/", h.GetAllAnnualLeaves)
	annualR.Post("/", validation.New[domain.CreateAnnualLeaveRequest](),h.CreateAnnualLeave)
	annualR.Get("/:id", h.GetAnnualLeaveByID)
	annualR.Put("/:id", validation.New[domain.CreateAnnualLeaveRequest](), h.UpdateAnnualLeave)
	annualR.Delete("/:id", h.DeleteAnnualLeave)
	annualR.Get("/status/:status", h.GetAnnualLeavesByStatus)
	annualR.Get("/date-range", h.GetAnnualLeavesByDateRange)
}

func (h *HttpAnnualLeaveHandler) GetAllAnnualLeaves(c *fiber.Ctx) error {
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
	startDate, _ := time.Parse("2006-01-02", query)
	endDate, _ := time.Parse("2006-01-02", query)
	filter := &domain.AnnualLeave{
		Status: query,
		User: domain.User{
			Email: query,
			FirstName: query,
			LastName: query,
		},
		StartDate: startDate,
		EndDate: endDate,
	}

	
	annualLeaves, nextPage, err := h.annualLeaveService.FindAll(uint(page), uint(size), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if annualLeaves == nil {
		return c.JSON([]domain.AnnualLeave{})
	}

	totalItem, err := h.annualLeaveService.Count(filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	maxPage := int(totalItem)/size


	if nextPage > 0 && nextPage <= uint(maxPage) {
		c.Set("X-Cursor", strconv.Itoa(int(nextPage)))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))

	return c.JSON(annualLeaves)
}

func (h *HttpAnnualLeaveHandler) CreateAnnualLeave(c *fiber.Ctx) error {
	annualLeaveReq := utilities.ExtractStructFromValidator[domain.CreateAnnualLeaveRequest](c)

	startDate, _ := time.Parse("2006-01-02", annualLeaveReq.StartDate)
	endDate, _ := time.Parse("2006-01-02", annualLeaveReq.EndDate)
	annualLeave := &domain.AnnualLeave{
		ID: uuid.NewString()	,
		UserID:    annualLeaveReq.UserID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    annualLeaveReq.Reason,
	}

	if err := h.annualLeaveService.Create(annualLeave); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Message{
		Code:    fiber.StatusCreated,
		Message: "Annual leave created successfully",
	})
}

func (h *HttpAnnualLeaveHandler) GetAnnualLeaveByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	annualLeave, err := h.annualLeaveService.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if annualLeave == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "annual leave not found",
		})
	}

	return c.JSON(annualLeave)
}

func (h *HttpAnnualLeaveHandler) UpdateAnnualLeave(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	annualLeaveReq := utilities.ExtractStructFromValidator[domain.CreateAnnualLeaveRequest](c)

	startDate, _ := time.Parse("2006-01-02", annualLeaveReq.StartDate)
	endDate, _ := time.Parse("2006-01-02", annualLeaveReq.EndDate)
	annualLeave := &domain.AnnualLeave{
		ID:        id,
		UserID:    annualLeaveReq.UserID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    annualLeaveReq.Reason,
		Status: annualLeaveReq.Status,
	}

	if err := h.annualLeaveService.Update(annualLeave); err != nil {
		xlogger.Logger.Error().Msgf("error updating annual leave: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Message{
		Code:    fiber.StatusOK,
		Message: "Annual leave updated successfully",
	})
}

func (h *HttpAnnualLeaveHandler) DeleteAnnualLeave(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	if err := h.annualLeaveService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *HttpAnnualLeaveHandler) GetAnnualLeavesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "status is required",
		})
	}

	annualLeaves, err := h.annualLeaveService.FindByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if annualLeaves == nil {
		return c.JSON([]domain.AnnualLeave{})
	}

	return c.JSON(annualLeaves)
}

func (h *HttpAnnualLeaveHandler) GetAnnualLeavesByDateRange(c *fiber.Ctx) error {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	if startDate == "" || endDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "startDate and endDate are required",
		})
	}

	annualLeaves, err := h.annualLeaveService.FindByDateRange(startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if annualLeaves == nil {
		return c.JSON([]domain.AnnualLeave{})
	}

	return c.JSON(annualLeaves)
}