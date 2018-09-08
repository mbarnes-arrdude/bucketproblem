package bucketproblem

const (
	InitialOp       SimulationOperation = 0
	Empty           SimulationOperation = 1
	Fill            SimulationOperation = 2
	Pour            SimulationOperation = 3
	FinalOp         SimulationOperation = 4
	Truncated       SimulationOperation = 5
	SimulationError SimulationOperation = 6
)

type SimulationOperation int

var SimulationOperations = [...]string{
	"New Buckets",
	"Empty",
	"Fill",
	"Pour",
	"Solved",
	"Truncated",
	"X Error",
}

var SimulationOperationChars = [...]rune{
	'_',
	'v',
	'^',
	'+',
	'=',
	'x',
	'!',
}

func (o SimulationOperation) String() string {
	if o < InitialOp || o > SimulationError {
		return "Unknown"
	}
	return SimulationOperations[o]
}

func (o SimulationOperation) Rune() rune {
	if o < InitialOp || o > SimulationError {
		return '.'
	}
	return SimulationOperationChars[o]
}
