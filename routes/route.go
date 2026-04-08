package routes

import (
	"kpi-backend/internal/handler"
	"kpi-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handlers *handler.Handlers) {
	// ─── Public routes ────────────────────────────────────────────────────────
	api := app.Group("/api")
	api.Post("/auth/login", handlers.Auth.Login)
	api.Post("/dev/seed", handlers.Seed.Seed)

	// ─── Protected routes ─────────────────────────────────────────────────────
	protected := api.Group("/", middleware.Auth())

	// Employees
	protected.Get("/employees", handlers.Employee.List)
	protected.Post("/employees", handlers.Employee.Create)
	protected.Put("/employees/:id", handlers.Employee.Update)
	protected.Delete("/employees/:id", handlers.Employee.Delete)

	// Users (accounts with system access)
	protected.Get("/users", handlers.User.List)
	protected.Post("/users", handlers.User.Create)
	protected.Put("/users/:id", handlers.User.Update)
	protected.Delete("/users/:id", handlers.User.Delete)

	// Departments
	protected.Get("/departments", handlers.Department.List)
	protected.Post("/departments", handlers.Department.Create)
	protected.Put("/departments/:id", handlers.Department.Update)
	protected.Delete("/departments/:id", handlers.Department.Delete)

	// KPI Templates
	protected.Get("/templates", handlers.Template.List)
	protected.Post("/templates", handlers.Template.Create)
	protected.Put("/templates/:id", handlers.Template.Update)
	protected.Delete("/templates/:id", handlers.Template.Delete)

	// Evaluations
	protected.Get("/evaluations", handlers.Evaluation.List)
	protected.Post("/evaluations", handlers.Evaluation.Create)
	protected.Delete("/evaluations/:id", handlers.Evaluation.Delete)

	// Dashboard Stats
	protected.Get("/dashboard/stats", handlers.Dashboard.GetStats)

	// Salary Settings (per-tenant formula)
	protected.Get("/salary-settings", handlers.SalarySettings.Get)
	protected.Put("/salary-settings", handlers.SalarySettings.Update)

	// Salary Adjustments (History/Approvals)
	protected.Get("/salary-adjustments", handlers.Adjustment.List)
	protected.Post("/salary-adjustments", handlers.Adjustment.Create)

	// Role Permissions
	protected.Get("/permissions", handlers.Permissions.GetForTenant)
	protected.Get("/permissions/all", handlers.Permissions.GetAllForTenant)
	protected.Put("/permissions", handlers.Permissions.Update)

	// Tenant Config (existing)
	protected.Patch("/tenants/:id/config", handlers.Tenant.UpdateConfig)

	// Period Configuration — requires manage_period permission
	protected.Get("/periods", handlers.Tenant.ListPeriods)
	protected.Post("/periods", handlers.Tenant.CreatePeriod)
	protected.Put("/periods/:id", handlers.Tenant.UpdatePeriod)
	protected.Put("/periods/:id/active", handlers.Tenant.SetActivePeriod)
	protected.Delete("/periods/:id", handlers.Tenant.DeletePeriod)

	// ─── Super Admin routes ────────────────────────────────────────────────────
	protected.Get("/admin/tenants", handlers.Admin.ListTenants)
	protected.Post("/admin/tenants", handlers.Admin.CreateTenant)
	protected.Get("/admin/tenants/:id", handlers.Admin.GetTenant)
	protected.Put("/admin/tenants/:id", handlers.Admin.UpdateTenant)
	protected.Delete("/admin/tenants/:id", handlers.Admin.DeleteTenant)
	protected.Get("/admin/tenants/:id/config", handlers.Admin.GetTenantConfig)
	protected.Put("/admin/tenants/:id/config", handlers.Admin.UpdateTenantConfig)
	protected.Get("/admin/tenants/:id/permissions", handlers.Permissions.GetByTenantParam)
	protected.Put("/admin/tenants/:id/permissions", handlers.Permissions.UpdateByTenantParam)
}
