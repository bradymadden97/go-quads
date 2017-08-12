// Main.go
package main

import (
	"container/heap"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/disintegration/imaging"
)

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
	headNode := initialize("tiger.jpg")

	//Begin minheap
	mh := make(MinHeap, 1)
	mh[0] = headNode
	heap.Init(&mh)

	//Loop
	iterations := 200
	for i := 0; i < iterations; i++ {
		a := heap.Pop(&mh).(*Img)
		a.c1, a.c2, a.c3, a.c4 = splitHistogram(a.hist, a.width, a.height)

		heap.Push(&mh, a.c1)
		heap.Push(&mh, a.c2)
		heap.Push(&mh, a.c3)
		heap.Push(&mh, a.c4)

		/*for z := 0; z < len(mh); z++ {
			fmt.Printf("%v, ", mh[z].error)
		}*/
		fmt.Printf("%v", a.error)
		fmt.Println()

		ifo := displayImage(headNode)
		imaging.Save(ifo, "./out/cheetah_"+strconv.Itoa(i)+".jpg")

	}

	fo := displayImage(headNode)
	imaging.Save(fo, "./out/cheetah_final.jpg")
}

func initialize(fn string) *Img {
	//Get input image and convert to NRGBA
	input_img, err := imaging.Open(fn)
	if err != nil {
		log.Fatalf("Open image failed: %v", err)
	}
	img, err := toNRGBA(input_img)
	if err != nil {
		log.Fatalf("Image error: %v", err)
	}

	//Create Img object, get initial error, start Quadtree
	headNode := Img{
		width:  img.Bounds().Max.X,
		height: img.Bounds().Max.Y,
	}
	headNode.hist, headNode.pix = histogram(img)
	analyzeImage(&headNode)

	return &headNode
}

func displayImage(head *Img) *image.NRGBA {
	base := color.RGBA{uint8(head.color[0]), uint8(head.color[1]), uint8(head.color[2]), 255}
	canvas := imaging.New(head.width, head.height, base)
	return traverseTree(canvas, head, image.Point{0, 0})
}

func traverseTree(canvas *image.NRGBA, node *Img, p image.Point) *image.NRGBA {
	if node.c1 == nil && node.c2 == nil && node.c3 == nil && node.c4 == nil {
		c := color.RGBA{uint8(node.color[0]), uint8(node.color[1]), uint8(node.color[2]), 255}
		a := imaging.New(node.width, node.height, c)
		canvas = imaging.Paste(canvas, a, p)
	} else {
		canvas = traverseTree(canvas, node.c1, p)
		canvas = traverseTree(canvas, node.c2, image.Point{p.X + int(node.width/2), p.Y})
		canvas = traverseTree(canvas, node.c3, image.Point{p.X, p.Y + int(node.height/2)})
		canvas = traverseTree(canvas, node.c4, image.Point{p.X + int(node.width/2), p.Y + int(node.height/2)})
	}
	return canvas
}

// analyzeImage takes in an Img and returns the error
func analyzeImage(i *Img) {
	avg := averageRGB(i.hist, i.pix)
	i.color = avg
	i.error = calculateError(i.hist, avg)
	return
}

/*
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
			img:    l[j],
			width:  l[j].Bounds().Max.X,
			height: l[j].Bounds().Max.Y,
		}
		analyzeImage(&newNode)
		nl = append(nl, &newNode)
	}
	return nl
}
*/

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

// splitHistogram takes in a double array and length and width of image and returns four arrays split by quadrants of subimage
func splitHistogram(h [][]int, w int, l int) (*Img, *Img, *Img, *Img) {
	c1, c2, c3, c4 := make([][]int, 0), make([][]int, 0), make([][]int, 0), make([][]int, 0)
	for i := 0; i < len(h); i++ {
		if i < int(l/2)*w {
			if i%w < w/2 {
				c1 = append(c1, h[i])
			} else {
				c2 = append(c2, h[i])
			}
		} else {
			if i%w < w/2 {
				c3 = append(c3, h[i])
			} else {
				c4 = append(c4, h[i])
			}
		}
	}
	return newNode(c1, w, l), newNode(c2, w, l), newNode(c3, w, l), newNode(c4, w, l)
}

func newNode(hist [][]int, w int, h int) *Img {
	newNode := Img{
		width:  int(w / 2),
		height: int(h / 2),
		hist:   hist,
	}
	newNode.pix = newNode.width * newNode.height
	analyzeImage(&newNode)
	return &newNode
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
