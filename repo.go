package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func RebuildDB() error {	
	err := os.Remove(Config.Dbname)
	if err != nil {
		fmt.Println(err)
	}
	//sqlite doesn't require autoincrement on id's (integer primary key auto assigns an unused random rowid if empty)
	sqlStmt := "create table user (id integer not null primary key, email text, hashword text, rank integer, since text);" +
		"create table character (id integer not null primary key, user_id integer, name text, image text, json_stats text, date_modified text);"

	_, err = Db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	fmt.Println("Database created.")
	return nil
}

//COMMON
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s) //returns both a time and an error so it can be returned directly.
}

//USERS
func GetUserById(id int) (*User, error) {
	query := "SELECT * FROM user WHERE Id = ?"
	row := Db.QueryRow(query, id)
	return MapUser(row)
}
func InsertUser(u *User) error {
	_, err := Db.Exec(fmt.Sprintf("INSERT into user (email, hashword, rank, since) VALUES ('%v', '%v', '%v', '%v')", u.Email, u.Hashword, u.Rank, TimeToString(u.Since)))
	if err != nil {
		return err
	}
	return nil
}
func GetUserByEmail(email string) (*User, error) {
	query := "SELECT * FROM user WHERE EMAIL = ?"
	row := Db.QueryRow(query, email)
	return MapUser(row)
}
func MapUser(r *sql.Row) (*User, error) {
	var u User
	var t string
	err := r.Scan(&u.Id, &u.Email, &u.Hashword, &u.Rank, &t)
	if err != nil {
		return nil, err
	}
	u.Since, err = StringToTime(t)
	return &u, err
}
//END USERS

//Characters
func GetCharById(id int) (*Character, error) {
	query := "SELECT * FROM character WHERE Id = ?"
	row := Db.QueryRow(query, id)
	return MapChar(row)
}
func InsertChar(c *Character) (sql.Result, error) {
	return Db.Exec(fmt.Sprintf("INSERT into character (user_id, name, stats_json, image, date_modified) VALUES ('%v', '%v', '%v', '%v')", c.UserId, c.Name, c.StatsJson, c.Image, c.DateModified))
}
func UpdateChar(c Character) (sql.Result, error) {
	return Db.Exec(fmt.Sprintf("UPDATE character SET name=%s, stats_json=%s, image=%s, date_modified=%v", c.Name, c.StatsJson, c.Image, c.DateModified))
}
func MapChar(r *sql.Row) (*Character, error) {
	var c Character
	var t string
	err := r.Scan(&c.Id, &c.UserId, &c.Name, &c.Image, &c.StatsJson, &t)
	if err != nil {
		return nil, err
	}
	c.DateModified, err = StringToTime(t)
	return &c, err
}
//End Characters
