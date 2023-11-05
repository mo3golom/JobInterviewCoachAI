package variables

import (
	"errors"
	"fmt"
	"job-interviewer/pkg/structs"
	"os"
	"strings"
)

var (
	errEmptyVariableName = errors.New("variable NameValue is empty")
)

type (
	Configuration struct {
		Repository structs.Singleton[Repository]
	}
)

func NewConfiguration() (*Configuration, error) {
	if err := validate(); err != nil {
		return nil, fmt.Errorf("failed to validate variables: %w", err)
	}

	return &Configuration{
		Repository: structs.NewSingleton(func() (Repository, error) {
			return NewDefaultRepository(), nil
		}),
	}, nil
}

func validate() error {
	for _, v := range registry {
		name := strings.TrimSpace(v.Name())
		if name == "" {
			return errEmptyVariableName
		}

		if v.Type() == VariableTypeEnvironment {
			_, ok := os.LookupEnv(name)
			if !ok {
				return fmt.Errorf("environment variable '%s' is absent", name)
			}
		}
	}

	return nil
}
