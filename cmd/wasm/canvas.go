package main

import (
	"image"
	"image/color"
	"syscall/js"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Canvas struct {
	// Window elements
	window js.Value
	doc    js.Value
	body   js.Value

	// Canvas properties
	canvas  js.Value
	ctx     js.Value
	imgData js.Value
	width   int
	height  int

	//
	gc    *draw2dimg.GraphicContext
	image *image.RGBA

	copyBuff js.Value

	canvasSize     float64
	simulationSize float64
}

func NewCanvas(simuSize float64, canvSize float64) (*Canvas, error) {
	var cv Canvas

	cv.window = js.Global()
	cv.doc = cv.window.Get("document")
	cv.body = cv.doc.Get("body")
	cv.canvas = cv.doc.Call("createElement", "canvas")

	cv.canvas.Set("height", canvSize)
	cv.canvas.Set("width", canvSize)

	cv.body.Call("appendChild", cv.canvas)

	cv.height = int(canvSize)
	cv.width = int(canvSize)

	cv.ctx = cv.canvas.Call("getContext", "2d")
	cv.imgData = cv.ctx.Call("createImageData", int(canvSize), int(canvSize))
	cv.image = image.NewRGBA(image.Rect(0, 0, int(canvSize), int(canvSize)))
	cv.copyBuff = js.Global().Get("Uint8Array").New(len(cv.image.Pix))

	cv.gc = draw2dimg.NewGraphicContext(cv.image)

	cv.canvasSize = canvSize
	cv.simulationSize = simuSize

	return &cv, nil
}

func (cv *Canvas) getRatio() float64 {
	return cv.canvasSize / cv.simulationSize
}

func (cv *Canvas) ClearGC(color color.RGBA) {
	cv.gc.SetFillColor(color)
	cv.gc.Clear()
}

func (cv *Canvas) Render() {
	js.CopyBytesToJS(cv.copyBuff, cv.image.Pix)
	cv.imgData.Get("data").Call("set", cv.copyBuff)
	cv.ctx.Call("putImageData", cv.imgData, 0, 0)
}
