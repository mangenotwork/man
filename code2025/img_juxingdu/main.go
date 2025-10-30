package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// 配置参数（可根据需求调整）
const (
	// 图像预处理参数（保持与之前流程兼容）
	cannySigma       = 1.5
	cannyHighThresh  = 100
	cannyLowThresh   = 50
	dilateKernelSize = 5
	dilateIterations = 2
	erodeKernelSize  = 5
	erodeIterations  = 2

	// 矩形度筛选参数
	minArea           = 50  // 最小区域面积（过滤过小区域）
	minRectangularity = 0.7 // 最小矩形度（越接近1越接近矩形）

	// 性能与路径配置
	maxProcessSize = 1024
	inputPath      = "input.png"
	edgeMaskPath   = "1_edge_mask.png"
	filledMaskPath = "3_filled_mask.png"
	labeledPath    = "labeled_rect.png"
)

// 8连通方向（用于连通域分析）
var dirs8 = []image.Point{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// 区域特征结构体（存储面积、最小外接矩形等信息）
type regionFeature struct {
	pixels         []image.Point // 区域像素点
	area           int           // 区域实际面积（像素数）
	minX, maxX     int           // 最小外接矩形边界
	minY, maxY     int
	mbrArea        int     // 最小外接矩形面积
	rectangularity float64 // 矩形度
}

// -------------------------- 工具函数 --------------------------

// 读取图像并转为灰度矩阵（支持缩放）
func ReadGrayImage(path string) ([][]uint8, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开图片失败: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("解码失败: %w", err)
	}

	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// 缩放超大图像
	scale := 1.0
	if w > maxProcessSize || h > maxProcessSize {
		scale = math.Min(float64(maxProcessSize)/float64(w), float64(maxProcessSize)/float64(h))
		w = int(float64(w) * scale)
		h = int(float64(h) * scale)
		fmt.Printf("图像缩放至 %dx%d（比例: %.2f）\n", w, h, scale)
	}

	gray := make([][]uint8, h)
	for y := 0; y < h; y++ {
		gray[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			srcX, srcY := int(float64(x)/scale), int(float64(y)/scale)
			r, g, b, _ := img.At(srcX, srcY).RGBA()
			gray[y][x] = uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
		}
	}
	return gray, nil
}

// 保存灰度图像
func SaveGrayImage(matrix [][]uint8, path string) error {
	h, w := len(matrix), len(matrix[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := matrix[y][x]
			img.SetRGBA(x, y, color.RGBA{val, val, val, 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()
	return png.Encode(file, img)
}

// 增强对比度
func EnhanceContrast(gray [][]uint8) [][]uint8 {
	h, w := len(gray), len(gray[0])
	minVal, maxVal := 255, 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := int(gray[y][x])
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}
	if maxVal == minVal {
		enhanced := make([][]uint8, h)
		for y := 0; y < h; y++ {
			enhanced[y] = make([]uint8, w)
			for x := 0; x < w; x++ {
				enhanced[y][x] = 128
			}
		}
		return enhanced
	}
	denominator := maxVal - minVal
	enhanced := make([][]uint8, h)
	for y := 0; y < h; y++ {
		enhanced[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			val := int(gray[y][x])
			enhanced[y][x] = uint8((val - minVal) * 255 / denominator)
		}
	}
	return enhanced
}

// -------------------------- 图像预处理（边缘检测+填充） --------------------------

// Canny边缘检测
func CannyEdgeMask(gray [][]uint8) [][]uint8 {
	//h, w := len(gray), len(gray[0])
	gaussKernel := generateGaussianKernel(cannySigma, 7)
	blurred := gaussianBlur(gray, gaussKernel)
	mag, dir := calculateGradients(blurred)
	suppressed := nonMaxSuppression(mag, dir)
	return hysteresis(suppressed)
}

// 生成高斯核
func generateGaussianKernel(sigma float64, size int) [][]float64 {
	kernel := make([][]float64, size)
	center := size / 2
	var sum float64
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := float64(x - center)
			dy := float64(y - center)
			kernel[y][x] = math.Exp(-(dx*dx+dy*dy)/(2*sigma*sigma)) / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}
	return kernel
}

// 高斯模糊
func gaussianBlur(gray [][]uint8, kernel [][]float64) [][]uint8 {
	h, w := len(gray), len(gray[0])
	ksize := len(kernel)
	radius := ksize / 2
	blurred := make([][]uint8, h)
	for y := 0; y < h; y++ {
		blurred[y] = make([]uint8, w)
	}
	offsets := make([]image.Point, 0, ksize*ksize)
	weights := make([]float64, 0, ksize*ksize)
	for ky := 0; ky < ksize; ky++ {
		for kx := 0; kx < ksize; kx++ {
			offsets = append(offsets, image.Point{kx - radius, ky - radius})
			weights = append(weights, kernel[ky][kx])
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sum float64
			for i, off := range offsets {
				nx, ny := x+off.X, y+off.Y
				if nx >= 0 && nx < w && ny >= 0 && ny < h {
					sum += float64(gray[ny][nx]) * weights[i]
				}
			}
			blurred[y][x] = clamp(sum)
		}
	}
	return blurred
}

// 计算梯度
func calculateGradients(blurred [][]uint8) (mag [][]float64, dir [][]float64) {
	h, w := len(blurred), len(blurred[0])
	mag = make([][]float64, h)
	dir = make([][]float64, h)
	for y := 0; y < h; y++ {
		mag[y] = make([]float64, w)
		dir[y] = make([]float64, w)
	}
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			p00 := int(blurred[y-1][x-1])
			p01 := int(blurred[y-1][x])
			p02 := int(blurred[y-1][x+1])
			p10 := int(blurred[y][x-1])
			p12 := int(blurred[y][x+1])
			p20 := int(blurred[y+1][x-1])
			p21 := int(blurred[y+1][x])
			p22 := int(blurred[y+1][x+1])
			gx := -p00 - 2*p10 - p20 + p02 + 2*p12 + p22
			gy := -p00 - 2*p01 - p02 + p20 + 2*p21 + p22
			mag[y][x] = math.Hypot(float64(gx), float64(gy))
			dir[y][x] = math.Atan2(float64(gy), float64(gx))
			if dir[y][x] < 0 {
				dir[y][x] += math.Pi
			}
		}
	}
	return
}

// 非极大值抑制
func nonMaxSuppression(mag [][]float64, dir [][]float64) [][]float64 {
	h, w := len(mag), len(mag[0])
	suppressed := make([][]float64, h)
	for y := 0; y < h; y++ {
		suppressed[y] = make([]float64, w)
	}
	pi4 := math.Pi / 4
	//pi2 := math.Pi / 2
	//pi34 := 3 * math.Pi / 4
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			d := dir[y][x]
			m := mag[y][x]
			if m == 0 {
				continue
			}
			var n1, n2 float64
			switch {
			case (d < pi4) || (d >= 7*pi4):
				n1, n2 = mag[y][x-1], mag[y][x+1]
			case d < 3*pi4:
				n1, n2 = mag[y-1][x+1], mag[y+1][x-1]
			case d < 5*pi4:
				n1, n2 = mag[y-1][x], mag[y+1][x]
			default:
				n1, n2 = mag[y-1][x-1], mag[y+1][x+1]
			}
			if m >= n1 && m >= n2 {
				suppressed[y][x] = m
			} else {
				suppressed[y][x] = 0
			}
		}
	}
	return suppressed
}

// hysteresis边缘连接
func hysteresis(suppressed [][]float64) [][]uint8 {
	h, w := len(suppressed), len(suppressed[0])
	edges := make([][]uint8, h)
	for y := range edges {
		edges[y] = make([]uint8, w)
	}
	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if suppressed[y][x] >= cannyHighThresh && !visited[y][x] {
				queue := []image.Point{{x, y}}
				visited[y][x] = true
				edges[y][x] = 255
				head := 0
				for head < len(queue) {
					p := queue[head]
					head++
					for _, d := range dirs8 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && !visited[ny][nx] && suppressed[ny][nx] >= cannyLowThresh {
							visited[ny][nx] = true
							edges[ny][nx] = 255
							queue = append(queue, image.Point{nx, ny})
						}
					}
				}
			}
		}
	}
	return edges
}

// 形态学膨胀
func Dilate(mask [][]uint8, kernelSize, iterations int) [][]uint8 {
	h, w := len(mask), len(mask[0])
	radius := kernelSize / 2
	current := make([][]uint8, h)
	for y := 0; y < h; y++ {
		current[y] = make([]uint8, w)
		copy(current[y], mask[y])
	}
	next := make([][]uint8, h)
	for y := 0; y < h; y++ {
		next[y] = make([]uint8, w)
	}
	for i := 0; i < iterations; i++ {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x < radius || x >= w-radius || y < radius || y >= h-radius {
					next[y][x] = current[y][x]
					continue
				}
				hasForeground := false
				for ky := -radius; ky <= radius; ky++ {
					for kx := -radius; kx <= radius; kx++ {
						if current[y+ky][x+kx] == 255 {
							hasForeground = true
							break
						}
					}
					if hasForeground {
						break
					}
				}
				next[y][x] = 255
				if !hasForeground {
					next[y][x] = 0
				}
			}
		}
		current, next = next, current
	}
	return current
}

// 填充空洞（获取完整区域）
func FillHoles(mask [][]uint8) [][]uint8 {
	h, w := len(mask), len(mask[0])
	inverted := make([][]uint8, h)
	for y := 0; y < h; y++ {
		inverted[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			if mask[y][x] == 0 {
				inverted[y][x] = 255
			}
		}
	}
	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}
	// 标记边缘连通域（非空洞）
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (y == 0 || y == h-1 || x == 0 || x == w-1) && inverted[y][x] == 255 && !visited[y][x] {
				queue := []image.Point{{x, y}}
				visited[y][x] = true
				head := 0
				for head < len(queue) {
					p := queue[head]
					head++
					for _, d := range dirs8 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && inverted[ny][nx] == 255 && !visited[ny][nx] {
							visited[ny][nx] = true
							queue = append(queue, image.Point{nx, ny})
						}
					}
				}
			}
		}
	}
	// 填充空洞（未被标记的区域）
	filled := make([][]uint8, h)
	for y := 0; y < h; y++ {
		filled[y] = make([]uint8, w)
		copy(filled[y], mask[y])
		for x := 0; x < w; x++ {
			if inverted[y][x] == 255 && !visited[y][x] {
				filled[y][x] = 255
			}
		}
	}
	return filled
}

// -------------------------- 核心：矩形度计算与筛选 --------------------------

// 提取所有区域的特征（面积、矩形度等）
func extractRegionFeatures(mask [][]uint8) []regionFeature {
	h, w := len(mask), len(mask[0])
	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}
	var features []regionFeature

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if mask[y][x] == 255 && !visited[y][x] {
				// 8连通域分析
				queue := []image.Point{{x, y}}
				visited[y][x] = true
				head := 0
				var pixels []image.Point
				minX, maxX := x, x
				minY, maxY := y, y

				for head < len(queue) {
					p := queue[head]
					head++
					pixels = append(pixels, p)

					// 更新外接矩形边界
					if p.X < minX {
						minX = p.X
					}
					if p.X > maxX {
						maxX = p.X
					}
					if p.Y < minY {
						minY = p.Y
					}
					if p.Y > maxY {
						maxY = p.Y
					}

					// 扩展连通域
					for _, d := range dirs8 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && mask[ny][nx] == 255 && !visited[ny][nx] {
							visited[ny][nx] = true
							queue = append(queue, image.Point{nx, ny})
						}
					}
				}

				// 计算区域特征
				area := len(pixels)
				mbrWidth := maxX - minX + 1
				mbrHeight := maxY - minY + 1
				mbrArea := mbrWidth * mbrHeight
				rectangularity := float64(area) / float64(mbrArea)

				features = append(features, regionFeature{
					pixels:         pixels,
					area:           area,
					minX:           minX,
					maxX:           maxX,
					minY:           minY,
					maxY:           maxY,
					mbrArea:        mbrArea,
					rectangularity: rectangularity,
				})
			}
		}
	}
	return features
}

// 标记符合矩形度条件的区域
func LabelByRectangularity(gray [][]uint8, mask [][]uint8) image.Image {
	//h, w := len(gray), len(gray[0])
	fmt.Println("开始计算区域矩形度...")

	// 提取所有区域特征
	features := extractRegionFeatures(mask)
	if len(features) == 0 {
		fmt.Println("未检测到任何区域")
		return grayToRGBA(gray)
	}

	// 筛选符合条件的区域（面积和矩形度达标）
	var validFeatures []regionFeature
	for i, feat := range features {
		if feat.area < minArea {
			fmt.Printf("区域 %d：面积过小（%d < %d），跳过\n", i+1, feat.area, minArea)
			continue
		}
		if feat.rectangularity < minRectangularity {
			fmt.Printf("区域 %d：矩形度不足（%.2f < %.2f），跳过\n", i+1, feat.rectangularity, minRectangularity)
			continue
		}
		validFeatures = append(validFeatures, feat)
		fmt.Printf("区域 %d：保留（面积: %d, 矩形度: %.2f）\n", i+1, feat.area, feat.rectangularity)
	}

	if len(validFeatures) == 0 {
		fmt.Println("无符合条件的矩形区域")
		return grayToRGBA(gray)
	}

	// 标记有效区域
	colors := []color.RGBA{
		{255, 0, 0, 128},   // 红
		{0, 255, 0, 128},   // 绿
		{0, 0, 255, 128},   // 蓝
		{255, 255, 0, 128}, // 黄
		{255, 0, 255, 128}, // 品红
	}
	result := grayToRGBA(gray)
	for i, feat := range validFeatures {
		colorIdx := i % len(colors)
		c := colors[colorIdx]
		// 标记区域像素
		for _, p := range feat.pixels {
			result.SetRGBA(p.X, p.Y, c)
		}
		// 绘制最小外接矩形（边框）
		drawRect(result, feat.minX, feat.minY, feat.maxX, feat.maxY, color.RGBA{255, 255, 255, 255})
		// 标记区域编号
		drawNumber(result, (feat.minX+feat.maxX)/2, (feat.minY+feat.maxY)/2, i+1)
	}

	return result
}

// 辅助函数：灰度转RGBA
func grayToRGBA(gray [][]uint8) *image.RGBA {
	h, w := len(gray), len(gray[0])
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := gray[y][x]
			rgba.SetRGBA(x, y, color.RGBA{val, val, val, 255})
		}
	}
	return rgba
}

// 辅助函数：绘制矩形边框
func drawRect(img *image.RGBA, minX, minY, maxX, maxY int, c color.RGBA) {
	// 上边框
	for x := minX; x <= maxX; x++ {
		img.SetRGBA(x, minY, c)
	}
	// 下边框
	for x := minX; x <= maxX; x++ {
		img.SetRGBA(x, maxY, c)
	}
	// 左边框
	for y := minY; y <= maxY; y++ {
		img.SetRGBA(minX, y, c)
	}
	// 右边框
	for y := minY; y <= maxY; y++ {
		img.SetRGBA(maxX, y, c)
	}
}

// 辅助函数：绘制数字标记
func drawNumber(img *image.RGBA, x, y, num int) {
	// 简单绘制十字标记（代表编号位置）
	size := 3
	for dx := -size; dx <= size; dx++ {
		img.SetRGBA(x+dx, y, color.RGBA{255, 255, 255, 255})
		img.SetRGBA(x, y+dx, color.RGBA{255, 255, 255, 255})
	}
}

func clamp(val float64) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

// -------------------------- 主函数 --------------------------

func main() {
	fmt.Println("矩形度特征区域标记程序启动...")

	// 1. 读取并预处理图像
	originalGray, err := ReadGrayImage(inputPath)
	if err != nil {
		fmt.Printf("读取图像失败: %v\n", err)
		return
	}
	h, w := len(originalGray), len(originalGray[0])
	fmt.Printf("图像尺寸: %dx%d\n", w, h)

	enhancedGray := EnhanceContrast(originalGray)

	// 2. 边缘检测与区域填充
	edgeMask := CannyEdgeMask(enhancedGray)
	if err := SaveGrayImage(edgeMask, edgeMaskPath); err != nil {
		fmt.Printf("保存边缘掩码失败: %v\n", err)
		return
	}
	fmt.Println("步骤1：边缘掩码已保存")

	dilatedMask := Dilate(edgeMask, dilateKernelSize, dilateIterations)
	filledMask := FillHoles(dilatedMask)
	if err := SaveGrayImage(filledMask, filledMaskPath); err != nil {
		fmt.Printf("保存填充掩码失败: %v\n", err)
		return
	}
	fmt.Println("步骤2：填充后的区域掩码已保存")

	// 3. 基于矩形度标记区域
	labeledImg := LabelByRectangularity(enhancedGray, filledMask)
	file, err := os.Create(labeledPath)
	if err != nil {
		fmt.Printf("保存标记结果失败: %v\n", err)
		return
	}
	defer file.Close()
	if err := png.Encode(file, labeledImg); err != nil {
		fmt.Printf("编码标记图像失败: %v\n", err)
		return
	}
	fmt.Printf("步骤3：矩形度标记结果已保存至 %s\n", labeledPath)

	fmt.Println("所有操作完成！")
}
