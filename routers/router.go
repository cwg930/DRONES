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
//	router.HandleFunc("/", controllers.Index).Methods("GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router = setCORSRoutes(router)
	router = setFlightPlanRoutes(router)
	router = setReportRoutes(router)
	router = setFileRoutes(router)
	return router
}

func setCORSRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/files", controllers.HandleCORS).Methods("OPTIONS")
	router.HandleFunc("/flightplans", controllers.HandleCORS).Methods("OPTIONS")
	router.HandleFunc("/flightplans/{planId}", controllers.HandleCORS).Methods("OPTIONS")
	return router
}

func setFlightPlanRoutes(router *mux.Router) *mux.Router {
	router.Handle("/flightplans", 
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ListPlans),
		)).Methods("GET")
	router.Handle("/flightplans/{planId}", 
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ShowPlan),
		)).Methods("GET")
	router.Handle("/flightplans", 
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.CreatePlan),
		)).Methods("POST")
	return router
}

func setReportRoutes(router *mux.Router) *mux.Router {
	router.Handle("/reports",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.ListReports),
		)).Methods("GET")
	router.Handle("/reports/{reportId}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetReport),
		)).Methods("GET")
	router.Handle("/reports",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.CreateReport),
		)).Methods("POST")
	return router
}

func setFileRoutes(router *mux.Router) *mux.Router {
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
