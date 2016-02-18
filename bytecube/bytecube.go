package bytecube

import (
	"errors"
	"log"
	"strconv"
)

//State is a "comparable" form of the cube's state for use in maps in the solver packages
type State struct {
	zero  uint32
	one   uint32
	two   uint32
	three uint32
	four  uint32
	five  uint32
}

type Cube struct {
	sides [6]uint32
}

var indexValue = [...]uint32{7, 56, 448, 3584, 28672, 229376, 1835008, 14680064, 117440512}

const MAX uint32 = 4294967295

var ErrIncorrectNumber = errors.New("The cube has an incorrect number of cubie sides")
var ErrCenterCubies = errors.New("Center cubies must be different colors")
var ErrIncorrectColorNumbers = errors.New("There must be 9 of each color")
var ErrIncorrectCorners = errors.New("The corner cubies are incorrect")
var ErrIncorrectSides = errors.New("The side cubies are incorrect")

func (c *Cube) setLocation(side, spot, intValue int) {
	if spot > 8 || spot < 0 {
		log.Println("Invalid spot ", spot)
		return
	}
	if side < 0 || side > 5 {
		log.Println("Invalid side ", side)
		return
	}
	if intValue < 0 || intValue > 5 {
		log.Println("Invalid value ", intValue)
	}
	value := uint32(intValue)
	mask := MAX - indexValue[spot]
	newValue := c.sides[side] & mask
	value = value << (3 * uint32(spot))
	newValue += value
	c.sides[side] = newValue
}

func (c *Cube) getLocation(side, spot int) int {
	if spot > 8 || spot < 0 {
		log.Println("Invalid spot ", spot)
		return 0
	}
	if side < 0 || side > 5 {
		log.Println("Invalid side ", side)
		return 0
	}
	v := indexValue[spot] & c.sides[side]
	v = v >> (3 * uint32(spot))
	return int(v)
}

func (c *Cube) String() string {
	result := ""
	for i := 0; i < 6; i++ {
		result += sideString(c.sides[i])
	}
	return result
}

func NewCube(s string) (*Cube, error) {
	if len(s) != 54 {
		return nil, ErrIncorrectNumber
	}
	c := new(Cube)
	c.sides = [6]uint32{}
	for i := 0; i < 6; i++ {
		for j := 0; j < 9; j++ {
			v, err := strconv.ParseUint(string(s[(i*9)+j]), 10, 32)

			if err != nil {
				log.Println("Invalid string: ", s, " at ", string(s[(i*9)+j]))
			}
			v = v << (uint(j) * 3)
			c.sides[i] += uint32(v)
		}
	}
	return c, nil
}

func (c *Cube) SolvedState() string {
	result := ""
	for i := 0; i < 6; i++ {
		for j := 0; j < 9; j++ {
			result += strconv.Itoa(c.getLocation(i, 4))
		}
	}
	return result
}

func (c *Cube) Validate() (bool, error) {
	if !checkCenters(c) {
		return false, ErrCenterCubies
	}
	if !checkColors(c) {
		return false, ErrIncorrectColorNumbers
	}
	if !checkCorners(c) {
		return false, ErrIncorrectCorners
	}
	if !checkSides(c) {
		return false, ErrIncorrectSides
	}
	return true, nil
}

func checkCenters(c *Cube) bool {
	colorsFound := make(map[int]bool)
	for i := range c.sides {
		value := c.getLocation(i, 4)
		if value > 5 || value < 0 {
			return false
		}
		if colorsFound[value] {
			return false
		}
		colorsFound[value] = true
	}
	return true
}

func checkColors(c *Cube) bool {
	colors := make(map[int]int)
	for i := range c.sides {
		for j := 0; j < 9; j++ {
			value := c.getLocation(i, j)
			if value > 5 || value < 0 {
				return false
			}
			colors[value]++
		}
	}
	for i := 0; i < 6; i++ {
		if colors[i] != 9 {
			return false
		}
	}
	return true
}

func checkCorners(c *Cube) bool {
	opposites := getOpposites(c)
	corners := getCorners(c)
	if !checkCornerColors(corners) {
		return false
	}
	for _, corner := range corners {
		if !checkCornerDifferent(corner) || !checkCornerNotOpposite(corner, opposites) {
			return false
		}
	}
	return true
}

func getOpposites(c *Cube) map[int]int {
	result := make(map[int]int)
	result[c.getLocation(0, 4)] = c.getLocation(2, 4)
	result[c.getLocation(2, 4)] = c.getLocation(0, 4)
	result[c.getLocation(1, 4)] = c.getLocation(3, 4)
	result[c.getLocation(3, 4)] = c.getLocation(1, 4)
	result[c.getLocation(4, 4)] = c.getLocation(5, 4)
	result[c.getLocation(5, 4)] = c.getLocation(4, 4)
	return result
}

func checkCornerColors(corners []map[int]bool) bool {
	colors := make(map[int]int)
	for _, corner := range corners {
		for key := range corner {
			colors[key]++
		}
	}
	for i := 0; i < 6; i++ {
		if colors[i] != 4 {
			return false
		}
	}
	return true
}

func checkCornerDifferent(corner map[int]bool) bool {
	colors := make(map[int]int)
	for color := range corner {
		colors[color]++
	}
	for _, x := range colors {
		if x != 1 {
			return false
		}
	}
	return true
}

func checkCornerNotOpposite(corner map[int]bool, opposites map[int]int) bool {
	colors := make([]int, 3)
	i := 0
	for x := range corner {
		colors[i] = x
		i++
	}
	if colors[1] == opposites[colors[0]] || colors[2] == opposites[colors[0]] || colors[2] == opposites[colors[1]] {
		return false
	}
	return true
}

func getCorners(c *Cube) []map[int]bool {
	result := make([]map[int]bool, 8)
	result[0] = corner(c, []int{0, 0}, []int{1, 2}, []int{4, 6})
	result[1] = corner(c, []int{0, 6}, []int{1, 8}, []int{5, 0})
	result[2] = corner(c, []int{0, 2}, []int{4, 8}, []int{3, 0})
	result[3] = corner(c, []int{0, 8}, []int{5, 2}, []int{3, 6})
	result[4] = corner(c, []int{2, 0}, []int{3, 2}, []int{4, 2})
	result[5] = corner(c, []int{2, 2}, []int{4, 0}, []int{1, 0})
	result[6] = corner(c, []int{2, 6}, []int{3, 8}, []int{5, 8})
	result[7] = corner(c, []int{2, 8}, []int{1, 6}, []int{5, 6})
	return result
}

func corner(c *Cube, side1, side2, side3 []int) map[int]bool {
	result := make(map[int]bool)
	result[c.getLocation(side1[0], side1[1])] = true
	result[c.getLocation(side2[0], side2[1])] = true
	result[c.getLocation(side3[0], side3[1])] = true
	return result
}

func checkSides(c *Cube) bool {
	opposites := getOpposites(c)
	sides := getSides(c)
	if !checkSideColors(sides) {
		return false
	}
	for _, side := range sides {
		if !checkSideNotOpposite(side, opposites) || !checkSideNotSame(side) {
			return false
		}
	}
	return true
}

func getSides(c *Cube) [][2]int {
	result := make([][2]int, 12)
	result[0] = [2]int{c.getLocation(0, 1), c.getLocation(4, 7)}
	result[1] = [2]int{c.getLocation(0, 5), c.getLocation(3, 3)}
	result[2] = [2]int{c.getLocation(0, 7), c.getLocation(5, 1)}
	result[3] = [2]int{c.getLocation(0, 3), c.getLocation(1, 6)}
	result[4] = [2]int{c.getLocation(1, 1), c.getLocation(4, 3)}
	result[5] = [2]int{c.getLocation(4, 5), c.getLocation(3, 1)}
	result[6] = [2]int{c.getLocation(3, 7), c.getLocation(5, 5)}
	result[7] = [2]int{c.getLocation(5, 3), c.getLocation(1, 7)}
	result[8] = [2]int{c.getLocation(2, 1), c.getLocation(4, 1)}
	result[9] = [2]int{c.getLocation(2, 5), c.getLocation(1, 3)}
	result[10] = [2]int{c.getLocation(2, 7), c.getLocation(5, 7)}
	result[11] = [2]int{c.getLocation(2, 3), c.getLocation(3, 5)}
	return result
}

func checkSideNotOpposite(side [2]int, opposites map[int]int) bool {
	return side[0] != opposites[side[1]]
}

func checkSideNotSame(side [2]int) bool {
	return side[0] != side[1]
}

func checkSideColors(sides [][2]int) bool {
	colors := make(map[int]int)
	for _, side := range sides {
		for _, color := range side {
			colors[color]++
		}
	}
	for i := 0; i < 6; i++ {
		if colors[i] != 4 {
			return false
		}
	}
	return true
}

func sideString(side uint32) string {
	result := ""
	for i := 0; i < 9; i++ {
		v := side & indexValue[i]
		v = v >> (uint(i) * 3)
		result += strconv.Itoa(int(v))
	}
	return result
}

func (c *Cube) getRow(side, rowNumber int) []int {
	results := make([]int, 3)
	for i := 0; i < 3; i++ {
		results[i] = c.getLocation(side, rowNumber*3+i)
	}
	return results
}

func (c *Cube) setRow(side, rowNumber int, row []int) {
	for i := 0; i < 3; i++ {
		c.setLocation(side, rowNumber*3+i, row[i])
	}
}

func (c *Cube) getColumn(side, columnNumber int) []int {
	results := make([]int, 3)
	for i := 0; i < 3; i++ {
		results[i] = c.getLocation(side, 3*i+columnNumber)
	}
	return results
}

func (c *Cube) setColumn(side, columnNumber int, column []int) {
	for i := 0; i < 3; i++ {
		c.setLocation(side, 3*i+columnNumber, column[i])
	}
}

func (c *Cube) rotateRowRight(row int) {
	last := c.getRow(3, row)
	for i := 3; i > 0; i-- {
		c.setRow(i, row, c.getRow(i-1, row))
	}
	c.setRow(0, row, last)
}
func (c *Cube) rotateRowLeft(row int) {
	last := c.getRow(0, row)
	for i := 0; i < 3; i++ {
		c.setRow(i, row, c.getRow(i+1, row))
	}
	c.setRow(3, row, last)
}

func (c *Cube) rotateColumnUp(col int) {
	last := c.getColumn(0, col)
	c.setColumn(0, col, c.getColumn(5, col))
	c.setColumn(5, col, reverse(c.getColumn(2, 2-col)))
	c.setColumn(2, 2-col, reverse(c.getColumn(4, col)))
	c.setColumn(4, col, last)
}
func (c *Cube) rotateColumnDown(col int) {
	last := c.getColumn(5, col)
	c.setColumn(5, col, c.getColumn(0, col))
	c.setColumn(0, col, c.getColumn(4, col))
	c.setColumn(4, col, reverse(c.getColumn(2, 2-col)))
	c.setColumn(2, 2-col, reverse(last))
}

func (c *Cube) rotateLevelClockwise(lvl int) {
	last := c.getColumn(1, 2-lvl)
	c.setColumn(1, 2-lvl, c.getRow(5, lvl))
	c.setRow(5, lvl, reverse(c.getColumn(3, lvl)))
	c.setColumn(3, lvl, c.getRow(4, 2-lvl))
	c.setRow(4, 2-lvl, reverse(last))
}

func (c *Cube) rotateLevelCounterClockwise(lvl int) {
	last := c.getColumn(1, 2-lvl)
	c.setColumn(1, 2-lvl, reverse(c.getRow(4, 2-lvl)))
	c.setRow(4, 2-lvl, c.getColumn(3, lvl))
	c.setColumn(3, lvl, reverse(c.getRow(5, lvl)))
	c.setRow(5, lvl, last)
}

func (c *Cube) rotateSideClockwise(side int) {
	tmp := c.getRow(side, 0)
	c.setRow(side, 0, reverse(c.getColumn(side, 0)))
	c.setColumn(side, 0, c.getRow(side, 2))
	//end of the bottom row will be wrong but gets fixed when the last column is set based on the stored value
	c.setRow(side, 2, reverse(c.getColumn(side, 2)))
	c.setColumn(side, 2, tmp)
}

func (c *Cube) rotateSideCounterClockwise(side int) {
	tmp := c.getRow(side, 0)
	c.setRow(side, 0, c.getColumn(side, 2))
	c.setColumn(side, 2, reverse(c.getRow(side, 2)))
	c.setRow(side, 2, c.getColumn(side, 0))
	c.setColumn(side, 0, reverse(tmp))
}

func reverse(start []int) []int {
	start[0], start[2] = start[2], start[0]
	return start
}

func (c *Cube) RotateClockwise(lvls int) {
	for i := 0; i <= lvls; i++ {
		c.rotateLevelClockwise(i)
	}
	c.rotateSideClockwise(0)
}

func (c *Cube) RotateCounterClockwise(lvls int) {
	for i := 0; i <= lvls; i++ {
		c.rotateLevelCounterClockwise(i)
	}
	c.rotateSideCounterClockwise(0)
}

func (c *Cube) RotateRight(rows int) {
	for i := 0; i <= rows; i++ {
		c.rotateRowRight(i)
	}
	c.rotateSideClockwise(4)
}

func (c *Cube) RotateLeft(rows int) {
	for i := 0; i <= rows; i++ {
		c.rotateRowLeft(i)
	}
	c.rotateSideCounterClockwise(4)
}

func (c *Cube) RotateUp(columns int) {
	for i := 0; i <= columns; i++ {
		c.rotateColumnUp(i)
	}
	c.rotateSideCounterClockwise(1)
}

func (c *Cube) RotateDown(columns int) {
	for i := 0; i <= columns; i++ {
		c.rotateColumnDown(i)
	}
	c.rotateSideClockwise(1)
}

//Separate functions for "official" notation

func (c *Cube) RotateR() {
	c.rotateSideClockwise(3)
	c.rotateColumnUp(2)
}

func (c *Cube) RotateRCounter() {
	c.rotateSideCounterClockwise(3)
	c.rotateColumnDown(2)
}

func (c *Cube) RotateL() {
	c.rotateSideClockwise(1)
	c.rotateColumnDown(0)
}

func (c *Cube) RotateLCounter() {
	c.rotateSideCounterClockwise(1)
	c.rotateColumnUp(0)
}

func (c *Cube) RotateU() {
	c.rotateSideClockwise(4)
	c.rotateRowRight(0)
}

func (c *Cube) RotateUCounter() {
	c.rotateSideCounterClockwise(4)
	c.rotateRowLeft(0)
}

func (c *Cube) RotateD() {
	c.rotateSideClockwise(5)
	c.rotateRowLeft(2)
}

func (c *Cube) RotateDCounter() {
	c.rotateSideCounterClockwise(5)
	c.rotateRowRight(2)
}

func (c *Cube) RotateF() {
	c.rotateSideClockwise(0)
	c.rotateLevelClockwise(0)
}

func (c *Cube) RotateFCounter() {
	c.rotateSideCounterClockwise(0)
	c.rotateLevelCounterClockwise(0)
}

func (c *Cube) RotateB() {
	c.rotateSideClockwise(2)
	c.rotateLevelCounterClockwise(2)
}

func (c *Cube) RotateBCounter() {
	c.rotateSideCounterClockwise(2)
	c.rotateLevelClockwise(2)
}

func (c *Cube) Solved() bool {
	for i := 0; i < 6; i++ {
		x := c.getLocation(i, 0)
		for j := 0; j < 9; j++ {
			if c.getLocation(i, j) != x {
				return false
			}
		}
	}
	return true
}

func (c *Cube) State() State {
	s := State{}
	s.zero = c.sides[0]
	s.one = c.sides[1]
	s.two = c.sides[2]
	s.three = c.sides[3]
	s.four = c.sides[4]
	s.five = c.sides[5]
	return s
}

func NewWithState(s State) *Cube {
	c := new(Cube)
	c.sides = [6]uint32{}
	c.sides[0] = s.zero
	c.sides[1] = s.one
	c.sides[2] = s.two
	c.sides[3] = s.three
	c.sides[4] = s.four
	c.sides[5] = s.five
	return c
}
