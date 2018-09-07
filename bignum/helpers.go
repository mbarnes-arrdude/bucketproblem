package bignum

import (
	bp "arrdude.com/bucketproblem"
	"fmt"
	"math/big"
)

//no reason to reinitialize these in memory each time they are used
var bigzero = big.NewInt(0)
var bigone = big.NewInt(1)
var bignegone = big.NewInt(-1)

//Used as final LU for generateSimulation()
var stepsPerAction = big.NewInt(2)

//Spew() prints to stdout a human meaningful representation of a Solution
// Output contains newlines.
// WARNING: Helper function only for development. Not versioned. Do not use in production.

func (s Solution) Spew() {
	fmt.Println("\n=====")
	s.Problem.Spew()
	fmt.Printf("Code: %s\n", s.Code)
	fmt.Printf("Denominator: %v\n", s.Denominator)
	fmt.Printf("MultInverseA: %v\n", s.MultInverseA)
	fmt.Printf("MultInverseB: %v\n", s.MultInverseB)
	fmt.Printf("CountFromA: %v\n", s.CountFromA)
	fmt.Printf("CountFromB: %v\n", s.CountFromB)
	fmt.Printf("TvolumeA: %v\n", s.TvolumeA)
	fmt.Printf("TvolumeB: %v\n", s.TvolumeB)
	fmt.Printf("FromB: %t\n", s.FromB)
	fmt.Printf("PredictedStateCount: %v\n", s.PredictedStateCount)
	s.Operations.Spew()
	fmt.Println("=====")
}

//Spew() prints to stdout a human meaningful representation of the Problem and the Solution
// Output contains newlines.
// WARNING: Helper function only for development. Not versioned. Do not use in production.
func (p *Problem) Spew() {
	fmt.Printf("=== Problem: A: %v, B: %v, Desired: %v ===\n", p.BucketA, p.BucketB, p.Desired)
}

//Hash() delivers a hash of the problem in a sort domain that favors the Desired as identity, and a prime number modulus
// of a resource identity
func (p *Problem) Hash() (hash int) {
	//Identity is most likely to be important for LU by job. The purpose of the algorithm is anchored on the Desired
	// outcome. When looking for "similar" problems this attribute is the most intuitively significant to application.
	// after all it is usually easier to adjust your envionmental parameters (water jugs are easy to find) than adjust
	// your goal.

	// 8bit randomish mod' of ^2 domain for uniqueness
	// As a constant for this process the resulting partition scatter will group by associated processes clumping on
	// resource associations resources closest to this one are most likely to use it. Scaled resources in micro-services
	// environments will generally be divided into clusters of resources that the same user will use. A process id will
	// generate a short-lived hash that will be seen most often by more closely clustered resources in horizontal scaling.
	//
	// The value is also likely the least significant in any purpose of the requestor being unpredictably arbitrary to
	// an applied purpose. But it is most significant in grouping for performance metrics and this reinforces the
	// clumpiness of it to significant resources when scaled and the system measured. Not only will it be most likely
	// close to other system resources the user would connect to to access the value but it also is an identity key for
	// resource management.
	ppartshort := bp.GetProcessHash() % (0x100 - 1)

	//bignum ints provide too large a potential domain space for encoding atomic identity. not only does it crowd entropy
	// for an index, it exceeds it. Any cache LU by parameters is likely going to require a secondary linear search.
	// for the potential bracket of parameters of the problem which fly towards infinity. The reduction of the domain
	// means that some hash collision will occur. The secondary search after hash lookup should be linear with most
	// significant

	// Identity for BucketA: ^2 linear domain for clumpiness
	// Least significant of paradomain for sort because it is a multiplier in the purpose
	// and larger size means smaller entropy bucket for scatter
	apartlong := p.BucketA.Int64()      //scale to 64 bits
	apartshort := int(apartlong)        //scale to 32 bits
	apartshort = apartshort / 0x1000000 //scale to 8 bits

	// Identity for BucketB: ^2 linear domain for clumpiness
	// Most significant of paradomain because it is a quotient in the purpose
	// and smaller size means larger scatter entropy than BucketA
	bpartlong := p.BucketA.Int64()      //scale to 64 bits
	bpartshort := int(bpartlong)        //scale to 32 bits
	bpartshort = bpartshort / 0x1000000 //scale to 8 bits

	// identity domain for Desired ^2 linear domain for clumpiness
	// Primary domain because it is most purposeful << most likely sorted for problems that are similar
	dpartlong := p.Desired.Int64()      //scale to 64 bits
	dpartshort := int(dpartlong)        //scale to 32 bits
	dpartshort = dpartshort / 0x1000000 //scale to 8 bits

	hash = ((dpartshort << 0x18) & 0xff000000) |
		((bpartshort << 0x10) & 0xff0000) |
		((apartshort << 0x8) & 0xff00) |
		(ppartshort & 0xff)
	// identity
	return hash
}

//Spew() prints to stdout a human meaningful representation of a BucketStateList
// Output contains newlines.
// WARNING: Helper function only for development. Not versioned. Do not use in production
func (blist *BucketStateList) Spew() {
	for idx, state := range *blist {
		fmt.Printf("%d) %s\nA: %v B: %v\n", idx, state.Operation, state.AmountBucketA, state.AmountBucketB)
	}
}

//Spew() prints to stdout a human meaningful representation of a BucketStateCache
// Output contains newlines.
// WARNING: Helper function only for development. Not versioned. Do not use in production
func (c *BucketStateCache) Spew() {
	fmt.Printf("FillCount: %v\n", c.FillCount)
	fmt.Printf("PourCount: %v\n", c.PourCount)
	fmt.Printf("EmptyCount: %v\n", c.EmptyCount)
	c.BucketStateList.Spew()
}
