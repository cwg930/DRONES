package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/context"
	auth "github.com/cwg930/drones-server/authentication"
)

func ListReports(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	usr := context.Get(r, auth.UserKey)
	reports, err := db.AllReportsForUser(int(usr.(float64)))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type","application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(reports)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetReport(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	vars := mux.Vars(r)
	reportId, err := strconv.ParseInt(vars["reportId"], 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	report, err := db.GetReport(int(reportId))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(report)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateReport(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	decoder := json.NewDecoder(r.Body)
	var report models.Report
	err := decoder.Decode(&report)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	owner := context.Get(r, auth.UserKey)
	report.OwnerID = int(owner.(float64))
	id, err := db.AddReport(report)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
