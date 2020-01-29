package client

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

/*
SetupRegression set up multiple data points in each runtime processing time table for regression
when the app / version is updated
*/
func SetupRegression(app string, version string) {
	dbVersion := QueryAppVersion(app)
	fmt.Println("DB version : " + dbVersion)
	fmt.Println("current version : " + version)
	if dbVersion == "0" {
		fmt.Printf("Insert app %s version %s to AppVersion table...\n", app, version)
		InsertAppVersion(app, version)
	}
	// Set up the table when app / version is updated
	result := CompareVersion(version, dbVersion)
	if result == 1 {
		fmt.Printf("Current version %s is greater than DB version %s ..\n", version, dbVersion)
		UpdateAppVersion(app, version)
		for _, runtime := range runtimes {
			for _, imageNum := range setupImageNums {
				Schedule(runtime, imageNum, app, version, false)
			}
		}
	} else {
		fmt.Printf("Current version %s equals to / is less than DB version %s .. \n", version, dbVersion)
		fmt.Println("Checking if at least two data points exist for each runtime...")
		for _, runtime := range runtimes {
			var (
				X    mat.Matrix
				rows int = 0
			)
			// At least 2 data points exist for each app & version
			X, _ = QueryDataSet(runtime, app, version, 2)
			if X != nil {
				rows, _ = X.Dims()
			}
			if rows < 2 {
				for _, imageNum := range setupImageNums {
					Schedule(runtime, imageNum, app, version, false)
				}
			}
		}
	}
}

// /*
// Regress function leveraging Bayesian Ridge Regression
// */
// func Regress(runtime string, app string, version string, numDP int) (float64, float64) {
// 	var (
// 		coef      float64
// 		intercept float64
// 	)
// 	X, Y := QueryDataSet(runtime, app, version, numDP)

// 	if X == nil && Y == nil {
// 		fmt.Printf("No data point of %s in DB...\n", runtime)
// 		return 0.0, 0.0
// 	}
// 	nSamples, nOutputs := X.Dims()
// 	if nSamples <= nOutputs {
// 		fmt.Printf("Single data point of %s in DB...\n", runtime)
// 		return Y.At(0, 0) / X.At(0, 0), 0
// 	}
// 	YPred := mat.NewDense(nSamples, nOutputs, nil)
// 	model := linearmodel.NewBayesianRidge()
// 	model.Fit(X, Y)
// 	model.Predict(X, YPred)
// 	coef = model.LinearModel.Coef.At(0, 0)
// 	intercept = model.LinearModel.Intercept.At(0, 0)
// 	if math.IsNaN(coef) || math.IsNaN(intercept) {
// 		return 0.0, 0.0
// 	}
// 	// r2score := metrics.R2Score(Y, YPred, nil, "variance_weighted").At(0, 0)
// 	// if r2score > .999 {
// 	// 	fmt.Println("BayesianRidge ok")
// 	// }
// 	return coef, intercept
// }

/*
Regress function leveraging Bayesian Ridge Regression
*/
func Regress(runtime string, app string, version string, numDP int) (float64, float64) {
	var (
		coef      float64
		intercept float64
		match     [][]byte
		output    []byte
		err       error
		cmd       *exec.Cmd
	)
	repoPATH := HomeDir() + "/go/src/github.com/heronalps/STOIC"
	X, Y := QueryDataSet(runtime, app, version, numDP)

	if X == nil && Y == nil {
		fmt.Printf("No data point of %s in DB...\n", runtime)
		return 0.0, 0.0
	}
	nSamples, nOutputs := X.Dims()
	if nSamples <= nOutputs {
		fmt.Printf("Single data point of %s in DB...\n", runtime)
		return Y.At(0, 0) / X.At(0, 0), 0
	}

	// Run RANSAC regression
	FILE := "./scripts/robust_regression.py "
	cmdRun := fmt.Sprintf("python3 %s %s %s %s %d", FILE, runtime, app, version, numDP)
	// fmt.Println(cmdRun)
	cmd = exec.Command("bash", "-c", cmdRun)
	cmd.Dir = repoPATH
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error running regression. msg: %s \n", err.Error())
	}
	//fmt.Println(string(output))
	reCoef := regexp.MustCompile(`Coefficient: (\d*\.\d*)`)
	match = reCoef.FindSubmatch(output)
	if len(match) > 0 {
		coef, _ = strconv.ParseFloat(string(match[1]), 64)
	}

	reInt := regexp.MustCompile(`Intercept: (\d*\.\d*)`)
	match = reInt.FindSubmatch(output)
	if len(match) > 0 {
		intercept, _ = strconv.ParseFloat(string(match[1]), 64)
	}

	// if math.IsNaN(coef) || math.IsNaN(intercept) {
	// 	return 0.0, 0.0
	// }

	return coef, intercept
}
