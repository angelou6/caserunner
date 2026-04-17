package testcase

import (
	"caserunner/internal/direction"
	"strings"
)

type TestCase struct {
	Input, Output []string
}

func NewTestCase(input, output []string) TestCase {
	return TestCase{Input: input, Output: output}
}

func (t *TestCase) AppendToCase(start int, testcase []string, dir direction.Direction) {
	insert := &t.Input
	if dir == direction.Output {
		insert = &t.Output
	}

	for _, c := range testcase[start:] {
		switch c {
		case string(dir), "":
			continue
		case string(dir.Oposite()):
			return
		}

		c = strings.ReplaceAll(c, `\`, "")
		if c == "" {
			continue
		}
		*insert = append(*insert, c)
	}
}
