package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type EvaluationHandler struct {
	evalUC domain.EvaluationUseCase
}

func NewEvaluationHandler(evalUC domain.EvaluationUseCase) *EvaluationHandler {
	return &EvaluationHandler{evalUC: evalUC}
}

func (h *EvaluationHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	userID, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	filter := domain.EvaluationFilter{
		TenantID:    tenantID,
		ViewerID:    userID,
		ViewerRole:  role,
		SearchQuery:  c.Query("q"),
		DepartmentID: c.Query("departmentId"),
		Period:       c.Query("period"),
	}

	results, err := h.evalUC.ListEvaluations(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(results)
}

func (h *EvaluationHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	evaluatorID, _ := c.Locals("userId").(string)
	evaluatorRole, _ := c.Locals("role").(string)
	evaluatorName, _ := c.Locals("fullName").(string)
	var req domain.EvaluationResult
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	// Inject evaluator attribution
	req.EvaluatorID = evaluatorID
	req.EvaluatorName = evaluatorName
	req.EvaluatorRole = evaluatorRole
	res, err := h.evalUC.CreateEvaluation(c.Context(), tenantID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *EvaluationHandler) Delete(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	if err := h.evalUC.DeleteEvaluation(c.Context(), id, tenantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
