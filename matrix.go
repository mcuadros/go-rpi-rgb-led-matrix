package main

/*
#cgo CFLAGS: -std=c99 -Ivendor/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -Lvendor/rpi-rgb-led-matrix/lib -lstdc++ -lm
#include <led-matrix-c.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>

struct RGBLedMatrix *matrix;
struct LedCanvas *offscreen_canvas;

void setCanvas(struct LedCanvas *offscreen_canvas, const uint32_t pixels[]) {
  int i, x, y;
	uint32_t color;
		for (y = 0; y < 32; ++y) {
			for (x = 0; x < 64; ++x) {
				i = x + (y * 64);

    		color = pixels[i];
    		led_canvas_set_pixel(offscreen_canvas, x, y,
					(color >> 16) & 255, (color >> 8) & 255, color & 255);
			}
	}
}

int example(struct RGBLedMatrixOptions *options) {
  struct LedCanvas *offscreen_canvas;
  int width, height;
  int x, y, i;

  matrix = led_matrix_create_from_options(options, NULL, NULL);
  if (matrix == NULL)
    return 1;

  offscreen_canvas = led_matrix_create_offscreen_canvas(matrix);

  led_canvas_get_size(offscreen_canvas, &width, &height);

  fprintf(stderr, "Size: %dx%d. Hardware gpio mapping: %s\n",
          width, height, options->hardware_mapping);

  for (i = 0; i < 1000; ++i) {
    for (y = 0; y < height; ++y) {
      for (x = 0; x < width; ++x) {
        led_canvas_set_pixel(offscreen_canvas, x, y, i & 0xff, x, y);
      }
    }

    offscreen_canvas = led_matrix_swap_on_vsync(matrix, offscreen_canvas);
  }

  led_matrix_delete(matrix);

  return 0;
}
*/
import "C"
import (
	"fmt"
	"image/color"
	"time"
	"unsafe"
)

func main() {
	o := &C.struct_RGBLedMatrixOptions{}
	o.rows = 32
	o.chain_length = 2

	s := time.Now()
	C.example(o)
	fmt.Println(time.Since(s))

	s = time.Now()
	example(o)
	fmt.Println(time.Since(s))

}

func example(o *C.struct_RGBLedMatrixOptions) {
	matrix := C.led_matrix_create_from_options(o, nil, nil)
	if matrix == nil {
		panic("jodido")
	}

	//offscreen := C.led_matrix_create_offscreen_canvas(matrix)
	for i := 1; i < 1000; i++ {
		p := make([]C.uint32_t, 2048)
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				pos := x + (y * 64)

				c := color.RGBA{uint8(i & 0xff), uint8(x), uint8(y), 0}
				p[pos] = C.uint32_t(colorToUint32(c))
			}
		}

		C.setCanvas(matrix, (*C.uint32_t)(unsafe.Pointer(&p[0])))
		//offscreen = C.led_matrix_swap_on_vsync(matrix, offscreen)
	}

	C.led_matrix_delete(matrix)
}

func colorToUint32(c color.Color) uint32 {
	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, _ := c.RGBA()
	//	r<<16 | g<<8 | b
	return (red>>8)<<16 | (green>>8)<<8 | blue>>8
}
