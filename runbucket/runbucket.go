package main

import (
	bp "arrdude.com/bucketproblem"
	"arrdude.com/bucketproblem/bignum"
	"fmt"
	"math/big"
	"sync"
)

type Chans struct {
	controlChannel *chan bignum.ProcessControlOperation
	stateChannel   chan bignum.ProcessControlOperation
	resultsChannel chan bignum.SimulationState
	running        bool
}

func newchans() (c *Chans) {
	r := new(Chans)
	r.stateChannel = make(chan bignum.ProcessControlOperation, 2)
	r.resultsChannel = make(chan bignum.SimulationState, 50000)
	return r
}

func main() {

	biga, _ := new(big.Int).SetString("1000000000000000000000000000000000000000000001", 10)
	bigb, _ := new(big.Int).SetString("100000000000000000000000000000000000000000000", 10)
	bigc, _ := new(big.Int).SetString("10003", 10)

	jobs := [][]*big.Int{
		//{big.NewInt(5),
		//	big.NewInt(3),
		//	big.NewInt(3)},

		//{big.NewInt(17),
		//	big.NewInt(12),
		//	big.NewInt(9)},

		//{big.NewInt(7),
		//	big.NewInt(5),
		//	big.NewInt(7)},

		//{big.NewInt(111),
		//	big.NewInt(11),
		//	big.NewInt(3)},

		//subtractive no empty
		//		{big.NewInt(8),
		//			big.NewInt(6),
		//			big.NewInt(2)},

		//additive no empty
		//{big.NewInt(101),
		//	big.NewInt(3),
		//	big.NewInt(12)},

		//[]*big.Int{big.NewInt(15),
		//	big.NewInt(5),
		//	big.NewInt(10)},

		{biga,
			bigb,
			bigc},

		//{big.NewInt(9),
		//	big.NewInt(12),
		//	big.NewInt(7)},
	}
	fmt.Println("\nStarting jobs")

	for idx, job := range jobs {
		name := fmt.Sprintf("%05dxxx", idx)
		fmt.Printf("Doing Job: %s\n", name)
		problem := bignum.NewProblem(job[0], job[1], job[2])

		chans := newchans()
		solution := bignum.NewSolution(problem)
		controller := bignum.GetIdleSolutionProcessor(name, solution, &chans.stateChannel, &chans.resultsChannel)
		controlChannel := controller.GetStopStartChannel()
		var wg sync.WaitGroup
		wg.Add(2)
		go chans.doListenSimulation(&wg)
		go chans.doProcStartStop(controlChannel, &wg)
		*controlChannel <- bignum.Start
		wg.Wait()
		solution.Spew()
	}
}

func (c *Chans) doProcStartStop(controlChannel *chan bignum.ProcessControlOperation, group *sync.WaitGroup) {
	var running = true
	defer func() {
		group.Done()
	}()
	for running {
		select {
		case op := <-c.stateChannel:
			fmt.Printf("FLOW: %s\n", op)
			if int(op)&int(bignum.Error) > 0 {
				running = false
				fmt.Println("ERROR!")
				return
			}
		default:
		}
	}
}

func (c *Chans) doListenSimulation(group *sync.WaitGroup) {
	var doit = true
	defer func() {
		fmt.Println("Done Listening Simulation")
		group.Done()
	}()
	for doit {
		//fmt.Println("Running")
		select {
		case bucket := <-c.resultsChannel:
			fmt.Printf("ACTION: %s\n", bucket)
			if bucket.Operation >= bp.FinalOp {
				fmt.Println("Done")
				doit = false
				break
			}
		default:
		}
	}
}
