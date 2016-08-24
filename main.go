package main

import (
	"log"
	"net/http"
	"os"
	//"io/ioutil"
	"github.com/joho/godotenv"
	"github.com/gorilla/handlers"
	"github.com/cwg930/drones-server/routers"
	"github.com/cwg930/drones-server/models"
	"github.com/cwg930/drones-server/controllers"
)

func main() {
	
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}
	connectStr := os.Getenv("CONNECTION_STR")
	secret := os.Getenv("AUTH_SECRET")
	err := models.InitDB(string(connectStr))
	if err != nil {
		log.Panic(err)
	}
	controllers.InitEnv(db, secret)
	router := routers.NewRouter()

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))
}
