package main

import (
	"fmt"
	"log"
	"os"

	"github.com/claydowling/genesyscharacter/genesys"
)

func main() {
	archetypes := genesys.ReadArchetypeFile("data/archetypes/base.yaml")

	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Executing %s\n", executable)
	for idx, a := range archetypes {
		fmt.Printf("%2d) %s\n", idx, a.Name)
	}
}
