package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	"log"
	"os"
	"io"

	"github.com/cwg930/imgapitest/models"
	"github.com/gorilla/mux"
)

type Env struct{
	db models.Datastore
}

var Envr Env

func InitEnv(db *models.DB) {
	Envr = Env{db}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("adduser.gtpl")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (env *Env) UserIndex(w http.ResponseWriter, r *http.Request) {
	usrs, err := env.db.AllUsers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(usrs)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	w.WriteHeader(http.StatusOK)
}

func (env *Env) ShowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["userId"], 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	usr, err := env.db.GetUser(int(userId))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	w.WriteHeader(http.StatusOK)
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

func SubmitIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("addfile.gtpl")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
		
func (env *Env) SubmitFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./files/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fMeta := models.FileMeta{FileName:"./files/" + handler.Filename}
	err = env.db.AddFile(fMeta)
	if err != nil {
		log.Printf("Error submitting file info to db: %v",err)
		http.Error(w, http.StatusText(500), 500)
		return
	}	
}

func (env *Env) ShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId, err := strconv.ParseInt(vars["fileId"], 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	fMeta, err := env.db.GetFile(int(fileId))
	http.ServeFile(w, r, fMeta.FileName)
}
