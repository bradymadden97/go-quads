// cli.go
package main

import (
	"flag"
)

func initializeFlags() (*string, *int, *bool, *int, *int, *bool) {
	inFile := flag.String("n", "", "Input image name")
	iterations := flag.Int("i", 20, "Number of quad iterations to perform")
	toGif := flag.Bool("g", false, "Convert the intermediate images to a GIF")
	gifDelay := flag.Int("gf", 20, "Frames per second for GIF")
	gifPause := flag.Int("gp", 2, "Pause in seconds at end of GIF loop")
	border := flag.Bool("b", false, "Adds 1px black border to quads")

	flag.Parse()

	return inFile, iterations, toGif, gifDelay, gifPause, border
}
