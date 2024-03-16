package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestBusiness(t *testing.T) {
	utils.AssertEqual(t, float64(100), ProfitPercent(50, 100))
}
