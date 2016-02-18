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

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"
)

// Statment for create tables on database
var groups = "CREATE TABLE IF NOT EXISTS `Groups` (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR(45) NOT NULL);"
var user = "CREATE TABLE IF NOT EXISTS `User` (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR(100) NOT NULL, `userName` VARCHAR(20) NOT NULL, `password` VARCHAR (100) NOT NULL, `email` VARCHAR(100) NOT NULL, `picture` VARCHAR(100) NULL, `Group_id` INT NOT NULL);"
var db *sql.DB

func main() {

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
	// important values userName = loiane and password = Packt123@
	//i2 := "INSERT INTO User (`name`, `userName`, `password`, `email`, `Group_id`) VALUES ('Loiane Groner', 'loiane', '$2a$10$2a4e8803c91cc5edca222evoNPfhdRyGEG9RZcg7.qGqTjuCgXKda', 'me@loiane.com', '1');"
	newPass := "Packt123@"
	savePass, err := GenerateCryptPass([]byte(newPass))
	if err != nil {
		panic("Error save pass encrypted " + err.Error())
	}
	i2 := "INSERT INTO User (`name`, `userName`, `password`, `email`, `Group_id`) VALUES ('Loiane Groner', 'loiane', '" + string(savePass) + "', 'me@loiane.com', '1');"
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

	sm := http.NewServeMux()
	sm.Handle("/", http.FileServer(http.Dir("./masteringextjs/")))
	s := http.NewServeMux()
	s.Handle("/", sm)
	s.HandleFunc("/security/signup", func(w http.ResponseWriter, r *http.Request) {
		rd := render.New()
		user := User{}
		r.ParseForm()
		user.Name = r.FormValue("user")
		user.Password = r.FormValue("password")
		log.Println("Signup request", user)
		name := stripSlashes(user.Name)
		pass := stripSlashes(user.Password)
		if name == "" || pass == "" {
			log.Println("Not parameters on the request")
			rd.JSON(w, http.StatusBadRequest, map[string]string{"Error": "Bad parameters"})
			return
		}

		sql := "SELECT * FROM USER WHERE userName='" + name + "'"
		result, err := db.Query(sql, nil)
		if err != nil {
			log.Println("Doesn't find user", err.Error())
			rd.JSON(w, http.StatusNotFound, map[string]string{"Error": err.Error()})
			return
		}
		defer result.Close()

		userDb := User{}
		for result.Next() {
			log.Println("On for loop")
			result.Scan(&userDb.Id, &userDb.Name, &userDb.UserName,
				&userDb.Password, &userDb.Email, &userDb.Picture, &userDb.Group)
			log.Println(userDb)
		}

		// FIXME we need encrypt password when the user submmit the login
		compareName := (name == userDb.UserName)
		comparePass := CheckPass(pass, userDb.Password)
		log.Println("Signup db", userDb, comparePass, compareName)

		if compareName && comparePass {
			userDb.Success = true
			userDb.Msg = "User authenticated!"
			userDb.Authenticated = "yes"

			send := Send{}
			// send.Authenticated = "yes"
			// send.UserName = userDb.UserName
			send.Success = true
			send.Msg = "User authenticated!"

			rd.JSON(w, http.StatusOK, send)
			return
		}

		send := Send{}
		send.Msg = "Not Authorized"
		rd.JSON(w, http.StatusNotFound, send)
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	fmt.Println("Starting server")
	err = server.ListenAndServe()
	log.Println("Trouble starting the server:" + err.Error())
}

func CheckPass(pass string, db string) bool {
	valid, err := CheckPassword([]byte(db), GenerateRandomSalt(), []byte(pass))
	if err != nil {
		log.Println("Error check pass", err)
		return false
	}
	return valid
}

type Send struct {
	// Implemented data struct from example php
	// Authenticated string `json:"authenticated"`
	// UserName      string `json:"username"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type User struct {
	// Implemented data struct from example php
	Id            int
	Name          string
	UserName      string `json:"username"`
	Password      string
	Email         string
	Picture       string
	Group         int
	Authenticated string `json:"authenticated"`
	Success       bool   `json:"success"`
	Msg           string `json:"msg"`
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

/////////////////////////////////////////////
// First aproach to encrypt user password
/////////////////////////////////////////////
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lenSalt = 10

// Checking the input password be the user's password
func CheckPassword(hash []byte, salt []byte, pass []byte) (valid bool, err error) {
	log.Println(hash, salt, pass)
	log.Println(string(hash), string(salt), string(pass))
	// use another algorithm that append method
	passSalt := append(pass[:], salt[:]...)
	err = bcrypt.CompareHashAndPassword(hash, passSalt)
	if err != nil {
		log.Println("CheckPass: the password isn't the same", err)
		return false, err
	}
	return true, nil
}

// Create a salt needed for create the hash used for register and login user
func GenerateRandomSalt() []byte {
	// rand.Seed(time.Now().UTC().UnixNano())
	// b := make([]byte, lenSalt)
	// for i := range b {
	// 	b[i] = letterBytes[rand.Intn(len(letterBytes))]
	// }
	// return b
	return []byte("hola")
}

// Create password encrypte to save on the database
func GenerateCryptPass(pass []byte) ([]byte, error) {
	salt := GenerateRandomSalt()
	passSalt := append(pass[:], salt[:]...)
	hashHex, err := bcrypt.GenerateFromPassword(passSalt, bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt broke")
		return pass, err
	}
	return hashHex, nil
}
