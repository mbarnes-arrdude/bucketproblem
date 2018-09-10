package bignum

import (
	"math/big"
)
import bp "arrdude.com/bucketproblem"

func (s *Solution) generateGCD(controller *ChannelController) {
	controller.state = controller.state | StageGcd | Running | Initialized

	if !controller.mayContinue() {
		s.Operations.appendErrorBucket(bigzero, bp.ProcessKilled, controller)
		return
	}

	if s.Problem.BucketA.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bigzero, bp.BucketATooSmall, controller)
		return
	}

	if s.Problem.BucketB.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bigzero, bp.BucketBTooSmall, controller)
		return
	}

	if s.Problem.Desired.Cmp(bigzero) != 1 {
		s.Operations.appendErrorBucket(bigzero, bp.DesiredTooSmall, controller)
		return
	}

	if s.Problem.Desired.Cmp(s.Problem.BucketA) > -1 {
		s.Operations.appendErrorBucket(bigzero, bp.DesiredTooBig, controller)
		return
	}

	if s.Problem.BucketB.Cmp(s.Problem.BucketA) > -1 {
		s.Operations.appendErrorBucket(bigzero, bp.BucketBTooBig, controller)
		return
	}

	s.Complexity.Add(s.Problem.BucketA, s.Problem.BucketB)

	s.Denominator.GCD(s.MultInverseA, s.MultInverseB, s.Problem.BucketA, s.Problem.BucketB)

	if s.Denominator.Cmp(bigzero) == -1 {
		s.Operations.appendErrorBucket(bigzero, bp.NoGCDFound, controller)
		return
	}

	s.TvolumeB.Mul(s.Problem.Desired, s.MultInverseB)
	s.CountFromB.Mod(s.TvolumeB, s.Problem.BucketA)

	moddesired := new(big.Int).Mod(s.Problem.Desired, s.Denominator)

	if moddesired.Cmp(bigzero) != 0 {
		s.Operations.appendErrorBucket(bigzero, bp.DenominatorNotMultiple, controller)
		return
	}

	s.CountFromA.Mul(s.CountFromB, bignegone)
	s.CountFromA.Mod(s.CountFromA, s.Problem.BucketA)

	identity := s.compareCountFromAandCountFromB()
	if identity == 0 {
		s.FromB = s.Problem.Desired.Cmp(s.Problem.BucketB) == 0
	} else {
		s.FromB = identity == 1
	}

	s.PredictedStateCount = new(big.Int)

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

		countempties.Mul(s.CountFromA, capto)
		countempties.Div(countempties, capfrom)

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
