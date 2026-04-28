package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Triangle struct {
	Base, Height float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Angles() (float64, float64, float64) {
	return 60, 60, 60 // equilateral for demo
}

func describeShape(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())

	// Type assertion to access Triangle-specific method
	if t, ok := s.(Triangle); ok {
		a1, a2, a3 := t.Angles()
		fmt.Printf("Triangle angles: %.0f, %.0f, %.0f\n", a1, a2, a3)
	}
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 3, Height: 4},
		Triangle{Base: 6, Height: 8},
	}

	for _, s := range shapes {
		fmt.Printf("Shape: %T\n", s)
		describeShape(s)
		fmt.Println()
	}
}
