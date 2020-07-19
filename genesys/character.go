package genesys

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

// Skill represents a skill, which may be applied to a Character to become a CharacterSkill
type Skill struct {
	Name    string
	Ability string
}

// CharacterSkill is a skill attached to a character
type CharacterSkill struct {
	ProficiencyDice int
	AbilityDice     int
	Skill
}

// Character represents the additions made to the character beyond the archetype,
// by the player.
type Character struct {
	Name       string
	Archetype  string
	Brawn      int
	Agility    int
	Intellect  int
	Cunning    int
	Will       int
	Presence   int
	Experience int
	Skills     map[string]int
}

// CalculatedCharacter is the result of Calculate, with archetype, skills, and feats applied
type CalculatedCharacter struct {
	Name       string
	Archetype  string
	Brawn      int
	Agility    int
	Intellect  int
	Cunning    int
	Will       int
	Presence   int
	Experience int
	Skills     []CharacterSkill
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

func calculateSkill(name string, level int, cc *CalculatedCharacter, skills *[]Skill) CharacterSkill {

	var cs CharacterSkill
	var skill Skill
	for _, s := range *skills {
		if s.Name == name {
			skill = s
			break
		}
	}

	cs.Ability = skill.Ability
	cs.Name = skill.Name

	var abilitylevel int
	switch skill.Ability {
	case "Agility":
		abilitylevel = cc.Brawn
	case "Brawn":
		abilitylevel = cc.Brawn
	case "Cunning":
		abilitylevel = cc.Brawn
	case "Intellect":
		abilitylevel = cc.Intellect
	case "Presence":
		abilitylevel = cc.Brawn
	case "Will":
		abilitylevel = cc.Brawn
	}

	cs.ProficiencyDice = level
	cs.AbilityDice = abilitylevel - level

	return cs
}

// Calculate takes the character and known archetypes, returns a fully calculated character
func Calculate(character Character, archetypes []Archetype, skills []Skill) (*CalculatedCharacter, error) {

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
	c.Intellect = a.Intellect + character.Intellect
	c.Presence = a.Presence + character.Presence
	c.Will = a.Will + character.Will

	for k, v := range character.Skills {
		cs := calculateSkill(k, v, &c, &skills)
		c.Skills = append(c.Skills, cs)
	}

	return &c, nil
}

func readYamlFile(filename string, dest interface{}) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(dat), dest)
	if err != nil {
		return err
	}

	return nil
}

// ReadCharacterFile reads a single character file
func ReadCharacterFile(filename string) (*Character, error) {
	var c Character
	err := readYamlFile(filename, &c)
	return &c, err
}

// ReadSkillFile loads skills from a file and returns them in an array.
// Returns an empty list and error on failure
func ReadSkillFile(filename string) ([]Skill, error) {
	var s []Skill
	err := readYamlFile(filename, &s)
	return s, err
}
