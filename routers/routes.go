package routers

import (
	"net/http"
	"github.com/cwg930/drones-server/controllers"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

//var Envr = controllers.Envr

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		controllers.Index,
	},
	Route{
		"UserIndex",
		"GET",
		"/users",
		controllers.Envr.UserIndex,
	},
	Route{
		"ShowUser",
		"GET",
		"/users/{userId}",
		controllers.Envr.ShowUser,
	},
	Route{
		"CreateUser",
		"POST",
		"/users",
		controllers.Envr.CreateUser,
	},
	Route{
		"FileIndex",
		"GET",
		"/files",
		controllers.Envr.ListFiles,
	},
	Route{
		"SubmitFile",
		"POST",
		"/files",
		controllers.Envr.SubmitFile,
	},
	Route{
		"ShowFile",
		"GET",
		"/files/{fileId}",
		controllers.Envr.ShowFile,
	},
}
