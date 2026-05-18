package configuration

import (
	"time"

	"github.com/gouef/gorm"
)

type GormDatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	TimeZone        string        `mapstructure:"timezone"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	Debug           bool          `mapstructure:"debug"`
}

// ToGormConfig vyexportuje data do struktury z balíčku gouef/gorm
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
	}
}
