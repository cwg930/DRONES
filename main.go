package main

import (
	"log"
	"net/http"
	"github.com/cwg930/imgapitest/controllers"
	"github.com/cwg930/imgapitest/models"
)

func main() {
	db, err := models.NewDB("root:@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Panic(err)
	}
	controllers.InitEnv(db)
	router := controllers.NewRouter() 

	log.Fatal(http.ListenAndServe(":8080", router))
}
