package main

import (
	"flag"
	"os"
	"time"

	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows       = flag.Int("led-rows", 32, "number of rows supported")
	chain      = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness = flag.Int("brightness", 100, "brightness (0-100)")
	image      = flag.String("image", "", "image path")
)

func main() {
	f, err := os.Open(*image)
	fatal(err)

	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.ChainLength = *chain
	config.Brightness = *brightness

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	close, err := tk.PlayGIF(f)
	fatal(err)

	time.Sleep(time.Second * 30)
	close <- true
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
