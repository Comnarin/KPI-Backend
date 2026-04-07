package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type EmployeeHandler struct {
	employeeUC domain.EmployeeUseCase
}

func NewEmployeeHandler(employeeUC domain.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{employeeUC: employeeUC}
}

func (h *EmployeeHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	filter := domain.EmployeeFilter{
		TenantID:               tenantID,
		Status:                 c.Query("status"),
		SearchQuery:            c.Query("q"),
		ExcludeEvaluatedPeriod: c.Query("excludeEvaluatedPeriod"),
		EvaluatorID:            c.Query("evaluatorId"),
	}
	employees, err := h.employeeUC.ListEmployees(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(employees)
}

func (h *EmployeeHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	var body domain.Employee
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.TenantID = tenantID
	emp, err := h.employeeUC.CreateEmployee(c.Context(), body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(emp)
}

func (h *EmployeeHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	var body domain.Employee
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.ID = id
	body.TenantID = tenantID
	emp, err := h.employeeUC.UpdateEmployee(c.Context(), body)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "employee not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(emp)
}

func (h *EmployeeHandler) Delete(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	if err := h.employeeUC.DeleteEmployee(c.Context(), id, tenantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
