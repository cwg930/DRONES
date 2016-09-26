package models

import "log"

type Report struct{
	Name string `json:"name"`
	ID int `json:"id"`
	OwnerID int `json:"owner"`
	PlanID int `json:"flightplan"`
}

type Entry struct{
	FileName string `json:"filename"`
	ID int `json:"id"`
}
