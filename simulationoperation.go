package bucketproblem

const (
	InitialOp       SimulationOperation = 0
	Empty           SimulationOperation = 1
	Fill            SimulationOperation = 2
	Pour            SimulationOperation = 3
	FinalOp         SimulationOperation = 4
	SimulationError SimulationOperation = 5
)

type SimulationOperation int

var SimulationOperations = [...]string{
	"New Buckets",
	"Empty",
	"Fill",
	"Pour",
	"Solved",
	"Error",
}

func (o SimulationOperation) String() string {
	if o < InitialOp || o > SimulationError {
		return "Unknown"
	}
	return SimulationOperations[o]
}
