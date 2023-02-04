package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env            string `env:"RR_ENV" envDefault:"development"`
	Port           int    `env:"RR_PORT" envDefault:"8080"`
	DBHost         string `env:"DB_HOST" envDefault:"db"`
	DBPort         int    `env:"DB_PORT" envDefault:"3306"`
	DBUser         string `env:"DB_USER" envDefault:"user"`
	DBPassword     string `env:"DB_PASSWORD" envDefault:"my-secret-pw"`
	DBName         string `env:"DB_NAME" envDefault:"reading_record"`
	DBMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	DBMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	DBMaxIdleTime  string `env:"DB_MAX_IDLE_TIME" envDefault:"15m"`
	SMTPHost       string `env:"SMTP_HOST" envDefault:"127.0.0.1"`
	SMTPPort       int    `env:"SMTP_PORT" envDefault:"22225"`
	SMTPUsername   string `env:"SMTP_USERNAME" envDefault:"smtpusername"`
	SMTPPassword   string `env:"SMTP_PASSWORD" envDefault:"smtppassword"`
	SMTPSender     string `env:"SMTP_SENDER" envDefault:"sender@example.com"`
	TrustedOrigins string `env:"TRUSTED_ORIGINS" envDefault:"http://localhost:3000"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
