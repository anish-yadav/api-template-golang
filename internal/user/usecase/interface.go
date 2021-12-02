package usecase

import (
	"context"

	"github.com/anish-yadav/api-template-golang/internal/pkg/user"
)

type HttpManager interface {
	CreateUser(user user.User) (string, error)
	ChangePassword(ctx context.Context, old string, pwd string) error
	DeleteUser(id string) error
	ResetPassword(ctx context.Context, new string) error
	RequestPasswordReset(email string) error
	Login(username string, password string) (*user.UserDb, string, error)
	GetMe(ctx context.Context) (*user.UserDb, string, error)
}
