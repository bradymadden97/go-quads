// util.go
package main

import (
	"image"
	"io"
	"os"
	"os/exec"
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

func concatName(name string, itr string) string {
	n, end := splitName(name)
	return itr + n + "." + end
}

func splitName(name string) (string, string) {
	splt := strings.Split(name, ".")
	return strings.Join(splt[:len(splt)-1], "."), splt[len(splt)-1]
}

func toGIF(name string, frames int, pause int) error {
	n, _ := splitName(name)
	cmd := exec.Command("python", "gif.py", "-n", n, "-f", strconv.Itoa(frames), "-p", strconv.Itoa(pause))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
