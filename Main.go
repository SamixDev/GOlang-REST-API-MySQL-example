package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//User Tag structure for your database
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Number  string `json:"number"`
}

var db *sql.DB
var err error

// function to open connection to mysql database
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""         // your password, leave it like this if there is no password
	dbName := "usertest" // your database name
	dbIP := "127.0.0.1:3306"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIP+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db

}

//Index func to view all the records
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()

	var users []User

	result, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post User
		err = result.Scan(&post.ID, &post.Name, &post.Country, &post.Number)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, post)
	}
	json.NewEncoder(w).Encode(users)
	defer db.Close()
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	vars := mux.Vars(r)
	Name := vars["name"]
	Country := vars["country"]
	Number := vars["number"]

	// perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO users(name, country, number) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Name, Country, Number)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
	defer db.Close()
}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Query("SELECT * FROM users WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	var post User
	for stmt.Next() {

		err = stmt.Scan(&post.ID, &post.Name, &post.Country, &post.Number)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
	defer db.Close()
}
func delUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
	defer db.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	Name := params["name"]
	Country := params["country"]
	Number := params["number"]

	// perform a db.Query insert
	stmt, err := db.Prepare("Update users SET name = ?, country = ?, number = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Name, Country, Number, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User with ID = %s was updated", params["id"])
	defer db.Close()
}
func main() {
	log.Println("Server started on: http://localhost:8080")
	router := mux.NewRouter()
	//On postman try http://localhost:8080/all with method GET
	router.HandleFunc("/all", Index).Methods("GET")
	//On postman try http://localhost:8080/add?name=Test&country=LEB&number=7777777 with metho POST
	router.HandleFunc("/add", insertUser).Methods("POST").Queries("name", "{name}", "country", "{country}", "number", "{number}")
	//On postman try http://localhost:8080/get/1 with method GET
	router.HandleFunc("/get/{id}", getUser).Methods("GET")
	//On postman try http://localhost:8080/update/1 with method PUT
	router.HandleFunc("/update/{id}", updateUser).Methods("PUT").Queries("name", "{name}", "country", "{country}", "number", "{number}")
	//On postman try http://localhost:8080/del/1 with method DELETE
	router.HandleFunc("/del/{id}", delUser).Methods("DELETE")

	http.ListenAndServe(":8080", router)

}
