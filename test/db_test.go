package test

import (
	"testing"

	"github.com/heronalps/STOIC/database"
)

func TestCreateDatabase(t *testing.T) {
	name := "test"
	result := database.CreateDatabase(name)
	if result {
		t.Errorf("%s was not created...", name)
	}
}
