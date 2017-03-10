package emulator

import (
	"image"
	"image/color"
	"log"

	"sync"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xff}
	red   = color.RGBA{0x7f, 0x00, 0x00, 0x7f}
)

var margin = 10

type Emulator struct {
	PixelPitch int
	Gutter     int
	Width      int
	Height     int

	leds    []color.Color
	w       screen.Window
	isReady bool

	wg sync.WaitGroup
}

func (e *Emulator) Init() {
	e.leds = make([]color.Color, 2048)

	e.wg.Add(1)
	go e.init()
	e.wg.Wait()
}

func (e *Emulator) init() {
	driver.Main(func(s screen.Screen) {
		var err error
		e.w, err = s.NewWindow(&screen.NewWindowOptions{
			Title: "Basic Shiny Example",
		})

		if err != nil {
			panic(err)
		}

		defer e.w.Release()

		var sz size.Event
		for {
			evn := e.w.NextEvent()
			switch evn := evn.(type) {
			case paint.Event:
				for _, r := range imageutil.Border(sz.Bounds(), margin) {
					e.w.Fill(r, red, screen.Src)
				}

				e.w.Fill(sz.Bounds().Inset(margin), black, screen.Src)
				e.w.Publish()
				if e.isReady {
					continue
				}

				e.Apply(make([]color.Color, 2048))
				e.wg.Done()
				e.isReady = true
			case size.Event:
				sz = evn

			case error:
				log.Print(e)
			}
		}
	})
}

func (e *Emulator) Geometry() (width, height int) {
	return e.Width, e.Height
}

func (e *Emulator) Apply(leds []color.Color) error {
	defer func() { e.leds = make([]color.Color, 2048) }()

	for col := 0; col < e.Width; col++ {
		for row := 0; row < e.Height; row++ {
			x := col * (e.PixelPitch + e.Gutter)
			y := row * (e.PixelPitch + e.Gutter)

			x += margin * 2
			y += margin * 2

			c := e.At(col + (row * e.Width))
			e.w.Fill(image.Rect(x, y, x+e.PixelPitch, y+e.PixelPitch), c, screen.Over)
		}
	}

	e.w.Publish()
	return nil
}

func (e *Emulator) Render() error {
	return e.Apply(e.leds)
}

func (e *Emulator) At(position int) color.Color {
	if e.leds[position] == nil {
		return color.Black
	}

	return e.leds[position]
}

func (e *Emulator) Set(position int, c color.Color) {
	e.leds[position] = color.RGBAModel.Convert(c)
}

func (e *Emulator) Close() error {
	return nil
}
