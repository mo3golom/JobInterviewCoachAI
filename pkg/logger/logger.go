package logger

import (
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type DefaultLogger struct {
	enableSentry bool
}

func NewLogger(env string, sentryDSN string) (DefaultLogger, error) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Environment: env,
	})

	return DefaultLogger{
		enableSentry: env != "local",
	}, err
}

func (l DefaultLogger) Info(msg string, fields ...Field) {
	fieldsMap := fieldsToMap(fields)

	log.WithFields(fieldsMap).Info(msg)

	if !l.enableSentry {
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)

		scope.SetExtras(fieldsMap)
		sentry.CaptureMessage(msg)
	})
}

func (l DefaultLogger) Error(msg string, err error, fields ...Field) {
	fieldsMap := fieldsToMap(fields)
	fieldsMap["error"] = err.Error()
	log.WithFields(fieldsMap).Error(msg)

	if !l.enableSentry {
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)

		scope.SetExtras(fieldsMap)
		sentry.CaptureMessage(msg)
	})
}

func (l DefaultLogger) Flush() {
	if !l.enableSentry {
		return
	}
	sentry.Recover()
	sentry.Flush(time.Second * 5)
}

func fieldsToMap(in []Field) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for _, field := range in {
		result[field.Key] = field.Value
	}

	return result
}
