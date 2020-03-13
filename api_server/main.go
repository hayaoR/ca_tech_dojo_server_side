package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Token struct {
	Token string
}

type tomlConfig struct {
	ServerURL      string
	SQLConfigParam string
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//register in DB
	if err := user.Create(); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	tokenJSON, err := json.Marshal(Token{tokenString})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tokenJSON)

}

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("x-token")

	token, err := auth(tokenString)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := User{}
		id := int(claims["id"].(float64))

		if err := user.Get(id); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenJSON, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	token, err := auth(tokenString)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := User{ID: int(claims["id"].(float64)), Name: tmpUser.Name}

		if err := user.Update(); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("token not valid")
	}
}

func DrawGachaHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("x-token")

	token, err := auth(tokenString)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ID int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ID = int(claims["id"].(float64))
	}

	characterList, err := GetProbabilityList()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ガチャをtime回引いて結果を返す。
	times := Time{}
	if err := json.Unmarshal(body, &times); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	results := Results{}
	for i := 0; i < times.Times; i++ {
		idx := DrawGacha(characterList)
		character := Character{CharacterID: strconv.Itoa(characterList[idx].ID), Name: characterList[idx].Name}
		results.Results = append(results.Results, character)
	}

	//getしたモンスターをデータベースに登録
	for _, v := range results.Results {
		//v.CharacterIDは数値を文字列に変換したものだから失敗しないはず
		characterID, _ := strconv.Atoi(v.CharacterID)
		poses := Posession{UserID: ID, CharacterID: characterID}
		if err := poses.RegistrateOwnership(); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	tokenJSON, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tokenJSON)

}

func GetCharactersHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("x-token")

	token, err := auth(tokenString)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.ID = int(claims["id"].(float64))
	}

	characters, err := user.GetCharacters()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenJSON, err := json.MarshalIndent(characters, "", "\t")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tokenJSON)
}

func execute() error {
	var config tomlConfig
	if _, err := toml.DecodeFile("setting/setting.toml", &config); err != nil {
		return err
	}

	var err error
	db, err = sql.Open("mysql", config.SQLConfigParam)
	if err != nil {
		return err
	}
	defer db.Close()

	server := http.Server{
		Addr: config.ServerURL,
	}
	r := mux.NewRouter()
	r.HandleFunc("/user/create", GetTokenHandler)
	r.HandleFunc("/user/get", GetNameHandler)
	r.HandleFunc("/user/update", UpdateHandler)

	r.HandleFunc("/gacha/draw", DrawGachaHandler)

	r.HandleFunc("/character/list", GetCharactersHandler)
	http.Handle("/", r)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func auth(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINGKEY")), nil
	})

	return token, err
}

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}
