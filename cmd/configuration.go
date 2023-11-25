package cmd

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"job-interviewer/pkg/logger"
	variables "job-interviewer/pkg/variables"
	"os"
)

func MustInitDB(ctx context.Context) *sqlx.DB {
	source := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		loadEnvValue("DB_USER"),
		loadEnvValue("DB_PASSWORD"),
		loadEnvValue("DB_HOST"),
		loadEnvValue("DB_PORT"),
		loadEnvValue("DB_NAME"),
		loadEnvValue("DB_SSL"),
	)

	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		source,
	)
	if err != nil {
		panic(err)
	}

	return db
}

func MustInitLogger() logger.Logger {
	sentryDsn := loadEnvValue("SENTRY_DSN")
	log, err := logger.NewLogger(string(variables.AppEnvironment()), sentryDsn)
	if err != nil {
		panic(err)
	}

	return log
}

func loadEnvValue(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("env value doesn't exists")
	}

	return value
}
