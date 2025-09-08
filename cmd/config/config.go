package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config — основной конфиг приложения
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Log      interface{} // логгер оставляем для передачи извне
	LogLevel string      `env:"LOG_LEVEL" env-default:"info"`
	Env      string      `env:"ENV" env-default:"local"`
}

// ServerConfig — конфиг сервера
type ServerConfig struct {
	Host         string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port         string        `env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"120s"`
}

// DBConfig — конфиг базы данных PostgreSQL
type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"db"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"subscriptions_user"`
	Password string `env:"DB_PASSWORD" env-default:"subscriptions_pass"`
	Name     string `env:"DB_NAME" env-default:"subscriptions_db"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

// LoadConfig загружает конфиг из .env
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// DSN строит строку подключения к PostgreSQL
func (db *DBConfig) DSN() string {
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.Name +
		" sslmode=" + db.SSLMode
}
