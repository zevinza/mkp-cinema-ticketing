package lib

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/crypto/bcrypt"
)

func TestCipherEncryptDecrypt(t *testing.T) {
	plaintext := "password"
	key := "CIPHER_SECRETKEY_MUST_HAVE_32BIT"

	_, err := CipherEncrypt(plaintext, key[:28])
	// Invalid Key just have 28 byte in Encrypt
	utils.AssertEqual(t, fmt.Sprint("crypto/aes: invalid key size ", len(key[:28])), err.Error())

	cipherEncrypt, _ := CipherEncrypt(plaintext, key)
	cipherDecrypt, _ := CipherDecrypt(cipherEncrypt, key)
	// Success Decrypt
	utils.AssertEqual(t, plaintext, string(cipherDecrypt))

	_, err = CipherDecrypt(cipherEncrypt, key[:28])
	// Invalid Key just have 28 byte in Decrypt
	utils.AssertEqual(t, fmt.Sprint("crypto/aes: invalid key size ", len(key[:28])), err.Error())

	_, err = CipherDecrypt([]byte(string(cipherEncrypt)[:5]), key)
	// Len byte is different
	utils.AssertEqual(t, "ciphertext too short", err.Error())
}

func TestGeneratePassword(t *testing.T) {
	password := GeneratePassword(20, 6, 6, 6)

	utils.AssertEqual(t, 20, len(password))
}

func TestRandomChars(t *testing.T) {
	utils.AssertEqual(t, 10, len(RandomChars(10)))
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword returned an error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Errorf("Hashed password does not match the expected value")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Failed to generate hash: %v", err)
	}

	match := CheckPasswordHash(password, string(hashedPassword))
	if !match {
		t.Error("CheckPasswordHash failed to verify the password hash")
	}
}
