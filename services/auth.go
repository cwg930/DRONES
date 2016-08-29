package services

import (
	"encoding/json"
	"net/http"
	"log"
	auth "github.com/cwg930/drones-server/authentication"
	"github.com/cwg930/drones-server/models"
)

func Login(user *models.User) (int, []byte){
	authBackend := auth.InitAuthBackend()
	log.Println("in login service username = " + user.Username)
	if authBackend.Authenticate(user) {
		token, err := authBackend.GenerateToken(user.ID)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		}
		response, err := json.Marshal(auth.Token{token})
		log.Println("in login service")
		log.Println(err)
		return http.StatusOK, response
	}
	
	return http.StatusUnauthorized, []byte("")
}
