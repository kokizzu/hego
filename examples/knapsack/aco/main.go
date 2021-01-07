package main

import (
	"fmt"
	"math"

	"github.com/ccssmnn/hego"
)

func knapsack(selection []bool, values, weights []float64, maxWeight float64) float64 {
	value := 0.0
	weight := 0.0
	for index, choice := range selection {
		if choice {
			value += values[index]
			weight += weights[index]
		}
	}
	if weight > maxWeight {
		value -= 10 * (weight - maxWeight)
	}
	return value
}

var values = []float64{69, 21, 33, 50, 89, 58, 27, 95, 52, 68, 26, 84, 46, 28, 25, 81, 82, 27, 50, 61}
var weights = []float64{6, 1, 1, 4, 9, 7, 3, 5, 7, 7, 9, 4, 4, 4, 8, 7, 7, 6, 5, 3}
var maxWeight = 30.0
var pheromones []float64

type ant struct {
	weight    float64
	value     float64
	selection []bool
}

func (a *ant) Init() {
	a.weight = 0.0
	a.value = 0.0
	a.selection = make([]bool, len(weights))
}

func (a *ant) Step(next int) bool {
	a.weight += weights[next]
	a.value += values[next]
	a.selection[next] = true
	done := true
	for i, choice := range a.selection {
		if !choice {
			if a.weight+weights[i] < maxWeight {
				done = false
			}
		}
	}
	return done
}

func (a *ant) PerceivePheromone() []float64 {
	res := make([]float64, len(pheromones))
	copy(res, pheromones)
	// do not take items that are already taken
	for i, choice := range a.selection {
		if choice {
			res[i] = 0.0
		}
	}
	// do not take items, if their weight would increase load too much
	for i := range pheromones {
		if a.weight+weights[i] > maxWeight {
			res[i] = 0.0
		}
	}
	return res
}

func (a *ant) DropPheromone(performance float64) {
	for index, choice := range a.selection {
		if choice {
			pheromones[index] += 0.2
		}
	}
}

func (a *ant) Evaporate(factor, min float64) {
	for i := range pheromones {
		pheromones[i] = math.Max(min, pheromones[i]*factor)
	}
}

func (a *ant) Performance() float64 {
	return -knapsack(a.selection, values, weights, maxWeight)
}

func main() {
	initialPheromone := 1.0
	pheromones = make([]float64, len(weights))
	for i := range pheromones {
		pheromones[i] = initialPheromone
	}
	population := make([]hego.Ant, 10)
	for i := range population {
		population[i] = &ant{}
	}

	settings := hego.ACOSettings{}
	settings.Evaporation = 0.95
	settings.MinPheromone = 0.01
	settings.MaxIterations = 100
	settings.Verbose = settings.MaxIterations / 10

	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
	} else {
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
