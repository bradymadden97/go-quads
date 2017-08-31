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
	gf *int    //Gif frames per second
	gp *int    //Gif pause before repeat
	ds *bool   //Don't save intermediate images
	c  *bool   //Modify quads to circles
}

func initializeFlags() *Flags {
	flags := Flags{
		f:  flag.String("f", "", "Input image name"),
		i:  flag.Int("i", 20, "Number of quad iterations to perform"),
		b:  flag.Bool("b", false, "Adds 1px black border to quads"),
		bc: flag.String("bc", "0,0,0", "Border/ background color between quads"),
		g:  flag.Bool("g", false, "Convert the intermediate images to a GIF"),
		gf: flag.Int("gf", 20, "Frames per second for GIF"),
		gp: flag.Int("gp", 2, "Pause in seconds at end of GIF loop"),
		ds: flag.Bool("ds", false, "Don't save subimages"),
		c:  flag.Bool("c", false, "Modify quads to circles"),
	}
	flag.Parse()

	return &flags
}
