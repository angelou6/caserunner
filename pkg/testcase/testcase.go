package testcase

import (
	"strings"
)

type TestCase struct {
	Input, Output []string
}

func NewTestCase(input, output []string) TestCase {
	return TestCase{Input: input, Output: output}
}

func (t *TestCase) AppendToCase(start int, testcase []string, dir Direction) {
	insert := &t.Input
	if dir == Output {
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
