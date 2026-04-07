package usecase

import (
	"context"
	"errors"
	"kpi-backend/internal/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {
	return u.userRepo.List(ctx, filter)
}

func (u *userUseCase) CreateUser(ctx context.Context, req domain.User) (*domain.User, error) {
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, errors.New("email, password and fullName are required")
	}
	validRoles := map[string]bool{
		"CEO": true, "HR": true, "HEAD_OF_DEPT": true, "EMPLOYEE": true,
	}
	roleStr := string(req.Role)
	if !validRoles[roleStr] {
		roleStr = "EMPLOYEE"
	}
	user := &domain.User{
		TenantID:     req.TenantID,
		Email:        req.Email,
		Password:     req.Password,
		FullName:     req.FullName,
		Role:         domain.Role(roleStr),
		DepartmentID: req.DepartmentID,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id, tenantID, fullName, role, departmentID string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if user.TenantID != tenantID {
		return nil, domain.ErrNotFound
	}
	validRoles := map[string]bool{
		"CEO": true, "HR": true, "HEAD_OF_DEPT": true, "EMPLOYEE": true,
	}
	if fullName != "" {
		user.FullName = fullName
	}
	if validRoles[role] {
		user.Role = domain.Role(role)
	}
	// Update department assignment (relevant for HEAD_OF_DEPT)
	user.DepartmentID = departmentID
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id, tenantID string) error {
	return u.userRepo.Delete(ctx, id, tenantID)
}
