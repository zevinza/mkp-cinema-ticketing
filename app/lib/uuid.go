package lib

import (
	guuid "github.com/google/uuid"
)

// GenUUIDString func
func GenUUIDString() string {
	id := guuid.New().String()
	return id
}

// StringToUUID func
func StringToUUID(s string) *guuid.UUID {
	res, _ := guuid.Parse(s)
	return &res
}

// GenUUID func
func GenUUID() *guuid.UUID {
	id := guuid.New()
	return &id
}

// NewUUID func
func NewUUID() *guuid.UUID {
	id, _ := guuid.NewRandom()
	return &id
}
