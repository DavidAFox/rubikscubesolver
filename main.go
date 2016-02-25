package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/davidafox/rubikscubesolver/bytecube"
	"github.com/davidafox/rubikscubesolver/combined"
	"github.com/davidafox/rubikscubesolver/rubikscuberunner"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var depth = flag.Int("depth", 6, "specify the depth to use breadth-first seach")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	fmt.Println("Enter cube state with numbers 0-5 representing the colors starting with the Side facing you and going clockwise around the cube followed by the top and then the bottom. Type quit to quit.")
	fmt.Println("Example: 000000000111111111222222222333333333444444444555555555")
	scanner := bufio.NewScanner(os.Stdin)
	valid := false
	var err error
	var c *bytecube.Cube
	for !valid {
		scanner.Scan()
		state := scanner.Text()
		if state == "quit" {
			return
		}
		c, err = bytecube.NewCube(state)
		if err != nil {
			fmt.Println(err)
		} else {
			valid, err = c.Validate()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if c.Solved() {
		fmt.Println("The cube is already solved.")
		return
	}
	cf := combined.NewFactory()
	r := rubikscuberunner.NewOfficialRunner(c)
	//	r.Run("R U' B' L F R' U2 F2 L' D R U L'")
	s := combined.NewSolver(c.String(), cf, *depth)
	startTime := time.Now()
	solution := s.Solve()
	runtime := time.Since(startTime)
	r.Run(solution)
	fmt.Println(solution)
	fmt.Println("Time: ", runtime)
	fmt.Println("Solved: ", c.Solved())

}
