package main

import (
	"caserunner/internal/parser"
	"fmt"
)

func main() {
	input := `
plain text outside tests is treated like a comment
exec: python $code
no setting time-limit results in unlimited time
time-limit: 3ms
--
output:
Fizz

input:
3
--

--
input:
3
5

output:
Fizz
Buzz

--
	`

	tf := parser.New()
	err := tf.ParseFile("main.py", input)
	if err != nil {
		fmt.Println(err)
	}
	for _, tc := range tf.Tests {
		fmt.Println(tc)
	}
	fmt.Printf("exec: %s, time limit: %d \n", tf.Exec, tf.TimeLimit)
}
