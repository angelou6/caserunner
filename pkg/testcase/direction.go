package testcase

type Direction string

const (
	Input  Direction = "input:"
	Output Direction = "output:"
)

func (d Direction) oposite() Direction {
	if d == Input {
		return Output
	}
	return Input
}
