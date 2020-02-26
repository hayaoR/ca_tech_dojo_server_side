package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id   int64
	Name string
}

func (user *User) Create(db *sql.DB) (err error) {
	stmtIns, err := db.Prepare("INSERT INTO users VALUES( 0, ? )") // ? = placeholder
	if err != nil {
		fmt.Println("Prepare")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(user.Name)
	if err != nil {
		fmt.Println("Exec")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		fmt.Println("Id")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return
}

func GetPost(id int64, db *sql.DB) (user User, err error) {
	user = User{}
	err = db.QueryRow("select id, name from users where id=?", id).Scan(&user.Name)
	return
}

func main() {
	db, err := sql.Open("mysql", "root:mysql@([localhost]:3306)/tech_dojo")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connection is success")
	user := User{Name: "koki honda"}

	fmt.Println(user)
	user.Create(db)
	fmt.Println(user)

	readUser, _ := GetPost(user.Id, db)
	fmt.Println(readUser)

}
