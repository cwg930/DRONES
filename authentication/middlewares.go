package authentication

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
	"net/http"
	"fmt"
)

type key int

const UserKey key = 0

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc){
	authBackend := InitAuthBackend()
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With, Authorization")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.HMACSecret, nil
		}
	})
	if err != nil{
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("\n%T   %v\n", claims["sub"], claims["sub"])
		context.Set(req, UserKey, claims["sub"])
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
	
