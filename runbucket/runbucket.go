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
	pos        int
}

type Chans struct {
	controlChannel *chan bignum.ProcessControlOperation
	stateChannel   chan bignum.ProcessControlOperation
	resultsChannel chan bignum.SimulationState
	running        bool
}

func newChans() (c *Chans) {
	r := new(Chans)
	r.stateChannel = make(chan bignum.ProcessControlOperation, 2)
	r.resultsChannel = make(chan bignum.SimulationState, 20)
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
var table = tview.NewTable().SetBorders(true).SetFixed(1, 1)
var normStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
var pauseStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
var runStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
var doneStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
var selStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
var currselStyle = new(tcell.Style).Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

type (
	JobList []*Job
)

var lastselected = -1

func main() {
	joblist := make([]Job, 3)
	table.Select(0, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyTab {
				if lastselected > -1 && lastselected < len(joblist) {
					pjob := joblist[lastselected]
					*pjob.controller.GetStopStartChannel() <- bignum.Pause
				}
			}
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, false)
			}
		}).
		SetSelectedFunc(
			func(row int, column int) {
				table.GetCell(row, 0).SetStyle(currselStyle)
				table.SetSelectable(false, false)
				if lastselected > -1 && lastselected < len(joblist) && row != lastselected {
					oldjob := joblist[lastselected]
					style := normStyle
					if &oldjob != nil {
						if oldjob.controller.IsTerminated() {
							style = doneStyle
						} else if oldjob.controller.IsPaused() {
							style = pauseStyle
						} else if oldjob.controller.IsRunning() {
							style = runStyle
						}
					}
					table.GetCell(lastselected, 0).SetStyle(style)
					table.GetCell(lastselected, 1).SetStyle(style)
				}
				lastselected = row
			}).
		SetSelectedStyle(selStyle.Decompose())

	cbucketa, _ := new(big.Int).SetString("57", 10)
	cbucketb, _ := new(big.Int).SetString("41", 10)
	cdesired, _ := new(big.Int).SetString("20", 10)

	bbucketa, _ := new(big.Int).SetString("100000013", 10)
	bbucketb, _ := new(big.Int).SetString("10000013", 10)
	bdesired, _ := new(big.Int).SetString("100003", 10)

	bucketa, _ := new(big.Int).SetString("1000013", 10)
	bucketb, _ := new(big.Int).SetString("10013", 10)
	desired, _ := new(big.Int).SetString("10003", 10)

	problem := bignum.NewProblem(bucketa, bucketb, desired)
	bproblem := bignum.NewProblem(bbucketa, bbucketb, bdesired)
	cproblem := bignum.NewProblem(cbucketa, cbucketb, cdesired)

	joba := NewJob(table, problem)
	jobb := NewJob(table, bproblem)
	jobc := NewJob(table, cproblem)

	joblist = []Job{*joba, *jobb, *jobc}

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func setSelected(sel int) {
	//r, c := table.GetSelection()
}

func addJob(j *Job) {
	newidx := table.GetRowCount()

	j.pos = newidx
	table.SetCell(newidx, 0,
		j.statecell.
			SetStyle(normStyle).
			SetAlign(tview.AlignCenter))

	table.SetCell(newidx, 1,
		j.simcell.
			SetStyle(normStyle).
			SetAlign(tview.AlignLeft).SetExpansion(40))
}

func NewJob(table *tview.Table, problem *bignum.Problem) (j *Job) {
	j = new(Job)
	j.statecell = tview.NewTableCell("hi")
	j.simcell = tview.NewTableCell("there")

	j.statecell.SetText("bye")
	j.simcell.SetText("here")

	solution := bignum.NewSolution(problem)

	j.chans = newChans()

	j.controller = bignum.GetIdleSolutionProcessor("DUMMY", solution, &j.chans.stateChannel, &j.chans.resultsChannel)
	j.startListenChannels()
	addJob(j)
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
	if job.pos != lastselected {
		job.simcell.SetStyle(doneStyle)
	}
	job.statecell.SetText(printState(job.controller))
	app.Draw()
}

func (j *Job) doListenProgress(display *tview.TableCell, group *sync.WaitGroup) {
	var running = true
	defer func() {
		if group != nil {
			display.SetStyle(doneStyle)
			group.Done()
		}
	}()
	for running {
		select {
		case op := <-j.chans.stateChannel:
			display.SetText(printState(j.controller))
			if j.pos == lastselected {
				display.SetStyle(currselStyle)
			} else if j.controller.IsTerminated() {
				display.SetStyle(doneStyle)
			} else if j.controller.IsPaused() {
				display.SetStyle(pauseStyle)
			} else {
				display.SetStyle(runStyle)
			}

			if int(op)&int(bignum.Error) > 0 {
				if j.pos == lastselected {
					display.SetStyle(currselStyle)
				} else {
					display.SetStyle(doneStyle)
				}
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
