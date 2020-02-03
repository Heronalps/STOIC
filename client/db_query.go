package client

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

/*
QueryDataSet queries a certain amount of data point of runtimes from ProcessingTime table
numDP - The number of data point in the data set
return X (nSamples, nFeatures), Y (nSamples)
*/
func QueryDataSet(runtime string, app string, version string, numDP int) (mat.Matrix, mat.Matrix) {
	var (
		XSlice []float64
		YSlice []float64
		X      mat.Matrix
		Y      mat.Matrix
	)

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	queryStr := fmt.Sprintf(`SELECT image_num, %s from ProcessingTime%s 
	WHERE application = ? and version = ? ORDER BY task_id DESC LIMIT ?;`, runtime, strings.Title(runtime))

	rows, err := db.Query(queryStr, app, version, numDP)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var (
			imageNum float64
			procTime float64
		)
		if err := rows.Scan(&imageNum, &procTime); err != nil {
			fmt.Println(err.Error())
		}
		XSlice = append(XSlice, imageNum)
		YSlice = append(YSlice, procTime)
	}
	nSamples, nFeatures, nOutputs := len(XSlice), 1, 1

	// Database has data points for current runtime
	if nSamples > 0 {
		X = mat.NewDense(nSamples, nFeatures, XSlice)
		Y = mat.NewDense(nSamples, nOutputs, YSlice)
	}

	return X, Y
}

/*
QueryDeploymentTime queries latest deployment time of specific runtime
*/
func QueryDeploymentTime(runtime string) float64 {
	var (
		deploymentTimes []float64
	)

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	// LIMIT 1 => Latest deployment time
	// LIMIT numDP => median on 10 latest deployment time
	queryStr := fmt.Sprintf("SELECT %s from DeploymentTime ORDER BY deployment_id DESC LIMIT ?", runtime)
	rows, err := db.Query(queryStr, windowSizes[runtime])
	fmt.Printf("WindowSizes: %v..\n", windowSizes)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var deploymentTime float64
		if err := rows.Scan(&deploymentTime); err != nil {
			fmt.Println(err.Error())
		}
		deploymentTimes = append(deploymentTimes, deploymentTime)
	}
	if len(deploymentTimes) == 0 {
		return defaultDeployTimes[runtime]
	}
	//fmt.Printf("deployment time : %f \n", deploymentTime)
	// return Median(deploymentTimes)
	return Average(deploymentTimes)
}

/*
QueryDeploymentTimeSeries queries a series of deployment time of specific runtime
*/
func QueryDeploymentTimeSeries(runtime string) []float64 {
	var (
		deploymentTimes []float64
	)

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	queryStr := fmt.Sprintf("SELECT %s from DeploymentTime ORDER BY deployment_id DESC LIMIT ?", runtime)
	rows, err := db.Query(queryStr, deploymentTimeNumDP)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var deploymentTime float64
		if err := rows.Scan(&deploymentTime); err != nil {
			fmt.Println(err.Error())
		}
		deploymentTimes = append(deploymentTimes, deploymentTime)
	}
	return deploymentTimes
}

/*
QueryAppVersion queries the latest version of an application
Returns 0 when app doesn't exist in AppVersion table
*/
func QueryAppVersion(app string) string {
	var (
		version string = "0"
	)

	db := connectDB(username, password, dbIP, dbPort)
	useDB(db, dbName)
	defer db.Close()
	queryStr := fmt.Sprintf("SELECT version from AppVersion WHERE app=?;")
	rows, err := db.Query(queryStr, app)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			fmt.Println(err.Error())
		}
	}
	return version
}
