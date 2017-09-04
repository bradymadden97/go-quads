// main.go
package main

import (
	"container/heap"
	"log"
	"os"
)

const outputFolder string = "./out/"

type Img struct {
	hist   [][]int   //Histogram of image stored as [R, G, B]
	pix    int       //Number of pixels in image
	color  []float64 //Average color stored as [R, G, B]
	error  float64   //Calculated error between average pixels and image
	width  int       //Picture width
	height int       //Picture height
	c1     *Img      //Pointer to child 1
	c2     *Img      //Pointer to child 2
	c3     *Img      //Pointer to child 3
	c4     *Img      //Pointer to child 4
}

func main() {
	flags := initializeFlags()
	if *flags.f == "" {
		log.Fatal(" -f <input image> required")
	}

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err := os.Mkdir(outputFolder, os.ModeDir); err != nil {
			log.Fatal(err)
		}
	}

	headNode, err := initialize(*flags.f)
	if err != nil {
		log.Fatal(err)
	}

	mh := make(MinHeap, 1)
	mh[0] = headNode
	heap.Init(&mh)

	imgs, err_itr := iterate(&mh, headNode, *flags.i, *flags.f, *flags.b, *flags.c, *flags.bc, *flags.s, *flags.g)
	if err_itr != nil {
		log.Fatal(err_itr)
	}

	if *flags.g {
		err = toGIF(imgs, *flags.f, *flags.gd, *flags.gp)
		if err != nil {
			log.Fatal(err)
		}
	}
}
