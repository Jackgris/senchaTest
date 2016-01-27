package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Statment for create tables on database
var groups = "CREATE TABLE IF NOT EXISTS `Groups` (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR(45) NOT NULL);"
var user = "CREATE TABLE IF NOT EXISTS `User` (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR(100) NOT NULL, `userName` VARCHAR(20) NOT NULL, `password` VARCHAR (100) NOT NULL, `email` VARCHAR(100) NOT NULL, `picture` VARCHAR(100) NULL, `Group_id` INT NOT NULL);"
var db *sql.DB

func main() {
	sm := http.NewServeMux()
	sm.Handle("/", http.FileServer(http.Dir("./masteringextjs/")))
	s := http.NewServeMux()
	s.Handle("/", sm)
	s.HandleFunc("/security/signup", func(rw http.ResponseWriter, r *http.Request) {
		user, err := DecodeUserData(r.Body)
		if err != nil {
			log.Println("Error signup" + err.Error())
			return
		}
		name := stripSlashes(user.name)
		pass := stripSlashes(user.pass)

		sql := "SELECT * FROM USER WHERE userName=" + name
		result, err := db.Query(sql, nil)
		// FIXME unimplmented
		log.Println(result, pass)
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	os.Remove("./sakila.db")
	// Create and check database connection
	db, err := sql.Open("sqlite3", "./sakila.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection created")
	defer db.Close()

	_, err = db.Exec(groups, nil)
	if err != nil {
		log.Printf("%q: %s\n", err, groups)
		return
	}

	_, err = db.Exec(user, nil)
	if err != nil {
		log.Printf("%q: %s\n", err, user)
		return
	}

	i1 := "INSERT INTO Groups (`name`) VALUES ('admin');"
	i2 := "INSERT INTO User (`name`, `userName`, `password`, `email`, `Group_id`) VALUES ('Loiane Groner', 'loiane', '$2a$10$2a4e8803c91cc5edca222evoNPfhdRyGEG9RZcg7.qGqTjuCgXKda', 'me@loiane.com', '1');"

	_, err = db.Exec(i1, nil)
	if err != nil {
		log.Printf("%q: %s\n", err, i1)
		return
	}
	_, err = db.Exec(i2, nil)
	if err != nil {
		log.Printf("%q: %s\n", err, i2)
		return
	}

	fmt.Println("Starting server")
	err = server.ListenAndServe()
	log.Println("Trouble starting the server:" + err.Error())
}

type User struct {
	name string
	pass string
}

func DecodeUserData(r io.ReadCloser) (*User, error) {
	defer r.Close()
	var u User
	err := json.NewDecoder(r).Decode(&u)
	return &u, err
}

func stripSlashes(i string) (o string) {
	o = strings.Replace(i, "\\", "", -1)
	return o
}
