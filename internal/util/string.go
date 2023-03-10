package util

import "github.com/google/uuid"

// GenerateUUID :nodoc:
func GenerateUUID() string {
	return uuid.New().String()
}
