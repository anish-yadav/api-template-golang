package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anish-yadav/api-template-golang/internal/constants"
	"github.com/anish-yadav/api-template-golang/internal/pkg/db"
	"github.com/anish-yadav/api-template-golang/internal/pkg/webservice"
	"github.com/anish-yadav/api-template-golang/internal/util"
	"github.com/google/uuid"
)

var (
	dbURI  = flag.String("dbAddr", "mongodb+srv://admin:admin@lms-cluster.mnoel.mongodb.net/test", "url of mongodb database")
	dbName = flag.String("db", "api-template-golang", "database name")
	port   = flag.String("port", "8080", "port of the server")
	log    = flag.String("log", "debug", "log level")
)

func init() {
	jwtSecret := os.Getenv(constants.JwtSecret)
	if len(jwtSecret) == 0 {
		jwtSecret = uuid.New().String()
		fmt.Println(jwtSecret)
		os.Setenv(constants.JwtSecret, jwtSecret)
	}
	osPort := os.Getenv("PORT")
	if len(osPort) != 0 {
		*port = osPort
	}
}

func main() {

	flag.Parse()
	db.Init(*dbURI, *dbName)
	util.InitLogger(*log)

	webservice.StartServer(*port)
}
