package wordcloud_go

import (
	"strconv"
	"fmt"
)

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
