package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// VALIDATOR validate request body
var VALIDATOR *validator.Validate = validator.New()

type Auth struct {
	UserID  *uuid.UUID `json:"user_id,omitempty"`
	IsAdmin *bool      `json:"is_admin,omitempty"`
	Exp     *int64     `json:"exp,omitempty"`
}

// init Register custom validation function
func init() {
	VALIDATOR.RegisterValidation("phone", validatePhone)
	VALIDATOR.RegisterValidation("email", validateEmail)
	VALIDATOR.RegisterValidation("website", validateWebsite)
	VALIDATOR.RegisterValidation("emptyString", validateEmptyString)
	VALIDATOR.RegisterValidation("noWhiteSpace", validateNoWhiteSpace)
}

// Custom validation function for phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phoneRegex := regexp.MustCompile(`^\d{10,12}$`)
	return phoneRegex.MatchString(fl.Field().String())
}

// Custom validation function for email format
func validateEmail(fl validator.FieldLevel) bool {
	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return emailRegex.MatchString(fl.Field().String())
}

// Custom validation function for website format
func validateWebsite(fl validator.FieldLevel) bool {
	websiteRegex := regexp.MustCompile(`^(http|https):\/\/[^\s/$.?#].[^\s]*$`)
	return websiteRegex.MatchString(fl.Field().String())
}

// Custom validation function for empty string
func validateEmptyString(fl validator.FieldLevel) bool {
	emptyString := regexp.MustCompile(`^\s*$`)
	return !emptyString.MatchString(fl.Field().String())
}

// Custom validation function for no white space
func validateNoWhiteSpace(fl validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), " ") &&
		!strings.HasPrefix(fl.Field().String(), " ") &&
		!strings.HasSuffix(fl.Field().String(), " ")
}

// ClaimsJWT func
func ClaimsJWT(accesToken *string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(*accesToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	timeNow := time.Now().Unix()
	timeSessions := int64(claims["exp"].(float64))
	if timeSessions < timeNow {
		return claims, err
	}
	return claims, nil
}

// GetXUserID provide user id from the authentication token
func GetXUserID(c *fiber.Ctx) *uuid.UUID {
	authData := c.Locals("auth").(jwt.MapClaims)
	if authData == nil {
		return nil
	}

	js, err := json.Marshal(authData)
	if nil != err {
		return nil
	}

	auth := Auth{}
	err = json.Unmarshal(js, &auth)
	if err != nil {
		return nil
	}

	return auth.UserID
}

func GetXIsAdmin(c *fiber.Ctx) bool {
	authData := c.Locals("auth").(jwt.MapClaims)
	if authData == nil {
		return false
	}

	js, err := json.Marshal(authData)
	if nil != err {
		return false
	}

	auth := Auth{}
	err = json.Unmarshal(js, &auth)
	if err != nil {
		return false
	}

	if auth.IsAdmin != nil {
		return *auth.IsAdmin
	}
	return false
}

// BodyParser with validation
func BodyParser(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); nil != err {
		return err
	}
	return VALIDATOR.Struct(payload)
}

// GetImageData retrieves the image data from the specified URL
func GetImageData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func ParamsID(c *fiber.Ctx) *uuid.UUID {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return nil
	}
	return &id
}
