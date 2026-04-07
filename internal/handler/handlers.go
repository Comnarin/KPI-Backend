package handler

import "github.com/gofiber/fiber/v2"

// tenantFromCtx extracts tenantId set by the Auth middleware.
// Panics on misconfiguration rather than silently returning wrong data.
func tenantFromCtx(c *fiber.Ctx) string {
	v, _ := c.Locals("tenantId").(string)
	return v
}

// roleFromCtx extracts role set by the Auth middleware.
func roleFromCtx(c *fiber.Ctx) string {
	v, _ := c.Locals("role").(string)
	return v
}

// SafeUser is the public-facing user DTO — password field is never included.
type SafeUser struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"fullName"`
	Role         string `json:"role"`
	TenantID     string `json:"tenantId"`
	DepartmentID string `json:"departmentId"`
}

type Handlers struct {
	Evaluation     *EvaluationHandler
	Template       *TemplateHandler
	Employee       *EmployeeHandler
	User           *UserHandler
	Department     *DepartmentHandler
	SalarySettings *SalarySettingsHandler
	Permissions    *PermissionsHandler
	Admin          *AdminHandler
	Auth           *AuthHandler
	Tenant         *TenantHandler
	Seed           *SeedHandler
	Dashboard      *DashboardHandler
	Adjustment     *SalaryAdjustmentHandler
}
