package variables

import "fmt"

const (
	VariableTypeEnvironment VariableType = "environment"
)

var registry []Variable

var (
	// AppEnvironmentVariable окружение, в котором запущено приложение
	AppEnvironmentVariable = Environment[string]("ENV", "prod")
)

type (
	Variable interface {
		Name() string
		Type() VariableType
	}

	VariableType string

	DefaultVariable[T any] struct {
		name         string
		defaultValue T
		t            VariableType
	}

	StringVariable = DefaultVariable[string]
	BoolVariable   = DefaultVariable[bool]
)

func (v DefaultVariable[T]) Name() string {
	return v.name
}

func (v DefaultVariable[T]) Type() VariableType {
	return v.t
}

func (v DefaultVariable[T]) String() string {
	return fmt.Sprintf("variable '%s', default '%v'", v.name, v.defaultValue)
}

func Environment[T any](name string, defaultValue T) DefaultVariable[T] {
	v := DefaultVariable[T]{name: name, defaultValue: defaultValue, t: VariableTypeEnvironment}
	register(v)

	return v
}

func register(v Variable) {
	registry = append(registry, v)
}
