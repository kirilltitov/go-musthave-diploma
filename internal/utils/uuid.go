package utils

import "github.com/google/uuid"

func NewUUID6() uuid.UUID {
	result, _ := uuid.NewV6()
	return result
}
