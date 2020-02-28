package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id   int64
	Name string
}

var db *sql.DB

func (user *User) Create(db *sql.DB) (err error) {
	stmtIns, err := db.Prepare("INSERT INTO users VALUES( 0, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(user.Name)
	if err != nil {
		return err
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return
}

func (user *User) Update(db *sql.DB) (err error) {
	_, err = db.Exec("update users set name = ? where id = ?", user.Name, user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return
}

