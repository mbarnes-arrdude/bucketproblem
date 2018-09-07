package bucketproblem

const (
	ResultsOK              ResultCode = 0
	ResultsTruncated       ResultCode = 1
	BucketSizeMin          ResultCode = 2
	DenominatorNotMultiple ResultCode = 3
	BucketATooSmall        ResultCode = 4
	BucketBTooSmall        ResultCode = 5
	DesiredTooSmall        ResultCode = 6
	BucketBTooBig          ResultCode = 7
	DesiredTooBig          ResultCode = 8
	NoGCDFound             ResultCode = 9
	ProcessKilled          ResultCode = 10
)

//ResultCode is a pseudo-enum having descriptions of ResultCodes
type ResultCode int

var ResultCodes = [...]string{
	"The problem was solved",
	"The problem was solved but the list of operations truncated.",
	"The problem could not be solved because the problem requires two buckets of positive size and at least one was less than 1 unit bignum.",
	"The problem could not be solved because the desired amount is not a factor of the greatest common denominator of the buckets.",
	"The problem could not be solved because the size of BucketA was too small. It must be in the range: 0 < BucketA.",
	"The problem could not be solved because the size of BucketB was too small. It must be in the range: 0 < BucketB < BucketA.",
	"The problem could not be solved because the desired amount was too small. It must be in the range: 0 < Desired < BucketA.",
	"The problem could not be solved because the size of BucketB too large. It must be in the range: 0 < BucketB < BucketA.",
	"The problem could not be solved because the desired amount too large. It must be in the range: 0 < Desired < BucketA.",
	"The problem could not be solved because there was no Greatest Common Denominator of the Inputs. A and B must be relatively prime.",
	"The problem could not be solved because there was an unexpected error or the process was otherwise killed.",
}

func (code ResultCode) String() string {
	if code < ResultsOK || code > ProcessKilled {
		return "Unknown Solution"
	}
	return ResultCodes[code]
}
