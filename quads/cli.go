// cli.go
package main

import (
	"flag"
)

type Flags struct {
	f  *string //Input filename
	i  *int    //Iterations
	b  *bool   //Borders
	bc *string //Border/background color
	g  *bool   //to Gif
	gd *int    //Gif delay per frame in 100th of a second
	gp *int    //Gif pause before repeat
	s  *bool   //Save intermediate images
	c  *bool   //Modify quads to circles
}

func initializeFlags() *Flags {
	flags := Flags{
		f:  flag.String("f", "", "Input image name"),
		i:  flag.Int("i", 200, "Number of quad iterations to perform"),
		b:  flag.Bool("b", false, "Adds 1px black border to quads"),
		bc: flag.String("bc", "0,0,0,255", "Border/ background color between quads"),
		g:  flag.Bool("g", false, "Convert the intermediate images to a GIF"),
		gd: flag.Int("gd", 5, "Delay per frame in GIF in 100th of a second"),
		gp: flag.Int("gp", 2, "Pause in seconds at end of GIF loop"),
		s:  flag.Bool("s", false, "Save subimages"),
		c:  flag.Bool("c", false, "Modify quads to circles"),
	}
	flag.Parse()

	return &flags
}
