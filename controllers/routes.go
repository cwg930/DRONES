package controllers

import (
	"net/http"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"UserIndex",
		"GET",
		"/users",
		Envr.UserIndex,
	},
	Route{
		"ShowUser",
		"GET",
		"/users/{userId}",
		Envr.ShowUser,
	},
	Route{
		"CreateUser",
		"POST",
		"/users",
		Envr.CreateUser,
	},
	Route{
		"SubmitIndex",
		"GET",
		"/files",
		SubmitIndex,
	},
	Route{
		"SubmitFile",
		"POST",
		"/files",
		Envr.SubmitFile,
	},
	Route{
		"ShowFile",
		"GET",
		"/files/{fileId}",
		Envr.ShowFile,
	},
}
