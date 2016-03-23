package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"github.com/cwg930/imgapitest/controllers"
	"github.com/cwg930/imgapitest/models"
)

func main() {
	connectStr, err := ioutil.ReadFile("server.conf")
	if err != nil {
		log.Panic(err)
	}
	db, err := models.NewDB(string(connectStr))
	if err != nil {
		log.Panic(err)
	}
	controllers.InitEnv(db)
	router := controllers.NewRouter() 

	log.Fatal(http.ListenAndServe(":8080", router))
}
