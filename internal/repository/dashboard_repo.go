package repository

import (
	"context"
	"math"
	"sync"

	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetDashboardStats fetches all dashboard data concurrently to minimise latency.
// Each independent query runs in its own goroutine; errors are collected and the
// first non-nil error is returned after all goroutines complete.
func (r *dashboardRepository) GetDashboardStats(ctx context.Context, tenantID string, period string) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{
		BarData:           []domain.BarChartData{},
		RadarData:         []domain.RadarChartData{},
		RatingDist:        []domain.RatingDistribution{},
		RecentEvaluations: []domain.RecentEvaluation{},
		Periods:           []domain.EvaluationPeriod{},
	}

	// Protect concurrent writes to stats
	var mu sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0, 6)

	addErr := func(err error) {
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	// ── 1. Total Employees ───────────────────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		err := r.db.WithContext(ctx).
			Model(&domain.Employee{}).
			Where("tenant_id = ?", tenantID).
			Count(&count).Error
		mu.Lock()
		stats.TotalEmployees = int(count)
		mu.Unlock()
		addErr(err)
	}()

	// ── 2. Evaluated Count ───────────────────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		query := r.db.WithContext(ctx).
			Model(&domain.EvaluationResult{}).
			Where("tenant_id = ?", tenantID)

		if period != "" {
			query = query.Where("period = ?", period)
		}

		err := query.Select("count(distinct employee_id)").Scan(&count).Error
		mu.Lock()
		stats.EvaluatedCount = int(count)
		mu.Unlock()
		addErr(err)
	}()

	// ── 3. Average Score ─────────────────────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		var avg float64
		// First, calculate the average for EACH employee (grouped by employee_id)
		subQuery := r.db.WithContext(ctx).
			Model(&domain.EvaluationResult{}).
			Select("AVG(total_score) as emp_avg").
			Where("tenant_id = ?", tenantID)

		if period != "" {
			subQuery = subQuery.Where("period = ?", period)
		}
		subQuery = subQuery.Group("employee_id")

		// Then, calculate the team-wide average of those per-employee averages
		err := r.db.WithContext(ctx).
			Table("(?) as sub", subQuery).
			Select("COALESCE(AVG(emp_avg), 0)").
			Scan(&avg).Error
		mu.Lock()
		stats.AverageScore = int(math.Round(avg))
		mu.Unlock()
		addErr(err)
	}()

	// ── 4. Top 8 Employees by Score (Bar Chart) ──────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		type barRow struct {
			FirstName  string
			TotalScore int
		}
		var rows []barRow
		// Use CAST to guarantee uuid comparison works even if column type drifts
		query := r.db.WithContext(ctx).
			Table("employees e").
			Select("e.first_name, CAST(COALESCE(AVG(ev.total_score), 0) AS INT) AS total_score")
		
		if period != "" {
			query = query.Joins("LEFT JOIN evaluation_results ev ON ev.employee_id = e.id AND ev.period = ?", period)
		} else {
			query = query.Joins("LEFT JOIN evaluation_results ev ON ev.employee_id = e.id")
		}

		err := query.Where("e.tenant_id = ?", tenantID).
			Group("e.id, e.first_name").
			Order("total_score DESC").
			Limit(8).
			Scan(&rows).Error

		mu.Lock()
		for _, row := range rows {
			stats.BarData = append(stats.BarData, domain.BarChartData{
				Name:  row.FirstName,
				Score: row.TotalScore,
			})
		}
		mu.Unlock()
		addErr(err)
	}()

	// ── 5. Rating Distribution ───────────────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		type distRow struct {
			RatingLevel string
			Count       int
		}
		var rows []distRow
		query := r.db.WithContext(ctx).
			Model(&domain.EvaluationResult{}).
			Select("rating_level, COUNT(*) AS count").
			Where("tenant_id = ?", tenantID)

		if period != "" {
			query = query.Where("period = ?", period)
		}

		err := query.Group("rating_level").Scan(&rows).Error

		// Build ordered distribution with defaults of 0
		levels := []string{"ดีเยี่ยม", "เกินเป้า", "ได้เป้า", "ต้องปรับปรุง", "ไม่ผ่านเกณฑ์"}
		ratingMap := make(map[string]int, len(levels))
		for _, l := range levels {
			ratingMap[l] = 0
		}
		for _, row := range rows {
			if row.RatingLevel != "" {
				ratingMap[row.RatingLevel] = row.Count
			}
		}

		dist := make([]domain.RatingDistribution, 0, len(levels))
		for _, l := range levels {
			dist = append(dist, domain.RatingDistribution{Level: l, Count: ratingMap[l]})
		}
		mu.Lock()
		stats.RatingDist = dist
		mu.Unlock()
		addErr(err)
	}()

	// ── 6. Recent Evaluations (last 6) ───────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		type recentRow struct {
			ID             string
			EmployeeID     string
			FirstName      string
			LastName       string
			Position       string
			DepartmentName string
			TotalScore     int
			RatingLevel    string
			EvaluatedAt    string // formatted from DB
		}
		var rows []recentRow
		query := r.db.WithContext(ctx).
			Table("evaluation_results ev").
			Select(`ev.id, ev.employee_id, e.first_name, e.last_name, e.position,
				COALESCE(d.name, '') AS department_name,
				ev.total_score, ev.rating_level,
				TO_CHAR(ev.evaluated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS evaluated_at`).
			Joins("JOIN employees e ON ev.employee_id = e.id").
			Joins("LEFT JOIN departments d ON d.id = e.department_id").
			Where("ev.tenant_id = ?", tenantID)

		if period != "" {
			query = query.Where("ev.period = ?", period)
		}

		err := query.Order("ev.evaluated_at DESC NULLS LAST").
			Limit(6).
			Scan(&rows).Error

		recents := make([]domain.RecentEvaluation, 0, len(rows))
		for _, rr := range rows {
			recents = append(recents, domain.RecentEvaluation{
				ID:             rr.ID,
				EmployeeID:     rr.EmployeeID,
				FirstName:      rr.FirstName,
				LastName:       rr.LastName,
				Position:       rr.Position,
				DepartmentName: rr.DepartmentName,
				TotalScore:     rr.TotalScore,
				RatingLevel:    rr.RatingLevel,
				EvaluatedAt:    rr.EvaluatedAt,
			})
		}
		mu.Lock()
		stats.RecentEvaluations = recents
		mu.Unlock()
		addErr(err)
	}()

	// ── 7. Total Salary Budget (Approved Adjustments) ────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		var totalIncrease float64
		// SUM(recommended_salary - current_salary) for approved adjustments
		query := r.db.WithContext(ctx).
			Model(&domain.SalaryAdjustment{}).
			Where("tenant_id = ? AND approved = ?", tenantID, true)

		if period != "" {
			query = query.Where("period = ?", period)
		}

		err := query.Select("COALESCE(SUM(recommended_salary - current_salary), 0)").Scan(&totalIncrease).Error
		mu.Lock()
		stats.TotalSalaryBudget = totalIncrease
		mu.Unlock()
		addErr(err)
	}()

	// ── 8. Periods (for dropdown filter) ─────────────────────────────────────
	wg.Add(1)
	go func() {
		defer wg.Done()
		var periods []domain.EvaluationPeriod
		err := r.db.WithContext(ctx).
			Where("tenant_id = ?", tenantID).
			Order("created_at DESC").
			Find(&periods).Error
		mu.Lock()
		if periods != nil {
			stats.Periods = periods
		}
		mu.Unlock()
		addErr(err)
	}()

	// ── Static radar data (KPI category breakdown — replace with real query when available) ──
	stats.RadarData = []domain.RadarChartData{
		{Subject: "ยอดขาย", Value: 75},
		{Subject: "คุณภาพงาน", Value: 85},
		{Subject: "ทำงานเป็นทีม", Value: 65},
		{Subject: "ตรงต่อเวลา", Value: 80},
		{Subject: "ความคิริเริ่ม", Value: 70},
	}

	wg.Wait()

	// Return first error encountered (all queries ran regardless)
	for _, err := range errs {
		if err != nil {
			return stats, err
		}
	}
	return stats, nil
}
