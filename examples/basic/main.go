package main

import (
	"flag"
	"image/color"

	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	chain      = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
)

func main() {
	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.ChainLength = *chain
	config.Brightness = *brightness

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

	bounds := c.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c.Set(x, y, color.RGBA{255, 0, 0, 255})
			c.Render()
		}
	}
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
