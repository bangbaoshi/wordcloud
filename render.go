package wordcloud_go

import (
	"github.com/fogleman/gg"
	"image/color"
	"strconv"
)

type WordCloudRender struct {
	MaxFontSize    float64
	MinFontSize    float64
	FontPath       string
	OutlineImgPath string
	MeasureDc      *gg.Context
	DrawDc         *gg.Context
	TextList       []string
	Angles         []int
	Colors         []*color.RGBA
	OutImgPath     string
	worldMap       *WorldMap
}

func NewWordCloudRender(maxFontSize, minFontSize float64, fontPath string,
	imgPath string, textList []string,
	angles []int, colors []*color.RGBA,
	outImgPath string) *WordCloudRender {

	render := &WordCloudRender{
		MaxFontSize:    maxFontSize,
		MinFontSize:    minFontSize,
		FontPath:       fontPath,
		OutlineImgPath: imgPath,
		TextList:       textList,
		Angles:         angles,
		Colors:         colors,
		OutImgPath:     outImgPath,
	}
	worldMap := TwoByBitmap(imgPath)
	render.worldMap = worldMap
	drawDc := gg.NewContext(worldMap.RealImageWidth, worldMap.RealImageHeight)
	drawDc.SetRGB(1, 1, 1)
	drawDc.Clear()
	drawDc.SetRGB(0, 0, 0)
	render.DrawDc = drawDc
	if err := drawDc.LoadFontFace(fontPath, render.MaxFontSize); err != nil {
		panic(err)
	}

	render.ResetMeasureDc(render.MaxFontSize)
	return render
}

func (this *WordCloudRender) Render() {
	fontSize := this.MaxFontSize
	currentTextIdx := 0
	colorIdx := 0
	checkRet := &CheckResult{}

	gridCache := make(map[string]*Grid)

	var itemGrid *Grid
	bigestSizeCnt := 0
	for {
		var msg string = this.TextList[currentTextIdx]
		key := strconv.Itoa(int(fontSize)) + msg
		if _, ok := gridCache[key]; ok {
			//取出缓存中的itemGrid
			itemGrid = gridCache[key]
		} else {
			itemGrid = &Grid{}
			//性能主要消耗点
			w, h, xscale, yscale := GetTextBound(this.MeasureDc, msg)
			itemGrid.XScale = int(xscale)
			itemGrid.YScale = int(yscale)
			if int(w)%2 != 0 {
				w += XUNIT
			}
			if int(h)%2 != 0 {
				h += YUNIT
			}
			positions, w1, h1 := TwoByBlock(int(w), int(h))
			itemGrid.Width = int(w1)
			itemGrid.Height = int(h1)
			itemGrid.positions = positions
			//把(fontSize+msg)当做一个组合缓存下来
			gridCache[key] = itemGrid
		}

		isFound := this.collisionCheck(
			0, this.worldMap, itemGrid, checkRet, this.Angles)
		if isFound {

			currentTextIdx++
			currentTextIdx = currentTextIdx % len(this.TextList)
			color := this.Colors[colorIdx]
			colorIdx++
			colorIdx = colorIdx % len(this.Colors)
			this.DrawDc.SetRGB(float64(color.R), float64(color.G), float64(color.B))

			DrawText(this.DrawDc, msg, float64(checkRet.Xpos),
				float64(checkRet.Ypos), Angle2Pi(float64(checkRet.Angle)))
			if fontSize == this.MaxFontSize {
				bigestSizeCnt++
				if bigestSizeCnt > len(this.TextList) {
					fontSize = 40
					this.UpdateFontSize(fontSize)
				}
			}
		} else {
			fontSize -= 3
			if fontSize < this.MinFontSize {
				break
			}
			this.UpdateFontSize(fontSize)
		}
	}
	this.DrawDc.SavePNG(this.OutImgPath)
}

func (this *WordCloudRender) UpdateFontSize(fontSize float64) {
	this.DrawDc.SetFontSize(fontSize)
	this.MeasureDc.SetFontSize(fontSize)
}

func (this *WordCloudRender) ResetMeasureDc(fontSize float64) {
	measureDc := gg.NewContext(this.worldMap.RealImageWidth, this.worldMap.RealImageHeight)
	measureDc.SetRGBA(0, 0, 0, 0)
	measureDc.Clear()
	this.MeasureDc = measureDc
	if err := measureDc.LoadFontFace(this.FontPath, fontSize); err != nil {
		panic(err)
	}
}

func (this *WordCloudRender) collisionCheck(lastCheckAngle float64, worldMap *WorldMap,
	itemGrid *Grid, ret *CheckResult, tryAngles []int) bool {

	centerX := worldMap.Width / 2
	centerY := worldMap.Height / 2
	isFound := true
	xDistanceToCenter := 0
	yDistanceToCenter := 0
	tempXpos := 0
	tempYpos := 0

	angleMark := 0
	currentAngleIdx := 0
	for angle := lastCheckAngle; angle <= DEGREE_360; angle += 1 {
		currentAngleIdx = 0
		angleMark = tryAngles[currentAngleIdx]
		currentAngleIdx++
		Rotate(itemGrid, float64(angleMark), centerX, centerY)
		xDiff := CosT(angle) * 1
		yDiff := SinT(angle) * 1
		tempXpos = 0
		tempYpos = 0
		xLeiji := xDiff
		yLeiji := yDiff
		xDistanceToCenter = 0
		yDistanceToCenter = 0
		result := IS_NOT_FIT
		for {
			result = IS_NOT_FIT
			if xDistanceToCenter != tempXpos || yDistanceToCenter != tempYpos {
				tempXpos = xDistanceToCenter
				tempYpos = yDistanceToCenter
				result = itemGrid.IsFit(xDistanceToCenter, yDistanceToCenter, worldMap.Width, worldMap.Height, worldMap.CollisionMap)
				if result == OUT_INDEX {
					if currentAngleIdx < len(tryAngles) {
						angleMark = tryAngles[currentAngleIdx]
						currentAngleIdx++
						Rotate(itemGrid, float64(angleMark), centerX, centerY)
						xLeiji = xDiff
						yLeiji = yDiff
						tempXpos = 0
						tempYpos = 0
						xDistanceToCenter = 0
						yDistanceToCenter = 0
					} else {
						ret.Angle = 0
						isFound = false
						break
					}
				} else if result == IS_FIT {
					isFound = true
					itemGrid.Fill(worldMap.Width, worldMap.Height, worldMap.CollisionMap)
					ret.Angle = angleMark
					ret.Xpos = (xDistanceToCenter + centerX) * XUNIT
					ret.Ypos = (yDistanceToCenter + centerY) * YUNIT
					ret.LastCheckAngle = int(angle)
					break
				}
			}
			xLeiji += xDiff
			yLeiji += yDiff
			xDistanceToCenter = int(CeilT(xLeiji))
			yDistanceToCenter = int(CeilT(yLeiji))
		}
		if angle >= DEGREE_360 {
			ret.Angle = 0
			isFound = false
			break
		}
		if result == IS_FIT {
			break
		}
	}
	return isFound
}
