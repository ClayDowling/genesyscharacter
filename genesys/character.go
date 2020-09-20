package genesys

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// SettingDefaultLocation will load settings from the main distribution location.
const SettingDefaultLocation = "__default__"

// Archetype represents a character archetype
type Archetype struct {
	Name       string
	Brawn      int
	Agility    int
	Intellect  int
	Cunning    int
	Will       int
	Presence   int
	Wound      int
	Strain     int
	Experience int
}

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

// Talent describes a talent from the rule book
type Talent struct {
	Name   string
	Tier   int
	Ranked bool
	Gives  string
}

// CharacterTalent describes a talent as applied to a character, taking into account multiple
// occurrances for ranked talents.
type CharacterTalent struct {
	Level int
	Talent
}

// Character represents the additions made to the character beyond the archetype,
// by the player.
type Character struct {
	Name       string
	Player     string
	Profession string
	Archetype  string
	Brawn      int
	Agility    int
	Intellect  int
	Cunning    int
	Will       int
	Presence   int
	Experience int
	Skills     map[string]int
	Talents    []string
}

// Setting provides the defined components of the setting the character was built in.
type Setting struct {
	Archetypes []Archetype
	Skills     []Skill
	Talents    []Talent
}

// CalculatedCharacter is the result of Calculate, with archetype, skills, and feats applied
type CalculatedCharacter struct {
	Name       string
	Archetype  string
	Player     string
	Profession string
	Brawn      int
	Agility    int
	Intellect  int
	Cunning    int
	Will       int
	Presence   int
	Experience int
	Skills     map[string]CharacterSkill
	Talents    map[string]CharacterTalent
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

func calculateSkill(name string, level int, cc *CalculatedCharacter, skills *[]Skill) (CharacterSkill, error) {

	var cs CharacterSkill
	var skill Skill
	for _, s := range *skills {
		if s.Name == name {
			skill = s
			break
		}
	}
	if skill.Name == "" {
		return cs, fmt.Errorf("Could not match skill '%s'", name)
	}

	cs.Ability = skill.Ability
	cs.Name = skill.Name

	var abilitylevel int
	switch strings.ToLower(skill.Ability) {
	case "agility":
		abilitylevel = cc.Agility
	case "brawn":
		abilitylevel = cc.Brawn
	case "cunning":
		abilitylevel = cc.Cunning
	case "intellect":
		abilitylevel = cc.Intellect
	case "presence":
		abilitylevel = cc.Presence
	case "will":
		abilitylevel = cc.Will
	}

	cs.ProficiencyDice = level
	cs.AbilityDice = abilitylevel - level

	return cs, nil
}

// Calculate takes the character and known archetypes, returns a fully calculated character
func Calculate(character Character, setting Setting) (CalculatedCharacter, error) {

	var c CalculatedCharacter

	a, err := FindArchetype(character.Archetype, setting.Archetypes)
	if err != nil {
		return c, err
	}

	c.Name = character.Name
	c.Player = character.Player
	c.Profession = character.Profession
	c.Agility = a.Agility + character.Agility
	c.Archetype = a.Name
	c.Brawn = a.Brawn + character.Brawn
	c.Cunning = a.Cunning + character.Cunning
	c.Experience = a.Experience + character.Experience
	c.Intellect = a.Intellect + character.Intellect
	c.Presence = a.Presence + character.Presence
	c.Will = a.Will + character.Will

	c.Skills = make(map[string]CharacterSkill)
	for _, s := range setting.Skills {
		level, ok := character.Skills[s.Name]
		if !ok {
			level = 0
		}
		cs, err := calculateSkill(s.Name, level, &c, &setting.Skills)
		if err != nil {
			return c, err
		}
		c.Skills[cs.Name] = cs
	}

	return c, nil
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
func ReadCharacterFile(filename string) (Character, error) {
	var c Character
	err := readYamlFile(filename, &c)
	return c, err
}

// ReadSkillFile loads skills from a file and returns them in an array.
// Returns an empty list and error on failure
func ReadSkillFile(filename string) ([]Skill, error) {
	var s []Skill
	err := readYamlFile(filename, &s)
	return s, err
}

// ReadArchetypeFile loads archetypes from the listed file, returns
// an empty list and error on failure
func ReadArchetypeFile(filename string) ([]Archetype, error) {
	var a []Archetype
	err := readYamlFile(filename, &a)
	return a, err
}

// ReadTalentFile loads talents from the listed file, returns
// an empty list and error on failure
func ReadTalentFile(filename string) ([]Talent, error) {
	var t []Talent
	err := readYamlFile(filename, &t)
	return t, err
}

// ReadSetting will find a folder with the necessary name and read the files
// which define the setting.
//
// name:      name of the setting, which matches a folder in the setting directory
// sourcedir: folder to search for setting directories.  In normal usage SettingDefaultLocation
//            will search under the data folder where the executable lives.
func ReadSetting(name string, sourcedir string) Setting {
	var s Setting
	var err interface{}
	var base string

	if sourcedir == SettingDefaultLocation {
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		base = filepath.Join(filepath.Dir(exe), "data")
	} else {
		base = sourcedir
	}
	archetypesfile := filepath.Join(base, name, "archetypes.yaml")
	skillsfile := filepath.Join(base, name, "skills.yaml")
	talentsfile := filepath.Join(base, name, "talents.yaml")

	s.Archetypes, err = ReadArchetypeFile(archetypesfile)
	if err != nil {
		log.Fatalf("Could not read %s", archetypesfile)
	}
	s.Skills, err = ReadSkillFile(skillsfile)
	if err != nil {
		log.Fatalf("Could not read %s", skillsfile)
	}
	s.Talents, err = ReadTalentFile(talentsfile)
	if err != nil {
		log.Fatalf("Could not read %s", talentsfile)
	}
	return s
}
