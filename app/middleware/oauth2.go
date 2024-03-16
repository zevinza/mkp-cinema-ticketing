package middleware

import (
	"api/app/lib"
	"api/app/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func Oauth2Authentication(c *fiber.Ctx) error {
	// do something here

	auth := c.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer") {
		return lib.ErrorUnauthorized(c, "Bearer token is needed")
	}

	claims, err := GetClaimsFromToken(strings.TrimPrefix(auth, "Bearer "))

	if err != nil {
		return lib.ErrorUnauthorized(c, "2"+err.Error())
	}

	c.Locals("auth", claims)

	return c.Next()
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// GenerateAccessToken func
func GenerateAccessToken(user *model.User) (model.ResponseToken, error) {
	result := model.ResponseToken{}
	exp := viper.GetInt("TOKEN_EXPIRE_IN")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.ID.String(),
			"is_admin": lib.RevBool(user.IsAdmin),
		})

	tokenString, err := token.SignedString([]byte(viper.GetString("TOKEN_KEY")))
	if err != nil {
		return result, err
	}

	return model.ResponseToken{
		AccessToken: &tokenString,
		ExpiresIn:   lib.Int64ptr(time.Now().Add(time.Hour * time.Duration(int64(exp))).Unix()),
		TokenType:   lib.Strptr("Bearer"),
	}, nil

}
