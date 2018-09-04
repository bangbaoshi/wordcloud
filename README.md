# golang版本的文字云算法实现

#### 效果图

<img width="50%" src="/imgs/template.png"/>
<img width="50%" src="/imgs/template2.png"/>

#### 测试步骤如下

````
git clone https://github.com/bangbaoshi/wordcloud.git

cd wordcloud

go run boot/main.go

````
通过以上三步即可在imgs目录中生成文字云图片(查看imgs/out.png)

#### 目录介绍

1. boot目录包含测试用例
2. fonts目录包含若干种字体(非商业使用)
3. imgs目录包含模板图片，文字云生成的效果图就是按照模板图片的样子来生成

#### 使用说明

boot/main.go中已经简单介绍了使用方法
```
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
		&color.RGBA{0x0, 0x60, 0x30, 0xff},
		&color.RGBA{0x60, 0x0, 0x0, 0xff},
		&color.RGBA{0x73, 0x73, 0x0, 0xff},
	}
	//设置对应的字体路径，和输出路径
	render := wordcloud_go.NewWordCloudRender(60, 8,
		"./fonts/xin_shi_gu_yin.ttf",
		"./imgs/tiger.png", textList, angles, colors, "./imgs/out.png")
	//开始渲染
	render.Render()
}

func main() {
	renderNow()
}


```

#### 项目介绍
1. 使用golang语言实现了文字云算法
2. 用golang实现一些有趣的想法
