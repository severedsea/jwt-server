package v1

import (
	"context"

	"github.com/severedsea/jwt-server/internal/service/auth"
)

var _ AuthService = (*auth.Service)(nil)

type AuthService interface {
	Login(ctx context.Context, authCode string) (auth.Token, error)
	Logout(ctx context.Context, userIDNo string) error
}
