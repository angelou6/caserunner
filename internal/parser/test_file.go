package parser

import (
	"reflect"
	"time"
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
