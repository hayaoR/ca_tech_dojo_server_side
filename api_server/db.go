package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID   int64
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

	user.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Update() error {
	_, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
