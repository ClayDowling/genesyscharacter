package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"github.com/claydowling/genesyscharacter/genesys"
)

func getPath(filename string) string {
	base, _ := os.Executable()
	startdir := filepath.Join(filepath.Dir(base), "data")
	return filepath.Join(startdir, filename)
}

func main() {
	archetypes, err := genesys.ReadArchetypeFile(getPath("base/archetypes.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	skills, err := genesys.ReadSkillFile(getPath("base/skills.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	character, err := genesys.ReadCharacterFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	sheet, err := ioutil.ReadFile(getPath("base/sheet.txt"))

	cc, err := genesys.Calculate(character, archetypes, skills)

	s := template.Must(template.New("sheet").Parse(string(sheet)))
	s.Execute(os.Stdout, cc)
}
