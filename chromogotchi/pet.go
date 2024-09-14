package main

type Pet struct {
	name      string
	face      string
	happiness float32
	hunger    float32
}

var faces []string = []string{"👾"}

func makePet(name string) Pet {
	return Pet{name, faces[0], float32(100), float32(100)}
}
