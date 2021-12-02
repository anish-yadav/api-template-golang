package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/anish-yadav/api-template-golang/internal/constants"
	"github.com/anish-yadav/api-template-golang/internal/pkg/permission"
	"github.com/anish-yadav/api-template-golang/internal/pkg/user"
	webresponse "github.com/anish-yadav/api-template-golang/internal/pkg/webservice/response"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

const (
	authHeaderKey = "Authorization"
	bearer        = "Bearer "
)

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(authHeaderKey)
		tokenString := strings.TrimPrefix(token, bearer)
		tokenString = strings.Trim(tokenString, " ")
		secret := os.Getenv(constants.JwtSecret)
		log.Debugf("auth.PermissionMiddleware: token : %s", tokenString)
		if len(tokenString) == 0 {
			webresponse.RespondWithError(w, http.StatusBadRequest, constants.InvalidToken)
			return
		}
		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(constants.InvalidSigningMethod)
			}

			if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
				return nil, errors.New(constants.TokenExpired)
			}
			// todo need to autogenerate this secret on first run
			return []byte(secret), nil
		})
		if err != nil {
			log.Debugf("auth.PermissionMiddleware: failed to parse/verify : %s", err.Error())
			webresponse.RespondWithError(w, http.StatusUnauthorized, constants.InvalidToken)
			return
		}
		claims := t.Claims.(jwt.MapClaims)
		data := claims["data"].(map[string]interface{})
		path := strings.TrimPrefix(r.URL.Path, "/api/v1")
		reqPermission := PermissionMap[r.Method][path]
		log.Debugf("permission required for %s : %s", r.URL.Path, reqPermission)
		userID := data["user_id"]
		if userID == nil {
			log.Debugf("invalid token, userID: %s", userID)
			webresponse.RespondWithError(w, http.StatusForbidden, constants.Forbidden)
			return
		}
		currUser := user.GetUserById(fmt.Sprintf("%s", userID))
		userPermission := permission.GetPermissionByName(currUser.Type)

		if !userPermission.HasPermission(reqPermission) {
			log.Debugf("user is not authorized")
			webresponse.RespondWithError(w, http.StatusForbidden, constants.Forbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "user", currUser)
		r = r.WithContext(ctx)
		if t.Valid {
			next.ServeHTTP(w, r)
			return
		}

		webresponse.RespondWithError(w, http.StatusUnauthorized, constants.InvalidToken)
		return
	})
}

func VerifyResetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(authHeaderKey)
		tokenString := strings.TrimPrefix(token, bearer)
		tokenString = strings.Trim(tokenString, " ")
		log.Debugf("auth.PermissionMiddleware: token : %s", tokenString)
		secret := os.Getenv(constants.JwtSecret)
		if len(tokenString) == 0 {
			webresponse.RespondWithError(w, http.StatusBadRequest, constants.InvalidToken)
			return
		}
		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(constants.InvalidSigningMethod)
			}

			if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
				return nil, errors.New(constants.TokenExpired)
			}
			// todo need to autogenerate this secret on first run
			return []byte(secret), nil
		})
		if err != nil {
			log.Debugf("auth.PermissionMiddleware: failed to parse/verify : %s", err.Error())
			webresponse.RespondWithError(w, http.StatusUnauthorized, constants.InvalidToken)
			return
		}
		claims := t.Claims.(jwt.MapClaims)
		data := claims["data"].(map[string]interface{})
		path := strings.TrimPrefix(r.URL.Path, "/api/v1")
		reqPermission := PermissionMap[r.Method][path]
		log.Debugf("permission required for %s : %s", r.URL.Path, reqPermission)
		username := data["username"]
		tokenID := data["token_id"]

		if username == nil || tokenID == nil {
			webresponse.RespondWithError(w, http.StatusForbidden, constants.Forbidden)
			return
		}

		currUser := user.GetUserByEmail(fmt.Sprintf("%s", username))
		resetReq := user.GetReqById(fmt.Sprintf("%s", tokenID))
		userPermission := permission.GetPermissionByName(currUser.Type)

		if !userPermission.HasPermission(reqPermission) || !resetReq.IsValid() {
			webresponse.RespondWithError(w, http.StatusForbidden, constants.Forbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "user", currUser)
		ctx = context.WithValue(ctx, "request", resetReq)
		r = r.WithContext(ctx)

		if t.Valid {
			next.ServeHTTP(w, r)
			return
		}

		webresponse.RespondWithError(w, http.StatusUnauthorized, constants.InvalidToken)
		return
	})
}
