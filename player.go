package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

var (
	//go:embed Player_notmove.png
	playerNotMove []byte
	//go:embed Player_jump.png
	playerJump []byte
)

var (
	playerNotMoveAnimation [4]*ebiten.Image
	playerJumpAnimation    [4]*ebiten.Image
)

func init() {
	playerNotMoveAnimation = decodeImageArray(playerNotMove)
	playerJumpAnimation = decodeImageArray(playerJump)
}

type Player struct {
	onGround  bool
	locationY float64
	velocityY float64
}

func updatePlayer(player Player) Player {
	player.onGround = player.locationY <= 0
	if player.onGround {
		player.velocityY = 0
		player.locationY = 0

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.GamepadAxis(0, 1) < -.5 {
			player.velocityY = 3
		}
	}

	player.velocityY -= .1
	player.locationY += player.velocityY
	return player
}

func drawPlayer(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(S*4, 96-g.player.locationY)
	if g.player.onGround {
		screen.DrawImage(playerNotMoveAnimation[g.frame/8%4], op)
	} else {
		state := sgn(math.Round(g.player.velocityY))
		stateToFrame := map[int]int{
			-1: 2,
			0:  1,
			1:  0,
		}
		screen.DrawImage(playerJumpAnimation[stateToFrame[state]], op)
	}
}

func sgn(a float64) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

func (p Player) Bounds() image.Rectangle {
	x, y := S*4+4, 96-int(p.locationY)+6
	return image.Rect(x, y, x+S-8, y+S-6)
}
