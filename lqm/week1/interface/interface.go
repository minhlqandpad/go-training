package main

import (
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	r float64
}

type Rectangle struct {
	h, w float64
}

func (c Circle) Area() float64 {
	return c.r * c.r * math.Pi
}

func (c Circle) Perimeter() float64 {
	return 2 * c.r * math.Pi
}

func (r Rectangle) Area() float64 {
	return r.h * r.w
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.h + r.w)
}
