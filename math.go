package wordcloud

import (
	"fmt"
	"github.com/bangbaoshi/gg"
	"image/png"
	"math"
	"os"
	"strconv"
)


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

type WorldMap struct {
	Width           int
	Height          int
	CollisionMap    []int
	RealImageWidth  int
	RealImageHeight int
}

func (this *WorldMap) PrintMap() {
	for y := 0; y < this.Height; y++ {
		str := ""
		for x := 0; x < this.Width; x++ {
			idx := y*this.Width + x
			str = str + strconv.Itoa(this.CollisionMap[idx])
		}
		fmt.Println(str)
	}
}

func TwoByBitmap(imgpath string) *WorldMap {
	worldMap := &WorldMap{
		CollisionMap: make([]int, 0),
	}
	file, err := os.Open(imgpath)
	if err != nil {
		fmt.Println(err)
	}
	img, err := png.Decode(file)
	file.Close()
	bounds := img.Bounds()
	w := bounds.Size().X
	h := bounds.Size().Y
	worldMap.Width = w / XUNIT
	worldMap.Height = h / YUNIT
	worldMap.RealImageWidth = w
	worldMap.RealImageHeight = h
	for y := 0; y < worldMap.Height; y++ {
		for x := 0; x < worldMap.Width; x++ {
			color := img.At(x*XUNIT, y*YUNIT)
			_, _, _, alpha := color.RGBA()
			if alpha == 0 {
				worldMap.CollisionMap = append(worldMap.CollisionMap, 1)
			} else {
				worldMap.CollisionMap = append(worldMap.CollisionMap, 0)
			}
		}
	}
	return worldMap
}

func TwoByBlock(width, height int) ([]*Position, int, int) {
	maxX := width / XUNIT
	maxY := height / YUNIT
	len := maxX * maxY
	positions := make([]*Position, len)
	for i := 0; i < len; i++ {
		positions[i] = NewPosition(0, 0, IS_NOT_FIT, 0, 0)
	}
	return positions, maxX, maxY
}

func DrawText(dc *gg.Context, text string, xpos, ypos, rotation float64) {
	if rotation != 0 {
		dc.RotateAbout(rotation, xpos, ypos)
	}
	dc.DrawStringAnchored(text, xpos, ypos, 0.5, 0.5)
	if rotation != 0 {
		dc.RotateAbout(-rotation, xpos, ypos)
	}
}

func GetTextBound(measureDc *gg.Context, text string) (w, h, xdiff, ydiff float64) {
	measureDc.SetRGBA(0, 0, 0, 0)
	measureDc.Clear()
	measureDc.SetRGBA(0, 0, 0, 1)
	measureDc.DrawStringAnchored(text, 375, 375, 0.5, 0.5)
	img := measureDc.Image()
	width := measureDc.Width()
	height := measureDc.Height()
	maxX := 0
	maxY := 0
	minX := 9999999
	minY := 9999999
	for y := 0; y < height; y++ {

		for x := 0; x < width; x++ {
			color := img.At(x, y)
			_, _, _, alpha := color.RGBA()

			if alpha != 0 {
				if minX > x {
					minX = x
				}
				if minY > y {
					minY = y
				}
				if maxX < x {
					maxX = x
				}
				if maxY < y {
					maxY = y
				}
			}
		}
	}
	w1, h1 := measureDc.MeasureString(text)
	wdiff := float64(maxX - minX)
	hdiff := float64(maxY - minY)
	xdiff = float64(w1 - wdiff)
	ydiff = float64(h1 - hdiff)

	return wdiff, hdiff, xdiff, ydiff
	//return w1, h1, xdiff, ydiff
}

/*
 *先设置清空颜色，再进行清空
 */
func Clear(dc *gg.Context) {
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()
}

func Rotate(grid *Grid, angle float64, centerX, centerY int) {
	maxX := grid.Width
	maxY := grid.Height
	width := maxX * XUNIT
	height := maxY * YUNIT
	halfX := width / 2
	halfY := height / 2
	tempX := 0
	tempY := 0
	gridData := grid.positions
	sinPi := SinT(angle)
	cosPi := CosT(angle)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			index := y*maxX + x
			position := gridData[index]
			position.Xpos = x
			position.Ypos = y
			position.Xpos = position.Xpos*XUNIT - halfX
			position.Ypos = position.Ypos*YUNIT - halfY
			tempX = position.Xpos
			tempY = position.Ypos
			position.Xpos = (int)(float64(tempX)*cosPi - float64(tempY)*sinPi)
			position.Ypos = (int)(float64(tempX)*sinPi + float64(tempY)*cosPi)
			position.Xpos /= XUNIT
			position.Ypos /= YUNIT

			position.Xpos += centerX
			position.Ypos += centerY

			position.XLeiji = position.Xpos
			position.YLeiji = position.Ypos
		}
	}
}

func CeilT(value float64) float64 {
	return math.Ceil(value)
}

func CosT(angle float64) float64 {
	angle = angle / DEGREE_180 * math.Pi
	return math.Cos(angle)
}

func SinT(angle float64) float64 {
	angle = angle / DEGREE_180 * math.Pi
	return math.Sin(angle)
}

func Angle2Pi(angle float64) float64 {
	return angle / DEGREE_180 * math.Pi
}
