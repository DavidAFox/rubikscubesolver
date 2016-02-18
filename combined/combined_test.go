package combined

import (
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
	"github.com/davidafox/rubikscubesolver/rubikscuberunner"
	"testing"
)

func TestSolveRealCube(t *testing.T) {
	data2 := []string{
		"L",
		"R U L F L",
		"R2 U L U2 R'",
		"R R D'",
		"F R B L U R F' R2",
	}
	for i, x := range data2 {
		c, _ := bytecube.NewCube("000000000111111111222222222333333333444444444555555555")
		r := rubikscuberunner.NewOfficialRunner(c)
		r.Run(x)
		s := NewSolver(c.String(), NewFactory(), 2)
		result := s.Solve()
		fmt.Println("Finished: ", i)
		fmt.Println("Solution:", result)
		r2 := rubikscuberunner.NewOfficialRunner(c)
		r2.Run(result)
		if !c.Solved() {
			t.Error("Failed to solve got: ", c.String())
		}
	}
}
