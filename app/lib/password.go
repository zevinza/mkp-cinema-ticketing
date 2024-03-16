package lib

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// PasswordEncrypt Password Encrypt
func PasswordEncrypt(plain, salt, key string, cost ...int) string {
	if len(cost) == 0 {
		cost = append(cost, bcrypt.DefaultCost)
	}
	password := plain + "//" + salt + "//" + key
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost[0])
	encoded := ""
	if nil == err {
		encoded = string(hashed)
	}

	return encoded
}

// PasswordCompare Password Compare
func PasswordCompare(encrypted, plain, salt, key string) bool {
	password := plain + "//" + salt + "//" + key
	return nil == bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
}

func GenPassword(str string) string {
	return PasswordEncrypt(str, viper.GetString("SALT"), viper.GetString("AES"))
}
