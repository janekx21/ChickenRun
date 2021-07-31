package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
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

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
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
		screen.DrawImage(playerJumpAnimation[g.frame/8%3], op)
	}
}

func (p Player) GetRectangle() image.Rectangle {
	x, y := S*4+4, 96-int(p.locationY)+6
	return image.Rect(x, y, x+S-8, y+S-6)
}
