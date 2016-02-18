package depthfirst

import (
	"github.com/davidafox/rubikscubesolver/bytecube"
	"github.com/davidafox/rubikscubesolver/rubikscuberunner"
	"testing"
)

func TestSolveRealCube(t *testing.T) {
	data := []string{
		"L0",
		"R1U0L0C1L0",
		"r1U0L0u1T0",
		"R0R0D1",
	}
	for _, x := range data {
		c, _ := bytecube.NewCube("000000000111111111222222222333333333444444444555555555")
		r := rubikscuberunner.NewRunner(c)
		r.Run(x)
		s := NewSolver(c.String(), NewFactory())
		result := s.Solve()
		r.Run(result)
		if !c.Solved() {
			t.Error("Failed to solve got: ", c.String())
		}
	}
}
