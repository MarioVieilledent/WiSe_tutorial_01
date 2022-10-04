package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const w int = 1280
const h int = 720

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
	/*
		for _, node := range Nodes {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(s.aX*node.LON+s.bX-1, s.aY*node.LAT+s.bY-1)
			// op.GeoM.Scale(1.5, 1)
			screen.DrawImage(node_1_img, op)
		}
	*/
	for key, way := range Ways {
		for id, _ := range way.list {
			if id < len(way.list)-1 {
				clr := getColor((way.properties["highway"]))
				// fmt.Println(s.aX*way[id].LON+s.bX, s.aX*way[id+1].LON+s.bX)
				if len(QueriedWaysId) == 0 {
					// Displays a line in the screen
					ebitenutil.DrawLine(screen,
						s.aX*way.list[id].LON+s.bX,
						s.aY*way.list[id].LAT+s.bY,
						s.aX*way.list[id+1].LON+s.bX,
						s.aY*way.list[id+1].LAT+s.bY,
						clr,
					)
				} else {
					if contains(QueriedWaysId, key) {
						// Displays a line in the screen
						ebitenutil.DrawLine(screen,
							s.aX*way.list[id].LON+s.bX,
							s.aY*way.list[id].LAT+s.bY,
							s.aX*way.list[id+1].LON+s.bX,
							s.aY*way.list[id+1].LAT+s.bY,
							clr,
						)
					}
				}
			}
		}
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

func getColor(highway string) color.Color {
	switch highway {
	case "residential":
		return color.RGBA64{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff}
	case "footway":
		return color.RGBA64{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff}
	case "tertiary":
		return color.RGBA64{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff}
	case "tertiary_link":
		return color.RGBA64{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff}
	case "path":
		return color.RGBA64{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff}
	case "pedestrian":
		return color.RGBA64{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff}
	case "steps":
		return color.RGBA64{R: 0x88ff, G: 0xffff, B: 0xffff, A: 0xffff}
	case "service":
		return color.RGBA64{R: 0xffff, G: 0x88ff, B: 0xffff, A: 0xffff}
	case "secondary":
		return color.RGBA64{R: 0xffff, G: 0xffff, B: 0x88ff, A: 0xffff}
	case "cycleway":
		return color.RGBA64{R: 0x88ff, G: 0x88ff, B: 0xffff, A: 0xffff}
	case "living_street":
		return color.RGBA64{R: 0xffff, G: 0x88ff, B: 0x88ff, A: 0xffff}
	default:
		return color.RGBA64{R: 0x88ff, G: 0xffff, B: 0x88ff, A: 0xffff}
	}
}

func contains(s []int, id int) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}
	return false
}
