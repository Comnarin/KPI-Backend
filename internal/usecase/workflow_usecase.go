package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type workflowUseCase struct {
	tenantRepo domain.TenantRepository
}

func NewWorkflowUseCase(tenantRepo domain.TenantRepository) domain.WorkflowUseCase {
	return &workflowUseCase{tenantRepo: tenantRepo}
}

func (u *workflowUseCase) GetWorkflowForTenant(ctx context.Context, tenantID string) ([]domain.EvaluationStage, error) {
	cfg, err := u.tenantRepo.GetConfig(ctx, tenantID)
	if err != nil {
		// If not found, default to CEO only
		return []domain.EvaluationStage{domain.StageCEO}, nil
	}

	stages := []domain.EvaluationStage{}
	if cfg.EnableHREval {
		stages = append(stages, domain.StageHR)
	}
	if cfg.EnableDeptHeadEval {
		stages = append(stages, domain.StageDeptHead)
	}
	if cfg.EnableCEOEval {
		stages = append(stages, domain.StageCEO)
	}

	if len(stages) == 0 {
		return []domain.EvaluationStage{domain.StageCEO}, nil
	}
	return stages, nil
}
