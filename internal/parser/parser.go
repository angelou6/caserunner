package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type direction string

const (
	Input  direction = "input"
	Output direction = "output"
)

type TestCase struct {
	input  []string
	output []string
}

type TestFile struct {
	Exec      string
	TimeLimit time.Duration
	Tests     []TestCase
}

func New() *TestFile {
	return &TestFile{}
}

func matchFirstRegex(input string, regex string) (string, error) {
	timeR := regexp.MustCompile(regex)
	matches := timeR.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", fmt.Errorf("no match found for pattern %q in input", regex)
	}

	return strings.TrimSpace(matches[0][1]), nil
}

func getIndexes(str string, arr []string) []int {
	foundIndexes := []int{}
	for i, s := range arr {
		if str == s {
			foundIndexes = append(foundIndexes, i)
		}
	}
	return foundIndexes
}

func appendToCase(indexes []int, tcinput []string, into direction) []string {
	res := []string{}

	otherInto := ""
	if into == "input" {
		otherInto = "output:"
	} else {
		otherInto = "input:"
	}

	for i := indexes[0]; i < len(tcinput); i++ {
		if len(tcinput[i]) == 0 || tcinput[i] == string(into)+":" {
			continue
		} else if tcinput[i] == otherInto {
			break
		}

		s := strings.TrimSpace(tcinput[i])
		s = strings.Replace(s, "\\", "", 1)
		res = append(res, s)
	}
	return res
}

// This function fucking sucks
// Everything arround this function fucking sucks
func parseTest(input string) (TestCase, error) {
	if len(input) == 0 {
		// TODO: better error message with line numbers and shit
		return TestCase{}, errors.New("Test case is empty.")
	}

	tcinput := strings.Split(input, "\n")
	inputIndexes := getIndexes("input:", tcinput)
	outputIndexes := getIndexes("output:", tcinput)

	// TODO: better error messages with line numbers and shit
	if len(inputIndexes) > 1 || len(outputIndexes) > 1 {
		return TestCase{}, errors.New("Multiple indexes or outputs found")
	} else if len(inputIndexes) == 0 || len(outputIndexes) == 0 {
		return TestCase{}, errors.New("Input or output not found")
	}

	tc := TestCase{}

	tc.input = appendToCase(inputIndexes, tcinput, "input")
	tc.output = appendToCase(outputIndexes, tcinput, "output")

	return tc, nil
}

func (t *TestFile) ParseFile(program string, input string) error {
	// Get exec
	exline, err := matchFirstRegex(input, "exec:(.*)")
	if err != nil {
		return errors.New("Exec field not found.")
	}
	t.Exec = strings.ReplaceAll(exline, "$code", program)

	// Get time limit
	timeline, err := matchFirstRegex(input, "time-limit:(.*)")
	t.TimeLimit, err = time.ParseDuration(timeline)
	if err != nil {
		t.TimeLimit = -1
	}

	// Parse test cases
	re := regexp.MustCompile(`(?sm)^--\n(.*?)^--`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		body := strings.TrimSpace(match[1])
		tc, err := parseTest(body)
		if err != nil {
			return err
		}
		t.Tests = append(t.Tests, tc)
	}

	return nil
}
