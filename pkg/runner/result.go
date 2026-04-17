package runner

import (
	"caserunner/internal/colors"
	"fmt"
)

type Result struct {
	success, incorrect, failed uint16
}

func (r Result) PrintResults() {
	fmt.Printf(
		"%s: %d, %s: %d, %s: %d\n",
		colors.Colorize(colors.Green, "Correctas"), r.success,
		colors.Colorize(colors.Yellow, "Incorrectas"), r.incorrect,
		colors.Colorize(colors.Red, "Fallos"), r.failed,
	)
}
