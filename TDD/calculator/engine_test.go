package calculator_test

import (
	"testing"
	"github.com/jfb0301/golang-testing-reference/TDD/calculator"
)



func TestAdd(t *testing.T) {
	// Arrange 
	e := calculator.Engine{}
	x, y := 2.5, 3.5 
	want := 6.0
	
	// Act 
	got := e.Add(x,y)

	// Assert 
	if got != 6.0 {
		t.Errorf("Add(%.2f, %.2f) incorrect, got: %.2f, want: %.2f", 2.5, 3.5, got, 6.0)
	}
}