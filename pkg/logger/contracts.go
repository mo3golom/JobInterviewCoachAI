package logger

type (
	Field struct {
		Key   string
		Value any
	}

	Logger interface {
		Info(msg string, fields ...Field)
		Error(msg string, err error, fields ...Field)
		Flush()
	}
)
