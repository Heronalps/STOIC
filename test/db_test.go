package test

import (
	"fmt"
	"testing"

	"github.com/heronalps/STOIC/database"
)

func TestCreateDatabase(t *testing.T) {
	err = database.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("TestCreateDatabase.. \n")
		t.Errorf("%s was not created...\n", dbName)
	}
}

func TestCreateProcessingTimeTable(t *testing.T) {
	err = database.CreateProcessingTimeTable(dbName, "edge")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = database.CreateProcessingTimeTable(dbName, "cpu")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = database.CreateProcessingTimeTable(dbName, "gpu1")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
	err = database.CreateProcessingTimeTable(dbName, "gpu2")
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
}

func TestCreateDeploymentTimeTable(t *testing.T) {
	err = database.CreateDeploymentTimeTable(dbName)
	if err != nil {
		t.Errorf("TestCreateDeploymentTimeTable...\n")
		t.Errorf("DeploymentTime table was not created...\n")
	}
}

func TestAppendRecordProcessing(t *testing.T) {
	err = database.AppendRecordProcessing(dbName, "edge", 10, 1.56)
	if err != nil {
		t.Errorf("TestAppendRecord...\n")
		t.Errorf("The record was not appended...\n")
	}
}

func TestAppendRecordDeployment(t *testing.T) {
	err = database.AppendRecordDeployment(dbName, 1.1, 1.2, 1.3)
	if err != nil {
		t.Errorf("TestAppendRecord...\n")
		t.Errorf("The record was not appended...\n")
	}
}

func TestQueryDeploymentTime(t *testing.T) {
	var deploymentTime interface{}
	deploymentTime = database.QueryDeploymentTime(1)
	defer database.QueryDeploymentTime(0)

	duration, ok := deploymentTime.(float64)
	if !ok {
		t.Errorf("The query of deployment time was not successful...\n")
	}
	fmt.Printf("Duration : %f seconds..\n", duration)
}

func TestUpdateDeploymentTimeTable(t *testing.T) {
	err = database.UpdateDeploymentTimeTable()
	if err != nil {
		t.Errorf("Updating DeploymentTime table was not successful ...\n")
	}
}
