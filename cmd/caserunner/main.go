package main

import (
	"caserunner/pkg/parser"
	"caserunner/pkg/runner"
	"fmt"
	"os"
)

func main() {
	input := `
this is a test comment
exec: python $code
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
	file, err := parser.ParseFile(input, os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res := runner.RunSuite(file)
	res.PrintResults()
}
