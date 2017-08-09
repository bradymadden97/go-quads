// Main.go

/*
** Using quadtree to keep track of order of images to piece back together recursively
** Using minheap to quickly push new subimages and pop highest-error image
 */

package main

import (
	"container/heap"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/disintegration/imaging"
)

type Img struct {
	img   *image.NRGBA //Pointer to image
	color []float64    //Average color stored as [R, G, B]
	error float64      //Calculated error between average pixels and image
	c1    *Img         //Pointer to child 1
	c2    *Img         //Pointer to child 2
	c3    *Img         //Pointer to child 3
	c4    *Img         //Pointer to child 4
}

type MinHeap []*Img

func main() {
	//Get input image and convert to NRGBA
	input_img, err := imaging.Open("tiger.jpg")
	if err != nil {
		log.Fatalf("Open image failed: %v", err)
	}
	img, err := toNRGBA(input_img)
	if err != nil {
		log.Fatalf("Image error: %v", err)
	}

	//Create Img object, get initial error, start Quadtree
	headNode := Img{
		img: img,
	}
	analyzeImage(&headNode)

	//Begin minheap
	mh := make(MinHeap, 1)
	mh[0] = &headNode
	heap.Init(&mh)

	//Loop
	iterations := 5
	for i := 0; i < iterations; i++ {
		a := mh.Pop().(*Img)
		c := splitImage(a.img)

		a.img = nil
		a.c1, a.c2, a.c3, a.c4 = c[0], c[1], c[2], c[3]

		heap.Push(&mh, a.c1)
		heap.Push(&mh, a.c2)
		heap.Push(&mh, a.c3)
		heap.Push(&mh, a.c4)
	}
}

func (mh MinHeap) Len() int { return len(mh) }

func (mh MinHeap) Less(i, j int) bool {
	return mh[i].error > mh[j].error
}

func (mh MinHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
}

func (mh *MinHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	img := old[n-1]
	*mh = old[0 : n-1]
	return img
}

func (mh *MinHeap) Push(x interface{}) {
	img := x.(*Img)
	*mh = append(*mh, img)
}

// analyzeImage takes in an Img and returns the error
func analyzeImage(i *Img) {
	img := i.img
	hist, pix := histogram(img)
	avg := averageRGB(hist, pix)
	i.color = avg
	i.error = calculateError(hist, avg)
	return
}

// splitImage splits the input image into 4 equal images by width and height
func splitImage(img *image.NRGBA) []*Img {
	img_min_x, img_max_x := img.Bounds().Min.X, img.Bounds().Max.X
	img_min_y, img_max_y := img.Bounds().Min.Y, img.Bounds().Max.Y
	img_width, img_height := img_max_x-img_min_x, img_max_y-img_min_y

	r1 := image.Rect(img_min_x, img_min_y, img_width/2, img_height/2)
	r2 := image.Rect(img_width/2, img_min_y, img_max_x, img_height/2)
	r3 := image.Rect(img_min_x, img_height/2, img_width/2, img_max_y)
	r4 := image.Rect(img_width/2, img_height/2, img_max_x, img_max_y)

	l := []*image.NRGBA{
		imaging.Clone(img.SubImage(r1)),
		imaging.Clone(img.SubImage(r2)),
		imaging.Clone(img.SubImage(r3)),
		imaging.Clone(img.SubImage(r4)),
	}

	nl := make([]*Img, 0)
	for j := 0; j < len(l); j++ {
		newNode := Img{
			img: l[j],
		}
		analyzeImage(&newNode)
		nl = append(nl, &newNode)
	}
	return nl
}

// calculateError takes in an array of pixels separated by R,G,B values and an array of R,G,B average values
// it returns an int64 for the total error
func calculateError(hist [][]int, avg []float64) float64 {
	re, ge, be := 0.0, 0.0, 0.0
	ravg, gavg, bavg := avg[0], avg[1], avg[2]
	for i := 0; i < len(hist); i++ {
		re += math.Pow(float64(hist[i][0])-ravg, 2)
		ge += math.Pow(float64(hist[i][1])-gavg, 2)
		be += math.Pow(float64(hist[i][2])-bavg, 2)
	}
	return (re + ge + be) / float64(len(hist))
}

// averageRGB takes in an array of pixels separated by R,G,B values and calculates the average R,G,B value
func averageRGB(hist [][]int, p int) []float64 {
	r, g, b := 0, 0, 0
	for i := 0; i < p; i++ {
		r += hist[i][0]
		g += hist[i][1]
		b += hist[i][2]
	}
	pix := float64(p)
	avg := []float64{float64(r) / pix, float64(g) / pix, float64(b) / pix}
	return avg
}

// histogram takes in an image and returns a list of pixels separated by R,G,B values
func histogram(img *image.NRGBA) ([][]int, int) {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	p := w * h
	hist := make([][]int, p)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			loc := y*img.Stride + x*4

			r := int(img.Pix[loc])
			g := int(img.Pix[loc+1])
			b := int(img.Pix[loc+2])

			i := []int{r, g, b}
			hist[y*w+x] = i
		}
	}

	return hist, p
}

// toNRGBA converts any image type to *image.NRGBA with min-point at (0, 0).
func toNRGBA(img image.Image) (*image.NRGBA, error) {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0, nil
		}
	}
	return nil, fmt.Errorf("Incorrect Bounds: NRGBA", nil)
}
