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

func Print(input string, color Color) {
	fmt.Printf("\x1b[%sm%s\x1b[0m", color, input)
}

func Println(input string, color Color) {
	Print(input+"\n", color)
}

func Colorize(input string, color Color) string {
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, input)
}

func Printf(color Color, format string, a ...any) (n int, err error) {
	colorized := Colorize(format, color)
	return fmt.Fprintf(os.Stdout, colorized, a...)
}
