package usecase

import (
	"context"
	"kpi-backend/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
)

type evaluationUseCase struct {
	evalRepo     domain.EvaluationRepository
	workflowUC   domain.WorkflowUseCase
	employeeRepo domain.EmployeeRepository
	templateRepo domain.TemplateRepository
	periodRepo   domain.PeriodRepository
	userRepo     domain.UserRepository
}

func NewEvaluationUseCase(
	evalRepo domain.EvaluationRepository,
	workflowUC domain.WorkflowUseCase,
	empRepo domain.EmployeeRepository,
	tplRepo domain.TemplateRepository,
	periodRepo domain.PeriodRepository,
	userRepo domain.UserRepository,
) domain.EvaluationUseCase {
	return &evaluationUseCase{
		evalRepo:     evalRepo,
		workflowUC:   workflowUC,
		employeeRepo: empRepo,
		templateRepo: tplRepo,
		periodRepo:   periodRepo,
		userRepo:     userRepo,
	}
}

func (u *evaluationUseCase) CreateEvaluation(ctx context.Context, tenantID string, req domain.EvaluationResult) (*domain.EvaluationResult, error) {
	// 1. Enforce system period — user cannot set period manually
	activePeriod, err := u.periodRepo.GetActive(ctx, tenantID)
	if err != nil || activePeriod == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ยังไม่ได้ตั้งค่ารอบการประเมิน กรุณาติดต่อผู้ดูแลระบบเพื่อกำหนดรอบการประเมิน")
	}
	// Always override with the active period label
	req.Period = activePeriod.Label
	req.EvaluatedAt = time.Now()

	// 2. Fetch Employee
	emp, err := u.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, err
	}
	req.EmployeeName = emp.FirstName + " " + emp.LastName

	// 3. Fetch Template
	tpl, err := u.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}

	// 4. Department validation (แม่แบบ)
	if tpl.DepartmentID != "" && emp.DepartmentID != tpl.DepartmentID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Employee department does not match template department")
	}

	// 5. Handle workflow
	stages, err := u.workflowUC.GetWorkflowForTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	first := stages[0]
	req.TenantID = tenantID
	req.CurrentStage = &first

	if err := u.evalRepo.Create(ctx, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (u *evaluationUseCase) ListEvaluations(ctx context.Context, filter domain.EvaluationFilter) ([]domain.EvaluationResult, error) {
	if filter.ViewerRole == string(domain.RoleEmployee) {
		user, err := u.userRepo.GetByID(ctx, filter.ViewerID)
		if err == nil {
			emp, empErr := u.employeeRepo.GetByEmail(ctx, filter.TenantID, user.Email)
			if empErr == nil {
				filter.SubjectEmployeeID = emp.ID
			} else {
				return []domain.EvaluationResult{}, nil
			}
		} else {
			return []domain.EvaluationResult{}, nil
		}
	}
	return u.evalRepo.List(ctx, filter)
}

func (u *evaluationUseCase) DeleteEvaluation(ctx context.Context, id string, tenantID string) error {
	return u.evalRepo.Delete(ctx, id, tenantID)
}
