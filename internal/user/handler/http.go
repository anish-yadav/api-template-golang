package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anish-yadav/api-template-golang/internal/constants"
	webresponse "github.com/anish-yadav/api-template-golang/internal/pkg/webservice/response"
	"github.com/anish-yadav/api-template-golang/internal/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userHttpHandler struct {
	manager usecase.HttpManager
}

func NewHttpHandler(manager usecase.HttpManager) *userHttpHandler {
	return &userHttpHandler{manager}
}

func (u *userHttpHandler) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	var userRequest UserRequest
	err = json.Unmarshal(data, &userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	id, err := u.manager.CreateUser(userRequest.toUser())
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			webresponse.RespondWithError(w, http.StatusConflict, constants.Conflict)
			return
		}
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}
	response := UserRequestResponse{id}
	webresponse.RespondWithSuccess(w, http.StatusCreated, response)
	return
}

func (u *userHttpHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	var userRequest ChangePasswordRequest
	err = json.Unmarshal(data, &userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	err = u.manager.ChangePassword(r.Context(), userRequest.OldPassword, userRequest.NewPassword)
	if err != nil {
		if err.Error() == constants.ItemNotFound {
			webresponse.RespondWithError(w, http.StatusNotFound, constants.ItemNotFound)
			return
		}
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}
	response := UserRequestResponse{userRequest.ID}
	webresponse.RespondWithSuccess(w, http.StatusCreated, response)
	return
}

func (u *userHttpHandler) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Debugf("deleteing user: %s", id)

	if !primitive.IsValidObjectID(id) {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.ItemNotFound)
		return
	}

	err := u.manager.DeleteUser(id)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadGateway, constants.IntervalServerError)
		return
	}

	webresponse.RespondWithSuccess(w, http.StatusOK, []byte("success"))
	return

}

func (u *userHttpHandler) HandleRequestReset(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	var req ReqResetPasswordRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	err = u.manager.RequestPasswordReset(req.Email)
	if err != nil {
		if err.Error() == constants.ItemNotFound {
			webresponse.RespondWithError(w, http.StatusNotFound, constants.ItemNotFound)
			return
		}
		log.Debugf("HandleRequestReset: %s", err.Error())
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}
	webresponse.RespondWithSuccess(w, http.StatusOK, "")
	return
}

func (u *userHttpHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	var userRequest ResetPasswordRequest
	err = json.Unmarshal(data, &userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(userRequest)
	if err != nil {
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	err = u.manager.ResetPassword(r.Context(), userRequest.NewPassword)
	if err != nil {
		if err.Error() == constants.ItemNotFound {
			webresponse.RespondWithError(w, http.StatusNotFound, constants.ItemNotFound)
			return
		}
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}
	webresponse.RespondWithSuccess(w, http.StatusOK, "")
	return
}

func (u *userHttpHandler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debugf("malformed data %s", err.Error())
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	var loginReq LoginRequest
	err = json.Unmarshal(data, &loginReq)
	if err != nil {
		log.Debugf("%s", data)
		log.Debugf("malformed data : %s ", err.Error())
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(loginReq)
	if err != nil {
		log.Debugf("malformed data : %s ", err.Error())
		webresponse.RespondWithError(w, http.StatusBadRequest, constants.BadRequest)
		return
	}
	user, token, err := u.manager.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		log.Debugf("loginUser: %s", err.Error())
		if err.Error() == constants.ItemNotFound || err == bcrypt.ErrMismatchedHashAndPassword {
			webresponse.RespondWithError(w, http.StatusUnauthorized, constants.PasswordMismatch)
			return
		}
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}
	response := LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}
	webresponse.RespondWithSuccess(w, http.StatusOK, response)
	return
}

func (u *userHttpHandler) HandleGetMeRequest(w http.ResponseWriter, r *http.Request) {

	user, token, err := u.manager.GetMe(r.Context())
	if err != nil {
		webresponse.RespondWithError(w, http.StatusInternalServerError, constants.IntervalServerError)
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}
	webresponse.RespondWithSuccess(w, http.StatusOK, response)
	return
}
