package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type SalarySettingsHandler struct {
	salaryUC domain.SalaryFormulaUseCase
}

func NewSalarySettingsHandler(salaryUC domain.SalaryFormulaUseCase) *SalarySettingsHandler {
	return &SalarySettingsHandler{salaryUC: salaryUC}
}

func (h *SalarySettingsHandler) Get(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	formula, err := h.salaryUC.GetFormula(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(formula)
}

func (h *SalarySettingsHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	var body domain.SalaryFormula
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.TenantID = tenantID
	formula, err := h.salaryUC.UpdateFormula(c.Context(), body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(formula)
}
