package common

import "net/url"

func Of[T any](v T) *T {
	return &v
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func RemoveDuplicates[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	encountered := make(map[K]bool)
	var result []T

	for _, v := range slice {
		key := keyFunc(v)
		if !encountered[key] {
			encountered[key] = true
			result = append(result, v)
		}
	}

	return result
}
