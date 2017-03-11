package rgbmatrix

import (
	"image"
	"image/draw"
	"image/gif"
	"io"
	"time"
)

// ToolKit is a convinient set of function to operate with a led of Matrix
type ToolKit struct {
	// Canvas is the Canvas wrapping the Matrix, if you want to instanciate
	// a ToolKit with a custom Canvas you can use directly the struct,
	// without calling NewToolKit
	Canvas *Canvas

	// Transform function if present is applied just before draw the image to
	// the Matrix, this is a small example:
	//	tk.Transform = func(img image.Image) *image.NRGBA {
	//		return imaging.Fill(img, 64, 96, imaging.Center, imaging.Lanczos)
	//	}
	Transform func(img image.Image) *image.NRGBA
}

// NewToolKit returns a new ToolKit wrapping the given Matrix
func NewToolKit(m Matrix) *ToolKit {
	return &ToolKit{
		Canvas: NewCanvas(m),
	}
}

// PlayImage draws the given image during the given delay
func (tk *ToolKit) PlayImage(i image.Image, delay time.Duration) error {
	start := time.Now()
	defer func() { time.Sleep(delay - time.Since(start)) }()

	if tk.Transform != nil {
		i = tk.Transform(i)
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), i, image.ZP, draw.Over)
	return tk.Canvas.Render()
}

type Animation interface {
	Next() (image.Image, <-chan time.Time, error)
}

// PlayAnimation play the image during the delay returned by Next, until an err
// is returned, if io.EOF is returned, PlayAnimation finish without an error
func (tk *ToolKit) PlayAnimation(a Animation) error {
	var err error
	var i image.Image
	var n <-chan time.Time

	for {
		i, n, err = a.Next()
		if err != nil {
			break
		}

		if err := tk.PlayImageUntil(i, n); err != nil {
			return err
		}
	}

	if err == io.EOF {
		return nil
	}

	return err
}

// PlayImageUntil draws the given image until is notified to stop
func (tk *ToolKit) PlayImageUntil(i image.Image, notify <-chan time.Time) error {
	defer func() {
		<-notify
	}()

	if tk.Transform != nil {
		i = tk.Transform(i)
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), i, image.ZP, draw.Over)
	return tk.Canvas.Render()
}

// PlayImages draws a sequence of images during the given delays, the len of
// images should be equal to the len of delay. If loop is true the function
// loops over images until a true is sent to the returned chan
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

// PlayGIF reads and draw a gif file from r. It use the contained images and
// delays and loops over it, until a true is sent to the returned chan
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

// Close close the toolkit and the inner canvas
func (tk *ToolKit) Close() error {
	return tk.Canvas.Close()
}
