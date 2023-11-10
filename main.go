package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
	_ "image/png"
	"log"
	"math/rand"
)

//go:generate goversioninfo -icon="icon.ico"

const S = 16
const WIDTH = 160
const HEIGHT = 144

var (
	//go:embed asset/background_0.png
	background0 []byte
	//go:embed asset/background_1.png
	background1 []byte
	//go:embed asset/background_2.png
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
	time        float64
	gameSpeed   float64
	spawnTimer  float64
	backgroundX float64
	player      Player
	blocks      Blocks
}

func NewGame() Game {
	return Game{player: Player{}, blocks: make(Blocks, 0), gameSpeed: 1, spawnTimer: 1}
}

func (g *Game) Update() error {
	deltaTime := 1.0 / 60.0 * g.gameSpeed // (1.0 / ebiten.CurrentTPS()) * g.gameSpeed

	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		deltaTime *= .05
	}

	g.time += deltaTime
	g.spawnTimer -= deltaTime

	g.player = updatePlayer(g.player)
	g.blocks = updateBlocks(g.blocks, deltaTime)

	for _, block := range g.blocks {
		if block.Bounds().Overlaps(g.player.Bounds()) {
			*g = NewGame()
		}
	}

	if g.spawnTimer <= 0 {
		spawnVariations := []float64{96, 96, 96, 96 - 16, 96 - 32}
		y := spawnVariations[rand.Intn(len(spawnVariations))]
		g.blocks = append(g.blocks, Block{pos: f64.Vec2{WIDTH - g.spawnTimer*60, y}})
		g.spawnTimer = 2*g.gameSpeed + rand.Float64()
	}
	g.backgroundX += deltaTime * 60

	g.gameSpeed += 1.0 / 60.0 / 16

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	drawBackground(screen, int(g.backgroundX), 0, 3)
	drawBackground(screen, int(g.backgroundX), 1, 2)

	drawPlayer(screen, g)
	drawBlocks(screen, g)

	drawBackground(screen, int(g.backgroundX), 2, 1)

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

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	w, h := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowSize(w, h)
	ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Chicken Run")
	game := NewGame()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
