package utils

import "time"

func DerefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}