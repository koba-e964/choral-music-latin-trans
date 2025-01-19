package main

import (
	"fmt"
)

type Conj struct {
	Original    string `toml:"original_form"`
	Translation string `toml:"translation"`
}

func conjFactory(conjs []Conj) func(prepOriginal string) string {
	return func(conjOriginal string) string {
		for _, conj := range conjs {
			if conj.Original == conjOriginal {
				return fmt.Sprintf("`%s`は接続詞(%s)です。", conjOriginal, conj.Translation)
			}
		}
		panic("unknown conjunction")
	}
}
