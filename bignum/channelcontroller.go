package bignum

import (
	bp "arrdude.com/bucketproblem"
	"fmt"
	"sync"
	"time"
)

const (
	pauserelaxedperiod              = 250
	SimulationCollectorChannelSize  = 2
	StateCollectorChannelBufferSize = 2
	PublicCommandChannelBufferSize  = 2
	//ProcessControlOperation Masks
	Noop            int = 0x0
	Initialized     int = 0x1
	Paused          int = 0x3
	Complete        int = 0x7
	Running         int = 0x10
	StageStats      int = 0xf00
	StageGcd        int = 0x100
	StageSimulation int = 0x200
	ProcessStats    int = 0xf000
	ProcessFinished int = 0x1000
	Error           int = 0x3000

	Start               ProcessControlOperation = ProcessControlOperation(Running)
	Pause               ProcessControlOperation = ProcessControlOperation(Paused)
	StageDone           ProcessControlOperation = ProcessControlOperation(Complete)
	Done                ProcessControlOperation = ProcessControlOperation(ProcessFinished)
	StartGcdOnly        ProcessControlOperation = ProcessControlOperation(StageGcd)
	StartSimulationOnly ProcessControlOperation = ProcessControlOperation(StageSimulation)
	Kill                ProcessControlOperation = ProcessControlOperation(Error)
)

type ProcessControlOperation int

type ChannelController struct {
	solution *Solution
	state    int
	autorun  bool

	StartStopCollector chan ProcessControlOperation
	stateCollector     chan ProcessControlOperation
	stateEmitters      map[string]*chan ProcessControlOperation
	stateEmittersMutex *sync.Mutex

	simulationOperationCollector     chan SimulationState
	simulationOperationEmitters      map[string]*chan SimulationState
	simulationOperationEmittersMutex *sync.Mutex
}

func NewChannelController(solution *Solution, autorun bool) (state *ChannelController) {
	state = new(ChannelController)
	state.solution = solution
	state.autorun = autorun
	state.state = int(Noop)

	state.StartStopCollector = make(chan ProcessControlOperation, PublicCommandChannelBufferSize)
	state.stateCollector = make(chan ProcessControlOperation, StateCollectorChannelBufferSize)
	state.stateEmitters = make(map[string]*chan ProcessControlOperation)
	state.stateEmittersMutex = new(sync.Mutex)

	state.simulationOperationCollector = make(chan SimulationState, SimulationCollectorChannelSize)
	state.simulationOperationEmitters = make(map[string]*chan SimulationState)
	state.simulationOperationEmittersMutex = new(sync.Mutex)

	if autorun {
		state.startListeners(true)
	} else {
		go state.startListeners(false)
	}
	return state
}

func (p *ChannelController) startListeners(autostart bool) {
	var wg sync.WaitGroup
	if !autostart {
		wg.Add(2)
	}
	go p.listenStateChanges(wg, autostart)
	if autostart {
		p.StartStopCollector <- Start
		p.listenSimulationEvents(wg, autostart)
	} else {
		go p.listenSimulationEvents(wg, autostart)
	}
	defer func() {
		if !autostart {
			wg.Wait()
		}
		close(p.simulationOperationCollector)
		close(p.StartStopCollector)
		close(p.stateCollector)
	}()
}

func (p *ChannelController) mayContinue() bool {
	if !p.IsRunning() || p.IsTerminated() {
		fmt.Println("not running or is terminated")
		return false
	}
	for p.IsPaused() {
		fmt.Println("PAUSED")
		time.Sleep(time.Duration(pauserelaxedperiod) * time.Millisecond)
	}
	if !p.IsRunning() || p.IsTerminated() {
		fmt.Println("post pause not running or is terminated")
		return false
	}
	return true
}

func (p *ChannelController) changeState(op ProcessControlOperation) {

	switch op {
	case Kill:
		p.state = p.state & ^ProcessStats & ^StageStats & ^Running & ^Complete
		p.state = p.state | Error
		return
		break
	case StartSimulationOnly:
		if p.IsTerminated() {
			return
		}

		p.state = p.state & ^StageStats & ^Complete
		p.state = p.state | StageSimulation | Running | Initialized
		defer func() {
			go p.solution.generateSimulation(p)
		}()

		return
		break
	case StartGcdOnly:
		if p.IsTerminated() {
			return
		}
		p.state = p.state & ^StageStats & ^Complete
		p.state = p.state | StageGcd | Running | Initialized
		defer func() {
			go func() {
				p.solution.generateGCD(p)
			}()
		}()
		return
		break
	case Done:
		p.state = p.state & ^Running
		p.state = p.state | ProcessFinished
		return
		break
	case StageDone:
		p.state = int(p.state) | Complete
		if p.GetStage() == StageGcd && p.IsAutorun() {
			p.stateCollector <- StartSimulationOnly
			return
			break
		}
		p.stateCollector <- Done
		return
		break
	case Pause:
		p.state = p.state | Paused
		return
		break
	case Start:
		if p.IsRunning() && p.IsStageComplete() {
			if p.GetStage() == StageGcd && p.IsAutorun() {
				p.autorun = true
				p.state = p.state | Complete
				p.stateCollector <- StartSimulationOnly
				return
			}
		} else if p.IsPaused() {
			p.state = int(p.state) | Running & ^Paused
		} else if !p.IsRunning() {
			if p.GetStage() == Noop {
				p.autorun = true
				p.stateCollector <- StartGcdOnly
				return
			}
		}
		break
	}
}

func (s *ChannelController) GetStopStartChannel() *chan ProcessControlOperation {
	return &s.StartStopCollector
}

func (p *ChannelController) IsAutorun() bool {
	return p.autorun
}

func (p *ChannelController) IsTerminated() bool {
	return int(p.state)&ProcessFinished == ProcessFinished
}

func (p *ChannelController) GetStage() int {
	return int(p.state) & StageStats
}

func (p *ChannelController) IsRunning() bool {
	return int(p.state)&Running == Running
}

func (p *ChannelController) IsStageInitialized() bool {
	return int(p.state)&Initialized == Initialized
}

func (p *ChannelController) IsPaused() bool {
	return int(p.state)&Paused == Paused
}

func (p *ChannelController) IsStageComplete() bool {
	return int(p.state)&Complete == Complete
}

//RegisterStateChannel will register a *chan ProcessControlOperation identified by name to receive
// ProcessControlOperation event notifications. Notifications will be added to the channel as the process control
// of the Solution changes by internal or external trigger. Registering a channel with a name collision for
// an existing channel will clobber.
//
// Arguments:
// name string
// statech *chan ProcessControlOperation
//
func (s *ChannelController) RegisterStateChannel(name string, statech *chan ProcessControlOperation) {
	s.stateEmittersMutex.Lock()
	s.stateEmitters[name] = statech
	s.stateEmittersMutex.Unlock()
}

//UnregisterStateChannel will unregister a *chan ProcessControlOperation identified by name. Misses fail silently.
//
// Arguments:
// name string
func (s *ChannelController) UnregisterStateChannel(name string) {
	s.stateEmittersMutex.Lock()
	delete(s.stateEmitters, name)
	s.stateEmittersMutex.Unlock()
}

//RegisterResultChannel will register a *chan SimulationState identified by name to receive state SimulationState event
// notifications. Notifications will be added to the channel as the simulator adds states.
//
// Arguments:
// name string
// resultch *chan SimulationState
func (s *ChannelController) RegisterResultChannel(name string, resultch *chan SimulationState) {
	if resultch == nil {
		return
	}
	s.simulationOperationEmittersMutex.Lock()
	s.simulationOperationEmitters[name] = resultch
	s.simulationOperationEmittersMutex.Unlock()
}

//UnregisterResultChannel will unregister a *chan ProcessControlOperation identified by name. Misses fail silently.
//
// Arguments:
// name string
func (s *ChannelController) UnregisterResultChannel(name string) {
	s.simulationOperationEmittersMutex.Lock()
	delete(s.simulationOperationEmitters, name)
	s.simulationOperationEmittersMutex.Unlock()
}

func (s *ChannelController) listenStateChanges(group sync.WaitGroup, autostart bool) {
	var running = true
	defer func() {
		if !autostart {
			group.Done()
		}
	}()
	for running {
		select {
		case o := <-s.stateCollector:
			s.changeState(o)
			for _, ch := range s.stateEmitters {
				//BLOCKING
				select {
				case *ch <- o:
					break
				}
			}
			if o&Done == Done {
				running = false
				return
			}
		case o := <-s.StartStopCollector:
			s.stateCollector <- o
		default:
		}
	}
}

func (s *ChannelController) listenSimulationEvents(group sync.WaitGroup, autostart bool) {
	var running = true
	defer func() {
		if !autostart {
			group.Done()
		}
	}()
	for running {
		select {
		case o := <-s.simulationOperationCollector:
			for _, ch := range s.simulationOperationEmitters {
				//BLOCKING
				select {
				case *ch <- o:
					break
				}
			}
			if int(o.Operation) >= int(bp.FinalOp) {
				fmt.Println("Finalized")
				running = false
				defer func() {
					if o.Operation == bp.FinalOp {
						s.stateCollector <- Done
					} else {
						s.stateCollector <- Kill
					}
				}()
				return
			}
			break
		default:
		}
	}
}

func (p *ChannelController) String() string {
	return fmt.Sprintf("Controller: %d Stage: %d Running %d Status %d", p.state&ProcessStats, p.state&StageStats, p.state&Running, p.state&Complete)
}
