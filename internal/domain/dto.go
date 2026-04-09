package domain

// UpdateEmployeeRequest defines the fields allowed for partial updates.
// Pointers are used to distinguish between a missing field and a zero-value field.
type UpdateEmployeeRequest struct {
	FirstName        *string  `json:"firstName,omitempty"`
	LastName         *string  `json:"lastName,omitempty"`
	Email            *string  `json:"email,omitempty"`
	DepartmentID     *string  `json:"departmentId,omitempty"`
	Position         *string  `json:"position,omitempty"`
	BaseSalary       *float64 `json:"baseSalary,omitempty"`
	PersonalCapacity *float64 `json:"personalCapacity,omitempty"`
	VariablePayBase  *float64 `json:"variablePayBase,omitempty"`
	Code             *string  `json:"code,omitempty"`
	Status           *string  `json:"status,omitempty"`
	YearsOfService   *int     `json:"yearsOfService,omitempty"`
	StartDate        *string  `json:"startDate,omitempty"`
}

// EmployeeLightDTO represents a simplified employee object for list views.
type EmployeeLightDTO struct {
	ID           string      `json:"id"`
	Code         string      `json:"code"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	DepartmentID string      `json:"departmentId"`
	Department   *Department `json:"department,omitempty"`
	Position     string      `json:"position"`
	Status       string      `json:"status"`
}

// EvaluationResponseDTO embeds EvaluationResult and adds the virtual departmentName field.
type EvaluationResponseDTO struct {
	EvaluationResult
	DepartmentName string `json:"departmentName"`
}
