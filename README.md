Rubikscubesolver is for solving 3x3 Rubik's cubes.

#### Use
To solve a cube start the program and then enter the cube in the form "000000000111111111222222222333333333444444444555555555" where each number represents one cube face.

Each side is expressed starting in the top left, from right to left and top to bottom.  Start with the side facing you and work around clockwise then the top side followed by the bottom.  The numbers represent each of the six colors found on the cube.  The example above represents a solved cube.

The result (if it finishes) will be a series of moves to transform the cube from its starting state to the solved state.  The letter represents the side of the cube to rotate.
* R - the right side of the cube
* L - the left side of the cube
* U - the top side of the cube
* D - the bottom side of the cube
* F - the front side of the cube
* B - the back side of the cube

A letter on its own means a clockwise turn and a letter followed by a ' means a counterclockwise turn.  A letter followed by a 2 means to turn that side clockwise 180 or a double turn.

#### Runtime
It uses a combination of a breadth first search out a depth specified by a -depth _number_ flag from the starting state and the solution  
If no depth is specified a default of 6 is used.  Greater depths will generally run faster but consume more memory.

The program doesn't use any heuristics and as a result will take a very long time (essentially forever) to solve cubes taking too many steps.

The limit is probably around 15 steps depending on hardware and how long you're willing to wait.

Some of the other packages include earlier implementations of the solver or of the cube.