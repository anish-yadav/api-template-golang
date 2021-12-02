package handler

import "github.com/anish-yadav/api-template-golang/internal/pkg/user"

type UserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"required"`
}

type UserRequestResponse struct {
	ID string `json:"id"`
}

type ChangePasswordRequest struct {
	ID          string `json:"id" validate:"required"`
	OldPassword string `json:"old" validate:"required"`
	NewPassword string `json:"new" validate:"required"`
}

type ReqResetPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *UserRequest) toUser() user.User {
	return user.User{
		Name:  u.Name,
		Email: u.Email,
		Type:  u.Type,
	}
}
