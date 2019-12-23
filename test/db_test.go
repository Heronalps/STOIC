package test

import (
	"testing"

	"github.com/heronalps/STOIC/database"
)

func TestCreateDatabase(t *testing.T) {
	name := "test"
	err := database.CreateDatabase(name)
	if err != nil {
		t.Errorf("TestCreateDatabase.. \n")
		t.Errorf("%s was not created...\n", name)
	}
}

func TestCreateProcessingTimeTable(t *testing.T) {
	dbName := "test"
	err := database.CreateProcessingTimeTable(dbName)
	if err != nil {
		t.Errorf("TestCreateProcessingTimeTable...\n")
		t.Errorf("ProcessingTime table was not created...\n")
	}
}

func TestCreateDeploymentTimeTable(t *testing.T) {
	dbName := "test"
	err := database.CreateDeploymentTimeTable(dbName)
	if err != nil {
		t.Errorf("TestCreateDeploymentTimeTable...\n")
		t.Errorf("DeploymentTime table was not created...\n")
	}
}
