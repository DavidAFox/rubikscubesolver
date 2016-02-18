package breadthfirst

import (
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
	"runtime"
	"strconv"
)

type Cube interface {
	RotateRight(int)
	RotateLeft(int)
	RotateUp(int)
	RotateDown(int)
	RotateClockwise(int)
	RotateCounterClockwise(int)
	String() string
	Solved() bool
	State() bytecube.State
}

const STARTING_SIDE = 0
const SOLVED_SIDE = 1

type CubeFactory interface {
	New(bytecube.State) Cube
}

type Factory struct {
	size int
}

func (f *Factory) New(s bytecube.State) Cube {
	return bytecube.NewWithState(s)
}

func NewFactory(size int) *Factory {
	f := new(Factory)
	f.size = size
	return f
}

type Solver struct {
	startingState bytecube.State
	solvedState   bytecube.State
	size          int
	factory       CubeFactory
	foundStates   []map[bytecube.State]string
}

type cubeState struct {
	state bytecube.State
	steps string
}

func NewSolver(startingState, solved string, size int, factory CubeFactory) *Solver {
	s := new(Solver)
	s.size = size
	s.factory = factory
	c, _ := bytecube.NewCube(startingState)
	s.startingState = c.State()
	c, _ = bytecube.NewCube(solved)
	s.solvedState = c.State()
	s.foundStates = make([]map[bytecube.State]string, 2, 2)
	s.foundStates[0] = make(map[bytecube.State]string)
	s.foundStates[1] = make(map[bytecube.State]string)
	return s
}

func newCubeState(state bytecube.State, steps string) *cubeState {
	c := new(cubeState)
	c.state = state
	c.steps = steps
	return c
}

func (s *Solver) Solve() string {
	states := make([][]*cubeState, 2, 2)
	states[0] = make([]*cubeState, 1, 1)
	states[1] = make([]*cubeState, 1, 1)
	states[0][0] = newCubeState(s.startingState, "")
	states[1][0] = newCubeState(s.solvedState, "")
	s.foundStates[0][s.startingState] = ""
	s.foundStates[1][s.solvedState] = ""
	res := make(chan *results)
	workers := s.spawnWorkers(res)
	currentStates := make([][]*cubeState, 2, 2)
	for i := 0; ; i++ {
		for L := 0; L < 2; L++ {
			currentStates[L] = currentStates[L][:0]
			workersRunning := 0
			jobsPerWorker := len(states[L]) / len(workers)
			if jobsPerWorker == 0 {
				workers[0] <- NewWorkerList(L, states[L])
				workersRunning = 1
			} else {
				for j, worker := range workers {
					if j == len(workers)-1 {

						worker <- NewWorkerList(L, states[L][j*jobsPerWorker:])
					} else {
						worker <- NewWorkerList(L, states[L][j*jobsPerWorker:(j+1)*jobsPerWorker])
					}
					workersRunning++
				}
			}
			for k := 0; k < workersRunning; k++ {
				result := <-res
				if result.solved {
					for _, x := range workers {
						close(x)
					}
					return result.states[0].steps
				}
				for _, x := range result.states {
					if y, ok := s.foundStates[(L+1)%2][x.state]; ok {
						fmt.Println(L)
						fmt.Println(x.steps, y)
						var solution string
						if L == 0 {
							solution = x.steps + y
						} else {
							solution = y + x.steps
						}
						for _, x := range workers {
							close(x)
						}
						return solution
					}
					_, ok := s.foundStates[L][x.state]
					if !ok {
						currentStates[L] = append(currentStates[L], x)
						s.foundStates[L][x.state] = x.steps
					}
				}
			}
			fmt.Println("On step ", i, "(", L, ")", "states: ", len(currentStates[L]))
			states[L], currentStates[L] = currentStates[L], states[L]
		}
	}
}

func (s *Solver) spawnWorkers(res chan *results) []chan *workerList {
	channels := make([]chan *workerList, 0, 0)
	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		c := make(chan *workerList)
		channels = append(channels, c)
		go s.NewWorker(res, c)
	}
	return channels
}

type workerList struct {
	side   int
	states []*cubeState
}

func NewWorkerList(side int, states []*cubeState) *workerList {
	wl := new(workerList)
	wl.side = side
	wl.states = states
	return wl
}

type results struct {
	states []*cubeState
	solved bool
}

func (s *Solver) NewWorker(res chan *results, in chan *workerList) {
	var nstates []*cubeState
	var solved bool
	for list := range in {
		resultStates := make([]*cubeState, 0)
		for _, state := range list.states {
			if list.side == STARTING_SIDE {
				nstates, solved = s.getNextStates(state, true)
			} else {
				nstates, solved = s.getNextStatesFromSolved(state, true)
			}
			if solved {
				res <- &results{nstates, solved}
				return
			}
			resultStates = append(resultStates, nstates...)
		}
		res <- &results{resultStates, false}
	}
}

func (s *Solver) getNextStates(cube *cubeState, concurrent bool) ([]*cubeState, bool) {
	var states []*cubeState
	var solved bool
	nextStates := make([]*cubeState, 0, 10)
	rotations := []struct {
		fun    func(Cube, int) (bytecube.State, bool)
		letter string
	}{
		{right, "R"},
		{left, "L"},
		{up, "U"},
		{down, "D"},
		{clockwise, "C"},
		{counterclockwise, "T"},
		{doubleRight, "r"},
		{doubleUp, "u"},
		{doubleClockwise, "c"},
	}
	for _, x := range rotations {
		states, solved = s.doOneTypeOfRotationConcurrent(cube, x.fun, x.letter)
		if solved {
			return states, true
		}
		nextStates = append(nextStates, states...)
	}
	return nextStates, false
}

func (s *Solver) getNextStatesFromSolved(cube *cubeState, concurrent bool) ([]*cubeState, bool) {
	var states []*cubeState
	var solved bool
	nextStates := make([]*cubeState, 0, 10)
	rotations := []struct {
		fun    func(Cube, int) (bytecube.State, bool)
		letter string
	}{
		{right, "L"},
		{left, "R"},
		{up, "D"},
		{down, "U"},
		{clockwise, "T"},
		{counterclockwise, "C"},
		{doubleRight, "r"},
		{doubleUp, "u"},
		{doubleClockwise, "c"},
	}
	for _, x := range rotations {
		states, solved = s.doOneTypeOfRotationConcurrentFromSolved(cube, x.fun, x.letter)
		if solved {
			return states, true
		}
		nextStates = append(nextStates, states...)
	}
	return nextStates, false
}

func (s *Solver) doOneTypeOfRotationConcurrent(cube *cubeState, f func(Cube, int) (bytecube.State, bool), rotationLetter string) ([]*cubeState, bool) {
	results := make([]*cubeState, 0)
	for i := 0; i < s.size-1; i++ {
		c := s.factory.New(cube.state)
		newState, solved := f(c, i)
		cs := newCubeState(newState, cube.steps+rotationLetter+strconv.Itoa(i))
		if solved {
			r := make([]*cubeState, 1)
			r[0] = cs
			return r, true
		}
		results = append(results, cs)
	}
	return results, false
}

func (s *Solver) doOneTypeOfRotationConcurrentFromSolved(cube *cubeState, f func(Cube, int) (bytecube.State, bool), rotationLetter string) ([]*cubeState, bool) {
	results := make([]*cubeState, 0)
	for i := 0; i < s.size-1; i++ {
		c := s.factory.New(cube.state)
		newState, _ := f(c, i)
		cs := newCubeState(newState, rotationLetter+strconv.Itoa(i)+cube.steps)
		results = append(results, cs)
	}
	return results, false
}

func right(c Cube, rows int) (bytecube.State, bool) {
	c.RotateRight(rows)
	return c.State(), c.Solved()
}

func left(c Cube, rows int) (bytecube.State, bool) {
	c.RotateLeft(rows)
	return c.State(), c.Solved()
}

func up(c Cube, rows int) (bytecube.State, bool) {
	c.RotateUp(rows)
	return c.State(), c.Solved()
}
func down(c Cube, rows int) (bytecube.State, bool) {
	c.RotateDown(rows)
	return c.State(), c.Solved()
}
func clockwise(c Cube, rows int) (bytecube.State, bool) {
	c.RotateClockwise(rows)
	return c.State(), c.Solved()
}
func counterclockwise(c Cube, rows int) (bytecube.State, bool) {
	c.RotateCounterClockwise(rows)
	return c.State(), c.Solved()
}
func doubleRight(c Cube, rows int) (bytecube.State, bool) {
	c.RotateRight(rows)
	c.RotateRight(rows)
	return c.State(), c.Solved()
}

func doubleUp(c Cube, rows int) (bytecube.State, bool) {
	c.RotateUp(rows)
	c.RotateUp(rows)
	return c.State(), c.Solved()
}
func doubleClockwise(c Cube, rows int) (bytecube.State, bool) {
	c.RotateClockwise(rows)
	c.RotateClockwise(rows)
	return c.State(), c.Solved()
}
