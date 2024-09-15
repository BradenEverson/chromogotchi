package main

import (
	"image/png"
	"math/rand"
	"os"
	"strconv"
)

var names = []string{"Peter", "Glorb", "Jrog", "Silbert", "Skleve"}

type Pet struct {
	name         string
	happiness    float32
	hunger       float32
	wakefullness float32

	sprite []byte

	depression float32
	hungerRate float32
	sleepyRate float32
}

func makePet() Pet {
	var petSprite []byte
    loc := rand.Intn(3) + 1;
    nameLoc := rand.Intn(len(names) - 1)

    name := names[nameLoc]

    // If fails, pet sprite is just blank
    defaultPet(&petSprite, "sprites/" + strconv.Itoa(loc) + "/idle.png")

	return Pet{name, 100.0, 100.0, 100.0, petSprite, 0.5, 1.0, 2.5}
}

func defaultPet(arr *[]byte, path string) error {
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

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			*arr = append(*arr, r8)
			*arr = append(*arr, g8)
			*arr = append(*arr, b8)
			*arr = append(*arr, a8)
		}
	}

	return nil
}
