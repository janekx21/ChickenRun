package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
)

func drawWiredRectangleRect(screen *ebiten.Image, rect image.Rectangle, color color.Color) {
	x := float64(rect.Min.X)
	y := float64(rect.Min.Y)
	w := float64(rect.Size().X)
	h := float64(rect.Size().Y)
	drawWiredRectangle(screen, x, y, w, h, color)
}

func drawWiredRectangle(screen *ebiten.Image, x, y, w, h float64, color color.Color) {
	ebitenutil.DrawLine(screen, x, y, x+w, y, color)
	ebitenutil.DrawLine(screen, x, y, x, y+h, color)
	ebitenutil.DrawLine(screen, x+w, y, x+w, y+h, color)
	ebitenutil.DrawLine(screen, x, y+h, x+w, y+h, color)
}
func drawDebug(screen *ebiten.Image, g *Game) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("cap(g.blocks) = %d", cap(g.blocks)))

	drawWiredRectangleRect(screen, g.player.GetRectangle(), colornames.Green)
	for _, block := range g.blocks {
		drawWiredRectangleRect(screen, block.GetRectangle(), colornames.Red)
		if block.GetRectangle().Overlaps(g.player.GetRectangle()) {
			ebitenutil.DebugPrint(screen, "Game Over!")
		}
	}
}
