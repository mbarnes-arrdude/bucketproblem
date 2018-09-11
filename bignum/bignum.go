package bignum

import "fmt"
import bp "arrdude.com/bucketproblem"

//GetRunningSolutionProcessor() creates a Solution running the algorithm and then returns the pointer to the running solution.
func GetRunningSolutionProcessor(s *Solution) *ChannelController {
	return NewChannelController(s, true)
}

//GetRunningSolutionProcessor() creates a Solution and ChannelController which waits in idle state for a signal to start.
//
//Parameters:
//poolid string NYI
//s *Solution is a pointer to the Solution object
//stateChannel *chan ProcessControlOperation is the caller's channel will be subscribed for reading controller events.
//  Writes to this channel are blocking and the controller will halt until they are read.
//resultChannel *chan SimulationState is the caller's channel which will be subscribed for simulator events. Writes to
// this channel are not blocking except for the last simulator event which contains the final state which will block
// until read.
//
//Returns:
//c *ChannelController which will have the Solution ready for solving and/or simulation. To direct the controller, the
// caller should request the Start/Stop channel from the controller using GetStopStartChannel() which will accept and
// a ProcessControlOperation eg. c.GetStartStopChannel() <- bignum.Start
func GetIdleSolutionProcessor(poolid string, s *Solution, stateChannel *chan ProcessControlOperation, resultChannel *chan SimulationState) (c *ChannelController) {
	c = NewChannelController(s, false)
	var hash = c.Solution.Problem.Hash()
	var name = fmt.Sprintf("%s%019d", bp.GetVersionId(), hash)
	if stateChannel != nil {
		c.RegisterStateChannel(name, stateChannel)
	}
	if resultChannel != nil {
		c.RegisterResultChannel(name, resultChannel)
	}
	return c
}
