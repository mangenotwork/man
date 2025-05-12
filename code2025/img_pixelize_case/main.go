package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

// go run main.go -input mv2.jpg -block 12 -colors 128 -stroke 12 -vibrancy 2.2 -edge 1600
// go run main.go -input mv2.jpg -block 12 -colors 132 -stroke 12 -vibrancy 2.2 -edge 1600
// go run main.go -input mv2.jpg -block 12 -colors 164 -stroke 12 -vibrancy 2.2 -edge 1600
// go run main.go -input mv3.jpg -block 14 -colors 187 -stroke 14 -vibrancy 2.0 -edge 1800
// go run main.go -input mv5.jpg -block 12 -colors 128 -stroke 12 -vibrancy 2.0 -edge 1600

func main() {
	// 定义命令行参数
	inputPath := flag.String("input", "", "输入图片路径")
	outputPath := flag.String("output", "", "输出图片路径（默认在输入文件名后添加_pixel）")
	blockSize := flag.Int("block", 10, "像素块大小（像素）")
	format := flag.String("format", "auto", "输出格式（auto/jpeg/png）")
	colorCount := flag.Int("colors", 8, "卡通化颜色数量")
	strokeWidth := flag.Int("stroke", 1, "描边宽度（像素）")
	enhanceVibrancy := flag.Float64("vibrancy", 1.5, "颜色鲜艳度增强因子")
	edgeThreshold := flag.Int("edge", 1000, "边缘检测阈值")
	flag.Parse()

	// 检查输入参数
	if *inputPath == "" {
		log.Fatal("请提供输入图片路径 (-input)")
	}

	// 处理输出路径
	if *outputPath == "" {
		ext := filepath.Ext(*inputPath)
		base := (*inputPath)[:len(*inputPath)-len(ext)]
		*outputPath = fmt.Sprintf("%s_pixel%s", base, ext)
	}

	// 读取输入图片
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		log.Fatalf("无法打开输入文件: %v", err)
	}
	defer inputFile.Close()

	// 解码图片
	src, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("无法解码图片: %v", err)
	}

	// 转换为像素图并卡通化
	pixelImage := pixelize(src, *blockSize, *colorCount, *strokeWidth, *enhanceVibrancy, *edgeThreshold)

	// 创建输出文件
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	// 根据格式保存图片
	switch *format {
	case "auto":
		ext := filepath.Ext(*outputPath)
		switch ext {
		case ".jpg", ".jpeg":
			err = jpeg.Encode(outputFile, pixelImage, &jpeg.Options{Quality: 90})
		case ".png":
			err = png.Encode(outputFile, pixelImage)
		default:
			log.Fatalf("不支持的文件格式: %s", ext)
		}
	case "jpeg":
		err = jpeg.Encode(outputFile, pixelImage, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outputFile, pixelImage)
	default:
		log.Fatalf("不支持的输出格式: %s", *format)
	}

	if err != nil {
		log.Fatalf("无法保存图片: %v", err)
	}

	fmt.Printf("成功生成卡通像素图: %s\n", *outputPath)
}

// HSV 颜色模型结构
type HSV struct {
	H, S, V float64
}

// RGBToHSV 将 RGB 颜色转换为 HSV
func RGBToHSV(r, g, b uint8) HSV {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	max := math.Max(rNorm, math.Max(gNorm, bNorm))
	min := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := max - min

	var h, s, v float64
	v = max

	if delta == 0 {
		h = 0
		s = 0
	} else {
		s = delta / max

		switch max {
		case rNorm:
			h = ((gNorm - bNorm) / delta)
			if h < 0 {
				h += 6
			}
		case gNorm:
			h = ((bNorm - rNorm) / delta) + 2
		case bNorm:
			h = ((rNorm - gNorm) / delta) + 4
		}

		h *= 60 // 转换为度
	}

	return HSV{H: h, S: s, V: v}
}

// HSVToRGB 将 HSV 颜色转换为 RGB
func HSVToRGB(hsv HSV) (uint8, uint8, uint8) {
	h, s, v := hsv.H, hsv.S, hsv.V

	if s == 0 {
		// 灰色
		r := uint8(v * 255)
		return r, r, r
	}

	h /= 60 // 转换为 0-6 区间
	sector := math.Floor(h)
	f := h - sector

	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float64

	switch sector {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

// enhanceColorVibrancy 增强颜色鲜艳度
func enhanceColorVibrancy(c color.RGBA, factor float64) color.RGBA {
	if factor <= 1.0 {
		return c // 无需增强
	}

	hsv := RGBToHSV(c.R, c.G, c.B)

	// 只增强不饱和颜色的饱和度，避免过饱和
	if hsv.S < 0.8 {
		hsv.S = math.Min(1.0, hsv.S*factor)
	}

	// 适当提高亮度
	hsv.V = math.Min(1.0, hsv.V*1.05)

	r, g, b := HSVToRGB(hsv)
	return color.RGBA{R: r, G: g, B: b, A: c.A}
}

// pixelize 将原始图片转换为带描边效果的卡通像素图
func pixelize(src image.Image, blockSize, colorCount, strokeWidth int, vibrancyFactor float64, edgeThreshold int) image.Image {
	// 获取图片边界
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建目标图像
	dst := image.NewRGBA(bounds)

	// 计算所有像素块的平均颜色，并存储在二维数组中
	blockColors := make([][]color.RGBA, (height+blockSize-1)/blockSize)
	for y := range blockColors {
		blockColors[y] = make([]color.RGBA, (width+blockSize-1)/blockSize)
	}

	// 遍历每个像素块，计算平均颜色
	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			// 计算当前块的边界
			endX := min(x+blockSize, width)
			endY := min(y+blockSize, height)

			// 计算块内所有像素的平均颜色
			r, g, b, a := 0, 0, 0, 0
			count := 0

			for y1 := y; y1 < endY; y1++ {
				for x1 := x; x1 < endX; x1++ {
					c := src.At(x1, y1)
					r1, g1, b1, a1 := c.RGBA()
					r += int(r1 >> 8)
					g += int(g1 >> 8)
					b += int(b1 >> 8)
					a += int(a1 >> 8)
					count++
				}
			}

			if count > 0 {
				// 计算平均颜色
				r /= count
				g /= count
				b /= count
				a /= count

				// 存储块颜色
				blockY := y / blockSize
				blockX := x / blockSize
				blockColors[blockY][blockX] = color.RGBA{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
					A: uint8(a),
				}
			}
		}
	}

	// 提取所有块颜色用于创建调色板
	allColors := make([]color.RGBA, 0, len(blockColors)*len(blockColors[0]))
	for y := range blockColors {
		for x := range blockColors[y] {
			allColors = append(allColors, blockColors[y][x])
		}
	}

	// 使用中位切割算法进行颜色量化，创建卡通调色板
	palette := medianCutQuantize(allColors, colorCount)

	// 增强调色板颜色的鲜艳度
	for i := range palette {
		palette[i] = enhanceColorVibrancy(palette[i], vibrancyFactor)
	}

	// 创建边缘检测结果数组
	edges := make([][]bool, len(blockColors))
	for y := range edges {
		edges[y] = make([]bool, len(blockColors[y]))
	}

	// 检测边缘
	for y := 0; y < len(blockColors); y++ {
		for x := 0; x < len(blockColors[y]); x++ {
			// 获取当前块的颜色
			currentColor := findNearestColor(blockColors[y][x], palette)

			// 检查周围块的颜色差异
			isEdge := false

			// 检查右邻居
			if x+1 < len(blockColors[y]) {
				rightColor := findNearestColor(blockColors[y][x+1], palette)
				if colorDistance(currentColor, rightColor) > edgeThreshold {
					isEdge = true
				}
			}

			// 检查下邻居
			if y+1 < len(blockColors) {
				bottomColor := findNearestColor(blockColors[y+1][x], palette)
				if colorDistance(currentColor, bottomColor) > edgeThreshold {
					isEdge = true
				}
			}

			edges[y][x] = isEdge
		}
	}

	// 绘制像素图
	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			blockY := y / blockSize
			blockX := x / blockSize

			if blockY < len(blockColors) && blockX < len(blockColors[blockY]) {
				// 获取调色板中最接近的颜色
				nearestColor := findNearestColor(blockColors[blockY][blockX], palette)

				// 计算当前块的边界
				endX := min(x+blockSize, width)
				endY := min(y+blockSize, height)

				// 填充当前块
				block := image.Rect(x, y, endX, endY)
				draw.Draw(dst, block, &image.Uniform{nearestColor}, image.Point{}, draw.Src)
			}
		}
	}

	// 绘制描边
	for y := 0; y < len(edges); y++ {
		for x := 0; x < len(edges[y]); x++ {
			if edges[y][x] {
				// 计算当前块在图像中的像素坐标
				px := x * blockSize
				py := y * blockSize
				pxEnd := min(px+blockSize, width)
				pyEnd := min(py+blockSize, height)

				// 获取当前块的颜色
				currentColor := findNearestColor(blockColors[y][x], palette)

				// 生成自适应描边颜色
				strokeColor := getAdaptiveStrokeColor(currentColor)

				// 绘制右侧描边
				if x+1 < len(edges[y]) && edges[y][x+1] {
					for dy := py; dy < pyEnd; dy++ {
						for sw := 0; sw < strokeWidth; sw++ {
							if pxEnd+sw < width {
								dst.Set(pxEnd+sw, dy, strokeColor)
							}
						}
					}
				}

				// 绘制下侧描边
				if y+1 < len(edges) && edges[y+1][x] {
					for dx := px; dx < pxEnd; dx++ {
						for sw := 0; sw < strokeWidth; sw++ {
							if pyEnd+sw < height {
								dst.Set(dx, pyEnd+sw, strokeColor)
							}
						}
					}
				}

				// 绘制右下角像素
				if x+1 < len(edges[y]) && y+1 < len(edges) && edges[y+1][x+1] {
					for sw := 0; sw < strokeWidth; sw++ {
						for sh := 0; sh < strokeWidth; sh++ {
							if pxEnd+sw < width && pyEnd+sh < height {
								dst.Set(pxEnd+sw, pyEnd+sh, strokeColor)
							}
						}
					}
				}
			}
		}
	}

	return dst
}

// getAdaptiveStrokeColor 根据色块颜色生成自适应描边颜色
func getAdaptiveStrokeColor(c color.RGBA) color.RGBA {
	hsv := RGBToHSV(c.R, c.G, c.B)

	// 深色块使用亮色描边，亮色块使用暗色描边
	if hsv.V < 0.5 {
		// 暗色块：增加亮度和饱和度
		hsv.V = math.Min(1.0, hsv.V+0.3)
		hsv.S = math.Min(1.0, hsv.S+0.2)
	} else {
		// 亮色块：降低亮度
		hsv.V = math.Max(0.0, hsv.V-0.3)
	}

	r, g, b := HSVToRGB(hsv)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// medianCutQuantize 使用中位切割算法进行颜色量化
func medianCutQuantize(colors []color.RGBA, targetCount int) []color.RGBA {
	if len(colors) == 0 || targetCount <= 0 {
		return []color.RGBA{}
	}

	// 如果颜色数量已经小于等于目标数量，直接返回
	if len(colors) <= targetCount {
		return uniqueColors(colors)
	}

	// 创建初始盒子
	boxes := []colorBox{{colors: colors}}

	// 不断分割盒子直到达到目标颜色数量
	for len(boxes) < targetCount {
		// 找到最大的盒子进行分割
		maxBoxIndex := 0
		maxRange := 0

		for i, box := range boxes {
			rRange, gRange, bRange := box.colorRange()
			maxBoxRange := max(rRange, max(gRange, bRange))
			if maxBoxRange > maxRange {
				maxRange = maxBoxRange
				maxBoxIndex = i
			}
		}

		// 分割最大的盒子
		bestBox := boxes[maxBoxIndex]
		box1, box2 := bestBox.split()

		// 替换原盒子为两个新盒子
		boxes = append(boxes[:maxBoxIndex], boxes[maxBoxIndex+1:]...)
		boxes = append(boxes, box1, box2)
	}

	// 计算每个盒子的平均颜色作为调色板
	palette := make([]color.RGBA, 0, len(boxes))
	for _, box := range boxes {
		palette = append(palette, box.averageColor())
	}

	return palette
}

// colorBox 表示一个颜色盒子，用于中位切割算法
type colorBox struct {
	colors []color.RGBA
}

// colorRange 计算盒子中颜色的RGB范围
func (b colorBox) colorRange() (int, int, int) {
	if len(b.colors) == 0 {
		return 0, 0, 0
	}

	minR, maxR := 255, 0
	minG, maxG := 255, 0
	minB, maxB := 255, 0

	for _, c := range b.colors {
		minR = min(int(c.R), minR)
		maxR = max(int(c.R), maxR)
		minG = min(int(c.G), minG)
		maxG = max(int(c.G), maxG)
		minB = min(int(c.B), minB)
		maxB = max(int(c.B), maxB)
	}

	return maxR - minR, maxG - minG, maxB - minB
}

// split 分割颜色盒子
func (b colorBox) split() (colorBox, colorBox) {
	if len(b.colors) <= 1 {
		return b, colorBox{}
	}

	// 计算RGB范围
	rRange, gRange, bRange := b.colorRange()

	// 确定要分割的轴（范围最大的颜色通道）
	axis := 0 // 0:R, 1:G, 2:B
	if gRange > rRange && gRange > bRange {
		axis = 1
	} else if bRange > rRange && bRange > gRange {
		axis = 2
	}

	// 复制颜色并按选定的轴排序
	sortedColors := make([]color.RGBA, len(b.colors))
	copy(sortedColors, b.colors)

	switch axis {
	case 0: // R
		for i := 0; i < len(sortedColors)-1; i++ {
			for j := i + 1; j < len(sortedColors); j++ {
				if sortedColors[i].R > sortedColors[j].R {
					sortedColors[i], sortedColors[j] = sortedColors[j], sortedColors[i]
				}
			}
		}
	case 1: // G
		for i := 0; i < len(sortedColors)-1; i++ {
			for j := i + 1; j < len(sortedColors); j++ {
				if sortedColors[i].G > sortedColors[j].G {
					sortedColors[i], sortedColors[j] = sortedColors[j], sortedColors[i]
				}
			}
		}
	case 2: // B
		for i := 0; i < len(sortedColors)-1; i++ {
			for j := i + 1; j < len(sortedColors); j++ {
				if sortedColors[i].B > sortedColors[j].B {
					sortedColors[i], sortedColors[j] = sortedColors[j], sortedColors[i]
				}
			}
		}
	}

	// 在中间分割
	mid := len(sortedColors) / 2
	return colorBox{colors: sortedColors[:mid]}, colorBox{colors: sortedColors[mid:]}
}

// averageColor 计算盒子中所有颜色的平均值
func (b colorBox) averageColor() color.RGBA {
	if len(b.colors) == 0 {
		return color.RGBA{}
	}

	var sumR, sumG, sumB, sumA uint64

	for _, c := range b.colors {
		sumR += uint64(c.R)
		sumG += uint64(c.G)
		sumB += uint64(c.B)
		sumA += uint64(c.A)
	}

	count := uint64(len(b.colors))
	return color.RGBA{
		R: uint8(sumR / count),
		G: uint8(sumG / count),
		B: uint8(sumB / count),
		A: uint8(sumA / count),
	}
}

// uniqueColors 从颜色列表中提取唯一的颜色
func uniqueColors(colors []color.RGBA) []color.RGBA {
	unique := make(map[color.RGBA]struct{})
	result := make([]color.RGBA, 0, len(colors))

	for _, c := range colors {
		if _, exists := unique[c]; !exists {
			unique[c] = struct{}{}
			result = append(result, c)
		}
	}

	return result
}

// findNearestColor 在调色板中找到最接近的颜色
func findNearestColor(c color.RGBA, palette []color.RGBA) color.RGBA {
	if len(palette) == 0 {
		return c
	}

	minDistance := math.MaxInt32
	nearest := palette[0]

	for _, p := range palette {
		dist := colorDistance(c, p)
		if dist < minDistance {
			minDistance = dist
			nearest = p
		}
	}

	return nearest
}

// colorDistance 计算两个颜色之间的欧几里得距离
func colorDistance(c1, c2 color.RGBA) int {
	dr := int(c1.R) - int(c2.R)
	dg := int(c1.G) - int(c2.G)
	db := int(c1.B) - int(c2.B)
	return dr*dr + dg*dg + db*db
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
