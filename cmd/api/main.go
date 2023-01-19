package main

import (
	"github.com/sirupsen/logrus"
)

type application struct {
	config *Config
	logger *logrus.Logger
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
