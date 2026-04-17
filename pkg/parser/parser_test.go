package parser

import (
	"math"
	"slices"
	"testing"
	"time"
)

func TestTestCase(t *testing.T) {
	input := `
input:
3
3
5
15

output:
Fizz
Buzz
FizzBuzz
`

	expectedInput := []string{"3", "3", "5", "15"}
	expectedOutput := []string{"Fizz", "Buzz", "FizzBuzz"}

	tc, err := ParseTestCase(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, tc.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", tc.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, tc.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", tc.Output, expectedOutput)
	}
}

func TestReverseOrder(t *testing.T) {
	input := `
output:
Fizz

input:
3
`

	expectedInput := []string{"3"}
	expectedOutput := []string{"Fizz"}

	tc, err := ParseTestCase(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, tc.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", tc.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, tc.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", tc.Output, expectedOutput)
	}
}

func TestEmpty(t *testing.T) {
	_, err := ParseTestCase("")
	if err == nil {
		t.FailNow()
	}
}

func TestIgnore(t *testing.T) {
	input := `
input:
3
output:
\output: Fizz
`

	expectedInput := []string{"3"}
	expectedOutput := []string{"output: Fizz"}

	tc, err := ParseTestCase(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, tc.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", tc.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, tc.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", tc.Output, expectedOutput)
	}
}

func TestInvalidToken(t *testing.T) {
	input := `
break_here:
3

output:
Fizz
	`
	_, err := ParseTestCase(input)
	if err == nil {
		t.FailNow()
	}
}

func TestMultipleTokens(t *testing.T) {
	duplicatedInput := `
input:
3

input:
3

output:
Fizz
	`
	_, err := ParseTestCase(duplicatedInput)
	if err == nil {
		t.Error("Duplicated input uncaught")
	}

	duplicatedOutput := `
input:
3

output:
Fizz

output:
Fizz
	`

	_, err = ParseTestCase(duplicatedOutput)
	if err == nil {
		t.Error("Duplicated output uncaught")
	}
}

func TestMissingToken(t *testing.T) {
	missingOutput := `
input:
3
	`
	_, err := ParseTestCase(missingOutput)
	if err == nil {
		t.Error("Missing output not caught")
	}

	missingInput := `
output:
Fizz
	`

	_, err = ParseTestCase(missingInput)
	if err == nil {
		t.Error("Missing input uncaught")
	}
}

func TestTestFile(t *testing.T) {
	input := `
plain text outside tests is treated like a comment
exec: python $code
time-limit: 3ms
--
input:
3
output:
Fizz
--
	`

	suite, err := ParseFile(input, "main.py")
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if suite.Exec != "python main.py" {
		t.Errorf("Wrong exec. Expected 'python main.py', got %s", suite.Exec)
	}

	if suite.TimeLimit != 3*time.Millisecond {
		t.Errorf("Wrong time limit. Expected '3ms', got %v", suite.TimeLimit)
	}
}

func TestMultipleTests(t *testing.T) {
	input := `
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

	suite, err := ParseFile(input, "main.py")
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if len(suite.Cases) != 2 {
		t.Fatalf("Expected 2 test cases, got %d", len(suite.Cases))
	}
}

func TestMissingConfig(t *testing.T) {
	missingTime := `
exec: python $code
--
input:
3
output:
Fizz
--
	`

	suite, err := ParseFile(missingTime, "main.py")
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}
	if suite.TimeLimit != time.Duration(math.MaxInt64) {
		t.Errorf("Expected infinite time limit, got %v", suite.TimeLimit)
	}

	missingExec := `
time-limit: 3ms
--
input:
3
output:
Fizz
--
	`

	_, err = ParseFile(missingExec, "main.py")
	if err == nil {
		t.Error("Missing exec not caught")
	}
}

func TestEmptyFile(t *testing.T) {
	_, err := ParseFile("", "main.py")
	if err == nil {
		t.FailNow()
	}
}
