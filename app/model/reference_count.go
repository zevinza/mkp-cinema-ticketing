package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ReferenceCount struct {
	Base
	DataOwner
	ReferenceCountAPI
}

type ReferenceCountAPI struct {
	Name   *string `json:"name,omitempty" example:"Transaction" gorm:"type:varchar(256);index:idx_reference_count_name_unique,unique,where:deleted_at is null;not null"` // Name
	Prefix *string `json:"prefix,omitempty" example:"INV"`                                                                                                               // Prefix
	Length *int    `json:"length,omitempty" example:"5"`                                                                                                                 // Length
	Count  *int64  `json:"count,omitempty" example:"1"`                                                                                                                  // Count
}

func (s *ReferenceCount) Seed() *[]ReferenceCount {
	contents := []string{
		"Transaction|INV|5",
		"Ticket|TC|8",
		"Product|PR|7",
	}

	c := []ReferenceCount{}
	for _, content := range contents {
		data := strings.Split(content, "|")

		length, _ := strconv.Atoi(data[2])
		count := int64(0)
		c = append(c, ReferenceCount{
			ReferenceCountAPI: ReferenceCountAPI{
				Name:   &data[0],
				Prefix: &data[1],
				Length: &length,
				Count:  &count,
			},
		})
	}
	return &c
}

func GenRefCount(name string, tx *gorm.DB) *string {
	ref := ReferenceCount{}
	if row := tx.Where(`name = ?`, name).Take(&ref); row.RowsAffected < 1 {
		return nil
	}
	if ref.Count == nil || ref.Prefix == nil || ref.Length == nil {
		return nil
	}

	count := fmt.Sprint(*ref.Count + 1)
	length := *ref.Length
	if len(count) > length {
		length = len(count) + 1
	}
	str := *ref.Prefix + "-" + strings.Repeat("0", length-len(count)) + count
	if name == "Transaction" {
		str = str + "-" + time.Now().Format("020106")
	}

	if err := tx.Model(&ref).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return nil
	}

	return &str
}
