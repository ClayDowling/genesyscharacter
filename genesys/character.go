package genesys

import (
	"fmt"
)

type Character struct {
	Name         string
	Archetype    string
	Brawn        int
	Agility      int
	Intelligence int
	Cunning      int
	Will         int
	Presence     int
	Experience   int
}

type CalculatedCharacter struct {
	Name         string
	Archetype    string
	Brawn        int
	Agility      int
	Intelligence int
	Cunning      int
	Will         int
	Presence     int
	Experience   int
}

// FindArchetype searches for needing in a haystack of archetypes.
// returns a pointer to the archetype if found, error if not
func FindArchetype(needle string, haystack []Archetype) (*Archetype, error) {
	for _, a := range haystack {
		if a.Name == needle {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("Unknown Archetype '%s'", needle)
}

// Calculate takes the character and known archetypes, returns a fully calculated character
func Calculate(character Character, archetypes []Archetype) (*CalculatedCharacter, error) {

	a, err := FindArchetype(character.Archetype, archetypes)
	if err != nil {
		return nil, err
	}

	var c CalculatedCharacter

	c.Name = character.Name

	c.Agility = a.Agility + character.Agility
	c.Archetype = a.Name
	c.Brawn = a.Brawn + character.Brawn
	c.Cunning = a.Cunning + character.Cunning
	c.Experience = a.Experience + character.Experience
	c.Intelligence = a.Intelligence + character.Intelligence
	c.Presence = a.Presence + character.Presence
	c.Will = a.Will + character.Will

	return &c, nil
}
