package bignum

import (
	"math/big"
)

//Problem defines the parameters of the algorithm
type Problem struct {
	BucketA *big.Int `json:"bucketa"`
	BucketB *big.Int `json:"bucketb"`
	Desired *big.Int `json:"desired"`
}

//NewProblem creates a new Problem from parameters describing each bucket size and the desired remainder.
//
// Parameters:
// a *bignum.Int Size of bucket A
// b *bignum.Int Size of bucket B
// d *bignum.Int Desired remainder from the solution
//
// Returns:
// p *Problem a pointer to the created object
func NewProblem(a, b, d *big.Int) (p *Problem) {
	p = new(Problem)
	p.BucketA = a
	p.BucketB = b
	p.Desired = d
	return p
}

//compareAandB returns the int result of A.Cmp(B) on the buckets
func (p *Problem) compareAandB() int {
	return p.BucketA.Cmp(p.BucketB)
}
