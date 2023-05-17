package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/liampulles/lekkerlog"
)

func main() {
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
