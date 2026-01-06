package utils

import "time"

func NowPtr() *time.Time {
	now := time.Now()
	return &now
}

func Ptr[T any](v T) *T {
	return &v
}
