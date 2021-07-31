package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

var (
	//go:embed Tile_7.png
	tile7 []byte
)

var (
	blockImage *ebiten.Image
)

func init() {
	blockImage = decodeImage(tile7)
}

type Block struct {
	locationX float64
	dead      bool
}

type Blocks []Block

func (b Blocks) filterDead() Blocks {
	tmp := Blocks{}
	for _, block := range b {
		if !block.dead {
			tmp = append(tmp, block)
		}
	}
	return tmp
}

func drawBlocks(screen *ebiten.Image, g *Game) {
	for _, block := range g.blocks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(block.locationX, 96)
		screen.DrawImage(blockImage, op)
	}
}

func updateBlocks(blocks Blocks) Blocks {
	for i := range blocks {
		block := &blocks[i]
		block.locationX -= 1
		if block.locationX < -16 {
			block.locationX = WIDTH
			block.dead = true
		}
	}

	return blocks.filterDead()
}

func (b Block) GetRectangle() image.Rectangle {
	return image.Rect(int(b.locationX), 96, S+int(b.locationX), S+96)
}
