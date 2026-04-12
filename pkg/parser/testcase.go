package parser

import (
	"reflect"
	"strings"
	"time"
)

type direction string

const (
	Input  direction = "input:"
	Output direction = "output:"
)

type TestCase struct {
	Input  []string
	Output []string
}

type TestFile struct {
	Exec      string
	TimeLimit time.Duration
	Tests     []TestCase
}

func (d direction) Opposite() direction {
	if d == Input {
		return Output
	}
	return Input
}

func New() *TestFile {
	return &TestFile{}
}

func (t *TestCase) JudgeOutput(output []string) bool {
	expected := []string{}
	for _, o := range t.Output {
		if o != "" {
			expected = append(expected, o)
		}
	}
	return reflect.DeepEqual(expected, output)
}

func (t *TestCase) AppendToCase(indexes int, caseInputs []string, into direction) {
	for i := indexes; i < len(caseInputs); i++ {
		if len(caseInputs[i]) == 0 || caseInputs[i] == string(into) {
			continue
		} else if caseInputs[i] == string(into.Opposite()) {
			break
		}

		s := strings.TrimSpace(caseInputs[i])
		s = strings.Replace(s, "\\", "", 1)

		if into == Input {
			t.Input = append(t.Input, s)
		} else {
			t.Output = append(t.Output, s)
		}
	}
}
