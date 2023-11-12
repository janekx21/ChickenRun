package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func decodeImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
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

func decodeImageStrip(data []byte, count int) []*ebiten.Image {
	img := decodeImage(data)
	var images = make([]*ebiten.Image, 0, 8)
	for i := 0; i < count; i++ {
		x, y := i, 0
		sub := img.SubImage(image.Rect(x*S, y*S, x*S+S, y*S+S)).(*ebiten.Image)
		images = append(images, sub)
	}
	return images
}
