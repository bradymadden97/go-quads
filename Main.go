// Main.go
package main

import (
	"fmt"
	"image"
	"log"

	"github.com/disintegration/imaging"
)

func main() {
	input_img, err := imaging.Open("tiger.jpg")
	if err != nil {
		log.Fatalf("Open image failed: %v", err)
	}
	img, err := toNRGBA(input_img)
	if err != nil {
		log.Fatalf("Image error: %v", err)
	}

	hist, pix := histogram(img)
	avg := averageRGB(hist, pix)

	//fmt.Printf("%v", len(hist))
	fmt.Printf("%v", avg)
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
