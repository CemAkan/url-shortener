package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Log = logrus.New()

func InitLogger() {
	//first output target
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)

	//formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
	})

	//logs folder checker & creator
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0775)
	}
}
