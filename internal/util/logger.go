package util

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogger(level string) {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	configureLevel(level)
}

func configureLevel(level string) {
	lvl, err := log.ParseLevel(level)

	if err == nil {
		log.SetLevel(lvl)
	} else {
		log.Errorf("error during parsing log level: %s", err)
	}

	log.Debugf("Log Level is: %s", lvl.String())
}