package colors

import "fmt"

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

func RGB(red int, green int, blue int) Color {
	return Color(fmt.Sprintf("38;2;%d;%d;%d", red, green, blue))
}

func Print(input string, color Color) {
	fmt.Printf("\x1b[%sm%s\x1b[0m", color, input)
}

func Println(input string, color Color) {
	Print(input+"\n", color)
}
