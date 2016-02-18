package depthfirst

import (
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
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
}

func (f *Factory) New(s bytecube.State) Cube {
	return bytecube.NewWithState(s)
}

func NewFactory() *Factory {
	f := new(Factory)
	return f
}

type Solver struct {
	funcs         []func(Cube, int) (bytecube.State, bool)
	funcsLetter   []string
	startingState string
	factory       CubeFactory
}

func NewSolver(state string, factory CubeFactory) *Solver {
	s := new(Solver)
	s.funcs = make([]func(Cube, int) (bytecube.State, bool), 9, 9)
	s.funcsLetter = make([]string, 9, 9)
	data := []struct {
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
	for i, x := range data {
		s.funcs[i] = x.fun
		s.funcsLetter[i] = x.letter
	}
	s.startingState = state
	s.factory = factory
	return s
}

func (s *Solver) Solve() string {
	c, _ := bytecube.NewCube(s.startingState)
	state := c.State()
	result, _ := s.SolveR(state, 0, 20, "")
	return result
}

func (s *Solver) SolveR(state bytecube.State, depth, maxDepth int, steps string) (string, int) {
	if depth >= maxDepth {
		return "", -1
	}
	currentMaxDepth := maxDepth
	currentSteps := steps
	foundSolution := false
	for i := 0; i < len(s.funcs); i++ {
		for j := 0; j < 2; j++ {
			cube := s.factory.New(state)
			rState, rSolved := s.funcs[i](cube, j)
			if rSolved {
				fmt.Println("Found Solution, depth: ", depth+1)
				return steps + s.funcsLetter[i] + strconv.Itoa(j), depth
			}
			rSteps, rDepth := s.SolveR(rState, depth+1, currentMaxDepth, steps+s.funcsLetter[i]+strconv.Itoa(j))
			if rDepth != -1 && rDepth < currentMaxDepth {
				currentMaxDepth = rDepth
				currentSteps = rSteps
				foundSolution = true
			}
		}
	}
	if foundSolution {
		return currentSteps, currentMaxDepth
	}
	return "", -1
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
