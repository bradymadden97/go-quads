// util.go
package main

import (
	"image"
	"image/gif"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andybons/gogif"
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

func concatName(name string, itr int) string {
	n, end := splitName(name)
	return n + strconv.Itoa(itr) + "." + end
}

func splitName(name string) (string, string) {
	splt := strings.Split(name, ".")
	return strings.Join(splt[:len(splt)-1], "."), splt[len(splt)-1]
}

func toGIF(fn string, delay int) error {
	name, _ := splitName(fn)
	outGIF := &gif.GIF{}
	imgList := []string{}

	err := filepath.Walk("./out", func(i string, f os.FileInfo, err error) error {
		imgList = append(imgList, i)
		return nil
	})
	if err != nil {
		return err
	}
	for _, i := range imgList[1:] {
		img, open_err := openImage(i)
		if open_err != nil {
			return open_err
		}
		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 256}
		quantizer.Quantize(palettedImage, bounds, img, image.ZP)
		outGIF.Image = append(outGIF.Image, palettedImage)
		outGIF.Delay = append(outGIF.Delay, delay)
	}

	out_gif, gif_err := os.OpenFile("./out/"+name+".gif", os.O_WRONLY|os.O_CREATE, 0600)
	if gif_err != nil {
		return gif_err
	}
	defer out_gif.Close()
	gif.EncodeAll(out_gif, outGIF)
	return nil
}
