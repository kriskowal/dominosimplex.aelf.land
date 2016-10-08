package main

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"github.com/ojrac/opensimplex-go"
)

var pal = color.Palette{
	color.Black,
	color.White,
}

var positive = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 1, 0, 0, 0, 1, 0, 0},
	{0, 0, 1, 0, 1, 0, 1, 0, 0},
	{1, 0, 1, 0, 0, 0, 1, 0, 1},
	{1, 0, 1, 0, 1, 0, 1, 0, 1},
	{1, 1, 1, 0, 0, 0, 1, 1, 1},
	{1, 1, 1, 0, 1, 0, 1, 1, 1},
	{1, 1, 1, 1, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1},
}

var negative = [][]int{
	{0, 0, 0, 0, 1, 0, 0, 0, 0},
	{1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 1},
	{0, 1, 0, 1, 0, 1, 0, 1, 0},
	{0, 1, 0, 1, 1, 1, 0, 1, 0},
	{1, 0, 1, 1, 0, 1, 1, 0, 1},
	{1, 0, 1, 1, 1, 1, 1, 0, 1},
	{1, 1, 1, 1, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1},
}

const pad = 1
const scale = 0.1
const seed1 = 50
const seed2 = 175
const length = 100
const width = 50
const height = 50
const delay = 10

func inner(img *image.Paletted, x, y int, bitmap []int) {
	i := 0
	for dy := 0; dy < 3; dy++ {
		for dx := 0; dx < 3; dx++ {
			img.Set(x*(3+pad)+dx, y*(3+pad)+dy, pal[bitmap[i]])
			i++
		}
	}
}

func outer(img *image.Paletted, x, y int, v float64) {
	mag := math.Abs(v)
	ceil := math.Ceil(mag)
	bit := int(ceil)
	rem := int((mag / ceil) * 10)
	if mag == ceil {
		rem = 0
	}
	var bitmap, nitmap []int
	if v < 0 {
		bitmap = negative[bit]
		nitmap = negative[rem]
	} else {
		bitmap = positive[bit]
		nitmap = positive[rem]
	}
	i := 0
	for dy := 0; dy < 3; dy++ {
		for dx := 0; dx < 3; dx++ {
			if bitmap[i] != 0 {
				inner(img, x*(3+pad)+dx, y*(3+pad)+dy, nitmap)
			}
			i++
		}
	}
}

func cell(noise *opensimplex.Noise, x, y, i int) float64 {
	ox := math.Cos(2*math.Pi*float64(i)/length) * 1
	oy := math.Sin(4*math.Pi*float64(i)/length) * 1
	a := noise.Eval2((float64(x)+ox)*scale, (float64(y)+oy)*scale)
	b := noise.Eval2((float64(x-width)+ox)*scale, (float64(y)+oy)*scale)
	ab := (a*(1.0-float64(x)/width) + b*(float64(x)/width))
	c := noise.Eval2((float64(x)+ox)*scale, (float64(y-height)+oy)*scale)
	d := noise.Eval2((float64(x-width)+ox)*scale, (float64(y-height)+oy)*scale)
	cd := (c*(1.0-float64(x)/width) + d*(float64(x)/width))
	return 10 * (ab*(1.0-float64(y)/height) + cd*(float64(y)/height))
}

func main() {
	fullwidth, fullheight := width*(3+pad)*(3+pad), height*(3+pad)*(3+pad)
	rec := image.Rectangle{image.ZP, image.Pt(fullwidth, fullheight)}
	images := make([]*image.Paletted, 0, length)
	delays := make([]int, 0, length)

	noise := opensimplex.NewWithSeed(seed1)

	for i := 0; i < length; i++ {
		img := image.NewPaletted(rec, pal)
		images = append(images, img)
		delays = append(delays, delay)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				outer(img, x, y, float64(cell(noise, x, y, i)))
			}
		}
	}

	f, err := os.OpenFile("kriskowaldominosimplex.gif", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
