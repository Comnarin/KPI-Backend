package handler

import (
	"encoding/json"
	"kpi-backend/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeedHandler struct {
	DB *gorm.DB
}

func NewSeedHandler(db *gorm.DB) *SeedHandler {
	return &SeedHandler{DB: db}
}

func (h *SeedHandler) Seed(c *fiber.Ctx) error {
	tenantID := "d86865e0-3f0e-4899-9b87-134d7947f6c7"

	tenant := domain.Tenant{
		ID:       tenantID,
		Name:     "Demo Company",
		Code:     "DEMO",
		Size:     "SMALL",
		IsActive: true,
	}

	tenantConfig := domain.TenantConfig{
		TenantID:      tenantID,
		EnableCEOEval: true,
	}

	user := domain.User{
		TenantID: tenantID,
		Email:    "ceo@demo.com",
		Password: "password",
		FullName: "Demo CEO",
		Role:     domain.RoleCEO,
	}

	hrUser := domain.User{
		TenantID: tenantID,
		Email:    "hr@demo.com",
		Password: "password",
		FullName: "Demo HR",
		Role:     domain.RoleHR,
	}

	superAdmin := domain.User{
		TenantID: tenantID,
		Email:    "admin@system.com",
		Password: "admin123",
		FullName: "Super Admin",
		Role:     domain.RoleSuperAdmin,
	}

	departments := []domain.Department{
		{TenantID: tenantID, Name: "ฝ่ายพัฒนาซอฟต์แวร์", Code: "DEV", Description: "ทีมพัฒนาซอฟต์แวร์"},
		{TenantID: tenantID, Name: "ฝ่ายขาย", Code: "SALES", Description: "ทีมขายและการตลาด"},
		{TenantID: tenantID, Name: "ฝ่ายทรัพยากรบุคคล", Code: "HR", Description: "ทีม HR"},
		{TenantID: tenantID, Name: "ฝ่ายการเงิน", Code: "FINANCE", Description: "ทีมการเงินและบัญชี"},
	}

	employees := []domain.Employee{
		{
			TenantID: tenantID, Code: "EMP001",
			FirstName: "สมชาย", LastName: "ทดลอง",
			Department: "ฝ่ายขาย", Position: "Sales Rep",
			BaseSalary: 30000, PersonalCapacity: 2000, VariablePayBase: 5000,
			YearsOfService: 3, StartDate: "2021-01-15",
		},
		{
			TenantID: tenantID, Code: "EMP002",
			FirstName: "สมหญิง", LastName: "ทดสอบ",
			Department: "ฝ่ายพัฒนาซอฟต์แวร์", Position: "Developer",
			BaseSalary: 40000, PersonalCapacity: 5000, VariablePayBase: 8000,
			YearsOfService: 5, StartDate: "2019-06-20",
		},
		{
			TenantID: tenantID, Code: "EMP003",
			FirstName: "มานี", LastName: "รักงาน",
			Department: "ฝ่ายทรัพยากรบุคคล", Position: "HR Officer",
			BaseSalary: 28000, PersonalCapacity: 1500, VariablePayBase: 3000,
			YearsOfService: 2, StartDate: "2022-03-01",
		},
	}

	templates := []domain.EvaluationTemplate{
		{
			TenantID: tenantID, Name: "แบบประเมินทีมพัฒนาซอฟต์แวร์",
			Department: "ฝ่ายพัฒนาซอฟต์แวร์", Period: "รายไตรมาส",
			Definition: datatypes.JSON([]byte(`[{"id":"k1","name":"อัตราการทำงานเสร็จใน Sprint","weight":30,"targetValue":100,"unit":"%"},{"id":"k2","name":"คุณภาพโค้ด (Bug Rate)","weight":25,"targetValue":2,"unit":"บั๊ก/ฟีเจอร์"},{"id":"k3","name":"การส่งงานตรงเวลา","weight":20,"targetValue":95,"unit":"%"},{"id":"k4","name":"การ Review Code","weight":15,"targetValue":100,"unit":"%"},{"id":"k5","name":"คุณภาพเอกสาร","weight":10,"targetValue":5,"unit":"คะแนน"}]`)),
			CreatedByID: user.ID,
		},
		{
			TenantID: tenantID, Name: "แบบประเมินทีมขาย",
			Department: "ฝ่ายขาย", Period: "รายไตรมาส",
			Definition: datatypes.JSON([]byte(`[{"id":"k1","name":"ยอดขายรายไตรมาส","weight":35,"targetValue":3000000,"unit":"บาท"},{"id":"k2","name":"อัตราการรักษาลูกค้า","weight":25,"targetValue":90,"unit":"%"},{"id":"k3","name":"จำนวนลูกค้าใหม่","weight":20,"targetValue":20,"unit":"ราย"},{"id":"k4","name":"ความพึงพอใจลูกค้า","weight":10,"targetValue":90,"unit":"%"},{"id":"k5","name":"การทำงานเป็นทีม","weight":10,"targetValue":5,"unit":"คะแนน"}]`)),
			CreatedByID: user.ID,
		},
	}

	// Default salary formula
	meritMap, _ := json.Marshal(map[string]int{
		"ดีเยี่ยม":     8,
		"เกินเป้า":     6,
		"ได้เป้า":      4,
		"ต้องปรับปรุง": 2,
		"ไม่ผ่านเกณฑ์": 0,
	})
	salaryFormula := domain.SalaryFormula{
		TenantID:          tenantID,
		TenureRatePerYear: 0.5,
		MaxTenureBonus:    10.0,
		MeritMap:          datatypes.JSON(meritMap),
		UpdatedAt:         time.Now(),
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Clear everything
		tx.Exec("DELETE FROM evaluation_results")
		tx.Exec("DELETE FROM evaluation_templates")
		tx.Exec("DELETE FROM employees")
		tx.Exec("DELETE FROM departments")
		tx.Exec("DELETE FROM salary_formulas")
		tx.Exec("DELETE FROM role_permissions")
		tx.Exec("DELETE FROM tenant_configs")
		tx.Exec("DELETE FROM users")
		tx.Exec("DELETE FROM tenants")

		// Create Tenant
		if err := tx.Omit(clause.Associations).Create(&tenant).Error; err != nil {
			return err
		}

		// Create Users
		user.TenantID = tenant.ID
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		hrUser.TenantID = tenant.ID
		if err := tx.Create(&hrUser).Error; err != nil {
			return err
		}
		superAdmin.TenantID = tenant.ID
		if err := tx.Create(&superAdmin).Error; err != nil {
			return err
		}

		// Create config
		tenantConfig.TenantID = tenant.ID
		if err := tx.Create(&tenantConfig).Error; err != nil {
			return err
		}

		// Create departments
		for i := range departments {
			departments[i].TenantID = tenant.ID
		}
		if err := tx.Create(&departments).Error; err != nil {
			return err
		}

		// Create employees
		for i := range employees {
			employees[i].TenantID = tenant.ID
		}
		if err := tx.Create(&employees).Error; err != nil {
			return err
		}

		// Create templates (use user.ID now that it exists)
		for i := range templates {
			templates[i].TenantID = tenant.ID
			templates[i].CreatedByID = user.ID
		}
		if err := tx.Create(&templates).Error; err != nil {
			return err
		}

		// Create salary formula
		salaryFormula.TenantID = tenant.ID
		if err := tx.Create(&salaryFormula).Error; err != nil {
			return err
		}

		// Create default permissions for each role
		var perms []domain.RolePermission
		roles := []string{
			string(domain.RoleCEO),
			string(domain.RoleHR),
			string(domain.RoleHead),
			string(domain.RoleEmployee),
		}
		defaults := map[string]map[string]bool{
			string(domain.RoleCEO):      {domain.PermViewDashboard: true, domain.PermViewSalary: true, domain.PermApproveSalary: true, domain.PermCRUDEmployee: true, domain.PermCRUDDepartment: true, domain.PermCRUDTemplate: true, domain.PermManageUsers: true, domain.PermManageSettings: true},
			string(domain.RoleHR):       {domain.PermViewDashboard: true, domain.PermViewSalary: false, domain.PermApproveSalary: false, domain.PermCRUDEmployee: true, domain.PermCRUDDepartment: true, domain.PermCRUDTemplate: true, domain.PermManageUsers: false, domain.PermManageSettings: false},
			string(domain.RoleHead):     {domain.PermViewDashboard: true, domain.PermViewSalary: true, domain.PermApproveSalary: false, domain.PermCRUDEmployee: false, domain.PermCRUDDepartment: false, domain.PermCRUDTemplate: true, domain.PermManageUsers: false, domain.PermManageSettings: false},
			string(domain.RoleEmployee): {domain.PermViewDashboard: true, domain.PermViewSalary: false, domain.PermApproveSalary: false, domain.PermCRUDEmployee: false, domain.PermCRUDDepartment: false, domain.PermCRUDTemplate: false, domain.PermManageUsers: false, domain.PermManageSettings: false},
		}
		for _, role := range roles {
			for permKey, allowed := range defaults[role] {
				perms = append(perms, domain.RolePermission{
					TenantID:      tenant.ID,
					Role:          role,
					PermissionKey: permKey,
					Allowed:       allowed,
				})
			}
		}
		if err := tx.Create(&perms).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"tenantId":    tenantID,
		"tenantCode":  "DEMO",
		"ceoEmail":    user.Email,
		"hrEmail":     hrUser.Email,
		"adminEmail":  superAdmin.Email,
		"status":      "success",
	})
}
