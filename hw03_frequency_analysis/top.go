package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Dict map[string]int

var reMatchWord = regexp.MustCompile(`(?m)[\p{L}\d-_]+`)

func Top10(s string) []string {
	s = strings.Trim(s, " \r\n")
	if s == "" {
		return nil
	}

	words := reMatchWord.FindAllString(s, -1)

	d := make(Dict)

	for _, w := range words {
		w = strings.ToLower(w)
		if w == "-" {
			continue
		}

		d[w]++
	}

	wordsUnique := make([]string, len(d))
	i := 0
	for w := range d {
		wordsUnique[i] = w
		i++
	}

	sort.Slice(wordsUnique, func(i, j int) bool {
		w1 := wordsUnique[i]
		w2 := wordsUnique[j]

		n1 := d[w1]
		n2 := d[w2]

		if n1 == n2 {
			return w1 < w2
		}

		return n1 > n2
	})

	if len(wordsUnique) > 10 {
		wordsUnique = wordsUnique[:10]
	}

	return wordsUnique
}
