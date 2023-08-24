package main

import (
	"fmt"
	"os"

	"github.com/severedsea/golang-kit/envvar"
	"github.com/severedsea/golang-kit/web/server"
	"github.com/severedsea/jwt-server/cmd/serverd/banner"
	"github.com/severedsea/jwt-server/cmd/serverd/router"
	"github.com/severedsea/jwt-server/internal/pkg/jwt"
)

func main() {
	banner.Print()

	// Validate envvars
	envvarValidate()

	// Auth
	jwt.InitKeyFiles(
		os.Getenv("JWT_PRIVATE_KEY_PATH"),
		os.Getenv("JWT_PUBLIC_KEY_PATH"),
	)

	// Start server
	s := server.New(fmt.Sprintf(":%s", envvar.Get("PORT", "3000")), router.Handler())
	s.Start()
}

func envvarValidate() {
	envvar.ValidateEitherNotEmpty("JWT_PUBLIC_KEY_PATH", "JWT_PUBLIC_KEY")
	envvar.ValidateEitherNotEmpty("JWT_PRIVATE_KEY_PATH", "JWT_PRIVATE_KEY")
}
