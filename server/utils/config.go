package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadFromDefaultEnvFile loads env file from default location (typically .env in the
// root project directory)
func LoadFromDefaultEnvFile() error {
	return godotenv.Load()
}

// LoadFromEnvFile loads specific env files
func LoadFromEnvFile(envFiles []string) error {
	return godotenv.Load(envFiles...)
}

// GetOsValueAsArray returns key value split by spaces
func GetOsValueAsArray(key string) []string {
	return strings.Split(os.Getenv(key), " ")
}

// GetOsValue returns the raw value from the env var table
func GetOsValue(key string) string {
	return os.Getenv(key)
}

// GetOsValueAsBoolean returns a parsed boolean, panics otherwise
func GetOsValueAsBoolean(key string) bool {
	sval := GetOsValue(key)
	val, err := strconv.ParseBool(sval)

	if err != nil {
		panic(fmt.Sprintf("Invalid boolean value %v for key %v", sval, key))
	}
	return val
}

// GetOsValueAsInt32 returns a parsed number, panics otherwise
func GetOsValueAsInt32(key string) int32 {
	sval := GetOsValue(key)
	val, err := strconv.ParseInt(sval, 10, 32)

	if err != nil {
		panic(fmt.Sprintf("Invalid boolean value %v for key %v", sval, key))
	}
	return int32(val)
}

// GetOsValueAsDuration returns a duration from a string (e.g. 1h, 2h45m, etc)
func GetOsValueAsDuration(key string) time.Duration {
	sval := GetOsValue(key)
	val, err := time.ParseDuration(sval)

	if err != nil {
		panic(fmt.Sprintf("Invalid duration value %v for key %v", sval, key))
	}
	return val
}

// InitializeEnvVars checks to see if an env file is the source of secrets
func InitializeEnvVars() {
	// load config from env by default, use NO_LOAD_ENV_FILE to use supplied env
	if _, noLoadEnvFile := os.LookupEnv("NO_LOAD_ENV_FILE"); !noLoadEnvFile {
		if err := LoadFromDefaultEnvFile(); err != nil {
			panic("Unable to load .env file")
		}
	}
}
