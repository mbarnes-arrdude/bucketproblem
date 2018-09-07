package bignum

import (
	bp "arrdude.com/bucketproblem"
	"math/big"
)

//SimulationState holds the state of a bucket at a particular step
// Values calculated are solution after accompanied operation.
type SimulationState struct {
	Operation     bp.SimulationOperation
	AmountBucketA *big.Int
	AmountBucketB *big.Int
}

//BucketStateList is an array of SimulationState pointers
type BucketStateList []*SimulationState

// Stores in order of operation.
type BucketStateCache struct {
	BucketStateList
	s          *Solution
	FillCount  *big.Int
	PourCount  *big.Int
	EmptyCount *big.Int
}

func newBucketState(from *big.Int, to *big.Int, operation bp.SimulationOperation, fromb bool) (s *SimulationState) {
	s = new(SimulationState)
	if fromb {
		s.AmountBucketA = new(big.Int).Set(to)
		s.AmountBucketB = new(big.Int).Set(from)
	} else {
		s.AmountBucketA = new(big.Int).Set(from)
		s.AmountBucketB = new(big.Int).Set(to)
	}
	s.Operation = operation
	return s
}

func newBucketStateCache() (b *BucketStateCache) {
	b = new(BucketStateCache)
	b.FillCount = big.NewInt(0)
	b.PourCount = big.NewInt(0)
	b.EmptyCount = big.NewInt(0)

	b.BucketStateList = make(BucketStateList, 0)
	return b
}

//generateSimulation() is an internal function that uses the results of the previously run extended Euclidian algorithm to generate a
// simulation of the solution.
//
// It derives transfer volume and approximated action count for each side then simulates the
// chosen solution to populate MinOperations *BucketStateList. The simulated number of "fill" operations is
// accumulated in FillCount. the Code represents the solution of running the algorithm and will be populated as
// ResultsOK if successful or another value of psuedo enum ResultCode if there is a problem or error.
//
// Output is truncated if the number of steps in the simulation exceeds MaxOperationsListSize
//
// Returns:
// ops BucketStateList has an underlying array of the bucket states per operation
func (s *Solution) generateSimulation(controller *ChannelController) {
	controller.state = controller.state | StageSimulation | Running | Initialized
	if !controller.mayContinue() {
		return
	}
	//s.PredictedStateCount = new(big.Int)
	//
	from := new(big.Int)
	to := new(big.Int)
	//
	capfrom := new(big.Int)
	capto := new(big.Int)
	//
	//countempties := new(big.Int)
	//
	////switch buckets for additive direction and set direction specific pretenses
	if s.FromB {
		capfrom.Set(s.Problem.BucketB)
		capto.Set(s.Problem.BucketA)
		//
		//	countempties.Mul(s.CountFromB, capfrom)
		//	countempties.Div(countempties, capto)
		//	//subtract a predicted empty if we will not overflow
		//	if s.Problem.Desired.Cmp(capfrom) != 1 {
		//		countempties.Sub(countempties, bigone)
		//	}
		//	s.PredictedStateCount.Set(s.CountFromB)
	} else {
		capfrom.Set(s.Problem.BucketA)
		capto.Set(s.Problem.BucketB)
		//
		//	countempties.Div(capfrom, capto)
		//	countempties.Mul(countempties, s.CountFromA)
		//
		//	if s.Problem.Desired.Cmp(s.Problem.BucketA) == 1 {
		//		countempties.Sub(countempties, bigone)
		//	}
		//	s.PredictedStateCount.Set(s.CountFromA)
	}
	//
	//s.PredictedStateCount.Add(s.PredictedStateCount, countempties)
	////multiply by steps per action
	//s.PredictedStateCount.Mul(s.PredictedStateCount, stepsPerAction)
	////add initial state
	//s.PredictedStateCount.Add(s.PredictedStateCount, big.NewInt(1))

	if !controller.mayContinue() {
		s.Operations.appendErrorBucket(bp.ProcessKilled, controller)
		return
	}

	s.Operations.appendInitialBucket(s.FromB, controller)

	//iterate the solution appending new states to the list
	for from.Cmp(s.Problem.Desired) != 0 && to.Cmp(s.Problem.Desired) != 0 {

		if !controller.mayContinue() {
			s.Operations.appendErrorBucket(bp.ProcessKilled, controller)
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

		if !controller.mayContinue() {
			s.Operations.appendErrorBucket(bp.ProcessKilled, controller)
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

func (blist BucketStateList) isFull() bool {
	return len(blist) >= bp.MaxOperationsListSize
}

func (s BucketStateCache) GetNextIndex() *big.Int {
	//index 0 is always InitialOp
	fullcount := new(big.Int)
	fullcount.Set(s.FillCount)
	fullcount.Add(fullcount, s.PourCount)
	fullcount.Add(fullcount, s.EmptyCount)
	fullcount.Add(fullcount, bigone) //next
	return fullcount
}

func (c *BucketStateCache) appendInitialBucket(fromb bool, controller *ChannelController) (full bool) {
	newstate := newBucketState(bigzero, bigzero, bp.InitialOp, fromb)
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}

func (c *BucketStateCache) appendNewBucket(index *big.Int, from *big.Int, to *big.Int, op bp.SimulationOperation, fromb bool, controller *ChannelController) (full bool) {
	newstate := newBucketState(new(big.Int).Set(from), new(big.Int).Set(to), op, fromb)
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() || op == bp.FinalOp {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}

func (c *BucketStateCache) appendSolvedBucket(index *big.Int, from *big.Int, to *big.Int, fromb bool, controller *ChannelController) (full bool) {
	newstate := newBucketState(new(big.Int).Set(from), new(big.Int).Set(to), bp.FinalOp, fromb)
	controller.simulationOperationCollector <- *newstate
	c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
	return true
}

func (c *BucketStateCache) appendErrorBucket(code bp.ResultCode, controller *ChannelController) (full bool) {
	newstate := newBucketState(bigzero, bigzero, bp.SimulationError, false)
	controller.solution.Code = code
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}
