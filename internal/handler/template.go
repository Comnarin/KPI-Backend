package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type TemplateHandler struct {
	templateUC domain.TemplateUseCase
}

func NewTemplateHandler(templateUC domain.TemplateUseCase) *TemplateHandler {
	return &TemplateHandler{templateUC: templateUC}
}

func (h *TemplateHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	filter := domain.TemplateFilter{
		TenantID:    tenantID,
		SearchQuery: c.Query("q"),
	}

	templates, err := h.templateUC.ListTemplates(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(templates)
}

func (h *TemplateHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	userID := c.Locals("userId").(string)
	var body domain.EvaluationTemplate
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.TenantID = tenantID
	body.CreatedByID = userID
	tpl, err := h.templateUC.CreateTemplate(c.Context(), body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(tpl)
}

func (h *TemplateHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	var body domain.EvaluationTemplate
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	body.ID = id
	body.TenantID = tenantID
	tpl, err := h.templateUC.UpdateTemplate(c.Context(), body)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "template not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(tpl)
}

func (h *TemplateHandler) Delete(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	if err := h.templateUC.DeleteTemplate(c.Context(), id, tenantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
