package wordcloud_go

const (
	IS_NOT_FIT = 1
	IS_FIT     = 2
	OUT_INDEX  = 3
	DEGREE_360 = 360
	DEGREE_180 = 180
	IS_EMPTY   = 0
	XUNIT      = 2
	YUNIT      = 2
)

type Position struct {
	Xpos   int
	Ypos   int
	Value  int
	XLeiji int
	YLeiji int
}

func NewPosition(xpos, ypos, value, xleiji, yleiji int) *Position {
	pos := &Position{
		Xpos:   xpos,
		Ypos:   ypos,
		Value:  value,
		XLeiji: xleiji,
		YLeiji: yleiji,
	}
	return pos
}

type Grid struct {
	Width     int
	Height    int
	positions []*Position
	XScale    int
	YScale    int
}

func (this *Grid) IsFit(xIncrement, yIncrement, width, height int, gridIntArray []int) int {
	for i := 0; i < this.Height; i++ {
		for j := 0; j < this.Width; j++ {
			index := i*this.Width + j
			position := this.positions[index]
			if position.Value != IS_EMPTY {
				position.Xpos = position.XLeiji + xIncrement
				position.Ypos = position.YLeiji + yIncrement
				if position.Xpos < 0 || position.Ypos < 0 || position.Xpos >= width || position.Ypos >= height {
					return OUT_INDEX
				}
				index = position.Ypos*width + position.Xpos
				if position.Value != 0 && gridIntArray[index] == position.Value {
					return IS_NOT_FIT
				}
			}
		}
	}
	return IS_FIT
}

func (this *Grid) SetCollisionMap(collisionMap []int, width, height int) {
	this.Width = width
	this.Height = height
	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := collisionMap[index]
			position := NewPosition(x, y, value, 0, 0)
			this.positions = append(this.positions, position)
			index++
		}
	}
}

func (this *Grid) Fill(gridIntArrayWidth, gridIntArrayHeight int, gridIntArray []int) {
	for y := 0; y < this.Height; y++ {
		for x := 0; x < this.Width; x++ {
			index := y*this.Width + x
			position := this.positions[index]
			index = position.Ypos*gridIntArrayWidth + position.Xpos
			if position.Value != IS_EMPTY {
				gridIntArray[index] = position.Value
			}
		}
	}
}
