package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TH-takahirohara/reading_record_api/internal/data"
	"github.com/TH-takahirohara/reading_record_api/internal/mailer"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type application struct {
	config *Config
	logger *logrus.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPSender),
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)

	duration, err := time.ParseDuration(cfg.DBMaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
