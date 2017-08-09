// Main.go
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
	img   *image.NRGBA
	error float64
}

type MinHeap []*Img

type ImgNode struct {
	node *Img
	c1   *ImgNode
	c2   *ImgNode
	c3   *ImgNode
	c4   *ImgNode
}

func main() {
	input_img, err := imaging.Open("tiger.jpg")
	if err != nil {
		log.Fatalf("Open image failed: %v", err)
	}
	img, err := toNRGBA(input_img)
	if err != nil {
		log.Fatalf("Image error: %v", err)
	}

	i := Img{
		img: img,
	}
	analyzeImage(&i)

	mh := make(MinHeap, 1)
	mh[0] = &i
	heap.Init(&mh)

	for mh.Len() > 0 {
		i := heap.Pop(&mh).(*Img)
		fmt.Printf("%v", i.error)
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
	i.error = calculateError(hist, avg)
	return
}

// splitImage splits the input image into 4 equal images by width and height
func splitImage(img *image.NRGBA) []*image.NRGBA {
	img_min_x, img_max_x := img.Bounds().Min.X, img.Bounds().Max.X
	img_min_y, img_max_y := img.Bounds().Min.Y, img.Bounds().Max.Y
	img_width, img_height := img_max_x-img_min_x, img_max_y-img_min_y

	r1 := image.Rect(img_min_x, img_min_y, img_width/2, img_height/2)
	r2 := image.Rect(img_width/2, img_min_y, img_max_x, img_height/2)
	r3 := image.Rect(img_min_x, img_height/2, img_width/2, img_max_y)
	r4 := image.Rect(img_width/2, img_height/2, img_max_x, img_max_y)

	s1, s2, s3, s4 := img.SubImage(r1), img.SubImage(r2), img.SubImage(r3), img.SubImage(r4)

	return []*image.NRGBA{imaging.Clone(s1), imaging.Clone(s2), imaging.Clone(s3), imaging.Clone(s4)}
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
