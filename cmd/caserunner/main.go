package main

import (
	"caserunner/internal/parser"
	"caserunner/internal/runner"
	"fmt"
	"os"
)

func main() {
	input := `
exec: python $code
time-limit: 22ms

--
input:
3

output:
Fizz
--

--
input:
1

output:
1
--
`
	testfile := parser.New()
	err := testfile.ParseFile(input, "internal/runner/tests/fizzbuzz.py")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runner.RunFile(*testfile)
}
