package hego

import (
	"math"
	"math/rand"
	"testing"
)

type genome float64

func (b genome) Crossover(other Genome) Genome {
	return b + genome(rand.Float64())*(b-other.(genome))
}

func (b genome) Fitness() float64 {
	return float64(b * b)
}

func (b genome) Mutate() Genome {
	return b + genome(rand.NormFloat64())
}

func TestVerifyGASettings(t *testing.T) {
	settings := GASettings{}
	settings.MutationRate = 1.1
	err := settings.Verify()
	if err == nil {
		t.Error("verifification should fail for invalid mutation rate")
	}
	settings.MutationRate = 0.5
	settings.Elitism = -1
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for negative elitism count")
	}
	settings.Elitism = 1
	settings.Selection = TournamentSelection
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail, when TournamentSelection is selected but tournament size not provided / <2")
	}
	settings.TournamentSize = 2
	err = settings.Verify()
	if err != nil {
		t.Errorf("for valid settings verification should pass, got: %v for %v", err, settings)
	}
}

func TestGA(t *testing.T) {
	populationSize := 10
	population := make([]Genome, populationSize)
	settings := GASettings{}
	settings.MutationRate = 1.1
	res, err := GA(population, settings)
	if err == nil {
		t.Error("verifification should fail for invalid mutation rate")
	}

	settings.MutationRate = 0.1
	settings.Elitism = 1
	settings.MaxIterations = 100
	settings.Verbose = 10
	settings.KeepIntermediateResults = true

	for i := range population {
		candidate := genome(-20.0 + 40.0*rand.Float64())
		population[i] = candidate
	}

	res, err = GA(population, settings)

	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	if res.Iterations != settings.MaxIterations {
		t.Errorf("unexpected number of iterations. Expected %v, got %v", settings.MaxIterations, res.Iterations)
	}
	if math.Abs(res.BestFitness) > 0.5 {
		t.Error("unexpected solution found")
	}
}

func TestWeightedChoice(t *testing.T) {
	weights := []float64{1.0, 2.0, 0.0}
	n := 20
	choices := weightedChoice(weights, n)
	if len(choices) != n {
		t.Errorf("expected number of choices to be %v, got %v", n, len(choices))
	}
	for _, choice := range choices {
		if choice == 2 {
			t.Error("2 should not be a choice")
		}
	}
	weights = []float64{0.0, 0.0, 0.0}
	choices = weightedChoice(weights, n)
	if choices[0] != -1 {
		t.Errorf("weightedChoice should return -1 if probability of every choice is 0, got: %v", choices[0])
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("weighted choice should fail for 0 or less choices")
		}
	}()
	weightedChoice(weights, 0)
}

func TestBinaryWeightedChoice(t *testing.T) {
	weights := []float64{1.0, 2.0, 0.0}
	n := 20
	choices := binaryWeightedChoice(weights, n)
	if len(choices) != n {
		t.Errorf("expected number of choices to be %v, got %v", n, len(choices))
	}
	for _, choice := range choices {
		if choice == 2 {
			t.Error("2 should not be a choice")
		}
	}
	weights = []float64{0.0, 0.0, 0.0}
	choices = binaryWeightedChoice(weights, n)
	if choices[0] != -1 {
		t.Errorf("binaryWeightedChoice should return -1 if probability of every choice is 0, got: %v", choices[0])
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("weighted choice should fail for 0 or less choices")
		}
	}()
	binaryWeightedChoice(weights, 0)
}

func TestTournament(t *testing.T) {
	weights := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	index := tournament(weights)
	if index != 0 {
		t.Errorf("expected index 0 to win the tournament, got %v", index)
	}
}

func TestSelections(t *testing.T) {
	pop := population{
		candidate{genome: genome(1.0), fitness: 4.0},
		candidate{genome: genome(1.0), fitness: 3.0},
		candidate{genome: genome(1.0), fitness: 2.0},
		candidate{genome: genome(1.0), fitness: 1.0},
	}
	settings := GASettings{}
	settings.Selection = FitnessProportionalSelection
	parendIds := pop.selectParents(&settings)
	if len(parendIds) != len(pop) {
		t.Errorf("expected length of parentIds after selection to be equal to population size, got: %v", len(parendIds))
	}
	settings.Selection = TournamentSelection
	settings.TournamentSize = 2
	parendIds = pop.selectParents(&settings)
	if len(parendIds) != len(pop) {
		t.Errorf("expected length of parentIds after selection to be equal to population size, got: %v", len(parendIds))
	}
	settings.Selection = RankBasedSelection
	parendIds = pop.selectParents(&settings)
	if len(parendIds) != len(pop) {
		t.Errorf("expected length of parentIds after selection to be equal to population size, got: %v", len(parendIds))
	}
	if pop[0].fitness != 1.0 {
		t.Error("expected population to be sorted after rank based selection")
	}
}
