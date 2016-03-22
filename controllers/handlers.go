package controllers

import (
//	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"strconv"

	"github.com/cwg930/imgapitest/models"
//	"github.com/gorilla/mux"
)

type Env struct{
	db models.Datastore
}

var Envr Env

func InitEnv(db *models.DB) {
	Envr = Env{db}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("adduser.gtpl")
	t.Execute(w, nil)
}

func (env *Env) UserIndex(w http.ResponseWriter, r *http.Request) {
	usrs, err := env.db.AllUsers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, usr := range usrs {
		fmt.Fprintf(w, "%d\tName: %s\tAge: %d\n", usr.ID, usr.Name, usr.Age)
	}
}

func (env *Env) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("username")
	age, err := strconv.ParseInt(r.Form.Get("age"), 10, 32)
	if err != nil {
		age = 0
	}
	usr := models.User{name, 0, int(age)}
	err = env.db.AddUser(usr)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}
		
