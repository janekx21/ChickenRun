package main

import (
	_ "embed"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

var (
	//go:embed asset/coin.png
	coin []byte
)

var (
	coinImage []*ebiten.Image
)

func init() {
	coinImage = decodeImageStrip(coin, 8)
}

type Coin struct {
	pos  f64.Vec2
	dead bool
	time float64
}

type Coins []Coin

func (c Coins) filterDead() Coins {
	tmp := Coins{}
	for _, coin := range c {
		if !coin.dead {
			tmp = append(tmp, coin)
		}
	}
	return tmp
}

func drawCoins(screen *ebiten.Image, g *Game) {
	for _, coin := range g.coins {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(coin.pos[0], coin.pos[1]+math.Cos(coin.time*math.Pi)*3)
		frame := int(math.Floor((coin.time * 8))) % 8
		screen.DrawImage(coinImage[frame], op)
	}
}

func updateCoins(coins Coins, dt float64) Coins {
	for i := range coins {
		coin := &coins[i]
		coin.pos[0] -= dt * 60
		coin.time += dt
		if coin.pos[0] < -16 {
			coin.pos[0] = WIDTH
			coin.dead = true
		}
	}

	return coins.filterDead()
}

func (b Coin) Bounds() image.Rectangle {
	return image.Rect(int(b.pos[0]), int(b.pos[1]), S+int(b.pos[0]), int(S+b.pos[1]))
}
