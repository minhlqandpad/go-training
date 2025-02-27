package main

import (
	"testing"
	"math"
)

func TestShapeCalculations(t *testing.T) {
	areaTests := []struct {
		name string
		shape Shape
		area float64
		perimeter float64
	}{
		{
			name: "Rectangle",
			shape: Rectangle{1.0, 2.0}, 
			area: 2.0, 
			perimeter: 6.0,
		},
		{
			name: "Circle",
			shape: Circle{1.5}, 
			area: math.Pi * 1.5 * 1.5,
			perimeter: 2 * math.Pi * 1.5,
		},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			gotArea := tt.shape.Area()
			wantArea := tt.area
			assertCorrectResult(t, gotArea, wantArea)

			gotPerimeter := tt.shape.Perimeter()
			wantPerimeter := tt.perimeter
			assertCorrectResult(t, gotPerimeter, wantPerimeter)
		})
	}
}

func assertCorrectResult(t testing.TB, got, want float64) {
	t.Helper()
	if got != want {
		t.Errorf("expected %g but got %g instead", want, got)
	}
}
