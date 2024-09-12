//nolint:gochecknoglobals,errorlint,fatcontext,forcetypeassert
package config

import (
	"avito/domain"
	"avito/log"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DbUsername    string `validate:"required" body:"db_username"`
	DbPassword    string `validate:"required" body:"db_password"`
	DbHost        string `validate:"required" body:"db_host"`
	DbPort        string `validate:"required" body:"db_port"`
	DbName        string `validate:"required" body:"db_name"`
	DbSchema      string `validate:"required" body:"db_schema"`
	DbSSL         string `validate:"required" body:"db_ssl"`
	AppMode       string `validate:"required" body:"app_mode"`
	ServerAddress string `validate:"required" body:"server_address"`
	ConnString    string `validate:"required" body:"conn_string"`
	UrlJDBC       string `validate:"required" body:"jdbc_url"`
	MigrationDir  string `validate:"required" body:"migration_dir"`
}

func (c *Config) WithSchema(schema string) *Config {
	c.DbSchema = schema
	return c
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
}

func LoadFromEnv(envFilePath string) *Config {
	_ = godotenv.Load(envFilePath)
	return &Config{
		DbUsername:    getEnv("POSTGRES_USERNAME", ""),
		DbPassword:    getEnv("POSTGRES_PASSWORD", ""),
		DbHost:        getEnv("POSTGRES_HOST", ""),
		DbPort:        getEnv("POSTGRES_PORT", ""),
		DbName:        getEnv("POSTGRES_DATABASE", ""),
		DbSSL:         getEnv("POSTGRES_SSL", "disable"),
		AppMode:       getEnv("APP_MODE", "dev"),
		DbSchema:      getEnv("POSTGRES_SCHEMA", "postgres"),
		ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ConnString:    getEnv("POSTGRES_CONN", ""),
		UrlJDBC:       getEnv("POSTGRES_JDBC_URL", ""),
		MigrationDir:  getEnv("MIGRATIONS_DIR", "migrations"),
	}
}

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func (c *Config) Validate(ctx context.Context) (context.Context, error) {
	if err := validate.Struct(c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			ctx = log.AddField(ctx, log.Field{Name: "Field", Value: err.StructField()})
			ctx = log.AddField(ctx, log.Field{Name: "ActualTag", Value: err.ActualTag()})
			ctx = log.AddField(ctx, log.Field{Name: "Tag", Value: err.Tag()})
			ctx = log.AddField(ctx, log.Field{Name: "Param", Value: err.Param()})
			ctx = log.AddField(ctx, log.Field{Name: "Value", Value: err.Value()})
		}
		return ctx, domain.ErrInvalidConfig
	}

	return ctx, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
