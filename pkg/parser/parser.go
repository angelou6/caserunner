package parser

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

type TestSuite struct {
	exec      string
	timeLimit time.Duration
	Cases     []TestCase
}

func NewTestSuite(exec string, limit time.Duration) TestSuite {
	return TestSuite{exec: exec, timeLimit: limit}
}

func firstRegexMatch(regex, input string) (string, error) {
	r := regexp.MustCompile(regex)
	if m := r.FindStringSubmatch(input); m != nil {
		return m[1], nil
	}
	return "", fmt.Errorf("Regex: %s not found in %s", regex, input)
}

func ParseFile(input, fileLocation string) (TestSuite, error) {
	exec, err := firstRegexMatch(`exec:\s*(.+)`, input)
	exec = strings.ReplaceAll(exec, "$code", fileLocation)
	if err != nil {
		return TestSuite{}, err
	}

	var limit time.Duration
	timeLimit, err := firstRegexMatch(`time-limit:\s*(\S+)`, input)
	if err != nil {
		limit = time.Duration(math.MaxInt64)
	} else {
		limit, err = time.ParseDuration(timeLimit)
		if err != nil {
			return TestSuite{}, err
		}
	}

	cases := regexp.MustCompile(`(?s)--\n(.*?)\n--`).FindAllStringSubmatch(input, -1)
	suite := NewTestSuite(exec, limit)
	for _, c := range cases {
		tc, err := ParseTestCase(c[1])
		if err != nil {
			return TestSuite{}, err
		}

		suite.Cases = append(suite.Cases, tc)
	}

	return suite, nil
}
