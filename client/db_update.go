package client

import (
	"errors"
	"fmt"
	"strings"
)

/*
AppendRecordProcessing appends a record of (image num, duration) to Processing Time table of specific runtime
*/
func AppendRecordProcessing(dbName string, runtime string, imageNum int,
	duration float64, application string, version string) error {

	fmt.Printf("Updating ProcessingTime table of %s duration %f...\n", runtime, duration)

	db := connectDB(username, password, dbIP, dbPort)
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
func AppendRecordDeployment(dbName string, cpu string, gpu1 string, gpu2 string) error {
	if cpu == "null" || gpu1 == "null" || gpu2 == "null" {
		return errors.New("unable to deploy runtime")
	}

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`INSERT INTO DeploymentTime (cpu, gpu1, gpu2) VALUES (?, ?, ?);`)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(NullString(cpu), NullString(gpu1), NullString(gpu2))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

/*
QueryDeploymentTimeNautilus queries the current deployment time on Nautilus for runtimes
*/
func QueryDeploymentTimeNautilus(numGPU int64, app string) (bool, float64) {
	var (
		duration   float64
		isDeployed bool
		err        error
	)

	_, isDeployed, duration, err = Deploy(namespace, timeQueryDeployment, numGPU, app)

	if err != nil {
		fmt.Println(err.Error())
	}
	return isDeployed, duration
}

/*
UpdateDeploymentTimeTable updates the DeploymentTime table
*/
func UpdateDeploymentTimeTable(app string) error {
	var (
		numGPU int64
		// 0 - cpu, 1 - gpu1, 2 - gpu2
		deploymentTimes [3]string
		err             error
	)

	// Start from the current GPU number + 1
	numGPU = QueryGPUNum(namespace, timeQueryDeployment)
	for i := 0; i < 3; i++ {
		numGPU = (numGPU + 1) % 3
		isDeployed, elasped := QueryDeploymentTimeNautilus(numGPU, app)
		if isDeployed {
			deploymentTimes[numGPU] = fmt.Sprintf("%f", elasped)
		} else {
			deploymentTimes[numGPU] = "null"
		}
		switch numGPU {
		case 0:
			runtimes["cpu"] = isDeployed
		case 1:
			runtimes["gpu1"] = isDeployed
		case 2:
			runtimes["gpu2"] = isDeployed
		}
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

/*
UpdateAppVersion updates the latest version of an application
*/
func UpdateAppVersion(app string, version string) error {
	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`UPDATE AppVersion SET Version=? WHERE app=?;`)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(version, app)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

/*
InsertAppVersion inserts the nonexistent app and its current version
*/
func InsertAppVersion(app string, version string) error {
	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`INSERT INTO AppVersion (App, Version) VALUES (?, ?);`)
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(app, version)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

/*
AppendRecordLogTime appends a record to LogTime table
*/
func AppendRecordLogTime(imageNum int, app string, version string, runtime string,
	predTotal float64, predTransfer float64, predDeploy float64, predProc float64,
	actTotal float64, actTransfer float64, actDeploy float64, actProc float64) error {

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	stmt, err := db.Prepare(`INSERT INTO LogTime (image_num, app, version, runtime, 
		pred_total, pred_transfer, pred_deploy, pred_proc,
		act_total, act_transfer, act_deploy, act_proc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	_, err = stmt.Exec(imageNum, app, version, runtime,
		predTotal, predTransfer, predDeploy, predProc,
		actTotal, actTransfer, actDeploy, actProc)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}
