package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	tenantUC domain.TenantUseCase
}

func NewAdminHandler(tenantUC domain.TenantUseCase) *AdminHandler {
	return &AdminHandler{tenantUC: tenantUC}
}

func requireSuperAdmin(c *fiber.Ctx) bool {
	role, ok := c.Locals("role").(string)
	return ok && role == string(domain.RoleSuperAdmin)
}

func (h *AdminHandler) ListTenants(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}
	tenants, err := h.tenantUC.ListTenants(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(tenants)
}

type createTenantWithAdminRequest struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Size          string `json:"size"`
	MaxUsers      int    `json:"maxUsers"`
	AdminEmail    string `json:"adminEmail"`
	AdminPassword string `json:"adminPassword"`
	AdminName     string `json:"adminName"`
}

func (h *AdminHandler) CreateTenant(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}

	var body createTenantWithAdminRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	tenant := domain.Tenant{
		Name:     body.Name,
		Code:     body.Code,
		Size:     body.Size,
		MaxUsers: body.MaxUsers,
	}

	// Set default MaxUsers based on org size if not provided
	if tenant.MaxUsers <= 0 {
		tenant.MaxUsers = domain.DefaultMaxUsers(body.Size)
	}

	result, err := h.tenantUC.CreateTenant(c.Context(), tenant, body.AdminEmail, body.AdminPassword, body.AdminName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *AdminHandler) GetTenant(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}
	id := c.Params("id")
	tenant, err := h.tenantUC.GetTenant(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "tenant not found")
	}
	return c.JSON(tenant)
}

type updateTenantRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Size     string `json:"size"`
	MaxUsers int    `json:"maxUsers"`
	IsActive bool   `json:"isActive"`
}

func (h *AdminHandler) UpdateTenant(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}
	id := c.Params("id")
	var body updateTenantRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	tenant, err := h.tenantUC.UpdateTenant(c.Context(), id, body.Name, body.Code, body.Size, body.IsActive)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "tenant not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(tenant)
}

func (h *AdminHandler) DeleteTenant(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}
	id := c.Params("id")
	if err := h.tenantUC.DeleteTenant(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

type updateTenantConfigAdminReq struct {
	EnableHREval       bool `json:"enableHrEval"`
	EnableDeptHeadEval bool `json:"enableDeptHeadEval"`
	EnableCEOEval      bool `json:"enableCeoEval"`
}

func (h *AdminHandler) UpdateTenantConfig(c *fiber.Ctx) error {
	if !requireSuperAdmin(c) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}
	id := c.Params("id")
	var body updateTenantConfigAdminReq
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	cfg := domain.TenantConfig{
		EnableHREval:       body.EnableHREval,
		EnableDeptHeadEval: body.EnableDeptHeadEval,
		EnableCEOEval:      body.EnableCEOEval,
	}
	result, err := h.tenantUC.UpdateTenantConfig(c.Context(), id, cfg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(result)
}
