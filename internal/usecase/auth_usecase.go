package usecase

import (
	"context"
	"os"
	"time"
	"kpi-backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type authUseCase struct {
	userRepo   domain.UserRepository
	tenantRepo domain.TenantRepository
	secret     string
}

func NewAuthUseCase(userRepo domain.UserRepository, tenantRepo domain.TenantRepository) domain.AuthUseCase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-change-me"
	}
	return &authUseCase{userRepo: userRepo, tenantRepo: tenantRepo, secret: secret}
}

func (u *authUseCase) Login(ctx context.Context, tenantCode, email, password string) (string, *domain.User, error) {
	// If tenantCode provided, validate it matches the user's tenant
	if tenantCode != "" {
		tenant, err := u.tenantRepo.GetByCode(ctx, tenantCode)
		if err != nil {
			return "", nil, domain.ErrNotFound
		}
		_ = tenant // Used for validation — we'll compare after fetching user
	}

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// If tenant code provided, validate user belongs to this tenant
	if tenantCode != "" {
		tenant, _ := u.tenantRepo.GetByCode(ctx, tenantCode)
		if tenant != nil && user.TenantID != tenant.ID {
			return "", nil, domain.ErrInvalidCredentials
		}
	}

	// Plain-text password check (demo mode — use bcrypt in production)
	if user.Password != password {
		return "", nil, domain.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"tid":  user.TenantID,
		"role": string(user.Role),
		"name": user.FullName,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.secret))
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}
