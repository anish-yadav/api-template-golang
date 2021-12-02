package webservice

import (
	"net/http"

	"github.com/anish-yadav/api-template-golang/internal/router"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

func StartServer(port string) {
	r := router.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handlers.CORS(headersOk, originsOk, methodsOk)(r),
	}

	log.Infof("webservice.StartServer: server starting at port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Debugf("webservice.StartServer: %s", err.Error())
	}
}
