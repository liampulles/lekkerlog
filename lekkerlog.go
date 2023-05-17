package lekkerlog

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func Prettify(jsonLine []byte) string {
	// Unmarshal JSON
	var m map[string]interface{}
	if err := json.Unmarshal(jsonLine, &m); err != nil {
		// Can't unmarshal -> resort to standard output.
		// TODO: Create a good secondary formatter
		return string(jsonLine)
	}

	// As sorted pairs
	pairs := toPairs(m)
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Key < pairs[j].Key })

	// Format
	yellow := color.New(color.FgYellow).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	var s strings.Builder
	for i, pair := range pairs {
		s.WriteString(yellow(strings.ToUpper(pair.Key)))
		s.WriteString(white("="))
		s.WriteString(white(pair.Value))
		// Not last one?
		if i != len(pairs)-1 {
			s.WriteString(white("\t|"))
		}
	}
	s.WriteByte('\n')

	return s.String()
}

type pair struct {
	Key   string
	Value interface{}
}

func toPairs(m map[string]interface{}) []pair {
	pairs := make([]pair, len(m))
	i := -1
	for k, v := range m {
		i++
		pairs[i] = pair{k, v}
	}
	return pairs
}
