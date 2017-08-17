// cli.go
package main

import (
	"flag"
)

func initializeFlags() (*string, *int, *bool, *int, *bool) {
	inFile := flag.String("f", "", "Input image name")
	iterations := flag.Int("i", 20, "Number of quad iterations to perform")
	toGif := flag.Bool("g", false, "Convert the intermediate images to a GIF")
	gifDelay := flag.Int("d", 20, "Delay between gif images")
	border := flag.Bool("b", false, "Adds 1px black border to quads")

	flag.Parse()

	return inFile, iterations, toGif, gifDelay, border
}
