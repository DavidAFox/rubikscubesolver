package rubikscuberunner

import (
	"strconv"
	"testing"
)

type fakeCube struct {
	result string
}

func (fc *fakeCube) RotateClockwise(row int) {
	fc.result += "C" + strconv.Itoa(row)
}

func (fc *fakeCube) RotateCounterClockwise(row int) {
	fc.result += "T" + strconv.Itoa(row)
}

func (fc *fakeCube) RotateUp(col int) {
	fc.result += "U" + strconv.Itoa(col)
}

func (fc *fakeCube) RotateDown(col int) {
	fc.result += "D" + strconv.Itoa(col)
}

func (fc *fakeCube) RotateRight(row int) {
	fc.result += "R" + strconv.Itoa(row)
}

func (fc *fakeCube) RotateLeft(row int) {
	fc.result += "L" + strconv.Itoa(row)
}

func TestRun(t *testing.T) {
	data := []string{
		"C1R4D7U1C8T7L9",
		"C9D2R1U8C9",
		"C9",
		"",
	}
	for _, x := range data {
		fc := new(fakeCube)
		r := NewRunner(fc)
		r.Run(x)
		if fc.result != x {
			t.Error("Failed Run got: ", fc.result, " expected: ", x)
		}
	}
}
