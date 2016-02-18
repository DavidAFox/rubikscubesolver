package rubikscuberunner

import (
	"fmt"
	"strconv"
	"strings"
)

//rubikscuberunner interprets a string of commands in letter number pairs
//C==RotateClockwise()
//T==RotateCounterCLockwise()
//U==RotateUp()
//D==RotateDown()
//R==RotateRight()
//L==RotateLeft()
//So C1 would result in RotateClockwise(1)
//and R3 would result in RotateCounterClockwise(3).

type cube interface {
	RotateClockwise(row int)
	RotateCounterClockwise(row int)
	RotateUp(col int)
	RotateDown(col int)
	RotateRight(row int)
	RotateLeft(row int)
}

type officialCube interface {
	RotateR()
	RotateRCounter()
	RotateL()
	RotateLCounter()
	RotateU()
	RotateUCounter()
	RotateD()
	RotateDCounter()
	RotateF()
	RotateFCounter()
	RotateB()
	RotateBCounter()
}

type Runner struct {
	c cube
}

type OfficialRunner struct {
	c officialCube
}

func NewOfficialRunner(c officialCube) *OfficialRunner {
	r := new(OfficialRunner)
	r.c = c
	return r
}

func NewRunner(c cube) *Runner {
	r := new(Runner)
	r.c = c
	return r
}

func (r *OfficialRunner) Run(s string) {
	s = strings.TrimSpace(s)
	steps := strings.Split(s, " ")
	for _, step := range steps {
		switch string(step[0]) {
		case "R":
			if len(step) < 2 {
				r.c.RotateR()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateRCounter()
				} else {
					r.c.RotateR()
					r.c.RotateR()
				}
			}
		case "L":
			if len(step) < 2 {
				r.c.RotateL()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateLCounter()
				} else {
					r.c.RotateL()
					r.c.RotateL()
				}
			}
		case "U":
			if len(step) < 2 {
				r.c.RotateU()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateUCounter()
				} else {
					r.c.RotateU()
					r.c.RotateU()
				}
			}
		case "D":
			if len(step) < 2 {
				r.c.RotateD()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateDCounter()
				} else {
					r.c.RotateD()
					r.c.RotateD()
				}
			}
		case "F":
			if len(step) < 2 {
				r.c.RotateF()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateFCounter()
				} else {
					r.c.RotateF()
					r.c.RotateF()
				}
			}
		case "B":
			if len(step) < 2 {
				r.c.RotateB()
			} else {
				if string(step[1]) == "'" {
					r.c.RotateBCounter()
				} else {
					r.c.RotateB()
					r.c.RotateB()
				}
			}
		default:
			fmt.Println("Invalid step: ", s)
		}
	}
}

func (r *Runner) Run(s string) {
	for i := 0; i+1 < len(s); i += 2 {
		x, err := strconv.Atoi(string(s[i+1]))
		if err != nil {
			fmt.Println("Invalid string: ", s, " at ", string(s[i+1]))
			return
		}
		switch string(s[i]) {
		case "C":
			r.c.RotateClockwise(x)
		case "T":
			r.c.RotateCounterClockwise(x)
		case "U":
			r.c.RotateUp(x)
		case "D":
			r.c.RotateDown(x)
		case "R":
			r.c.RotateRight(x)
		case "L":
			r.c.RotateLeft(x)
		case "r":
			r.c.RotateRight(x)
			r.c.RotateRight(x)
		case "u":
			r.c.RotateUp(x)
			r.c.RotateUp(x)
		case "c":
			r.c.RotateClockwise(x)
			r.c.RotateClockwise(x)
		default:
			fmt.Println("Invalid string: ", s, " at ", string(s[i]))
		}
	}
}
