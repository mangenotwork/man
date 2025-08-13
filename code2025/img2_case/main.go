package main

import (
	"encoding/base64"
	"fmt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	_ "image/gif" // 使用 _ 导入 image/gif、image/jpeg 和 image/png 包，这会调用这些包的 init 函数，从而注册相应的图像格式
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {

	// case1()

	// case2()

	// case3()

	// case4()

	// case5()

	// case6()

	// case7()

	// case8()

	// case9()

	// case10()

	// case11()

	// case12()

	// case13()

	// case14()

	// case15()

	case16()
}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// convert image to NRGBA
func convertToNRGBA(src image.Image) *image.NRGBA {
	srcBounds := src.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

func DecodeImage(filePath string) (img image.Image, err error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err = image.Decode(reader)

	return
}

// 根据文件名打开图片,并编码,返回编码对象和文件类型
func LoadImage(path string) (img image.Image, filetype string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, filetype, err = image.Decode(file)
	return
}

func case1() {
	_, fileType, err := LoadImage("./test2.jpg")
	log.Println("err = ", err)
	log.Println("fileType = ", fileType)
}

// 如果您有未知格式的图像数据，请使用图像。解码功能可以检测格式。公认的格式集是在运行时构建的，不限于标准包库中的格式。
// 图像格式包通常在init函数中注册其格式，主包将“下划线导入”此类包，仅用于格式注册的副作用。
// convertToPNG converts from any recognized format to PNG.
func convertToPNG(w io.Writer, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

//	=========================================================================================================  颜色融合
//
// Color 表示一个RGBA颜色
type Color struct {
	R, G, B, A float64 // 范围从0.0到1.0
}

// NewColor 创建一个新的RGBA颜色
func NewColor(r, g, b, a float64) Color {
	return Color{
		R: clamp(r, 0.0, 1.0),
		G: clamp(g, 0.0, 1.0),
		B: clamp(b, 0.0, 1.0),
		A: clamp(a, 0.0, 1.0),
	}
}

// 辅助函数：将值限制在min和max之间
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Add 颜色加法混合
func (c Color) Add(other Color) Color {
	return Color{
		R: clamp(c.R+other.R, 0.0, 1.0),
		G: clamp(c.G+other.G, 0.0, 1.0),
		B: clamp(c.B+other.B, 0.0, 1.0),
		A: clamp(c.A+other.A, 0.0, 1.0),
	}
}

// Subtract 颜色减法混合
func (c Color) Subtract(other Color) Color {
	return Color{
		R: clamp(c.R-other.R, 0.0, 1.0),
		G: clamp(c.G-other.G, 0.0, 1.0),
		B: clamp(c.B-other.B, 0.0, 1.0),
		A: clamp(c.A-other.A, 0.0, 1.0),
	}
}

// Multiply 颜色乘法混合（常用于变暗）
func (c Color) Multiply(other Color) Color {
	return Color{
		R: clamp(c.R*other.R, 0.0, 1.0),
		G: clamp(c.G*other.G, 0.0, 1.0),
		B: clamp(c.B*other.B, 0.0, 1.0),
		A: clamp(c.A*other.A, 0.0, 1.0),
	}
}

// Blend 颜色混合，根据比例混合两种颜色
func (c Color) Blend(other Color, ratio float64) Color {
	ratio = clamp(ratio, 0.0, 1.0)
	invRatio := 1.0 - ratio
	return Color{
		R: c.R*invRatio + other.R*ratio,
		G: c.G*invRatio + other.G*ratio,
		B: c.B*invRatio + other.B*ratio,
		A: c.A*invRatio + other.A*ratio,
	}
}

// String 返回颜色的字符串表示形式（用于调试）
func (c Color) String() string {
	return fmt.Sprintf("RGBA(%d, %d, %d, %.2f)",
		int(c.R*255), int(c.G*255), int(c.B*255), c.A)
}

// ToHex 返回颜色的十六进制表示形式
func (c Color) ToHex() string {
	return fmt.Sprintf("#%02X%02X%02X",
		int(c.R*255), int(c.G*255), int(c.B*255))
}

// 从HSV颜色空间转换到RGB
func HSVtoRGB(h, s, v float64) Color {
	h = math.Mod(h, 360.0)
	if h < 0 {
		h += 360.0
	}
	s = clamp(s, 0.0, 1.0)
	v = clamp(v, 0.0, 1.0)

	if s == 0 {
		return NewColor(v, v, v, 1.0)
	}

	h /= 60.0
	hi := int(h) % 6
	f := h - float64(int(h))
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	var r, g, b float64
	switch hi {
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

	return NewColor(r, g, b, 1.0)
}

func case2() {
	// 定义几个基本颜色变量
	red := NewColor(1.0, 0.0, 0.0, 1.0)
	green := NewColor(0.0, 1.0, 0.0, 1.0)
	blue := NewColor(0.0, 0.0, 1.0, 1.0)
	yellow := NewColor(1.0, 1.0, 0.0, 1.0)

	// 创建一个从HSV转换的颜色
	purple := HSVtoRGB(270, 0.8, 0.8)

	// 演示颜色混合
	fmt.Println("基本颜色:")
	fmt.Println("红色:", red.ToHex())
	fmt.Println("绿色:", green.ToHex())
	fmt.Println("蓝色:", blue.ToHex())
	fmt.Println("黄色:", yellow.ToHex())
	fmt.Println("紫色:", purple.ToHex())

	fmt.Println("\n颜色混合示例:")
	// 混合红色和绿色
	yellowMix := red.Blend(green, 0.5)
	fmt.Println("红色 + 绿色 =", yellowMix.ToHex(), yellowMix)

	// 混合蓝色和黄色
	greenMix := blue.Blend(yellow, 0.5)
	fmt.Println("蓝色 + 黄色 =", greenMix.ToHex(), greenMix)

	// 混合所有基本颜色
	mixed := red.Blend(green, 0.25).Blend(blue, 0.25).Blend(yellow, 0.25)
	fmt.Println("所有颜色混合 =", mixed.ToHex(), mixed)

	// 颜色减法
	darker := red.Subtract(NewColor(0.3, 0.0, 0.0, 1.0))
	fmt.Println("红色变暗 =", darker.ToHex(), darker)

	// 颜色乘法
	darker2 := red.Multiply(NewColor(0.7, 0.7, 0.7, 1.0))
	fmt.Println("红色乘暗 =", darker2.ToHex(), darker2)
}

//	=========================================================================================================  计算图片的饱和度

// 计算单个RGB像素的饱和度
// 返回值范围: 0.0 (灰度) 到 1.0 (最大饱和度)
func calculatePixelSaturation(r, g, b uint8) float64 {
	// 将uint8转换为0.0-1.0范围的浮点数
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	// 找到RGB中的最大值和最小值
	maxVal := max3(rNorm, gNorm, bNorm)
	minVal := min3(rNorm, gNorm, bNorm)

	// 如果是灰度（max == min），饱和度为0
	if maxVal == minVal {
		return 0.0
	}

	// 计算亮度
	luminance := (maxVal + minVal) / 2.0

	var saturation float64
	if luminance <= 0.5 {
		saturation = (maxVal - minVal) / (maxVal + minVal)
	} else {
		saturation = (maxVal - minVal) / (2.0 - maxVal - minVal)
	}

	return saturation
}

// 辅助函数：返回三个数中的最大值
func max3(a, b, c float64) float64 {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}

// 辅助函数：返回三个数中的最小值
func min3(a, b, c float64) float64 {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

// 计算图片的平均饱和度
func calculateImageSaturation(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 计算所有像素的饱和度并取平均值
	var totalSaturation float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的饱和度并累加
			totalSaturation += calculatePixelSaturation(r8, g8, b8)
		}
	}

	// 返回平均饱和度
	averageSaturation := totalSaturation / float64(totalPixels)
	return averageSaturation, nil
}

func case3() {
	filePath := "./output_case3.jpg"
	saturation, err := calculateImageSaturation(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 以百分比形式显示结果
	fmt.Printf("图片的平均饱和度: %.2f%%\n", saturation*100)
}

//	=========================================================================================================  计算图片的亮度值
// ITU-R BT.709 标准的亮度计算公式：
// Y = 0.2126*R + 0.7152*G + 0.0722*B
// 这个公式考虑了人眼对不同颜色的敏感度差异：
// 绿色对亮度感知贡献最大（71.52%）
// 红色次之（21.26%）
// 蓝色贡献最小（7.22%）

// 计算单个RGB像素的亮度
// 返回值范围: 0.0 (最暗) 到 1.0 (最亮)
func calculatePixelBrightness(r, g, b uint8) float64 {
	// 使用ITU-R BT.709标准的亮度计算公式
	// 人眼对绿色最敏感，红色次之，蓝色最不敏感
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	// 亮度公式：Y = 0.2126*R + 0.7152*G + 0.0722*B
	brightness := 0.2126*rNorm + 0.7152*gNorm + 0.0722*bNorm
	return brightness
}

// 计算图片的平均亮度
func calculateImageBrightness(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 计算所有像素的亮度并取平均值
	var totalBrightness float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的亮度并累加
			totalBrightness += calculatePixelBrightness(r8, g8, b8)
		}
	}

	// 返回平均亮度
	averageBrightness := totalBrightness / float64(totalPixels)
	return averageBrightness, nil
}

func case4() {
	filePath := "./output_case3.jpg"
	brightness, err := calculateImageBrightness(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 以百分比形式显示结果
	fmt.Printf("图片的平均亮度: %.2f%%\n", brightness*100)
}

//	=========================================================================================================  计算图片的对比度值
// 两种计算对比度的方法：
//
//1. Michelson 对比度（默认使用）：
//	公式：(最大亮度 - 最小亮度) / (最大亮度 + 最小亮度)
//	这种方法简单直观，通过计算图像中最亮和最暗像素的差异来衡量对比度
//	值范围：0.0（无对比度，完全灰阶）到 1.0（最大对比度）
//2. 标准差对比度（注释中提供）：
//	基于亮度值的统计分布，使用标准差除以平均值（变异系数）
//	这种方法能更好地反映整体图像的对比度分布情况
//	适用于需要更精确对比度评估的场景

// 计算单个RGB像素的亮度（使用ITU-R BT.709标准）
func getLuminance(r, g, b uint8) float64 {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0
	return 0.2126*rNorm + 0.7152*gNorm + 0.0722*bNorm
}

// 计算图片的对比度
// 对比度公式: (最大亮度 - 最小亮度) / (最大亮度 + 最小亮度)
// 返回值范围: 0.0 (无对比度) 到 1.0 (最大对比度)
func calculateImageContrast(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 初始化最大和最小亮度
	maxLuminance := 0.0
	minLuminance := 1.0
	var totalLuminance float64
	var totalLuminanceSquared float64

	// 遍历所有像素计算亮度统计值
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的亮度
			luminance := getLuminance(r8, g8, b8)

			// 更新最大和最小亮度
			if luminance > maxLuminance {
				maxLuminance = luminance
			}
			if luminance < minLuminance {
				minLuminance = luminance
			}

			// 累加亮度值用于计算标准差
			totalLuminance += luminance
			totalLuminanceSquared += luminance * luminance
		}
	}

	//// 方法1: 使用最大最小亮度计算对比度 (Michelson对比度)
	//// 这种方法适合简单场景
	//contrast := (maxLuminance - minLuminance) / (maxLuminance + minLuminance)

	// 方法2: 使用标准差计算对比度 (更精确但计算稍复杂)
	mean := totalLuminance / float64(totalPixels)
	variance := (totalLuminanceSquared / float64(totalPixels)) - (mean * mean)
	stdDev := math.Sqrt(variance)
	contrast := stdDev / mean // 变异系数作为对比度度量

	return contrast, nil
}

func case5() {
	filePath := "./output_case3.jpg"
	contrast, err := calculateImageContrast(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 以百分比形式显示结果
	fmt.Printf("图片的对比度: %.2f%%\n", contrast*100)
}

//	=========================================================================================================  计算图片的锐度值
// 使用 Sobel 算子来计算图像的锐度，工作原理如下：
//	首先将图像转换为亮度矩阵，忽略色彩信息，只关注明暗变化
//	使用 Sobel 边缘检测算子，该算子包含两个卷积核：
//	水平方向卷积核（检测垂直边缘）
//	垂直方向卷积核（检测水平边缘）
//	对每个像素应用这两个卷积核，计算梯度幅度（边缘强度）
//	所有像素的平均边缘强度作为图像的锐度值

// 计算单个像素的亮度
//func getLuminance(r, g, b uint8) float64 {
//	rNorm := float64(r) / 255.0
//	gNorm := float64(g) / 255.0
//	bNorm := float64(b) / 255.0
//	return 0.2126*rNorm + 0.7152*gNorm + 0.0722*bNorm
//}

// 使用Sobel算子计算锐度
// 锐度值越高，表示图像越清晰
func calculateImageSharpness(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width < 3 || height < 3 {
		return 0.0, fmt.Errorf("图片尺寸过小，无法计算锐度")
	}

	// 创建亮度矩阵
	luminance := make([][]float64, height)
	for y := 0; y < height; y++ {
		luminance[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			luminance[y][x] = getLuminance(r8, g8, b8)
		}
	}

	// Sobel算子 - 水平和垂直方向的卷积核
	sobelX := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	sobelY := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	// 计算边缘强度
	var totalEdgeStrength float64
	edgeCount := 0

	// 遍历每个像素（跳过边界像素）
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 应用Sobel算子
			var gx, gy float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					gx += luminance[y+ky][x+kx] * sobelX[ky+1][kx+1]
					gy += luminance[y+ky][x+kx] * sobelY[ky+1][kx+1]
				}
			}

			// 计算梯度幅度（边缘强度）
			edgeStrength := math.Sqrt(gx*gx + gy*gy)
			totalEdgeStrength += edgeStrength
			edgeCount++
		}
	}

	// 计算平均边缘强度作为锐度度量
	averageSharpness := totalEdgeStrength / float64(edgeCount)
	return averageSharpness, nil
}

func case6() {
	filePath := "./output_case3.jpg"
	sharpness, err := calculateImageSharpness(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 锐度值没有固定范围，相对值越高表示图像越清晰
	fmt.Printf("图片的锐度值: %.4f\n", sharpness)
}

//	=========================================================================================================  计算图片的曝光度值
//
// 曝光度分析结果
type ExposureResult struct {
	OverexposedRatio  float64 // 过曝像素比例 (0.0-1.0)
	UnderexposedRatio float64 // 欠曝像素比例 (0.0-1.0)
	AverageLuminance  float64 // 平均亮度 (0.0-1.0)
	ExposureRating    string  // 曝光评级: "过曝", "欠曝", "正常"
}

// 计算图片的曝光度
func calculateImageExposure(filePath string) (ExposureResult, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return ExposureResult{}, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return ExposureResult{}, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return ExposureResult{}, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return ExposureResult{}, fmt.Errorf("图片尺寸为零")
	}

	// 阈值定义
	overexposureThreshold := 0.9  // 亮度超过此值视为过曝
	underexposureThreshold := 0.1 // 亮度低于此值视为欠曝

	// 统计变量
	var overexposedCount, underexposedCount int
	var totalLuminance float64

	// 遍历所有像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算亮度
			luminance := getLuminance(r8, g8, b8)
			totalLuminance += luminance

			// 判断是否过曝或欠曝
			if luminance >= overexposureThreshold {
				overexposedCount++
			} else if luminance <= underexposureThreshold {
				underexposedCount++
			}
		}
	}

	// 计算比例
	overexposedRatio := float64(overexposedCount) / float64(totalPixels)
	underexposedRatio := float64(underexposedCount) / float64(totalPixels)
	averageLuminance := totalLuminance / float64(totalPixels)

	// 确定曝光评级
	var exposureRating string
	if overexposedRatio > 0.2 { // 超过20%像素过曝
		exposureRating = "过曝"
	} else if underexposedRatio > 0.4 { // 超过40%像素欠曝
		exposureRating = "欠曝"
	} else {
		exposureRating = "正常"
	}

	return ExposureResult{
		OverexposedRatio:  overexposedRatio,
		UnderexposedRatio: underexposedRatio,
		AverageLuminance:  averageLuminance,
		ExposureRating:    exposureRating,
	}, nil
}

func case7() {
	filePath := "./output_case3.jpg"
	result, err := calculateImageExposure(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	fmt.Printf("图片曝光度分析:\n")
	fmt.Printf("  过曝像素比例: %.2f%%\n", result.OverexposedRatio*100)
	fmt.Printf("  欠曝像素比例: %.2f%%\n", result.UnderexposedRatio*100)
	fmt.Printf("  平均亮度: %.2f%%\n", result.AverageLuminance*100)
	fmt.Printf("  曝光评级: %s\n", result.ExposureRating)
}

//	=========================================================================================================  计算图片的色温值
// 这个程序计算色温的原理如下：
//	分析图像中所有非暗部像素的 RGB 值，忽略过暗像素（避免黑色影响判断）
//	计算红、绿、蓝三通道的平均值，并以绿色为基准进行归一化处理
//	通过红蓝光比例计算色温值，使用经验模型将比例转换为开尔文温度
//	根据色温值判断图片的色调偏向：
//	低于 5000K：暖色调（偏黄 / 红色）
//	5000K-7000K：中性色调
//	高于 7000K：冷色调（偏蓝色）

// 计算图片的色温（开尔文）
func calculateImageColorTemperature(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 统计所有像素的RGB值总和（忽略极暗像素，避免影响计算）
	var totalR, totalG, totalB float64
	validPixels := 0
	var darkThreshold uint8 = 30 // 忽略暗部像素（避免黑色影响色温判断）

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 忽略暗部像素
			if r8 < darkThreshold && g8 < darkThreshold && b8 < darkThreshold {
				continue
			}

			totalR += float64(r8)
			totalG += float64(g8)
			totalB += float64(b8)
			validPixels++
		}
	}

	if validPixels == 0 {
		return 0.0, fmt.Errorf("图片过暗，无法计算色温")
	}

	// 计算平均RGB值
	avgR := totalR / float64(validPixels)
	avgG := totalG / float64(validPixels)
	avgB := totalB / float64(validPixels)

	// 归一化处理（以绿色为基准）
	rNorm := avgR / avgG
	bNorm := avgB / avgG

	// 计算红蓝光比例
	rRatio := rNorm / (rNorm + bNorm)
	bRatio := bNorm / (rNorm + bNorm)

	// 色温转换公式（基于经验模型）
	// 色温范围大致在2000K（暖黄）到10000K（冷蓝）之间
	var temperature float64
	if rRatio > bRatio {
		// 暖色调
		ratio := rRatio / bRatio
		temperature = 6500 - 4500*math.Min(1.0, (ratio-1.0)/2.0)
	} else {
		// 冷色调
		ratio := bRatio / rRatio
		temperature = 6500 + 3500*math.Min(1.0, (ratio-1.0)/2.0)
	}

	// 确保色温在合理范围内
	if temperature < 2000 {
		temperature = 2000
	} else if temperature > 10000 {
		temperature = 10000
	}

	return temperature, nil
}

func case8() {
	filePath := "./output_case3.jpg"
	temp, err := calculateImageColorTemperature(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 判断色温偏向
	var tone string
	if temp < 5000 {
		tone = "暖色调（偏黄/红）"
	} else if temp > 7000 {
		tone = "冷色调（偏蓝）"
	} else {
		tone = "中性色调"
	}

	fmt.Printf("图片的色温: %.0fK\n", temp)
	fmt.Printf("色调偏向: %s\n", tone)
}

//	=========================================================================================================  计算图片的色调值
// 计算图片的色调（Hue）需要将 RGB 颜色空间转换到 HSV（色相 - 饱和度 - 明度）颜色空间，
//其中 Hue（色相）代表了颜色的基本属性，如红色、绿色、蓝色等。色调通常用角度（0°-360°）表示不同的颜色
//
// 原理如下：
//	将每个像素的 RGB 值转换为 HSV 颜色空间，提取色调（Hue）值
//	色调用 0°-360° 的角度表示，不同角度对应不同颜色：
//	0°/360°：红色
//	60°：黄色
//	120°：绿色
//	180°：青色
//	240°：蓝色
//	300°：品红色
//	忽略过暗和低饱和度的像素（接近灰色的像素），避免影响色调判断
//	统计所有有效像素的色调分布，计算平均色调角度和主要色调类别
//
//	程序将色调分为 13 种常见类别，包括红色、橙色、黄色、绿色等，最终输出：
//	平均色调角度（0°-360°）
//	主要色调类别（图片中最占优势的颜色）

// 色调范围定义
type HueRange struct {
	Start    float64
	End      float64
	Category string
}

// 常见色调范围分类
var hueRanges = []HueRange{
	{345, 360, "红色"},
	{0, 15, "红色"},
	{15, 45, "橙色"},
	{45, 75, "黄色"},
	{75, 105, "绿色"},
	{105, 135, "青色"},
	{135, 165, "蓝色"},
	{165, 195, "靛蓝色"},
	{195, 225, "紫色"},
	{225, 255, "粉红色"},
	{255, 285, "品红色"},
	{285, 315, "紫红色"},
	{315, 345, "深红色"},
}

// 将RGB转换为HSV颜色空间，返回色调值（0-360）
func rgbToHue(r, g, b uint8) float64 {
	// 归一化到0.0-1.0范围
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
	minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := maxVal - minVal

	var hue float64

	// 如果delta为0，说明是灰色，没有色调
	if delta == 0 {
		return -1 // 用-1表示无色调（灰色）
	}

	// 计算色调
	switch {
	case maxVal == rNorm:
		hue = math.Mod((gNorm-bNorm)/delta, 6)
	case maxVal == gNorm:
		hue = (bNorm-rNorm)/delta + 2
	case maxVal == bNorm:
		hue = (rNorm-gNorm)/delta + 4
	}

	// 转换为0-360度
	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}

// 确定色调所属类别
func getHueCategory(hue float64) string {
	if hue < 0 {
		return "灰色/无色调"
	}

	for _, r := range hueRanges {
		if (hue >= r.Start && hue <= r.End) ||
			(r.Start > r.End && (hue >= r.Start || hue <= r.End)) {
			return r.Category
		}
	}
	return "未知"
}

// 计算图片的主色调
func calculateImageHue(filePath string) (float64, string, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, "", fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, "", fmt.Errorf("图片尺寸为零")
	}

	// 统计变量
	var totalHue float64
	hueCount := 0
	hueDistribution := make(map[string]int)
	nonGrayPixels := 0

	// 忽略过暗像素的阈值
	var darkThreshold uint8 = 30
	// 忽略低饱和度像素的阈值
	minSaturation := 0.15

	// 遍历所有像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 忽略过暗像素
			if r8 < darkThreshold && g8 < darkThreshold && b8 < darkThreshold {
				continue
			}

			// 计算饱和度（用于过滤低饱和度像素）
			rNorm := float64(r8) / 255.0
			gNorm := float64(g8) / 255.0
			bNorm := float64(b8) / 255.0
			maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
			minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
			var saturation float64
			if maxVal > 0 {
				saturation = (maxVal - minVal) / maxVal
			}

			// 忽略低饱和度像素（接近灰色）
			if saturation < minSaturation {
				continue
			}

			// 计算色调
			hue := rgbToHue(r8, g8, b8)
			if hue >= 0 {
				totalHue += hue
				hueCount++
				nonGrayPixels++

				// 统计色调分布
				category := getHueCategory(hue)
				hueDistribution[category]++
			}
		}
	}

	if nonGrayPixels == 0 {
		return 0.0, "灰色/无明显色调", nil
	}

	// 计算平均色调
	averageHue := totalHue / float64(hueCount)

	// 找到最主要的色调类别
	mainCategory := "未知"
	maxCount := 0
	for cat, count := range hueDistribution {
		if count > maxCount {
			maxCount = count
			mainCategory = cat
		}
	}

	return averageHue, mainCategory, nil
}

func case9() {
	filePath := "./output_case3.jpg"
	hue, category, err := calculateImageHue(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	fmt.Printf("图片色调分析:\n")
	if hue >= 0 {
		fmt.Printf("  平均色调角度: %.1f°\n", hue)
	}
	fmt.Printf("  主要色调类别: %s\n", category)
}

//	=========================================================================================================  计算图片的噪点值
//	原理如下：
//
//	首先将图像转换为亮度矩阵，专注于亮度通道的变化
//	使用 3x3 高斯模糊核创建图像的平滑版本，高斯模糊能有效去除高频噪声
//	计算原始图像与平滑图像之间的差异，这种差异主要来自于噪点
//	对所有像素的差异取平均值，作为整体噪点水平的度量

// 计算图片的噪点值
// 返回值越高，表示噪点越多
func calculateImageNoise(filePath string) (float64, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0.0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return 0.0, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return 0.0, fmt.Errorf("无法解码图片: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width < 3 || height < 3 {
		return 0.0, fmt.Errorf("图片尺寸过小，无法计算噪点")
	}

	// 创建亮度矩阵
	luminance := make([][]float64, height)
	for y := 0; y < height; y++ {
		luminance[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			luminance[y][x] = getLuminance(r8, g8, b8)
		}
	}

	// 使用3x3高斯模糊核创建平滑图像
	gaussianKernel := [3][3]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	kernelSum := 16.0 // 高斯核的总和

	// 创建平滑后的亮度矩阵
	smoothed := make([][]float64, height)
	for y := 0; y < height; y++ {
		smoothed[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 边界像素直接使用原始值
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				smoothed[y][x] = luminance[y][x]
				continue
			}

			// 应用高斯模糊
			var sum float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					sum += luminance[y+ky][x+kx] * gaussianKernel[ky+1][kx+1]
				}
			}
			smoothed[y][x] = sum / kernelSum
		}
	}

	// 计算原始图像与平滑图像的差异（噪点估计）
	var totalNoise float64
	noiseCount := 0

	// 忽略边缘像素
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 计算亮度差异的绝对值
			diff := math.Abs(luminance[y][x] - smoothed[y][x])
			totalNoise += diff
			noiseCount++
		}
	}

	// 计算平均噪点值
	averageNoise := totalNoise / float64(noiseCount)
	return averageNoise, nil
}

func case10() {
	filePath := "./output_case3.jpg"
	noise, err := calculateImageNoise(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 判断噪点等级
	var noiseLevel string
	switch {
	case noise < 0.01:
		noiseLevel = "极低"
	case noise < 0.03:
		noiseLevel = "低"
	case noise < 0.06:
		noiseLevel = "中等"
	case noise < 0.1:
		noiseLevel = "高"
	default:
		noiseLevel = "极高"
	}

	fmt.Printf("图片噪点分析:\n")
	fmt.Printf("  噪点值: %.4f\n", noise)
	fmt.Printf("  噪点等级: %s\n", noiseLevel)
}

//	=========================================================================================================  图片反锯齿
// 核心思路是通过超采样（Super Sampling）来实现 - 先创建更高分辨率的图像，绘制后再缩小到原始尺寸，从而自然产生平滑边缘。

// 超采样因子，值越大抗锯齿效果越好，但性能消耗也越大
const sampleFactor = 4

// 绘制一条线并应用反锯齿
func drawLineWithAA(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	// 创建高分辨率的临时图像用于超采样
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	superImg := image.NewRGBA(image.Rect(0, 0, width*sampleFactor, height*sampleFactor))

	// 在高分辨率图像上绘制线条
	drawLine(superImg, x0*sampleFactor, y0*sampleFactor,
		x1*sampleFactor, y1*sampleFactor, c)

	// 将高分辨率图像缩小到原始尺寸，实现反锯齿
	downsample(img, superImg)
}

// 基本的 Bresenham 线绘制算法
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1

	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy

	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// 将高分辨率图像缩小到目标尺寸
func downsample(dst, src *image.RGBA) {
	bounds := dst.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 对每个目标像素，取对应高分辨率区域的平均值
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a uint32
			count := 0

			// 计算高分辨率图像中的对应区域
			startX := x * sampleFactor
			startY := y * sampleFactor
			endX := startX + sampleFactor
			endY := startY + sampleFactor

			// 累加区域内所有像素的颜色值
			for sy := startY; sy < endY; sy++ {
				for sx := startX; sx < endX; sx++ {
					r1, g1, b1, a1 := src.At(sx, sy).RGBA()
					r += r1
					g += g1
					b += b1
					a += a1
					count++
				}
			}

			// 计算平均值并设置到目标像素
			avgR := uint8((r / uint32(count)) >> 8)
			avgG := uint8((g / uint32(count)) >> 8)
			avgB := uint8((b / uint32(count)) >> 8)
			avgA := uint8((a / uint32(count)) >> 8)

			dst.SetRGBA(x, y, color.RGBA{avgR, avgG, avgB, avgA})
		}
	}
}

// 辅助函数：计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func case11() {
	// 创建一个 400x400 的图像
	width, height := 400, 400
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充白色背景
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}

	// 绘制一条对角线，应用反锯齿
	red := color.RGBA{255, 0, 0, 255}
	drawLineWithAA(img, 50, 50, 350, 350, red)

	// 保存图像
	outputFile, err := os.Create("anti_aliasing_example.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	png.Encode(outputFile, img)
}

// =========================================================================================================  golang 实现处理输入图片的反锯齿

// 加载图片（支持PNG和JPEG）
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// 保存图片（根据扩展名自动选择格式）
func saveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch ext := filepath.Ext(path); ext {
	case ".png":
		return png.Encode(file, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		return png.Encode(file, img) // 默认保存为PNG
	}
}

// 将图片转换为RGBA格式以便像素级操作
func toRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

// 计算像素亮度（用于边缘检测）
func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	// 转换为0-1范围
	r = r >> 8
	g = g >> 8
	b = b >> 8
	// 亮度公式（ITU-R BT.601）
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}

// 检测是否为边缘像素（通过与周围像素的亮度差判断）
func isEdgePixel(img *image.RGBA, x, y int, threshold float64) bool {
	bounds := img.Bounds()
	if x <= bounds.Min.X || x >= bounds.Max.X-1 ||
		y <= bounds.Min.Y || y >= bounds.Max.Y-1 {
		return false // 边缘像素不处理
	}

	// 当前像素亮度
	current := luminance(img.At(x, y))

	// 检查上下左右四个方向的像素
	dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]
		neighbor := luminance(img.At(nx, ny))
		if absf64(current-neighbor) > threshold {
			return true
		}
	}
	return false
}

// 对边缘像素应用3x3高斯模糊平滑处理
func applyAntiAliasing(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	width, height := bounds.Max.X, bounds.Max.Y

	// 3x3高斯核（已归一化）
	gaussianKernel := [3][3]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	total := 16.0 // 核总和，用于归一化

	// 边缘检测阈值（可根据需要调整）
	edgeThreshold := 30.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 非边缘像素直接复制
			if !isEdgePixel(img, x, y, edgeThreshold) {
				result.Set(x, y, img.At(x, y))
				continue
			}

			// 对边缘像素应用高斯模糊
			var r, g, b, a float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					nx, ny := x+kx, y+ky
					// 处理边界
					if nx < 0 {
						nx = 0
					} else if nx >= width {
						nx = width - 1
					}
					if ny < 0 {
						ny = 0
					} else if ny >= height {
						ny = height - 1
					}

					// 获取像素值并应用权重
					pixel := img.At(nx, ny)
					pr, pg, pb, pa := pixel.RGBA()
					weight := gaussianKernel[ky+1][kx+1] / total

					r += float64(pr>>8) * weight
					g += float64(pg>>8) * weight
					b += float64(pb>>8) * weight
					a += float64(pa>>8) * weight
				}
			}

			// 设置处理后的像素值
			result.SetRGBA(x, y, color.RGBA{
				R: uint8(clamp(r, 0, 255)),
				G: uint8(clamp(g, 0, 255)),
				B: uint8(clamp(b, 0, 255)),
				A: uint8(clamp(a, 0, 255)),
			})
		}
	}
	return result
}

// 辅助函数：计算绝对值
func absf64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func case12() {
	// 输入输出文件路径
	inputPath := "output_case2.jpg"    // 替换为你的输入图片路径
	outputPath := "output_case2_1.jpg" // 处理后的输出路径

	// 加载图片
	img, err := loadImage(inputPath)
	if err != nil {
		panic("无法加载图片: " + err.Error())
	}

	// 转换为RGBA以便处理
	rgbaImg := toRGBA(img)

	// 应用反锯齿处理
	processedImg := applyAntiAliasing(rgbaImg)

	// 保存处理后的图片
	err = saveImage(processedImg, outputPath)
	if err != nil {
		panic("无法保存图片: " + err.Error())
	}

	println("反锯齿处理完成，结果已保存到", outputPath)
}

// ========================================================================================================= 实现对图像执行伽玛校正
// 伽玛校正（Gamma Correction）是一种用于调整图像亮度和对比度的数字图像处理技术，其核心目的是补偿人类视觉系统对光线强度的非线性感知，
// 以及设备（如显示器、相机）在图像采集或显示过程中产生的非线性失真
//
// 伽玛校正的本质是通过幂函数 “反向补偿” 非线性失真，让图像的亮度表现更贴近人类视觉的预期

// 对图像执行伽玛校正
func applyGammaCorrection(img image.Image, gamma float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建输出图像
	output := image.NewRGBA(bounds)

	// 计算伽玛的倒数（优化计算效率）
	invGamma := 1.0 / gamma

	// 预计算伽玛校正查找表（0-255范围）
	var gammaTable [256]uint8
	for i := 0; i < 256; i++ {
		// 将像素值归一化到0.0-1.0范围
		normalized := float64(i) / 255.0
		// 应用伽玛校正公式：I' = I^(1/γ)
		corrected := math.Pow(normalized, invGamma)
		// 转换回0-255范围并确保在有效区间内
		gammaTable[i] = clampGM(corrected * 255.0)
	}

	// 遍历每个像素并应用伽玛校正
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取当前像素的RGBA值
			r, g, b, a := img.At(x, y).RGBA()

			// 将16位值转换为8位（RGBA()返回的是0-65535范围）
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 应用预计算的伽玛校正表
			rCorrected := gammaTable[r8]
			gCorrected := gammaTable[g8]
			bCorrected := gammaTable[b8]

			// 设置校正后的像素值（alpha通道保持不变）
			output.SetRGBA(x, y, color.RGBA{
				R: rCorrected,
				G: gCorrected,
				B: bCorrected,
				A: a8,
			})
		}
	}

	return output
}

// 确保值在0-255范围内
func clampGM(value float64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value + 0.5) // 四舍五入
}

func case13() {
	// 输入输出文件路径
	inputPath := "output_case2.jpg"    // 替换为你的输入图像路径
	outputPath := "output_case2_2.jpg" // 校正后的输出路径
	gammaValue := 0.2                  // 伽玛值（<1使图像变暗，>1使图像变亮）

	// 加载图像
	img, err := loadImage(inputPath)
	if err != nil {
		panic("无法加载图像: " + err.Error())
	}

	// 应用伽玛校正
	correctedImg := applyGammaCorrection(img, gammaValue)

	// 保存校正后的图像
	err = saveImage(correctedImg, outputPath)
	if err != nil {
		panic("无法保存图像: " + err.Error())
	}

	println("伽玛校正完成，伽玛值为", gammaValue)
	println("输出文件:", outputPath)
}

// ========================================================================================================= 直方图

// 计算直方图数据
func calculateHistogram(img image.Image) (rHist, gHist, bHist [256]int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 遍历每个像素，统计各通道像素值出现次数
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// 将16位值转换为8位（0-255范围）
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 对应通道的直方图计数加1
			rHist[r8]++
			gHist[g8]++
			bHist[b8]++
		}
	}

	return rHist, gHist, bHist
}

// 找到直方图中的最大值（用于归一化）
func findMax(hist [256]int) int {
	max := 0
	for _, v := range hist {
		if v > max {
			max = v
		}
	}
	return max
}

// 绘制直方图图像
func drawHistogram(rHist, gHist, bHist [256]int, width, height int) image.Image {
	// 创建输出图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充白色背景
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}

	// 找到三个通道的最大计数值
	maxR := findMax(rHist)
	maxG := findMax(gHist)
	maxB := findMax(bHist)
	maxCount := maxR
	if maxG > maxCount {
		maxCount = maxG
	}
	if maxB > maxCount {
		maxCount = maxB
	}

	// 计算每个bin的宽度
	binWidth := width / 256

	// 绘制三个通道的直方图
	for i := 0; i < 256; i++ {
		// 计算每个通道的高度（归一化到图像高度）
		rHeight := int(float64(rHist[i]) / float64(maxCount) * float64(height-10))
		gHeight := int(float64(gHist[i]) / float64(maxCount) * float64(height-10))
		bHeight := int(float64(bHist[i]) / float64(maxCount) * float64(height-10))

		// 计算x坐标
		x := i * binWidth

		// 绘制红色通道（R）
		drawBar(img, x, height-1, binWidth-1, rHeight, color.RGBA{255, 0, 0, 200})

		// 绘制绿色通道（G）
		drawBar(img, x, height-1-rHeight, binWidth-1, gHeight, color.RGBA{0, 255, 0, 200})

		// 绘制蓝色通道（B）
		drawBar(img, x, height-1-rHeight-gHeight, binWidth-1, bHeight, color.RGBA{0, 0, 255, 200})
	}

	return img
}

// 绘制单个直方图条
func drawBar(img *image.RGBA, x, y, width, height int, c color.RGBA) {
	// 确保高度不为负
	if height <= 0 {
		return
	}

	// 计算起始y坐标（图像坐标系y轴向下为正）
	startY := y - height

	// 绘制矩形条
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			imgX := x + dx
			imgY := startY + dy
			// 确保坐标在图像范围内
			if imgX < img.Bounds().Max.X && imgY >= 0 && imgY < img.Bounds().Max.Y {
				img.Set(imgX, imgY, c)
			}
		}
	}
}

func case14() {
	// 配置参数
	inputPath := "test2.jpg"      // 输入图像路径
	outputPath := "histogram.png" // 直方图输出路径
	histWidth := 800              // 直方图图像宽度
	histHeight := 400             // 直方图图像高度

	// 加载图像
	img, err := loadImage(inputPath)
	if err != nil {
		panic("无法加载图像: " + err.Error())
	}

	// 计算直方图数据
	rHist, gHist, bHist := calculateHistogram(img)

	// 绘制直方图
	histImage := drawHistogram(rHist, gHist, bHist, histWidth, histHeight)

	// 保存直方图
	err = saveImage(histImage, outputPath)
	if err != nil {
		panic("无法保存直方图: " + err.Error())
	}

	println("直方图生成完成，已保存到", outputPath)
}

// ========================================================================================================= 返回直方图的值

// 修改图像的直方图数值（即调整像素值的分布），本质上是改变图像中像素的亮度、颜色或对比度分布，最终会直接影响图像的视觉效果
// 图像的直方图本质是 “像素值分布的数学描述”，修改直方图的过程，其实是通过数学映射（如灰度变换）调整每个像素的实际值。这种调整直接反
// 映在图像的亮度、对比度、颜色偏向和细节表现上，是图像增强、校正（如伽玛校正）、风格化处理的核心原理。例如，相机的 “对比度调节”“亮度
// 调节” 功能，本质就是在动态修改图像的直方图分布。

// 计算图像的RGB三通道直方图，返回归一化后的float64切片
// 返回值依次为：红色通道直方图、绿色通道直方图、蓝色通道直方图
func CalculateNormalizedHistograms(img image.Image) (rHist, gHist, bHist []float64) {
	// 初始化三个通道的直方图数组（0-255共256个区间）
	rCounts := [256]int{}
	gCounts := [256]int{}
	bCounts := [256]int{}

	// 获取图像边界并计算总像素数
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height
	if totalPixels == 0 {
		// 空图像返回全零切片
		return make([]float64, 256), make([]float64, 256), make([]float64, 256)
	}

	// 遍历每个像素，统计各通道像素值出现次数
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值（返回的是0-65535范围的16位值）
			r, g, b, _ := img.At(x, y).RGBA()

			// 转换为8位值（0-255范围）
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 对应通道的计数加1
			rCounts[r8]++
			gCounts[g8]++
			bCounts[b8]++
		}
	}

	// 归一化：将计数值转换为占总像素数的比例（0.0-1.0）
	rHist = make([]float64, 256)
	gHist = make([]float64, 256)
	bHist = make([]float64, 256)

	for i := 0; i < 256; i++ {
		rHist[i] = float64(rCounts[i]) / float64(totalPixels)
		gHist[i] = float64(gCounts[i]) / float64(totalPixels)
		bHist[i] = float64(bCounts[i]) / float64(totalPixels)
	}

	return rHist, gHist, bHist
}

func case15() {
	// 示例用法
	inputPath := "test2.jpg" // 替换为你的图像路径

	// 加载图像
	file, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开图像文件: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("无法解码图像: " + err.Error())
	}

	// 计算直方图
	rHist, gHist, bHist := CalculateNormalizedHistograms(img)

	// 打印部分结果（示例）
	println("红色通道直方图部分值:")
	for i := 0; i < 10; i++ { // 打印前10个值
		fmt.Printf("值 %d: %.6f\n", i, rHist[i])
	}

	println("\n绿色通道直方图部分值:")
	for i := 245; i < 256; i++ { // 打印最后11个值
		fmt.Printf("值 %d: %.6f\n", i, gHist[i])
	}

	println("\n蓝色通道直方图部分值:")
	for i := 245; i < 256; i++ { // 打印最后11个值
		fmt.Printf("值 %d: %.6f\n", i, bHist[i])
	}
}

// ========================================================================================================= 实现提高图片亮度，使用并发来提高处理速度

// 并发调整图像亮度
// brightnessDelta: 亮度增量（-255到255之间，正值提高亮度，负值降低亮度）
func adjustBrightnessConcurrent(img image.Image, brightnessDelta int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds)

	// 转换为RGBA用于快速像素访问
	rgbaImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgbaImg.Set(x, y, img.At(x, y))
		}
	}

	// 根据CPU核心数确定并发goroutine数量
	numWorkers := runtime.NumCPU()
	rowsPerWorker := height / numWorkers
	var wg sync.WaitGroup

	// 分割图像行，分配给不同的goroutine处理
	for i := 0; i < numWorkers; i++ {
		startY := i * rowsPerWorker
		endY := startY + rowsPerWorker

		// 最后一个worker处理剩余的行
		if i == numWorkers-1 {
			endY = height
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			// 处理当前范围内的所有行
			for y := start; y < end; y++ {
				for x := 0; x < width; x++ {
					// 获取像素的RGBA值（0-65535范围）
					r, g, b, a := rgbaImg.At(x, y).RGBA()

					// 转换为8位值（0-255）并调整亮度
					r8 := clampInt(int(r>>8)+brightnessDelta, 0, 255)
					g8 := clampInt(int(g>>8)+brightnessDelta, 0, 255)
					b8 := clampInt(int(b>>8)+brightnessDelta, 0, 255)
					a8 := uint8(a >> 8) // Alpha通道不变

					// 设置调整后的像素值
					output.SetRGBA(x, y, color.RGBA{
						R: uint8(r8),
						G: uint8(g8),
						B: uint8(b8),
						A: a8,
					})
				}
			}
		}(startY, endY)
	}

	// 等待所有goroutine完成
	wg.Wait()
	return output
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func case16() {
	// 配置参数
	inputPath := "test2.jpg"    // 输入图像路径
	outputPath := "test2_1.jpg" // 输出图像路径
	brightnessDelta := 50       // 亮度增量（可调整范围：-255到255）

	// 加载图像
	img, err := loadImage(inputPath)
	if err != nil {
		panic("无法加载图像: " + err.Error())
	}

	// 并发调整亮度
	adjustedImg := adjustBrightnessConcurrent(img, brightnessDelta)

	// 保存调整后的图像
	err = saveImage(adjustedImg, outputPath)
	if err != nil {
		panic("无法保存图像: " + err.Error())
	}

	println("亮度调整完成，增量为", brightnessDelta)
	println("输出文件:", outputPath)
}
