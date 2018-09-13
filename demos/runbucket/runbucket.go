package main

import (
	bp "arrdude.com/bucketproblem"
	"arrdude.com/bucketproblem/biglib"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"math/big"
	"sync"
	"time"
)

type Job struct {
	problem    biglib.Problem
	startts    int64
	endts      int64
	controller *biglib.ChannelController
	statecell  *tview.TableCell
	simcell    *tview.TableCell
	chans      *Chans
	wg         sync.WaitGroup
	pos        int
}

type Chans struct {
	controlChannel *chan biglib.ProcessControlOperation
	stateChannel   chan biglib.ProcessControlOperation
	resultsChannel chan biglib.SimulationState
	running        bool
}

func newChans() (c *Chans) {
	r := new(Chans)
	r.stateChannel = make(chan biglib.ProcessControlOperation, 2)
	r.resultsChannel = make(chan biglib.SimulationState, 20)
	return r
}

var app = tview.NewApplication()

var newBucketA = big.NewInt(57)
var newBucketB = big.NewInt(37)
var newDesired = big.NewInt(7)

var normStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
var pauseStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
var runStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
var doneStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
var currselStyle = new(tcell.Style).Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
var selStyle = new(tcell.Style).Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

var table = tview.NewTable().SetBorders(true).SetFixed(1, 1)
var insttext = tview.NewTextView()
var problemtext = tview.NewTextView()
var solutiontext = tview.NewTextView()

var infopane = tview.NewFlex().SetDirection(tview.FlexColumn).
	AddItem(problemtext, 0, 1, false).
	AddItem(solutiontext, 0, 1, false)
var layout = tview.NewFlex().SetDirection(tview.FlexRow).
	AddItem(table, 0, 6, true).
	AddItem(infopane, 12, 6, false).
	AddItem(insttext, 2, 1, false)
var form = tview.NewForm().
	AddInputField("Bigger Bucket", newBucketA.Text(10), 30, nil, func(newsval string) {
		newBucketA = big.NewInt(0)
		newBucketA, _ = newBucketA.SetString(newsval, 10)
	}).
	AddInputField("Smaller Bucket", newBucketB.Text(10), 30, nil, func(newsval string) {
		newBucketB = big.NewInt(0)
		newBucketB, _ = newBucketB.SetString(newsval, 10)
	}).
	AddInputField("Desired", newDesired.Text(10), 30, nil, func(newsval string) {
		newDesired = big.NewInt(0)
		newDesired, _ = newDesired.SetString(newsval, 10)
	}).
	AddButton("Cancel", func() {
		app.SetRoot(layout, true)
	}).
	AddButton("Save", func() {
		submitNewProblem()
		app.SetRoot(layout, true)
	})

func submitNewProblem() {
	bucketA := new(big.Int).Set(newBucketA)
	bucketB := new(big.Int).Set(newBucketB)
	desired := new(big.Int).Set(newDesired)
	problem := biglib.NewProblem(bucketA, bucketB, desired)
	newJob(table, problem)
}

type (
	JobList []*Job
)

var lastselected = -1
var joblist = make([]Job, 3)
var joblistmutex = new(sync.Mutex)

func main() {
	insttext.SetText(getInitialInstructions())
	problemtext.SetText("Select a Job From the Table")
	solutiontext.SetText("")
	table.Select(0, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyTab {
				joblistmutex.Lock()
				if lastselected > -1 && lastselected < len(joblist) {
					pjob := joblist[lastselected]
					*pjob.controller.GetStopStartChannel() <- biglib.Pause
				}
				joblistmutex.Unlock()
			}
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, false)
				insttext.SetText(getSelectingInstructions())
			}
		}).
		SetSelectedFunc(
			func(row int, column int) {
				selectRow(row)
			}).
		SetSelectedStyle(selStyle.Decompose())

	ebucketa, _ := new(big.Int).SetString("5", 10)
	ebucketb, _ := new(big.Int).SetString("3", 10)
	edesired, _ := new(big.Int).SetString("4", 10)

	dbucketa, _ := new(big.Int).SetString("12", 10)
	dbucketb, _ := new(big.Int).SetString("4", 10)
	ddesired, _ := new(big.Int).SetString("3", 10)

	cbucketa, _ := new(big.Int).SetString("57", 10)
	cbucketb, _ := new(big.Int).SetString("41", 10)
	cdesired, _ := new(big.Int).SetString("20", 10)

	bbucketa, _ := new(big.Int).SetString("100000013", 10)
	bbucketb, _ := new(big.Int).SetString("10000013", 10)
	bdesired, _ := new(big.Int).SetString("100003", 10)

	bucketa, _ := new(big.Int).SetString("1000013", 10)
	bucketb, _ := new(big.Int).SetString("10013", 10)
	desired, _ := new(big.Int).SetString("10003", 10)

	problem := biglib.NewProblem(bucketa, bucketb, desired)
	bproblem := biglib.NewProblem(bbucketa, bbucketb, bdesired)
	cproblem := biglib.NewProblem(cbucketa, cbucketb, cdesired)
	dproblem := biglib.NewProblem(dbucketa, dbucketb, ddesired)
	eproblem := biglib.NewProblem(ebucketa, ebucketb, edesired)

	joblist = []Job{}

	newJob(table, problem)
	newJob(table, bproblem)
	newJob(table, cproblem)
	newJob(table, dproblem)
	newJob(table, eproblem)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			app.SetRoot(form, true)
		}
		return event
	})

	if len(joblist) > 0 {
		selectRow(0)
	}

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func selectRow(row int) {
	table.GetCell(row, 0).SetStyle(currselStyle)
	table.SetSelectable(false, false)
	if lastselected > -1 && lastselected < len(joblist) && row != lastselected {
		joblistmutex.Lock()
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
		joblistmutex.Unlock()
	}
	if row > -1 && row < len(joblist) {
		joblistmutex.Lock()
		lastselected = row
		newjob := joblist[row]
		solutiontext.SetText(printSolution(newjob.controller.Solution))
		problemtext.SetText(printProblem(newjob.controller.Solution.Problem))
		insttext.SetText(getSelectedInstructions())
		joblistmutex.Unlock()
	}
}

func addJob(j *Job) {
	joblistmutex.Lock()
	newidx := len(joblist)
	j.pos = newidx
	joblist = append(joblist, []Job{*j}...)
	joblistmutex.Unlock()
	table.SetCell(newidx, 0,
		j.statecell.
			SetStyle(normStyle).
			SetAlign(tview.AlignCenter))

	table.SetCell(newidx, 1,
		j.simcell.
			SetStyle(normStyle).
			SetAlign(tview.AlignLeft).SetExpansion(40))
}

func newJob(table *tview.Table, problem *biglib.Problem) (j *Job) {
	j = new(Job)
	j.statecell = tview.NewTableCell("hi")
	j.simcell = tview.NewTableCell("there")

	j.statecell.SetText("bye")
	j.simcell.SetText("here")

	solution := biglib.NewSolution(problem)

	j.chans = newChans()

	j.controller = biglib.GetIdleSolutionProcessor("DUMMY", solution, &j.chans.stateChannel, &j.chans.resultsChannel)
	j.startListenChannels()
	addJob(j)
	defer func() {
		controlChannel := j.controller.GetStopStartChannel()
		*controlChannel <- biglib.Start
	}()

	selectRow(j.pos)
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
			if j.pos == lastselected {
				display.SetStyle(currselStyle)
			} else {
				display.SetStyle(doneStyle)
			}
			insttext.SetText(getSelectedInstructions())
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

			if int(op)&int(biglib.Error) > 0 {
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
			if lastselected == j.pos {
				solutiontext.SetText(printSolution(j.controller.Solution))
			}
			time.Sleep(time.Millisecond * time.Duration(1))
			app.Draw()
			if bucket.Operation >= bp.FinalOp {
				doit = false
				break
			}
		default:
		}
	}
}

func printLastTableEntry(lastbucket biglib.SimulationState) string {
	return fmt.Sprintf("%15v) | %16v | %c | %v\n", lastbucket.Idx, lastbucket.AmountBucketA, lastbucket.Operation.Rune(), lastbucket.AmountBucketB)
}

func printState(c *biglib.ChannelController) string {
	pausepart := "Sim"
	if c.IsTerminated() {
		pausepart = "Done"
	} else if c.IsPaused() {
		pausepart = "PAUSED"
	} else {
		stage := c.GetStage()
		if stage == biglib.StageGcd {
			pausepart = "GCD"
		} else if stage == biglib.Noop {
			pausepart = "Ready"
		}
	}
	direction := "A->B (-)"
	if c.Solution.FromB {
		direction = "B->A (+)"
	}
	problempart := fmt.Sprintf("N: %12v", c.Solution.Problem.Desired)
	prediction := "Unsolvable"
	if c.Solution.PredictedStateCount != nil {
		prediction = fmt.Sprintf("%12v", c.Solution.PredictedStateCount)
	}
	solutionpart := fmt.Sprintf("C:%2x D:%8s P: %12s", c.Solution.GetComplexityScale(), direction, prediction)
	return fmt.Sprintf("%s %s %10s", problempart, solutionpart, pausepart)
}

func printProblem(problem *biglib.Problem) (r string) {
	r = fmt.Sprintf("Problem (Hash: %x)\n", problem.Hash()) +
		fmt.Sprintf("- Bucket A: %v\n", problem.BucketA) +
		fmt.Sprintf("- Bucket B: %v\n", problem.BucketB) +
		fmt.Sprintf("- Desired: %v\n", problem.Desired)
	return r
}

func printSolution(solution *biglib.Solution) (r string) {
	sdirection := "Subtractive (A -> B)"
	if solution.FromB {
		sdirection = "Additive (A <- B)"
	}
	prediction := "Unsolvable"
	if solution.PredictedStateCount != nil {
		prediction = fmt.Sprintf("%v", solution.PredictedStateCount)
	}
	r = fmt.Sprintf("Solution\n") +
		fmt.Sprintf("- Result: %s\n", solution.Code) +
		fmt.Sprintf("- Complexity: %v\n", solution.Complexity) +
		fmt.Sprintf("- GCD: %v\n", solution.Denominator) +
		fmt.Sprintf("- GCDNanoTime: %v\n", solution.GCDNanoTime) +
		fmt.Sprintf("- Direction: %s\n", sdirection) +
		fmt.Sprintf("- CountFromA: %v\n", solution.CountFromA) +
		fmt.Sprintf("- CountFromB: %v\n", solution.CountFromB) +
		fmt.Sprintf("- PredictedSteps: %v\n", prediction) +
		fmt.Sprintf("- Simulated Steps: %v\n", solution.Operations.GetNextIndex())
	return r
}

func getInitialInstructions() (s string) {
	return "[Enter] Begin Selection, [Ctl+N] New Problem, [esc] Exit"
}

func getSelectingInstructions() (s string) {
	return "[Enter] Select Job, [UpArrow] Move Selection Up, [DownArrow] Move Selection Down, [Ctl+N] New Problem, [esc] Exit"
}

func getSelectedInstructions() (s string) {
	if lastselected > -1 && lastselected < len(joblist) && &joblist[lastselected] != nil &&
		!joblist[lastselected].controller.IsTerminated() {
		return getInitialInstructions() + ", " +
			"[Tab] Pause/UnPause"
	}
	return getInitialInstructions()
}
