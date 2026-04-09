package domain

import (
	"context"
	"time"
)

// ─── Repository Interfaces ───────────────────────────────────────────────────

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetByCode(ctx context.Context, code string) (*Tenant, error)
	List(ctx context.Context) ([]Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	GetConfig(ctx context.Context, tenantID string) (*TenantConfig, error)
	UpdateConfig(ctx context.Context, config *TenantConfig) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, filter UserFilter) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string, tenantID string) error
}

type EvaluationRepository interface {
	Create(ctx context.Context, evaluation *EvaluationResult) error
	List(ctx context.Context, filter EvaluationFilter) ([]EvaluationResponseDTO, error)
	Delete(ctx context.Context, id string, tenantID string) error
}

type TemplateRepository interface {
	List(ctx context.Context, filter TemplateFilter) ([]EvaluationTemplate, error)
	GetByID(ctx context.Context, id string) (*EvaluationTemplate, error)
	Create(ctx context.Context, template *EvaluationTemplate) error
	Update(ctx context.Context, template *EvaluationTemplate) error
	Delete(ctx context.Context, id string, tenantID string) error
}

type EmployeeRepository interface {
	List(ctx context.Context, filter EmployeeFilter) ([]Employee, error)
	GetByID(ctx context.Context, id string) (*Employee, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*Employee, error)
	Create(ctx context.Context, employee *Employee) error
	Update(ctx context.Context, employee *Employee) error
	UpdateFields(ctx context.Context, id string, tenantID string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string, tenantID string) error
}

type DepartmentRepository interface {
	List(ctx context.Context, filter DepartmentFilter) ([]Department, error)
	GetByID(ctx context.Context, id string) (*Department, error)
	Create(ctx context.Context, dept *Department) error
	Update(ctx context.Context, dept *Department) error
	Delete(ctx context.Context, id string, tenantID string) error
}

type SalaryFormulaRepository interface {
	GetByTenant(ctx context.Context, tenantID string) (*SalaryFormula, error)
	Upsert(ctx context.Context, formula *SalaryFormula) error
}

type AdjustmentRepository interface {
	Create(ctx context.Context, adj *SalaryAdjustment) error
	List(ctx context.Context, tenantID string) ([]SalaryAdjustment, error)
}

type PermissionRepository interface {
	GetForTenant(ctx context.Context, tenantID string) ([]RolePermission, error)
	BulkUpsert(ctx context.Context, perms []RolePermission) error
	HasPermission(ctx context.Context, tenantID string, role string, permKey string) (bool, error)
	InitDefaults(ctx context.Context, tenantID string) error
}

// ─── UseCase Interfaces ───────────────────────────────────────────────────────

type AuthUseCase interface {
	Login(ctx context.Context, tenantCode, email, password string) (string, *User, error)
}

type EvaluationUseCase interface {
	ListEvaluations(ctx context.Context, filter EvaluationFilter) ([]EvaluationResponseDTO, error)
	CreateEvaluation(ctx context.Context, tenantID string, req EvaluationResult) (*EvaluationResult, error)
	DeleteEvaluation(ctx context.Context, id string, tenantID string) error
}

type TenantUseCase interface {
	CreateTenant(ctx context.Context, tenant Tenant, adminEmail, adminPassword, adminName string) (*Tenant, error)
	GetTenant(ctx context.Context, id string) (*Tenant, error)
	ListTenants(ctx context.Context) ([]Tenant, error)
	UpdateTenant(ctx context.Context, id string, name, code, size string, maxUsers int, isActive bool) (*Tenant, error)
	DeleteTenant(ctx context.Context, id string) error
	GetTenantConfig(ctx context.Context, tenantID string) (*TenantConfig, error)
	UpdateTenantConfig(ctx context.Context, tenantID string, config TenantConfig) (*TenantConfig, error)
	// Period management
	GetActivePeriod(ctx context.Context, tenantID string) (*EvaluationPeriod, error)
	ListPeriods(ctx context.Context, tenantID string) ([]EvaluationPeriod, error)
	CreatePeriod(ctx context.Context, tenantID string, label string, startDate, endDate *time.Time) (*EvaluationPeriod, error)
	UpdatePeriod(ctx context.Context, tenantID, periodID string, label string, startDate, endDate *time.Time) (*EvaluationPeriod, error)
	DeletePeriod(ctx context.Context, tenantID, periodID string) error
	SetActivePeriod(ctx context.Context, tenantID, periodID string) error
}

// ─── PeriodRepository ───
type PeriodRepository interface {
	ListByTenant(ctx context.Context, tenantID string) ([]EvaluationPeriod, error)
	GetActive(ctx context.Context, tenantID string) (*EvaluationPeriod, error)
	GetByID(ctx context.Context, tenantID, id string) (*EvaluationPeriod, error)
	Create(ctx context.Context, period *EvaluationPeriod) error
	Update(ctx context.Context, period *EvaluationPeriod) error
	Delete(ctx context.Context, tenantID, id string) error
	SetActive(ctx context.Context, tenantID, id string) error
}

type WorkflowUseCase interface {
	GetWorkflowForTenant(ctx context.Context, tenantID string) ([]EvaluationStage, error)
}

type TemplateUseCase interface {
	ListTemplates(ctx context.Context, filter TemplateFilter) ([]EvaluationTemplate, error)
	GetTemplate(ctx context.Context, id, tenantID string) (*EvaluationTemplate, error)
	CreateTemplate(ctx context.Context, template EvaluationTemplate) (*EvaluationTemplate, error)
	UpdateTemplate(ctx context.Context, template EvaluationTemplate) (*EvaluationTemplate, error)
	DeleteTemplate(ctx context.Context, id string, tenantID string) error
}

type EmployeeUseCase interface {
	ListEmployees(ctx context.Context, filter EmployeeFilter) ([]Employee, error)
	CreateEmployee(ctx context.Context, emp Employee) (*Employee, error)
	UpdateEmployee(ctx context.Context, emp Employee) (*Employee, error)
	PartialUpdateEmployee(ctx context.Context, id string, tenantID string, req UpdateEmployeeRequest) (*Employee, error)
	DeleteEmployee(ctx context.Context, id string, tenantID string) error
}

type DepartmentUseCase interface {
	ListDepartments(ctx context.Context, filter DepartmentFilter) ([]Department, error)
	CreateDepartment(ctx context.Context, dept Department) (*Department, error)
	UpdateDepartment(ctx context.Context, dept Department) (*Department, error)
	DeleteDepartment(ctx context.Context, id string, tenantID string) error
}

type AdjustmentUseCase interface {
	CreateAdjustment(ctx context.Context, adj SalaryAdjustment) (*SalaryAdjustment, error)
	ListAdjustments(ctx context.Context, tenantID string) ([]SalaryAdjustment, error)
}

type SalaryFormulaUseCase interface {
	GetFormula(ctx context.Context, tenantID string) (*SalaryFormula, error)
	UpdateFormula(ctx context.Context, formula SalaryFormula) (*SalaryFormula, error)
}

type PermissionUseCase interface {
	GetPermissions(ctx context.Context, tenantID string) ([]RolePermission, error)
	UpdatePermissions(ctx context.Context, tenantID string, perms []RolePermission) error
	HasPermission(ctx context.Context, tenantID string, role string, permKey string) (bool, error)
}

type UserUseCase interface {
	ListUsers(ctx context.Context, filter UserFilter) ([]User, error)
	CreateUser(ctx context.Context, req User) (*User, error)
	UpdateUser(ctx context.Context, id, tenantID, fullName, role, departmentID string) (*User, error)
	DeleteUser(ctx context.Context, id, tenantID string) error
	ChangePassword(ctx context.Context, userID, tenantID, newPassword string) error
}
