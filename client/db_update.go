package client

import (
	"fmt"
	"strings"
)

/*
AppendRecordProcessing appends a record of (image num, duration) to Processing Time table of specific runtime
*/
func AppendRecordProcessing(dbName string, runtime string, imageNum int,
	duration float64, application string, version string) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`INSERT INTO ProcessingTime%s (image_num, application, version, %s) VALUES (?, ?, ?, ?);`, strings.Title(runtime), runtime)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(imageNum, application, version, duration)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

/*
AppendRecordDeployment appends a record of (image num, duration) to Processing Time table of specific runtime
*/
func AppendRecordDeployment(dbName string, cpu float64, gpu1 float64, gpu2 float64) error {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`INSERT INTO DeploymentTime (cpu, gpu1, gpu2) VALUES (?, ?, ?);`)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(cpu, gpu1, gpu2)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

/*
QueryDeploymentTime queries the current deployment time on Nautilus for runtimes
*/
func QueryDeploymentTime(numGPU int64) float64 {
	var (
		duration float64
		err      error
	)

	_, duration, err = Deploy(namespace, deployment, numGPU)

	if err != nil {
		fmt.Println(err.Error())
	}
	return duration
}

/*
UpdateDeploymentTimeTable updates the DeploymentTime table
*/
func UpdateDeploymentTimeTable() error {
	var (
		numGPU int64
		// 0 - cpu, 1 - gpu1, 2 - gpu2
		deploymentTimes [3]float64
		err             error
	)

	// Start from the current GPU number + 1
	numGPU = QueryGPUNum(namespace, deployment)
	for i := 0; i < 3; i++ {
		numGPU = (numGPU + 1) % 3
		deploymentTimes[numGPU] = QueryDeploymentTime(numGPU)
	}
	err = AppendRecordDeployment(dbName, deploymentTimes[0], deploymentTimes[1], deploymentTimes[2])
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

/*
UpdateRegressionTimeTable updates coefficient and intercept of specific runtime in the Regression Table

Caution: Updating Regression Table is subject to disk I/O. Keep this functionality for future necessary use case
*/
func UpdateRegressionTimeTable(runtime string) error {
	return nil
}
