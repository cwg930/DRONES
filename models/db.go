package models

import (
	"fmt"
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

var dbInstance *DB

func InitDB(connectionStr string) (*DB, error) {
	if dbInstance == nil {
		db, err := sql.Open("mysql", connectionStr)
		if err != nil {
			return nil, err
		}

		if err = db.Ping(); err != nil {
			return nil, err
		}
		dbInstance = &DB{db}
		fmt.Printf("inside if %+v\n", dbInstance)
	}
	fmt.Printf("outside if %+v\n", dbInstance)
	return dbInstance, nil
}
