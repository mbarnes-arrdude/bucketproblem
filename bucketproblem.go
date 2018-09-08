package bucketproblem

import (
	"fmt"
)

const (
	// build managed values
	PROCESS_NAME  = "BUCKALG"
	VERSION_MAJOR = 0
	VERSION_MINOR = 0
	VERSION_PATCH = 1
)

var MaxOperationsListSize = 100

func GetVersionId() string {
	return fmt.Sprintf("%2d%2d%2d%s", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH, PROCESS_NAME)
}
