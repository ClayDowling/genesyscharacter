package genesys

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v3"
)

// Archetype represents a character archetype, used to calculate starting values for a character
type Archetype struct {
	Name         string
	Brawn        int
	Agility      int
	Intelligence int
	Cunning      int
	Will         int
	Presence     int
	Wound        int
	Strain       int
	Experience   int
}

// ReadFile reads all of the Archetypes in filename and returns that as a slice
func ReadFile(filename string) []Archetype {

	var archetypes []Archetype

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s: %v\n", filename, err)
		return archetypes
	}

	err = yaml.Unmarshal([]byte(dat), &archetypes)
	if err != nil {
		log.Printf("Unmarshaling %s: %v\n", filename, err)
	}

	return archetypes
}
