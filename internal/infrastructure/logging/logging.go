package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel string
	LogFile  string
}

func NewLog(cfg Config) (*logrus.Logger, error) {

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	logger.SetLevel(level)
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	return logger, nil
}
