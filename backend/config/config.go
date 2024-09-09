package config

import (
	"avito/log"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DbUsername string `validate:"required" body:"db_username"`
	DbPassword string `validate:"required" body:"db_password"`
	DbHost     string `validate:"required" body:"db_host"`
	DbPort     string `validate:"required" body:"db_port"`
	DbName     string `validate:"required" body:"db_name"`
	DbSchema   string `validate:"required" body:"db_schema"`
	DbSSL      string `validate:"required" body:"db_ssl"`
	AppMode    string `validate:"required" body:"app_mode"`
}

func (c *Config) WithSchema(schema string) *Config {
	c.DbSchema = schema
	return c
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSSL)
}

func LoadFromEnv() *Config {
	_ = godotenv.Load("dev.env")
	return &Config{
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		DbSSL:      os.Getenv("DB_SSL"),
		AppMode:    os.Getenv("APP_MODE"),
		DbSchema:   os.Getenv("DB_SCHEMA"),
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
		return ctx, fmt.Errorf("invalid configuration")
	}
	return ctx, nil
}
