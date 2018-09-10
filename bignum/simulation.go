package bignum

import (
	bp "arrdude.com/bucketproblem"
	"math/big"
)

//SimulationState holds the state of a bucket at a particular step
// Values calculated are solution after accompanied operation.
type SimulationState struct {
	Idx           *big.Int               `json: idx`
	Operation     bp.SimulationOperation `json:"operation"`
	AmountBucketA *big.Int               `json:"bucketa"`
	AmountBucketB *big.Int               `json:"bucketb"`
}

//BucketStateList is an array of SimulationState pointers
type BucketStateList []*SimulationState

// Stores in order of operation.
type BucketStateCache struct {
	BucketStateList `json:"bucketstatelist"`
	s               *Solution
	FillCount       *big.Int `json:"fillcount"`
	PourCount       *big.Int `json:"pourcount"`
	EmptyCount      *big.Int `json:"emptycount"`
}

func newBucketState(idx *big.Int, from *big.Int, to *big.Int, operation bp.SimulationOperation, fromb bool) (s *SimulationState) {
	s = new(SimulationState)
	s.Idx = idx
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

func (b *BucketStateCache) GetLastOperation() (s *SimulationState) {
	s = b.BucketStateList[len(b.BucketStateList)-1]
	return s
}

//generateSimulation() simulates the chosen solution populating MinOperations *BucketStateList. The simulated number of
// operations is accumulated in FillCount PourCount and . the Code represents the solution of running the algorithm and
// will be populated as ResultsOK if successful or another value of psuedo enum ResultCode if there is a problem or error.
//
// Output is truncated if the number of steps in the simulation exceeds MaxOperationsListSize
//
// parameter:
// controller *ChannelController emits simulation events on channels as it traverses the solution
// ops BucketStateList has an underlying array of the bucket states per operation
func (s *Solution) generateSimulation(controller *ChannelController) {
	newsize := int(s.PredictedStateCount.Int64())
	if newsize < SimulationCollectorChannelSizeSmall {
		newsize = SimulationCollectorChannelSizeSmall
	} else if newsize > SimulationCollectorChannelSizeLarge {
		newsize = SimulationCollectorChannelSizeLarge
	}
	controller.simulationOperationCollector = make(chan SimulationState, newsize)
	controller.state = controller.state | StageSimulation | Running | Initialized
	if !controller.mayContinue() {
		return
	}
	from := new(big.Int)
	to := new(big.Int)

	capfrom := new(big.Int)
	capto := new(big.Int)

	////switch buckets for additive direction and set direction specific pretenses
	if s.FromB {
		capfrom.Set(s.Problem.BucketB)
		capto.Set(s.Problem.BucketA)
	} else {
		capfrom.Set(s.Problem.BucketA)
		capto.Set(s.Problem.BucketB)
	}

	if !controller.mayContinue() {
		s.Operations.appendErrorBucket(s.Operations.GetNextIndex(), bp.ProcessKilled, controller)
		return
	}

	s.Operations.appendInitialBucket(s.FromB, controller)

	//iterate the solution appending new states to the list
	for from.Cmp(s.Problem.Desired) != 0 && to.Cmp(s.Problem.Desired) != 0 {

		if !controller.mayContinue() {
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

		if !controller.mayContinue() {
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
	newstate := newBucketState(bigzero, bigzero, bigzero, bp.InitialOp, fromb)
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}

func (c *BucketStateCache) appendNewBucket(index *big.Int, from *big.Int, to *big.Int, op bp.SimulationOperation, fromb bool, controller *ChannelController) (full bool) {
	newstate := newBucketState(index, new(big.Int).Set(from), new(big.Int).Set(to), op, fromb)
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() || op == bp.FinalOp {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}

func (c *BucketStateCache) appendSolvedBucket(index *big.Int, from *big.Int, to *big.Int, fromb bool, controller *ChannelController) (full bool) {
	newstate := newBucketState(index, new(big.Int).Set(from), new(big.Int).Set(to), bp.FinalOp, fromb)
	controller.simulationOperationCollector <- *newstate
	c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
	return true
}

func (c *BucketStateCache) appendErrorBucket(index *big.Int, code bp.ResultCode, controller *ChannelController) (full bool) {
	newstate := newBucketState(index, bigzero, bigzero, bp.SimulationError, false)
	controller.Solution.Code = code
	controller.simulationOperationCollector <- *newstate
	if !c.isFull() {
		c.BucketStateList = append(c.BucketStateList, []*SimulationState{newstate}...)
		return false
	}
	return true
}
