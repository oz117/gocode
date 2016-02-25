/*

 */
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

// Prepares a query and returns, if successful *sql.Stmt
func prepareQuery(query string) (stmt *sql.Stmt) {
	stmt, err := db.Prepare(query)
	handleErrors(err)
	return
}

// Used to access the database returns, if successful *sql.DB
func accessDb(driver string, user string, password string, host string, port string, dbName string) (db *sql.DB) {
	logs := user + ":"

	if len(password) != 0 {
		logs = logs + password
	}
	db, err := sql.Open(driver, logs+"@tcp("+host+":"+port+")/"+dbName)
	handleErrors(err)
	log.Print("Database connection ok!")
	return
}

// Does a log.Fatal if anything bad occurs. It s not the best way of doing things.
// I only do this that way because this is a test.
func handleErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var str string
	var id int
	var all = make(map[int]string)

	db = accessDb("mysql", "root", "", "127.0.0.1", "3306", "test")
	defer db.Close()
	stmt := prepareQuery("SELECT id, name FROM user")
	defer stmt.Close()
	rows, err := stmt.Query()
	handleErrors(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &str)
		handleErrors(err)
		all[id] = str
	}
	err = rows.Err()
	handleErrors(err)
	for id, name := range all {
		log.Printf("[%d][%s]", id, name)
	}
}
