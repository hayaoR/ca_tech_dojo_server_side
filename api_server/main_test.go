package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func TestE2E(t *testing.T) {
	var config tomlConfig
	if _, err := toml.DecodeFile("setting/setting.toml", &config); err != nil {
		t.Fatal(err)
	}

	var err error
	db, err = sql.Open("mysql", config.SQLConfigParam)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	serveMux := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/user/create", GetTokenHandler)
	r.HandleFunc("/user/get", GetNameHandler)
	r.HandleFunc("/user/update", UpdateHandler)
	serveMux.Handle("/", r)

	// post name
	name := "yurina hirate"
	jsonStr := []byte(`
	{
		"name": "yurina hirate"
	}
	`)

	writer, request := postName(jsonStr)

	serveMux.ServeHTTP(writer, request)

	token, err := getToken(writer)

	if err != nil {
		t.Fatal(err)
	}

	//get name
	writer, request = getName(*token)

	serveMux.ServeHTTP(writer, request)

	user, err := getUser(writer)

	if err != nil {
		t.Fatal(err)
	}

	if user.Name != name {
		t.Fatal("Wrong content, was expecting", name, "but got", user.Name)
	}

	// update name
	name = "nana komatu"
	jsonStr = []byte(`
	{
		"name": "nana komatu"
	}
	`)

	writer, request = updateUser(jsonStr, *token)

	serveMux.ServeHTTP(writer, request)

	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}

	//get name
	writer, request = getName(*token)

	serveMux.ServeHTTP(writer, request)

	user, err = getUser(writer)

	if err != nil {
		t.Fatal(err)
	}

	if user.Name != name {
		t.Fatal("Wrong content, was expecting", name, "but got", user.Name)
	}
}

func TestGacha(t *testing.T) {
	var config tomlConfig
	if _, err := toml.DecodeFile("setting/setting.toml", &config); err != nil {
		t.Fatal(err)
	}

	var err error
	db, err = sql.Open("mysql", config.SQLConfigParam)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	serveMux := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/user/create", GetTokenHandler)
	r.HandleFunc("/user/get", GetNameHandler)
	r.HandleFunc("/user/update", UpdateHandler)

	r.HandleFunc("/gacha/draw", DrawGachaHandler)

	r.HandleFunc("/character/list", GetCharactersHandler)
	serveMux.Handle("/", r)

	// post name
	jsonStr := []byte(`
	{
		"name": "yurina hirate"
	}
	`)

	writer, request := postName(jsonStr)

	serveMux.ServeHTTP(writer, request)

	token, err := getToken(writer)

	if err != nil {
		t.Fatal(err)
	}

	jsonStr = []byte(`
	{
		"name": "nana komatu"
	}
	`)

	writer, request = postName(jsonStr)

	serveMux.ServeHTTP(writer, request)

	token2, err := getToken(writer)

	if err != nil {
		t.Fatal(err)
	}

	// tokenで二回、token2で一回ガチャを引く
	jsonStr = []byte(`
	{
		"times": 2
	}
	`)
	writer, request = drawGacha(*token, jsonStr)
	serveMux.ServeHTTP(writer, request)
	results, err := getGacha(writer)

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results.Results {
		t.Log(v.CharacterID, v.Name)
	}
	writer, request = drawGacha(*token, jsonStr)
	serveMux.ServeHTTP(writer, request)
	results, err = getGacha(writer)

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results.Results {
		t.Log(v.CharacterID, v.Name)
	}

	writer, request = drawGacha(*token2, jsonStr)
	serveMux.ServeHTTP(writer, request)
	results, err = getGacha(writer)

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results.Results {
		t.Log(v.CharacterID, v.Name)
	}

	writer, request = requestCharacters(*token)
	serveMux.ServeHTTP(writer, request)
	characters, err := getCharacters(writer)

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range characters.Characters {
		t.Log(v.UserCharacterID, v.CharacterID, v.Name)
	}
}

func postName(jsonStr []byte) (*httptest.ResponseRecorder, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/create", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return writer, request

}

func getName(token Token) (*httptest.ResponseRecorder, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/get", nil)
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	return writer, request
}

func updateUser(jsonStr []byte, token Token) (*httptest.ResponseRecorder, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/update", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	return writer, request
}
func drawGacha(token Token, jsonStr []byte) (*httptest.ResponseRecorder, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/gacha/draw", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	return writer, request
}

func requestCharacters(token Token) (*httptest.ResponseRecorder, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/character/list", nil)
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	return writer, request
}

func getGacha(writer *httptest.ResponseRecorder) (*Results, error) {
	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status is not OK")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := Results{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func getToken(writer *httptest.ResponseRecorder) (*Token, error) {
	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status is not OK")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func getUser(writer *httptest.ResponseRecorder) (*User, error) {
	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status is not OK")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getCharacters(writer *httptest.ResponseRecorder) (*Characters, error) {
	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status is not OK")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	characters := Characters{}
	if err := json.Unmarshal(body, &characters); err != nil {
		return nil, err
	}

	return &characters, nil
}
