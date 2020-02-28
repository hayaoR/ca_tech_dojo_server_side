package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id   int64
	Name string
}

type Token struct {
	Token string
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

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("can't decode")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(user.Name)
	//fmt.Println(os.Getenv("SIGNINGKEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"nbf":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	if err != nil {
		fmt.Println("tokenString Error")
		fmt.Println(err.Error())
	}
	tokenJSON, err := json.Marshal(Token{tokenString})

	if err != nil {
		fmt.Printf("Token Error")
		fmt.Println(err.Error())
	} else {
		//fmt.Printf(tokenString)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(tokenJSON)
	}
}

var JwtMiddleWare = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNINGKEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

var GetName = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Clame content \n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t %#v\n", k, v)
	}
})

func main() {
	/*
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
	*/

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	authname := JwtMiddleWare.Handler(GetName)
	http.HandleFunc("/user/create", GetTokenHandler)
	http.Handle("/user/get", authname)
	server.ListenAndServe()
}
