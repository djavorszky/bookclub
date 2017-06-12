package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

var out io.Writer = os.Stdout

func echoTwo(args []string) {
	var s, sep string
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
	fmt.Fprintln(out, s)
}

func echoThree(args []string) {
	fmt.Fprintln(out, strings.Join(args, " "))
}
