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
//	commands := map[string]func(){
//		"insert": insert,
//		"update": update,
//	}
	readConfig()
	fmt.Println(Config)
	http.HandleFunc("/api", apiHandler)

	fmt.Println("Listening")
	http.ListenAndServe(":"+Config.Port, nil)
	fmt.Println("Listening")
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
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	var ops []Operation
	err = json.Unmarshal(body, ops)
	fmt.Fprint(w, "Hello")
}
