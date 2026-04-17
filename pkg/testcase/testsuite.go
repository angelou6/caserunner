package testcase

import "time"

type TestSuite struct {
	Exec      string
	TimeLimit time.Duration
	Cases     []TestCase
}

func NewSuite(exec string, limit time.Duration, cases []TestCase) TestSuite {
	return TestSuite{Exec: exec, TimeLimit: limit, Cases: cases}
}
