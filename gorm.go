package web_bootstrap

import (
	"time"

	"github.com/gouef/gorm"
)

type GormDatabaseLoggerConfig struct {
	SlowThreshold             time.Duration `json:"slow_threshold" yaml:"slow_threshold" mapstructure:"slow_threshold"`
	LogLevel                  string        `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
	IgnoreRecordNotFoundError bool          `json:"ignore_record_not_found_error" yaml:"ignore_record_not_found_error" mapstructure:"ignore_record_not_found_error"`
	ParameterizedQueries      bool          `json:"parameterized_queries" yaml:"parameterized_queries" mapstructure:"parameterized_queries"`
	Colorful                  bool          `json:"colorful" yaml:"colorful" mapstructure:"colorful"`
}

type GormDatabaseConfig struct {
	Driver          string                   `json:"driver" yaml:"driver" mapstructure:"driver"`
	Host            string                   `json:"host" yaml:"host" mapstructure:"host"`
	Port            int                      `json:"port" yaml:"port" mapstructure:"port"`
	User            string                   `json:"user" yaml:"user" mapstructure:"user"`
	Password        string                   `json:"password" yaml:"password" mapstructure:"password"`
	DBName          string                   `json:"dbname" yaml:"dbname" mapstructure:"dbname"`
	SSLMode         string                   `json:"sslmode" yaml:"sslmode" mapstructure:"sslmode"`
	TimeZone        string                   `json:"timezone" yaml:"timezone" mapstructure:"timezone"`
	MaxIdleConns    int                      `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int                      `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration            `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	Debug           bool                     `json:"debug" yaml:"debug" mapstructure:"debug"`
	Logger          GormDatabaseLoggerConfig `json:"logger" yaml:"logger" mapstructure:"logger"`
}

func (db GormDatabaseConfig) ToGormConfig() *gorm.Config {
	return &gorm.Config{
		Driver:          db.Driver,
		Host:            db.Host,
		Port:            db.Port,
		User:            db.User,
		Password:        db.Password,
		DBName:          db.DBName,
		SSLMode:         db.SSLMode,
		TimeZone:        db.TimeZone,
		MaxIdleConns:    db.MaxIdleConns,
		MaxOpenConns:    db.MaxOpenConns,
		ConnMaxLifetime: db.ConnMaxLifetime,
		Debug:           db.Debug,
		Logger:          loggerConfigToGouefLoggerConfig(db.Logger),
	}
}

func loggerConfigToGouefLoggerConfig(cfg GormDatabaseLoggerConfig) gorm.LoggerConfig {
	return gorm.LoggerConfig{
		SlowThreshold:             cfg.SlowThreshold,
		LogLevel:                  cfg.LogLevel,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		Colorful:                  cfg.Colorful,
	}
}
