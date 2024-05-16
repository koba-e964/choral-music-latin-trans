package main

import (
	"fmt"

	"github.com/koba-e964/choral-music-latin-trans/decl"
)

type Prep struct {
	Original    string `toml:"original_form"`
	Translation string `toml:"translation"`
	Takes       string `toml:"takes"`
}

func prepFactory(preps []Prep) func(prepOriginal string) string {
	return func(prepOriginal string) string {
		for _, prep := range preps {
			if prep.Original == prepOriginal {
				caseText := ""
				for i, case_ := range prep.Takes {
					if i > 0 {
						caseText += "または"
					}
					caseText += decl.RuneToCase(case_)
				}
				return fmt.Sprintf("`%s`は%sをとる前置詞(%s)です。", prepOriginal, caseText, prep.Translation)
			}
		}
		panic("unknown preposition")
	}
}
