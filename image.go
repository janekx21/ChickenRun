package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func decodeImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func loadImageArray(path string) [4]*ebiten.Image {
	img := loadImage(path)
	var images [4]*ebiten.Image
	for i := 0; i < 4; i++ {
		x, y := i%2, i/2
		sub := img.SubImage(image.Rect(x*S, y*S, x*S+S, y*S+S)).(*ebiten.Image)
		images[i] = sub
	}
	return images
}

func decodeImageArray(data []byte) [4]*ebiten.Image {
	img := decodeImage(data)
	var images [4]*ebiten.Image
	for i := 0; i < 4; i++ {
		x, y := i%2, i/2
		sub := img.SubImage(image.Rect(x*S, y*S, x*S+S, y*S+S)).(*ebiten.Image)
		images[i] = sub
	}
	return images
}
