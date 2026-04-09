package registry

import (
	"kpi-backend/internal/handler"
	"kpi-backend/internal/repository"
	"kpi-backend/internal/usecase"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db}
}

func (r *Registry) NewAppHandlers() *handler.Handlers {
	// Repositories
	tenantRepo := repository.NewTenantRepository(r.db)
	userRepo := repository.NewUserRepository(r.db)
	evalRepo := repository.NewEvaluationRepository(r.db)
	empRepo := repository.NewEmployeeRepository(r.db)
	tplRepo := repository.NewTemplateRepository(r.db)
	deptRepo := repository.NewDepartmentRepository(r.db)
	salaryRepo := repository.NewSalaryFormulaRepository(r.db)
	permRepo := repository.NewPermissionRepository(r.db)
	periodRepo := repository.NewPeriodRepository(r.db)
	dashRepo := repository.NewDashboardRepository(r.db)
	adjRepo := repository.NewAdjustmentRepository(r.db)

	// UseCases
	workflowUC := usecase.NewWorkflowUseCase(tenantRepo)
	authUC := usecase.NewAuthUseCase(userRepo, tenantRepo)
	evalUC := usecase.NewEvaluationUseCase(evalRepo, workflowUC, empRepo, tplRepo, periodRepo, userRepo)
	tenantUC := usecase.NewTenantUseCase(tenantRepo, periodRepo, userRepo, permRepo)
	empUC := usecase.NewEmployeeUseCase(empRepo)
	tplUC := usecase.NewTemplateUseCase(tplRepo)
	deptUC := usecase.NewDepartmentUseCase(deptRepo, empRepo)
	salaryUC := usecase.NewSalaryFormulaUseCase(salaryRepo)
	permUC := usecase.NewPermissionUseCase(permRepo)
	dashUC := usecase.NewDashboardUseCase(dashRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	adjUC := usecase.NewAdjustmentUseCase(adjRepo)

	// Handlers
	return &handler.Handlers{
		Evaluation:     handler.NewEvaluationHandler(evalUC),
		Template:       handler.NewTemplateHandler(tplUC),
		Employee:       handler.NewEmployeeHandler(empUC),
		User:           handler.NewUserHandler(userUC),
		Department:     handler.NewDepartmentHandler(deptUC),
		SalarySettings: handler.NewSalarySettingsHandler(salaryUC),
		Permissions:    handler.NewPermissionsHandler(permUC),
		Admin:          handler.NewAdminHandler(tenantUC),
		Auth:           handler.NewAuthHandler(authUC),
		Tenant:         handler.NewTenantHandler(tenantUC, permUC),
		Seed:           handler.NewSeedHandler(r.db),
		Dashboard:      handler.NewDashboardHandler(dashUC),
		Adjustment:     handler.NewSalaryAdjustmentHandler(adjUC),
	}
}
