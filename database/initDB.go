package database

import (
	"database/sql"
	"fmt"
	"strings"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var username string = "root"
var password string = "123456"
var ip string = "127.0.0.1"
var port int = 3306

func connectDB(username string, password string, ip string, port int) *sql.DB {
	// Define Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, ip, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func useDB(db *sql.DB, dbName string) {
	_, err := db.Exec("USE " + dbName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Select the %s ...\n", dbName)
}

/*
CreateDatabase creates a database in MySQL instance. CREATE operation is idempotent.
*/
func CreateDatabase(dbName string) error {
	db := connectDB(username, password, ip, port)
	defer db.Close()
	_, err := db.Exec(fmt.Sprintf("CREATE database %s", dbName))
	if err != nil {
		fmt.Println("In CreateDatabase function ... ")
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Successfully created database %s ...\n", dbName)
	return err
}

/*
CreateProcessingTimeTable creates a table for recording processing time of image batches
*/
func CreateProcessingTimeTable(dbName string, runtime string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`CREATE TABLE ProcessingTime%s (
		task_id INT NOT NULL AUTO_INCREMENT, 
		time_stamp TIMESTAMP NOT NULL, 
		image_num INT NOT NULL, 
		%s FLOAT, 
		primary key(task_id));`, strings.Title(runtime), runtime)

	stmt, err := db.Prepare(stmtStr)
	if err != nil {
		fmt.Println("Error in Preparing stmt...")
		fmt.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Error in Executing stmt...")
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("%s ProcessingTime table is created successfully...\n", runtime)
	return err
}

/*
CreateDeploymentTimeTable creates a table for monitoring deployment time of runtimes
*/
func CreateDeploymentTimeTable(dbName string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE DeploymentTime(
		deployment_id INT NOT NULL AUTO_INCREMENT, 
		time_stamp TIMESTAMP NOT NULL, 
		edge FLOAT, 
		cpu FLOAT, 
		gpu1 FLOAT, 
		gpu2 FLOAT, 
		primary key(deployment_id));`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("DeploymentTime table is created successfully...")
	return err
}
