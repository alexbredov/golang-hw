package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(line string) []string {
	if line == "" {
		return []string{}
	}
	r := regexp.MustCompile(`\s+`)
	line2 := r.ReplaceAllString(line, " ")
	pairs := rankByWordCounter(line2)
	keys := []string{}
	for _, pair := range pairs[:10] {
		keys = append(keys, pair.Key)
	}
	return keys
}

type Pair struct {
	Key   string
	Value int
}

func rankByWordCounter(line string) []Pair {
	wordsKey := make(map[string]int)
	words := strings.Split(line, " ")
	for _, word := range words {
		wordsKey[word]++
	}
	sortedslice := []Pair{}
	for k, v := range wordsKey {
		sortedslice = append(sortedslice, Pair{k, v})
	}
	sort.Slice(sortedslice, func(i, j int) bool {
		if sortedslice[i].Value < sortedslice[j].Value {
			return false
		}
		if sortedslice[i].Value == sortedslice[j].Value {
			return sortedslice[i].Key < sortedslice[j].Key
		}
		return true
	})
	return sortedslice
}
