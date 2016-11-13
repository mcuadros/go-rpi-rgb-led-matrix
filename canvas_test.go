package rgbmatrix

import (
	"image/color"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CanvasSuite struct{}

var _ = Suite(&CanvasSuite{})

func (s *CanvasSuite) TestNewCanvas(c *C) {
	canvas := NewCanvas(NewMatrixMock())
	c.Assert(canvas, NotNil)
	c.Assert(canvas.w, Equals, 64)
	c.Assert(canvas.h, Equals, 32)
}

func (s *CanvasSuite) TestRender(c *C) {
	m := NewMatrixMock()
	canvas := &Canvas{m: m}
	canvas.Render()

	c.Assert(m.called["Render"], Equals, true)
}

func (s *CanvasSuite) TestColorModel(c *C) {
	canvas := &Canvas{}

	c.Assert(canvas.ColorModel(), Equals, color.RGBAModel)
}

func (s *CanvasSuite) TestBounds(c *C) {

	canvas := &Canvas{w: 10, h: 20}

	b := canvas.Bounds()
	c.Assert(b.Min.X, Equals, 0)
	c.Assert(b.Min.Y, Equals, 0)
	c.Assert(b.Max.X, Equals, 10)
	c.Assert(b.Max.Y, Equals, 20)
}

func (s *CanvasSuite) TestAt(c *C) {
	m := NewMatrixMock()
	canvas := &Canvas{w: 10, h: 20, m: m}
	canvas.At(5, 15)

	c.Assert(m.called["At"], Equals, 155)
}

func (s *CanvasSuite) TestSet(c *C) {
	m := NewMatrixMock()
	canvas := &Canvas{w: 10, h: 20, m: m}
	canvas.Set(5, 15, color.White)

	c.Assert(m.called["Set"], Equals, 155)
	c.Assert(m.colors[155], Equals, color.White)
}

func (s *CanvasSuite) TestClear(c *C) {
	m := NewMatrixMock()

	canvas := &Canvas{w: 10, h: 20, m: m}
	err := canvas.Clear()
	c.Assert(err, IsNil)

	for _, px := range m.colors {
		c.Assert(px, Equals, color.Black)
	}

	c.Assert(m.called["Render"], Equals, true)
}

func (s *CanvasSuite) TestClose(c *C) {
	m := NewMatrixMock()
	canvas := &Canvas{w: 10, h: 20, m: m}
	err := canvas.Close()
	c.Assert(err, IsNil)

	for _, px := range m.colors {
		c.Assert(px, Equals, color.Black)
	}

	c.Assert(m.called["Render"], Equals, true)
}

type MatrixMock struct {
	called map[string]interface{}
	colors []color.Color
}

func NewMatrixMock() *MatrixMock {
	return &MatrixMock{
		called: make(map[string]interface{}, 0),
		colors: make([]color.Color, 200),
	}
}

func (m *MatrixMock) Geometry() (width, height int) {
	return 64, 32
}

func (m *MatrixMock) Initialize() error {
	m.called["Initialize"] = true
	return nil
}

func (m *MatrixMock) At(position int) color.Color {
	m.called["At"] = position
	return color.Black
}

func (m *MatrixMock) Set(position int, c color.Color) {
	m.called["Set"] = position
	m.colors[position] = c
}

func (m *MatrixMock) Apply(leds []color.Color) error {
	for position, l := range leds {
		m.Set(position, l)
	}

	return m.Render()
}

func (m *MatrixMock) Render() error {
	m.called["Render"] = true
	return nil
}

func (m *MatrixMock) Close() error {
	m.called["Close"] = true
	return nil
}
