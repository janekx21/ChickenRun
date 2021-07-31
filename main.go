package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	_ "image/png"
	"log"
)

const S = 16
const WIDTH = 160
const HEIGHT = 144

var (
	//go:embed Tile_7.png
	tile7 []byte
	//go:embed background_0.png
	background0 []byte
	//go:embed background_1.png
	background1 []byte
	//go:embed background_2.png
	background2 []byte
)

var (
	backgroundImages [3]*ebiten.Image
	blockImage       *ebiten.Image
)

func init() {
	log.SetPrefix("ChickenRun: ")
	log.SetFlags(0)
	backgroundImages[0] = decodeImage(background0)
	backgroundImages[1] = decodeImage(background1)
	backgroundImages[2] = decodeImage(background2)
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

type Game struct {
	frame  int
	player Player
	blocks Blocks
}

func NewGame() Game {
	return Game{player: Player{}, blocks: make(Blocks, 0)}
}

func (g *Game) Update() error {
	g.frame += 1
	g.player = updatePlayer(g.player)
	g.blocks = updateBlocks(g.blocks)

	for _, block := range g.blocks {
		if block.GetRectangle().Overlaps(g.player.GetRectangle()) {
			*g = NewGame()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || g.frame%120 == 0 {
		g.blocks = append(g.blocks, Block{locationX: WIDTH})
	}

	return nil
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

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	drawBackground(screen, g.frame, 0, 3)
	drawBackground(screen, g.frame, 1, 2)
	drawBackground(screen, g.frame, 2, 1)

	drawPlayer(screen, g)
	drawBlocks(screen, g)

	// debug
	/*
		ebitenutil.DebugPrint(screen, fmt.Sprintf("cap(g.blocks) = %d", cap(g.blocks)))

		drawWiredRectangleRect(screen, g.player.GetRectangle(), colornames.Green)
		for _, block := range g.blocks {
			drawWiredRectangleRect(screen, block.GetRectangle(), colornames.Red)
			if block.GetRectangle().Overlaps(g.player.GetRectangle()) {
				ebitenutil.DebugPrint(screen, "Game Over!")
			}
		}
	*/
}

func drawBackground(screen *ebiten.Image, frame int, index int, divide int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame/divide%WIDTH), 0)
	screen.DrawImage(backgroundImages[index], op)
	op.GeoM.Translate(WIDTH, 0)
	screen.DrawImage(backgroundImages[index], op)
}

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

func drawBlocks(screen *ebiten.Image, g *Game) {
	for _, block := range g.blocks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(block.locationX, 96)
		screen.DrawImage(blockImage, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Chicken Run")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&Game{player: Player{}, blocks: make(Blocks, 0)}); err != nil {
		log.Fatal(err)
	}
}
