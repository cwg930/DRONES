package routers

import (
//	"net/http"
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"github.com/cwg930/drones-server/controllers"
	"github.com/cwg930/drones-server/authentication"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.Index).Methods("GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.Handle("/users",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.UserIndex),
		)).Methods("GET")
	router.Handle("/users/{userId}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ShowUser),
		)).Methods("GET")
	router.Handle("/files", 
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ListFiles),
		)).Methods("GET")
	router.Handle("/files",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.SubmitFile),
		)).Methods("POST")
	router.Handle("/files/{fileId}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ShowFile),
		)).Methods("GET")
	return router
}
