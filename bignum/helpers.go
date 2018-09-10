package bignum

import (
	"fmt"
	"math/big"
	"strconv"
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
	fmt.Printf("Complexity: %v\n", s.Denominator)
	fmt.Printf("Denominator: %v\n", s.Denominator)
	fmt.Printf("MultInverseA: %v\n", s.MultInverseA)
	fmt.Printf("CountFromA: %v\n", s.CountFromA)
	fmt.Printf("CountFromB: %v\n", s.CountFromB)
	fmt.Printf("TvolumeA: %v\n", s.TvolumeA)
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

//Hash() delivers a hash of the problem in a sort domain that favors scale in sort and the Desired as identity
func (p *Problem) Hash() (hash int) {
	var apartbig string
	var apartsmall string
	var bpartbig string
	var bpartsmall string
	var dpartbig string
	var dpartsmall string

	ahex := p.BucketA.Text(16)
	if len(ahex) == 1 {
		apartbig = "0"
		apartsmall = ahex[:1]
	} else if len(ahex) == 2 {
		apartbig = ahex[:1]
		apartsmall = ahex[len(ahex)-1:]
	} else {
		apartbig = ahex[:1]
		apartsmall = ahex[len(ahex)-1:]
	}

	bhex := p.BucketB.Text(16)
	if len(bhex) == 1 {
		bpartbig = "0"
		bpartsmall = bhex[:1]
	} else {
		bpartbig = bhex[:1]
		bpartsmall = bhex[len(bhex)-1:]
	}

	dhex := p.Desired.Text(16)
	if len(dhex) == 1 {
		dpartbig = "0"
		dpartsmall = dhex[:1]
	} else {
		dpartbig = dhex[:1]
		dpartsmall = dhex[len(dhex)-1:]
	}

	scale := len(ahex)
	if scale > 0xff {
		scale = 0xff
	}
	spart := fmt.Sprintf("%02x", scale)
	fmt.Printf("scale %02x\n", scale)

	hexstr := spart + dpartbig + apartbig + bpartbig + bpartsmall + apartsmall + dpartsmall

	hashlong, err := strconv.ParseInt(hexstr, 16, 32)
	hash = int(hashlong)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	fmt.Println("Hash: ", hash)
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
