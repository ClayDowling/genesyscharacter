package genesys

import (
	"os"
	"testing"
)

var archetypes []Archetype

func TestMain(m *testing.M) {
	archetypes = ReadArchetypeFile("testfile.arc")
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

func IsExpectedCalculatedCharacter(cc *CalculatedCharacter, expected *CalculatedCharacter, t *testing.T) {
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
	if cc.Intelligence != expected.Intelligence {
		t.Errorf("Expected Intelligence %d, got %d", expected.Intelligence, cc.Intelligence)
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
	expected.Intelligence = a.Intelligence
	expected.Presence = a.Presence
	expected.Will = a.Will

	cc, _ := Calculate(c, archetypes)

	IsExpectedCalculatedCharacter(cc, &expected, t)
}

func Test_Calculate_CopiesCharacterNameToCalculatedCharacter(t *testing.T) {

	var c Character
	c.Name = "Wilberforce"
	c.Archetype = "The Intellectual"

	cc, _ := Calculate(c, archetypes)

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
	c.Intelligence = 5
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
	expected.Intelligence = c.Intelligence + a.Intelligence
	expected.Name = c.Name
	expected.Presence = c.Presence + a.Presence
	expected.Will = c.Will + a.Will

	cc, _ := Calculate(c, archetypes)

	IsExpectedCalculatedCharacter(cc, &expected, t)

}

func Test_CalculateGivenBogusArchetypeReturnsError(t *testing.T) {
	var c Character
	c.Archetype = "Bogus"

	_, err := Calculate(c, archetypes)
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

	if c.Skills["cool"] != 3 {
		t.Errorf("Expected Cool 3, got %d", c.Skills["cool"])
	}
	if c.Name != "J. Marcus Hart" {
		t.Errorf("Expected name 'J. Marcus Hart', got '%s'", c.Name)
	}
	if c.Experience != 50 {
		t.Errorf("Expected Experience 50, got %d", c.Experience)
	}
}
