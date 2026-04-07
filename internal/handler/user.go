package handler

import (
	"kpi-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUC domain.UserUseCase
}

func NewUserHandler(userUC domain.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

type createUserRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"fullName"`
	Role         string `json:"role"`
	DepartmentID string `json:"departmentId"`
}

type updateUserRequest struct {
	FullName     string `json:"fullName"`
	Role         string `json:"role"`
	DepartmentID string `json:"departmentId"`
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	filter := domain.UserFilter{
		TenantID:    tenantID,
		SearchQuery: c.Query("q"),
	}

	users, err := h.userUC.ListUsers(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Map to SafeUser DTO — password is never serialised (json:"-" in domain.User)
	safe := make([]SafeUser, len(users))
	for i, u := range users {
		safe[i] = SafeUser{
			ID:           u.ID,
			Email:        u.Email,
			FullName:     u.FullName,
			Role:         string(u.Role),
			TenantID:     u.TenantID,
			DepartmentID: u.DepartmentID,
		}
	}
	return c.JSON(safe)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	var body createUserRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	userReq := domain.User{
		TenantID:     tenantID,
		Email:        body.Email,
		Password:     body.Password,
		FullName:     body.FullName,
		Role:         domain.Role(body.Role),
		DepartmentID: body.DepartmentID,
	}
	user, err := h.userUC.CreateUser(c.Context(), userReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(SafeUser{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         string(user.Role),
		TenantID:     user.TenantID,
		DepartmentID: user.DepartmentID,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	var body updateUserRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	user, err := h.userUC.UpdateUser(c.Context(), id, tenantID, body.FullName, body.Role, body.DepartmentID)
	if err != nil {
		if err == domain.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(SafeUser{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         string(user.Role),
		TenantID:     user.TenantID,
		DepartmentID: user.DepartmentID,
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	tenantID := tenantFromCtx(c)
	id := c.Params("id")
	if err := h.userUC.DeleteUser(c.Context(), id, tenantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
