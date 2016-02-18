package rubikscube

import (
	"fmt"
	"math"
	"strconv"
)

type Cube struct {
	state  string
	size   int
	sizeSq int
}

func (c *Cube) String() string {
	return c.state
}

func (c *Cube) SetState(s string) {
	c.state = s
}

func NewCube(state string, size int) *Cube {
	c := new(Cube)
	c.state = state
	c.size = size
	c.sizeSq = int(math.Pow(float64(c.size), 2))
	return c
}

func NewCubeSlice(slice [][][]int, size int) *Cube {
	c := new(Cube)
	c.state = convertToString(slice)
	c.size = size
	c.sizeSq = int(math.Pow(float64(c.size), 2))
	return c
}

func (c *Cube) print() {
	fmt.Println(c.state)
}

func (c *Cube) getRow(side, rowNumber int) string {
	starting := c.sizeSq*side + rowNumber*c.size
	return c.state[starting : starting+c.size]
}

func (c *Cube) setRow(side, rowNumber int, row string) {
	if len(row) != c.size {
		fmt.Println("Invalid Row: ", row)
		return
	}
	starting := c.sizeSq*side + rowNumber*c.size
	c.state = c.state[:starting] + row + c.state[starting+len(row):]
}

func (c *Cube) getColumn(side, columnNumber int) string {
	result := ""
	starting := c.sizeSq * side
	for i := 0; i < c.size; i++ {
		result += string(c.state[starting+(i*c.size)+columnNumber])
	}
	return result
}

func (c *Cube) setColumn(side, columnNumber int, column string) {
	starting := c.sizeSq * side
	result := c.state[:starting]
	for i := 0; i < c.size; i++ {
		if columnNumber > 0 {
			result += c.state[starting+(i*c.size) : starting+(i*c.size)+columnNumber]
		}
		result += string(column[i])
		if columnNumber < c.size {
			result += c.state[starting+(i*c.size)+columnNumber+1 : starting+(i*c.size)+c.size]
		}
	}
	result += c.state[starting+c.sizeSq : len(c.state)]
	c.state = result
}

func (c *Cube) rotateRowRight(row int) {
	last := c.getRow(3, row)
	for i := 3; i > 0; i-- {
		c.setRow(i, row, c.getRow(i-1, row))
	}
	c.setRow(0, row, last)
}

func (c *Cube) RotateRight(rows int) {
	for i := 0; i <= rows; i++ {
		c.rotateRowRight(i)
	}
	newSide4 := rotateSideClockwise(c.getSide(4), c.size)
	c.state = c.state[:c.sizeSq*4] + newSide4 + c.getSide(5)
}

func (c *Cube) rotateRowLeft(row int) {
	last := c.getRow(0, row)
	for i := 0; i < 3; i++ {
		c.setRow(i, row, c.getRow(i+1, row))
	}
	c.setRow(3, row, last)
}

func (c *Cube) RotateLeft(rows int) {
	for i := 0; i <= rows; i++ {
		c.rotateRowLeft(i)
	}
	newSide4 := rotateSideCounterClockwise(c.getSide(4), c.size)
	c.state = c.state[:c.sizeSq*4] + newSide4 + c.getSide(5)
}

func (c *Cube) rotateColumnUp(col int) {
	last := c.getColumn(0, col)
	c.setColumn(0, col, c.getColumn(5, col))
	c.setColumn(5, col, reverseString(c.getColumn(2, c.size-1-col)))
	c.setColumn(2, c.size-1-col, reverseString(c.getColumn(4, col)))
	c.setColumn(4, col, last)
}

func (c *Cube) RotateUp(columns int) {
	for i := 0; i <= columns; i++ {
		c.rotateColumnUp(i)
	}
	newSide1 := rotateSideCounterClockwise(c.getSide(1), c.size)
	c.state = c.state[:c.sizeSq*1] + newSide1 + c.state[c.sizeSq*2:]
}

func (c *Cube) rotateColumnDown(col int) {
	last := c.getColumn(5, col)
	c.setColumn(5, col, c.getColumn(0, col))
	c.setColumn(0, col, c.getColumn(4, col))
	c.setColumn(4, col, reverseString(c.getColumn(2, c.size-1-col)))
	c.setColumn(2, c.size-1-col, reverseString(last))
}

func (c *Cube) RotateDown(columns int) {
	for i := 0; i <= columns; i++ {
		c.rotateColumnDown(i)
	}
	newSide1 := rotateSideClockwise(c.getSide(1), c.size)
	c.state = c.state[:c.sizeSq*1] + newSide1 + c.state[c.sizeSq*2:]
}

func (c *Cube) RotateClockwise(level int) {
	side := make([]string, 6)
	side[0] = rotateSideClockwise(c.getSide(0), c.size)
	for i := 1; i < 6; i++ {
		side[i] = c.getSide(i)
	}
	for i := 0; i < level+1; i++ {
		side[1] = setSideColumn(side[1], (c.size-1)-i, c.getRow(5, i), c.size)
		side[3] = setSideColumn(side[3], i, c.getRow(4, c.size-1-i), c.size)
		side[4] = setSideRow(side[4], (c.size-1)-i, reverseString(c.getColumn(1, c.size-1-i)), c.size)
		side[5] = setSideRow(side[5], i, reverseString(c.getColumn(3, i)), c.size)
	}
	c.state = ""
	for _, x := range side {
		c.state += x
	}
}
func (c *Cube) getSide(n int) string {
	return c.state[c.sizeSq*n : c.sizeSq*(n+1)]
}
func (c *Cube) RotateCounterClockwise(level int) {
	side := make([]string, 6)
	side[0] = rotateSideCounterClockwise(c.getSide(0), c.size)
	for i := 1; i < 6; i++ {
		side[i] = c.getSide(i)
	}
	for i := 0; i < level+1; i++ {
		side[1] = setSideColumn(side[1], (c.size-1)-i, reverseString(c.getRow(4, c.size-1-i)), c.size)
		side[3] = setSideColumn(side[3], i, reverseString(c.getRow(5, i)), c.size)
		side[4] = setSideRow(side[4], (c.size-1)-i, c.getColumn(3, i), c.size)
		side[5] = setSideRow(side[5], i, c.getColumn(1, c.size-1-i), c.size)
	}
	c.state = ""
	for _, x := range side {
		c.state += x
	}
}

func rotateSideClockwise(side string, size int) string {
	for i := 0; i < size/2; i++ {
		tmp := getSideRow(side, i, size)
		side = setSideRow(side, i, getSideRow(side, i, size)[0:i]+reverseString(getSideColumn(side, i, size)[i:size-i])+getSideRow(side, i, size)[size-i:size], size)
		side = setSideColumn(side, i, getSideColumn(side, i, size)[0:i]+getSideRow(side, size-i-1, size)[i:size-i]+getSideColumn(side, i, size)[size-i:size], size)
		//end of the bottom row will be wrong but gets fixed when the last column is set based on the stored value
		side = setSideRow(side, size-i-1, getSideRow(side, size-i-1, size)[0:i]+reverseString(getSideColumn(side, size-i-1, size)[i:size-i])+getSideRow(side, size-i-1, size)[size-i:size], size)
		side = setSideColumn(side, size-i-1, getSideColumn(side, size-i-1, size)[0:i]+tmp[i:size-i]+getSideColumn(side, size-i-1, size)[size-i:size], size)
	}
	return side
}
func rotateSideCounterClockwise(side string, size int) string {
	for i := 0; i < size/2; i++ {
		tmp := getSideRow(side, i, size)
		side = setSideRow(side, i, getSideRow(side, i, size)[0:i]+getSideColumn(side, size-i-1, size)[i:size-i]+getSideRow(side, i, size)[size-i:size], size)
		side = setSideColumn(side, size-i-1, getSideColumn(side, size-i-1, size)[0:i]+reverseString(getSideRow(side, size-i-1, size)[i:size-i])+getSideColumn(side, size-i-1, size)[size-i:size], size)
		side = setSideRow(side, size-i-1, getSideRow(side, size-i-1, size)[0:i]+getSideColumn(side, i, size)[i:size-i]+getSideRow(side, size-i-1, size)[size-i:size], size)
		side = setSideColumn(side, i, getSideColumn(side, i, size)[0:i]+reverseString(tmp[i:size-i])+getSideColumn(side, i, size)[size-i:size], size)
	}
	return side
}

func getSideColumn(side string, columnNumber, size int) string {
	result := ""
	for i := 0; i < size; i++ {
		result += string(side[(i*size)+columnNumber])
	}
	return result
}

func setSideColumn(side string, columnNumber int, column string, size int) string {
	result := ""
	for i := 0; i < size; i++ {
		if columnNumber > 0 {
			result += side[(i * size) : (i*size)+columnNumber]
		}
		result += string(column[i])
		if columnNumber < size {
			result += side[(i*size)+columnNumber+1 : (i*size)+size]
		}
	}
	return result
}

func getSideRow(side string, rowNumber, size int) string {
	starting := rowNumber * size
	return side[starting : starting+size]

}

func setSideRow(side string, rowNumber int, row string, size int) string {
	if len(row) != size {
		fmt.Println("Invalid Row: ", row)
		return ""
	}
	starting := rowNumber * size
	result := side[:starting] + row + side[starting+len(row):]
	return result
}

func reverseString(s string) string {
	r := ""
	for i := len(s) - 1; i >= 0; i-- {
		r += string(s[i])
	}
	return r
}

func (c *Cube) Solved() bool {
	for i := 0; i < 6; i++ {
		x := c.state[i*c.sizeSq]
		for j := 0; j < c.sizeSq; j++ {
			if c.state[i*c.sizeSq+j] != x {
				return false
			}
		}
	}
	return true
}

func convertToString(slices [][][]int) string {
	result := ""
	for _, side := range slices {
		for _, row := range side {
			for _, square := range row {
				result += strconv.Itoa(square)
			}
		}
	}
	return result
}
