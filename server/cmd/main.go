package main

import (
	"os"

	"github.com/klynxe/word-of-wisdom/server/internal/config"
	"github.com/klynxe/word-of-wisdom/server/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	logger := initLogger()

	logger.Info("Starting Word of Wisdom Server...")

	app, err := server.New(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create server")

		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		logger.WithError(err).Fatal("Failed to run server")

		os.Exit(1)
	}
}

func initLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetOutput(os.Stdout)

	logger.SetLevel(logrus.InfoLevel)

	return logger
}
