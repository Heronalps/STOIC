package test

import (
	"fmt"
	"testing"

	"github.com/heronalps/STOIC/client"
)

func TestCreateDatabase(t *testing.T) {
	err = client.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("TestCreateDatabase.. \n")
		t.Errorf("%s was not created...\n", dbName)
	}
}

func TestCreateProcessingTimeTable(t *testing.T) {
	err = client.CreateProcessingTimeTable(dbName, "edge")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = client.CreateProcessingTimeTable(dbName, "cpu")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = client.CreateProcessingTimeTable(dbName, "gpu1")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = client.CreateProcessingTimeTable(dbName, "gpu2")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
}

func TestCreateDeploymentTimeTable(t *testing.T) {
	err = client.CreateDeploymentTimeTable(dbName)
	if err != nil {
		t.Errorf("TestCreateDeploymentTimeTable...\n")
		t.Errorf("DeploymentTime table was not created...\n")
	}
}

func TestAppendRecordProcessing(t *testing.T) {
	err = client.AppendRecordProcessing(dbName, "edge", 10, 1.56, "wtb", "1.0")
	if err != nil {
		t.Errorf("TestAppendRecord...\n")
		t.Errorf("The record was not appended...\n")
	}
}

func TestAppendRecordDeployment(t *testing.T) {
	// err = client.AppendRecordDeployment(dbName, "1.1", "1.2", "1.3")
	// if err != nil {
	// 	t.Errorf("TestAppendRecord...\n")
	// 	t.Errorf("The record was not appended...\n")
	// }
	err = client.AppendRecordDeployment(dbName, "null", "1.2", "1.3")
	if err != nil {
		t.Errorf("TestAppendRecord...\n")
		t.Errorf("The record was not appended...\n")
	}
}

func TestQueryDeploymentTimeNautilus(t *testing.T) {
	var deploymentTime interface{}
	_, deploymentTime = client.QueryDeploymentTimeNautilus(1, app)
	defer client.QueryDeploymentTimeNautilus(0, app)

	duration, ok := deploymentTime.(float64)
	if !ok {
		t.Errorf("The query of deployment time was not successful...\n")
	}
	fmt.Printf("Duration : %f seconds..\n", duration)
}

func TestUpdateDeploymentTimeTable(t *testing.T) {
	err = client.UpdateDeploymentTimeTable(app)
	if err != nil {
		t.Errorf("Updating DeploymentTime table was not successful ...\n")
	}
}

func TestCreateRegressionTable(t *testing.T) {
	err = client.CreateRegressionTable(dbName)
	if err != nil {
		t.Errorf("Creating Regression Table was not successful ...\n")
	}
}

func TestQueryDataSet(t *testing.T) {
	X, Y := client.QueryDataSet(runtime, app, version, numDP)
	fmt.Printf("X: %v \n", X)
	fmt.Printf("Y: %v \n", Y)
}

func TestCreateAppVersionTable(t *testing.T) {
	err := client.CreateAppVersionTable(dbName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestQueryAppVersion(t *testing.T) {
	version := client.QueryAppVersion(app)
	fmt.Println("version : " + version)
}

func TestInsertAppVersion(t *testing.T) {
	err := client.InsertAppVersion(app, version)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestUpdateAppVersion(t *testing.T) {
	err := client.UpdateAppVersion(app, updateVersion)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCreateLogTimeTable(t *testing.T) {
	err := client.CreateLogTimeTable(dbName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestAppendRecordLogTime(t *testing.T) {
	err := client.AppendRecordLogTime(imageNum, app, version, runtime, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0)
	if err != nil {
		fmt.Println(err.Error())
	}
}
