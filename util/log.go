package util

import (
	"github.com/momper14/rotatefilehook"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	"github.com/sirupsen/logrus"
	"time"
)

func InitLogger() {
	var logLevel = logrus.InfoLevel

	// if config is debug, will set latter
	//if debug {
	//	logLevel = logrus.DebugLevel
	//}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
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
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logrus.AddHook(rotateFileHook)
}
