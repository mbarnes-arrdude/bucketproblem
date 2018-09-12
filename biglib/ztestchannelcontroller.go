package biglib

//with (op *ProcessControlOperation) String()
// ensure acceptance test using bracket of output for all states 0x0 <= int(op) <= 0xffff
//with NewChannelController(solution *Solution, autorun bool)
// ensure return not nil
// ensure return has Solution == solution
// ensure return has autorun == autorun
// ensure return has StartStopCollector != nil
// ensure return has stateCollector != nil
// ensure return has stateEmitters != nil
// ensure return has stateEmittersMutex != nil
// ensure return has simulationOperationCollector != nil
// ensure return has simulationOperationEmitters != nil
// ensure return has simulationOperationEmittersMutex != nil
// when autorun ensure return blocks and return solution is complete
// when !autorun ensure return is immediate and return is started
//with (s *ChannelController) GetStopStartChannel()
// ensure return is not nil
//with (p *ChannelController) MayContinue()
// when running and not terminated and not paused it returns true
// when not running returns false
// when running and pause it blocks
// when running and not terminated returns false
// when terminated returns false
// when running and pause it blocks when unpaused it returns true
// when running then paused (it blocks) then when terminated it returns false
//with (p *ChannelController) IsAutorun() when initialized
// ensure return != nil
//with (p *ChannelController) GetState() when initialized
// ensure return == NoOp
//with (p *ChannelController) IsTerminated() when initialized
// ensure return false
//with (p *ChannelController) GetStage() when initialized
// ensure return == NoOp
//with (p *ChannelController) IsRunning() when initialized
// ensure return false when not running
//with (p *ChannelController) IsStageInitialized() when initialized
// ensure return false
//with (p *ChannelController) IsPaused() when initialized
// ensure return false
//with (p *ChannelController) IsStageComplete() when initialized
// ensure return false
//with (s *ChannelController) RegisterStateChannel(name string, statech *chan ProcessControlOperation)
// ensure channel is added to map as name
//with (s *ChannelController) UnregisterStateChannel(name string)
// ensure channel is removed from map
//with (s *ChannelController) RegisterResultChannel(name string, resultch *chan SimulationState)
// ensure channel is added to map as name
//with (s *ChannelController) UnregisterResultChannel(name string)
// ensure channel is removed from map
//with (p *ChannelController) changeState(op ProcessControlOperation)
// when op == Kill no signals are sent
// when op == Kill state is adjusted
// when op == StartSimulationOnly no signals are sent
// when op == StartSimulationOnly state is adjusted
// when op == StartSimulationOnly and terminated state is not changed, simulation is not started
// when op == StartSimulationOnly simulation is started
// when op == StartGcdOnly state is adjusted
// when op == StartGcdOnly no signals are sent
// when op == StartGcdOnly and terminated state is not changed, simulation is not started
// when op == StartGcdOnly simulation is not started
// when op == Done state is adjusted
// when op == Done no signals are sent
// when op == StageDone state is adjusted
// when op == StageDone and !gcd and autorun Done is sent
// when op == StageDone and gcd and autorun simulation is started
// when op == StageDone and gcd and not autorun simulation is Not started
// when op == Pause state is flipped and initialied and running ensured
// when op == Pause no signals are sent
// when op == Start and gcd and autorun simulation is started
// when op == Start and !gcd or !autorun and pause simulation is unpause and state adjusts
// when op == Start and !isRunning and stage is Noop autorun becomes true and gcd is started
//with (p *ChannelController) startListeners(autostart bool)
// when autostart return blocks, start signal is sent
// when !autostart return does not block no signal is sent
// when returned channels are closed
//with (s *ChannelController) listenStateChanges(group sync.WaitGroup, autostart bool)
// function blocks forever
// when autostart waitgroup.Done() is called
// when <-s.stateCollector read does not block
// when <-s.stateCollector registered channels are notified with blocking
// when <-s.stateCollector == Done function returns
// when <-s.StartStopCollector then s.stateCollector <-
//with (s *ChannelController) listenSimulationEvents(group sync.WaitGroup, autostart bool)
// function blocks forever
// when autostart waitgroup.Done is called
// read <-s.simulationOperationCollector does not block
// registered channel <- are not blocking
// when <- last object (FinalOp and Error) write to subscribed channel blocking
// when <- FinalOp s.stateCollector <- Done
// when <- Error s.stateCollector <- Kill
