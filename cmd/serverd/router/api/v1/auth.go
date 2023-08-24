package v1

import (
	"net/http"

	"github.com/severedsea/golang-kit/web"
	"github.com/severedsea/jwt-server/internal/service/auth"
)

type AuthHandler struct {
	auth AuthService
}

func NewAuthHandler(a AuthService) AuthHandler {
	return AuthHandler{
		auth: a,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Login will generate an access_token for the provided subject and return as a session cookie
func (h AuthHandler) Login() http.HandlerFunc {
	return web.WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		subject := r.URL.Query().Get("subject")
		token, err := h.auth.Login(ctx, subject)
		if err != nil {
			return err
		}

		// Set token as cookie for web clients
		http.SetCookie(w, token.Cookie())

		web.RespondJSON(ctx, w, TokenResponse{AccessToken: token.AccessToken}, nil)

		return nil
	})
}

// Verify will return 200 if the provided access_token from either header or cookie is valid
func (h AuthHandler) Verify() http.HandlerFunc {
	return web.WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		// Do nothing, let the middleware do the validation

		return nil
	})
}

// Logout invalidates the access_token and cookie
func (h AuthHandler) Logout() http.HandlerFunc {
	return web.WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		claims, err := auth.ClaimsFromContext(ctx)
		if err != nil {
			return err
		}

		if err := h.auth.Logout(ctx, claims.Subject); err != nil {
			return err
		}

		// Invalidate cookie
		auth.InvalidateCookie(w)

		return nil
	})
}
