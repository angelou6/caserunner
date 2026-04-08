package parser

import (
	"slices"
	"testing"
)

func TestTestCase(t *testing.T) {
	input :=
		`input:
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

	cases, err := parseTest(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, cases.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", cases.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, cases.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", cases.Output, expectedOutput)
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

	cases, err := parseTest(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, cases.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", cases.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, cases.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", cases.Output, expectedOutput)
	}
}

func TestEmpty(t *testing.T) {
	_, error := parseTest("")
	if error == nil {
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

	cases, err := parseTest(input)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if !slices.Equal(expectedInput, cases.Input) {
		t.Errorf("Inputs are different, Got: %v, expected %v", cases.Input, expectedInput)
	}
	if !slices.Equal(expectedOutput, cases.Output) {
		t.Errorf("Outputs are different, Got: %v, expected %v", cases.Output, expectedOutput)
	}
}

func TestInvalidToken(t *testing.T) {
	input := `
break_here:
3

output:
Fizz
	`
	_, error := parseTest(input)
	if error == nil {
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
	_, error := parseTest(duplicatedInput)
	if error == nil {
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

	_, error = parseTest(duplicatedOutput)
	if error == nil {
		t.Error("Duplicated output uncaught")
	}
}

func TestMissingToken(t *testing.T) {
	missingOutput := `
input:
3
	`
	_, error := parseTest(missingOutput)
	if error == nil {
		t.Error("Missing output not cought")
	}

	missingInput := `
output:
Fizz
	`

	_, error = parseTest(missingInput)
	if error == nil {
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

	testFile := New()
	error := testFile.ParseFile(input, "main.py")

	if error != nil {
		t.Errorf("Got error: %v", error)
	}

	if testFile.Exec != "python main.py" {
		t.Errorf("Wrong exec. Expected 'python main.py', got %s", testFile.Exec)
	}

	if testFile.TimeLimit != 3_000_000 {
		t.Errorf("Wrong time limit. Expected '3,000,000 nanosecons', got %d", testFile.TimeLimit)
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

	testFile := New()
	error := testFile.ParseFile(input, "main.py")

	if error != nil {
		t.Errorf("Got error: %v", error)
	}

	if len(testFile.Tests) != 2 {
		t.FailNow()
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

	testFile := New()
	error := testFile.ParseFile(missingTime, "main.py")

	if error != nil {
		t.Errorf("Got error: %v", error)
	}

	if testFile.TimeLimit != -1 {
		t.Errorf("Wrong time limit. Expected '-1', got %d", testFile.TimeLimit)
	}

	missingCode := `
--
input:
3
output:
Fizz
--
	`

	testFile = New()
	error = testFile.ParseFile(missingCode, "main.py")

	if error == nil {
		t.Error("Exec missing uncaught")
	}

	emptyCode := `
--
input:
3
output:
Fizz
--
	`

	testFile = New()
	error = testFile.ParseFile(emptyCode, "main.py")

	if error == nil {
		t.Error("Exec empty uncaught")
	}
}

func TestEmptyFile(t *testing.T) {
	testFile := New()
	error := testFile.ParseFile("", "main.py")

	if error == nil {
		t.FailNow()
	}
}
