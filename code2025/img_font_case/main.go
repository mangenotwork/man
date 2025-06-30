package main

import (
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"os"
)

// 配置参数
type Config struct {
	InputImagePath  string
	OutputImagePath string
	Text            string
	ShadowColor     color.RGBA // 阴影颜色
	PositionX       int
	PositionY       int
	FontSize        float64
	ShadowOffsetX   int
	ShadowOffsetY   int
	ShadowBlur      int
}

func main() {
	// 解析命令行参数
	config := parseFlags()

	// 打开输入图片
	inputFile, err := os.Open(config.InputImagePath)
	if err != nil {
		log.Fatalf("无法打开输入图片: %v", err)
	}
	defer inputFile.Close()

	// 解码图片
	img, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("无法解码图片: %v", err)
	}

	// 创建可绘制的RGBA图像
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, image.Point{}, draw.Src)

	// 在图片上添加带阴影的文字（文字透明，仅显示阴影）
	addTransparentTextWithShadow(dst, config)

	// 创建输出文件
	outputFile, err := os.Create(config.OutputImagePath)
	if err != nil {
		log.Fatalf("无法创建输出图片: %v", err)
	}
	defer outputFile.Close()

	// 保存为PNG格式
	if err := png.Encode(outputFile, dst); err != nil {
		log.Fatalf("无法编码图片: %v", err)
	}

	fmt.Printf("已成功生成图片: %s\n", config.OutputImagePath)
}

func parseFlags() Config {
	var config Config

	// 配置参数（灰色半透明阴影）
	config.InputImagePath = "./img.png"     // 输入图片路径
	config.OutputImagePath = "./output.png" // 输出图片路径
	config.Text = "Hello, World! 这是透明文字的阴影"
	config.ShadowColor = color.RGBA{100, 100, 100, 160} // 阴影颜色(RGBA)
	config.PositionX = 500
	config.PositionY = 160
	config.FontSize = 36
	config.ShadowOffsetX = 2
	config.ShadowOffsetY = 2
	config.ShadowBlur = 20

	return config
}

func addTransparentTextWithShadow(dst draw.Image, config Config) {
	// 加载字体文件（示例使用思源黑体，需提前下载到同目录）
	fontFile, err := os.Open("./SourceHanSansCN-Regular.otf")
	if err != nil {
		log.Printf("使用默认字体: %v", err)
		//addWithDefaultFont(dst, config)
		return
	}
	defer fontFile.Close()

	// 读取字体数据
	fontData, err := io.ReadAll(fontFile)
	if err != nil {
		log.Fatalf("读取字体数据失败: %v", err)
	}

	// 解析字体
	fontObj, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatalf("解析字体失败: %v", err)
	}

	// 创建字体对象
	face, err := opentype.NewFace(fontObj, &opentype.FaceOptions{
		Size:    config.FontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("创建字体对象失败: %v", err)
	}

	// 绘制透明文字的阴影
	drawTransparentTextShadow(dst, config, face)
}

//// 使用默认字体绘制
//func addWithDefaultFont(dst draw.Image, config Config) {
//	// 使用基本字体作为回退
//	face := basicfont.Face7x13
//	if config.FontSize > 12 {
//		face = basicfont.Face9x15
//	}
//	drawTransparentTextShadow(dst, config, face)
//}

// 绘制透明文字的阴影核心逻辑
func drawTransparentTextShadow(dst draw.Image, config Config, face font.Face) {
	// 计算文本边界并调整位置
	textBounds := calculateTextBounds(config.Text, face)
	x, y := config.PositionX, config.PositionY+textBounds.Dy()

	// 创建透明文字的Alpha掩码（仅保留文字形状）
	textMask := createTransparentTextMask(config.Text, face, x, y, dst.Bounds())

	// 绘制阴影
	drawShadow(dst, textMask, image.Point{x, y}, config)
}

// 计算文本边界
func calculateTextBounds(text string, face font.Face) image.Rectangle {
	var width fixed.Int26_6
	height := face.Metrics().Height

	for _, r := range text {
		advance, ok := face.GlyphAdvance(r)
		if ok {
			width += advance
		}
	}

	w := int(width.Round())
	h := int(height.Round())
	return image.Rect(0, 0, w, h)
}

// 创建透明文字的Alpha掩码（文字本身透明，仅保留形状）
func createTransparentTextMask(text string, face font.Face, x, y int, imgBounds image.Rectangle) image.Image {
	textBounds := calculateTextBounds(text, face)
	padding := 10
	maskBounds := image.Rect(0, 0, textBounds.Dx()+padding*2, textBounds.Dy()+padding*2)
	mask := image.NewAlpha(maskBounds)

	// 绘制文字形状到Alpha掩码（完全不透明，仅用于生成形状）
	d := &font.Drawer{
		Dst:  mask,
		Src:  image.NewUniform(color.Alpha{255}), // 仅生成文字形状
		Face: face,
		Dot:  fixed.Point26_6{fixed.Int26_6(padding * 64), fixed.Int26_6((textBounds.Dy() - 2) * 64)},
	}
	d.DrawString(text)

	return mask
}

// 绘制阴影
func drawShadow(dst draw.Image, textMask image.Image, pos image.Point, config Config) {
	maskBounds := textMask.Bounds()
	imgBounds := dst.Bounds()

	// 创建阴影图像
	shadow := image.NewRGBA(imgBounds)

	// 绘制阴影（未模糊）
	for sy := maskBounds.Min.Y; sy < maskBounds.Max.Y; sy++ {
		for sx := maskBounds.Min.X; sx < maskBounds.Max.X; sx++ {
			// 获取文字形状的Alpha值
			a := textMask.At(sx, sy).(color.Alpha).A
			if a == 0 {
				continue
			}

			// 计算阴影在目标图像中的位置
			dstX := pos.X - maskBounds.Min.X + sx + config.ShadowOffsetX
			dstY := pos.Y - maskBounds.Min.Y + sy + config.ShadowOffsetY

			// 确保位置在图像范围内
			if dstX >= imgBounds.Min.X && dstX < imgBounds.Max.X &&
				dstY >= imgBounds.Min.Y && dstY < imgBounds.Max.Y {
				// 使用配置的阴影颜色和文字Alpha值
				shadow.Set(dstX, dstY, color.RGBA{
					R: config.ShadowColor.R,
					G: config.ShadowColor.G,
					B: config.ShadowColor.B,
					A: a, // 文字形状控制阴影透明度
				})
			}
		}
	}

	// 应用高斯模糊
	if config.ShadowBlur > 0 {
		shadow = gaussianBlur(shadow, config.ShadowBlur)
	}

	// 合并阴影到目标图像
	draw.Draw(dst, imgBounds, shadow, image.Point{}, draw.Over)
}

// 高斯模糊算法（优化版）
func gaussianBlur(src *image.RGBA, radius int) *image.RGBA {
	if radius <= 0 {
		return src
	}

	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	kernel := calculateGaussianKernel(radius)

	// 水平模糊
	horizontalBlur(src, dst, kernel)

	// 垂直模糊（使用临时图像）
	temp := image.NewRGBA(bounds)
	verticalBlur(dst, temp, kernel)
	return temp
}

// 计算高斯核
func calculateGaussianKernel(radius int) []float64 {
	kernel := make([]float64, 2*radius+1)
	sigma := float64(radius) / 3.0
	sum := 0.0

	for i := -radius; i <= radius; i++ {
		x := float64(i)
		kernel[i+radius] = math.Exp(-(x * x) / (2 * sigma * sigma))
		sum += kernel[i+radius]
	}

	// 归一化
	for i := range kernel {
		kernel[i] /= sum
	}
	return kernel
}

// 水平方向模糊
func horizontalBlur(src, dst *image.RGBA, kernel []float64) {
	bounds := src.Bounds()
	radius := (len(kernel) - 1) / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64
			for i, w := range kernel {
				ix := x + (i - radius)
				if ix < bounds.Min.X {
					ix = bounds.Min.X
				}
				if ix >= bounds.Max.X {
					ix = bounds.Max.X - 1
				}

				sr, sg, sb, sa := src.At(ix, y).RGBA()
				r += float64(sr>>8) * w
				g += float64(sg>>8) * w
				b += float64(sb>>8) * w
				a += float64(sa>>8) * w
			}
			dst.Set(x, y, color.RGBA{
				R: uint8(r + 0.5),
				G: uint8(g + 0.5),
				B: uint8(b + 0.5),
				A: uint8(a + 0.5),
			})
		}
	}
}

// 垂直方向模糊
func verticalBlur(src, dst *image.RGBA, kernel []float64) {
	bounds := src.Bounds()
	radius := (len(kernel) - 1) / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64
			for i, w := range kernel {
				iy := y + (i - radius)
				if iy < bounds.Min.Y {
					iy = bounds.Min.Y
				}
				if iy >= bounds.Max.Y {
					iy = bounds.Max.Y - 1
				}

				sr, sg, sb, sa := src.At(x, iy).RGBA()
				r += float64(sr>>8) * w
				g += float64(sg>>8) * w
				b += float64(sb>>8) * w
				a += float64(sa>>8) * w
			}
			dst.Set(x, y, color.RGBA{
				R: uint8(r + 0.5),
				G: uint8(g + 0.5),
				B: uint8(b + 0.5),
				A: uint8(a + 0.5),
			})
		}
	}
}
