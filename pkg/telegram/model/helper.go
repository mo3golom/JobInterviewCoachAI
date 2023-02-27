package model

import "strings"

const (
	dataSep            = ":"
	commandWithDataSep = "#"
)

func CommandWithDataToString(command string, data []string) string {
	return strings.Join([]string{
		command,
		DataToString(data),
	},
		commandWithDataSep,
	)
}

func StringToCommandWithData(in string) (string, []string) {
	split := strings.Split(in, commandWithDataSep)
	if len(split) == 0 {
		return "", nil
	}
	if len(split) == 1 {
		return split[0], nil
	}

	return split[0], StringToData(split[1])
}

func DataToString(data []string) string {
	return strings.Join(data, dataSep)
}

func StringToData(in string) []string {
	return strings.Split(in, dataSep)
}
