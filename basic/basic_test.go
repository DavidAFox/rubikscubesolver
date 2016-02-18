package basic

import (
	"github.com/davidafox/rubikscubesolver/rubikscube"
	"github.com/davidafox/rubikscubesolver/rubikscuberunner"
	"testing"
)

func TestSolve(t *testing.T) {
	c := rubikscube.NewCube("000000000111111111222222222333333333444444444555555555", 3)
	r := rubikscuberunner.NewRunner(c)
	r.Run("R1L1U0C1L0L0C1L0")
	s := NewSolver(c.String(), 3, NewFactory(3))
	result := s.Solve()
	r.Run(result)
	if !c.Solved() {
		t.Error("Failed to solve got: ", c.String())
	}
}
