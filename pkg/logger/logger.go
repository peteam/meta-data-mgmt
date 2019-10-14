package logger

import (
	"io"
	"os"
	"strings"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"github.com/sirupsen/logrus"
)

var BootstrapLogger *logrus.Logger
var AccessLogger *logrus.Logger
var Logger *logrus.Logger

func init() {
	initBootstrapLogger()
	initAccessLogger()
	initLogger()
}

func initBootstrapLogger() {
	BootstrapLogger = logrus.New()
	BootstrapLogger.Level = getLogLevel(config.Viper.GetString("logging.logfile.bootstrap.loglevel"))
	BootstrapLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	file, err := os.OpenFile(
		config.Viper.GetString("logging.logfile.bootstrap.path")+
			strings.ToLower(config.Viper.GetString("logging.logfile.bootstrap.name")),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		BootstrapLogger.SetOutput(os.Stdout)
	} else {
		BootstrapLogger.SetOutput(file)
	}
}

func initAccessLogger() {
	AccessLogger = logrus.New()
	AccessLogger.Level = getLogLevel(config.Viper.GetString("logging.logfile.access.loglevel"))
	AccessLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	file, err := os.OpenFile(
		config.Viper.GetString("logging.logfile.access.path")+
			strings.ToLower(config.Viper.GetString("logging.logfile.access.name")),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		AccessLogger.SetOutput(os.Stdout)
	} else {
		AccessLogger.SetOutput(file)
	}
}

func initLogger() {
	Logger = logrus.New()
	Logger.Level = getLogLevel(config.Viper.GetString("logging.logfile.service.loglevel"))
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	file, err := os.OpenFile(
		config.Viper.GetString("logging.logfile.service.path")+
			strings.ToLower(config.Viper.GetString("logging.logfile.service.name")),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.SetOutput(os.Stdout)
	} else {
		logwriter := io.MultiWriter(os.Stdout, file)
		Logger.SetOutput(logwriter)
	}
}

func getLogLevel(loglevel string) logrus.Level {
	switch loglevel {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.ErrorLevel
	}
}
