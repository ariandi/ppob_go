package util

import (
	"github.com/momper14/rotatefilehook"
	"github.com/sirupsen/logrus"
	"time"
)

func InitLogger() {
	var logLevel = logrus.InfoLevel
	t := time.Now()
	logFIleName := "log_" + t.Format("20060102")

	// if config is debug, will set latter
	//if debug {
	//	logLevel = logrus.DebugLevel
	//}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/" + logFIleName + ".log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)
	//logrus.SetOutput(formatter.ColorableStdOut)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logrus.AddHook(rotateFileHook)
}
