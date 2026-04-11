package runner

import (
	"bufio"
	"bytes"
	"caserunner/internal/colors"
	"caserunner/internal/parser"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func runTest(t parser.TestCase, command string, timeLimit time.Duration) ([]string, error) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	lines := make(chan string)
	done := make(chan error, 1)

	// Read the stderr
	var stderrBuf bytes.Buffer
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrBuf.WriteString(scanner.Text() + "\n")
		}
	}()

	// Read the stdout
	// If the program ends with error we send what is recorded in stderr
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
					return []string{}, errors.New(stderrLine)
				}
				return []string{}, errors.New("stdout se cerro inesperadamente.")
			}
			elapsed := time.Since(start)
			if timeLimit == -1 || elapsed <= timeLimit {
				output = append(output, strings.TrimRight(response, " "))
			} else {
				return []string{}, errors.New("Tiempo excedido.")
			}
		case err := <-done:
			return []string{}, err
		}
	}

	stdin.Close()
	<-done
	return output, nil
}

func RunFile(testcases *parser.TestFile, verbose bool, halt bool) {
	var correct, incorrect, failure int

	for i, test := range testcases.Tests {
		res, err := runTest(test, testcases.Exec, testcases.TimeLimit)
		if err != nil {
			failure++

			colors.Printf(colors.Red, "Error en problema %d:\n", i+1)
			fmt.Printf("%v\n\n", err)

			if halt {
				break
			}
			continue
		}

		result := test.JudgeOutput(res)
		if !result {
			incorrect++

			colors.Printf(colors.Yellow, "Problema %d incorrecto\n", i+1)
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
				fmt.Printf("Input: %q, output: %q\n\n", test, res)
			}
		}
	}

	correctaStr := "correctas"
	if correct == 1 {
		correctaStr = "correcta"
	}
	incorrectaStr := "incorrectas"
	if incorrect == 1 {
		incorrectaStr = "incorrecta"
	}
	falloStr := "fallos"
	if failure == 1 {
		falloStr = "fallo"
	}
	fmt.Printf("%d %s, %d %s, %d %s\n",
		correct,
		colors.Colorize(correctaStr, colors.Green),
		incorrect,
		colors.Colorize(incorrectaStr, colors.Yellow),
		failure,
		colors.Colorize(falloStr, colors.Red),
	)
}
