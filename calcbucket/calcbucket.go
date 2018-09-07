package main

import (
	bucketproblem "arrdude.com/bucketproblem/bignum"
	"fmt"
	"math/big"
)

func main() {

	//biga, _ := new(big.Int).SetString("1000000000000000000000000000000000000000000001", 10)
	//bigb, _ := new(big.Int).SetString("100000000000000000000000000000000000000000000", 10)
	//bigc, _ := new(big.Int).SetString("1000003", 10)

	jobs := [][]*big.Int{
		{big.NewInt(5),
			big.NewInt(3),
			big.NewInt(4)},

		//[]*bignum.Int{bignum.NewInt(17),
		//	bignum.NewInt(12),
		//	bignum.NewInt(9)},

		//[]*bignum.Int{bignum.NewInt(11),
		//	bignum.NewInt(3),
		//	bignum.NewInt(9)},

		//subtractive no empty
		//		[]*bignum.Int{bignum.NewInt(2),
		//			bignum.NewInt(8),
		//			bignum.NewInt(6)},

		//additive no empty
		//[]*bignum.Int{bignum.NewInt(101),
		//	bignum.NewInt(3),
		//	bignum.NewInt(12)},

		//[]*bignum.Int{bignum.NewInt(15),
		//	bignum.NewInt(5),
		//	bignum.NewInt(5)},

		//[]*big.Int{biga,
		//	bigb,
		//	bigc},

		//[]*bignum.Int{bignum.NewInt(9),
		//	bignum.NewInt(12),
		//	bignum.NewInt(7)},
	}
	fmt.Println("\nStarting jobs")

	for idx, job := range jobs {
		fmt.Printf("Doing Job: %d\n", idx)
		problem := bucketproblem.NewProblem(job[0], job[1], job[2])
		solution := bucketproblem.NewSolution(problem)
		bucketproblem.GetRunningSolutionProcessor(solution)
		//controlchannel := controller.GetStopStartChannel()
		//
		//for !solution.IsDone(){
		//	fmt.Printf("Status: %s for %v\n", solution.GetState(), solution)
		//	time.Sleep(5 * time.Second)
		//}
		solution.Spew()
	}
	fmt.Println("Jobs Done\n")
}
