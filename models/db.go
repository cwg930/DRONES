package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Datastore interface{
	AllUsers() ([]*User, error)
	GetUser(int) (*User, error)
	AddUser(User) (error)
	AllFiles() ([]*FileMeta, error)
	AddFile(FileMeta) (error)
	GetFile(int) (*FileMeta, error)
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
