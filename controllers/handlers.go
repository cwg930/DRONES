package controllers

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	"log"
	"os"
	"io"
	"io/ioutil"
	
	"github.com/cwg930/drones-server/models"
	"github.com/cwg930/drones-server/services"
	auth "github.com/cwg930/drones-server/authentication"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
)
/*
type Env struct{
	db models.Datastore
	secret string
}

var Envr Env

func InitEnv(db *models.DB, secret string) {
	Envr = Env{db, secret}
}
*/

var db *models.DB

func Init() error {
	var err error
	db, err = models.InitDB(string(os.Getenv("CONNECTION_STR")))
	if err != nil {
		return err
	}
	fmt.Printf("in controllers %+v", db)
	return nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.gtpl")
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


func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	log.Println("in login handler username = " + username)
	user := &models.User{Username: username, Password: password}
	log.Println("in login handler user.Username = " + user.Username)
	responseStatus, token := services.Login(user)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(responseStatus)
	log.Println("in login handler " + string(token))
	w.Write(token)
}

func UserIndex(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	usrs, err := db.AllUsers()
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

func ShowUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["userId"], 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	usr, err := db.GetUser(int(userId))
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	authBackend := auth.InitAuthBackend()
	log.Printf("\nIn CreateUser username=%v", username)
	success, err := authBackend.Register(username, password)
	if err != nil || !success {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		log.Println(success)
	}
	w.WriteHeader(http.StatusOK)
}

func SubmitIndex(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		
func SubmitFile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	usr := context.Get(r, auth.UserKey)
	reportID, err := strconv.Atoi(r.FormValue("reportID"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	pointID, err := strconv.Atoi(r.FormValue("pointID"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fMeta := models.FileMeta{FileName:"./files/" + handler.Filename, OwnerID: int(usr.(float64)), ReportID: reportID, PointID: pointID}
	err = db.AddFile(fMeta)
	if err != nil {
		log.Printf("Error submitting file info to db: %v",err)
		http.Error(w, http.StatusText(500), 500)
		return
	}	
}

func ListFiles(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	usr := context.Get(r, auth.UserKey)
	files, err := db.AllFiles(int(usr.(float64)))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func ShowFile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	fileId, err := strconv.ParseInt(vars["fileId"], 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	fMeta, err := db.GetFile(int(fileId))
	http.ServeFile(w, r, fMeta.FileName)
}

func SendFileBase64(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	vars := mux.Vars(r)
	fileId, err := strconv.ParseInt(vars["fileId"], 10,32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fMeta, err := db.GetFile(int(fileId))
	file, err := ioutil.ReadFile(fMeta.FileName)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	enc := base64.NewEncoder(base64.StdEncoding, w)
	enc.Write(file)
}

func GetToken(w http.ResponseWriter, r *http.Request){

}
