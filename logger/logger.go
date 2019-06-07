package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Logger ...
var Logger = logrus.New()

// Configure ...
func Configure(wd string, file string, level string) {
	// Default logging to standard output
	Logger.SetOutput(os.Stdout)

	// Format
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	formatter.ForceColors = true
	Logger.Formatter = formatter

	// Loglevel (from configuration file)
	l, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.Errorf("Failed to parse log level from configuration: %v", err)
	} else {
		Logger.Level = l
	}

	Logger.Println(Logger.Level)

	// Log File (from configuration file)
	if file != "" {
		logAbsFile := filepath.Join(GetCurrentPath(), wd, file)

		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   logAbsFile,
			MaxSize:    50,
			MaxBackups: 3,
			MaxAge:     7,
			Level:      logrus.DebugLevel,
			Formatter: &prefixed.TextFormatter{
				ForceColors:   false,
				DisableColors: true},
		})

		if err != nil {
			Logger.Errorf("Failed to initialize file rotate hook: %v", err)
		}

		Logger.AddHook(rotateFileHook)
	}
}

// Get ...
func Get(prefix string) *logrus.Entry {
	l := Logger.WithFields(logrus.Fields{"prefix": prefix})
	return l
}

// GetCurrentPath ...
func GetCurrentPath() string {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return ex
}
