package main

import (
	bp "arrdude.com/bucketproblem"
	"arrdude.com/bucketproblem/bignum"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"math/big"
	"sync"
	"time"
)

type Job struct {
	problem    bignum.Problem
	startts    int64
	endts      int64
	controller *bignum.ChannelController
	statecell  *tview.TableCell
	simcell    *tview.TableCell
	chans      *Chans
	wg         sync.WaitGroup
}

type Chans struct {
	controlChannel *chan bignum.ProcessControlOperation
	stateChannel   chan bignum.ProcessControlOperation
	resultsChannel chan bignum.SimulationState
	running        bool
}

func newchans() (c *Chans) {
	r := new(Chans)
	r.stateChannel = make(chan bignum.ProcessControlOperation, 2)
	r.resultsChannel = make(chan bignum.SimulationState, 200)
	return r
}

func (c *Chans) updateResultsChannelBuffer(newsize int) {
	if newsize < bignum.SimulationCollectorChannelSizeSmall {
		newsize = bignum.SimulationCollectorChannelSizeSmall
	} else if newsize > bignum.SimulationCollectorChannelSizeLarge {
		newsize = bignum.SimulationCollectorChannelSizeLarge
	}
	c.resultsChannel = make(chan bignum.SimulationState, newsize)
}

var app = tview.NewApplication()

type (
	JobList []*Job
)

func main() {
	joblist := make(JobList, 1)
	table := tview.NewTable().SetBorders(true).SetFixed(1, 2)

	bucketa, _ := new(big.Int).SetString("57", 10)
	bucketb, _ := new(big.Int).SetString("41", 10)
	desired, _ := new(big.Int).SetString("20", 10)
	problem := bignum.NewProblem(bucketa, bucketb, desired)

	bbucketa, _ := new(big.Int).SetString("100000013", 10)
	bbucketb, _ := new(big.Int).SetString("10000013", 10)
	bdesired, _ := new(big.Int).SetString("100003", 10)
	bproblem := bignum.NewProblem(bbucketa, bbucketb, bdesired)

	cbucketa, _ := new(big.Int).SetString("1000013", 10)
	cbucketb, _ := new(big.Int).SetString("10013", 10)
	cdesired, _ := new(big.Int).SetString("10003", 10)
	cproblem := bignum.NewProblem(cbucketa, cbucketb, cdesired)

	joba := NewJob(table, problem)
	jobb := NewJob(table, bproblem)
	jobc := NewJob(table, cproblem)

	joblist = append(joblist, []*Job{joba, jobb, jobc}...)

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func NewJob(table *tview.Table, problem *bignum.Problem) (j *Job) {
	j = new(Job)
	j.statecell = tview.NewTableCell("hi")
	j.simcell = tview.NewTableCell("there")

	color := tcell.ColorYellow

	newidx := table.GetRowCount()

	table.SetCell(newidx, 0,
		j.statecell.
			SetTextColor(color).
			SetAlign(tview.AlignCenter))

	table.SetCell(newidx, 1,
		j.simcell.
			SetTextColor(color).
			SetAlign(tview.AlignLeft).SetExpansion(40))

	j.statecell.SetText("bye")
	j.simcell.SetText("here")

	solution := bignum.NewSolution(problem)

	j.chans = newchans()

	j.controller = bignum.GetIdleSolutionProcessor("DUMMY", solution, &j.chans.stateChannel, &j.chans.resultsChannel)
	j.startListenChannels()
	defer func() {
		controlChannel := j.controller.GetStopStartChannel()
		*controlChannel <- bignum.Start
	}()

	return j
}

func (job *Job) startListenChannels() {
	job.wg.Add(2)
	go job.doListenSimulation(job.simcell, &job.wg)
	//job.simcell.SetText("Started")
	go job.doListenProgress(job.statecell, &job.wg)
	//job.statecell.SetText("Started")
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			*job.controller.GetStopStartChannel() <- bignum.Pause
		} else if event.Key() == tcell.KeyTab {
			*job.controller.GetStopStartChannel() <- bignum.Done
		}
		return event
	})
	defer func() {
		go func() {
			job.wg.Wait()
			job.showFinal()
		}()
	}()
	return
}

func (job *Job) showFinal() {
	job.simcell.SetText(printLastTableEntry(*job.controller.Solution.Operations.GetLastOperation()))
	job.statecell.SetText(printState(job.controller))
	app.Draw()
}

func (j *Job) doListenProgress(display *tview.TableCell, group *sync.WaitGroup) {
	var running = true
	defer func() {
		if group != nil {
			group.Done()
		}
	}()
	for running {
		select {
		case op := <-j.chans.stateChannel:
			display.SetText(printState(j.controller))
			if int(op)&int(bignum.Error) > 0 {
				running = false
				return
			}
		default:
		}
	}
}

func (j *Job) doListenSimulation(display *tview.TableCell, group *sync.WaitGroup) {
	var doit = true
	defer func() {
		if group != nil {
			group.Done()
		}
	}()
	for doit {
		select {
		case bucket := <-j.chans.resultsChannel:
			display.SetText(printLastTableEntry(bucket))
			app.Draw()
			time.Sleep(time.Millisecond * time.Duration(1))
			if bucket.Operation >= bp.FinalOp {
				doit = false
				break
			}
		default:
		}
	}
}

func printLastTableEntry(lastbucket bignum.SimulationState) string {
	return fmt.Sprintf("%15v) | %16v | %c | %v\n", lastbucket.Idx, lastbucket.AmountBucketA, lastbucket.Operation.Rune(), lastbucket.AmountBucketB)
}

func printState(c *bignum.ChannelController) string {
	pausepart := "Sim"
	if c.IsTerminated() {
		pausepart = "Done"
	} else if c.IsPaused() {
		pausepart = "PAUSED"
	} else {
		stage := c.GetStage()
		if stage == bignum.StageGcd {
			pausepart = "GCD"
		} else if stage == bignum.Noop {
			pausepart = "Ready"
		}
	}
	direction := "A->B (+)"
	if c.Solution.FromB {
		direction = "B->A (-)"
	}
	problempart := fmt.Sprintf("N: %12v", c.Solution.Problem.Desired)
	solutionpart := fmt.Sprintf("C:%2x D:%8s P: %12v", c.Solution.GetComplexityScale(), direction, c.Solution.PredictedStateCount)
	return fmt.Sprintf("%s %s %s", problempart, solutionpart, pausepart)
}
