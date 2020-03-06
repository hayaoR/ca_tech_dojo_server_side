package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

	
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/create", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	serveMux.ServeHTTP(writer, request)

	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	
	token := Token{}
	if err := json.Unmarshal(body, &token); err != nil {
		t.Fatal(err)
	}

	//get name
	writer = httptest.NewRecorder()
	request = httptest.NewRequest("Get", "/user/get", nil)
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	serveMux.ServeHTTP(writer, request)

	resp = writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	user := User{}
	if err := json.Unmarshal(body, &user); err != nil {
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

	writer = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/user/update", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	serveMux.ServeHTTP(writer, request)

	resp = writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}

	//get name
	writer = httptest.NewRecorder()
	request = httptest.NewRequest("Get", "/user/get", nil)
	request.Header.Set("accept", "application/json")
	request.Header.Set("x-token", token.Token)

	serveMux.ServeHTTP(writer, request)

	resp = writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	user = User{}
	if err := json.Unmarshal(body, &user); err != nil {
		t.Fatal(err)
	}

	if user.Name != name {
		t.Fatal("Wrong content, was expecting", name, "but got", user.Name)
	}
}
