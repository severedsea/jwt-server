package jwt

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/severedsea/golang-kit/projectpath"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// InitKeys initialised the auth keys
func InitKeys(privKeyBytes []byte, pubKeyBytes []byte) {
	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "private key"))
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "public key"))
	}
}

// InitKeyFiles initialises the auth keys from provided file paths
func InitKeyFiles(privKeyPath string, pubKeyPath string) {
	privKeyBytes, err := os.ReadFile(projectpath.Abs(privKeyPath))
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "private key"))
	}
	pubKeyBytes, err := os.ReadFile(projectpath.Abs(pubKeyPath))
	if err != nil {
		log.Fatalf("%s", errors.Wrap(err, "public key"))
	}

	InitKeys(privKeyBytes, pubKeyBytes)
}
