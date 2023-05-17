package lekkerlog

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

func Prettify(jsonLine []byte) string {
	l, err := parse(jsonLine)
	if err != nil {
		// Can't parse -> resort to standard output.
		// TODO: Create a good secondary formatter
		return fmt.Sprint(red(err.Error()), string(jsonLine), "\n")
	}
	return format(l)
}

type log struct {
	Level   string  `json:"level"`
	Message string  `json:"message"`
	Time    LogTime `json:"time"`
	More    map[string]interface{}
}

func parse(jsonLine []byte) (l log, err error) {
	// Base
	err = json.Unmarshal(jsonLine, &l)
	if err != nil {
		return
	}

	// More
	json.Unmarshal(jsonLine, &l.More)
	delete(l.More, "level")
	delete(l.More, "message")
	delete(l.More, "time")

	return
}

var (
	white        = color.New(color.FgWhite).SprintFunc()
	boldWhite    = color.New(color.FgWhite, color.Bold).SprintFunc()
	ulineWhite   = color.New(color.FgWhite, color.Underline).SprintFunc()
	boldGreen    = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow       = color.New(color.FgYellow).SprintFunc()
	yellowOnRed  = color.New(color.FgYellow, color.BgRed).SprintFunc()
	red          = color.New(color.FgRed).SprintFunc()
	blackOnRed   = color.New(color.FgBlack, color.BgRed).SprintFunc()
	blackOnWhite = color.New(color.FgBlack, color.BgWhite, color.Bold).SprintFunc()
	cyan         = color.New(color.FgCyan).SprintFunc()
	blue         = color.New(color.FgBlue).SprintFunc()
	boldMagenta  = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

var emptyTime = time.Time{}

const timeFormat = "2006/01/02 15:04:05"

func format(l log) string {
	var segs []string

	t := time.Time(l.Time)
	if t != emptyTime {
		segs = append(segs, blue(t.Local().Format(timeFormat)))
	}
	if l.Level != "" {
		segs = append(segs, boldWhite("[")+formatLevel(l.Level)+boldWhite("]"))
	}
	if l.Message != "" {
		segs = append(segs, boldWhite(l.Message))
	}

	if len(l.More) > 0 {
		moreSegs := make([]string, len(l.More))
		i := -1
		for k, v := range l.More {
			i++
			if m, ok := v.(map[string]interface{}); ok {
				j, _ := json.Marshal(m)
				v = string(j)
			}
			moreSegs[i] = fmt.Sprintf("%s=%v", boldGreen(k), v)
		}
		sort.Strings(moreSegs)
		segs = append(segs, boldMagenta("|>"), strings.Join(moreSegs, " "))
	}

	return strings.Join(segs, " ") + "\n"
}

func formatLevel(lvl string) string {
	// TODO: Find other ways / more options for categorizing.
	switch upt := strings.ToUpper(strings.TrimSpace(lvl)); upt {
	case "TRACE", "TRC", "TRCE":
		return boldWhite("TRC")
	case "DEBUG", "DBG", "DBUG":
		return yellow("DBG")
	case "INFO", "INF":
		return boldGreen("INF")
	case "WARN", "WRN", "WARNING":
		return yellowOnRed("WRN")
	case "ERROR", "ERR", "ERRO", "E":
		return red("ERR")
	case "FATAL", "FTL", "FATALITY", "FTLERROR":
		return blackOnRed("FTL")
	default:
		// Default case
		return boldWhite(upt[:3])
	}
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

type LogTime time.Time

func (lt *LogTime) UnmarshalJSON(data []byte) error {
	// Try as unix time
	i, err := unmarshalBigInt(data)
	if err == nil {
		// Could be different levels of precision
		t := time.UnixMicro(i.Int64())
		if reasonableTime(t) {
			*lt = LogTime(t)
			return nil
		}
		t = time.UnixMilli(i.Int64())
		if reasonableTime(t) {
			*lt = LogTime(t)
			return nil
		}
		t = time.Unix(i.Int64(), 0)
		if reasonableTime(t) {
			*lt = LogTime(t)
			return nil
		}
	}

	// No other options - ignore then.
	return nil
}

func unmarshalBigInt(data []byte) (big.Int, error) {
	var z big.Int
	_, ok := z.SetString(string(data), 10)
	if !ok {
		return z, fmt.Errorf("not a valid big integer: %s", string(data))
	}
	return z, nil
}

const year = time.Hour * 24 * 365

func reasonableTime(t time.Time) bool {
	// In the last 10 years either side
	now := time.Now()
	past, future := now.Add(-10*year), now.Add(10*year)
	return past.Before(t) && future.After(t)
}
