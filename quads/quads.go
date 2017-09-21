// quads.go
package main

import (
	"bytes"
	"container/heap"
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/disintegration/imaging"
)

func initialize(fn string) (*Img, error) {
	img, err := openImage(fn)
	if err != nil {
		return nil, err
	}
	headNode := Img{
		width:  img.Bounds().Max.X,
		height: img.Bounds().Max.Y,
		point:  image.Point{0, 0},
	}
	headNode.hist, headNode.pix = histogram(img)
	headNode.color, headNode.error = analyzeImage(&headNode)
	return &headNode, nil
}

func iterate(mh *MinHeap, hn *Img, itr int, fn string, b bool, c bool, bc string, s bool, g bool) ([]image.Image, error) {
	//imgs := make([]image.Image, itr)
	cl, err := decodeColor(bc)
	if err != nil {
		return nil, err
	}
	past_img := createImage(hn, b, c, cl)

	for i := 0; i < itr; i++ {
		if s {
			err := saveImage(past_img, fn, i, itr)
			if err != nil {
				return nil, err
			}
		}

		a := heap.Pop(mh).(*Img)
		if a.width <= 1 || a.height <= 1 {
			heap.Push(mh, a)
			break
		}
		a.c1, a.c2, a.c3, a.c4 = splitHistogram(a.hist, a.width, a.height, a.point)

		heap.Push(mh, a.c1)
		heap.Push(mh, a.c2)
		heap.Push(mh, a.c3)
		heap.Push(mh, a.c4)

		past_img = updateImage(past_img, []*Img{a.c1, a.c2, a.c3, a.c4}, b, c, cl)
	}
	err = saveImage(past_img, fn, itr, itr)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

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
			o := int(img.Pix[loc+3])

			i := []int{r, g, b, o}
			hist[y*w+x] = i
		}
	}
	return hist, p
}

func analyzeImage(i *Img) ([]float64, float64) {
	avg := averageRGB(i.hist, i.pix)
	err := calculateError(i.hist, avg)
	return avg, err
}

func averageRGB(hist [][]int, p int) []float64 {
	r, g, b, o := 0, 0, 0, 0
	for i := 0; i < p; i++ {
		r += hist[i][0]
		g += hist[i][1]
		b += hist[i][2]
		o += hist[i][3]
	}
	pix := float64(p)
	avg := []float64{float64(r) / pix, float64(g) / pix, float64(b) / pix, float64(o) / pix}
	return avg
}

func calculateError(hist [][]int, avg []float64) float64 {
	re, ge, be := 0.0, 0.0, 0.0
	ravg, gavg, bavg := avg[0], avg[1], avg[2]
	for i := 0; i < len(hist); i++ {
		re += math.Pow(float64(hist[i][0])-ravg, 2)
		ge += math.Pow(float64(hist[i][1])-gavg, 2)
		be += math.Pow(float64(hist[i][2])-bavg, 2)
	}
	return (re + ge + be)
}

func splitHistogram(h [][]int, w int, l int, p image.Point) (*Img, *Img, *Img, *Img) {
	c1, c2, c3, c4 := make([][]int, 0), make([][]int, 0), make([][]int, 0), make([][]int, 0)
	p1, p2, p3, p4 := image.Point{p.X, p.Y}, image.Point{p.X + w/2, p.Y}, image.Point{p.X, p.Y + l/2}, image.Point{p.X + w/2, p.Y + l/2}
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
	return newNode(c1, w, l, p1), newNode(c2, w, l, p2), newNode(c3, w, l, p3), newNode(c4, w, l, p4)
}

func newNode(hist [][]int, w int, h int, p image.Point) *Img {
	newNode := Img{
		width:  int(w / 2),
		height: int(h / 2),
		hist:   hist,
		point:  p,
	}
	newNode.pix = newNode.width * newNode.height
	newNode.color, newNode.error = analyzeImage(&newNode)
	return &newNode
}

func createImage(head *Img, border bool, circle bool, colorlist []uint8) *image.NRGBA {
	base := color.RGBA{uint8(head.color[0]), uint8(head.color[1]), uint8(head.color[2]), uint8(head.color[3])}
	canvas := imaging.New(head.width, head.height, base)
	if border {
		canvas = addBorder(canvas, head.width, head.height, head.point, colorlist)
	}
	if circle {
		canvas = addCircle(canvas, head.width, head.height, head.point, colorlist)
	}

	return canvas
}

func updateImage(img *image.NRGBA, sub_imgs []*Img, border bool, circle bool, colorlist []uint8) *image.NRGBA {
	for _, i := range sub_imgs {
		c := []uint8{uint8(i.color[0]), uint8(i.color[1]), uint8(i.color[2]), uint8(i.color[3])}
		new_img := pasteImage(img, i.width, i.height, i.point, c)
		if border {
			new_img = addBorder(new_img, i.width, i.height, i.point, colorlist)
		}
		if circle {
			new_img = addCircle(new_img, i.width, i.height, i.point, colorlist)
		}
	}
	return img
}

func pasteImage(bg *image.NRGBA, w int, h int, point image.Point, c []uint8) *image.NRGBA {
	for i := point.Y; i < point.Y+h; i++ {
		for j := point.X; j < point.X+w; j++ {
			copy(bg.Pix[i*bg.Stride+j*4:i*bg.Stride+(j+1)*4], c)
		}
	}
	return bg
}

func saveImage(fo *image.NRGBA, in string, itr int, max int) error {
	il, ml := len(strconv.Itoa(itr)), len(strconv.Itoa(max))
	var num bytes.Buffer
	for i := il; i < ml; i++ {
		num.WriteString("0")
	}
	num.WriteString(strconv.Itoa(itr))
	n := concatName(in, num.String())
	imaging.Save(fo, outputFolder+n)

	return nil
}

func addCircle(img *image.NRGBA, w int, h int, point image.Point, cl []uint8) *image.NRGBA {
	for y := point.Y; y < point.Y+h; y++ {
		for x := point.X; x < point.X+w; x++ {
			if euclideanDistance(w/2, x-point.X, h/2, y-point.Y) >= ovalRadius(w/2, h/2, getAngle(w, h, x-point.X, y-point.Y)) {
				copy(img.Pix[x*4+y*img.Stride:(x+1)*4+y*img.Stride], cl)
			}
		}
	}
	return img
}

func addBorder(img *image.NRGBA, w int, h int, point image.Point, bor []uint8) *image.NRGBA {
	for x := point.X; x < point.X+w; x++ {
		copy(img.Pix[point.Y*img.Stride+x*4:point.Y*img.Stride+(x+1)*4], bor)
		copy(img.Pix[(point.Y+(h-1))*img.Stride+x*4:(point.Y+(h-1))*img.Stride+(x+1)*4], bor)
	}
	for y := point.Y + 1; y < point.Y+h-1; y++ {
		copy(img.Pix[y*img.Stride+point.X*4:y*img.Stride+4+point.X*4], bor)
		copy(img.Pix[y*img.Stride+w*4-4+point.X*4:y*img.Stride+w*4+point.X*4], bor)
	}
	return img
}
