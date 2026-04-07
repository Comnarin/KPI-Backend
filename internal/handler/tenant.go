package handler

import (
	"kpi-backend/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TenantHandler struct {
	tenantUC domain.TenantUseCase
	permUC   domain.PermissionUseCase
}

func NewTenantHandler(tenantUC domain.TenantUseCase, permUC domain.PermissionUseCase) *TenantHandler {
	return &TenantHandler{tenantUC: tenantUC, permUC: permUC}
}

type updateConfigRequest struct {
	EnableHREval       *bool `json:"enableHrEval"`
	EnableDeptHeadEval *bool `json:"enableDeptHeadEval"`
	EnableCEOEval      *bool `json:"enableCeoEval"`
}

func (h *TenantHandler) UpdateConfig(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	if tenantID == "" {
		tenantID = tenantFromCtx(c)
	}

	var body updateConfigRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	cfg := domain.TenantConfig{}
	if body.EnableHREval != nil {
		cfg.EnableHREval = *body.EnableHREval
	}
	if body.EnableDeptHeadEval != nil {
		cfg.EnableDeptHeadEval = *body.EnableDeptHeadEval
	}
	if body.EnableCEOEval != nil {
		cfg.EnableCEOEval = *body.EnableCEOEval
	}

	res, err := h.tenantUC.UpdateTenantConfig(c.Context(), tenantID, cfg)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "config not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(res)
}

// ─── Period Management Endpoints ───

func (h *TenantHandler) ListPeriods(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	periods, err := h.tenantUC.ListPeriods(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(periods)
}

type periodRequest struct {
	Label     string  `json:"label"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

func parseDateOpt(d *string) (*time.Time, error) {
	if d != nil && *d != "" {
		t, err := time.Parse("2006-01-02", *d)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, nil
}

func (h *TenantHandler) CreatePeriod(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)
	if allowed, _ := h.permUC.HasPermission(c.Context(), tenantID, role, domain.PermManagePeriod); !allowed {
		return fiber.NewError(fiber.StatusForbidden, "คุณไม่มีสิทธิ์จัดการรอบการประเมิน")
	}

	var body periodRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	start, err := parseDateOpt(body.StartDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid start date format")
	}
	end, err := parseDateOpt(body.EndDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid end date format")
	}

	period, err := h.tenantUC.CreatePeriod(c.Context(), tenantID, body.Label, start, end)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(period)
}

func (h *TenantHandler) UpdatePeriod(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)
	if allowed, _ := h.permUC.HasPermission(c.Context(), tenantID, role, domain.PermManagePeriod); !allowed {
		return fiber.NewError(fiber.StatusForbidden, "คุณไม่มีสิทธิ์จัดการรอบการประเมิน")
	}

	id := c.Params("id")
	var body periodRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	start, err := parseDateOpt(body.StartDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid start date format")
	}
	end, err := parseDateOpt(body.EndDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid end date format")
	}

	period, err := h.tenantUC.UpdatePeriod(c.Context(), tenantID, id, body.Label, start, end)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "period not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(period)
}

func (h *TenantHandler) DeletePeriod(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)
	if allowed, _ := h.permUC.HasPermission(c.Context(), tenantID, role, domain.PermManagePeriod); !allowed {
		return fiber.NewError(fiber.StatusForbidden, "คุณไม่มีสิทธิ์จัดการรอบการประเมิน")
	}

	id := c.Params("id")
	if err := h.tenantUC.DeletePeriod(c.Context(), tenantID, id); err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "period not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TenantHandler) SetActivePeriod(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)
	if allowed, _ := h.permUC.HasPermission(c.Context(), tenantID, role, domain.PermManagePeriod); !allowed {
		return fiber.NewError(fiber.StatusForbidden, "คุณไม่มีสิทธิ์จัดการรอบการประเมิน")
	}

	id := c.Params("id")
	if err := h.tenantUC.SetActivePeriod(c.Context(), tenantID, id); err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "period not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
