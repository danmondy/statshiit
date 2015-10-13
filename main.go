package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Configuration struct {
	Port       string
	Dbname     string
	Dbpassword string
}
type Operation struct {
	Command string
	Content string
}

var Config Configuration

func main() {
	commands := map[string]func(){
		"insert": insert,
		"update": update,
	}

	readConfig()

	http.HandleFunc("/api", apiHandler)
	http.ListenAndServe(":"+Config.port, nil)
}

func readConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusExceptionFailed)
		w.Write([]byte(err.Error()))
		return
	}
	var ops []Operations
	err := json.Unmarshal(body, ops)

}
