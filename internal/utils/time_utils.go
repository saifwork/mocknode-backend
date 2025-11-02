package utils

import (
	"time"
)

// GetExpiryTimestamp returns current UTC time + duration (in hours)
func GetExpiryTimestamp(hours int) time.Time {
	return time.Now().UTC().Add(time.Duration(hours) * time.Hour)
}

func FormatTimeISO(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

// RemainingTTL returns remaining duration until expiry
func RemainingTTL(expiry time.Time) time.Duration {
	return time.Until(expiry)
}
