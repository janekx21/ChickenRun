package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
	"gopkg.in/yaml.v3"
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
	time            float64
	gameSpeed       float64
	blockSpawnTimer float64
	coinSpawnTimer  float64
	backgroundX     float64
	player          Player
	blocks          Blocks
	coins           Coins
	saveGame        SaveGame
}

type SaveGame struct {
	Score int64
}

func NewGame() Game {
	saveGame := loadGame()
	return Game{player: Player{}, blocks: make(Blocks, 0), gameSpeed: 1, blockSpawnTimer: 1, coins: make(Coins, 0), saveGame: saveGame}
}

func (g *Game) Update() error {
	// deltaTime := 1.0 / 60.0 * g.gameSpeed // (1.0 / ebiten.CurrentTPS()) * g.gameSpeed
	deltaTime := (1.0 / math.Max(ebiten.CurrentTPS(), 60.0)) * g.gameSpeed

	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		deltaTime *= .05
	}

	if ebiten.IsKeyPressed(ebiten.KeyF4) {
		g.coins = append(g.coins, Coin{pos: f64.Vec2{100, 100}})
	}

	g.time += deltaTime
	g.blockSpawnTimer -= deltaTime
	g.coinSpawnTimer -= deltaTime

	g.player = updatePlayer(g.player)
	g.blocks = updateBlocks(g.blocks, deltaTime)
	g.coins = updateCoins(g.coins, deltaTime)

	for _, block := range g.blocks {
		if block.Bounds().Overlaps(g.player.Bounds()) {
			*g = NewGame()
		}
	}

	for i, coin := range g.coins {
		if coin.Bounds().Overlaps(g.player.Bounds()) {
			g.coins[i].dead = true
			g.saveGame.Score += 5
			saveGame(g.saveGame)
		}
	}

	if g.blockSpawnTimer <= 0 {
		spawnVariations := []float64{96, 96, 96, 96 - 16, 96 - 32}
		y := spawnVariations[rand.Intn(len(spawnVariations))]
		g.blocks = append(g.blocks, Block{pos: f64.Vec2{WIDTH - g.blockSpawnTimer*60, y}})
		g.blockSpawnTimer = 2*g.gameSpeed + rand.Float64()
	}

	if g.coinSpawnTimer <= 0 {
		spawnVariations := []float64{96, 96, 96, 96 - 16, 96 - 32}
		y := spawnVariations[rand.Intn(len(spawnVariations))]
		g.coins = append(g.coins, Coin{pos: f64.Vec2{WIDTH - g.coinSpawnTimer*60, y}})
		g.coinSpawnTimer = g.gameSpeed + rand.Float64()
	}
	g.backgroundX += deltaTime * 60

	g.gameSpeed += 1.0 / 60.0 / 16

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	drawBackground(screen, int(g.backgroundX), 0, 3)
	drawBackground(screen, int(g.backgroundX), 1, 2)

	drawCoins(screen, g)
	drawBlocks(screen, g)
	drawPlayer(screen, g)

	drawBackground(screen, int(g.backgroundX), 2, 1)

	if ebiten.IsKeyPressed(ebiten.KeyF2) {
		drawDebug(screen, g)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("$%d", g.saveGame.Score))
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

func saveGame(saveGame SaveGame) {
	data, err := yaml.Marshal(&saveGame)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	dir = filepath.Join(dir, ".chicken_run.yml")
	os.WriteFile(dir, data, 0644)
}

func loadGame() SaveGame {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: %v using default\n", err)
		return SaveGame{}
	}
	dir = filepath.Join(dir, ".chicken_run.yml")
	data, err := os.ReadFile(dir)

	saveData := SaveGame{}
	err = yaml.Unmarshal(data, &saveData)
	if err != nil {
		log.Fatalf("error: %v using default\n", err)
		return SaveGame{}
	}
	return saveData
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
