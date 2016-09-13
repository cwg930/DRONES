package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	auth "github.com/cwg930/drones-server/authentication"
	"github.com/cwg930/drones-server/models"
)

func ListPlans(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	usr := context.Get(r, auth.UserKey)
	plans, err := db.AllPlansForUser(int(usr.(float64)))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(plans)
	if err != nil {
		log.Println(err) 
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ShowPlan(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	vars := mux.Vars(r)
	log.Println(vars)
	planId, err := strconv.ParseInt(vars["planId"], 10, 32)
	if err != nil {
		log.Println("1")
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	plan, err := db.GetPlan(int(planId))
	if err != nil { 
		log.Println("2")
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	err = json.NewEncoder(w).Encode(plan)
	if err != nil {
		log.Println("3")
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	log.Println("????")
	w.WriteHeader(http.StatusOK)
}

func CreatePlan(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	log.Printf("request body: %+v", r.Body)
	decoder := json.NewDecoder(r.Body)
	var p models.FlightPlan
	err := decoder.Decode(&p)
	owner := context.Get(r, auth.UserKey)
	log.Printf("plan received: %+v", p)
	p.OwnerID = int(owner.(float64))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	id, err := db.AddFlightPlan(p)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = db.AddAllPoints(int(id.(int64)), p.Points)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}
