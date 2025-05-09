package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetOutput(os.Stdout)

	Log.SetLevel(logrus.InfoLevel)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
	})
}
