package main

import (
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
)

// Define an objective function we want to minimize.
func objective(trial goptuna.Trial) (float64, error) {
	// Define a search space of the input values.
	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)

	// Here is a two-dimensional quadratic function.
	// F(x1, x2) = (x1 - 2)^2 + (x2 + 5)^2
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
	)
	if err != nil {
		panic(err)
	}

	// Run an objective function 100 times to find a global minimum.
	err = study.Optimize(objective, 100)
	if err != nil {
		panic(err)
	}

	// Print the best evaluation value and the parameters.
	// Mathematically, argmin F(x1, x2) is (x1, x2) = (+2, -5).
	v, _ := study.GetBestValue()
	p, _ := study.GetBestParams()
	log.Printf("Best evaluation value=%f (x1=%f, x2=%f)",
		v, p["x1"].(float64), p["x2"].(float64))
}
