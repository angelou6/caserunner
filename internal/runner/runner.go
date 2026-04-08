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
	for _, input := range t.Input {
		start := time.Now()
		fmt.Fprintln(stdin, input)

		select {
		case response, ok := <-lines:
			if !ok {
				return []string{}, fmt.Errorf("stdout closed unexpectedly\nstderr: %s", stderrBuf.String())
			}
			elapsed := time.Since(start)
			if timeLimit == -1 || elapsed <= timeLimit {
				output = append(output, response)
			} else {
				return []string{}, errors.New("timeout exceeded")
			}
		case err := <-done:
			return []string{}, err
		}
	}

	stdin.Close()
	<-done
	return output, nil
}

func RunFile(testcases parser.TestFile) {
	for i, test := range testcases.Tests {
		res, err := runTest(test, testcases.Exec, testcases.TimeLimit)
		if err != nil {
			// TODO: halt flag to stop execution on error
			// alse verbose flag for verbose output
			colors.Println(fmt.Sprintf("Test %d:", i), colors.Blue)
			colors.Print("Error: ", colors.Red)
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
}
