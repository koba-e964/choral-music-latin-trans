package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func discoverDirectories(configFilename string) ([]string, error) {
	directories := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		basename := filepath.Base(path)
		if !info.IsDir() && basename == configFilename {
			log.Println(path)
			dirname := filepath.Dir(path)
			directories = append(directories, dirname)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return directories, nil
}

type Config struct {
	Version     string `toml:"version"`
	FullText    string `toml:"full_text"`
	Translation string `toml:"translation"`
}

type DocData struct {
	WithMacrons    string
	WithoutMacrons string
	Translation    string
}

const configFilename = "article.toml"

func (config Config) convert() DocData {
	rule := func(r rune) rune {
		conversion := map[rune]rune{
			'ā': 'a',
			'ē': 'e',
			'ī': 'i',
			'ō': 'o',
			'ū': 'u',
			'Ā': 'A',
			'Ē': 'E',
			'Ī': 'I',
			'Ō': 'O',
			'Ū': 'U',
		}
		if val, ok := conversion[r]; ok {
			return val
		}
		return r
	}
	return DocData{
		WithMacrons:    config.FullText,
		WithoutMacrons: strings.Map(rule, config.FullText),
		Translation:    config.Translation,
	}
}

func main() {
	directories, err := discoverDirectories(configFilename)
	check(err)

	for _, directory := range directories {
		dat, err := os.ReadFile(directory + "/" + configFilename)
		if err != nil {
			log.Panic(err)
		}
		var config Config
		str := string(dat)
		toml.Decode(str, &config)
		// discover all .template files in the directory
		templates := []string{}
		err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			basename := filepath.Base(path)
			if !info.IsDir() && filepath.Ext(basename) == ".template" {
				log.Println(path)
				templates = append(templates, path)
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		for _, templateFilename := range templates {
			outputFilename := strings.TrimSuffix(templateFilename, ".template")
			content, err := os.ReadFile(templateFilename)
			if err != nil {
				log.Panic(err)
			}
			tmpl, err := template.New(directory).Parse(string(content))
			if err != nil {
				log.Panic(err)
			}
			func() {
				file, err := os.Create(outputFilename)
				if err != nil {
					log.Panic(err)
				}
				defer file.Close()
				if err := tmpl.Execute(file, config.convert()); err != nil {
					log.Panic(err)
				}
			}()
		}
	}
}
