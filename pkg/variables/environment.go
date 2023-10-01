package variables

import (
	"os"
)

type (
	EnvironmentValue string
)

const (
	EnvironmentProd    EnvironmentValue = "prod"
	EnvironmentStaging EnvironmentValue = "staging"
	EnvironmentLocal   EnvironmentValue = "local"
	EnvironmentDev     EnvironmentValue = "dev"
	EnvironmentTest    EnvironmentValue = "test"
)

var (
	allowedEnvironmentValues = map[EnvironmentValue]struct{}{
		EnvironmentProd:    {},
		EnvironmentStaging: {},
		EnvironmentDev:     {},
		EnvironmentLocal:   {},
		EnvironmentTest:    {},
	}
)

func AppEnvironment() EnvironmentValue {
	env, exists := os.LookupEnv(AppEnvironmentVariable.name)
	if !exists {
		return EnvironmentProd
	}

	value := EnvironmentValue(env)
	_, allowed := allowedEnvironmentValues[value]
	if !allowed {
		return EnvironmentProd
	}

	return value
}
