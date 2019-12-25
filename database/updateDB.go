package database

import (
	"fmt"
	"strings"
)

/*
AppendRecord appends a record of (image num, duration) to Processing Time table of specific runtime
*/
func AppendRecord(dbName string, runtime string, imageNum int, duration float64) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`INSERT INTO ProcessingTime%s (image_num, %s) VALUES (?, ?);`, strings.Title(runtime), runtime)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(imageNum, duration)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}
