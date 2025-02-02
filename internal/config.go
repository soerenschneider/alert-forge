package internal

import (
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var (
	DefaultSeverityMapping = map[string]int64{
		"fatal":    15,
		"critical": 14,
		"high":     12,
		"major":    10,
		"error":    8,
		"warning":  5,
		"low":      3,
		"info":     1,
	}
)

type ReportConfig struct {
	Enabled         bool         `yaml:"enabled"`
	Email           *EmailConfig `yaml:"email" envPrefix:"EMAIL_"`
	AwtrixInstances []string     `yaml:"awtrix_instances" validate:"omitempty,dive,http_url"`
	Timezone        string       `yaml:"timezone" env:"TIMEZONE" validate:"omitempty,timezone"`
	CronExpression  string       `yaml:"cron_expression" env:"CRON" validate:"omitempty,cron"`
}

type Config struct {
	AlertmanagerInstances  []string          `yaml:"alertmanager_instances" validate:"dive,http_url"`
	IgnoredAlertLabelPairs map[string]string `yaml:"blacklist"`
	Reports                ReportConfig      `yaml:"reports" envPrefix:"REPORTS_"`
	SqliteDatabaseFile     string            `yaml:"sqlite_db_file" validate:"filepath"`
	SeverityMapping        map[string]int64  `yaml:"severity_mapping"`
	HttpServer             HttpServer        `yaml:"http_server"`
	MetricsAddr            string            `yaml:"metrics_addr" validate:"omitempty,hostname_port"`
}

type HttpServer struct {
	Address string `yaml:"address" validate:"hostname_port"`
}

type EmailConfig struct {
	From         string   `yaml:"from" env:"FROM" validate:"omitempty,email"`
	FromFile     string   `yaml:"from_file" env:"FROM_FILE" validate:"omitempty,file"`
	To           []string `yaml:"to" env:"TO" validate:"omitempty,dive,email"`
	ToFile       string   `yaml:"to_file" env:"TO_FILE" validate:"omitempty,file"`
	Password     string   `yaml:"password" env:"PASSWORD"`
	PasswordFile string   `yaml:"password_file" env:"PASSWORD_FILE" validate:"omitempty,file"`
	Server       string   `yaml:"server" env:"SERVER" validate:"omitempty,hostname"`
	Port         int      `yaml:"port" env:"PORT" validate:"omitempty,gt=0,lt=65535"`
	Username     string   `yaml:"username" env:"USERNAME"`
	UsernameFile string   `yaml:"username_file" env:"USERNAME_FILE" validate:"omitempty,file"`
}

func (e *EmailConfig) GetFrom() (string, error) {
	if e.From != "" {
		return e.From, nil
	}

	data, err := os.ReadFile(e.FromFile)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e *EmailConfig) GetTo() ([]string, error) {
	if len(e.To) > 0 {
		return e.To, nil
	}

	data, err := os.ReadFile(e.ToFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), ","), nil
}

func (e *EmailConfig) GetPassword() (string, error) {
	if e.Password != "" {
		return e.Password, nil
	}

	data, err := os.ReadFile(e.PasswordFile)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e *EmailConfig) GetUsername() (string, error) {
	if e.Username != "" {
		return e.Username, nil
	}

	data, err := os.ReadFile(e.UsernameFile)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DefaultConfig() *Config {
	return &Config{
		SqliteDatabaseFile: "database.sql",
		SeverityMapping:    DefaultSeverityMapping,
		HttpServer:         HttpServer{Address: "0.0.0.0:8888"},
		MetricsAddr:        "0.0.0.0:9169",
	}
}

func ReadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	opts := env.Options{
		Prefix: "AM_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, err
}
