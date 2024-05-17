package decl

import (
	"fmt"
	"log"
)

type Adj struct {
	Original       string            `toml:"original_form"`
	DeclensionType int               `toml:"declension_type"`
	Translation    string            `toml:"translation"`
	Explanation    string            `toml:"explanation"`
	Declensions    map[string]string `toml:"declensions"`
}

// AdjFactory は名詞の活用について説明する関数を返す。
func AdjFactory(adjs []Adj) func(adjOriginal string, case_ string) string {
	return func(adjOriginal string, case_ string) string {
		adjEntry := Adj{}
		declined := ""
		adjEntry.DeclensionType = -1
		for _, adj := range adjs {
			if adj.Original == adjOriginal {
				if val, ok := adj.Declensions[case_]; ok {
					adjEntry = adj
					declined = val
				}
			}
		}
		if adjEntry.DeclensionType == -1 {
			panic("unknown noun")
		}

		declensionText := ""
		switch adjEntry.DeclensionType {
		case 12:
			declensionText = "第1・第2変化形容詞"
		case 3:
			declensionText = "第3変化形容詞"
		default:
			log.Panicf("unknown declension type: %d", adjEntry.DeclensionType)
		}
		caseText := getAdjCaseText(case_)
		return fmt.Sprintf("`%s`は%s`%s`(%s)の%sです。", declined, declensionText, adjEntry.Explanation, adjEntry.Translation, caseText)
	}
}

func getAdjCaseText(case_ string) string {
	if len(case_) <= 1 || len(case_) >= 4 {
		panic("unknown case")
	}
	caseText := RuneToCase(rune(case_[0]))
	gender := ""
	switch case_[1] {
	case 'm':
		gender = "男性"
	case 'f':
		gender = "女性"
	case 'n':
		gender = "中性"
	default:
		panic("unknown gender")
	}
	number := "単数"
	if len(case_) == 3 {
		if case_[2] == 'p' {
			number = "複数"
		} else {
			panic("unknown number")
		}
	}
	return fmt.Sprintf("%s・%s・%s", gender, number, caseText)
}
