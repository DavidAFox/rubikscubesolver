package combined

import (
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
	"runtime"
)

type Cube interface {
	RotateR()
	RotateRCounter()
	RotateL()
	RotateLCounter()
	RotateU()
	RotateUCounter()
	RotateD()
	RotateDCounter()
	RotateF()
	RotateFCounter()
	RotateB()
	RotateBCounter()
	Solved() bool
	State() bytecube.State
}

type rotations struct {
	fun    func(Cube) (bytecube.State, bool)
	letter string
}

const STARTING_SIDE = 0
const SOLVED_SIDE = 1

type CubeFactory interface {
	New(bytecube.State) Cube
}

type Factory struct {
}

func (f *Factory) New(s bytecube.State) Cube {
	return bytecube.NewWithState(s)
}

func NewFactory() *Factory {
	f := new(Factory)
	return f
}

type Solver struct {
	startingState bytecube.State
	solvedState   bytecube.State
	factory       CubeFactory
	foundStates   []map[bytecube.State]string
	rotations     []rotations
	depth         int
}

type cubeState struct {
	state bytecube.State
	steps string
}

func NewSolver(startingState string, factory CubeFactory, depth int) *Solver {
	s := new(Solver)
	s.factory = factory
	c, err := bytecube.NewCube(startingState)
	if err != nil {
		fmt.Println("Invalid State in solver")
		return nil
	}
	s.startingState = c.State()
	solvedStateCube, _ := bytecube.NewCube(c.SolvedState())
	s.solvedState = solvedStateCube.State()
	s.foundStates = make([]map[bytecube.State]string, 2, 2)
	s.foundStates[0] = make(map[bytecube.State]string)
	s.foundStates[1] = make(map[bytecube.State]string)
	s.rotations = []rotations{
		{right, "R"},
		{rightCounter, "R'"},
		{left, "L"},
		{leftCounter, "L'"},
		{up, "U"},
		{upCounter, "U'"},
		{down, "D"},
		{downCounter, "D'"},
		{front, "F"},
		{frontCounter, "F'"},
		{back, "B"},
		{backCounter, "B'"},
		{doubleRight, "R2"},
		{doubleLeft, "L2"},
		{doubleUp, "U2"},
		{doubleDown, "D2"},
		{doubleFront, "F2"},
		{doubleBack, "B2"},
	}
	s.depth = depth
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
	for i := 0; i < s.depth; i++ {
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
						var solution string
						if L == 0 {
							solution = x.steps + " " + y
						} else {
							solution = y + " " + x.steps
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
	for i := 1; i <= 20-(s.depth*2); i++ {
		fmt.Println("Depth: ", i)
		for _, state := range states[0] {
			result, depth := s.SolveR(state.state, 0, i, state.steps)
			if depth > 0 {
				return result
			}
		}
	}
	return ""
}

type action struct {
	fun      func(Cube) (bytecube.State, bool)
	callback func(bytecube.State, int, int, string) (string, int)
	letter   string
}

func (s *Solver) SolveR(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromRight(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromLeft(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromUp(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDown(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromFront(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromBack(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleRight(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleLeft(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleUp(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleDown(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{front, s.solveRFromFront, "F"},
		{frontCounter, s.solveRFromFront, "F'"},
		{back, s.solveRFromBack, "B"},
		{backCounter, s.solveRFromBack, "B'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleFront, s.solveRFromDoubleFront, "F2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleFront(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
		{doubleBack, s.solveRFromDoubleBack, "B2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) solveRFromDoubleBack(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	actions := []action{
		{right, s.solveRFromRight, "R"},
		{rightCounter, s.solveRFromRight, "R'"},
		{left, s.solveRFromLeft, "L"},
		{leftCounter, s.solveRFromLeft, "L'"},
		{up, s.solveRFromUp, "U"},
		{upCounter, s.solveRFromUp, "U'"},
		{down, s.solveRFromDown, "D"},
		{downCounter, s.solveRFromDown, "D'"},
		{doubleRight, s.solveRFromDoubleRight, "R2"},
		{doubleLeft, s.solveRFromDoubleLeft, "L2"},
		{doubleUp, s.solveRFromDoubleUp, "U2"},
		{doubleDown, s.solveRFromDoubleDown, "D2"},
	}
	return s.genericSolveR(actions, state, depth, maxDepth, steps)
}

func (s *Solver) genericSolveR(actions []action, state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	if depth >= maxDepth {
		return "", -1
	}
	for i := 0; i < len(actions); i++ {
		cube := s.factory.New(state)
		rState, rSolved := actions[i].fun(cube)
		if _, ok := s.foundStates[1][rState]; ok {
			fmt.Println("Found Solution, depth: ", depth+1)
			return steps + " " + actions[i].letter + " " + s.foundStates[1][rState], depth + 1
		}
		if _, ok := s.foundStates[0][rState]; ok {
			return "", -1
		}
		if rSolved {
			fmt.Println("Found Solution not in map, depth: ", depth+1)
			return steps + " " + actions[i].letter, depth
		}
		rSteps, rDepth := actions[i].callback(rState, depth+1, maxDepth, steps+" "+actions[i].letter)
		if rDepth != -1 {
			return rSteps, rDepth
		}
	}
	return "", -1
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
	for _, x := range s.rotations {
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
	//separate rotations with reverse letters
	rotations := []struct {
		fun    func(Cube) (bytecube.State, bool)
		letter string
	}{
		{right, "R'"},
		{rightCounter, "R"},
		{left, "L'"},
		{leftCounter, "L"},
		{up, "U'"},
		{upCounter, "U"},
		{down, "D'"},
		{downCounter, "D"},
		{front, "F'"},
		{frontCounter, "F"},
		{back, "B'"},
		{backCounter, "B"},
		{doubleRight, "R2"},
		{doubleLeft, "L2"},
		{doubleUp, "U2"},
		{doubleDown, "D2"},
		{doubleFront, "F2"},
		{doubleBack, "B2"},
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

func (s *Solver) doOneTypeOfRotationConcurrent(cube *cubeState, f func(Cube) (bytecube.State, bool), rotationLetter string) ([]*cubeState, bool) {
	results := make([]*cubeState, 0)
	c := s.factory.New(cube.state)
	newState, solved := f(c)
	cs := newCubeState(newState, cube.steps+" "+rotationLetter)
	if solved {
		r := make([]*cubeState, 1)
		r[0] = cs
		return r, true
	}
	results = append(results, cs)
	return results, false
}

func (s *Solver) doOneTypeOfRotationConcurrentFromSolved(cube *cubeState, f func(Cube) (bytecube.State, bool), rotationLetter string) ([]*cubeState, bool) {
	results := make([]*cubeState, 0)
	c := s.factory.New(cube.state)
	newState, _ := f(c)
	cs := newCubeState(newState, rotationLetter+" "+cube.steps)
	results = append(results, cs)
	return results, false
}

func right(c Cube) (bytecube.State, bool) {
	c.RotateR()
	return c.State(), c.Solved()
}

func rightCounter(c Cube) (bytecube.State, bool) {
	c.RotateRCounter()
	return c.State(), c.Solved()
}

func left(c Cube) (bytecube.State, bool) {
	c.RotateL()
	return c.State(), c.Solved()
}

func leftCounter(c Cube) (bytecube.State, bool) {
	c.RotateLCounter()
	return c.State(), c.Solved()
}

func up(c Cube) (bytecube.State, bool) {
	c.RotateU()
	return c.State(), c.Solved()
}

func upCounter(c Cube) (bytecube.State, bool) {
	c.RotateUCounter()
	return c.State(), c.Solved()
}

func down(c Cube) (bytecube.State, bool) {
	c.RotateD()
	return c.State(), c.Solved()
}

func downCounter(c Cube) (bytecube.State, bool) {
	c.RotateDCounter()
	return c.State(), c.Solved()
}

func front(c Cube) (bytecube.State, bool) {
	c.RotateF()
	return c.State(), c.Solved()
}

func frontCounter(c Cube) (bytecube.State, bool) {
	c.RotateFCounter()
	return c.State(), c.Solved()
}

func back(c Cube) (bytecube.State, bool) {
	c.RotateB()
	return c.State(), c.Solved()
}

func backCounter(c Cube) (bytecube.State, bool) {
	c.RotateBCounter()
	return c.State(), c.Solved()
}

func doubleRight(c Cube) (bytecube.State, bool) {
	c.RotateR()
	c.RotateR()
	return c.State(), c.Solved()
}

func doubleLeft(c Cube) (bytecube.State, bool) {
	c.RotateL()
	c.RotateL()
	return c.State(), c.Solved()
}

func doubleUp(c Cube) (bytecube.State, bool) {
	c.RotateU()
	c.RotateU()
	return c.State(), c.Solved()
}

func doubleDown(c Cube) (bytecube.State, bool) {
	c.RotateD()
	c.RotateD()
	return c.State(), c.Solved()
}

func doubleFront(c Cube) (bytecube.State, bool) {
	c.RotateF()
	c.RotateF()
	return c.State(), c.Solved()
}

func doubleBack(c Cube) (bytecube.State, bool) {
	c.RotateB()
	c.RotateB()
	return c.State(), c.Solved()
}
