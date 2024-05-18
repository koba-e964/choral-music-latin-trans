package main

import (
	"fmt"
	"log"
	"strconv"
)

type Verb struct {
	Original        string                       `toml:"original_form"`
	ConjugationType int                          `toml:"conjugation_type"`
	Translation     string                       `toml:"translation"`
	Explanation     string                       `toml:"explanation"`
	Conjugations    map[string]map[string]string `toml:"conjugations"`
}

func getVerbTenseText(tense string, kind string) string {
	tenseText := ""
	switch tense {
	case "inf/pres":
		tenseText = "不定法・能動態・現在"
	case "inf/perf":
		tenseText = "不定法・能動態・完了"
	default:
		panic("unknown tense")
	}
	kindText := ""
	switch kind {
	case "":
		kindText = ""
	default:
		panic("unknown kind")
	}
	return fmt.Sprintf("%s%s", tenseText, kindText)
}

func verbFactory(verbs []Verb) func(verbOriginal string, tense string, kind string) string {
	return func(verbOriginal string, tense string, kind string) string {
		verbEntry := Verb{}
		conjugated := ""
		verbEntry.ConjugationType = -1
		for _, verb := range verbs {
			if verb.Original == verbOriginal {
				if val, ok := verb.Conjugations[tense]; ok {
					if val2, ok2 := val[kind]; ok2 {
						verbEntry = verb
						conjugated = val2
					}
				}
			}
		}
		if verbEntry.ConjugationType == -1 {
			panic("unknown noun")
		}

		conjugationText := ""
		switch verbEntry.ConjugationType {
		case 1, 2, 3, 4:
			conjugationText = "第" + strconv.Itoa(verbEntry.ConjugationType) + "変化動詞"
		case 0:
			conjugationText = "不規則動詞"
		case 5:
			conjugationText = "形式受動相動詞"
		default:
			log.Panicf("unknown conjugation type: %d", verbEntry.ConjugationType)
		}
		caseText := getVerbTenseText(tense, kind)
		return fmt.Sprintf("`%s`は%s`%s`(%s)の%sです。", conjugated, conjugationText, verbEntry.Explanation, verbEntry.Translation, caseText)
	}
}
