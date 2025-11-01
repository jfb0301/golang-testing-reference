package calculator_test

import (
	"testing"
	"github.com/jfb0301/golang-testing-reference/TDD/calculator"
	"os"
	"log"
)


func TestMain(m *testing.M) {
	// Setup statements
	setup()

	// run the tests
	e := m.Run()

	// Clean up statements
	teardown()

	// report the exit code
	os.Exit(e)
}

func setup() {
	log.Println("Setting up.")
}

func teardown() {
	log.Println("Tearing down.")
}

func TestAdd(t *testing.T) {
	// Arrange 
	e := calculator.Engine{}
	x, y := 2.5, 3.5 
	want := 6.0
	
	// Act 
	got := e.Add(x,y)

	// Assert 
	if got != want {
		t.Errorf("Add(%.2f, %.2f) incorrect, got: %.2f, want: %.2f", 2.5, 3.5, got, 6.0)
	}
}