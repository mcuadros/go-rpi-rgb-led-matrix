package main

import (
	"flag"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

var (
	rows                     = flag.Int("led-rows", 32, "number of rows supported")
	cols                     = flag.Int("led-cols", 32, "number of columns supported")
	parallel                 = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain                    = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness               = flag.Int("brightness", 100, "brightness (0-100)")
	hardware_mapping         = flag.String("led-gpio-mapping", "regular", "Name of GPIO mapping used.")
	show_refresh             = flag.Bool("led-show-refresh", false, "Show refresh rate.")
	inverse_colors           = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
	disable_hardware_pulsing = flag.Bool("led-no-hardware-pulse", false, "Don't use hardware pin-pulse generation.")
	img                      = flag.String("image", "", "image path")

	rotate = flag.Int("rotate", 0, "rotate angle, 90, 180, 270")
)

func main() {
	f, err := os.Open(*img)
	fatal(err)

	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.HardwareMapping = *hardware_mapping
	config.ShowRefreshRate = *show_refresh
	config.InverseColors = *inverse_colors
	config.DisableHardwarePulsing = *disable_hardware_pulsing

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	switch *rotate {
	case 90:
		tk.Transform = imaging.Rotate90
	case 180:
		tk.Transform = imaging.Rotate180
	case 270:
		tk.Transform = imaging.Rotate270
	}

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
