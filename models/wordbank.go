/*
Copyright © 2022 THINKINGATOMS LLC <atom@thinkingatoms.com>
*/

package models

import (
	"github.com/thinkingatoms/apibase/models"
	"unicode/utf8"
)

type WordBank struct {
	Words []string
	cache *models.TenureCache
	Infos map[string]*WordInfo
}

type WordInfo struct {
	Word    string
	Length  int
	NUniq   int
	Letters map[rune]bool
	Repeats map[rune]bool
}

func NewWordBank(words []string, cache *models.TenureCache) *WordBank {
	infos := make(map[string]*WordInfo)
	for _, w := range words {
		infos[w] = NewWordInfo(w)
	}
	return &WordBank{
		Words: words,
		cache: cache,
		Infos: infos,
	}
}

func NewWordInfo(word string) *WordInfo {
	letters := make(map[rune]bool)
	repeats := make(map[rune]bool)
	last := rune(0)
	for i, w := 0, 0; i < len(word); i += w {
		r, width := utf8.DecodeRuneInString(word[i:])
		letters[r] = true
		if last == r {
			repeats[r] = true
		}
		last = r
		w = width
	}
	return &WordInfo{
		Word:    word,
		Length:  utf8.RuneCountInString(word),
		NUniq:   len(letters),
		Letters: letters,
		Repeats: repeats,
	}
}

func (self *WordBank) GetSubset(runes map[rune]bool, filter func(wi *WordInfo) bool) map[string]*WordInfo {
	cacheKey := RunesMap2Key(runes)
	if ret, ok := self.cache.Get(models.TenureMedium, cacheKey); ok {
		return ret.(map[string]*WordInfo)
	}
	words := make(map[string]*WordInfo)
	if filter == nil {
		filter = func(wi *WordInfo) bool {
			return true
		}
	}
	required := 0
	for r := range runes {
		if runes[r] {
			required++
		}
	}
	for w, i := range self.Infos {
		if !filter(i) {
			continue
		}
		found := 0
		for r := range i.Letters {
			if v, ok := runes[r]; !ok {
				found = -1
				break
			} else {
				if v {
					found++
				}
			}
		}
		if found == required {
			words[w] = i
		}
	}
	self.cache.Set(models.TenureMedium, cacheKey, words)
	return words
}
