package models

import "sort"

type KillerWasp struct {
	Runes     map[rune]bool
	MinLength int
}

const DefaultKillerWaspMinLength = 4

func NewKillerWasp(required, optional []rune, minLength int) *KillerWasp {
	m := make(map[rune]bool)
	for _, r := range optional {
		m[r] = false
	}
	for _, r := range required {
		m[r] = true
	}
	if minLength == 0 {
		minLength = DefaultKillerWaspMinLength
	}
	return &KillerWasp{
		Runes:     m,
		MinLength: minLength,
	}
}

func (self *KillerWasp) Solve(wb *WordBank) []string {
	words := wb.GetSubset(self.Runes, func(wi *WordInfo) bool {
		return wi.Length >= len(self.Runes)
	})
	tmp := make([]WordCount, 0, len(words))
	for w, i := range words {
		count := 0
		for r := range self.Runes {
			if i.Letters[r] {
				count++
			}
		}
		tmp = append(tmp, WordCount{
			Word:  w,
			Count: count,
		})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Count > tmp[j].Count
	})
	ret := make([]string, 0, len(words))
	for _, v := range tmp {
		ret = append(ret, v.Word)
	}
	return ret
}
