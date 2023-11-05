package variables

import (
	"context"
	"fmt"
	"os"
	"time"
)

type (
	ToggleRepository interface {
		GetBool(ctx context.Context, t DefaultVariable[bool]) bool
		GetString(ctx context.Context, v DefaultVariable[string]) string
	}

	DefaultRepository struct{}
)

func NewDefaultRepository() *DefaultRepository {
	return &DefaultRepository{}
}

func (d *DefaultRepository) GetEnvironment() EnvironmentValue {
	result, err := getConverted[string](AppEnvironmentVariable, String)
	if err != nil {
		return EnvironmentValue(AppEnvironmentVariable.defaultValue)
	}

	envVal := EnvironmentValue(result)
	if _, exists := allowedEnvironmentValues[envVal]; !exists {
		return EnvironmentValue(AppEnvironmentVariable.defaultValue)
	}
	return envVal
}

func (d *DefaultRepository) GetBool(v DefaultVariable[bool]) bool {
	if v.Type() != VariableTypeEnvironment {
		return v.defaultValue
	}

	env, ok := os.LookupEnv(v.name)
	if !ok {
		return v.defaultValue
	}

	extracted, err := Bool.Extract(env)
	if err != nil {
		return v.defaultValue
	}

	return extracted
}

func (d *DefaultRepository) GetInt64(v DefaultVariable[string]) int64 {
	result, err := getConverted[int64](v, Int64)
	if err != nil {
		return 0
	}

	return result
}

func (d *DefaultRepository) GetInt64s(v DefaultVariable[string]) []int64 {
	result, err := getConverted[[]int64](v, Int64s)
	if err != nil {
		return nil
	}

	return result
}

func (d *DefaultRepository) GetString(v DefaultVariable[string]) string {
	result, err := getConverted[string](v, String)
	if err != nil {
		return ""
	}

	return result
}

func (d *DefaultRepository) GetStrings(v DefaultVariable[string]) []string {
	result, err := getConverted[[]string](v, Strings)
	if err != nil {
		return nil
	}

	return result
}

func (d *DefaultRepository) GetDuration(v DefaultVariable[string]) time.Duration {
	result, err := getConverted[time.Duration](v, Duration)
	if err != nil {
		return 0
	}

	return result
}

func (d *DefaultRepository) GetConvertedDuration(
	v DefaultVariable[string],
) (time.Duration, error) {
	return getConverted[time.Duration](v, Duration)
}

func getConverted[T any](v DefaultVariable[string], e extractor[T]) (T, error) {
	var (
		t   T
		raw string
		ok  bool
	)

	switch v.Type() {
	case VariableTypeEnvironment:
		raw, ok = os.LookupEnv(v.Name())
	default:
		return t, fmt.Errorf("variable type '%s' is unsupproted", v.Type())
	}

	if !ok {
		raw = v.defaultValue
	}

	converted, err := e.Extract(raw)
	if err != nil {
		fallbackValue, fallbackErr := e.Extract(v.defaultValue)
		if fallbackErr != nil {
			return t, fmt.Errorf("%s: %w", fallbackErr, err)
		}

		return fallbackValue, err
	}

	return converted, nil
}
