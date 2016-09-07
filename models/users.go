package models

import (
	"log"
	"database/sql"
)
type User struct{
	Username string `json:"username"`
	Password string `json:"password"`
	ID int `json:"id"`
}

func (db *DB) AllUsers() ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	usrs := make([]*User, 0)
	for rows.Next() {
		usr := new(User)
		err := rows.Scan(&usr.ID, &usr.Username, &usr.Password)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, usr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return usrs, nil
}

func (db *DB) GetUser(id int) (*User, error) {
	usr := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&usr.ID, &usr.Username, &usr.Password)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (db *DB) AddUser(usr User) error {
	stmt, err := db.Prepare("INSERT INTO users(username,password) VALUES(?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(usr.Username, usr.Password)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Added user: %s with ID %d. Rows affected: %d", usr.Username, lastID, rowCnt)
	return nil
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	usr := &User{}
	log.Println("In GetUserByUsername, username = " + username)
	err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&usr.ID, &usr.Username, &usr.Password)
	switch{
	case err == sql.ErrNoRows:
		log.Println("no rows")
		return nil, nil
	case err != nil: 
		log.Println(err)
		return nil, err
	default:
		log.Printf("%+v\n", usr)
		return usr, nil 
	}
}
