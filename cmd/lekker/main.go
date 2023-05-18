package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/liampulles/lekkerlog"
)

func main() {
	// Only continue if being piped in
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Try piping in some JSON logs to lekker! zerolog format is ideal.")
		os.Exit(0)
	}

	// Stream input, process, write output.
	r := bufio.NewReader(os.Stdin)
	s, e := readln(r)
	for e == nil {
		out := lekkerlog.Prettify(s)
		fmt.Fprint(color.Error, out)
		s, e = readln(r)
	}
}

func readln(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return ln, err
}
