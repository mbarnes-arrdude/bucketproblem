package main

import (
	bp "arrdude.com/bucketproblem"
	bucketproblem "arrdude.com/bucketproblem/bignum"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"math/big"
	"os"
)

type Job struct {
}

func main() {
	app := cli.NewApp()
	app.Usage = "calcbucket bucketa bucketb desired"
	app.Action = func(c *cli.Context) error {
		//job := []*big.Int{big.NewInt(5),
		//	big.NewInt(3),
		//	big.NewInt(4)}
		sbucketa := c.Args().Get(0)
		sbucketb := c.Args().Get(1)
		sdesired := c.Args().Get(2)

		bucketa, avalid := new(big.Int).SetString(sbucketa, 10)
		bucketb, bvalid := new(big.Int).SetString(sbucketb, 10)
		desired, dvalid := new(big.Int).SetString(sdesired, 10)

		if len(c.Args()) < 3 {
			return errors.New("not enough arguments")
		}

		if avalid && bvalid && dvalid {
			problem := bucketproblem.NewProblem(bucketa, bucketb, desired)
			solution := bucketproblem.NewSolution(problem)

			fmt.Println("Running GCD Solution and Simulation")

			bucketproblem.GetRunningSolutionProcessor(solution)

			sdirection := "Subtractive (A -> B)"
			if solution.FromB {
				sdirection = "Additive (A <- B)"
			}

			fmt.Printf("Problem (Hash: %x)\n", problem.Hash())
			fmt.Printf("- Bucket A: %v\n", problem.BucketA)
			fmt.Printf("- Bucket B: %v\n", problem.BucketB)
			fmt.Printf("- Desired: %v\n", problem.Desired)
			fmt.Println("Solution")
			fmt.Printf("- Result: %s\n", solution.Code)
			fmt.Printf("- GCD: %v\n", solution.Denominator)
			fmt.Printf("- Direction: %s\n", sdirection)
			fmt.Printf("- Total Steps: %v\n", new(big.Int).Sub(solution.Operations.GetNextIndex(), big.NewInt(1)))
			fmt.Println("Simulation Table")
			truncated := false
			lastbucket := solution.Operations.BucketStateList[0]
			fmt.Printf("      Idx        | %16s |   | B\n", "A")
			fmt.Println("----------------------------------------------------------")
			for idx, bucket := range solution.Operations.BucketStateList {
				if idx == len(solution.Operations.BucketStateList)-1 {
					if lastbucket.AmountBucketA.Cmp(bucket.AmountBucketA) != 0 || lastbucket.AmountBucketB.Cmp(bucket.AmountBucketB) != 0 {
						lastbucket.Operation = bp.Truncated
						truncated = true
					}
				}
				if idx != 0 {
					printTableEntry(idx, lastbucket, solution.FromB)
				}
				lastbucket = bucket
			}
			if truncated {
				fmt.Println("             ...                ...       ...")
			}
			fmt.Println("----------------------------------------------------------")
			finalSimulationIdx := new(big.Int).Sub(solution.Operations.GetNextIndex(), big.NewInt(1))
			printLastTableEntry(finalSimulationIdx, lastbucket)
			if truncated {
				fmt.Printf("Content Truncated to %d entries (max entries). Use `runbucket` for better handling of large simulations.\n", bp.MaxOperationsListSize)
			}

			//solution.Spew()
		} else {
			message := fmt.Sprint("Arguments could not be parsed:\n")

			if !avalid {
				message += fmt.Sprintf("bucketa invalid: could not parse big.Int from %s\n")
			}

			if !bvalid {
				message += fmt.Sprintf("bucketa invalid: could not parse big.Int from %s\n")
			}

			if !dvalid {
				message += fmt.Sprintf("bucketa invalid: could not parse big.Int from %s\n")
			}
			return errors.New(message)
		}
		return nil
	}
	app.Run(os.Args)
}

func printTableEntry(idx int, lastbucket *bucketproblem.SimulationState, fromb bool) (int, error) {
	var dira = ' '
	var dirx = ' '
	var dirb = ' '
	if lastbucket.Operation == bp.Fill {
		if fromb {
			dirb = lastbucket.Operation.Rune()
		} else {
			dira = lastbucket.Operation.Rune()
		}
	} else if lastbucket.Operation == bp.Empty {
		if fromb {
			dira = lastbucket.Operation.Rune()
		} else {
			dirb = lastbucket.Operation.Rune()
		}
	} else if lastbucket.Operation == bp.Pour {
		if fromb {
			dira = lastbucket.Operation.Rune()
		} else {
			dirb = lastbucket.Operation.Rune()
		}
	} else {
		dirx = lastbucket.Operation.Rune()
	}
	return fmt.Printf("%15d) | %16v |%c%c%c| %v\n", idx-1, lastbucket.AmountBucketA, dira, dirx, dirb, lastbucket.AmountBucketB)
}

func printLastTableEntry(idx *big.Int, lastbucket *bucketproblem.SimulationState) (int, error) {
	return fmt.Printf("%15v) | %16v | %c | %v\n", idx, lastbucket.AmountBucketA, lastbucket.Operation.Rune(), lastbucket.AmountBucketB)
}
