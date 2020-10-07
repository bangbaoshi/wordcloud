package wordcloud

import (
	"fmt"
	"testing"
	"time"
)

func renderNow() {
	textList := []string{"恭喜", "发财", "万事", "如意"}
	angles := []int{0}
	colors := []*Color{
		&Color{0x0, 0x60, 0x30},
		&Color{0x60, 0x0, 0x0},
		// &color.RGBA{0x73, 0x73, 0x0, 0xff},
	}
	render := NewWordCloudRender(60, 8,
		"./fonts/xin_shi_gu_yin.ttf",
		"./imgs/foot.png", textList, angles, colors, "./imgs/foot_template2.png")
	render.Render()
}

func TestNewWordCloudRender(t *testing.T) {
	startedAt := time.Now().UnixNano()
	renderNow()
	endAt := time.Now().UnixNano()
	fmt.Printf("时间消耗:%d\n", endAt-startedAt)
}
