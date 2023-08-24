package envvar

import "os"

// Mock mocks the environment by setting the provided and returns a func to revert to previous value
//
// Sample:
// defer envvar.Mock("SOME_ENV_VAR", "value")()
//
func Mock(key, value string) func() {
	oldValue, ok := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if ok {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	}
}
