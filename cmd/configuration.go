package cmd

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"job-interviewer/pkg/logger"
	"os"
)

func MustInitDB(ctx context.Context) *sqlx.DB {
	envType := os.Getenv("ENV")

	sslModeValue := "require"
	if envType != "prod" {
		sslModeValue = "disable"
	}

	source := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		sslModeValue,
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
	envType := os.Getenv("ENV")
	sentryDsn := os.Getenv("SENTRY_DSN")
	log, err := logger.NewLogger(envType, sentryDsn)
	if err != nil {
		panic(err)
	}

	return log
}
