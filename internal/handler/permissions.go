package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PermissionsHandler struct {
	permUC domain.PermissionUseCase
}

func NewPermissionsHandler(permUC domain.PermissionUseCase) *PermissionsHandler {
	return &PermissionsHandler{permUC: permUC}
}

func (h *PermissionsHandler) GetForTenant(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)

	perms, err := h.permUC.GetPermissions(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var filtered []domain.RolePermission
	for _, p := range perms {
		if string(p.Role) == role {
			filtered = append(filtered, p)
		}
	}

	return c.JSON(filtered)
}

// GetAllForTenant returns all role permissions for the tenant.
// Used by the Permissions management page (manage_settings permission required).
func (h *PermissionsHandler) GetAllForTenant(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	perms, err := h.permUC.GetPermissions(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(perms)
}

type updatePermsRequest struct {
	Permissions []domain.RolePermission `json:"permissions"`
}

func (h *PermissionsHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	role := roleFromCtx(c)

	// Tenant-level update: CEO and HR can configure their own tenant's permissions.
	// SuperAdmin can configure any tenant's permissions.
	allowedRoles := map[string]bool{
		string(domain.RoleSuperAdmin): true,
		string(domain.RoleCEO):        true,
		string(domain.RoleHR):         true,
	}
	if !allowedRoles[role] {
		return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
	}

	var body updatePermsRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.permUC.UpdatePermissions(c.Context(), tenantID, body.Permissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PermissionsHandler) GetByTenantParam(c *fiber.Ctx) error {
	role := roleFromCtx(c)
	if role != string(domain.RoleSuperAdmin) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}

	tenantID := c.Params("id")
	perms, err := h.permUC.GetPermissions(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(perms)
}

func (h *PermissionsHandler) UpdateByTenantParam(c *fiber.Ctx) error {
	role := roleFromCtx(c)
	if role != string(domain.RoleSuperAdmin) {
		return fiber.NewError(fiber.StatusForbidden, "superadmin only")
	}

	tenantID := c.Params("id")
	var body updatePermsRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.permUC.UpdatePermissions(c.Context(), tenantID, body.Permissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
