package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
}

type Circle struct {
	x, y, r float64
}

func (c *Circle) area() float64 {
	return math.Pi * c.r * c.r
}

type Rectangle struct {
	x1, y1, x2, y2 float64
}

func (r *Rectangle) area() float64 {
	return math.Abs((r.x1 - r.x2) * (r.y1 - r.y2))
}

func totalArea(shapes ...Shape) (area float64) {
	area = 0
	for _, s := range shapes {
		area += s.area()
	}
	return
}

func main() {
	fmt.Println("I'm working!")

	c := &Circle{0, 0, 5}
	fmt.Printf("%f\n", c.area())

	r := &Rectangle{0, 0, 3, 4}
	fmt.Printf("%f\n", r.area())

	fmt.Printf("total area: %f\n", totalArea(r, c))
}
