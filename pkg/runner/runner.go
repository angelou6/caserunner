package runner

import (
	"bufio"
	"bytes"
	"caserunner/pkg/colors"
	"caserunner/pkg/testcase"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"time"
)

func process(program string, inputs []string, timeLimit time.Duration) ([]string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	args := strings.Split(program, " ")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, "", err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, "", err
	}

	var stdinErr error
	writeDone := make(chan struct{})
	go func() {
		defer close(writeDone)
		defer stdin.Close()
		for _, input := range inputs {
			if _, err := fmt.Fprintln(stdin, input); err != nil {
				stdinErr = err
				return
			}
		}
	}()

	var output []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		output = append(output, strings.TrimSpace(scanner.Text()))
	}

	<-writeDone
	waitErr := cmd.Wait()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return output, stderr.String(), errors.New("Programa excedió el tiempo limite")
	}
	if waitErr != nil {
		return output, stderr.String(), errors.New("Programa teminó con un error")
	}
	if stdinErr != nil {
		return output, stderr.String(), errors.New("Stdin se cerro inesperadamente")
	}
	return output, stderr.String(), nil
}

func RunSuite(suite testcase.TestSuite, verbose, halt bool) Result {
	var res Result
	for i, c := range suite.Cases {
		r, stderr, err := process(suite.Exec, c.Input, suite.TimeLimit)
		if err != nil {
			res.failed++
			colors.Printf(colors.Red, "Caso %d falló: %s\n", i+1, err)
			if stderr != "" {
				fmt.Print(stderr)
				if !strings.HasSuffix(stderr, "\n") {
					fmt.Println()
				}
			}
			if halt {
				break
			}
			continue
		}

		if slices.Equal(r, c.Output) {
			res.success++
			if verbose {
				colors.Printf(colors.Green, "Caso %d correcto:\n", i+1)
				fmt.Printf("Se esperaba %q, se obtuvo %q\n", r, c.Output)
			}
		} else {
			res.incorrect++
			colors.Printf(colors.Yellow, "Caso %d incorrecto:\n", i+1)
			fmt.Printf("Se esperaba %q, se obtuvo %q\n", r, c.Output)
			if halt {
				break
			}
		}
	}
	return res
}
