package parser

import (
	"fmt"
	"strings"
)

type TestCase struct {
	Input, Output []string
}

func NewTestCase(input, output []string) TestCase {
	return TestCase{Input: input, Output: output}
}

func findSingleIndex(input []string, find string) (int, error) {
	var found []int
	for i, line := range input {
		if line == find {
			found = append(found, i)
		}
	}

	if len(found) == 1 {
		return found[0], nil
	} else if len(found) > 1 {
		return 0, fmt.Errorf("More than one '%s' was found", find)
	}
	return 0, fmt.Errorf("Could not find %s", find)
}

func (t *TestCase) appendToCase(start int, testcase []string, dir direction) {
	insert := &t.Input
	if dir == output {
		insert = &t.Output
	}

	for _, c := range testcase[start:] {
		switch c {
		case string(dir), "":
			continue
		case string(dir.oposite()):
			return
		}
		c = strings.ReplaceAll(c, `\`, "")
		*insert = append(*insert, c)
	}
}

func ParseTestCase(caseString string) (TestCase, error) {
	testcase := strings.Split(caseString, "\n")
	inIdx, err := findSingleIndex(testcase, string(input))
	if err != nil {
		return TestCase{}, err
	}
	outIdx, err := findSingleIndex(testcase, string(output))
	if err != nil {
		return TestCase{}, err
	}

	tc := TestCase{}
	tc.appendToCase(inIdx, testcase, input)
	tc.appendToCase(outIdx, testcase, output)

	return tc, nil
}
