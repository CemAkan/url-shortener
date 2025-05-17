package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var Log = logrus.New()
var MainLogFile io.Writer

func InitLogger() {

	//logs folder checker & creator
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0775)
	}

	var err error

	// app.log file opener and creator if it not exist
	MainLogFile, err = FileOpener("app")

	if err != nil {
		Log.Warn("failed to open log file, defaulting open stdout only")
		Log.SetOutput(os.Stdout)
	} else {
		Log.SetOutput(io.MultiWriter(MainLogFile, os.Stdout))
	}

	Log.SetLevel(logrus.InfoLevel)

	//formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     false,
	})

	Log.SetLevel(logrus.InfoLevel)

	//formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     false,
	})

}

// SpecialLogger creates special *logrus.Logger objects with fileName and outputType (file|stdout) parameters
func SpecialLogger(fileName string, outputType string) *logrus.Logger {
	logger := logrus.New()
	var output io.Writer

	switch outputType {
	case "file":
		// Use mainLogFile as fallback
		output = MainLogFile

		// If fileName is given, try to open that file
		if fileName != "" {
			file, err := FileOpener(fileName)
			if err != nil {
				Log.WithError(err).Warnf("Failed to open %s, using default main log file", fileName)
			} else {
				output = file
			}
		}

	case "stdout":
		output = os.Stdout

	default:
		Log.Warnf("Unknown outputType '%s', falling back to mainLogFile", outputType)
		output = Log.Out
	}

	// Configure logger
	logger.SetOutput(output)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     false,
	})

	return logger
}

// FileOpener open if file is exist or not create a new one
func FileOpener(fileName string) (io.Writer, error) {
	logFilePath := fmt.Sprintf("logs/" + fileName + ".log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return nil, err
	}
	return file, nil
}
