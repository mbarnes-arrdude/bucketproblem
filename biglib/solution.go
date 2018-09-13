package biglib

import (
	p "arrdude.com/bucketproblem"
	"math/big"
)

//Solution is the data and metadata of a completed Problem
// Return object of Problem.Solution()
type Solution struct {
	Problem             *Problem     `json:"problem"`
	Complexity          *big.Int     `json:"complexity"`
	Denominator         *big.Int     `json:"denominator"`
	MultInverseA        *big.Int     `json:"multinva"`
	MultInverseB        *big.Int     `json:"multinvb"`
	Code                p.ResultCode `json:"result"`
	TvolumeA            *big.Int     `json:"tvolumea"`
	TvolumeB            *big.Int     `json:"tvolumeb"`
	CountFromA          *big.Int     `json:"countfroma"`
	CountFromB          *big.Int     `json:"countfromb"`
	FromB               bool         `json:"fromb"`
	PredictedStateCount *big.Int     `json:"predictedstatecount"`
	GCDNanoTime         int64        `json:"gcdnanotime"`

	Operations *BucketStateCache `json:"operations"`
}

//NewSolution solves and creates a completed *Solution using the values of a *Problem for its parameters.
// It calculates GCD and multiplicative inverse values for each side using the extended Euclidean
// algorithm as implemented in math/biglib. It determines fill count for each by projecting MulInv against the desired
// amount moding against the bucket size. Internally it then calls generateSimulation() to simulate and record the solution.
//
// Arguments:
// problem *Problem
//
// Returns:
// r *Solution
func NewSolution(problem *Problem) (r *Solution) {
	r = new(Solution)
	r.Problem = problem
	r.Complexity = new(big.Int)
	r.Denominator = new(big.Int)
	r.MultInverseA = new(big.Int)
	r.MultInverseB = new(big.Int)
	r.CountFromA = new(big.Int)
	r.CountFromB = new(big.Int)
	r.TvolumeA = new(big.Int)
	r.TvolumeB = new(big.Int)

	r.Operations = newBucketStateCache()

	return r
}

//compareCountFromAandCountFromB returns the int result of comparing the predicted counts if pouring in the simulation
// were "From" bucket A or "From" bucket B in Problem. The predictions are calculated during (s *Solution)generateGCD().
func (s *Solution) compareCountFromAandCountFromB() int {
	return s.CountFromA.Cmp(s.CountFromB)
}

func (s *Solution) GetComplexityScale() int {
	text := s.Complexity.Text(10)
	return len(text)
}
