package config

// Environment variables
var Environment = map[string]interface{}{
	"app_id":           "6A6EF734-F557-4F7E-9C2A-32CD28E43420",
	"app_version":      "v1.0.0",
	"app_name":         "unknown",
	"app_description":  "",
	"port":             8000,
	"timezone":         "UTC",
	"enable_migration": true, // always set true this value for production usage
	"endpoint":         "/api/v1",
	"environment":      "development",
	"db_host":          "postgres",
	"db_port":          5432,
	"db_user":          "postgres",
	"db_pass":          "postgres",
	"db_name":          "postgres",
	"db_table_prefix":  "",
	"redis_host":       "redis",
	"redis_port":       6379,
	"redis_pass":       "",
	"redis_index":      0,
	"prefork":          false,
	"json_engine":      "sonic",                            // available options: sonic, go-json, encoding/json
	"aes":              "mB53IvZupVcalBnlEPmLyl4xJD4YN6g4", // AES 256-bit must have key at least 32 byte
	"salt":             "salt",
	"header_token_key": "x-Token",
	"value_token_key":  "v0x37KYbJqKodL0393Xa6jXaRTTN2eD1",
	"token_key":        "mkp_cinema",
	"token_expire_in":  24,
	"access_token":     "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFjNjdkYWFkLTVhODYtNDBlNi1hNzlhLTQ5ZGVhYWE2ZDg0YSIsImJ1c2luZXNzX2lkIjoiNmExYzEyZmQtNzEzNC00YTVhLTkwYmQtZGM1YzEyNmZkNThhIiwiYXVkIjoibXktY2xpZW50LWlkIiwiZXhwIjoxNjg2Njc3NzkyLCJzdWIiOiIxYzY3ZGFhZC01YTg2LTQwZTYtYTc5YS00OWRlYWFhNmQ4NGEifQ.-qbG3YQn6WsOvn5dEzIVCgUfA_wXmXjQcgWeBoK62KUyFtvyiOw5dAN9zWmZcnhBr0jMpME4nVVKlrP3gDvi0A",
	"user_id":          "1c67daad-5a86-40e6-a79a-49deaaa6d84a",
	"cart_duration":    600,
	"cart_ticker":      3,
	"admin_fee":        4000,
}
