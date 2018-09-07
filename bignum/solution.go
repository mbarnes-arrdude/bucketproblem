package bignum

import (
	p "arrdude.com/bucketproblem"
	"math/big"
)

//Solution is the data and metadata of a completed Problem
// Return object of Problem.Solution()
type Solution struct {
	Problem             *Problem
	Denominator         *big.Int
	MultInverseA        *big.Int
	MultInverseB        *big.Int
	Code                p.ResultCode
	TvolumeA            *big.Int
	TvolumeB            *big.Int
	CountFromA          *big.Int
	CountFromB          *big.Int
	FromB               bool
	PredictedStateCount *big.Int

	Operations *BucketStateCache
}

//NewSolution solves and creates a completed *Solution using the values of a *Problem for its parameters.
// It calculates GCD and multiplicative inverse values for each side using the extended Euclidean
// algorithm as implemented in math/bignum. It determines fill count for each by projecting MulInv against the desired
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
	r.Denominator = new(big.Int)
	r.MultInverseA = new(big.Int)
	r.MultInverseB = new(big.Int)
	r.CountFromA = new(big.Int)
	r.CountFromB = new(big.Int)
	r.TvolumeB = new(big.Int)
	r.TvolumeA = new(big.Int)

	r.Operations = newBucketStateCache()

	return r
}

func (s *Solution) compareCountFromAandCountFromB() int {
	return s.CountFromA.Cmp(s.CountFromB)
}
