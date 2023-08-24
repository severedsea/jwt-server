package middleware

import "net/http"

// Adapter is the function signature to chain middleware handlers
type Adapter func(next http.Handler) http.Handler

// Adapt chains the provided list of adapters to the provided handler
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
