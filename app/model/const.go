package model

import "github.com/google/uuid"

const (
	movie1    string = "0f02aec5-b741-4d3e-8fb4-87ac4961a401"
	movie2    string = "0f02aec5-b741-4d3e-8fb4-87ad4961a402"
	location1 string = "0f02aec5-b741-4e3e-8fb4-87ad4961a403"
	location2 string = "0f02aec5-a741-4d3e-8fb4-87ad4961a404"
	theater1  string = "0f02aec5-a741-4d3e-8fb4-87ad4961a405"
	theater2  string = "0f02aec5-a741-4d3e-8fb4-87ad4961a406"
	theater3  string = "0f02aec5-a741-4d3e-8fb4-87ad4961a407"
	theater4  string = "0f02aec5-a741-4d3e-8fb4-87ad4961a408"
)

func ParseStrToUUID(str string) *uuid.UUID {
	u, err := uuid.Parse(str)
	if err != nil {
		return nil
	}
	return &u
}
