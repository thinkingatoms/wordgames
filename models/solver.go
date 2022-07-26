package models

import (
	"github.com/thinkingatoms/apibase/models"
	"sort"
	"strings"
)

type WordGames struct {
	WordBank *WordBank
}

func WordGamesFromConfig(config map[string]any, cache *models.TenureCache) *WordGames {
	tmp := config["words"].([]any)
	words := make([]string, len(tmp))
	for i, v := range tmp {
		words[i] = v.(string)
	}
	return &WordGames{
		WordBank: NewWordBank(words, cache),
	}
}

type WordGame interface {
	Solve(wb WordBank) []string
}

func SortMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func RunesMap2Key(runes map[rune]bool) string {
	required := make([]string, 0)
	optional := make([]string, 0)
	for r, v := range runes {
		if v {
			required = append(required, string(r))
		} else {
			optional = append(optional, string(r))
		}
	}
	sort.Strings(required)
	sort.Strings(optional)
	return strings.Join(required, ",") + ":" + strings.Join(optional, ",")
}

type WordCount struct {
	Word  string
	Count int
}
