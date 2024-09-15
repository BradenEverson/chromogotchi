package main

import (
	"image/png"
	"os"
)

type Pet struct {
	name         string
	happiness    float32
	hunger       float32
	wakefullness float32

	sprite []int32

	depression float32
	hungerRate float32
	sleepyRate float32
}

func makePet(name string) Pet {
	var petSprite []int32
	// If fails, pet sprite is just blank
	defaultPet(&petSprite, "sprites/1/idle.png")

	return Pet{name, 100.0, 100.0, 100.0, petSprite, 0.5, 1.0, 2.5}
}

func defaultPet(arr *[]int32, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := int32(r >> 8)
			g8 := int32(g >> 8)
			b8 := int32(b >> 8)
			a8 := int32(a >> 8)

			pixel := (r8 << 24) | (g8 << 16) | (b8 << 8) | a8
			*arr = append(*arr, pixel)
		}
	}

	return nil
}
