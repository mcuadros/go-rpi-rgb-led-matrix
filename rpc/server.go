package rpc

import (
	"fmt"
	"image/color"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

type RPCMatrix struct {
	m rgbmatrix.Matrix
}

type GeometryArgs struct{}
type GeometryReply struct{ Width, Height int }

func (m *RPCMatrix) Geometry(_ *GeometryArgs, reply *GeometryReply) error {
	w, h := m.m.Geometry()
	reply.Width = w
	reply.Height = h

	return nil
}

type ApplyArgs struct{ Colors []color.Color }
type ApplyReply struct{}

func (m *RPCMatrix) Apply(args *ApplyArgs, reply *ApplyReply) error {
	return m.m.Apply(args.Colors)
}

type CloseArgs struct{}
type CloseReply struct{}

func (m *RPCMatrix) Close(_ *CloseArgs, _ *CloseReply) error {
	return m.m.Close()
}

func Serve(m rgbmatrix.Matrix) {
	rpc.Register(&RPCMatrix{m})

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println(l)
	http.Serve(l, nil)
}
