package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestCurrentTime(t *testing.T) {
	timeWithFormat := CurrentTime("2006-01-02 15:04:05")
	timeDefault := CurrentTime("")
	utils.AssertEqual(t, timeWithFormat, timeDefault)
}

func TestTimeNow(t *testing.T) {
	TimeNow()
}
