package helper

func CopySlice[T any](in []T) []T {
	tempSlice := make([]T, 0, len(in)+1)
	tempSlice = append(tempSlice, in...)

	return tempSlice
}

func CopyMap[T comparable, V any](in map[T]V) map[T]V {
	tempMap := make(map[T]V, len(in)+1)
	for key, value := range in {
		tempMap[key] = value
	}

	return tempMap
}
