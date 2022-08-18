package utils

func Filter[T any](s []T, test func(T) bool) []T {
	var filtered []T
	for _, v := range s {
		if test(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Unique(s []string) []string {
	var unique []string
	for _, val := range s {
		if !Contains(unique, val) {
			unique = append(unique, val)
		}
	}
	return unique
}

func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func Any[T any](s []T, test func(T) bool) bool {
	for _, v := range s {
		if test(v) {
			return true
		}
	}
	return false
}

func All[T any](s []T, test func(T) bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, v := range s {
		if !test(v) {
			return false
		}
	}
	return true
}

func ZipMap(keys []string, values []string) map[string]string {
	if len(keys) != len(values) {
		panic("keys and values must be the same length")
	}

	m := make(map[string]string, len(keys))
	for i, key := range keys {
		m[key] = values[i]
	}
	return m
}
