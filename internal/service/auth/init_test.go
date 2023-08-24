package auth

import (
	"os"

	"github.com/severedsea/jwt-server/internal/pkg/jwt"
)

func init() {
	jwt.InitKeyFiles(os.Getenv("JWT_PRIVATE_KEY_PATH"), os.Getenv("JWT_PUBLIC_KEY_PATH"))
}
