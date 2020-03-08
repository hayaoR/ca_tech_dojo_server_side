package main

import (
	"database/sql"
)

type User struct {
	ID   int
	Name string
}

var db *sql.DB

func (user *User) Create() error {
	stmtIns, err := db.Prepare("INSERT INTO users VALUES( 0, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(user.Name)
	if err != nil {
		return err
	}

	var ID int64
	ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(ID)
	return nil
}

func (user *User) Update() error {
	_, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Get(id int) error {
	if err := db.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name); err != nil {
		return err
	}
	return nil
}
