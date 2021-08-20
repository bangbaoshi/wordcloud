package main

import (
	"image/color"

	"github.com/bangbaoshi/wordcloud"
)

func renderNow() {
	//需要写入的文本数组
	textList := []string{"恭喜", "发财", "万事", "如意"}
	//文本角度数组
	angles := []int{0, 15, -15, 90}
	//文本颜色数组
	colors := []*color.RGBA{
		{0x0, 0x60, 0x30, 0xff},
		{0x60, 0x0, 0x0, 0xff},
		{0x73, 0x73, 0x0, 0xff},
	}
	//设置对应的字体路径，和输出路径
	render := wordcloud.NewWordCloudRender(60, 8,
		"./fonts/xin_shi_gu_yin.ttf",
		"./imgs/china.png", textList, angles, colors, "./imgs/out.png")
	//开始渲染
	render.Render()
}

func main() {
	renderNow()
}
