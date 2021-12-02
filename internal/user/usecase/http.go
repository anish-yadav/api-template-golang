package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/anish-yadav/api-template-golang/internal/constants"
	"github.com/anish-yadav/api-template-golang/internal/pkg/user"
	"github.com/anish-yadav/api-template-golang/internal/util"
	log "github.com/sirupsen/logrus"
)

type httpManager struct {
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewHttpManager() *httpManager {
	rand.Seed(time.Now().UnixNano())
	return &httpManager{}
}

func (m *httpManager) CreateUser(u user.User) (string, error) {
	password := randStringRunes(8)
	log.Debugf("email: %s, password: %s", u.Email, password)
	newUser := user.NewUser(u.Name, u.Email, password, u.Type)
	// TODO need to send mail to user to reset its password
	id, err := newUser.AddToDB()
	if err != nil {
		return "", err
	}
	req := user.NewPasswordResetRequest(newUser)
	token, err := req.SendRequest()
	if err != nil {
		return "", err
	}
	log.Debugf("CreateUser: reset password req token : %s", token)
	return id, err
}

func (m *httpManager) RequestPasswordReset(email string) error {
	currUser := user.GetUserByEmail(email)
	if currUser == nil {
		return errors.New(constants.ItemNotFound)
	}
	req := user.NewPasswordResetRequest(currUser)
	token, err := req.SendRequest()
	if err != nil {
		return err
	}
	log.Debugf("reqPassReset: reset password req token : %s", token)
	return nil
}

func (m *httpManager) ChangePassword(ctx context.Context, old string, new string) error {
	currUser := ctx.Value("user").(*user.UserDb)
	return currUser.ChangePassword(old, new)
}

func (m *httpManager) ResetPassword(ctx context.Context, new string) error {
	currUser := ctx.Value("user").(*user.UserDb)
	err := currUser.ResetPassword(new)
	if err != nil {
		return err
	}
	resetReq := ctx.Value("request").(*user.ResetRequest)
	return resetReq.Close()
}

func (m *httpManager) Login(username string, pass string) (*user.UserDb, string, error) {
	currUser := user.GetUserByEmail(username)
	if currUser == nil {
		return nil, "", errors.New(constants.ItemNotFound)
	}
	token, err := currUser.Login(pass)
	if err != nil {
		return nil, "", err
	}
	return currUser, token, nil
}

func (m *httpManager) GetMe(ctx context.Context) (*user.UserDb, string, error) {
	currUser := ctx.Value("user").(*user.UserDb)
	data := map[string]string{
		"user_id": currUser.ID.Hex(),
	}
	token, err := util.CreateToken(data)
	if err != nil {
		return nil, "", err
	}
	return currUser, token, nil
}

// TODO: change to context
func (m *httpManager) DeleteUser(id string) error {
	return user.DeleteUserByID(id)
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
