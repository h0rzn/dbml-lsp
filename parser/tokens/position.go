package tokens

import "fmt"

type Position struct {
	Line   uint32
	Offset uint32
	Len    uint32
}

func (p *Position) String() string {
	return fmt.Sprintf("[%d:%d-%d]", p.Line, p.Offset, p.Offset+p.Len)
}
