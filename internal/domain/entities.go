package domain

import (
	"time"

	"gorm.io/datatypes"
)

type Role string

const (
	RoleSuperAdmin Role = "SUPERADMIN"
	RoleCEO        Role = "CEO"
	RoleHead       Role = "HEAD_OF_DEPT"
	RoleHR         Role = "HR"
	RoleEmployee   Role = "EMPLOYEE"
)

// Permission keys for RBAC
const (
	PermViewDashboard  = "view_dashboard"
	PermViewSalary     = "view_salary"
	PermApproveSalary  = "approve_salary"
	PermCRUDEmployee   = "crud_employee"
	PermCRUDDepartment = "crud_department"
	PermCRUDTemplate   = "crud_template"
	PermManageUsers    = "manage_users"
	PermManageSettings = "manage_settings"
	PermEvalKPI        = "eval_kpi"
	PermManagePeriod   = "manage_period"
)

// AllPermissions lists every permission key
var AllPermissions = []string{
	PermViewDashboard,
	PermViewSalary,
	PermApproveSalary,
	PermCRUDEmployee,
	PermCRUDDepartment,
	PermCRUDTemplate,
	PermManageUsers,
	PermManageSettings,
	PermEvalKPI,
	PermManagePeriod,
}

type Tenant struct {
	ID          string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string               `gorm:"not null" json:"name"`
	Code        string               `gorm:"uniqueIndex;not null" json:"code"`
	Size        string               `gorm:"size:32" json:"size"`        // Small, Medium, Enterprise
	MaxUsers    int                  `gorm:"default:20" json:"maxUsers"` // Configurable user limit per tenant
	LogoURL     string               `gorm:"size:512" json:"logoUrl"`
	IsActive    bool                 `gorm:"default:true" json:"isActive"`
	Users       []User               `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	Templates   []EvaluationTemplate `gorm:"foreignKey:TenantID" json:"templates,omitempty"`
	Results     []EvaluationResult   `gorm:"foreignKey:TenantID" json:"results,omitempty"`
	Config      TenantConfig         `gorm:"foreignKey:TenantID" json:"config"`
	Departments []Department         `gorm:"foreignKey:TenantID" json:"departments,omitempty"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

// DefaultMaxUsers returns the default user limit for a given org size
func DefaultMaxUsers(size string) int {
	switch size {
	case "Small":
		return 20
	case "Medium":
		return 200
	case "Enterprise":
		return 1000
	default:
		return 20
	}
}

type TenantConfig struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID           string    `gorm:"uniqueIndex;not null" json:"tenantId"`
	EnableHREval       bool      `gorm:"column:enable_hr_eval;default:false" json:"enableHrEval"`
	EnableDeptHeadEval bool      `gorm:"column:enable_dept_head_eval;default:false" json:"enableDeptHeadEval"`
	EnableCEOEval      bool      `gorm:"column:enable_ceo_eval;default:true" json:"enableCeoEval"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// EvaluationPeriod represents a scheduled assessment cycle
type EvaluationPeriod struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID  string     `gorm:"index;not null" json:"tenantId"`
	Label     string     `gorm:"size:128;not null" json:"label"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	IsActive  bool       `gorm:"default:false" json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// Department is a DB-backed department per tenant
type Department struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID    string    `gorm:"index;not null" json:"tenantId"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Code        string    `gorm:"size:32;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DepartmentFilter struct {
	TenantID    string
	SearchQuery string
}

// SalaryFormula stores per-tenant salary calculation logic
type SalaryFormula struct {
	ID                string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID          string         `gorm:"uniqueIndex;not null" json:"tenantId"`
	TenureRatePerYear float64        `gorm:"default:0.5" json:"tenureRatePerYear"`
	MaxTenureBonus    float64        `gorm:"default:10" json:"maxTenureBonus"`
	MeritMap          datatypes.JSON `gorm:"type:jsonb;not null" json:"meritMap"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// RolePermission stores which role has which permissions per tenant
type RolePermission struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID      string    `gorm:"uniqueIndex:idx_tenant_role_perm;not null" json:"tenantId"`
	Role          string    `gorm:"uniqueIndex:idx_tenant_role_perm;size:32;not null" json:"role"`
	PermissionKey string    `gorm:"uniqueIndex:idx_tenant_role_perm;size:64;not null" json:"permissionKey"`
	Allowed       bool      `gorm:"default:false" json:"allowed"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID     string    `gorm:"index;not null" json:"tenantId"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"not null" json:"-"`
	FullName     string    `gorm:"not null" json:"fullName"`
	Role         Role      `gorm:"type:varchar(32);not null;default:'EMPLOYEE'" json:"role"`
	DepartmentID string    `gorm:"size:128;default:''" json:"departmentId"` // set for HEAD_OF_DEPT
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type EmployeeFilter struct {
	TenantID               string
	Status                 string
	SearchQuery            string
	ExcludeEvaluatedPeriod string
	EvaluatorID            string
	DepartmentID           string
}

type Employee struct {
	ID               string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID         string      `gorm:"index;not null" json:"tenantId"`
	Code             string      `gorm:"size:64;not null" json:"code"`
	FirstName        string      `gorm:"size:128;not null" json:"firstName"`
	LastName         string      `gorm:"size:128;not null" json:"lastName"`
	Email            string      `gorm:"size:128;default:''" json:"email"`
	DepartmentID     string      `gorm:"type:uuid;index" json:"departmentId"`
	Department       *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position         string      `gorm:"size:128;not null" json:"position"`
	BaseSalary       float64     `gorm:"not null" json:"baseSalary"`
	PersonalCapacity float64     `gorm:"default:0" json:"personalCapacity"`
	VariablePayBase  float64     `gorm:"default:0" json:"variablePayBase"`
	YearsOfService   int         `gorm:"default:0" json:"yearsOfService"`
	Status           string      `gorm:"size:32;default:'ACTIVE'" json:"status"`
	StartDate        string      `gorm:"size:32" json:"startDate"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

type TemplateFilter struct {
	TenantID     string
	SearchQuery  string
	DepartmentID string
}

type EvaluationTemplate struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID     string         `gorm:"index;not null" json:"tenantId"`
	Name         string         `gorm:"not null" json:"name"`
	DepartmentID string         `gorm:"type:uuid;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"departmentId"`
	Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Visibility   string         `gorm:"size:32;default:'GENERAL'" json:"visibility"` // PERSONAL | GENERAL
	Period       string         `gorm:"size:64;default:'รายไตรมาส'" json:"period"`
	Definition   datatypes.JSON `gorm:"type:jsonb" json:"definition"`
	CreatedByID  string         `gorm:"not null" json:"createdById"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type EvaluationStage string

const (
	StageHR       EvaluationStage = "HR"
	StageDeptHead EvaluationStage = "DEPT_HEAD"
	StageCEO      EvaluationStage = "CEO"
)

type EvaluationFilter struct {
	TenantID          string
	ViewerID          string
	ViewerRole        string
	SubjectEmployeeID string
	DepartmentID      string
	SearchQuery       string
	Period            string
	SummaryOnly       bool
}

type EvaluationResult struct {
	ID           string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID     string           `gorm:"index;not null" json:"tenantId"`
	EmployeeID   string           `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"employeeId"`
	TemplateID   string           `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"templateId"`
	EmployeeName string           `gorm:"size:255" json:"employeeName"`
	Period       string           `gorm:"not null" json:"period"`
	PeriodType   string           `gorm:"not null" json:"periodType"`
	TotalScore   float64          `gorm:"default:0" json:"totalScore"`
	RatingLevel  string           `gorm:"size:64" json:"ratingLevel"`
	Notes        string           `gorm:"type:text" json:"notes"`
	Details      datatypes.JSON   `gorm:"type:jsonb" json:"details"`
	CurrentStage *EvaluationStage `gorm:"type:varchar(32)" json:"currentStage"`
	Completed    bool             `gorm:"default:false" json:"completed"`
	// Evaluator attribution
	EvaluatorID   string `gorm:"type:uuid;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"evaluatorId"`
	EvaluatorName string `gorm:"size:255;default:''" json:"evaluatorName"`
	EvaluatorRole string `gorm:"size:32;default:''" json:"evaluatorRole"`
	// Composite index for dashboard recent-evaluations query (tenant + time desc)
	EvaluatedAt time.Time `gorm:"index:idx_eval_tenant_time,priority:2" json:"evaluatedAt"`
	CreatedAt   time.Time `gorm:"index:idx_eval_tenant_time,priority:1" json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserFilter struct {
	TenantID    string
	SearchQuery string
}

type SalaryAdjustment struct {
	ID                     string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID               string    `gorm:"index;not null" json:"tenantId"`
	EmployeeID             string    `gorm:"index;not null" json:"employeeId"`
	Period                 string    `gorm:"index;not null" json:"period"`
	KPIScore               int       `json:"kpiScore"`
	RatingLevel            string    `json:"ratingLevel"`
	MeritPercent           float64   `json:"meritPercent"`
	TenureBonus            float64   `json:"tenureBonus"`
	MarketCorrection       float64   `json:"marketCorrection"`
	CostOfLiving           float64   `json:"costOfLiving"`
	TotalAdjustmentPercent float64   `json:"totalAdjustmentPercent"`
	CurrentSalary          float64   `json:"currentSalary"`
	RecommendedSalary      float64   `json:"recommendedSalary"`
	P3Amount               float64   `json:"p3Amount"`
	Approved               bool      `gorm:"default:false" json:"approved"`
	CreatedAt              time.Time `json:"createdAt"`
}
