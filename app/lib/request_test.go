package lib

import (
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
)

func TestBodyParser(t *testing.T) {
	type sample struct {
		Name *string `validate:"required,gte=9"`
	}

	app := fiber.New()
	app.Post("/validate", func(c *fiber.Ctx) error {
		data := new(sample)
		return ErrorBadRequest(c, BodyParser(c, data))
	})

	res, body, err := PostTest(app, "/validate", nil, "")
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, false, nil == body)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, _, err = PostTest(app, "/validate", nil, `{"name":"john"}`)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 400, res.StatusCode)
}

func TestClaimsJWT(t *testing.T) {
	// Mock claims
	mockClaims := jwt.MapClaims{
		"exp": float64(time.Now().Unix() + 1000), // Set expiration time in the future
		// Add other mock claims if needed
	}

	// Create a new token with the mock claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mockClaims)
	tokenString, _ := token.SignedString([]byte("secret")) // Sign the token with a secret key

	// Call the ClaimsJWT function with the mock token
	claims, err := ClaimsJWT(&tokenString)

	// Assert that there was no error
	if err != nil {
		t.Errorf("ClaimsJWT returned an error: %v", err)
	}

	// Assert that the claims match the mock claims
	if !reflect.DeepEqual(claims, mockClaims) {
		t.Errorf("ClaimsJWT did not return the expected claims")
	}

	// Test case for claims that cannot be converted to MapClaims
	invalidTokenString := "invalid-token"
	_, err = ClaimsJWT(&invalidTokenString)
	if err == nil {
		t.Error("ClaimsJWT should return an error for a token with invalid claims")
	}

	// Test case for claims that cannot be converted to MapClaims
	tokenStringWithoutClaims := "token-without-claims"
	_, err = ClaimsJWT(&tokenStringWithoutClaims)
	if err == nil {
		t.Error("ClaimsJWT should return an error for a token without claims")
	}

	// Test case for expired token
	expiredClaims := jwt.MapClaims{
		"exp": float64(time.Now().Unix() - 1000), // Set expiration time in the past
		// Add other mock claims if needed
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte("secret"))
	_, err = ClaimsJWT(&expiredTokenString)
	if err != nil {
		t.Errorf("ClaimsJWT should not return an error for an expired token: %v", err)
	}
}
