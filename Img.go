// Img.go
package main

import (
	"image"
)

type Img struct {
	img    *image.NRGBA //Pointer to image
	width  int          //Node width
	height int          //Node height
	color  []float64    //Average color stored as [R, G, B]
	error  float64      //Calculated error between average pixels and image
	c1     *Img         //Pointer to child 1
	c2     *Img         //Pointer to child 2
	c3     *Img         //Pointer to child 3
	c4     *Img         //Pointer to child 4
}
