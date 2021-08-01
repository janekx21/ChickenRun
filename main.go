package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)

const S = 16
const WIDTH = 160
const HEIGHT = 144

var (
	//go:embed background_0.png
	background0 []byte
	//go:embed background_1.png
	background1 []byte
	//go:embed background_2.png
	background2 []byte
)

var (
	backgroundImages [3]*ebiten.Image
)

func init() {
	log.SetPrefix("ChickenRun: ")
	log.SetFlags(0)
	backgroundImages[0] = decodeImage(background0)
	backgroundImages[1] = decodeImage(background1)
	backgroundImages[2] = decodeImage(background2)
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
		if block.Bounds().Overlaps(g.player.Bounds()) {
			*g = NewGame()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || g.frame%120 == 0 {
		g.blocks = append(g.blocks, Block{locationX: WIDTH})
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	drawBackground(screen, g.frame, 0, 3)
	drawBackground(screen, g.frame, 1, 2)

	drawPlayer(screen, g)
	drawBlocks(screen, g)

	drawBackground(screen, g.frame, 2, 1)

	if ebiten.IsKeyPressed(ebiten.KeyF2) {
		drawDebug(screen, g)
	}
}

func drawBackground(screen *ebiten.Image, frame int, index int, divide int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame/divide%WIDTH), 0)
	screen.DrawImage(backgroundImages[index], op)
	op.GeoM.Translate(WIDTH, 0)
	screen.DrawImage(backgroundImages[index], op)
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
