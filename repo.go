package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func RebuildDB() error {
	fmt.Println("Database created.")
	err := os.Remove(Config.Dbname)
	if err != nil {
		fmt.Println(err)
	}
	//sqlite doesn't require autoincrement on id's (integer primary key auto assigns an unused random rowid if empty)
	sqlStmt := "create table user (id integer not null primary key, email text, hashword text, rank integer, since text);" +
		"create table character (id integer not null primary key, user_id text, name text, img text, json_stats text);"

	_, err = Db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	fmt.Println("Database created.")
	return nil
}
func GetUserById(id int) (*User, error) {
	query := "SELECT * FROM USERS WHERE Id = ?"
	row := Db.QueryRow(query, id)
	return MapUser(row)
}
func InsertUser(u *User) error {
	_, err := Db.Exec(fmt.Sprintf("INSERT into User (email, hashword, rank, since) VALUES (%v, %v, %v, %v)"), u.Email, u.Hashword, u.Rank, u.Since)
	if err != nil {
		return err
	}
	return nil
}
func GetUserByEmail(email string) (*User, error) {
	query := "SELECT * FROM USERS WHERE EMAIL = ?"
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
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s) //returns both a time and an error so it can be returned directly.
}
