package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
)

var names = []string{
	"Blorbo", "Flern", "Zizzle", "Krungus", "Snorp", "Wibble", "Glef", "Plorn", "Zogbert", "Snerk",
	"Glarn", "Vibbit", "Jorf", "Snibble", "Tronk", "Bibble", "Skibble", "Frong", "Jilbo", "Splork",
	"Norf", "Dronk", "Skorb", "Flib", "Grundle", "Mibbit", "Splarn", "Gribble", "Norbit", "Wobbert",
	"Plink", "Skrogg", "Tweb", "Jorb", "Zlorn", "Brindle", "Sklorm", "Flarn", "Globbo", "Trilk",
	"Frob", "Snarg", "Zindle", "Cribble", "Gorf", "Brug", "Dribbit", "Slibber", "Flibble", "Kloof",
	"Trog", "Glimp", "Borf", "Nimpl", "Fribbit", "Quorn", "Glibble", "Drong", "Spliff", "Gribber",
	"Woggin", "Jubb", "Florp", "Drob", "Skorn", "Glimble", "Flung", "Wormp", "Trorb", "Flink",
	"Brog", "Splim", "Zorp", "Nork", "Grob", "Flunk", "Skrob", "Glarn", "Prindle", "Brorf",
	"Nubb", "Sklon", "Frigg", "Jimble", "Dragg", "Klarg", "Vibble", "Plog", "Splorb", "Wibber",
	"Gronk", "Slibble", "Twirp", "Frogbert", "Blip", "Drongus", "Snig", "Blurg", "Twonk", "Splurb",
	"Grilk", "Morb", "Klimp", "Jibble", "Peter", "Glorb", "Jrog", "Silbert", "Skleve",
}

type Pet struct {
	Name         string  `bson:"name"`
	Happiness    float32 `bson:"happiness"`
	Hunger       float32 `bson:"hunger"`
	Wakefullness float32 `bson:"wakefullness"`

	Sprite int    `bson:"sprite"`
	State  string `bson:"state"`

	Depression float32 `bson:"depression"`
	HungerRate float32 `bson:"hungerRate"`
	SleepyRate float32 `bson:"sleepyRate"`
	Id         string  `bson:"id"`
}

func makePet(id string) Pet {
	loc := rand.Intn(3) + 1
	nameLoc := rand.Intn(len(names) - 1)

	name := names[nameLoc]

	return Pet{name, 100.0, 100.0, 100.0, loc, "idle", 1.5, 1.0, 2.5, id}
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

func (pet *Pet) getSprite() []byte {
	var sprite []byte

	state := pet.getUpdatedState()
	path := "./sprites/" + strconv.Itoa(pet.Sprite) + "/" + state + ".png"
	defaultPet(&sprite, path)

	return sprite
}

func (pet *Pet) updateHunger(newVal float32) {
	pet.Hunger += newVal
	if pet.Hunger <= 0 {
		pet.Hunger = 0
	}
}

func (pet *Pet) updateSleep(newVal float32) {
	pet.Wakefullness += newVal
	if pet.Wakefullness <= 0 {
		pet.Wakefullness = 0
	}

}

func (pet Pet) updateHappy(newVal float32) {
	pet.Happiness += newVal
	if pet.Happiness <= 0 {
		pet.Happiness = 0
	}

}

func (pet Pet) getUpdatedState() string {
	avg := (pet.Wakefullness + pet.Happiness + pet.Hunger) / 3.0
	fmt.Println(pet.Wakefullness)
	fmt.Println(pet.Happiness)
	fmt.Println(pet.Hunger)

	var state string

	switch {
	case avg == 0.0:
		state = "dead"
	case avg <= 10.0:
		state = "dying"
	case avg <= 30.0:
		state = "bad"
	case avg <= 50.0:
		state = "sad"
	case pet.Wakefullness <= 30.0:
		state = "tired"
	case pet.Happiness <= 30.0:
		state = "depressed"
	case pet.Hunger <= 30.0:
		state = "hungry"
	default:
		state = "idle"
	}

	pet.State = state
	return pet.State
}
