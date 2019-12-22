package database

import (
	"database/sql"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

/*
CreateProcessingTimeTable creates a table for recording processing time of image batches
*/
func CreateProcessingTimeTable() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

/*
CreateDeploymentTimeTable creates a table for monitoring deployment time of runtimes
*/
func CreateDeploymentTimeTable() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}
