package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Datastore interface{
	AllUsers() ([]*User, error)
	AddUser(User) (error)
}

type DB struct{
	*sql.DB
}

func NewDB(connectionStr string) (*DB, error) {
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
