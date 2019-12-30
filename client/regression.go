package client

import (
	"fmt"
	"math/rand"
	"strings"

	linearmodel "github.com/pa-m/sklearn/linear_model"
	"github.com/pa-m/sklearn/metrics"
	"gonum.org/v1/gonum/mat"
)

/*
QueryDataSet queries a certain amount of data point of runtimes from ProcessingTime table
*/
func QueryDataSet(runtime string, numDP int) {
	db := connectDB(username, password, ip, port)
	useDB(db, dbName)
	defer db.Close()
	stmtStr := fmt.Sprintf(`SELECT image_num, %s from ProcessingTime%s order by task_id desc limit ?;`, runtime, strings.Title(runtime))
	stmt, err := db.Prepare(stmtStr)
	defer stmt.Close()

	_, err = stmt.Exec(numDP)
	if err != nil {
		fmt.Println(err.Error())
	}

}

/*
Regress function
*/
func Regress(dataset float64) {

}

/*
Regress2 tests Bayesian Ridge Regression
*/
func Regress2() {
	nSamples, nFeatures, nOutputs := 10000, 5, 1
	X := mat.NewDense(nSamples, nFeatures, nil)
	X.Apply(func(i int, j int, v float64) float64 {
		return rand.NormFloat64() * 20
	}, X)
	f := func(X mat.Matrix, i, o int) float64 {
		if o == 0 {
			return 1. + 2.*X.At(i, 0) + 3.*X.At(i, 1) + 4.*X.At(i, 2)
		}
		return 1. - 2.*X.At(i, 0) + 3.*X.At(i, 1) + float64(o)*X.At(i, 2)

	}
	Y := mat.NewDense(nSamples, nOutputs, nil)
	Y.Apply(func(i int, o int, v float64) float64 {
		return f(X, i, o)
	}, Y)
	fmt.Printf("(0,0): %f \n", X.At(0, 0))
	fmt.Printf("(0,1): %f \n", X.At(0, 1))
	fmt.Printf("(0,2): %f \n", X.At(0, 2))
	fmt.Printf("(0,3): %f \n", X.At(0, 3))
	fmt.Printf("(0,4): %f \n", X.At(0, 4))

	fmt.Printf("(0,0): %f \n", Y.At(0, 0))
	// fmt.Printf("(0,1): %f \n", Y.At(0, 1))
	// fmt.Printf("(0,2): %f \n", Y.At(0, 2))
	// fmt.Printf("(0,3): %f \n", Y.At(0, 3))
	// fmt.Printf("(0,4): %f \n", Y.At(0, 4))
	m := linearmodel.NewBayesianRidge()
	//start := time.Now()
	m.Fit(X, Y)
	//elapsed := time.Since(start)
	Ypred := mat.NewDense(nSamples, nOutputs, nil)
	m.Predict(X, Ypred)
	fmt.Printf("YPred (0, 0) : %f \n", Ypred.At(0, 0))
	r2score := metrics.R2Score(Y, Ypred, nil, "variance_weighted").At(0, 0)
	if r2score > .999 {
		fmt.Println("BayesianRidge ok")
	}
}
