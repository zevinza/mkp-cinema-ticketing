package lib

import (
	"testing"
)

func TestGenUUIDString(t *testing.T) {
	GenUUIDString()
}
func TestStringToUUID(t *testing.T) {
	StringToUUID("0f02aec5-b741-4d3e-8fb4-87ac4961a495")
}

func TestGenUUID(t *testing.T) {
	GenUUID()
	NewUUID()
}
