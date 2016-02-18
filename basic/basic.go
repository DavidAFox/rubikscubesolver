package basic

import (
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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
	size          int
	factory       CubeFactory
	foundStates   map[bytecube.State]bool
}

type cubeState struct {
	state bytecube.State
	steps string
}

func NewSolver(startingState string, size int, factory CubeFactory) *Solver {
	s := new(Solver)
	s.size = size
	s.factory = factory
	c, _ := bytecube.NewCube(startingState)
	s.startingState = c.State()
	s.foundStates = make(map[bytecube.State]bool)
	return s
}

func newCubeState(state bytecube.State, steps string) *cubeState {
	c := new(cubeState)
	c.state = state
	c.steps = steps
	return c
}

func (s *Solver) Solve() string {
	states := make([]*cubeState, 1, 1)
	states[0] = newCubeState(s.startingState, "")
	s.foundStates[s.startingState] = true
	currentStates := make([]*cubeState, 0, 0)
	for i := 0; ; i++ {
		currentStates := currentStates[:0]
		for _, x := range states {
			nstates, solved := s.nextStates(x, false)
			if solved {
				return nstates[0].steps
			}
			currentStates = append(currentStates, nstates...)
		}
		if i%1 == 0 {
			fmt.Println("On step ", i, "states: ", len(currentStates))
		}
		states, currentStates = currentStates, states
	}
}

func (s *Solver) SolveConcurrent() string {
	states := make([]*cubeState, 1, 1)
	states[0] = newCubeState(s.startingState, "")
	s.foundStates[s.startingState] = true
	res := make(chan *results)
	workers := s.spawnWorkers(res)
	for i := 0; ; i++ {
		currentStates := make([]*cubeState, 0, 10)
		workersRunning := 0
		jobsPerWorker := len(states) / len(workers)
		if jobsPerWorker == 0 {
			workers[0] <- states
			workersRunning = 1
		} else {
			for j, worker := range workers {
				if j == len(workers)-1 {
					worker <- states[j*jobsPerWorker:]
				} else {
					worker <- states[j*jobsPerWorker : (j+1)*jobsPerWorker]
				}
				workersRunning++
			}
		}
		for k := 0; k < workersRunning; k++ {
			result := <-res
			if result.solved {
				f, err := os.Create("main.mprof")
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(f)
				f.Close()
				for _, x := range workers {
					close(x)
				}
				return result.states[0].steps
			}
			for _, x := range result.states {
				if ok := s.foundStates[x.state]; !ok {
					currentStates = append(currentStates, x)
					s.foundStates[x.state] = true
				}
			}
		}
		fmt.Println("On step ", i, "states: ", len(currentStates))
		states, currentStates = currentStates, states
	}
}

func (s *Solver) spawnWorkers(res chan *results) []chan []*cubeState {
	channels := make([]chan []*cubeState, 0, 0)
	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		c := make(chan []*cubeState)
		channels = append(channels, c)
		go s.NewWorker(res, c)
	}
	return channels
}

type results struct {
	states []*cubeState
	solved bool
}

func (s *Solver) NewWorker(res chan *results, in chan []*cubeState) {
	for states := range in {
		resultStates := make([]*cubeState, 0)
		for _, state := range states {
			nstates, solved := s.nextStates(state, true)
			if solved {
				res <- &results{nstates, solved}
				return
			}
			resultStates = append(resultStates, nstates...)
		}
		res <- &results{resultStates, false}
	}
}

func (s *Solver) nextStates(cube *cubeState, concurrent bool) ([]*cubeState, bool) {
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
	}
	for _, x := range rotations {
		if concurrent {
			states, solved = s.doOneTypeOfRotationConcurrent(cube, x.fun, x.letter)
		} else {
			states, solved = s.doOneTypeOfRotation(cube, x.fun, x.letter)
		}
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

func (s *Solver) doOneTypeOfRotation(cube *cubeState, f func(Cube, int) (bytecube.State, bool), rotationLetter string) ([]*cubeState, bool) {
	results := make([]*cubeState, 0)
	for i := 0; i < s.size-1; i++ {
		c := s.factory.New(cube.state)
		newState, solved := f(c, i)
		ok := s.foundStates[newState]
		if !ok {
			cs := newCubeState(newState, cube.steps+rotationLetter+strconv.Itoa(i))
			if solved {
				r := make([]*cubeState, 1)
				r[0] = cs
				return r, true
			}
			results = append(results, cs)
			s.foundStates[newState] = true
		}
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
