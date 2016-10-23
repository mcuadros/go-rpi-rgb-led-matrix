package rgbmatrix

import (
	"image"
	"image/draw"
	"image/gif"
	"io"
	"time"
)

type ToolKit struct {
	Canvas *Canvas
}

func (tk *ToolKit) PlayImage(i image.Image, delay time.Duration) error {
	start := time.Now()
	defer func() { time.Sleep(delay - time.Since(start)) }()

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), i, image.ZP, draw.Over)
	return tk.Canvas.Render()
}

func (tk *ToolKit) PlayImages(images []image.Image, delay []time.Duration, loop int) chan bool {
	quit := make(chan bool, 0)

	go func() {
		l := len(images)
		i := 0
		for {
			select {
			case <-quit:
				return
			default:
				tk.PlayImage(images[i], delay[i])
			}

			i++
			if i >= l {
				if loop == 0 {
					i = 0
					continue
				}

				break
			}
		}
	}()

	return quit
}

func (tk *ToolKit) PlayGIF(r io.Reader) (chan bool, error) {
	gif, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}

	delay := make([]time.Duration, len(gif.Delay))
	images := make([]image.Image, len(gif.Image))
	for i, image := range gif.Image {
		images[i] = image
		delay[i] = time.Millisecond * time.Duration(gif.Delay[i]) * 10
	}

	return tk.PlayImages(images, delay, gif.LoopCount), nil
}
