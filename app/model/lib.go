/*
this file contains local library for package model
all func should be limited for this package only
*/
package model

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
)

// boolInt func convert 0 and 1 any into false and true boolean
func boolInt(b interface{}) *bool {
	var t, f bool = true, false
	s := fmt.Sprint(b)
	if s == "1" {
		return &t
	}
	return &f
}

func now() *strfmt.DateTime {
	now := strfmt.DateTime(time.Now())
	return &now
}

func strptr(s string) *string { return &s }
