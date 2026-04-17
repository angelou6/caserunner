package colors

import (
	"fmt"
	"os"
)

type Color string

const (
	Black   Color = "30"
	Red     Color = "31"
	Green   Color = "32"
	Yellow  Color = "33"
	Blue    Color = "34"
	Magenta Color = "35"
	Cyan    Color = "36"
	White   Color = "37"
)

func RGB(red, green, blue int) Color {
	return Color(fmt.Sprintf("38;2;%d;%d;%d", red, green, blue))
}

func Print(color Color, input string) {
	fmt.Printf("\x1b[%sm%s\x1b[0m", color, input)
}

func Println(color Color, input string) {
	Print(color, input+"\n")
}

func Colorize(color Color, input string) string {
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, input)
}

func Printf(color Color, format string, a ...any) (n int, err error) {
	colorized := Colorize(color, format)
	return fmt.Fprintf(os.Stdout, colorized, a...)
}
