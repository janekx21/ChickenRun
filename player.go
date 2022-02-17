package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

var (
	//go:embed asset/Player_notmove.png
	playerNotMove []byte
	//go:embed asset/Player_jump.png
	playerJump []byte
	//go:embed asset/Player_hover.png
	playerHover []byte
)

var (
	playerNotMoveAnimation [4]*ebiten.Image
	playerJumpAnimation    [4]*ebiten.Image
	playerHoverImg         *ebiten.Image
)

func init() {
	playerNotMoveAnimation = decodeImageArray(playerNotMove)
	playerJumpAnimation = decodeImageArray(playerJump)
	playerHoverImg = decodeImage(playerHover)
}

type Player struct {
	onGround   bool
	blockHover bool
	hover      bool
	locationY  float64
	velocityY  float64
}

func jumpDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.GamepadAxis(0, 1) < -.5
}

func updatePlayer(player Player) Player {
	player.onGround = player.locationY <= 0
	player.hover = false

	if player.onGround {
		player.velocityY = 0
		player.locationY = 0

		if jumpDown() {
			player.velocityY = 4.6
			player.blockHover = true
		}
	} else {
		if jumpDown() {
			if !player.blockHover {
				player.hover = true
			}
		} else {
			player.blockHover = false
		}
	}

	if player.hover {
		player.velocityY = -.2
	} else {
		player.velocityY -= .2
	}
	player.locationY += player.velocityY
	return player
}

func drawPlayer(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(S*4, 96-g.player.locationY)
	if g.player.onGround {
		screen.DrawImage(playerNotMoveAnimation[int(g.time*60)/8%4], op)
	} else {
		if g.player.hover {
			screen.DrawImage(playerHoverImg, op)
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
