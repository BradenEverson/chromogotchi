package main

type Pet struct {
	name         string
	happiness    float32
	hunger       float32
	wakefullness float32

	depression float32
	hungerRate float32
	sleepyRate float32
}

func makePet(name string) Pet {
	return Pet{name, 100.0, 100.0, 100.0, 0.5, 1.0, 2.5}
}
