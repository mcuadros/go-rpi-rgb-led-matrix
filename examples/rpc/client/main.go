package main

import (
	"flag"
	"os"
	"time"

	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/mcuadros/go-rpi-rgb-led-matrix/rpc"
)

var (
	img = flag.String("image", "", "image path")
)

func main() {
	f, err := os.Open(*img)
	fatal(err)

	m, err := rpc.NewClient("tcp", "10.20.20.20:1234")
	fatal(err)

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

	tk := &rgbmatrix.ToolKit{c}
	close, err := tk.PlayGIF(f)
	fatal(err)

	time.Sleep(time.Second * 3)
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
