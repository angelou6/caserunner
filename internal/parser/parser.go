package parser

import (
	"caserunner/internal/colors"
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

func matchFirstRegex(input string, regex string) (string, error) {
	timeR := regexp.MustCompile(regex)
	matches := timeR.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", fmt.Errorf("Patron %q no encontrado en input", regex)
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

func parseTest(input string) (TestCase, error) {
	if len(input) == 0 {
		return TestCase{}, errors.New("Caso de prueba vacío.")
	}

	tcinput := strings.Split(input, "\n")
	inputIndexes := getIndexes("input:", tcinput)
	outputIndexes := getIndexes("output:", tcinput)

	if len(inputIndexes) > 1 || len(outputIndexes) > 1 {
		return TestCase{}, errors.New("Multiples indexes o outputs encontrados.")
	} else if len(inputIndexes) == 0 || len(outputIndexes) == 0 {
		return TestCase{}, errors.New("Input o output no encontrados.")
	}

	tc := TestCase{}

	tc.Input = appendToCase(inputIndexes, tcinput, "input")
	tc.Output = appendToCase(outputIndexes, tcinput, "output")

	return tc, nil
}

func (t *TestFile) ParseFile(input string, program string) error {
	// Get exec
	exline, err := matchFirstRegex(input, "exec:(.*)")
	if err != nil {
		return errors.New("Exec no encotrado.")
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

	for i, match := range matches {
		body := strings.TrimSpace(match[1])
		tc, err := parseTest(body)
		if err != nil {
			colors.Println(fmt.Sprintf("Error parseando la prueba %d:", i+1), colors.Red)
			return err
		}
		t.Tests = append(t.Tests, tc)
	}

	return nil
}
