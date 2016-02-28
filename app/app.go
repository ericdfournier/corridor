package main

import (
	"flag"
	"fmt"
	"github.com/ericdfournier/corridor"
	"os"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "corridor"
	app.Usage = "generate near optimal solutions to multi-objective corridor location problems using a parallel genetic algorithm"
	app.Author = "eric daniel fournier"
	app.Email = "me@ericdfournier.com"
	app.Commands = []cli.Command{
		{
			Name:      "test",
			ShortName: "t",
			Usage:     "run corridor package test suite",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:      "suppress",
					ShortName: "s",
					Usage:     "surpress detailed test results output",
					Action:    testSuppress,
				},
				cli.BoolFlag{
					Name:      "verbose",
					ShortName: "v",
					Usage:     "print detailed test results output",
					Action:    testVerbose,
				},
			},
		},
		{
			Name:      "benchmark",
			ShortName: "b",
			Usage:     "run corridor package benchmark suite",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:      "small",
					ShortName: "s",
					Usage:     "small population size (P = 1,000)",
					Action:    benchmarkSmall,
				},
				cli.BoolFlag{
					Name:      "medium",
					ShortName: "m",
					Usage:     "medium population size (P = 10,000)",
					Action:    benchmarkMedium,
				},
				cli.BoolFlag{
					Name:      "large",
					ShortName: "l",
					Usage:     "large population size (P = 100,000)",
					Action:    benchmarkLarge,
				},
				cli.BoolFlag{
					Name:      "monte_carlo",
					ShortName: "mc",
					Usage:     "perform monte carlo simulation (N = 100)",
					Action:    benchmarkMonteCarlo,
				},
			},
		},
		{
			Name:      "solve",
			ShortName: "s",
			Usage:     "solve corridor location problem",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:      "source_row",
					ShortName: "sr",
					Usage:     "source location row subscript [INT]",
					Value:     nil,
				},
				cli.IntFlag{
					Name:      "source_column",
					ShortName: "sc",
					Usage:     "source location column subscript [INT]",
					Value:     nil,
				},
				cli.IntFlag{
					Name:      "destination_row",
					ShortName: "dr",
					Usage:     "destination location row subscript [INT]",
					Value:     nil,
				},
				cli.IntFlag{
					Name:      "destination_column",
					ShortName: "dc",
					Usage:     "destination location column subscript [INT]",
					Value:     nil,
				},
				cli.StringFlag{
					Name:      "search_domain",
					ShortName: "sd",
					Usage:     ".CSV formatted binary representation of feasible search domain [FILEPATH]",
					Value:     nil,
				},
				cli.StringFlag{
					Name:      "objective_function",
					ShortName: "of",
					Usage:     ".CSV formatted floating point representation of objective function values [FILEPATH]",
					Value:     nil,
				},
				Action: solveCorridorLocationProblem,
			},
		},
	}
	app.Run(os.Args)
}
