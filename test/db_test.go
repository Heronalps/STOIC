package test

import (
	"testing"

	"github.com/heronalps/STOIC/database"
)

var dbName string = "test"
var err error

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

func TestAppendRecord(t *testing.T) {
	err = database.AppendRecord(dbName, "edge", 10, 1.56)
	if err != nil {
		t.Errorf("TestAppendRecord...\n")
		t.Errorf("The record was not appended...\n")
	}
}
