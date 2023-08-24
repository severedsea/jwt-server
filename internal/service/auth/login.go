package auth

import (
	"context"

	"github.com/severedsea/golang-kit/web"
)

// Login creates a session and generates an access_token based on the subject provided
func (s Service) Login(ctx context.Context, subject string) (Token, error) {

	t, err := s.GenerateToken(ctx, subject)
	if err != nil {
		return Token{}, web.WithStack(err)
	}

	return t, nil
}
