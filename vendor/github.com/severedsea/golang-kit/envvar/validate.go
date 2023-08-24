/*
Package envvar is a helper for managing application environment variables.

Validation helpers can be used to validate environment values' format in order to fail fast.
Sample usage:
```
package main

import (
	"github.com/severedsea/golang-kit/envvar"
)

func main() {
	// Validate env variables
	envvar.ValidateTimeF("SOME_DATE")
	envvar.ValidateNotEmptyF("SOME_IMPORTANT_CONFIG")
	envvar.ValidateDurationF("SOME_DURATION_LIKE_EXPIRY")

	// Do something here
}
```

*/
package envvar

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ValidateEitherNotEmpty validates whether if all the values in the environment variable keys is empty
func ValidateEitherNotEmpty(keys ...string) bool {
	var ok bool
	for _, k := range keys {
		ok = ok || (strings.TrimSpace(os.Getenv(k)) != "")
	}
	return ok
}

// ValidateNotEmpty validates whether if the value in the environment variable key is empty
func ValidateNotEmpty(key string) bool {
	return strings.TrimSpace(os.Getenv(key)) != ""
}

// ValidateDuration validates whether if the value in the environment variable key is not a valid Go duration
func ValidateDuration(key string) (bool, error) {
	if _, err := time.ParseDuration(strings.TrimSpace(os.Getenv(key))); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateTime validates whether if the value in the environment variable key is not a valid time.Time
func ValidateTime(key string, layout string) (bool, error) {
	if _, err := time.Parse(layout, os.Getenv(key)); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateJSON validates whether if the value in the environment variable key is not a valid JSON
func ValidateJSON(key string) (bool, error) {
	j := map[string]interface{}{}
	if err := json.Unmarshal([]byte(os.Getenv(key)), &j); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateJSONArray validates whether if the value in the environment variable key is not a valid JSON array
func ValidateJSONArray(key string) (bool, error) {
	jarr := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(os.Getenv(key)), &jarr); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateEitherNotEmptyF exits/stops the app if all the values in the environment variable keys is empty
func ValidateEitherNotEmptyF(keys ...string) {
	if ok := ValidateEitherNotEmpty(keys...); !ok {
		log.Fatalf("[env] %s - not found", strings.Join(keys, ","))
	}
}

// ValidateNotEmptyF exits/stops the app if the value in the environment variable key is empty
func ValidateNotEmptyF(key string) {
	if ok := ValidateNotEmpty(key); !ok {
		log.Fatalf("[env] %s - not found", key)
	}
}

// ValidateDurationF exits/stops the app if the value in the environment variable key is not a valid Go duration
func ValidateDurationF(key string) {
	if _, err := ValidateDuration(key); err != nil {
		log.Fatalf("[env] %s - invalid duration: %s", key, err.Error())
	}
}

// ValidateTimeF exits/stops the app if the value in the environment variable key is not a valid time.Time
func ValidateTimeF(key string, layout string) {
	if _, err := ValidateTime(key, layout); err != nil {
		log.Fatalf("[env] %s - invalid time/date %s format: %s", key, layout, err.Error())
	}
}

// ValidateJSONF exits/stops the app if the value in the environment variable key is not a valid JSON
func ValidateJSONF(key string) {
	if _, err := ValidateJSON(key); err != nil {
		log.Fatalf("[env] %s - invalid JSON %s", key, err.Error())
	}
}

// ValidateJSONArrayF exits/stops the app if the value in the environment variable key is not a valid JSON array
func ValidateJSONArrayF(key string) {
	if _, err := ValidateJSONArray(key); err != nil {
		log.Fatalf("[env] %s - invalid JSON array %s", key, err.Error())
	}
}

// ValidateInt validates whether if the value in the environment variable key is not a valid int
func ValidateInt(key string) (bool, error) {
	if _, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key))); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateIntF exits/stops the app if the value in the environment variable key is not a valid int
func ValidateIntF(key string) {
	if _, err := ValidateInt(key); err != nil {
		log.Fatalf("[env] %s - invalid int %s", key, err.Error())
	}
}

// ValidateFloat64 validates whether if the value in the environment variable key is not a valid int
func ValidateFloat64(key string) (bool, error) {
	if _, err := strconv.ParseFloat(strings.TrimSpace(os.Getenv(key)), 64); err != nil {
		return false, err
	}
	return true, nil
}

// ValidateFloat64F exits/stops the app if the value in the environment variable key is not a valid int
func ValidateFloat64F(key string) {
	if _, err := ValidateFloat64(key); err != nil {
		log.Fatalf("[env] %s - invalid float64 %s", key, err.Error())
	}
}
