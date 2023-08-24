package envvar

import "os"

// Get returns the value using os.Getenv if it exist, else returns the default value provided
func Get(key string, defaultValue string) string {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue
	}
	return s
}
