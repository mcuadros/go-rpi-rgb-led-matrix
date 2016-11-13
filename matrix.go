package rgbmatrix

/*
#cgo CFLAGS: -std=c99 -I${SRCDIR}/vendor/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -L${SRCDIR}/vendor/rpi-rgb-led-matrix/lib -lstdc++ -lm
#include <led-matrix-c.h>

void led_matrix_swap(struct RGBLedMatrix *matrix,
                     int width, int height, const uint32_t pixels[]) {
  struct LedCanvas *offscreen_canvas;
  offscreen_canvas = led_matrix_create_offscreen_canvas(matrix);

  int i, x, y;
  uint32_t color;
  for (x = 0; x < width; ++x) {
    for (y = 0; y < height; ++y) {
      i = x + (y * width);
      color = pixels[i];

      led_canvas_set_pixel(offscreen_canvas, x, y,
        (color >> 16) & 255, (color >> 8) & 255, color & 255);
    }
  }

  led_matrix_swap_on_vsync(matrix, offscreen_canvas);
}
*/
import "C"
import (
	"fmt"
	"image/color"
	"unsafe"
)

// DefaultConfig default WS281x configuration
var DefaultConfig = HardwareConfig{
	Rows:              32,
	ChainLength:       1,
	Parallel:          1,
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
	Brightness:        100,
	ScanMode:          Progressive,
}

// HardwareConfig rgb-led-matrix configuration
type HardwareConfig struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int
	// ChainLengthis the number of displays daisy-chained together
	// (output of one connected to input of next).
	ChainLength int
	// Parallel is the number of parallel chains connected to the Pi; in old Pis
	// with 26 GPIO pins, that is 1, in newer Pis with 40 interfaces pins, that
	// can also be 2 or 3. The effective number of pixels in vertical direction is
	// then thus rows * parallel.
	Parallel int
	// Set PWM bits used for output. Default is 11, but if you only deal with
	// limited comic-colors, 1 might be sufficient. Lower require less CPU and
	// increases refresh-rate.
	PWMBits int
	// Change the base time-unit for the on-time in the lowest significant bit in
	// nanoseconds.  Higher numbers provide better quality (more accurate color,
	// less ghosting), but have a negative impact on the frame rate.
	PWMLSBNanoseconds int // the DMA channel to use
	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int
	// ScanMode progressive or interlaced
	ScanMode ScanMode // strip color layout
	// Disable the PWM hardware subsystem to create pulses. Typically, you don't
	// want to disable hardware pulsing, this is mostly for debugging and figuring
	// out if there is interference with the sound system.
	// This won't do anything if output enable is not connected to GPIO 18 in
	// non-standard wirings.
	DisableHardwarePulsing bool

	ShowRefreshRate bool
	InverseColors   bool
}

func (c *HardwareConfig) geometry() (width, height int) {
	return c.Rows * c.ChainLength, c.Rows * c.Parallel
}

func (c *HardwareConfig) toC() *C.struct_RGBLedMatrixOptions {
	o := &C.struct_RGBLedMatrixOptions{}
	o.rows = C.int(c.Rows)
	o.chain_length = C.int(c.ChainLength)
	o.parallel = C.int(c.Parallel)
	o.pwm_bits = C.int(c.PWMBits)
	o.pwm_lsb_nanoseconds = C.int(c.PWMLSBNanoseconds)
	o.brightness = C.int(c.Brightness)
	o.scan_mode = C.int(c.ScanMode)
	//	o.disable_hardware_pulsing = c.DisableHardwarePulsing
	//	o.show_refresh_rate = c.ShowRefreshRate
	//	o.inverse_colors = c.InverseColors

	return o
}

type ScanMode int8

const (
	Progressive ScanMode = 0
	Interlaced  ScanMode = 1
)

// RGBLedMatrix matrix representation for ws281x
type RGBLedMatrix struct {
	Config *HardwareConfig

	height int
	width  int
	matrix *C.struct_RGBLedMatrix
	leds   []C.uint32_t
}

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewRGBLedMatrix(config *HardwareConfig) (*RGBLedMatrix, error) {
	w, h := config.geometry()

	c := &RGBLedMatrix{
		Config: config,
		width:  w, height: h,
		matrix: C.led_matrix_create_from_options(config.toC(), nil, nil),
		leds:   make([]C.uint32_t, w*h),
	}

	if c.matrix == nil {
		return nil, fmt.Errorf("unable to allocate memory")
	}

	return c, nil
}

// Initialize initialize library, must be called once before other functions are
// called.
func (c *RGBLedMatrix) Initialize() error {
	return nil
}

// Geometry returns the width and the height of the matrix
func (c *RGBLedMatrix) Geometry() (width, height int) {
	return c.width, c.height
}

// Apply set all the pixels to the values contained in leds
func (c *RGBLedMatrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		c.Set(position, l)
	}

	return c.Render()
}

// Render update the display with the data from the LED buffer
func (c *RGBLedMatrix) Render() error {
	w, h := c.Config.geometry()

	C.led_matrix_swap(
		c.matrix,
		C.int(w), C.int(h),
		(*C.uint32_t)(unsafe.Pointer(&c.leds[0])),
	)

	c.leds = make([]C.uint32_t, w*h)
	return nil
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *RGBLedMatrix) At(position int) color.Color {
	return uint32ToColor(c.leds[position])
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *RGBLedMatrix) Set(position int, color color.Color) {
	c.leds[position] = C.uint32_t(colorToUint32(color))
}

// Close finalizes the ws281x interface
func (c *RGBLedMatrix) Close() error {
	C.led_matrix_delete(c.matrix)
	return nil
}

func colorToUint32(c color.Color) uint32 {
	if c == nil {
		return 0
	}

	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, _ := c.RGBA()
	return (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func uint32ToColor(u C.uint32_t) color.Color {
	return color.RGBA{
		uint8(u>>16) & 255,
		uint8(u>>8) & 255,
		uint8(u>>0) & 255,
		0,
	}
}
