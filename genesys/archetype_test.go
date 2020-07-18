package genesys

import (
	"testing"
)

func Test_ReadFileReturnsCorrectNumberOfArchetypes(t *testing.T) {
	actual := ReadFile("testfile.arc")
	if len(actual) != 2 {
		t.Errorf("Expected 2 archetypes, found %d", len(actual))
	}
}

func does_archetype_match(expected Archetype, actual Archetype, t *testing.T) {
	if actual.Name != expected.Name {
		t.Errorf("Expected name '%s', got '%s'", expected.Name, actual.Name)
	}
	if actual.Brawn != expected.Brawn {
		t.Errorf("Expected brawn %d, got %d", expected.Brawn, actual.Brawn)
	}
	if actual.Agility != expected.Agility {
		t.Errorf("Expected agility %d, got %d", expected.Agility, actual.Agility)
	}
	if actual.Intelligence != expected.Intelligence {
		t.Errorf("Expected intelligence %d, got %d", expected.Intelligence, actual.Intelligence)
	}
	if actual.Cunning != expected.Cunning {
		t.Errorf("Expected cunning %d, got %d", expected.Cunning, actual.Cunning)
	}
	if actual.Will != expected.Will {
		t.Errorf("Expected will %d, got %d", expected.Will, actual.Will)
	}
	if actual.Presence != expected.Presence {
		t.Errorf("Expected presence %d, got %d", expected.Presence, actual.Presence)
	}
	if actual.Wound != expected.Wound {
		t.Errorf("Expected wound %d, got %d", expected.Wound, actual.Wound)
	}
	if actual.Strain != expected.Strain {
		t.Errorf("Expected strain %d, got %d", expected.Strain, actual.Strain)
	}
	if actual.Experience != expected.Experience {
		t.Errorf("Expected experience %d, got %d", expected.Experience, actual.Experience)
	}

}

func Test_ReadFileReturnsExpectedArchetypes(t *testing.T) {
	actual := ReadFile("testfile.arc")

	intellectual := Archetype{
		Name:         "The Intellectual",
		Brawn:        2,
		Agility:      1,
		Intelligence: 3,
		Cunning:      2,
		Will:         2,
		Presence:     2,
		Wound:        8,
		Strain:       12,
		Experience:   100}
	aristocrat := Archetype{
		Name:         "The Aristocrat",
		Brawn:        1,
		Agility:      2,
		Intelligence: 2,
		Cunning:      2,
		Will:         2,
		Presence:     3,
		Wound:        10,
		Strain:       10,
		Experience:   100}

	does_archetype_match(intellectual, actual[0], t)
	does_archetype_match(aristocrat, actual[1], t)
}

func Test_ReadFileReturnsEmptyListWhenBadFileName(t *testing.T) {
	actual := ReadFile("bogus.arc")
	if len(actual) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(actual))
	}
}

func Test_ReadFileReturnsEmptyListWhenGetsNonArchetypeFile(t *testing.T) {
	actual := ReadFile("archetype_test.go")
	if len(actual) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(actual))
	}
}
