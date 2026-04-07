package domain

import "context"

type RatingDistribution struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type RecentEvaluation struct {
	ID          string `json:"id"`
	EmployeeID  string `json:"employeeId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Position    string `json:"position"`
	TotalScore  int    `json:"totalScore"`
	RatingLevel string `json:"ratingLevel"`
	EvaluatedAt string `json:"evaluatedAt"`
}

type BarChartData struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type RadarChartData struct {
	Subject string `json:"subject"`
	Value   int    `json:"value"`
}

type DashboardStats struct {
	TotalEmployees    int                  `json:"totalEmployees"`
	EvaluatedCount    int                  `json:"evaluatedCount"`
	AverageScore      int                  `json:"averageScore"`
	TotalSalaryBudget float64              `json:"totalSalaryBudget"`
	BarData           []BarChartData       `json:"barData"`
	RadarData         []RadarChartData     `json:"radarData"`
	RatingDist        []RatingDistribution `json:"ratingDist"`
	RecentEvaluations []RecentEvaluation   `json:"recentEvaluations"`
}

type DashboardRepository interface {
	GetDashboardStats(ctx context.Context, tenantID string, period string) (*DashboardStats, error)
}

type DashboardUseCase interface {
	GetDashboardStats(ctx context.Context, tenantID string, period string) (*DashboardStats, error)
}
