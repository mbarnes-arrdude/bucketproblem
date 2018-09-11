package bignum

import (
	bp "arrdude.com/bucketproblem"
	"math/big"
)

//generateSimulation() simulates the solution populating MinOperations *BucketStateList. The simulated number of
// operations is accumulated in FillCount PourCount and EmptyCounty. the Code represents the result of running the algorithm and
// will be populated as ResultsOK if successful or another value of ResultCode if there is a problem or error.
//
// Output is truncated if the number of steps in the simulation exceeds MaxOperationsListSize
//
//Parameters:
//controller *ChannelController directs the simulation through channels as it traverses the solution. The cache on the
// controller receives events synchronously and is responsible for relaying the states if such functionality exists.
func (s *Solution) generateSimulation(controller *ChannelController) {
	newsize := int(s.PredictedStateCount.Int64())
	if newsize < SimulationCollectorChannelSizeSmall {
		newsize = SimulationCollectorChannelSizeSmall
	} else if newsize > SimulationCollectorChannelSizeLarge {
		newsize = SimulationCollectorChannelSizeLarge
	}
	controller.simulationOperationCollector = make(chan SimulationState, newsize)
	controller.state = controller.state | StageSimulation | Running | Initialized
	if !controller.MayContinue() {
		return
	}
	from := new(big.Int)
	to := new(big.Int)

	capfrom := new(big.Int)
	capto := new(big.Int)

	////switch buckets for additive direction
	if s.FromB {
		capfrom.Set(s.Problem.BucketB)
		capto.Set(s.Problem.BucketA)
	} else {
		capfrom.Set(s.Problem.BucketA)
		capto.Set(s.Problem.BucketB)
	}

	if !controller.MayContinue() {
		s.Operations.appendErrorBucket(s.Operations.GetNextIndex(), bp.ProcessKilled, controller)
		return
	}

	s.Operations.appendInitialBucket(s.FromB, controller)

	//iterate the solution appending new states to the list
	for from.Cmp(s.Problem.Desired) != 0 && to.Cmp(s.Problem.Desired) != 0 {

		if !controller.MayContinue() {
			s.Operations.appendErrorBucket(s.Operations.GetNextIndex(), bp.ProcessKilled, controller)
			return
		}

		toremain := new(big.Int)
		quanttx := new(big.Int)
		toremain.Sub(capto, to)

		//determine quantity to transfer
		if from.Cmp(toremain) != -1 {
			quanttx.Set(toremain)
		} else {
			quanttx.Set(from)
		}

		//if there is any available and room pour
		// skips initial condition of A:0 B:0
		if quanttx.Cmp(bigzero) == 1 {
			to.Add(to, quanttx)
			from.Sub(from, quanttx)
			s.Operations.appendNewBucket(s.Operations.GetNextIndex(), from, to, bp.Pour, s.FromB, controller)
			s.Operations.PourCount.Add(s.Operations.PourCount, bigone)
		}

		//done?
		if from.Cmp(s.Problem.Desired) == 0 || to.Cmp(s.Problem.Desired) == 0 {
			break
		}

		if !controller.MayContinue() {
			s.Operations.appendErrorBucket(s.Operations.GetNextIndex(), bp.ProcessKilled, controller)
			return
		}

		//fill from if empty
		if from.Cmp(big.NewInt(0)) == 0 {
			from.Set(capfrom)
			s.Operations.appendNewBucket(s.Operations.GetNextIndex(), from, to, bp.Fill, s.FromB, controller)
			s.Operations.FillCount.Add(s.Operations.FillCount, bigone)
		}

		//empty to if full
		if to.Cmp(capto) == 0 {
			to.SetInt64(0)
			s.Operations.appendNewBucket(s.Operations.GetNextIndex(), from, to, bp.Empty, s.FromB, controller)
			s.Operations.EmptyCount.Add(s.Operations.EmptyCount, bigone)
		}
	}
	s.Operations.appendSolvedBucket(s.Operations.GetNextIndex(), from, to, s.FromB, controller)
}
