package main

import (
	"caserunner/pkg/parser"
	"fmt"
	"os"
)

func main() {
	input := `
this is a test comment
exec: python $code
time-limit: 3ms
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
	file, err := parser.ParseFile(input, "main.py")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(file)
}
