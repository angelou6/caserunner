package parser

import (
	"caserunner/internal/direction"
	"caserunner/pkg/testcase"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

func firstRegexMatch(regex, input string) (string, error) {
	r := regexp.MustCompile(regex)
	if m := r.FindStringSubmatch(input); m != nil {
		return m[1], nil
	}
	return "", fmt.Errorf("%s no fué encontrado en %s", regex, input)
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
		return 0, fmt.Errorf("Más de un '%s' fué encontrado", find)
	}
	return 0, fmt.Errorf("No se pudo encontrar %s", find)
}

func ParseTestCase(caseString string) (testcase.TestCase, error) {
	lines := strings.Split(caseString, "\n")
	inIdx, err := findSingleIndex(lines, string(direction.Input))
	if err != nil {
		return testcase.TestCase{}, err
	}
	outIdx, err := findSingleIndex(lines, string(direction.Output))
	if err != nil {
		return testcase.TestCase{}, err
	}

	tc := testcase.TestCase{}
	tc.AppendToCase(inIdx, lines, direction.Input)
	tc.AppendToCase(outIdx, lines, direction.Output)

	return tc, nil
}

func ParseFile(file, code string) (testcase.TestSuite, error) {
	exec, err := firstRegexMatch(`exec:\s*(.+)`, file)
	exec = strings.ReplaceAll(exec, "$code", code)
	if err != nil {
		return testcase.TestSuite{}, err
	}

	var limit time.Duration
	timeLimit, err := firstRegexMatch(`time-limit:\s*(\S+)`, file)
	if err != nil {
		limit = time.Duration(math.MaxInt64)
	} else {
		limit, err = time.ParseDuration(timeLimit)
		if err != nil {
			return testcase.TestSuite{}, err
		}
	}

	cases := regexp.MustCompile(`(?s)--\n(.*?)\n--`).FindAllStringSubmatch(file, -1)
	suite := testcase.TestSuite{Exec: exec, TimeLimit: limit}
	for _, c := range cases {
		tc, err := ParseTestCase(c[1])
		if err != nil {
			return testcase.TestSuite{}, err
		}

		suite.Cases = append(suite.Cases, tc)
	}

	return suite, nil
}
