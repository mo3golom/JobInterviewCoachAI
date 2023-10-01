package variables

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

var (
	String     = &stringExtractor{}
	Strings    = &stringsExtractor{}
	Int64      = &int64Extractor{}
	Int64s     = &int64sExtractor{}
	Minutes    = &minutesExtractor{}
	Bool       = &boolExtractor{}
	Duration   = &durationExtractor{}
	Dictionary = &dictionaryExtractor{}
)

type (
	extractor[T any] interface {
		Extract(value string) (T, error)
		TargetType() string
	}

	stringExtractor     struct{}
	int64Extractor      struct{}
	int64sExtractor     struct{}
	minutesExtractor    struct{}
	boolExtractor       struct{}
	stringsExtractor    struct{}
	durationExtractor   struct{}
	dictionaryExtractor struct{}
)

func (s *stringExtractor) Extract(value string) (string, error) {
	return value, nil
}

func (s *stringExtractor) TargetType() string {
	return "string"
}

func (i *int64Extractor) Extract(value string) (int64, error) {
	vInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return int64(vInt), nil
}

func (i *int64Extractor) TargetType() string {
	return "int64"
}

func (i *int64sExtractor) Extract(value string) ([]int64, error) {
	if value == "" {
		return nil, nil
	}

	intStrings := strings.Split(value, ",")
	result := make([]int64, 0, len(intStrings))
	for _, intStr := range intStrings {
		vInt, err := strconv.Atoi(intStr)
		if err != nil {
			return nil, err
		}
		result = append(result, int64(vInt))
	}

	return result, nil
}

func (i *int64sExtractor) TargetType() string {
	return "[]int64"
}

func (i *stringsExtractor) Extract(value string) ([]string, error) {
	if value == "" {
		return nil, nil
	}
	return strings.Split(value, ","), nil
}

func (i *stringsExtractor) TargetType() string {
	return "[]string"
}

func (m *minutesExtractor) Extract(value string) (time.Duration, error) {
	vInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return time.Duration(vInt) * time.Minute, nil
}

func (m *minutesExtractor) TargetType() string {
	return "minutes"
}

func (m *boolExtractor) Extract(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func (m *boolExtractor) TargetType() string {
	return "bool"
}

func (d *durationExtractor) Extract(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}

func (d *durationExtractor) TargetType() string {
	return "duration"
}

func (d *dictionaryExtractor) Extract(value string) (map[string]any, error) {
	content := make(map[string]any)
	err := json.Unmarshal([]byte(value), &content)
	if err != nil {
		return content, err
	}
	return content, nil
}

func (d *dictionaryExtractor) TargetType() string {
	return "map[string]any"
}
