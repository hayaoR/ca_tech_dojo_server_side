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

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("can't decode")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//register in DB
	if err := user.Create(); err != nil {
		log.Println(err.Error())
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	if err != nil {
		log.Println(err.Error())
		return
	}
	tokenJSON, err := json.Marshal(Token{tokenString})

	if err != nil {
		log.Println(err.Error())
		return
	}

	//fmt.Printf(tokenString)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(tokenJSON)

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
		log.Println(err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := User{}
		id := int64(claims["id"].(float64))

		if err := user.Get(id); err != nil {
			log.Println(err.Error())
			return
		}
		tokenJSON, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(tokenJSON)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var tmpUser User
	if err := json.NewDecoder(r.Body).Decode(&tmpUser); err != nil {
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
		log.Println(err.Error())
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["id"], claims["nbf"])
		user := User{ID: int64(claims["id"].(float64)), Name: tmpUser.Name}

		if err := user.Update(); err != nil {
			fmt.Println("failed to update")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	} else {
		log.Println("token not valid")
	}
}

func execute() error {
	var err error
	db, err = sql.Open("mysql", "root:mysql@([localhost]:3306)/tech_dojo")
	if err != nil {
		return err
	}
	defer db.Close()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/user/create", GetTokenHandler)
	http.HandleFunc("/user/get", GetNameHandler)
	http.HandleFunc("/user/update", UpdateHandler)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}
