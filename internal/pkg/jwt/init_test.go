package jwt

import (
	"os"
)

func init() {
	InitKeyFiles(os.Getenv("JWT_PRIVATE_KEY_PATH"), os.Getenv("JWT_PUBLIC_KEY_PATH"))
}
