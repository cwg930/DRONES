package authentication

import (
	jwt "github.com/dgrijalva/jwt-go"
//	"bufio"
	"crypto/rsa"
//	"crypto/x509"
//	"encoding/pem"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"regexp"
	"errors"
	"fmt"
	"github.com/cwg930/drones-server/models"
)

type AuthBackend struct {
	HMACSecret []byte
}

type Token struct {
	Token string `json:"token"`
}

const (
	tokenDuration = 72
	expireOffset = 3600
	unameExpr = "^[0-9A-Za-z_]{3,16}$" //alphanumeric strings 3-16 chars
	passExpr = "^[0-9A-Za-z_!@#$%^&*()?+-]{8,32}$" //alphanumeric + symbols 8-32 chars
)

var authBackendInstance *AuthBackend = nil

var ErrNotFound = errors.New("auth: user not found")
var ErrUserExists = errors.New("auth: user already exists")
var ErrInvalidPass = errors.New("auth: password not valid")
var ErrInvalidUser = errors.New("auth: username not valid")

func InitAuthBackend() *AuthBackend {
	if authBackendInstance == nil {
		authBackendInstance = &AuthBackend{[]byte(os.Getenv("AUTH_SECRET"))}
	}
	return authBackendInstance
}

func (backend *AuthBackend) GenerateToken(userID int) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":time.Now().Unix(),
		"sub":userID,
	})

	tokenString, err := token.SignedString(backend.HMACSecret)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func (backend *AuthBackend) Authenticate(user *models.User) bool {
	db, err := models.InitDB(string(os.Getenv("CONNECTION_STR")))
	if err != nil {
		return false
	}
	fmt.Println("In authenticate username = " + user.Username)
	foundUser, err := db.GetUserByUsername(user.Username)
	if err != nil || foundUser == nil {
		return false
	}
	if user.Username == foundUser.Username && bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)) == nil {
		user.ID = foundUser.ID
		return true
	}
	return false
}

func (backend *AuthBackend) Register(username string, password string) (bool, error) {
	valid, err := regexp.MatchString(unameExpr, username)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, ErrInvalidUser
	}
	valid, err = regexp.MatchString(passExpr, password)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, ErrInvalidPass
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hashedPass))
	if err != nil {
		return false, err
	}
	db, err := models.InitDB(string(os.Getenv("CONNECTION_STR")))
	if err != nil {
		return false, err
	}
	usr, err := db.GetUserByUsername(username)
	if err != nil {
		return false, err
	}else if err == nil && usr != nil {
		return false, ErrUserExists
	}else {
		usr := models.User{Username: username, Password: string(hashedPass)}
		err := db.AddUser(usr)
		if err != nil {
			return false, err
		}else {
			return true, nil
		}
	}
}

//NYI
func getPrivateKey() *rsa.PrivateKey {
	return nil
}
//NYI
func getPublicKey() *rsa.PublicKey {
	return nil
}
