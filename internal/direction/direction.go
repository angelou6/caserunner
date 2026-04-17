package direction

type Direction string

const (
	Input  Direction = "input:"
	Output Direction = "output:"
)

func (d Direction) Oposite() Direction {
	if d == Input {
		return Output
	}
	return Input
}
