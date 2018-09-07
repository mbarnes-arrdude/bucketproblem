package bucketproblem

import (
	"fmt"
	"os"
)

const (
	// build managed values
	PROCESS_NAME  = "BUCKALG"
	VERSION_MAJOR = 0
	VERSION_MINOR = 0
	VERSION_PATCH = 1
)

//used for identity
func GetProcessHash() int {
	return os.Getpid()
}

var MaxOperationsListSize = 10023

func GetVersionId() string {
	return fmt.Sprintf("%2d%2d%2d%s", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH, PROCESS_NAME)
}
