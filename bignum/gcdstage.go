package bignum

import (
	"fmt"
	"math/big"
)
import bp "arrdude.com/bucketproblem"

func (s *Solution) generateGCD(controller *ChannelController) {
	controller.state = controller.state | StageGcd | Running | Initialized

	if !controller.mayContinue() {
		s.Operations.appendErrorBucket(bp.ProcessKilled, controller)
		return
	}

	if s.Problem.BucketA.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bp.BucketATooSmall, controller)
		return
	}

	if s.Problem.BucketB.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bp.BucketBTooSmall, controller)
		return
	}

	if s.Problem.Desired.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bp.DesiredTooSmall, controller)
		return
	}

	if s.Problem.Desired.Cmp(s.Problem.BucketA) != -1 {
		s.Operations.appendErrorBucket(bp.DesiredTooBig, controller)
		return
	}

	if s.Problem.BucketB.Cmp(s.Problem.BucketA) != -1 {
		s.Operations.appendErrorBucket(bp.BucketBTooBig, controller)
		return
	}

	s.Denominator.GCD(s.MultInverseA, s.MultInverseB, s.Problem.BucketA, s.Problem.BucketB)

	//Just in case - allowed parameters should not allow this
	if s.Denominator.Cmp(bigzero) == -1 {
		s.Operations.appendErrorBucket(bp.NoGCDFound, controller)
		return
	}

	moddesired := new(big.Int).Mod(s.Problem.Desired, s.Denominator)

	if moddesired.Cmp(bigzero) != 0 {
		fmt.Printf("DenominatorNotMultiple!!! (A:%v B:%v) %v %% %v = %v\n", s.Problem.BucketA, s.Problem.BucketB, s.Problem.Desired, s.Denominator, moddesired)
		s.Operations.appendErrorBucket(bp.DenominatorNotMultiple, controller)
		return
	}

	s.TvolumeB.Mul(s.Problem.Desired, s.MultInverseB)
	s.CountFromB.Mod(s.TvolumeB, s.Problem.BucketA)

	//s.TvolumeA.Mul(s.Problem.Desired, s.MultInverseA)
	//s.CountFromA.Mod(s.TvolumeA, s.Problem.BucketA)
	//s.CountFromA.Div(s.CountFromA, s.Denominator)
	//s.CountFromB.Div(s.CountFromB, s.Denominator)

	//s.TvolumeA.Mul(s.Problem.Desired, s.MultInverseA)
	//s.CountFromA.Mod(s.TvolumeA, s.Problem.BucketB)

	s.CountFromA.Mul(s.CountFromB, bignegone)
	s.CountFromA.Mod(s.CountFromA, s.Problem.BucketA)

	//s.FromB = s.compareCountFromAandCountFromB() == 1
	identity := s.compareCountFromAandCountFromB()
	if identity == 0 {
		s.FromB = s.Problem.Desired.Cmp(s.Problem.BucketB) == 0
	} else {
		s.FromB = identity == 1
	}

	s.PredictedStateCount = new(big.Int)

	//from := new(big.Int)
	//to := new(big.Int)

	capfrom := new(big.Int)
	capto := new(big.Int)

	countempties := new(big.Int)

	//switch buckets for additive direction and set direction specific pretenses
	if s.FromB {
		capfrom.Set(s.Problem.BucketB)
		capto.Set(s.Problem.BucketA)

		countempties.Mul(s.CountFromB, capfrom)
		countempties.Div(countempties, capto)
		//subtract a predicted empty if we will not overflow
		if s.Problem.Desired.Cmp(capfrom) != 1 {
			countempties.Sub(countempties, bigone)
		}
		s.PredictedStateCount.Set(s.CountFromB)
	} else {
		capfrom.Set(s.Problem.BucketA)
		capto.Set(s.Problem.BucketB)

		countempties.Div(capfrom, capto)
		countempties.Mul(countempties, s.CountFromA)

		if s.Problem.Desired.Cmp(s.Problem.BucketA) == 1 {
			countempties.Sub(countempties, bigone)
		}
		s.PredictedStateCount.Set(s.CountFromA)
	}

	s.PredictedStateCount.Add(s.PredictedStateCount, countempties)
	//multiply by steps per action
	s.PredictedStateCount.Mul(s.PredictedStateCount, stepsPerAction)
	//add initial state
	s.PredictedStateCount.Add(s.PredictedStateCount, big.NewInt(1))

	controller.stateCollector <- (StageDone)
}
