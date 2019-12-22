package database

import (
	"database/sql"
	"fmt"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var username string = "root"
var password string = "123456"
var ip string = "127.0.0.1"
var port int = 3306

func connectDB(username string, password string, ip string, port int) *sql.DB {
	// Define Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/\n", username, password, ip, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/*
CreateDatabase creates a database in MySQL instance
*/
func CreateDatabase(name string) bool {
	db := connectDB(username, password, ip, port)
	defer db.Close()
	_, err := db.Exec("CREATE database " + name)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Printf("Successfully created database %s ...", name)
	return true
}

/*
CreateProcessingTimeTable creates a table for recording processing time of image batches
*/
func CreateProcessingTimeTable() {
	db := connectDB(username, password, ip, port)
	defer db.Close()
}

/*
CreateDeploymentTimeTable creates a table for monitoring deployment time of runtimes
*/
func CreateDeploymentTimeTable() {
	db := connectDB(username, password, ip, port)
	defer db.Close()
}
