package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 1. 打开图片文件
	file, err := os.Open("img.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()

	// 3. 创建可修改的图像副本
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 4. 定义毛玻璃效果区域
	frostedArea := image.Rect(520, 200, 1040, 380) // x1,y1,x2,y2

	// 5. 应用毛玻璃效果
	FrostedGlassEffect(rgba, frostedArea, 15) // 效果强度=8

	// 4. 定义毛玻璃效果区域
	frostedArea = image.Rect(600, 50, 800, 140) // x1,y1,x2,y2

	// 5. 应用毛玻璃效果
	FrostedGlassEffect(rgba, frostedArea, 8) // 效果强度=8

	// 6. 保存结果
	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	jpeg.Encode(out, rgba, &jpeg.Options{Quality: 90})
}

// FrostedGlassEffect 实现高效毛玻璃效果
func FrostedGlassEffect(img *image.RGBA, area image.Rectangle, intensity int) {
	if intensity < 1 {
		return
	}

	// 调整区域到图像边界内
	area = area.Intersect(img.Bounds())
	if area.Empty() {
		return
	}

	// 计算实际需要处理的区域（包含边缘扩展）
	processArea := image.Rect(
		max(0, area.Min.X-intensity),
		max(0, area.Min.Y-intensity),
		min(img.Bounds().Max.X, area.Max.X+intensity),
		min(img.Bounds().Max.Y, area.Max.Y+intensity),
	)

	// 创建处理区域的副本
	src := img.SubImage(processArea).(*image.RGBA)

	// 并行处理
	procs := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(procs)

	segment := area.Dy() / procs
	for i := 0; i < procs; i++ {
		go func(seg int) {
			defer wg.Done()
			startY := area.Min.Y + seg*segment
			endY := startY + segment
			if seg == procs-1 {
				endY = area.Max.Y
			}

			frostedGlassPass(img, src, area, processArea, startY, endY, intensity)
		}(i)
	}
	wg.Wait()
}

// frostedGlassPass 毛玻璃效果处理通道
func frostedGlassPass(dst, src *image.RGBA, area, processArea image.Rectangle, startY, endY, intensity int) {
	// 计算源图像和处理区域的偏移
	offsetX := processArea.Min.X
	offsetY := processArea.Min.Y

	// 处理指定范围内的每一行
	for y := startY; y < endY; y++ {
		for x := area.Min.X; x < area.Max.X; x++ {
			// 随机选择采样点
			rx := x + rand.Intn(2*intensity+1) - intensity
			ry := y + rand.Intn(2*intensity+1) - intensity

			// 确保采样点在处理区域内
			rx = clamp(rx, processArea.Min.X, processArea.Max.X-1)
			ry = clamp(ry, processArea.Min.Y, processArea.Max.Y-1)

			// 计算源图像中的位置
			srcX := rx - offsetX
			srcY := ry - offsetY

			// 安全获取源像素
			srcIdx := srcY*src.Stride + srcX*4
			if srcIdx < 0 || srcIdx+3 >= len(src.Pix) {
				continue // 跳过无效索引
			}

			r := src.Pix[srcIdx]
			g := src.Pix[srcIdx+1]
			b := src.Pix[srcIdx+2]
			a := src.Pix[srcIdx+3]

			// 设置目标像素
			dstIdx := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
			if dstIdx < 0 || dstIdx+3 >= len(dst.Pix) {
				continue // 跳过无效索引
			}

			dst.Pix[dstIdx] = r
			dst.Pix[dstIdx+1] = g
			dst.Pix[dstIdx+2] = b
			dst.Pix[dstIdx+3] = a
		}
	}
}

// clamp 确保值在指定范围内
func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
