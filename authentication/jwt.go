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
	"regexp"
	"github.com/cwg930/drones-server/models"
)

type AuthBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset = 3600
	unameExpr = "^[0-9A-Za-z_]+{3,16}$" //alphanumeric strings 3-16 chars
	passExpr = "^[0-9A-Za-z_!@#$%^&*()?+-]{8,32}$" //alphanumeric + symbols 8-32 chars
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

func (backend *AuthBackend) Authenticate(user *models.User) bool {
}

func (backend *AuthBackend) Register(username string, password string) (bool, error) {
	valid, err := regexp.MatchString(unameExpr, username)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}
	valid, err := regexp.MatchString(passExpr, password)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}
	hashedPass := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usr, err := models.DBConn.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	else if err == nil && usr != nil {
		return false, nil
	}
	else {
		usr := &models.User{Username: username, Password: hashedpass}
		err := models.DBconn.AddUser(usr)
		if err != nil {
			return false, err
		}
		else {
			return true, nil
		}
	}
}
