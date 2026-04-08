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
	return reflect.DeepEqual(t.Output, output)
}
