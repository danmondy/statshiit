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
	Content interface{}
}
type Response struct {
	Success bool
	Data string
}

//Making db global so I don't have to pass it around.
//I might change my mind but for now it seems simple. And simple is good.
var (
	Config Configuration
	Db     *sql.DB
	commands map[string]func(User, interface{}, chan Response)
)

func main() {
	commands = map[string]func(User, interface{}, chan Response){			
			"addCharacter": handleAddChar,
	}
	readConfig()
	var err error
	Db, err = sql.Open("sqlite3", Config.Dbname)
	if err != nil {
		panic(err)
	}
	
	RebuildDB();PopulateDB(); //uncomment if you make a schema change - will wipe data
	

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
		
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	var ops []Operation
	err = json.Unmarshal(body, &ops)	
	if err != nil{
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	results := make(chan Response, len(ops))
	for _, op := range ops{
		if handlr, ok := commands[op.Command]; ok{
			go handlr(*user, op.Content, results)
		}
	}
	responses := make([]Response, len(ops))
	for i, _ := range responses{
		responses[i] = <- results
	}
	fmt.Fprintf(w, "%v", responses)
}
func handleAddChar(u User, content interface{}, results chan Response){
	data := content.(map[string]interface{})
	_, okey := data["name"]; _, do := data["stats"]; _, key := data["image"]
	if !(okey && do && key) {
		results <- Response{false, "Content was missing either name, stats, or image"}
		return
	}

	stats, err := json.Marshal(data["stats"])
	c := NewCharacter(u.Id, data["name"].(string), string(stats), data["image"].(string))
	_, err = InsertChar(c)
	if err != nil{
		results <- Response{false, err.Error()}
		return
	}
	res, err := json.Marshal(c)
	results <- Response{true, string(res)}
}
