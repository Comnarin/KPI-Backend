package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUC domain.AuthUseCase
}

func NewAuthHandler(authUC domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type loginRequest struct {
	TenantCode string `json:"tenantCode"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	token, user, err := h.authUC.Login(c.Context(), body.TenantCode, body.Email, body.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials || err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"token":        token,
		"id":           user.ID,
		"fullName":     user.FullName,
		"tenantId":     user.TenantID,
		"role":         user.Role,
		"departmentId": user.DepartmentID,
	})
}
