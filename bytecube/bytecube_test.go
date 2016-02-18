package bytecube

import (
	"strconv"
	"strings"
	"testing"
)

var solvedCube = "000000000111111111222222222333333333444444444555555555"

func TestNewCubeAndString(t *testing.T) {
	data := []string{
		"000000000111111111222222222333333333344444444455555555",
		"555555555444444444333333333222222222111111111000000000",
	}
	for _, x := range data {
		c, _ := NewCube(x)
		if c.String() != x {
			t.Error("Failed NewCube/String \ngot: \t\t", c.String(), "\nexpected: \t", x)
		}
	}
}

func TestSideString(t *testing.T) {
	data := []struct {
		binary string
		result string
	}{
		{
			"001001001001001001001001001",
			"111111111",
		},
		{
			"000000000000000000000000000",
			"000000000",
		},
	}
	for _, x := range data {
		v64, err := strconv.ParseInt(x.binary, 2, 32)
		if err != nil {
			t.Error(err)
		}
		v := uint32(v64)
		got := sideString(v)
		if got != x.result {
			t.Error("Failed sideString got: ", got, " expected: ", x.result)
		}
	}
}

func TestSetLocation(t *testing.T) {
	data := []struct {
		side     int
		spot     int
		value    int
		starting string
		result   string
	}{
		{
			0,
			0,
			1,
			"000000000111111111222222222333333333444444444555555555",
			"100000000111111111222222222333333333444444444555555555",
		},
		{
			2,
			5,
			0,
			"000000000111111111222222222333333333444444444555555555",
			"000000000111111111222220222333333333444444444555555555",
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.starting)
		c.setLocation(x.side, x.spot, x.value)
		got := c.String()
		if got != x.result {
			t.Error("Failed setLocation got: ", got, " expected: ", x.result)
		}
	}
}

func TestGetLocation(t *testing.T) {
	data := []struct {
		side   int
		spot   int
		state  string
		result int
	}{
		{
			0,
			0,
			"000000000111111111222222222333333333444444444555555555",
			0,
		},
		{
			2,
			2,
			"000000000111111111222222222333333333444444444555555555",
			2,
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		got := c.getLocation(x.side, x.spot)
		if got != x.result {
			t.Error("Failed getLocation got: ", got, " expected: ", x.result)
		}
	}
}

func TestGetRow(t *testing.T) {
	data := []struct {
		side      int
		rowNumber int
		state     string
		result    []int
	}{
		{
			3,
			1,
			"000000000111111111222222222333333333444444444555555555",
			[]int{3, 3, 3},
		},
		{
			5,
			2,
			"000000000111111111222222222333333333444444444555555555",
			[]int{5, 5, 5},
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		got := c.getRow(x.side, x.rowNumber)
		failed := false
		if len(got) != len(x.result) {
			failed = true
		} else {
			for i, y := range got {
				if y != x.result[i] {
					failed = true
				}
			}
		}
		if failed {
			t.Error("Failed getRow got: ", got, " expected: ", x.result)
		}
	}
}
func TestGetColumn(t *testing.T) {
	c, _ := NewCube("011011011221221221222222222333333333141141141005005005")
	data := []struct {
		side         int
		columnNumber int
		result       []int
	}{
		{0, 0, []int{0, 0, 0}},
		{1, 2, []int{1, 1, 1}},
		{4, 1, []int{4, 4, 4}},
		{5, 2, []int{5, 5, 5}},
	}
	for _, x := range data {
		got := c.getColumn(x.side, x.columnNumber)
		failed := false
		if len(got) != len(x.result) {
			failed = true
		} else {
			for i, y := range got {
				if y != x.result[i] {
					failed = true
				}
			}
		}
		if failed {
			t.Error("Failed getColumn got: ", got, " expected: ", x.result)
		}
	}
}

func TestSetRow(t *testing.T) {
	data := []struct {
		side      int
		rowNumber int
		row       []int
		result    string
	}{
		{0, 0, []int{1, 1, 1}, "111000000111111111222222222333333333444444444555555555"},
		{1, 2, []int{3, 3, 3}, "000000000111111333222222222333333333444444444555555555"},
		{4, 1, []int{0, 0, 0}, "000000000111111111222222222333333333444000444555555555"},
		{5, 2, []int{2, 2, 2}, "000000000111111111222222222333333333444444444555555222"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.setRow(x.side, x.rowNumber, x.row)
		if c.String() != x.result {
			t.Error("Failed setRow got: ", c.String(), " expected: ", x.result)
		}
	}
}
func TestSetColumn(t *testing.T) {
	data := []struct {
		side         int
		columnNumber int
		column       []int
		result       string
	}{
		{0, 0, []int{1, 1, 1}, "100100100111111111222222222333333333444444444555555555"},
		{1, 2, []int{3, 3, 3}, "000000000113113113222222222333333333444444444555555555"},
		{4, 1, []int{0, 0, 0}, "000000000111111111222222222333333333404404404555555555"},
		{5, 2, []int{2, 2, 2}, "000000000111111111222222222333333333444444444552552552"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.setColumn(x.side, x.columnNumber, x.column)
		if c.String() != x.result {
			t.Error("Failed setColumn got: ", c.String(), " expected: ", x.result)
		}
	}
}
func TestRotateRowRight(t *testing.T) {
	data := []struct {
		row    int
		result string
	}{
		{0, "333000000000111111111222222222333333444444444555555555"},
		{1, "000333000111000111222111222333222333444444444555555555"},
		{2, "000000333111111000222222111333333222444444444555555555"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.rotateRowRight(x.row)
		if c.String() != x.result {
			t.Error("Failed RotateRight \ngot: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}
func TestRotateRowLeft(t *testing.T) {
	data := []struct {
		row    int
		result string
	}{
		{0, "111000000222111111333222222000333333444444444555555555"},
		{1, "000111000111222111222333222333000333444444444555555555"},
		{2, "000000111111111222222222333333333000444444444555555555"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.rotateRowLeft(x.row)
		if c.String() != x.result {
			t.Error("Failed RotateLeft \ngot: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}
func TestRotateColumnUp(t *testing.T) {
	data := []struct {
		column int
		result string
	}{
		{0, "500500500111111111224224224333333333044044044255255255"},
		{1, "050050050111111111242242242333333333404404404525525525"},
		{2, "005005005111111111422422422333333333440440440552552552"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.rotateColumnUp(x.column)
		if c.String() != x.result {
			t.Error("Failed RotateUp \ngot: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}

func TestRotateColumnDown(t *testing.T) {
	data := []struct {
		column int
		result string
	}{
		{0, "400400400111111111225225225333333333244244244055055055"},
		{1, "040040040111111111252252252333333333424424424505505505"},
		{2, "004004004111111111522522522333333333442442442550550550"},
	}
	for _, x := range data {
		c, _ := NewCube(solvedCube)
		c.rotateColumnDown(x.column)
		if c.String() != x.result {
			t.Error("Failed RotateDown \ngot: \t\t", c.String(), "\nexpected: \t", x.result)
		}

	}
}

func TestRotateClockwise(t *testing.T) {
	data := []struct {
		starting string
		result   string
		level    int
	}{
		{
			"000000000111111111222222222333333333444444444555555555",
			"000000000115115115222222222433433433444444111333555555",
			0,
		},
		{
			"111000111 111222333 222333444 555444333 111222111 111333222",
			"101101101 131231331 222333444 125124123 111321321 345345222",
			1,
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateClockwise(x.level)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateClockwise \n got: \t\t", c.String(), "\nexpected: \t", expected)
		}
	}
}

func TestRotateSideClockwise(t *testing.T) {
	data := []struct {
		start  string
		result string
		side   int
	}{
		{
			"123455555111111111222222222333333333444444444555555555",
			"541552553111111111222222222333333333444444444555555555",
			0,
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.start)
		c.rotateSideClockwise(x.side)
		if c.String() != x.result {
			t.Error("Failed rotateSideClockwise got: ", c.String(), " expected: ", x.result)
		}
	}
}

func TestRotateCounterClockwise(t *testing.T) {
	data := []struct {
		starting string
		result   string
		level    int
	}{
		{
			"000000000111111111222222222333333333444444444555555555",
			"000000000144144144222222222553553553444333333111111555",
			1,
		},
		{
			"111000111 111222333 222333444 555444333 111222111 111333222",
			"101101101 111221331 222333444 155144133 111222543 123333222",
			0,
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateCounterClockwise(x.level)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateCounterClockwise \n got: \t\t", c.String(), "\nexpected: \t", expected)
		}
	}
}

func TestRotateRight(t *testing.T) {
	data := []struct {
		starting string
		result   string
		rows     int
	}{
		{
			"000000000111111111222222222333333333111444111555555555",
			"333333000000000111111111222222222333141141141555555555",
			1,
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.starting)
		c.RotateRight(x.rows)
		if c.String() != x.result {
			t.Error("Failed RotateRight \n got: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}

func TestRotateLeft(t *testing.T) {
	data := []struct {
		starting string
		result   string
		rows     int
		size     int
	}{
		{
			"000000000 111111111 222222222 333333333 111444111 555555555",
			"111111000 222222111 333333222 000000333 141141141 555555555",
			1,
			3,
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateLeft(x.rows)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateLeft \n got: \t\t", c.String(), "\nexpected: \t", expected)
		}
	}
}

func TestRotateUp(t *testing.T) {
	data := []struct {
		starting string
		result   string
		columns  int
	}{
		{
			"000000000 111444111 222222222 333333333 444444444 555555555",
			"550550550 141141141 244244244 333333333 004004004 225225225",
			1,
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateUp(x.columns)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateUp \n got: \t\t", c.String(), "\nexpected: \t", expected)
		}
	}
}

func TestRotateDown(t *testing.T) {
	data := []struct {
		starting string
		result   string
		columns  int
	}{
		{
			"000000000 111444111 222222222 333333333 444444444 555555555",
			"440440440 141141141 255255255 333333333 224224224 005005005",
			1,
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateDown(x.columns)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateDown \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateR(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 111333111 444444444 555555555",
			"005005005 111111111 422422422 131131131 440440440 552552552",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateR()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateR \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateRCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 111333111 444444444 555555555",
			"004004004 111111111 522522522 131131131 442442442 550550550",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateRCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateRCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateL(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111333111 222222222 111333111 444444444 555555555",
			"400400400 131131131 225225225 111333111 244244244 055055055",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateL()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateL \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateLCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111333111 222222222 111333111 444444444 555555555",
			"500500500 131131131 224224224 111333111 044044044 255255255",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateLCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateLCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateU(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 333333333 111444111 111555111",
			"333000000 000111111 111222222 222333333 141141141 111555111",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateU()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateU \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateUCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 333333333 111444111 111555111",
			"111000000 222111111 333222222 000333333 141141141 111555111",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateUCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateUCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateD(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 333333333 111444111 111555111",
			"000000111 111111222 222222333 333333000 111444111 151151151",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateD()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateD \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateDCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"000000000 111111111 222222222 333333333 111444111 111555111",
			"000000333 111111000 222222111 333333222 111444111 151151151",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateDCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateDCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateF(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"111000111 111111111 111222111 333333333 444444444 555555555",
			"101101101 115115115 111222111 433433433 444444111 333555555",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateF()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateF \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateFCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"111000111 111111111 111222111 333333333 444444444 555555555",
			"101101101 114114114 111222111 533533533 444444333 111555555",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateFCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateFCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateB(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"111000111 111111111 111222111 333333333 444444444 555555555",
			"111000111 411411411 121121121 335335335 333444444 555555111",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateB()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateB \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestRotateBCounter(t *testing.T) {
	data := []struct {
		starting string
		result   string
	}{
		{
			"111000111 111111111 111222111 333333333 444444444 555555555",
			"111000111 511511511 121121121 334334334 111444444 555555333",
		},
	}
	for _, x := range data {
		c, _ := NewCube(strings.Replace(x.starting, " ", "", -1))
		c.RotateBCounter()
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateBCounter \n got: \t\t", c.String(), " \nexpected: \t", expected)
		}
	}
}

func TestSolved(t *testing.T) {
	data := []struct {
		state  string
		result bool
	}{
		{"000000000111111111222222222333333333444444444555555555", true},
		{"400400400111111111522522522333333333244244244055055055", false},
		{"500500500111111111422422422333333333044044044255255255", false},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if c.Solved() != x.result {
			t.Error("Failed Solved got: ", c.Solved(), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestCheckCenters(t *testing.T) {
	data := []struct {
		state  string
		result bool
	}{
		{"000000000111111111222222222333333333444444444555555555", true},
		{"444444444444444444444444444444444444444444444444444444", false},
		{"444444444444444444444444444455555555544442222222222222", false},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if checkCenters(c) != x.result {
			t.Error("Failed checkCenters got: ", checkCenters(c), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestCheckColors(t *testing.T) {
	data := []struct {
		state  string
		result bool
	}{
		{"000000000111111111222222222333333333444444444555555555", true},
		{"444444444444444444444444444444444444444444444444444444", false},
		{"444444444444444444444444444455555555544442222222222222", false},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if checkColors(c) != x.result {
			t.Error("Failed checkColors got: ", checkColors(c), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestCheckCorners(t *testing.T) {
	data := []struct {
		state  string
		result bool
	}{
		{"000000000111111111222222222333333333444444444555555555", true},
		{"000000001111111111222222222333333333444444444555555555", false},
		{"444444444444444444444444444444444444444444444444444444", false},
		{"444444444444444444444444444455555555544442222222222222", false},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if checkCorners(c) != x.result {
			t.Error("Failed checkColors got: ", checkCorners(c), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestCheckSides(t *testing.T) {
	data := []struct {
		state  string
		result bool
	}{
		{"000000000111111111222222222333333333444444444555555555", true},
		{"000000010111111111222222222333333333444444444555555555", false},
		{"444444444444444444444444444444444444444444444444444444", false},
		{"444444444444444444444444444455555555544442222222222222", false},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if checkSides(c) != x.result {
			t.Error("Failed checkColors got: ", checkSides(c), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestSolvedState(t *testing.T) {
	data := []struct {
		state  string
		result string
	}{
		{
			"111101111000010000333323333222232222555545555444454444",
			"000000000111111111222222222333333333444444444555555555",
		},
	}
	for _, x := range data {
		c, _ := NewCube(x.state)
		if c.SolvedState() != x.result {
			t.Error("Failed SolvedState got: ", c.SolvedState(), "expected: ", x.result)
		}
	}

}
