package controller

import (
	"api/app/lib"
	"bufio"
	"bytes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// PostLogs func
func PostLogs(c *fiber.Ctx) error {
	readFileLogs := ""
	if viper.GetString("ENVIRONTMENT") == "production" || viper.GetString("ENVIRONTMENT") == "stagging" {
		file, err := os.Open("logger.log")
		if err != nil {
			lib.Logs.Println(err)
			lib.ErrorInternal(c, err.Error())
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() { // internally, it advances token based on sperator
			readFileLogs += scanner.Text() // token in unicode-char
			readFileLogs += "\n"
		}
	} else {
		file, err := os.Open("logs/logger.log")
		if err != nil {
			lib.Logs.Println(err)
			lib.ErrorInternal(c, err.Error())
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() { // internally, it advances token based on sperator
			readFileLogs += scanner.Text() // token in unicode-char
			readFileLogs += "\n"
		}
	}
	return c.SendStream(bufio.NewReader(bytes.NewReader([]byte(readFileLogs))))
}
