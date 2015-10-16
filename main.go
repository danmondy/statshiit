package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

//Making db global so I don't have to pass it around.
//I might change my mind but for now it seems simple. And simple is good.
var (
	Config Configuration
	Db     *sql.DB
)

func main() {
	/*commands := map[string]func(chan int){
			"insert": insert,
			"update": update,
	}*/
	readConfig()
	var err error
	Db, err = sql.Open("sqlite3", Config.Dbname)
	if err != nil {
		panic(err)
	}
	
	RebuildDB() //uncomment if you make a schema change - will wipe data
	populateDB()

	fmt.Println(Config)
	http.HandleFunc("/api", apiHandler)
	fmt.Println("Listening")
	http.ListenAndServe(":"+Config.Port, nil)
	fmt.Println("Good bye:)")
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
	//AUTH
	email, p, ok := r.BasicAuth()
	if !ok{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("There was an error reading your credentials"))
		return
	}
	user, err := GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User not found:("+err.Error()))
		return
	}
	err = user.CompareHash([]byte(p))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("The password does not match."))
		return
	}
	//END AUTH - so continue with user object
	w.Write([]byte("Authed!"))
	return

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	var ops []Operation
	err = json.Unmarshal(body, ops)
	fmt.Fprintf(w, "%v\n %v", user, ops)
}
func populateDB(){
	users := []User{NewUser("danmondy@gmail.com", "happy", 3), NewUser("josh@gmail.com", "happy", 3)}
	for _, u := range users{
		err := InsertUser(&u);		
		if err != nil{
			fmt.Println(err)
		}
	}
}
