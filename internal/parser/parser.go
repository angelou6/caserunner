package parser

import (
	"caserunner/internal/colors"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func matchFirstRegex(input string, regex string) (string, error) {
	timeR := regexp.MustCompile(regex)
	matches := timeR.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", fmt.Errorf("Patron %q no encontrado en input", regex)
	}

	return strings.TrimSpace(matches[0][1]), nil
}

func getIndexes(to direction, arr []string) (int, error) {
	foundIndexes := []int{}
	for i, s := range arr {
		if s == string(to) {
			foundIndexes = append(foundIndexes, i)
		}
	}

	if len(foundIndexes) > 1 {
		return 0, fmt.Errorf("Multiples '%s' encontrados", to)
	}
	if len(foundIndexes) == 0 {
		return 0, fmt.Errorf("'%s' no encontrado", to)
	}
	return foundIndexes[0], nil
}

func (t *TestFile) parseTest(input string) error {
	if len(input) == 0 {
		return errors.New("Caso de prueba vacío.")
	}

	caseInputs := strings.Split(input, "\n")

	inputIdx, err := getIndexes(Input, caseInputs)
	if err != nil {
		return err
	}

	outputIdx, err := getIndexes(Output, caseInputs)
	if err != nil {
		return err
	}

	tc := TestCase{}
	tc.AppendToCase(inputIdx, caseInputs, Input)
	tc.AppendToCase(outputIdx, caseInputs, Output)

	t.Tests = append(t.Tests, tc)
	return nil
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
		err := t.parseTest(body)
		if err != nil {
			return fmt.Errorf(
				"%s: \n%v",
				colors.Colorize(fmt.Sprintf("Error parseando la prueba %d:", i+1), colors.Red),
				err,
			)
		}
	}

	return nil
}
