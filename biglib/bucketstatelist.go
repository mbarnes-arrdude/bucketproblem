package biglib

import (
	bp "arrdude.com/bucketproblem"
	"math/big"
)

//SimulationState holds the state of a bucket at a particular step of the simulation.
// Bucket values are the result of the accompanied operation.
type SimulationState struct {
	Idx           *big.Int               `json:"idx"`
	Operation     bp.SimulationOperation `json:"operation"`
	AmountBucketA *big.Int               `json:"bucketa"`
	AmountBucketB *big.Int               `json:"bucketb"`
}

//BucketStateList is an array of SimulationState pointers. It is populated by the simulation as run by the controller.
type BucketStateList []*SimulationState

//BucketStateCache holds the list of BucketStates generated by the simulation as well as actual counts of operations.
// Stores in order of operation.
type BucketStateCache struct {
	BucketStateList `json:"bucketstatelist"`
	//aggregators for simulation event counts
	FillCount  *big.Int `json:"fillcount"`
	PourCount  *big.Int `json:"pourcount"`
	EmptyCount *big.Int `json:"emptycount"`
}

func (blist BucketStateList) isFull() bool {
	return len(blist) >= bp.MaxOperationsListSize
}

//GetLastOperation() returns the last entry in the simulation state cache
//
//Returns:
//s *SimulationState //last state object in the cache
//nil if the cache is empty
func (b *BucketStateCache) GetLastOperation() (s *SimulationState) {
	if len(b.BucketStateList) < 1 {
		return nil
	}
	s = b.BucketStateList[len(b.BucketStateList)-1]
	return s
}

//Returns the actual sum of operations run in the simulation.
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
