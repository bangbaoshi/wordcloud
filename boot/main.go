package main

import (
	"image/color"
	"github.com/bangbaoshi/wordcloud"
	"time"
	"fmt"
)

func renderNow() {
	textList := []string{"恭喜", "发财", "万事", "如意"}
	angles := []int{0, 15, -15, 90}
	colors := []*color.RGBA{
		&color.RGBA{0x0, 0x60, 0x30, 0xff},
		&color.RGBA{0x60, 0x0, 0x0, 0xff},
		// &color.RGBA{0x73, 0x73, 0x0, 0xff},
	}
	render := wordcloud_go.NewWordCloudRender(60, 8,
		"./fonts/xin_shi_gu_yin.ttf",
		"./imgs/foot.png", textList, angles, colors, "./imgs/foot_template.png")
	render.Render()
}

func main() {
	startedAt := time.Now().Unix()
	renderNow()
	endAt := time.Now().Unix()
	fmt.Printf("时间消耗:%d\n", endAt-startedAt);
}
