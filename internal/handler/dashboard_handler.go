package handler

import (
	"log"
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	Usecase domain.DashboardUseCase
}

func NewDashboardHandler(uc domain.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{Usecase: uc}
}

func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	if tenantID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	period := c.Query("period")
	log.Printf("GetDashboardStats called for tenantID: %s, periodFilter: %q", tenantID, period)
	stats, err := h.Usecase.GetDashboardStats(c.Context(), tenantID, period)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(stats)
}

