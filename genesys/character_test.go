package genesys

import (
	"log"
	"os"
	"strings"
	"testing"
)

var archetypes []Archetype
var skills []Skill

func TestMain(m *testing.M) {
	var err error
	archetypes, _ = ReadArchetypeFile("testfile.arc")
	skills, err = ReadSkillFile("testskills.skl")
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(m.Run())
}

func findArchetype(needle string, haystack []Archetype, t *testing.T) Archetype {
	for _, a := range haystack {
		if a.Name == needle {
			return a
		}
	}

	t.Errorf("Could not find archetype '%s'", needle)

	var e Archetype
	return e
}

func IsExpectedCalculatedCharacter(cc CalculatedCharacter, expected CalculatedCharacter, t *testing.T) {
	if cc.Agility != expected.Agility {
		t.Errorf("Expected Agility %d, got %d", expected.Agility, cc.Agility)
	}
	if cc.Brawn != expected.Brawn {
		t.Errorf("Expected Brawn %d, got %d", expected.Brawn, cc.Brawn)
	}
	if cc.Cunning != expected.Cunning {
		t.Errorf("Expected Cunning %d, got %d", expected.Cunning, cc.Cunning)
	}
	if cc.Experience != expected.Experience {
		t.Errorf("Expected Experience %d, got %d", expected.Experience, cc.Experience)
	}
	if cc.Intellect != expected.Intellect {
		t.Errorf("Expected Intellect %d, got %d", expected.Intellect, cc.Intellect)
	}
	if cc.Presence != expected.Presence {
		t.Errorf("Expected Presence %d, got %d", expected.Presence, cc.Presence)
	}
	if cc.Will != expected.Will {
		t.Errorf("Expected Will %d, got %d", expected.Will, cc.Will)
	}
}

func Test_CalculateCharacter_WithArchetype_GivesCharacterWithArchetype(t *testing.T) {
	var c Character
	c.Archetype = "The Intellectual"

	a := findArchetype(c.Archetype, archetypes, t)

	var expected CalculatedCharacter
	expected.Archetype = a.Name
	expected.Agility = a.Agility
	expected.Brawn = a.Brawn
	expected.Cunning = a.Cunning
	expected.Experience = a.Experience
	expected.Intellect = a.Intellect
	expected.Presence = a.Presence
	expected.Will = a.Will

	cc, _ := Calculate(c, archetypes, skills)

	IsExpectedCalculatedCharacter(cc, expected, t)
}

func Test_Calculate_CopiesCharacterNameToCalculatedCharacter(t *testing.T) {

	var c Character
	c.Name = "Wilberforce"
	c.Archetype = "The Intellectual"

	cc, _ := Calculate(c, archetypes, skills)

	if cc.Name != c.Name {
		t.Errorf("Expected name '%s', got '%s'", c.Name, cc.Name)
	}

}

func Test_CalculateAddsCharacterTraitsToArchetypeTraits(t *testing.T) {

	var c Character
	c.Archetype = "The Aristocrat"
	c.Agility = 1
	c.Brawn = 2
	c.Cunning = 3
	c.Experience = 4
	c.Intellect = 5
	c.Name = "Floyd"
	c.Presence = 7
	c.Will = 8

	a := findArchetype(c.Archetype, archetypes, t)

	var expected CalculatedCharacter
	expected.Archetype = c.Archetype
	expected.Agility = c.Agility + a.Agility
	expected.Brawn = c.Brawn + a.Brawn
	expected.Cunning = c.Cunning + a.Cunning
	expected.Experience = c.Experience + a.Experience
	expected.Intellect = c.Intellect + a.Intellect
	expected.Name = c.Name
	expected.Presence = c.Presence + a.Presence
	expected.Will = c.Will + a.Will

	cc, _ := Calculate(c, archetypes, skills)

	IsExpectedCalculatedCharacter(cc, expected, t)

}

func Test_CalculateGivenBogusArchetypeReturnsError(t *testing.T) {
	var c Character
	c.Archetype = "Bogus"

	_, err := Calculate(c, archetypes, skills)
	if err == nil {
		t.Errorf("Expected earth shattering kaboom.  There was no earth shattering kaboom.")
	}
	if err.Error() != "Unknown Archetype 'Bogus'" {
		t.Errorf("Expected error Unknown Archetype 'Bogus', got \"%s\"", err.Error())
	}
}

func Test_ReadCharacterFileReturnsCharacterFromFile(t *testing.T) {
	c, err := ReadCharacterFile("testcharacter.gcr")
	if err != nil {
		t.Fatalf("Error loading file: %v", err)
	}

	if c.Skills["Athletics"] != 1 {
		t.Errorf("Expected Athletics 1, got %d", c.Skills["cool"])
	}
	if c.Name != "J. Marcus Hart" {
		t.Errorf("Expected name 'J. Marcus Hart', got '%s'", c.Name)
	}
	if c.Experience != 50 {
		t.Errorf("Expected Experience 50, got %d", c.Experience)
	}
}

func Test_readCharacterFileReturnsErrorGivenBogusFile(t *testing.T) {
	c, err := ReadCharacterFile("bogus.file")
	if err == nil {
		t.Errorf("There was no Earth-shattering kaboom.  There was supposed to be an Earth-shattering kaboom.")
	}
	expectedmessage := "open bogus.file:"
	if strings.Contains(err.Error(), expectedmessage) == false {
		t.Errorf("Expected '%s', got '%s'", expectedmessage, err.Error())
	}
	if c.Name != "" {
		t.Errorf("Expected empty character, got %v", c)
	}
}

func Test_readCharacterFileReturnsErrorGivenNonCharacterFile(t *testing.T) {
	_, err := ReadCharacterFile("character_test.go")
	if err == nil {
		t.Errorf("There was no Earth-shattering kaboom.  There was supposed to be an Earth-shattering kaboom")
	}
	expectedMessage := "yaml:"
	if err.Error()[:5] != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func Test_ReadSkillFileReturnsListOfSkills(t *testing.T) {
	c, _ := ReadSkillFile("testskills.skl")

	if len(c) != 2 {
		t.Fatalf("Expected 2 skills, found %d", len(c))
	}

	first := c[0]
	if first.Name != "Athletics" || first.Ability != "Brawn" {
		t.Errorf("Expected first skill to be Athletics (Brawn), got %s (%s)", first.Name, first.Ability)
	}

	second := c[1]
	if second.Name != "Computers" || second.Ability != "Intellect" {
		t.Errorf("Expected second skill to be Computers (Intellect), got %s (%s)", second.Name, second.Ability)
	}
}

func checkSkill(actual CharacterSkill, name string, ability string, proficiencydice int, abilitydice int, t *testing.T) {
	if actual.Name != name || actual.Ability != ability || actual.ProficiencyDice != proficiencydice || actual.AbilityDice != abilitydice {
		t.Errorf("Expected %s (%s) %d/%d, got %s (%s) %d/%d",
			name, ability, proficiencydice, abilitydice,
			actual.Name, actual.Ability, actual.ProficiencyDice, actual.AbilityDice)
	}
}

func Test_CalculateGivenValidSkillsCalculatesCharacterSkills(t *testing.T) {
	c, err := ReadCharacterFile("testcharacter.gcr")
	if err != nil {
		t.Fatal(err)
	}

	cc, err := Calculate(c, archetypes, skills)

	if len(cc.Skills) != 2 {
		t.Fatalf("Expected 2 skills, got %d", len(cc.Skills))
	}

	checkSkill(cc.Skills["Athletics"], "Athletics", "Brawn", 1, 0, t)
	checkSkill(cc.Skills["Computers"], "Computers", "Intellect", 2, 2, t)

}

func findSkill(cc CalculatedCharacter, name string) CharacterSkill {
	for _, s := range cc.Skills {
		if s.Name == name {
			return s
		}
	}
	var emptySkill CharacterSkill
	return emptySkill
}

func Test_CalculateGivenSkillAssignsLevelBasedOnCorrectAttribute(t *testing.T) {
	c, err := ReadCharacterFile("testcharacter.gcr")
	if err != nil {
		t.Fatal(err)
	}
	cunningSkill := Skill{
		Name:    "Cunning Skill",
		Ability: "Cunning",
	}
	agilitySkill := Skill{
		Name:    "Agility Skill",
		Ability: "Agility",
	}
	presenceSkill := Skill{
		Name:    "Presence Skill",
		Ability: "Presence",
	}
	willSkill := Skill{
		Name:    "Will Skill",
		Ability: "Will",
	}

	skills = append(skills, cunningSkill, agilitySkill, presenceSkill, willSkill)

	c.Skills["Cunning Skill"] = 1
	c.Skills["Agility Skill"] = 2
	c.Skills["Presence Skill"] = 3
	c.Skills["Will Skill"] = 4
	c.Agility = 1
	c.Brawn = 2
	c.Cunning = 3
	c.Intellect = 4
	c.Presence = 5
	c.Will = 6

	cc, err := Calculate(c, archetypes, skills)
	if err != nil {
		t.Fatal(err)
	}

	cs := findSkill(cc, "Cunning Skill")
	if cs.ProficiencyDice != 1 || cs.ProficiencyDice+cs.AbilityDice != cc.Cunning {
		t.Errorf("Bad %s (%s: %d) %d/%d", cs.Name, cs.Ability, cc.Cunning, cs.ProficiencyDice, cs.AbilityDice)
	}
	as := findSkill(cc, "Agility Skill")
	if as.ProficiencyDice != 2 || as.ProficiencyDice+as.AbilityDice != cc.Agility {
		t.Errorf("Bad %s (%s: %d) %d/%d", as.Name, as.Ability, cc.Agility, as.ProficiencyDice, as.AbilityDice)
	}
	ps := findSkill(cc, "Presence Skill")
	if ps.ProficiencyDice != 3 || ps.ProficiencyDice+ps.AbilityDice != cc.Presence {
		t.Errorf("Bad %s (%s: %d) %d/%d", ps.Name, ps.Ability, cc.Presence, ps.ProficiencyDice, ps.AbilityDice)
	}
	ws := findSkill(cc, "Will Skill")
	if ws.ProficiencyDice != 4 || (ws.ProficiencyDice+ws.AbilityDice != cc.Will) {
		t.Errorf("Bad %s (%s: %d) %d/%d", ws.Name, ws.Ability, cc.Will, ws.ProficiencyDice, ws.AbilityDice)
	}

}

func Test_CalculateGivesDefaultsForSkillsNotPresentOnCharacter(t *testing.T) {
	c, err := ReadCharacterFile("testcharacter.gcr")
	if err != nil {
		t.Fatal(err)
	}
	presenceSkill := Skill{
		Name:    "Presence Skill",
		Ability: "Presence",
	}
	skills = append(skills, presenceSkill)
	cc, err := Calculate(c, archetypes, skills)
	s := cc.Skills["Presence Skill"]
	if s.ProficiencyDice != 0 {
		t.Errorf("Expected 0 profeciency dice, got %d", s.ProficiencyDice)
	}
	if s.AbilityDice != 4 {
		t.Errorf("Expected 4 ability dice, got %d", s.AbilityDice)
	}
}
