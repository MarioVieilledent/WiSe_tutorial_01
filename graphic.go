package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const w int = 1920
const h int = 1080

var node_16_img *ebiten.Image
var node_4_img *ebiten.Image
var node_1_img *ebiten.Image

func init() {
	var err error
	node_16_img, _, err = ebitenutil.NewImageFromFile("assets/node_16.png")
	if err != nil {
		log.Fatal(err)
	}
	node_4_img, _, err = ebitenutil.NewImageFromFile("assets/node_4.png")
	if err != nil {
		log.Fatal(err)
	}
	node_1_img, _, err = ebitenutil.NewImageFromFile("assets/node_1.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Sys struct {
	width  int // Window width
	height int // Window height

	minX float64 // Minimum longitude
	maxX float64 // Maximum longitude
	minY float64 // Minimum latitude
	maxY float64 // Maximum latitude

	aX float64 // Slope of longitude function
	bX float64 // Y-intercept of longitude function

	aY float64 // Slope of latitude function
	bY float64 // Y-intercept of latitude function
}

func newSys() Sys {
	// Retrieve OSM information
	minX, maxX, minY, maxY := GetMinMaxPoints()

	aX := (float64(w) / (maxX - minX))
	bX := -aX * minX

	aY := (float64(h) / (minY - maxY))
	bY := -aY * maxY

	// fmt.Println(aX, bX, aY, bY)

	return Sys{
		w,
		h,
		minX,
		maxX,
		minY,
		maxY,
		aX,
		bX,
		aY,
		bY,
	}
}

func (s *Sys) Update() error {
	return nil
}

func (s *Sys) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	for _, node := range Nodes {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.aX*node.LON+s.bX-1, s.aY*node.LAT+s.bY-1)
		// op.GeoM.Scale(1.5, 1)
		screen.DrawImage(node_1_img, op)
	}
}

func (s *Sys) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.width, s.height
}

func openWindow() {
	// Create the graphic system
	sys := newSys()

	// Setting up window parameters
	ebiten.SetWindowSize(sys.width, sys.height)
	ebiten.SetWindowTitle("Hello, World!")

	// Starting the rendering
	if err := ebiten.RunGame(&sys); err != nil {
		log.Fatal(err)
	}
}
