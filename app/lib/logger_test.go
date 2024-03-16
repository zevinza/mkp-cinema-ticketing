package lib

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSetUser(t *testing.T) {
	username := "john_doe"
	result := SetUser(username)

	if result != username {
		t.Errorf("SetUser returned an incorrect value. Expected: %s, Got: %s", username, result)
	}
}

func TestGetUser(t *testing.T) {
	username := "john_doe"
	SetUser(username)

	result := GetUser()

	if result != username {
		t.Errorf("GetUser returned an incorrect value. Expected: %s, Got: %s", username, result)
	}
}

func TestRecordLog(t *testing.T) {
	logg := RecordLog("user1")

	if logg == nil {
		t.Error("RecordLog returned a nil logger")
	}

	if logg != nil {
		if logg.filename != "logs/logger.log" {
			t.Errorf("RecordLog returned a logger with incorrect filename. Expected: logs/logger.log, Got: %s", logg.filename)
		}
		if !strings.Contains(logg.Prefix(), "SYSTEMS -") {
			t.Errorf("RecordLog returned a logger with incorrect prefix. Expected: SYSTEMS -, Got: %s", logg.Prefix())
		}
	}
}

func TestCreateLogger(t *testing.T) {
	fname := "test_logger.log"
	user := "test_user"
	SetUser(user)

	// Ensure the file is closed and removed after the test completes
	defer func() {
		if err := os.Remove(fname); err != nil {
			t.Errorf("Failed to remove test logger file: %v", err)
		}
	}()

	logg := CreateLogger(fname, GetUser())

	if logg == nil {
		t.Error("CreateLogger returned a nil logger")
	}

	if logg != nil {
		if logg.filename != fname {
			t.Errorf("CreateLogger returned a logger with incorrect filename. Expected: %s, Got: %s", fname, logg.filename)
		}

		prefix := logg.Prefix()
		expectedPrefix := time.Now().Format(time.RFC3339) + " " + GetUser()
		if !strings.Contains(prefix, expectedPrefix) {
			t.Errorf("CreateLogger returned a logger with incorrect prefix. Expected: %s, Got: %s", expectedPrefix, prefix)
		}
	}
}

func TestDebug(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := ioutil.TempFile("", "test_logger_*.log")
	if err != nil {
		t.Fatalf("Failed to create temporary log file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Redirect the standard logger to the temporary log file
	log.SetOutput(tmpFile)

	// Set up the logger
	logg = &logger{
		filename: tmpFile.Name(),
		Logger:   log.New(tmpFile, "", 0),
	}

	cmd := "Debug message"
	result := Debug(cmd)

	if !result {
		t.Error("Debug function returned false, expected true")
	}

	// Read the log file content
	content, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Verify if the log message is present in the file
	if !strings.Contains(string(content), cmd) {
		t.Errorf("Log message not found in the log file. Expected: %s, File Content: %s", cmd, string(content))
	}
}
