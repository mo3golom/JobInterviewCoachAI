package logger

import (
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

type DefaultLogger struct {
}

func NewLogger(env string, sentryDSN string) (*DefaultLogger, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Environment: env,
	})
	if err != nil {
		return nil, err
	}

	return &DefaultLogger{}, nil
}

func (l *DefaultLogger) Info(msg string) {
	log.Println("INFO:", msg)
	sentry.CaptureMessage(msg)
}

func (l *DefaultLogger) Error(msg string, err error) {
	log.Println("ERROR:", msg, "EXT:", err.Error())
	sentry.CaptureException(err)
}

func (l *DefaultLogger) Flush() {
	sentry.Recover()
	sentry.Flush(time.Second * 5)
}
