package client

import (
	"database/sql"
	"fmt"
	"strings"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

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
	//fmt.Printf("Select the %s database ...\n", dbName)
}

/*
InitDB initializes DB and all necessary tables
*/
func InitDB() {
	CreateDatabase(dbName)
	for _, runtime := range runtimes {
		CreateProcessingTimeTable(dbName, runtime)
	}
	CreateDeploymentTimeTable(dbName)
	CreateAppVersionTable(dbName)
	CreateLogTimeTable(dbName)
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
		application VARCHAR(64),
		version VARCHAR(16),
		%s FLOAT, 
		primary key(task_id));`, strings.Title(runtime), runtime)

	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	if err != nil {
		fmt.Println("Error in Preparing stmt...")
		fmt.Println(err.Error())
		return err
	}

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
		cpu FLOAT, 
		gpu1 FLOAT, 
		gpu2 FLOAT, 
		primary key(deployment_id));`)
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("DeploymentTime table is created successfully...")
	return err
}

/*
CreateRegressionTable creates a table for Bayesian Ridge Regression coefficient and intercept
*/
func CreateRegressionTable(dbName string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE Regression(
		runtime VARCHAR(32) NOT NULL,
		coeff FLOAT NOT NULL, 
		intercept FLOAT NOT NULL);`)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	stmt, err = db.Prepare(`INSERT INTO Regression VALUE (?, ?, ?);`)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, runtime := range runtimes {
		_, err = stmt.Exec(runtime, 1.0, 1.0)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return err
}

/*
CreateAppVersionTable create table that maps application and latest version
*/
func CreateAppVersionTable(dbName string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE AppVersion(
		app VARCHAR(64) NOT NULL,
		version VARCHAR(16) NOT NULL);`)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("AppVersion table is created successfully...")
	return err
}

/*
CreateLogTimeTable create table for logging total response time and its component
*/
func CreateLogTimeTable(dbName string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE LogTime (
		task_id INT NOT NULL AUTO_INCREMENT,
		time_stamp TIMESTAMP NOT NULL,
		image_num INT NOT NULL,
		app varchar(64) NOT NULL,
		version varchar(16) NOT NULL,
		runtime VARCHAR(32) NOT NULL,
		pred_total FLOAT NOT NULL,
		pred_transfer FLOAT NOT NULL,
		pred_deploy FLOAT NOT NULL,
		pred_proc FLOAT NOT NULL,
		act_total FLOAT NOT NULL,
		act_transfer FLOAT NOT NULL,
		act_deploy FLOAT NOT NULL,
		act_proc FLOAT NOT NULL,
		primary key(task_id));`)
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
	fmt.Println("LogTime table is created successfully...")
	return err
}
