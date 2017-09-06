// util.go
package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func openImage(filename string) (*image.NRGBA, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := decodeImage(f)
	if err != nil {
		return nil, err
	}
	return cropImage(img), err
}

func decodeImage(i io.Reader) (*image.NRGBA, error) {
	img, _, err := image.Decode(i)
	if err != nil {
		return nil, err
	}
	nrgba, err := toNRGBA(img)
	if err != nil {
		return nil, err
	}
	return nrgba, nil
}

func cropImage(img *image.NRGBA) *image.NRGBA {
	mid_x, mid_y := int(img.Bounds().Max.X/2), int(img.Bounds().Max.Y/2)
	new_x, new_y := resizeBounds(img)
	r := image.Rect(mid_x-new_x/2, mid_y-new_y/2, mid_x+new_x/2, mid_y+new_y/2)
	return imaging.Clone(img.SubImage(r))
}

func resizeBounds(img *image.NRGBA) (int, int) {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y

	return nearestSquare(w), nearestSquare(h)
}

func toNRGBA(img image.Image) (*image.NRGBA, error) {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0, nil
		}
	}
	return imaging.Clone(img), nil
}

func nearestSquare(n int) int {
	i := 1
	for i <= n {
		i *= 2
	}
	return i / 2
}

func euclideanDistance(x1 int, x2 int, y1 int, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

func getAngle(w int, h int, x int, y int) float64 {
	dx, dy := math.Abs(float64(x-w/2)), math.Abs(float64(y-h/2))
	return math.Atan(dy / dx)
}

func ovalRadius(a int, b int, theta float64) float64 {
	return float64(a*b) / math.Sqrt(math.Pow(float64(a), 2)*math.Pow(math.Sin(theta), 2)+math.Pow(float64(b), 2)*math.Pow(math.Cos(theta), 2))
}

func decodeColor(bc string) ([]uint8, error) {
	l := strings.Split(bc, ",")
	if len(l) != 3 {
		return nil, fmt.Errorf("Error: backgroundcolor length %d not = 3", len(l))
	}
	cl := make([]uint8, 4)
	for i := 0; i < 3; i++ {
		s, err := strconv.ParseUint(l[i], 10, 64)
		if err != nil || s < 0 || s > 255 {
			return nil, err
		}
		cl[i] = uint8(s)
	}
	cl[3] = 255
	return cl, nil
}

func concatName(name string, itr string) string {
	n, end := splitName(name)
	return itr + n + "." + end
}

func splitName(name string) (string, string) {
	splt := strings.Split(name, ".")
	return strings.Join(splt[:len(splt)-1], "."), splt[len(splt)-1]
}

//Referenced https://github.com/esimov/stackblur-go/blob/master/cmd/main.go
func toGIF(imgs *[]image.Image, name string, delay int, pause int) error {
	outGif := &gif.GIF{}
	for _, i := range *imgs {
		inGif := image.NewPaletted(i.Bounds(), palette.Plan9)
		draw.Draw(inGif, i.Bounds(), i, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, inGif)
		outGif.Delay = append(outGif.Delay, delay)
	}
	n, _ := splitName(name)
	f, err := os.OpenFile(outputFolder+n+".gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, outGif)
}
