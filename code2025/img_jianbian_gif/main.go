package main

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

// 创建高质量256色调色板（增加过渡色细腻度）
func createGlobalPalette() color.Palette {
	palette := make(color.Palette, 0, 256)
	// 1. 增加基础色数量（16种核心色，覆盖更全）
	baseColors := []color.RGBA{
		{255, 0, 0, 255},     // 红
		{255, 127, 0, 255},   // 橙红
		{255, 255, 0, 255},   // 黄
		{127, 255, 0, 255},   // 黄绿
		{0, 255, 0, 255},     // 绿
		{0, 255, 127, 255},   // 青绿
		{0, 255, 255, 255},   // 青
		{0, 127, 255, 255},   // 青蓝
		{0, 0, 255, 255},     // 蓝
		{127, 0, 255, 255},   // 蓝紫
		{255, 0, 255, 255},   // 紫
		{255, 0, 127, 255},   // 紫红
		{255, 127, 127, 255}, // 浅红
		{127, 255, 127, 255}, // 浅绿
		{127, 127, 255, 255}, // 浅蓝
		{255, 255, 127, 255}, // 浅黄
	}
	for _, c := range baseColors {
		palette = append(palette, color.Color(c))
	}

	// 2. 生成高密度过渡色（相邻色间生成更多过渡，提升渐变细腻度）
	remaining := 256 - len(baseColors)
	step := remaining / (len(baseColors)) // 每个基础色周围分配过渡色
	for i := 0; i < len(baseColors); i++ {
		c1 := baseColors[i]
		// 与下一个颜色生成过渡（循环衔接）
		c2 := baseColors[(i+1)%len(baseColors)]
		for t := 0; t < step; t++ {
			factor := float64(t) / float64(step)
			r := uint8(float64(c1.R)*(1-factor) + float64(c2.R)*factor)
			g := uint8(float64(c1.G)*(1-factor) + float64(c2.G)*factor)
			b := uint8(float64(c1.B)*(1-factor) + float64(c2.B)*factor)
			palette = append(palette, color.Color(color.RGBA{r, g, b, 255}))
		}
	}
	// 确保严格256色
	for len(palette) < 256 {
		palette = append(palette, palette[len(palette)-1])
	}
	return palette
}

// 高精度颜色匹配（检查所有调色板颜色，提升匹配精度）
func closestColorIndex(palette color.Palette, c color.RGBA) int {
	minDist := 1 << 30
	idx := 0
	// 检查所有颜色（牺牲一点性能换质量）
	for i, pc := range palette {
		pr, pg, pb, _ := pc.RGBA()
		// 更精确的距离计算（加入gamma校正感知）
		diffR := int(c.R) - int(pr>>8)
		diffG := int(c.G) - int(pg>>8)
		diffB := int(c.B) - int(pb>>8)
		// 人眼对绿色更敏感，增加权重
		dist := diffR*diffR*2 + diffG*diffG*3 + diffB*diffB*2
		if dist < minDist {
			minDist = dist
			idx = i
			if minDist == 0 {
				break
			}
		}
	}
	return idx
}

// 扩散点结构（增加柔和度参数）
type Point struct {
	x, y       int
	color      color.RGBA
	startFrame int
	lifetime   int
	maxRadius  int
	softness   float64 // 边缘柔和度（0-1）
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 质量优先的参数调整
	width, height := 600, 600 // 适度增大尺寸（比500x500多44%像素，细节更丰富）
	frames := 150             // 增加帧数（提升流畅度）
	delay := 8                // 缩短延迟（保持动画节奏）
	totalPoints := 96         // 增加点数量（画面更饱满）
	batchSize := 24           // 每批点数量同步增加
	interval := 20            // 缩短间隔（点生成更密集）

	palette := createGlobalPalette()

	// 生成更高质量的扩散点（增加边缘柔和度）
	points := make([]Point, 0, totalPoints)
	for i := 0; i < totalPoints; i++ {
		startFrame := (i / batchSize) * interval
		points = append(points, Point{
			x:          rand.Intn(width),
			y:          rand.Intn(height),
			startFrame: startFrame,
			lifetime:   100,                      // 延长生命周期（扩散更充分）
			maxRadius:  90 + rand.Intn(180),      // 增大最大半径（覆盖范围更广）
			softness:   0.7 + rand.Float64()*0.3, // 随机柔和度（0.7-1.0）
			color: color.RGBA{
				// 从调色板中随机选择更丰富的颜色
				palette[rand.Intn(32)].(color.RGBA).R,
				palette[rand.Intn(32)].(color.RGBA).G,
				palette[rand.Intn(32)].(color.RGBA).B,
				255,
			},
		})
	}

	gifImg := &gif.GIF{
		Image:     make([]*image.Paletted, frames),
		Delay:     make([]int, frames),
		LoopCount: 0,
	}

	var wg sync.WaitGroup
	wg.Add(frames)

	for frame := 0; frame < frames; frame++ {
		go func(f int) {
			defer wg.Done()
			img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

			// 高质量背景（使用更细腻的深色）
			bgColor := color.RGBA{10, 10, 15, 255} // 深灰偏蓝，更显色彩
			bgIdx := closestColorIndex(palette, bgColor)
			for y := 0; y < height; y++ {
				row := img.Pix[y*img.Stride : (y+1)*img.Stride]
				for x := 0; x < width; x++ {
					row[x] = uint8(bgIdx)
				}
			}

			// 更高质量的扩散计算
			for _, p := range points {
				if f < p.startFrame || f > p.startFrame+p.lifetime {
					continue
				}

				elapsed := f - p.startFrame
				growthSpeed := float64(p.maxRadius) / float64(p.lifetime)
				currentRadius := float64(elapsed) * growthSpeed
				r := int(currentRadius)

				// 精确计算扩散范围
				minX := p.x - r
				if minX < 0 {
					minX = 0
				}
				maxX := p.x + r
				if maxX >= width {
					maxX = width - 1
				}
				minY := p.y - r
				if minY < 0 {
					minY = 0
				}
				maxY := p.y + r
				if maxY >= height {
					maxY = height - 1
				}

				// 细腻的颜色混合（保留更多淡色细节）
				for y := minY; y <= maxY; y++ {
					for x := minX; x <= maxX; x++ {
						dx := float64(x - p.x)
						dy := float64(y - p.y)
						distSq := dx*dx + dy*dy
						radiusSq := currentRadius * currentRadius
						if distSq > radiusSq {
							continue
						}

						// 更自然的衰减曲线（基于柔和度参数）
						dist := math.Sqrt(distSq)
						normDist := dist / currentRadius
						// 边缘柔和度调整：softness越高，边缘过渡越平缓
						weight := 1 - math.Pow(normDist, 2/p.softness)
						if weight < 0.05 { // 保留更多淡色（之前是0.1）
							continue
						}

						// 高精度颜色混合
						currentIdx := img.ColorIndexAt(x, y)
						currentColor := palette[currentIdx].(color.RGBA)

						// 线性混合（保留更多颜色细节）
						mixedR := uint8(float64(currentColor.R)*(1-weight) + float64(p.color.R)*weight)
						mixedG := uint8(float64(currentColor.G)*(1-weight) + float64(p.color.G)*weight)
						mixedB := uint8(float64(currentColor.B)*(1-weight) + float64(p.color.B)*weight)
						mixedColor := color.RGBA{mixedR, mixedG, mixedB, 255}

						// 精确映射到调色板
						mixedIdx := closestColorIndex(palette, mixedColor)
						img.SetColorIndex(x, y, uint8(mixedIdx))
					}
				}
			}

			gifImg.Image[f] = img
			gifImg.Delay[f] = delay
		}(frame)
	}

	wg.Wait()

	file, err := os.Create("high_quality_diffusion_4.gif")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := gif.EncodeAll(file, gifImg); err != nil {
		panic(err)
	}
}
