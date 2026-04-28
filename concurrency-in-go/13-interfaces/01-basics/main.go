package main

import "fmt"

type Dense interface {
	Density() float64
}

type Metal struct {
	Mass   float64
	Volume float64
}

func (m Metal) Density() float64 {
	return m.Mass / m.Volume
}

type Gas struct {
	Pressure      float64
	Temperature   float64
	MolecularMass float64
}

func (g Gas) Density() float64 {
	return (g.Pressure * g.MolecularMass) / (8.314 * g.Temperature)
}

func printDensity(d Dense) {
	fmt.Printf("Density: %.4f\n", d.Density())
}

func main() {
	iron := Metal{Mass: 7874, Volume: 1}
	oxygen := Gas{Pressure: 101325, Temperature: 300, MolecularMass: 0.032}

	printDensity(iron)
	printDensity(oxygen)
}
