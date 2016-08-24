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

var DBConn DB

func InitDB(connectionStr string) error {
	DBConn, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return err
	}
	if err := DBConn.Ping(); err != nil {
		return err
	}
	return nil
}
