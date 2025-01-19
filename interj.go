package main

import (
	"fmt"
)

type Interj struct {
	Original    string `toml:"original_form"`
	Translation string `toml:"translation"`
}

func interjFactory(interjs []Interj) func(prepOriginal string) string {
	return func(interjOriginal string) string {
		for _, interj := range interjs {
			if interj.Original == interjOriginal {
				return fmt.Sprintf("`%s`は間投詞(%s)です。", interjOriginal, interj.Translation)
			}
		}
		panic("unknown interjection")
	}
}
