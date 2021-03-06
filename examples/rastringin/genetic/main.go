package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
	"github.com/ccssmnn/hego/mutate"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

// genome is a vector of float values
type genome []float64

// Crossover returns a child genome which is a combination of the current and other
// genome. Here an the Arithmetic crossover operation is used
func (g genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.Arithmetic(g, other.(genome), [2]float64{-0.5, 1.5}))
}

// Mutate adds variation to a genome. In this case we add gaussian noise
func (g genome) Mutate() hego.Genome {
	return genome(mutate.Gauss(g, 0.5))
}

// Fitness is called to evaluate the objective functino value. Lower is better
func (g genome) Fitness() float64 {
	return rastringin(g[0], g[1])
}

func main() {
	// initialize population
	population := make([]hego.Genome, 100)
	for i := range population {
		population[i] = genome{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}
	}

	// set algorithm behaviour here
	settings := hego.GASettings{}
	settings.MutationRate = 0.5
	settings.Elitism = 5
	settings.MaxIterations = 100
	settings.Verbose = 10

	// call genetic algorithm
	result, err := hego.GA(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
		return
	}

	// extract solution and print result
	solution := result.BestGenome.(genome)
	fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	fmt.Printf("Minimum found at x = [%v, %v] with f(x) = %v\n", solution[0], solution[1], rastringin(solution[0], solution[1]))
}
