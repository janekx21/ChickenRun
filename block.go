package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
	"image"
)

var (
	//go:embed asset/Tile_7.png
	tile7 []byte
)

var (
	blockImage *ebiten.Image
)

func init() {
	blockImage = decodeImage(tile7)
}

type Block struct {
	pos  f64.Vec2
	dead bool
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
		op.GeoM.Translate(block.pos[0], block.pos[1])
		screen.DrawImage(blockImage, op)
	}
}

func updateBlocks(blocks Blocks, dt float64) Blocks {
	for i := range blocks {
		block := &blocks[i]
		block.pos[0] -= dt * 60
		if block.pos[0] < -16 {
			block.pos[0] = WIDTH
			block.dead = true
		}
	}

	return blocks.filterDead()
}

func (b Block) Bounds() image.Rectangle {
	return image.Rect(int(b.pos[0]), int(b.pos[1]), S+int(b.pos[0]), int(S+b.pos[1]))
}
