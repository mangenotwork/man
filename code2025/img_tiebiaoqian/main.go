package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// 配置参数（重点：标记区域的最小面积阈值）
const (
	// 边缘检测参数（保持不变）
	cannySigma      = 1.5
	cannyHighThresh = 100
	cannyLowThresh  = 50

	// 形态学参数（保持不变）
	dilateKernelSize = 5
	dilateIterations = 2
	erodeKernelSize  = 5
	erodeIterations  = 2
	erodeThreshold   = 50

	// 核心：最小标记面积（过小的区域不标记）
	minLabelArea = 500 // 可根据需求调整，如80、100等

	// 性能优化
	maxProcessSize = 1024

	// 路径（移除第4步相关路径）
	inputPath       = "input.png"
	edgeMaskPath    = "1_edge_mask.png"
	dilatedMaskPath = "2_dilated_mask.png"
	closedMaskPath  = "2.5_closed_mask.png"
	filledMaskPath  = "3_filled_mask.png"
	labeledPath     = "5_labeled.png" // 直接从第3步到第5步
)

// 8连通方向
var dirs8 = []image.Point{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// 4连通方向
var dirs4 = []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

// -------------------------- 工具函数 --------------------------

// ReadGrayImage 读取并转为灰度矩阵，超过最大尺寸则缩放
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
		fmt.Printf("图像过大，按比例 %.2f 缩放至 %dx%d\n", scale, w, h)
	}

	gray := make([][]uint8, h)
	for y := 0; y < h; y++ {
		gray[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			r, g, b, _ := img.At(srcX, srcY).RGBA()
			gray[y][x] = uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
		}
	}
	return gray, nil
}

// SaveGrayImage 保存灰度矩阵
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

// EnhanceContrast 增强对比度
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

// -------------------------- 步骤1：边缘检测 --------------------------

// CannyEdgeMask 生成边缘掩膜
func CannyEdgeMask(gray [][]uint8) [][]uint8 {
	//h, w := len(gray), len(gray[0])
	fmt.Println("步骤1：开始边缘检测...")

	gaussKernel := generateGaussianKernel(cannySigma, 7)
	blurred := gaussianBlur(gray, gaussKernel)
	mag, dir := calculateGradients(blurred)
	suppressed := nonMaxSuppression(mag, dir)
	edges := hysteresis(suppressed)

	return edges
}

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

func calculateGradients(blurred [][]uint8) (mag [][]float64, dir [][]float64) {
	h, w := len(blurred), len(blurred[0])
	mag = make([][]float64, h)
	dir = make([][]float64, h)
	for y := 0; y < h; y++ {
		mag[y] = make([]float64, w)
		dir[y] = make([]float64, w)
	}

	for y := 0; y < h; y++ {
		mag[y][0], dir[y][0] = 0, 0
		mag[y][w-1], dir[y][w-1] = 0, 0
	}
	for x := 0; x < w; x++ {
		mag[0][x], dir[0][x] = 0, 0
		mag[h-1][x], dir[h-1][x] = 0, 0
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

func nonMaxSuppression(mag [][]float64, dir [][]float64) [][]float64 {
	h, w := len(mag), len(mag[0])
	suppressed := make([][]float64, h)
	for y := 0; y < h; y++ {
		suppressed[y] = make([]float64, w)
	}

	for y := 0; y < h; y++ {
		suppressed[y][0] = 0
		suppressed[y][w-1] = 0
	}
	for x := 0; x < w; x++ {
		suppressed[0][x] = 0
		suppressed[h-1][x] = 0
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
		if y%100 == 0 && y > 0 {
			fmt.Printf("  边缘检测进度：%d/%d行\n", y, h)
		}

		for x := 0; x < w; x++ {
			if suppressed[y][x] >= cannyHighThresh && !visited[y][x] {
				queue := make([]image.Point, 0, 1024)
				queue = append(queue, image.Point{x, y})
				visited[y][x] = true
				edges[y][x] = 255
				head := 0

				for head < len(queue) {
					p := queue[head]
					head++

					for _, d := range dirs8 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && !visited[ny][nx] {
							if suppressed[ny][nx] >= cannyLowThresh {
								visited[ny][nx] = true
								edges[ny][nx] = 255
								queue = append(queue, image.Point{nx, ny})
							}
						}
					}
				}
			}
		}
	}

	return edges
}

// -------------------------- 步骤2：形态学操作 --------------------------

// Erode 腐蚀操作
func Erode(mask [][]uint8, kernelSize, iterations int) [][]uint8 {
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

				allForeground := true
				for ky := -radius; ky <= radius; ky++ {
					for kx := -radius; kx <= radius; kx++ {
						if current[y+ky][x+kx] == 0 {
							allForeground = false
							break
						}
					}
					if !allForeground {
						break
					}
				}

				if allForeground {
					next[y][x] = 255
				} else {
					next[y][x] = 0
				}
			}
		}

		current, next = next, current
	}

	return current
}

// Dilate 膨胀操作
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

				if hasForeground {
					next[y][x] = 255
				} else {
					next[y][x] = 0
				}
			}
		}

		current, next = next, current
	}

	return current
}

// MorphologicalClose 闭运算
func MorphologicalClose(mask [][]uint8) [][]uint8 {
	dilated := Dilate(mask, dilateKernelSize, dilateIterations)
	closed := Erode(dilated, erodeKernelSize, erodeIterations)
	return closed
}

// -------------------------- 步骤3：填充空洞 --------------------------

// FillHoles 填充空洞
func FillHoles(mask [][]uint8, originalGray [][]uint8) [][]uint8 {
	h, w := len(mask), len(mask[0])
	fmt.Println("步骤3：开始填充空洞...")

	inverted := make([][]uint8, h)
	for y := 0; y < h; y++ {
		inverted[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			if mask[y][x] == 0 {
				inverted[y][x] = 255
			} else {
				inverted[y][x] = 0
			}
		}
	}

	holes := findHoles(inverted)
	validHoles := filterValidHoles(holes, originalGray, mask)

	filled := make([][]uint8, h)
	for y := 0; y < h; y++ {
		filled[y] = make([]uint8, w)
		copy(filled[y], mask[y])
	}
	for _, hole := range validHoles {
		for _, p := range hole.pixels {
			filled[p.Y][p.X] = 255
		}
	}

	return filled
}

type holeRegion struct {
	pixels      []image.Point
	touchesEdge bool
}

func findHoles(bin [][]uint8) []holeRegion {
	h, w := len(bin), len(bin[0])
	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}
	var holes []holeRegion

	for y := 0; y < h; y++ {
		if y%100 == 0 && y > 0 {
			fmt.Printf("  空洞查找进度：%d/%d行\n", y, h)
		}

		for x := 0; x < w; x++ {
			if bin[y][x] == 255 && !visited[y][x] {
				queue := make([]image.Point, 0, 1024)
				queue = append(queue, image.Point{x, y})
				visited[y][x] = true
				var region []image.Point
				touchesEdge := false
				head := 0

				for head < len(queue) {
					p := queue[head]
					head++
					region = append(region, p)

					if p.X == 0 || p.X == w-1 || p.Y == 0 || p.Y == h-1 {
						touchesEdge = true
					}

					for _, d := range dirs4 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && bin[ny][nx] == 255 && !visited[ny][nx] {
							visited[ny][nx] = true
							queue = append(queue, image.Point{nx, ny})
						}
					}
				}

				if !touchesEdge {
					holes = append(holes, holeRegion{pixels: region, touchesEdge: touchesEdge})
				}
			}
		}
	}
	return holes
}

func filterValidHoles(holes []holeRegion, originalGray [][]uint8, mask [][]uint8) []holeRegion {
	var validHoles []holeRegion
	h, w := len(originalGray), len(originalGray[0])
	var bgSum, bgCount int

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if mask[y][x] == 0 {
				bgSum += int(originalGray[y][x])
				bgCount++
			}
		}
	}
	if bgCount == 0 {
		return holes
	}
	bgMean := float64(bgSum) / float64(bgCount)

	for _, hole := range holes {
		var holeSum int
		for _, p := range hole.pixels {
			holeSum += int(originalGray[p.Y][p.X])
		}
		holeMean := float64(holeSum) / float64(len(hole.pixels))
		if math.Abs(holeMean-bgMean) >= 5 {
			validHoles = append(validHoles, hole)
		}
	}
	return validHoles
}

// -------------------------- 步骤5：贴标签（核心：过滤小区域） --------------------------

// LabelRegions 基于第3步结果贴标签，过滤过小区域
func LabelRegions(originalGray [][]uint8, filledMask [][]uint8) image.Image {
	//h, w := len(filledMask), len(filledMask[0])
	fmt.Println("步骤5：开始标记区域（过滤过小区域）...")

	// 1. 查找第3步结果中的所有连通域
	components := findConnectedComponents(filledMask)
	if len(components) == 0 {
		fmt.Println("未检测到任何区域")
		return grayToRGBA(originalGray)
	}

	// 2. 过滤过小的区域（核心：仅保留面积≥minLabelArea的区域）
	var validComponents []connectedComponent
	for i, comp := range components {
		area := len(comp.pixels)
		if area < minLabelArea {
			fmt.Printf("  区域 %d：面积过小（%d < %d），不标记\n", i+1, area, minLabelArea)
			continue
		}
		validComponents = append(validComponents, comp)
		fmt.Printf("  区域 %d：面积合格（%d ≥ %d），标记\n", i+1, area, minLabelArea)
	}

	if len(validComponents) == 0 {
		fmt.Println("所有区域均过小，无标记结果")
		return grayToRGBA(originalGray)
	}

	// 3. 对有效区域进行标记
	colors := []color.RGBA{
		{255, 0, 0, 128},
		{0, 255, 0, 128},
		{0, 0, 255, 128},
		{255, 255, 0, 128},
		{255, 0, 255, 128},
		{0, 255, 255, 128},
	}

	result := grayToRGBA(originalGray)
	for i, comp := range validComponents {
		colorIdx := i % len(colors)
		c := colors[colorIdx]
		for _, p := range comp.pixels {
			result.SetRGBA(p.X, p.Y, c)
		}
		drawLabelNumber(result, comp, i+1)
	}

	return result
}

// findConnectedComponents 查找连通域（用于第3步结果分析）
func findConnectedComponents(bin [][]uint8) []connectedComponent {
	h, w := len(bin), len(bin[0])
	visited := make([][]bool, h)
	for y := range visited {
		visited[y] = make([]bool, w)
	}
	var components []connectedComponent

	for y := 0; y < h; y++ {
		if y%100 == 0 && y > 0 {
			fmt.Printf("  区域分析进度：%d/%d行\n", y, h)
		}

		for x := 0; x < w; x++ {
			if bin[y][x] == 255 && !visited[y][x] {
				queue := make([]image.Point, 0, 1024)
				queue = append(queue, image.Point{x, y})
				visited[y][x] = true
				var comp connectedComponent
				head := 0

				for head < len(queue) {
					p := queue[head]
					head++
					comp.pixels = append(comp.pixels, p)

					for _, d := range dirs8 {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < w && ny >= 0 && ny < h && bin[ny][nx] == 255 && !visited[ny][nx] {
							visited[ny][nx] = true
							queue = append(queue, image.Point{nx, ny})
						}
					}
				}

				components = append(components, comp)
			}
		}
	}
	return components
}

type connectedComponent struct {
	pixels []image.Point
}

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

func drawLabelNumber(img *image.RGBA, comp connectedComponent, label int) {
	var sumX, sumY int
	for _, p := range comp.pixels {
		sumX += p.X
		sumY += p.Y
	}
	cx := sumX / len(comp.pixels)
	cy := sumY / len(comp.pixels)

	crossSize := 5
	for dx := -crossSize; dx <= crossSize; dx++ {
		img.SetRGBA(cx+dx, cy, color.RGBA{255, 255, 255, 255})
		img.SetRGBA(cx, cy+dx, color.RGBA{255, 255, 255, 255})
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

// -------------------------- 主函数（移除第4步） --------------------------

func main() {
	fmt.Println("程序开始运行...")
	originalGray, err := ReadGrayImage(inputPath)
	if err != nil {
		fmt.Printf("读取原图失败: %v\n", err)
		return
	}
	h, w := len(originalGray), len(originalGray[0])
	if h == 0 || w == 0 {
		fmt.Println("无效图像尺寸")
		return
	}
	fmt.Printf("处理图像尺寸: %dx%d\n", w, h)

	enhancedGray := EnhanceContrast(originalGray)

	// 步骤1：边缘掩膜
	edgeMask := CannyEdgeMask(enhancedGray)
	if err := SaveGrayImage(edgeMask, edgeMaskPath); err != nil {
		fmt.Printf("保存边缘掩膜失败: %v\n", err)
		return
	}
	fmt.Println("步骤1：边缘梯度二值掩膜已保存")

	// 步骤2：膨胀与闭运算
	dilatedMask := Dilate(edgeMask, dilateKernelSize, dilateIterations)
	if err := SaveGrayImage(dilatedMask, dilatedMaskPath); err != nil {
		fmt.Printf("保存膨胀掩膜失败: %v\n", err)
		return
	}
	fmt.Println("步骤2：膨胀梯度掩膜已保存")

	closedMask := MorphologicalClose(dilatedMask)
	if err := SaveGrayImage(closedMask, closedMaskPath); err != nil {
		fmt.Printf("保存闭运算结果失败: %v\n", err)
		return
	}
	fmt.Println("步骤2.5：闭运算合并区域已保存")

	// 步骤3：填充空洞（保留良好效果）
	filledMask := FillHoles(closedMask, enhancedGray)
	if err := SaveGrayImage(filledMask, filledMaskPath); err != nil {
		fmt.Printf("保存填充空洞图像失败: %v\n", err)
		return
	}
	fmt.Println("步骤3：填充空洞后的二值图像已保存")

	// 直接从第3步到第5步（跳过第4步）
	labeledImg := LabelRegions(enhancedGray, filledMask)
	file, err := os.Create(labeledPath)
	if err != nil {
		fmt.Printf("保存贴标签图像失败: %v\n", err)
		return
	}
	defer file.Close()
	if err := png.Encode(file, labeledImg); err != nil {
		fmt.Printf("编码标签图像失败: %v\n", err)
		return
	}
	fmt.Println("步骤5：贴标签图像已保存（已过滤过小区域）")

	fmt.Println("所有步骤完成！")
}
