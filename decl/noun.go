package decl

import (
	"fmt"
	"log"
	"strconv"
)

type Noun struct {
	Original       string            `toml:"original_form"`
	DeclensionType int               `toml:"declension_type"`
	Translation    string            `toml:"translation"`
	Explanation    string            `toml:"explanation"`
	Declensions    map[string]string `toml:"declensions"`
}

// NounFactory は名詞の活用について説明する関数を返す。
func NounFactory(nouns []Noun) func(nounOriginal string, case_ string) string {
	return func(nounOriginal string, case_ string) string {
		nounEntry := Noun{}
		declined := ""
		nounEntry.DeclensionType = -1
		for _, noun := range nouns {
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

func getNounCaseText(case_ string) string {
	if len(case_) <= 1 || len(case_) >= 4 {
		panic("unknown case")
	}
	caseText := RuneToCase(rune(case_[0]))
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
