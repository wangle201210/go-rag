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
