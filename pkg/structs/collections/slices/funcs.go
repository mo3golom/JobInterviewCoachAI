package slices

import (
	"errors"
)

func Filter[T any](in []T, predicate func(T) bool) []T {
	if len(in) == 0 {
		return nil
	}

	results := make([]T, 0, len(in))

	for _, element := range in {
		if predicate(element) {
			results = append(results, element)
		}
	}

	return results
}

var ErrSingleExpectedButWasNone = errors.New("expected single element, but was none")
var ErrSingleExpectedButWasMore = errors.New("expected single element, but was more")

func Single[T any](in []T, predicate func(T) bool) (T, error) {
	var t T
	var result *T

	for _, element := range in {
		element := element

		if predicate(element) {
			if result != nil {
				return t, ErrSingleExpectedButWasMore
			}

			result = &element
		}
	}

	if result == nil {
		return t, ErrSingleExpectedButWasNone
	}

	return *result, nil
}

func Contains[T any](in []T, predicate func(T) bool) bool {
	for _, element := range in {
		if predicate(element) {
			return true
		}
	}

	return false
}

func ContainsValue[T comparable](in []T, value T) bool {
	return Contains(in, func(t T) bool {
		return t == value
	})
}

func UniqueByProperty[T any, K comparable](in []T, extractor func(T) K) []T {
	if len(in) == 0 {
		return nil
	}

	unique := make(map[K]struct{}, len(in))
	results := make([]T, 0, len(in))

	for _, element := range in {
		property := extractor(element)
		if _, exists := unique[property]; exists {
			continue
		}

		unique[property] = struct{}{}
		results = append(results, element)
	}

	return results
}
