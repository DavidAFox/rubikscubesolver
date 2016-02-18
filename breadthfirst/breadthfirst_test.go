package breadthfirst

import (
	"fmt"
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
		"C1L0R1U0T1T0R1U0C1R0T0",
	}
	for i, x := range data {
		c, _ := bytecube.NewCube("000000000111111111222222222333333333444444444555555555")
		r := rubikscuberunner.NewRunner(c)
		r.Run(x)
		s := NewSolver(c.String(), "000000000111111111222222222333333333444444444555555555", 3, NewFactory(3))
		result := s.Solve()
		fmt.Println("Finished: ", i)
		fmt.Println("Solution length: ", len(result)/2)
		r.Run(result)
		if !c.Solved() {
			t.Error("Failed to solve got: ", c.String())
		}
	}
}
