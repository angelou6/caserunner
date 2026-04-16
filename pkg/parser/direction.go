package parser

type direction string

const (
	input  direction = "input:"
	output direction = "output:"
)

func (d direction) oposite() direction {
	if d == input {
		return output
	}
	return input
}
