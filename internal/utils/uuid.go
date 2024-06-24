package utils

import "github.com/google/uuid"

func NewUUID6() uuid.UUID {
	result, err := uuid.NewV6()
	if err != nil {
		panic(err)
	}
	return result
}
