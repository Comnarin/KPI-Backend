package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type DepartmentHandler struct {
	deptUC domain.DepartmentUseCase
}

func NewDepartmentHandler(deptUC domain.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{deptUC: deptUC}
}

func (h *DepartmentHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	filter := domain.DepartmentFilter{
		TenantID:    tenantID,
		SearchQuery: c.Query("q"),
	}
	depts, err := h.deptUC.ListDepartments(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(depts)
}

func (h *DepartmentHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	var body domain.Department
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.TenantID = tenantID
	dept, err := h.deptUC.CreateDepartment(c.Context(), body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(dept)
}

func (h *DepartmentHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	var body domain.Department
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.ID = id
	body.TenantID = tenantID
	dept, err := h.deptUC.UpdateDepartment(c.Context(), body)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "department not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(dept)
}

func (h *DepartmentHandler) Delete(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	if err := h.deptUC.DeleteDepartment(c.Context(), id, tenantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
