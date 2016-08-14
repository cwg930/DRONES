package authentication

import (
	jwt "github.com/dgrijalva/jwt-go"
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset = 3600
)

var authBackendInstance *AuthBackend = nil

func InitAuthBackend() *AuthBackend {
	if authBackendInstance = nil {
		authBackendInstance = &AuthBackend{
			privateKey: getPrivateKey()
			PublicKey: getPublicKey()
		}
	}
	return authBackendInstance
}

func (backend *AuthBackend) GenerateToken(userUUID string) (string, error){
	token := jwt.New(jwt.SigningMethodRS512)
	
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userUUID
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

