package lib

import (
	"log"
	"os"
	"sync"
	"time"
)

// Logs Global variable
var Logs = RecordLog("SYSTEMS -")

type logger struct {
	filename string
	*log.Logger
}

var logg *logger
var once sync.Once
var user string

// SetUser func
func SetUser(username string) string {
	user = username
	return user
}

// GetUser func
func GetUser() string {
	return user
}

// RecordLog func
func RecordLog(userlog string) *logger {
	SetUser(userlog)
	once.Do(func() {
		logg = CreateLogger("logs/logger.log", GetUser())
	})
	return logg
}

// CreateLogger func
func CreateLogger(fname string, user string) *logger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return &logger{
		filename: fname,
		Logger:   log.New(file, time.Now().Format(time.RFC3339)+" "+GetUser()+" ", log.Lshortfile),
	}
}

// Debug func
func Debug(cmd string) bool {
	RecordLog(GetUser()).Println(cmd)
	return true
}
