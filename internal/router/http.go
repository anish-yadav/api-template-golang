package router

import (
	"net/http"

	"github.com/anish-yadav/api-template-golang/internal/auth"
	userHttpHandler "github.com/anish-yadav/api-template-golang/internal/user/handler"
	userUsecase "github.com/anish-yadav/api-template-golang/internal/user/usecase"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	userManager := userUsecase.NewHttpManager()
	userHandler := userHttpHandler.NewHttpHandler(userManager)

	// routes
	v1 := router.PathPrefix("/api/v1").Subrouter()
	publicPost := v1.Methods(http.MethodPost).Subrouter()
	publicGet := v1.Methods(http.MethodGet).Subrouter()

	publicGet.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("success")); err != nil {
			log.Errorf("failed to send health reposnse")
		}
	})

	publicPost.HandleFunc("/users/request-password-reset", userHandler.HandleRequestReset)
	publicPost.HandleFunc("/users/login", userHandler.HandleLoginRequest)

	//with user authenticated routes
	post := v1.Methods(http.MethodPost).Subrouter()
	post.Use(auth.VerifyResetMiddleware)
	post.HandleFunc("/users/reset-password", userHandler.HandleResetPassword)

	// jwt
	// permission authenticated routes
	v1 = v1.PathPrefix("/").Subrouter()
	v1.Use(auth.PermissionMiddleware)
	gets := v1.Methods(http.MethodGet).Subrouter()
	posts := v1.Methods(http.MethodPost).Subrouter()
	del := v1.Methods(http.MethodDelete).Subrouter()

	gets.HandleFunc("/users/me", userHandler.HandleGetMeRequest)

	posts.HandleFunc("/users", userHandler.handlePostUser)
	posts.HandleFunc("/users/change-password", userHandler.HandleChangePassword)

	del.HandleFunc("/users/{id}", userHandler.HandleUserDelete)
	return router
}
