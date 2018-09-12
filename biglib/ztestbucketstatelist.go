package biglib

//with SimulationState
// ensure JSON acceptance
//with BucketStateCache
// ensure JSON acceptance
//with (b *BucketStateCache) GetLastOperation()
// when empty ensure nil
// when len=1 ensure returns first
// when len>1 ensure returns last
// when len>Max ensure returns last
//with (s BucketStateCache) GetNextIndex()
// when empty ensure returns big.Int(0)
// ensure returns sum s.FillCount + s.PourCount + s.EmptyCount +1
//with (c *BucketStateCache) appendInitialBucket(fromb bool, controller *ChannelController)
// ensure inserted object has Idx == 0
// ensure inserted object has FromB == fromb
// ensure inserted object has BucketA == 0
// ensure inserted object has BucketB == 0
// ensure inserted object has Operation == biglib.Initial
// when empty ensure is inserted at idx 0 and returns true
// when partially populated ensure state is inserted in last position and returns true (should never happen)
// when full is not inserted into cache and returns false (should never happen)
// when empty it sends the state to the controller.simulationOperationCollector channel
// when full it sends the state to the controller.simulationOperationCollector channel
// when partially populated it sends the state to the controller.simulationOperationCollector channel
//with (c *BucketStateCache) appendNewBucket(index *big.Int, from *big.Int, to *big.Int, op bp.SimulationOperation, fromb bool, controller *ChannelController)
// ensure inserted object has Idx == index
// ensure inserted object has FromB == fromb
// when fromb ensure inserted object has BucketB == from
// when fromb ensure inserted object has BucketA == to
// when !fromb ensure inserted object has BucketA == from
// when !fromb ensure inserted object has BucketB == to
// ensure inserted object has Operation == op
// when empty ensure is inserted at idx 0 and returns true (should never happen)
// when partially populated ensure state is inserted in last position and returns true
// when full is not inserted into cache and returns false
// when empty it sends the state to the controller.simulationOperationCollector channel
// when full it sends the state to the controller.simulationOperationCollector channel
// when partially populated it sends the state to the controller.simulationOperationCollector channel
//with (c *BucketStateCache) appendSolvedBucket(index *big.Int, from *big.Int, to *big.Int, fromb bool, controller *ChannelController)
// ensure inserted object has Idx == index
// ensure inserted object has FromB == fromb
// when fromb ensure inserted object has BucketB == from
// when fromb ensure inserted object has BucketA == to
// when !fromb ensure inserted object has BucketA == from
// when !fromb ensure inserted object has BucketB == to
// ensure inserted object has Operation == biglib.FinalOp
// when empty ensure is inserted at idx 0 and returns true (should never happen)
// when partially populated ensure state is inserted in last position and returns true
// when full is not inserted into cache and returns false
// when empty it sends the state to the controller.simulationOperationCollector channel
// when full it sends the state to the controller.simulationOperationCollector channel
// when partially populated it sends the state to the controller.simulationOperationCollector channel
//with (c *BucketStateCache) appendErrorBucket(index *big.Int, code bp.ResultCode, controller *ChannelController)
// ensure inserted object has Idx == index
// ensure inserted object has FromB == fromb
// ensure inserted object has BucketB == 0
// ensure inserted object has BucketA == 0
// ensure inserted object has Operation == biglib.SimulationError
// when empty ensure is inserted at idx 0 and returns true (should never happen)
// when partially populated ensure state is inserted in last position and returns true
// when full is not inserted into cache and returns false
// when empty it sends the state to the controller.simulationOperationCollector channel
// when full it sends the state to the controller.simulationOperationCollector channel
// when partially populated it sends the state to the controller.simulationOperationCollector channel
//with newBucketState(idx *big.Int, from *big.Int, to *big.Int, operation bp.SimulationOperation, fromb bool)
// ensure return has Idx == idx
// when fromb ensure return has BucketB == from
// when fromb ensure return has BucketA == to
// when !fromb ensure return has BucketA == from
// when !fromb ensure return has BucketB == to
// ensure return has Operation == operation
//with newBucketStateCache() (b *BucketStateCache) {
//	ensure return has FillCount == big.NewInt(0)
//	ensure return has PourCount == big.NewInt(0)
//	ensure return has EmptyCount == big.NewInt(0)
//	ensure return has BucketStateList !=nil len() == 0
//	return b
//}
