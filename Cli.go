// Cli
package main

import (
	"flag"
)

func initializeFlags() (*string, *int, *bool, *int) {
	inFile := flag.String("f", "", "Input image name")
	iterations := flag.Int("i", 20, "Number of quad iterations to perform")
	toGif := flag.Bool("g", false, "Convert the intermediate images to a GIF")
	gifDelay := flag.Int("d", 20, "Delay between gif images")

	flag.Parse()

	return inFile, iterations, toGif, gifDelay
}
