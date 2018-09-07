package bignum

import "fmt"
import bp "arrdude.com/bucketproblem"

//GetRunningSolutionProcessor() creates a Solution running the algorithm and then returns the pointer to the running solution.
func GetRunningSolutionProcessor(s *Solution) *ChannelController {
	return NewChannelController(s, true)
}

//GetRunningSolutionProcessor() creates a Solution which waits in idle state and returns the pointer.
func GetIdleSolutionProcessor(poolid string, c *ChannelController, stateChannel *chan ProcessControlOperation, resultChannel *chan SimulationState) *ChannelController {
	var hash = c.solution.Problem.Hash()
	var name = fmt.Sprintf("%s%019d", bp.GetVersionId(), hash)
	if stateChannel != nil {
		c.RegisterStateChannel(name, stateChannel)
	}
	fmt.Printf("Registering result channel %v\n", resultChannel)
	if resultChannel != nil {
		c.RegisterResultChannel(name, resultChannel)
	}
	return c
}
