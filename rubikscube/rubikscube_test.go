package rubikscube

import (
	"strings"
	"testing"
)

var solvedCube = [][][]int{
	{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	},
	{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	},
	{
		{2, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
	},
	{
		{3, 3, 3},
		{3, 3, 3},
		{3, 3, 3},
	},
	{
		{4, 4, 4},
		{4, 4, 4},
		{4, 4, 4},
	},
	{
		{5, 5, 5},
		{5, 5, 5},
		{5, 5, 5},
	},
}

func TestConvertToString(t *testing.T) {
	expected := "000000000111111111222222222333333333444444444555555555"
	got := convertToString(solvedCube)
	if expected != got {
		t.Error("Failed to convert to string got: ", got, " expected: ", expected)
	}
}

func TestGetRow(t *testing.T) {
	c := NewCubeSlice(solvedCube, 3)
	data := []struct {
		side      int
		rowNumber int
		result    string
	}{
		{0, 0, "000"},
		{1, 2, "111"},
		{4, 1, "444"},
		{5, 2, "555"},
	}
	for _, x := range data {
		got := c.getRow(x.side, x.rowNumber)
		if got != x.result {
			t.Error("Failed getRow got: ", got, " expected: ", x.result)
		}
	}
}

func TestSetRow(t *testing.T) {
	data := []struct {
		side      int
		rowNumber int
		row       string
		result    string
	}{
		{0, 0, "111", "111000000111111111222222222333333333444444444555555555"},
		{1, 2, "333", "000000000111111333222222222333333333444444444555555555"},
		{4, 1, "000", "000000000111111111222222222333333333444000444555555555"},
		{5, 2, "222", "000000000111111111222222222333333333444444444555555222"},
	}
	for _, x := range data {
		c := NewCubeSlice(solvedCube, 3)
		c.setRow(x.side, x.rowNumber, x.row)
		if c.String() != x.result {
			t.Error("Failed setRow got: ", c.String(), " expected: ", x.result)
		}
	}
}

func TestGetColumn(t *testing.T) {
	c := NewCube("011011011221221221222222222333333333141141141005005005", 3)
	data := []struct {
		side         int
		columnNumber int
		result       string
	}{
		{0, 0, "000"},
		{1, 2, "111"},
		{4, 1, "444"},
		{5, 2, "555"},
	}
	for _, x := range data {
		got := c.getColumn(x.side, x.columnNumber)
		if got != x.result {
			t.Error("Failed getColumn got: ", got, " expected: ", x.result)
		}
	}
}

func TestSetColumn(t *testing.T) {
	data := []struct {
		side         int
		columnNumber int
		column       string
		result       string
	}{
		{0, 0, "111", "100100100111111111222222222333333333444444444555555555"},
		{1, 2, "333", "000000000113113113222222222333333333444444444555555555"},
		{4, 1, "000", "000000000111111111222222222333333333404404404555555555"},
		{5, 2, "222", "000000000111111111222222222333333333444444444552552552"},
	}
	for _, x := range data {
		c := NewCubeSlice(solvedCube, 3)
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
		c := NewCubeSlice(solvedCube, 3)
		c.rotateRowRight(x.row)
		if c.String() != x.result {
			t.Error("Failed RotateRight \ngot: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}

func TestRotateRowRight5(t *testing.T) {
	data := []struct {
		row    int
		result string
	}{
		{0, "333330000000000000000000000000111111111111111111111111122222222222222222222222223333333333333333333344444444444444444444444445555555555555555555555555"},
		{1, "000003333300000000000000011111000001111111111111112222211111222222222222222333332222233333333333333344444444444444444444444445555555555555555555555555"},
		{4, "000000000000000000003333311111111111111111111000002222222222222222222211111333333333333333333332222244444444444444444444444445555555555555555555555555"},
	}
	for _, x := range data {
		c := NewCube("000000000000000000000000011111111111111111111111112222222222222222222222222333333333333333333333333344444444444444444444444445555555555555555555555555", 5)
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
		c := NewCubeSlice(solvedCube, 3)
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
		c := NewCubeSlice(solvedCube, 3)
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
		c := NewCubeSlice(solvedCube, 3)
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
		{
			"123456789 123456789 123456789 123456789 123456789 123456789",
			"741852963 141452763 123456789 743856969 123852963 741852789",
			1,
		},
	}
	for _, x := range data {
		c := NewCube(strings.Replace(x.starting, " ", "", -1), 3)
		c.RotateClockwise(x.level)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateClockwise \n got: \t\t", c.String(), "\nexpected: \t", expected)
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
			"123456789 123456789 123456789 123456789 123456789 123456789",
			"369258147 169458747 123456789 363256149 123258147 369258789",
			1,
		},
		{
			"111000111 111222333 222333444 555444333 111222111 111333222",
			"101101101 111221331 222333444 155144133 111222543 123333222",
			0,
		},
	}
	for _, x := range data {
		c := NewCube(strings.Replace(x.starting, " ", "", -1), 3)
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
		size     int
	}{
		{
			"000000000111111111222222222333333333123456789123456789",
			"333333000000000111111111222222222333741852963123456789",
			1,
			3,
		},
	}
	for _, x := range data {
		c := NewCube(x.starting, x.size)
		c.RotateRight(x.rows)
		if c.String() != x.result {
			t.Error("Failed RotateRight \n got: \t\t", c.String(), "\nexpected: \t", x.result)
		}
	}
}
func TestRightAndLeft(t *testing.T) {
	startingState := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQR"
	c := NewCube(startingState, 3)
	c.RotateRight(0)
	c.RotateLeft(0)
	if c.String() != startingState {
		t.Error("Failed right and left \n got: \t\t", c.String(), "\nexpected: \t", startingState)
	}
}

func TestUpAndDown(t *testing.T) {
	startingState := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQR"
	c := NewCube(startingState, 3)
	c.RotateUp(1)
	c.RotateDown(1)
	if c.String() != startingState {
		t.Error("Failed up and down \n got: \t\t", c.String(), "\nexpected: \t", startingState)
	}
}

func TestTimes4(t *testing.T) {
	startingState := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQR"
	c := NewCube(startingState, 3)
	data := []struct {
		fun  func(int)
		name string
	}{
		{
			c.RotateUp,
			"RotateUp",
		},
		{
			c.RotateDown,
			"RotateDown",
		},
		{
			c.RotateLeft,
			"RotateLeft",
		},
		{
			c.RotateRight,
			"RotateRight",
		},
		{
			c.RotateClockwise,
			"RotateClockwise",
		},
		{
			c.RotateCounterClockwise,
			"RotateCounterClockwise",
		},
	}
	for _, x := range data {
		for i := 0; i < 4; i++ {
			x.fun(0)
		}
		if c.String() != startingState {
			t.Error("Failed running a rotation x4: ", x.name)
		}
		c.SetState(startingState)
	}
}

func TestClockwiseAndCounterClockwise(t *testing.T) {
	startingState := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQR"
	c := NewCube(startingState, 3)
	c.RotateClockwise(1)
	c.RotateCounterClockwise(1)
	if c.String() != startingState {
		t.Error("Failed counterclockwise and clockwise \n got: \t\t", c.String(), "\nexpected: \t", startingState)
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
			"000000000 111111111 222222222 333333333 123456789 123456789",
			"111111000 222222111 333333222 000000333 369258147 123456789",
			1,
			3,
		},
	}
	for _, x := range data {
		c := NewCube(strings.Replace(x.starting, " ", "", -1), x.size)
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
		size     int
	}{
		{
			"000000000 123456789 222222222 333333333 444444444 555555555",
			"550550550 369258147 244244244 333333333 004004004 225225225",
			1,
			3,
		},
		{
			"123456789 123456789 123456789 123456789 123456789 123456789",
			"123456789 369258147 127454781 123456789 123456789 923656389",
			0,
			3,
		},
	}
	for _, x := range data {
		c := NewCube(strings.Replace(x.starting, " ", "", -1), x.size)
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
		size     int
	}{
		{
			"000000000 123456789 222222222 333333333 444444444 555555555",
			"440440440 741852963 255255255 333333333 224224224 005005005",
			1,
			3,
		},
		{
			"123456789 123456789 123456789 123456789 123456789 123456789",
			"123456789 741852963 127454781 123456789 923656389 123456789",
			0,
			3,
		},
	}
	for _, x := range data {
		c := NewCube(strings.Replace(x.starting, " ", "", -1), x.size)
		c.RotateDown(x.columns)
		expected := strings.Replace(x.result, " ", "", -1)
		if c.String() != expected {
			t.Error("Failed RotateDown \n got: \t\t", c.String(), " \nexpected: \t", expected)
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
		c := NewCube(x.state, 3)
		if c.Solved() != x.result {
			t.Error("Failed Solved got: ", c.Solved(), "expected: ", x.result, "c: ", c.String())
		}
	}
}

func TestGetSide(t *testing.T) {
	data := []struct {
		state      string
		sideNumber int
		result     string
	}{
		{
			"000000000111111111222222222333333333444444444555555555",
			0,
			"000000000",
		},
		{
			"000000000111111111222222222333333333444444444555555555",
			5,
			"555555555",
		},
		{
			"000000000111111111222222222333333333444444444555555555",
			2,
			"222222222",
		},
	}
	for _, x := range data {
		c := NewCube(x.state, 3)
		got := c.getSide(x.sideNumber)
		if x.result != got {
			t.Error("Failed getSide got: ", got, " expected: ", x.result)
		}
	}
}

func TestReverseString(t *testing.T) {
	data := []struct {
		s      string
		result string
	}{
		{"12345", "54321"},
		{"a", "a"},
		{"", ""},
	}
	for _, x := range data {
		got := reverseString(x.s)
		if x.result != got {
			t.Error("Failed reverseString got: ", got, " expected: ", x.result)
		}
	}
}

func TestGetSideColumn(t *testing.T) {
	data := []struct {
		side   string
		col    int
		size   int
		result string
	}{
		{
			"123123123",
			2,
			3,
			"333",
		},
		{
			"123456789",
			1,
			3,
			"258",
		},
		{
			"123456789",
			0,
			3,
			"147",
		},
		{
			"1111122222333334444455555",
			4,
			5,
			"12345",
		},
	}
	for _, x := range data {
		got := getSideColumn(x.side, x.col, x.size)
		if got != x.result {
			t.Error("Failed getSideColumn got: ", got, "expected: ", x.result)
		}
	}
}

func TestSetSideColumn(t *testing.T) {
	data := []struct {
		side         string
		columnNumber int
		column       string
		size         int
		result       string
	}{
		{
			"111111111",
			2,
			"222",
			3,
			"112112112",
		},
		{
			"123456789",
			0,
			"000",
			3,
			"023056089",
		},
		{
			"1111122222333334444455555",
			3,
			"00000",
			5,
			"1110122202333034440455505",
		},
		{
			"1111222233334444",
			1,
			"7777",
			4,
			"1711272237334744",
		},
	}
	for _, x := range data {
		got := setSideColumn(x.side, x.columnNumber, x.column, x.size)
		if got != x.result {
			t.Error("Failed setSideColumn got: ", got, " expected: ", x.result)
		}
	}
}

func TestGetSideRow(t *testing.T) {
	data := []struct {
		side      string
		rowNumber int
		size      int
		result    string
	}{
		{
			"123456789",
			1,
			3,
			"456",
		},
		{
			"1111122222333334444455555",
			3,
			5,
			"44444",
		},
		{
			"11112222333344444",
			0,
			4,
			"1111",
		},
	}
	for _, x := range data {
		got := getSideRow(x.side, x.rowNumber, x.size)
		if got != x.result {
			t.Error("Failed getSideRow got: ", got, " expected: ", x.result)
		}
	}
}

func TestSetSideRow(t *testing.T) {
	data := []struct {
		side      string
		rowNumber int
		row       string
		size      int
		result    string
	}{
		{
			"123456789",
			1,
			"000",
			3,
			"123000789",
		},
		{
			"1111122222333334444455555",
			0,
			"88888",
			5,
			"8888822222333334444455555",
		},
		{
			"1234123412341234",
			3,
			"0000",
			4,
			"1234123412340000",
		},
	}
	for _, x := range data {
		got := setSideRow(x.side, x.rowNumber, x.row, x.size)
		if got != x.result {
			t.Error("Failed setSideRow got: ", got, " expected: ", x.result)
		}
	}
}

func TestRotateSideClockwise(t *testing.T) {
	data := []struct {
		side   string
		result string
		size   int
	}{
		{
			"123456789",
			"741852963",
			3,
		},
		{
			"1111122222333334444455555",
			"5432154321543215432154321",
			5,
		},
		{
			"1111222233334444",
			"4321432143214321",
			4,
		},
	}
	for _, x := range data {
		got := rotateSideClockwise(x.side, x.size)
		if got != x.result {
			t.Error("Failed rotateSideClockwise got: ", got, " expected: ", x.result)
		}
	}

}

func TestRotateSideCounterClockwise(t *testing.T) {
	data := []struct {
		side   string
		result string
		size   int
	}{
		{
			"123456789",
			"369258147",
			3,
		},
		{
			"1111122222333334444455555",
			"1234512345123451234512345",
			5,
		},
		{
			"1111222233334444",
			"1234123412341234",
			4,
		},
	}
	for _, x := range data {
		got := rotateSideCounterClockwise(x.side, x.size)
		if got != x.result {
			t.Error("Failed rotateSideCounterClockwise got: ", got, " expected: ", x.result)
		}
	}
}
