package main

// Functions for use Go Graphics github.com/fogleman/gg
// Simplifying position of images and text using zero based positive grid
// x, y grid axis references. Designed to insure images position
// themselves within the canvas padding. Includes functions
// to show grid lines and reference x,y text overlays and demo how
// images would position when assigned to each reference point.

import (
	"fmt"
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Struct with information about the image
// canvas or background layer.
type Canvas struct {
	ctx     *gg.Context
	width   int
	height  int
	padding float64
	color   string
	gridX   int
	gridY   int
	offsetX float64
	offsetY float64
	X       float64
	Y       float64
}

// Initializes various values in Canvas to make
// grid reference calculations simpler.
func (c *Canvas) init() {
	c.ctx = gg.NewContext(c.width, c.height)
	c.ctx.SetHexColor(c.color)
	c.ctx.Clear()
	w := float64(c.width)
	h := float64(c.width)
	p := float64(c.padding)
	c.offsetX = (w - p - p) / float64(c.gridX)
	c.offsetY = (h - p - p) / float64(c.gridY)
	c.X = (w / 2) - c.offsetX*(float64(c.gridX)/2)
	c.Y = (h / 2) + c.offsetY*(float64(c.gridY)/2)
}

// Draws vertical and horizontal grid lines on Canvas
func drawGridLines(c *Canvas) {

	w := float64(c.width)
	h := float64(c.width)
	p := float64(c.padding)

	vLines, hLines := c.gridX, c.gridY
	sw := (w - p - p) / float64(vLines)
	sh := (h - p - p) / float64(hLines)

	// Horizontal grid lines
	c.ctx.SetLineWidth(1)
	c.ctx.SetRGBA255(0, 0, 0, 30)

	shiftW := 0.0
	shiftH := 0.0
	for x := 0; x <= hLines/2; x++ {
		c.ctx.DrawLine(p, h/2+shiftH, w-p, h/2+shiftH)
		c.ctx.DrawLine(p, h/2-shiftH, w-p, h/2-shiftH)
		shiftW += sw
		shiftH += sh
	}
	c.ctx.Stroke()

	// Vertical grid lines
	c.ctx.SetLineWidth(1)
	c.ctx.SetRGBA255(0, 0, 0, 30)
	shiftW = 0.0
	shiftH = 0.0
	for x := 0; x <= vLines/2; x++ {
		c.ctx.DrawLine(w/2+shiftW, p, w/2+shiftW, h-p)
		c.ctx.DrawLine(w/2-shiftW, p, w/2-shiftW, h-p)
		shiftW += sw
		shiftH += sh
	}
	c.ctx.Stroke()

	// Outside rectangle and cross intersection grid lines
	c.ctx.SetLineWidth(3)
	c.ctx.SetRGBA255(230, 0, 0, 30)
	c.ctx.DrawLine(w/2, p, w/2, h-p)
	c.ctx.DrawLine(p, h/2, w-p, h/2)
	c.ctx.DrawRectangle(p, p, w-p-p, h-p-p)
	c.ctx.Stroke()
}

// Takes in a Canvas and draws grid references at grid axis
// Can also show images at grid reference points
func drawGridRef(c *Canvas, gridText bool, gridImage bool) {

	// if gridText true setup default go TTF font
	if gridText {
		font, err := truetype.Parse(goregular.TTF)
		if err != nil {
			log.Fatal(err)
		}
		c.ctx.SetRGBA255(0, 0, 0, 255)
		fontFace := truetype.NewFace(font, &truetype.Options{Size: 14})
		c.ctx.SetFontFace(fontFace)
	}

	// if gridImage true load a the image file
	var im image.Image
	var err error
	if gridImage {
		im, err = gg.LoadPNG("gopher.png")
		if err != nil {
			panic(err)
		}
	}

	var setX, setY int
	var smallGrid, largeGrid int
	var setMax, setMin int
	if c.gridX > c.gridY {
		setMax = c.gridX
		setMin = c.gridY
	} else {
		setMax = c.gridY
		setMin = c.gridX
	}

	for min := 0; min <= setMin; min++ {
		largeGrid = 0
		for max := 0; max <= setMax; max++ {
			if c.gridX > c.gridY {
				setX = largeGrid
				setY = smallGrid
			} else {
				setX = smallGrid
				setY = largeGrid
			}
			px, py, ax, ay := setGridPos(c, setX, setY)
			if gridImage {
				c.ctx.DrawImageAnchored(im, px, py, ax, ay)
			}
			if gridText {
				c.ctx.DrawStringWrapped(fmt.Sprintf("%d,%d", setX, setY), float64(px), float64(py)-5, ax, ay, 0, 1, gg.AlignCenter)
			}
			largeGrid++
		}
		smallGrid++
	}
}

// Sets start position 0,0 in lower left of image/grid
// Increments positions based on positive passed in x, y values
// Returns values required by gg draw functions
func setGridPos(c *Canvas, x, y int) (px, py int, ax, ay float64) {
	startX := int(c.X)
	startY := int(c.Y)

	if x > c.gridX {
		x = c.gridX
	}

	if y > c.gridY {
		y = c.gridY
	}

	switch x {
	case 0:
		px, ax = startX, 0.0
	case c.gridX:
		px, ax = startX+int(c.offsetX)*c.gridX, 1.0
	}

	switch y {
	case 0:
		py, ay = startY, 1.0
	case c.gridY:
		py, ay = startY-int(c.offsetY)*c.gridY, 0.0
	}

	if x > 0 && x < c.gridX {
		px, ax = startX+int(c.offsetX)*x, 0.5
	}

	if y > 0 && y < c.gridY {
		py, ay = startY-int(c.offsetY)*y, 0.5
	}

	return px, py, ax, ay
}

// Demo of text and image based grid references
// Outputs two images files gridref.png and gridimg.png
func main() {
	fmt.Println("Running gridgg demo...")

	layer1 := &Canvas{
		width:   1024,
		height:  1024,
		padding: 30,
		gridX:   8,
		gridY:   8,
		color:   "#FAFAFA",
	}

	layer2 := &Canvas{
		width:   1024,
		height:  1024,
		padding: 15,
		gridX:   8,
		gridY:   8,
		color:   "#FAFAFA",
	}

	// layer1 create grid text reference points
	layer1.init()
	drawGridRef(layer1, true, false)
	drawGridLines(layer1)
	fmt.Println("Saving file gridref.png...")
	layer1.ctx.SavePNG("gridref.png")

	// layer2 create grid image reference points
	// Demo image is sent in function to gopher.png
	layer2.init()
	drawGridRef(layer2, false, true)
	drawGridLines(layer2)
	fmt.Println("Saving file gridimg.png...")
	layer2.ctx.SavePNG("gridimg.png")
}
