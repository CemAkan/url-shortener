package infrastructure

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var Log = logrus.New()
var ServerLogOutput io.Writer

func InitLogger() {

	//logs folder checker & creator
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0775)
	}

	// app.log file opener and creator if it not exist
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		Log.Warn("failed to open log file, defaulting open stdout only")
		Log.SetOutput(os.Stdout)
	} else {
		Log.SetOutput(io.MultiWriter(file, os.Stdout))
	}

	Log.SetLevel(logrus.InfoLevel)

	//formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     false,
	})

	serverFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Log.Warn("Failed to open server log file, use app log file")
		ServerLogOutput = file
	} else {
		ServerLogOutput = serverFile
	}

}
