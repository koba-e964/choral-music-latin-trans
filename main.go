package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

type Noun struct {
	Original       string            `toml:"original_form"`
	DeclensionType int               `toml:"declension_type"`
	Translation    string            `toml:"translation"`
	Explanation    string            `toml:"explanation"`
	Declensions    map[string]string `toml:"declensions"`
}

type Prep struct {
	Original    string `toml:"original_form"`
	Translation string `toml:"translation"`
	Takes       string `toml:"takes"`
}

type Config struct {
	Version     string            `toml:"version"`
	FullText    string            `toml:"full_text"`
	Translation string            `toml:"translation"`
	Lines       map[string]string `toml:"lines"`
	Nouns       []Noun            `toml:"nouns"`
	Preps       []Prep            `toml:"preps"`
}

type DocData struct {
	WithMacrons    string
	WithoutMacrons string
	Translation    string
	Lines          map[string]string
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
	tmpl, err := template.New("fulltext").Parse(config.FullText)
	if err != nil {
		log.Panic(err)
	}
	docData := DocData{
		Translation: config.Translation,
		Lines:       config.Lines,
	}
	buf := strings.Builder{}
	err = tmpl.Execute(&buf, docData)
	if err != nil {
		log.Panic(err)
	}
	fulltext := buf.String()
	docData.WithMacrons = fulltext
	docData.WithoutMacrons = strings.Map(rule, fulltext)
	return docData
}

func runeToCase(r rune) string {
	caseText := ""
	switch r {
	case '1':
		caseText = "主格"
	case '2':
		caseText = "属格"
	case '3':
		caseText = "与格"
	case '4':
		caseText = "対格"
	case 'a':
		caseText = "奪格"
	case 'v':
		caseText = "呼格"
	default:
		panic("unknown case")
	}
	return caseText
}

func getNounCaseText(case_ string) string {
	if len(case_) <= 1 || len(case_) >= 4 {
		panic("unknown case")
	}
	caseText := runeToCase(rune(case_[0]))
	number := "単数"
	if len(case_) == 3 {
		if case_[2] == 'p' {
			number = "複数"
		} else {
			panic("unknown number")
		}
	}
	return fmt.Sprintf("%s・%s", number, caseText)
}

func nounFactory(config Config) func(nounOriginal string, case_ string) string {
	return func(nounOriginal string, case_ string) string {
		nounEntry := Noun{}
		declined := ""
		nounEntry.DeclensionType = -1
		for _, noun := range config.Nouns {
			if noun.Original == nounOriginal {
				if val, ok := noun.Declensions[case_]; ok {
					nounEntry = noun
					declined = val
				}
			}
		}
		if nounEntry.DeclensionType == -1 {
			panic("unknown noun")
		}

		declensionText := ""
		switch nounEntry.DeclensionType {
		case 1, 2, 3, 4, 5:
			declensionText = "第" + strconv.Itoa(nounEntry.DeclensionType) + "変化名詞"
		case 0:
			declensionText = "不変化名詞"
		default:
			log.Panicf("unknown declension type: %d", nounEntry.DeclensionType)
		}
		caseText := getNounCaseText(case_)
		return fmt.Sprintf("`%s`は%s`%s`(%s)の%sです。", declined, declensionText, nounEntry.Explanation, nounEntry.Translation, caseText)
	}
}

func prepFactory(config Config) func(prepOriginal string) string {
	return func(prepOriginal string) string {
		for _, prep := range config.Preps {
			if prep.Original == prepOriginal {
				caseText := ""
				for i, case_ := range prep.Takes {
					if i > 0 {
						caseText += "または"
					}
					caseText += runeToCase(case_)
				}
				return fmt.Sprintf("`%s`は%sをとる前置詞(%s)です。", prepOriginal, caseText, prep.Translation)
			}
		}
		panic("unknown preposition")
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
			tmpl, err := template.New(directory).Funcs(template.FuncMap{"noun": nounFactory(config), "prep": prepFactory(config)}).Parse(string(content))
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
