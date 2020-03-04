package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

type Token struct {
	Token string
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("can't decode")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//register in DB
	err = user.Create(db)
	if err != nil {
		fmt.Println("can't create")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	if err != nil {
		fmt.Println("tokenString Error")
		fmt.Println(err.Error())
		return
	}
	tokenJSON, err := json.Marshal(Token{tokenString})

	if err != nil {
		fmt.Printf("Token Error")
		fmt.Println(err.Error())
		return
	} else {
		//fmt.Printf(tokenString)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(tokenJSON)
	}
}

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("x-token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINGKEY")), nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["id"], claims["nbf"])
		user := User{}
		err := db.QueryRow("select id, name from users where id = ?", int64(claims["id"].(float64))).Scan(&user.ID, &user.Name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tokenJSON, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(tokenJSON)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var tmp User
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		fmt.Println("can't decode")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("x-token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINGKEY")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["id"], claims["nbf"])
		user := User{ID: int64(claims["id"].(float64)), Name: tmp.Name}

		err := user.Update(db)
		if err != nil {
			fmt.Println("failed to update")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	} else {
		fmt.Println("token not valid")
	}
}

func main() {

	var err error
	db, err = sql.Open("mysql", "root:mysql@([localhost]:3306)/tech_dojo")
	if err != nil {
		log.Fatal("unable to use data source name")
	}
	defer db.Close()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/user/create", GetTokenHandler)
	http.HandleFunc("/user/get", GetNameHandler)
	http.HandleFunc("/user/update", UpdateHandler)

	server.ListenAndServe()
}
