package runner

import (
	"bufio"
	"bytes"
	"caserunner/pkg/colors"
	"caserunner/pkg/parser"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func runTest(t parser.TestCase, command string, timeLimit time.Duration) ([]string, time.Duration, error) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	testStart := time.Now()
	cmd.Start()

	lines := make(chan string)
	done := make(chan error, 1)

	// Leer el stderr
	var stderrBuf bytes.Buffer
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrBuf.WriteString(scanner.Text() + "\n")
		}
	}()

	// Leer el stdout
	// Si el programa termina con error, enviamos lo que se registró en stderr
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
		err := cmd.Wait()
		if err != nil {
			done <- fmt.Errorf("%w:\nstderr: %s", err, stderrBuf.String())
		} else {
			done <- nil
		}
	}()

	output := []string{}
	for i, input := range t.Input {
		start := time.Now()
		fmt.Fprintln(stdin, input)

		if i < len(t.Output) && t.Output[i] == "" {
			continue
		}

		select {
		case response, ok := <-lines:
			if !ok {
				stderrLine := stderrBuf.String()

				if len(stderrLine) > 0 {
					return []string{}, 0, errors.New(stderrLine)
				}
				return []string{}, 0, errors.New("stdout se cerro inesperadamente.")
			}
			elapsed := time.Since(start)
			if timeLimit == -1 || elapsed <= timeLimit {
				output = append(output, strings.TrimRight(response, " "))
			} else {
				return []string{}, 0, errors.New("Tiempo excedido.")
			}
		case err := <-done:
			return []string{}, 0, err
		}
	}

	stdin.Close()
	<-done
	return output, time.Since(testStart), nil
}

func printResult(caseCount, correct, incorrect, failure int, totalCaseTime time.Duration) {
	plural := func(n int, singular string) string {
		if n == 1 {
			return singular
		}
		return singular + "s"
	}

	if caseCount > 0 {
		avg := totalCaseTime / time.Duration(caseCount)
		fmt.Printf("Tiempo total: %s, promedio por caso: %s\n", totalCaseTime, avg)
	} else {
		fmt.Printf("Tiempo total: %s\n", totalCaseTime)
	}

	fmt.Printf("%d %s, %d %s, %d %s\n",
		correct, colors.Colorize(plural(correct, "correcta"), colors.Green),
		incorrect, colors.Colorize(plural(incorrect, "incorrecta"), colors.Yellow),
		failure, colors.Colorize(plural(failure, "fallo"), colors.Red),
	)
}

func RunFile(testcases *parser.TestFile, verbose, halt bool) {
	var correct, incorrect, failure int
	var totalCaseTime time.Duration
	var caseCount int

	for i, test := range testcases.Tests {
		res, elapsed, err := runTest(test, testcases.Exec, testcases.TimeLimit)
		if err != nil {
			failure++

			colors.Printf(colors.Red, "Error en problema %d:\n", i+1)
			fmt.Printf("%v\n\n", err)

			if halt {
				break
			}
			continue
		}

		totalCaseTime += elapsed
		caseCount++

		result := test.JudgeOutput(res)
		if !result {
			incorrect++

			colors.Printf(colors.Yellow, "Problema %d incorrecto\n", i+1)
			if verbose {
				fmt.Printf("Tiempo: %s\n", elapsed)
			}
			expected := []string{}
			for _, o := range test.Output {
				if o != "" {
					expected = append(expected, o)
				}
			}
			fmt.Printf("Se esperaba %q, se obtuvo %q\n\n", expected, res)
			if halt {
				break
			}
		} else {
			correct++
			if verbose {
				colors.Printf(colors.Green, "Problema %d correcto\n", i+1)
				fmt.Printf("Tiempo: %s, input: %q, output: %q\n\n", elapsed, test, res)
			}
		}
	}

	printResult(caseCount, correct, incorrect, failure, totalCaseTime)
}
