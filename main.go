package main

import (
	"fmt"

	"github.com/claydowling/genesyscharacter/genesys"
)

func main() {
	archetypes := genesys.ReadArchetypeFile("base.arc")

	for idx, a := range archetypes {
		fmt.Printf("%2d) %s\n", idx, a.Name)
	}
}
