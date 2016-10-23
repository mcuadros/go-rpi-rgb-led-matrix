package rpc

import (
	"encoding/gob"
	"image/color"
	"net/rpc"

	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

func init() {
	gob.Register(color.RGBA{})
}

// RGBLedMatrix matrix representation for ws281x
type Client struct {
	network string
	addr    string
	client  *rpc.Client
	leds    []color.Color
}

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewClient(network, addr string) (rgbmatrix.Matrix, error) {
	client, err := rpc.DialHTTP(network, addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		network: network,
		addr:    addr,
		client:  client,
		leds:    make([]color.Color, 2048),
	}, nil
}

// Geometry returns the width and the height of the matrix
func (c *Client) Geometry() (width, height int) {
	var reply *GeometryReply
	err := c.client.Call("RPCMatrix.Geometry", &GeometryArgs{}, &reply)
	if err != nil {
		panic(err)
	}

	return reply.Width, reply.Height
}

func (c *Client) Apply(leds []color.Color) error {
	defer func() { c.leds = make([]color.Color, 2048) }()

	var reply *ApplyReply
	return c.client.Call("RPCMatrix.Apply", &ApplyArgs{Colors: leds}, &reply)
}

// Render update the display with the data from the LED buffer
func (c *Client) Render() error {
	return c.Apply(c.leds)
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *Client) At(position int) color.Color {
	if c.leds[position] == nil {
		return color.Black
	}

	return c.leds[position]
}

// Set set LED at position x,y to the provided 24-bit color value.
func (m *Client) Set(position int, c color.Color) {
	m.leds[position] = color.RGBAModel.Convert(c)
}

// Close finalizes the ws281x interface
func (c *Client) Close() error {
	return c.Apply(make([]color.Color, 2048))
}
