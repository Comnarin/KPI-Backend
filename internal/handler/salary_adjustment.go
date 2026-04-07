package handler

import (
	"kpi-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type SalaryAdjustmentHandler struct {
	Usecase domain.AdjustmentUseCase
}

func NewSalaryAdjustmentHandler(uc domain.AdjustmentUseCase) *SalaryAdjustmentHandler {
	return &SalaryAdjustmentHandler{Usecase: uc}
}

func (h *SalaryAdjustmentHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	adjustments, err := h.Usecase.ListAdjustments(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(adjustments)
}

func (h *SalaryAdjustmentHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	var adj domain.SalaryAdjustment
	if err := c.BodyParser(&adj); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	adj.TenantID = tenantID
	res, err := h.Usecase.CreateAdjustment(c.Context(), adj)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
