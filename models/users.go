package models

import "log"

type User struct{
	Name string
	ID int
	Age int
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
		err := rows.Scan(&usr.ID, &usr.Name, &usr.Age)
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

func (db *DB) AddUser(usr User) error {
	stmt, err := db.Prepare("INSERT INTO users(name,age) VALUES(?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(usr.Name, usr.Age)
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
	log.Printf("Added user: %s with ID %d. Rows affected: %d", usr.Name, lastID, rowCnt)
	return nil
}
