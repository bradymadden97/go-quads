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
	}
	headNode.hist, headNode.pix = histogram(img)
	headNode.color, headNode.error = analyzeImage(&headNode)
	return &headNode, nil
}

func iterate(mh *MinHeap, hn *Img, itr int, fn string, b bool, c bool, bc string, ds bool) error {
	cl, err := decodeColor(bc)
	if err != nil {
		return err
	}
	for i := 0; i < itr-1; i++ {
		a := heap.Pop(mh).(*Img)
		a.c1, a.c2, a.c3, a.c4 = splitHistogram(a.hist, a.width, a.height)

		heap.Push(mh, a.c1)
		heap.Push(mh, a.c2)
		heap.Push(mh, a.c3)
		heap.Push(mh, a.c4)

		if !ds {
			err := saveImage(hn, fn, i, b, c, itr, cl)
			if err != nil {
				return err
			}
		}
	}
	err = saveImage(hn, fn, itr, b, c, itr, cl)
	if err != nil {
		return err
	}
	return nil
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

			i := []int{r, g, b}
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
	newNode.color, newNode.error = analyzeImage(&newNode)
	return &newNode
}

func displayImage(head *Img, border bool, circle bool, colorlist []uint8) *image.NRGBA {
	base := color.RGBA{uint8(head.color[0]), uint8(head.color[1]), uint8(head.color[2]), 255}
	canvas := imaging.New(head.width, head.height, base)
	if border {
		canvas = addBorder(head.width, head.height, canvas, colorlist)
	}
	return traverseTree(canvas, head, image.Point{0, 0}, border, circle, colorlist)
}

func traverseTree(canvas *image.NRGBA, node *Img, p image.Point, border bool, circle bool, colorlist []uint8) *image.NRGBA {
	if node.c1 == nil && node.c2 == nil && node.c3 == nil && node.c4 == nil {
		c := color.RGBA{uint8(node.color[0]), uint8(node.color[1]), uint8(node.color[2]), 255}
		a := imaging.New(node.width, node.height, c)
		if border {
			a = addBorder(node.width, node.height, a, colorlist)
		}
		if circle {
			a = addCircle(node.width, node.height, a, colorlist)
		}
		canvas = imaging.Paste(canvas, a, p)
	} else {
		canvas = traverseTree(canvas, node.c1, p, border, circle, colorlist)
		canvas = traverseTree(canvas, node.c2, image.Point{p.X + int(node.width/2), p.Y}, border, circle, colorlist)
		canvas = traverseTree(canvas, node.c3, image.Point{p.X, p.Y + int(node.height/2)}, border, circle, colorlist)
		canvas = traverseTree(canvas, node.c4, image.Point{p.X + int(node.width/2), p.Y + int(node.height/2)}, border, circle, colorlist)
	}
	return canvas
}

func addCircle(w int, h int, img *image.NRGBA, cl []uint8) *image.NRGBA {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if euclideanDistance(w/2, x, h/2, y) >= ovalRadius(w/2, h/2, getAngle(w, h, x, y)) {
				copy(img.Pix[x*4+y*img.Stride:(x+1)*4+y*img.Stride], cl)
			}
		}
	}
	return img
}

func addBorder(w int, h int, img *image.NRGBA, bor []uint8) *image.NRGBA {
	for x := 0; x < w; x++ {
		copy(img.Pix[x*4:(x+1)*4], bor)
		copy(img.Pix[x*4+(h-1)*img.Stride:(x+1)*4+(h-1)*img.Stride], bor)
	}

	for y := 1; y < h-1; y++ {
		copy(img.Pix[y*img.Stride:y*img.Stride+4], bor)
		copy(img.Pix[y*img.Stride+w*4-4:y*img.Stride+w*4], bor)
	}

	return img
}

func saveImage(i *Img, in string, itr int, border bool, color bool, max int, colorlist []uint8) error {
	fo := displayImage(i, border, color, colorlist)
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
