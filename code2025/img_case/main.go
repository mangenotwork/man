package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
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

	// case17()

	// case18()

	// case19()

	// case20()

	// case21()

	// case22()

	// case23()

	// case24()

	// case25()

	// case26()

	// case27()

	// case28()

	// case29()

	// case30()

	// case31()

	// case32()

	// case33()

	// case34()

	// case35()

	// case36()

	// case37()

	// case38()

	// case39()

	// case40()

	case41()

	// case42()

	// case43()

	// case44()

	// case45()

	// case46()

	// case47()

	// case48()

	// case49()

	// case50()

	// case51()

	// case52()

	// case53()

	// case54()

	// case55()

	// case56()

	// case57()

	// case58()

	// case59()

	// case60()
}

func getTestImg() image.Image {
	filePath := "./test.png"

	// 打开图像文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ext := strings.ToLower(filePath[strings.LastIndex(filePath, ".")+1:])

	var img image.Image

	// 解码图像
	switch ext {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported image format: %s", ext)
	}
	return img
}

func getTest2Img() image.Image {
	filePath := "./test2.jpg"

	// 打开图像文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ext := strings.ToLower(filePath[strings.LastIndex(filePath, ".")+1:])

	var img image.Image

	// 解码图像
	switch ext {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported image format: %s", ext)
	}
	return img
}

func getImg(filePath string) image.Image {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ext := strings.ToLower(filePath[strings.LastIndex(filePath, ".")+1:])

	var img image.Image

	// 解码图像
	switch ext {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported image format: %s", ext)
	}
	return img
}

// ========================================================================

// 图像点的运算通常涉及到读取、修改和处理图像中每个像素点的颜色值
func case1() {
	filePath := "./test.png"

	// 打开图像文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ext := strings.ToLower(filePath[strings.LastIndex(filePath, ".")+1:])

	var img image.Image

	// 解码图像
	switch ext {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported image format: %s", ext)
	}

	// 获取图像的边界
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	log.Println("Image size: %dx%d", width, height)

	// 遍历图像的每个像素点
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取像素点的颜色
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			// 这里可以对r, g, b, a进行运算
			log.Printf("Pixel at (%d, %d): R=%d, G=%d, B=%d, A=%d", x, y, r>>8, g>>8, b>>8, a>>8)
		}
	}
}

// ========================================================================

func case2() {
	img := getTestImg()

	// 创建一个新的图像对象
	bounds := img.Bounds()
	grayImg := image.NewRGBA(bounds)

	// 遍历图像的每个像素点并进行灰度化处理
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取像素点的颜色
			r, g, b, a := img.At(x, y).RGBA()
			// 计算灰度值
			gray := uint8((float64(r>>8)*0.299 + float64(g>>8)*0.587 + float64(b>>8)*0.114))
			// 设置新的颜色
			grayImg.Set(x, y, color.RGBA{gray, gray, gray, uint8(a >> 8)})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, grayImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}

}

// ========================================================================

// case3 图像点的亮度调整
func case3() {
	img := getTestImg()

	brightness := 50 // 亮度调整值

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newR := int(r>>8) + brightness
			newG := int(g>>8) + brightness
			newB := int(b>>8) + brightness

			// 确保颜色值在0-255范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			newImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a >> 8)})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case3.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}

}

// ========================================================================

// 图像加法可以用于合成图像或增加图像的亮度
func case4() {
	img1 := getTestImg()
	img2 := getTest2Img()

	bounds := img1.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			// 计算新的颜色值
			r := uint8((r1>>8 + r2>>8) / 2)
			g := uint8((g1>>8 + g2>>8) / 2)
			b := uint8((b1>>8 + b2>>8) / 2)
			a := uint8((a1>>8 + a2>>8) / 2)

			newImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case4.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// 图像减法可以用于检测图像中的变化或突出差异
func case5() {
	img1 := getTestImg()
	img2 := getTest2Img()

	bounds := img1.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			// 计算新的颜色值
			r := uint8((r1>>8 - r2>>8))
			g := uint8((g1>>8 - g2>>8))
			b := uint8((b1>>8 - b2>>8))
			a := uint8((a1>>8 - a2>>8))

			newImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case5.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// 图像乘法可以用于调整图像的亮度或对比度
func case6() {
	img := getTestImg()
	factor := 1.5 // 亮度调整因子

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// 计算新的颜色值
			newR := uint8(float64(r>>8) * factor)
			newG := uint8(float64(g>>8) * factor)
			newB := uint8(float64(b>>8) * factor)
			newA := uint8(float64(a>>8) * factor)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, newA})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case6.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// 图像除法可以用于调整图像的对比度
func case7() {
	img := getTestImg()
	divisor := 1.5 // 对比度调整因子

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// 计算新的颜色值
			newR := uint8(float64(r>>8) / divisor)
			newG := uint8(float64(g>>8) / divisor)
			newB := uint8(float64(b>>8) / divisor)
			newA := uint8(float64(a>>8) / divisor)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, newA})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case7.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// 图像的逻辑运算通常指的是对图像像素进行按位逻辑操作，例如与（AND）、或（OR）、异或（XOR）和非（NOT）操作。
// 这些操作可以用于图像的特征提取、图像掩码处理等场景
func case8() {
	img1 := getTestImg()
	img2 := getTest2Img()

	// 获取图像的边界
	bounds := img1.Bounds()
	// 创建一个新的 RGBA 图像
	newImg := image.NewRGBA(bounds)

	// 遍历图像的每个像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取当前像素在第一幅图像中的颜色值
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			// 获取当前像素在第二幅图像中的颜色值
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			//// 对红、绿、蓝和透明度通道进行按位与运算
			//r := uint8((r1 >> 8) & (r2 >> 8))
			//g := uint8((g1 >> 8) & (g2 >> 8))
			//b := uint8((b1 >> 8) & (b2 >> 8))
			//a := uint8((a1 >> 8) & (a2 >> 8))

			r := uint8((r1 >> 8) | (r2 >> 8))
			g := uint8((g1 >> 8) | (g2 >> 8))
			b := uint8((b1 >> 8) | (b2 >> 8))
			a := uint8((a1 >> 8) | (a2 >> 8))

			// 设置新图像中对应像素的颜色值
			newImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("output_case8.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 编码并保存新图像
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// case9 图像缩放  使用 golang.org/x/image/draw 包实现图像缩放
// 除了 draw.BiLinear，golang.org/x/image/draw 包还提供了其他的缩放算法，例如：
// draw.NearestNeighbor：最近邻插值，速度快，但可能会导致图像出现锯齿。
// draw.ApproxBiLinear：近似双线性插值，速度比 draw.BiLinear 快，但质量稍低。
// draw.CatmullRom：Catmull-Rom 插值，质量较高，但速度较慢。
func case9() {
	img1 := getTestImg()
	width := 800
	height := 600

	// 创建一个新的图像对象，用于存储缩放后的图像
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// 使用 draw.BiLinear 进行双线性插值缩放
	draw.BiLinear.Scale(dst, dst.Bounds(), img1, img1.Bounds(), draw.Over, nil)

	// 创建输出图像文件
	out, err := os.Create("output_case9.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// 编码并保存缩放后的图像
	err = jpeg.Encode(out, dst, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}
}

// ========================================================================

// case10 旋转图像
func case10() {
	src := getTestImg()

	rotated90 := Rotate90Clockwise(src)

	// 保存结果
	outputFile90, _ := os.Create("output_case10_90.jpg")
	jpeg.Encode(outputFile90, rotated90, nil)
	defer outputFile90.Close()

	// 旋转 180 度
	rotated180 := Rotate180(src)

	// 保存旋转后的图像
	outputFile180, err := os.Create("output_case10_180.jpg")
	if err != nil {
		panic(err)
	}
	defer outputFile180.Close()
	err = jpeg.Encode(outputFile180, rotated180, nil)
	if err != nil {
		panic(err)
	}

	// 顺时针旋转 270 度
	rotated270 := Rotate270Clockwise(src)

	// 保存旋转后的图像
	outputFile270, err := os.Create("output_case10_270.jpg")
	if err != nil {
		panic(err)
	}
	defer outputFile270.Close()
	err = jpeg.Encode(outputFile270, rotated270, nil)
	if err != nil {
		panic(err)
	}

}

// Rotate90Clockwise 顺时针旋转 90 度
func Rotate90Clockwise(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcHeight, srcWidth))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := y
			newY := srcWidth - 1 - x
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// Rotate180 旋转 180 度
func Rotate180(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := srcWidth - 1 - x
			newY := srcHeight - 1 - y
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// Rotate270Clockwise 顺时针旋转 270 度
func Rotate270Clockwise(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcHeight, srcWidth))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := srcHeight - 1 - y
			newY := x
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// ========================================================================

// case11 自定义旋转角度
func case11() {
	src := getTestImg()

	// 自定义旋转角度
	angle := 44.4
	rotated := Rotate(src, angle)

	// 保存旋转后的图像
	outputFileName := fmt.Sprintf("output_case11_%f.jpg", angle)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, rotated, nil)
	if err != nil {
		panic(err)
	}
}

// Rotate 自定义旋转函数，支持任意角度旋转
func Rotate(src image.Image, angle float64) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// 计算旋转后的图像大小
	rad := angle * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)

	x1 := math.Abs(float64(srcWidth)*cos) + math.Abs(float64(srcHeight)*sin)
	y1 := math.Abs(float64(srcWidth)*sin) + math.Abs(float64(srcHeight)*cos)

	dstWidth := int(math.Ceil(x1))
	dstHeight := int(math.Ceil(y1))

	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// 计算旋转中心
	srcCenterX := float64(srcWidth) / 2
	srcCenterY := float64(srcHeight) / 2
	dstCenterX := float64(dstWidth) / 2
	dstCenterY := float64(dstHeight) / 2

	// 遍历目标图像的每个像素
	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// 计算目标像素相对于旋转中心的坐标
			dx := float64(x) - dstCenterX
			dy := float64(y) - dstCenterY

			// 逆向旋转得到源图像中的坐标
			srcX := cos*dx + sin*dy + srcCenterX
			srcY := -sin*dx + cos*dy + srcCenterY

			// 检查源坐标是否在源图像范围内
			if srcX >= 0 && srcX < float64(srcWidth) && srcY >= 0 && srcY < float64(srcHeight) {
				// 双线性插值
				x0 := int(math.Floor(srcX))
				y0 := int(math.Floor(srcY))
				x1 := x0 + 1
				y1 := y0 + 1

				if x1 >= srcWidth {
					x1 = srcWidth - 1
				}
				if y1 >= srcHeight {
					y1 = srcHeight - 1
				}

				srcColor00 := src.At(x0, y0)
				srcColor01 := src.At(x0, y1)
				srcColor10 := src.At(x1, y0)
				srcColor11 := src.At(x1, y1)

				// 计算插值权重
				u := srcX - float64(x0)
				v := srcY - float64(y0)

				// 双线性插值计算颜色
				r0, g0, b0, a0 := interpolateColor(srcColor00, srcColor10, u)
				r1, g1, b1, a1 := interpolateColor(srcColor01, srcColor11, u)
				r, g, b, a := interpolateColor(color.RGBA{r0, g0, b0, a0}, color.RGBA{r1, g1, b1, a1}, v)

				dst.Set(x, y, color.RGBA{r, g, b, a})
			}
		}
	}

	return dst
}

// interpolateColor 双线性插值计算颜色
func interpolateColor(c1, c2 color.Color, t float64) (uint8, uint8, uint8, uint8) {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := uint8((1-t)*float64(r1>>8) + t*float64(r2>>8))
	g := uint8((1-t)*float64(g1>>8) + t*float64(g2>>8))
	b := uint8((1-t)*float64(b1>>8) + t*float64(b2>>8))
	a := uint8((1-t)*float64(a1>>8) + t*float64(a2>>8))

	return r, g, b, a
}

// ========================================================================

// case12 图片平移
func case12() {
	src := getTestImg()

	// 定义平移量
	dx := 50
	dy := 30

	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{image.Transparent}, image.Point{}, draw.Src)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			newX := x + dx
			newY := y + dy
			if newX >= 0 && newX < bounds.Dx() && newY >= 0 && newY < bounds.Dy() {
				dst.Set(newX, newY, src.At(x, y))
			}
		}
	}

	outputFileName := "output_case12.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, dst, nil)
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case13 图像裁剪
func case13() {
	src := getTestImg()

	// 定义裁剪区域
	cropRect := image.Rect(100, 100, 300, 300)

	dst := image.NewRGBA(cropRect)
	draw.Draw(dst, dst.Bounds(), src, cropRect.Min, draw.Src)

	outputFileName := "output_case13.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, dst, nil)
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case14 图像转置
func case14() {
	src := getTestImg()

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(y, x, src.At(x, y))
		}
	}

	outputFileName := "output_case14.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, dst, nil)
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case15 图像镜像
func case15() {
	src := getTestImg()

	// 水平镜像变换
	horizontalMirrored := horizontalMirror(src)
	// 创建水平镜像输出文件
	horizontalOutputFile, err := os.Create("output_case15_horizontal.jpg")
	if err != nil {
		log.Fatalf("无法创建水平镜像输出文件: %v", err)
	}
	defer horizontalOutputFile.Close()
	// 编码并保存水平镜像后的图像
	err = jpeg.Encode(horizontalOutputFile, horizontalMirrored, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("无法保存水平镜像图像: %v", err)
	}

	// 垂直镜像变换
	verticalMirrored := verticalMirror(src)
	// 创建垂直镜像输出文件
	verticalOutputFile, err := os.Create("output_case15_vertical.jpg")
	if err != nil {
		log.Fatalf("无法创建垂直镜像输出文件: %v", err)
	}
	defer verticalOutputFile.Close()
	// 编码并保存垂直镜像后的图像
	err = jpeg.Encode(verticalOutputFile, verticalMirrored, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("无法保存垂直镜像图像: %v", err)
	}

}

// horizontalMirror 函数用于实现图像水平镜像
func horizontalMirror(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(width-1-x, y, src.At(x, y))
		}
	}
	return dst
}

// verticalMirror 函数用于实现图像垂直镜像
func verticalMirror(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, height-1-y, src.At(x, y))
		}
	}
	return dst
}

// ========================================================================

// case16 计算图像的灰度直方图
// 人类的视觉系统对不同颜色的敏感度是不一样的，相较于蓝色，人眼对绿色更为敏感，对红色的敏感度则处于两者之间。所以在转换过程中，
// 不能简单地取 RGB 三个值的平均值，而是要依据它们的重要程度赋予不同的权重。具体来说，常用的转换公式为：
// Gray=0.299×R+0.587×G+0.114×B
// 这里，红色（R）的权重是 0.299，绿色（G）的权重是 0.587，蓝色（B）的权重是 0.114。
func calculateHistogram(img image.Image) [256]int {
	var histogram [256]int
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 转换为灰度值
			gray := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256
			index := int(gray)
			if index >= 0 && index < 256 {
				histogram[index]++
			}
		}
	}
	return histogram
}

// ========================================================================

// case17 使用双线性插值进行图像增强
func case17() {
	src := getTestImg()

	// 调整图像大小
	newWidth := src.Bounds().Dx() * 2
	newHeight := src.Bounds().Dy() * 2
	resized := resizeImage(src, newWidth, newHeight)

	// 调整对比度进行增强
	contrastFactor := 1.5
	enhanced := adjustContrast(resized, contrastFactor)

	outputFileName := "output_case17.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, enhanced, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// bilinearInterpolate 双线性插值函数
func bilinearInterpolate(img image.Image, x, y float64) color.RGBA {
	bounds := img.Bounds()
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	x1 := x0 + 1
	y1 := y0 + 1

	if x0 < bounds.Min.X {
		x0 = bounds.Min.X
	}
	if y0 < bounds.Min.Y {
		y0 = bounds.Min.Y
	}
	if x1 >= bounds.Max.X {
		x1 = bounds.Max.X - 1
	}
	if y1 >= bounds.Max.Y {
		y1 = bounds.Max.Y - 1
	}

	dx := x - float64(x0)
	dy := y - float64(y0)

	c00 := img.At(x0, y0)
	c01 := img.At(x0, y1)
	c10 := img.At(x1, y0)
	c11 := img.At(x1, y1)

	r0, g0, b0, a0 := c00.RGBA()
	r1, g1, b1, a1 := c01.RGBA()
	r2, g2, b2, a2 := c10.RGBA()
	r3, g3, b3, a3 := c11.RGBA()

	r := uint8((1-dx)*(1-dy)*float64(r0)/65535 + dx*(1-dy)*float64(r2)/65535 +
		(1-dx)*dy*float64(r1)/65535 + dx*dy*float64(r3)/65535)
	g := uint8((1-dx)*(1-dy)*float64(g0)/65535 + dx*(1-dy)*float64(g2)/65535 +
		(1-dx)*dy*float64(g1)/65535 + dx*dy*float64(g3)/65535)
	b := uint8((1-dx)*(1-dy)*float64(b0)/65535 + dx*(1-dy)*float64(b2)/65535 +
		(1-dx)*dy*float64(b1)/65535 + dx*dy*float64(b3)/65535)
	a := uint8((1-dx)*(1-dy)*float64(a0)/65535 + dx*(1-dy)*float64(a2)/65535 +
		(1-dx)*dy*float64(a1)/65535 + dx*dy*float64(a3)/65535)

	return color.RGBA{r, g, b, a}
}

// resizeImage 使用双线性插值缩放图像
func resizeImage(src image.Image, newWidth, newHeight int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := float64(x) * float64(width) / float64(newWidth)
			srcY := float64(y) * float64(height) / float64(newHeight)
			c := bilinearInterpolate(src, srcX, srcY)
			dst.Set(x, y, c)
		}
	}

	return dst
}

// adjustContrast 调整图像对比度
func adjustContrast(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r = uint32(math.Min(math.Max(0, ((float64(r)/65535-0.5)*factor+0.5)*65535), 65535))
			g = uint32(math.Min(math.Max(0, ((float64(g)/65535-0.5)*factor+0.5)*65535), 65535))
			b = uint32(math.Min(math.Max(0, ((float64(b)/65535-0.5)*factor+0.5)*65535), 65535))
			dst.Set(x, y, color.RGBA{
				uint8(r / 256),
				uint8(g / 256),
				uint8(b / 256),
				uint8(a / 256),
			})
		}
	}

	return dst
}

// ========================================================================

// case18 将彩色图转换为灰度图
func case18() {
	src := getTest2Img()

	// 转换为灰度图
	bounds := src.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取彩色图像该点的颜色
			r, g, b, _ := src.At(x, y).RGBA()
			// 将颜色值从 0 - 65535 转换为 0 - 255
			r = r >> 8
			g = g >> 8
			b = b >> 8
			// 计算灰度值，使用标准的加权公式
			grayValue := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			// 设置灰度图像对应点的灰度值
			grayImg.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	outputFileName := "output_case18.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, grayImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case19 转二值图像
// 二值图像是一种特殊的图像类型，其每个像素点只有两种取值，通常用 0 和 1 来表示，也可以用黑色和白色来直观地表示，所以二值图像也常被称为黑白图像
// 应用场景
// 文字识别：在光学字符识别（OCR）系统中，常常将输入的文字图像转换为二值图像，这样可以清晰地提取文字的轮廓信息，便于后续的特征提取和字符分类，提高识别准确率。
// 图像分割：通过将图像二值化，可以将图像中的目标物体与背景分离，从而实现图像分割。例如，在医学图像分析中，可将二值图像用于分割出人体器官或病变区域，辅助医生进行诊断。
// 计算机视觉：在计算机视觉领域，二值图像可用于目标检测、跟踪和识别等任务。例如，通过对监控视频中的图像进行二值化处理，可以快速检测出运动的物体，实现目标跟踪和行为分析。
// 文档处理：在文档扫描和处理中，将文档图像转换为二值图像可以减少存储空间，同时便于进行文字提取、格式转换等操作。此外，二值图像还可用于文档的数字水印嵌入和检测，保护文档的版权和完整性。
func case19() {
	src := getTest2Img()

	// 定义阈值
	threshold := uint8(128)

	bounds := src.Bounds()
	binaryImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取当前像素的灰度值
			gray := color.GrayModel.Convert(src.At(x, y)).(color.Gray)
			if gray.Y > threshold {
				// 大于阈值的像素设为白色
				binaryImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				// 小于阈值的像素设为黑色
				binaryImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	outputFileName := "output_case19.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, binaryImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case20 图像刚性变换、仿射变换和透视变换
// 这些实现使用最近邻插值，对于高质量图像处理建议添加双线性插值等抗锯齿技术

// RigidTransform 刚性变换（旋转、缩放、平移）
func RigidTransform(img image.Image, angle float64, scale float64, tx float64, ty float64) *image.RGBA {
	radian := angle * math.Pi / 180
	cos := math.Cos(radian)
	sin := math.Sin(radian)

	// 构建仿射变换矩阵参数
	mat := [6]float64{
		scale * cos, -scale * sin, tx,
		scale * sin, scale * cos, ty,
	}

	return affineTransform(img, mat)
}

// AffineTransform 仿射变换
func AffineTransform(img image.Image, mat [6]float64) *image.RGBA {
	return affineTransform(img, mat)
}

// PerspectiveTransform 透视变换
func PerspectiveTransform(img image.Image, mat [9]float64) *image.RGBA {
	a, b, c := mat[0], mat[1], mat[2]
	d, e, f := mat[3], mat[4], mat[5]
	g, h, i := mat[6], mat[7], mat[8]

	det := a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
	if det == 0 {
		return cloneImage(img)
	}

	invDet := 1.0 / det
	h11 := (e*i - f*h) * invDet
	h12 := (c*h - b*i) * invDet
	h13 := (b*f - c*e) * invDet
	h21 := (f*g - d*i) * invDet
	h22 := (a*i - c*g) * invDet
	h23 := (c*d - a*f) * invDet
	h31 := (d*h - e*g) * invDet
	h32 := (b*g - a*h) * invDet
	h33 := (a*e - b*d) * invDet

	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			xH := float64(x)
			yH := float64(y)
			wH := 1.0

			X := h11*xH + h12*yH + h13*wH
			Y := h21*xH + h22*yH + h23*wH
			W := h31*xH + h32*yH + h33*wH

			if W == 0 {
				dest.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			srcX := X / W
			srcY := Y / W
			srcXInt := int(math.Round(srcX))
			srcYInt := int(math.Round(srcY))

			if inBounds(img.Bounds(), srcXInt, srcYInt) {
				dest.Set(x, y, img.At(srcXInt, srcYInt))
			} else {
				dest.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return dest
}

// 仿射变换核心实现
func affineTransform(img image.Image, mat [6]float64) *image.RGBA {
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)

	a, b, c := mat[0], mat[1], mat[2]
	d, e, f := mat[3], mat[4], mat[5]

	det := a*e - b*d
	if det == 0 {
		return cloneImage(img)
	}

	invDet := 1.0 / det
	aPrime := e * invDet
	bPrime := -b * invDet
	cPrime := (-e*c + b*f) * invDet
	dPrime := -d * invDet
	ePrime := a * invDet
	fPrime := (d*c - a*f) * invDet

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcX := aPrime*float64(x) + bPrime*float64(y) + cPrime
			srcY := dPrime*float64(x) + ePrime*float64(y) + fPrime

			srcXInt := int(math.Round(srcX))
			srcYInt := int(math.Round(srcY))

			if inBounds(bounds, srcXInt, srcYInt) {
				dest.Set(x, y, img.At(srcXInt, srcYInt))
			} else {
				dest.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return dest
}

// 辅助函数：克隆图像
func cloneImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	clone := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			clone.Set(x, y, img.At(x, y))
		}
	}
	return clone
}

// 辅助函数：检查坐标是否在边界内
func inBounds(bounds image.Rectangle, x, y int) bool {
	return x >= bounds.Min.X && x < bounds.Max.X &&
		y >= bounds.Min.Y && y < bounds.Max.Y
}

func case20() {
	src := getTest2Img()

	dst1 := RigidTransform(src, 10, 1, 0, 0)

	outputFileName1 := "output_case20-1.jpg"
	outputFile1, err := os.Create(outputFileName1)
	if err != nil {
		panic(err)
	}
	defer outputFile1.Close()
	err = jpeg.Encode(outputFile1, dst1, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

	dst2 := AffineTransform(src, [6]float64{0, 50, 100, 100, 50, 0})
	outputFileName2 := "output_case20-2.jpg"
	outputFile2, err := os.Create(outputFileName2)
	if err != nil {
		panic(err)
	}
	defer outputFile2.Close()
	err = jpeg.Encode(outputFile2, dst2, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

	dst3 := PerspectiveTransform(src, [9]float64{0, 50, 100, 100, 0, 0, 50, 100, 100})
	outputFileName3 := "output_case20-3.jpg"
	outputFile3, err := os.Create(outputFileName3)
	if err != nil {
		panic(err)
	}
	defer outputFile3.Close()
	err = jpeg.Encode(outputFile3, dst3, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// 例子效果失败

// case21 图像的刚性变换，仿射变换，透视变换，添加双线性插值等抗锯齿技术

// 双线性插值用于根据浮点坐标计算像素值。
func bilinearInterpolation(img *image.RGBA, x, y float64) color.RGBA {
	x1, y1 := int(x), int(y)
	x2, y2 := x1+1, y1+1
	if x2 >= img.Bounds().Max.X {
		x2 = img.Bounds().Max.X - 1
	}
	if y2 >= img.Bounds().Max.Y {
		y2 = img.Bounds().Max.Y - 1
	}

	q11 := img.RGBAAt(x1, y1)
	q12 := img.RGBAAt(x1, y2)
	q21 := img.RGBAAt(x2, y1)
	q22 := img.RGBAAt(x2, y2)

	xfrac, yfrac := x-float64(x1), y-float64(y1)

	r := uint8((1-xfrac)*(1-yfrac)*float64(q11.R) +
		xfrac*(1-yfrac)*float64(q21.R) +
		(1-xfrac)*yfrac*float64(q12.R) +
		xfrac*yfrac*float64(q22.R))
	g := uint8((1-xfrac)*(1-yfrac)*float64(q11.G) +
		xfrac*(1-yfrac)*float64(q21.G) +
		(1-xfrac)*yfrac*float64(q12.G) +
		xfrac*yfrac*float64(q22.G))
	b := uint8((1-xfrac)*(1-yfrac)*float64(q11.B) +
		xfrac*(1-yfrac)*float64(q21.B) +
		(1-xfrac)*yfrac*float64(q12.B) +
		xfrac*yfrac*float64(q22.B))

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// 平移
func translate(img *image.RGBA, dx, dy int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newX, newY := x-dx, y-dy
			if newX >= bounds.Min.X && newX < bounds.Max.X &&
				newY >= bounds.Min.Y && newY < bounds.Max.Y {
				newImg.Set(x, y, img.At(newX, newY))
			} else {
				newImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // 背景黑色
			}
		}
	}
	return newImg
}

// 缩放
func scale(img *image.RGBA, scaleX, scaleY float64) *image.RGBA {
	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * scaleX)
	newHeight := int(float64(bounds.Dy()) * scaleY)
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX, srcY := float64(x)/scaleX, float64(y)/scaleY
			newImg.Set(x, y, bilinearInterpolation(img, srcX, srcY))
		}
	}
	return newImg
}

// 旋转
func rotate(img *image.RGBA, angle float64) *image.RGBA {
	bounds := img.Bounds()
	centerX, centerY := float64(bounds.Dx())/2, float64(bounds.Dy())/2
	cosA, sinA := math.Cos(angle), math.Sin(angle)

	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 将 (x, y) 映射回原图
			rotatedX := (float64(x)-centerX)*cosA + (float64(y)-centerY)*sinA + centerX
			rotatedY := -(float64(x)-centerX)*sinA + (float64(y)-centerY)*cosA + centerY

			if rotatedX >= 0 && rotatedX < float64(bounds.Dx()) &&
				rotatedY >= 0 && rotatedY < float64(bounds.Dy()) {
				newImg.Set(x, y, bilinearInterpolation(img, rotatedX, rotatedY))
			} else {
				newImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // 背景黑色
			}
		}
	}
	return newImg
}

// 仿射变换通过 2x3 矩阵实现
func affineTransform23(img *image.RGBA, matrix [2][3]float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 计算逆变换
			det := matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
			if det == 0 {
				continue
			}
			invM := [2][3]float64{
				{matrix[1][1] / det, -matrix[0][1] / det, (matrix[0][1]*matrix[1][2] - matrix[1][1]*matrix[0][2]) / det},
				{-matrix[1][0] / det, matrix[0][0] / det, (matrix[1][0]*matrix[0][2] - matrix[0][0]*matrix[1][2]) / det},
			}

			srcX := invM[0][0]*float64(x) + invM[0][1]*float64(y) + invM[0][2]
			srcY := invM[1][0]*float64(x) + invM[1][1]*float64(y) + invM[1][2]

			if srcX >= 0 && srcX < float64(bounds.Dx()) &&
				srcY >= 0 && srcY < float64(bounds.Dy()) {
				newImg.Set(x, y, bilinearInterpolation(img, srcX, srcY))
			} else {
				newImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // 背景黑色
			}
		}
	}
	return newImg
}

// 透视变换通过 3x3 矩阵实现
func perspectiveTransform(img *image.RGBA, matrix [3][3]float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 计算逆变换
			det := matrix[0][0]*(matrix[1][1]*matrix[2][2]-matrix[1][2]*matrix[2][1]) -
				matrix[0][1]*(matrix[1][0]*matrix[2][2]-matrix[1][2]*matrix[2][0]) +
				matrix[0][2]*(matrix[1][0]*matrix[2][1]-matrix[1][1]*matrix[2][0])
			if det == 0 {
				continue
			}
			invM := [3][3]float64{
				{(matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]) / det,
					(matrix[0][2]*matrix[2][1] - matrix[0][1]*matrix[2][2]) / det,
					(matrix[0][1]*matrix[1][2] - matrix[0][2]*matrix[1][1]) / det},
				{(matrix[1][2]*matrix[2][0] - matrix[1][0]*matrix[2][2]) / det,
					(matrix[0][0]*matrix[2][2] - matrix[0][2]*matrix[2][0]) / det,
					(matrix[0][2]*matrix[1][0] - matrix[0][0]*matrix[1][2]) / det},
				{(matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]) / det,
					(matrix[0][1]*matrix[2][0] - matrix[0][0]*matrix[2][1]) / det,
					(matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]) / det},
			}

			w := invM[2][0]*float64(x) + invM[2][1]*float64(y) + invM[2][2]
			srcX := (invM[0][0]*float64(x) + invM[0][1]*float64(y) + invM[0][2]) / w
			srcY := (invM[1][0]*float64(x) + invM[1][1]*float64(y) + invM[1][2]) / w

			if srcX >= 0 && srcX < float64(bounds.Dx()) &&
				srcY >= 0 && srcY < float64(bounds.Dy()) {
				newImg.Set(x, y, bilinearInterpolation(img, srcX, srcY))
			} else {
				newImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // 背景黑色
			}
		}
	}
	return newImg
}

func case21() {
	src := getTestImg()

	rgbaImg := image.NewRGBA(src.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), src, image.Point{}, draw.Src)

	translated := translate(rgbaImg, 50, 30)
	scaled := scale(rgbaImg, 0.5, 0.5)
	rotated := rotate(rgbaImg, math.Pi/4)

	outputFileName1 := "output_case21-1.jpg"
	outputFile1, err := os.Create(outputFileName1)
	if err != nil {
		panic(err)
	}
	defer outputFile1.Close()
	err = jpeg.Encode(outputFile1, translated, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

	outputFileName2 := "output_case21-2.jpg"
	outputFile2, err := os.Create(outputFileName2)
	if err != nil {
		panic(err)
	}
	defer outputFile2.Close()
	err = jpeg.Encode(outputFile2, scaled, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

	outputFileName3 := "output_case21-3.jpg"
	outputFile3, err := os.Create(outputFileName3)
	if err != nil {
		panic(err)
	}
	defer outputFile3.Close()
	err = jpeg.Encode(outputFile3, rotated, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// 例子效果失败

// case22 实现基于三角形区域变换的图像变形

// 定义三角形及其相关的辅助函数
type Point struct {
	X, Y float64
}

type Triangle struct {
	P1, P2, P3 Point
}

// 仿射变换矩阵计算  根据源三角形和目标三角形计算仿射变换矩阵。
func computeAffineTransform(src, dst Triangle) [2][3]float64 {
	x1, y1 := src.P1.X, src.P1.Y
	x2, y2 := src.P2.X, src.P2.Y
	x3, y3 := src.P3.X, src.P3.Y

	u1, v1 := dst.P1.X, dst.P1.Y
	u2, v2 := dst.P2.X, dst.P2.Y
	u3, v3 := dst.P3.X, dst.P3.Y

	// 计算仿射变换矩阵
	det := (y2-y3)*(x1-x3) + (x3-x2)*(y1-y3)
	if det == 0 {
		return [2][3]float64{}
	}

	a11 := ((v2-v3)*(x1-x3) + (u3-u2)*(y1-y3)) / det
	a21 := ((v3-v1)*(x1-x3) + (u1-u3)*(y1-y3)) / det
	a12 := ((y2-y3)*(u1-u3) + (x3-x2)*(v1-v3)) / det
	a22 := ((y3-y1)*(u1-u3) + (x1-x3)*(v1-v3)) / det
	b1 := u1 - a11*x1 - a12*y1
	b2 := v1 - a21*x1 - a22*y1

	return [2][3]float64{
		{a11, a12, b1},
		{a21, a22, b2},
	}
}

// 判断点是否在三角形内  使用重心坐标法判断一个点是否在三角形内。
func isPointInTriangle(p Point, t Triangle) bool {
	v0 := Point{t.P3.X - t.P1.X, t.P3.Y - t.P1.Y}
	v1 := Point{t.P2.X - t.P1.X, t.P2.Y - t.P1.Y}
	v2 := Point{p.X - t.P1.X, p.Y - t.P1.Y}

	dot00 := v0.X*v0.X + v0.Y*v0.Y
	dot01 := v0.X*v1.X + v0.Y*v1.Y
	dot02 := v0.X*v2.X + v0.Y*v2.Y
	dot11 := v1.X*v1.X + v1.Y*v1.Y
	dot12 := v1.X*v2.X + v1.Y*v2.Y

	invDenom := 1.0 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v < 1)
}

// 三角形区域变换  对每个三角形区域进行仿射变换并绘制到目标图像中
func transformTriangle(img *image.RGBA, srcTri, dstTri Triangle, newImg *image.RGBA) {
	matrix := computeAffineTransform(srcTri, dstTri)

	for y := newImg.Bounds().Min.Y; y < newImg.Bounds().Max.Y; y++ {
		for x := newImg.Bounds().Min.X; x < newImg.Bounds().Max.X; x++ {
			p := Point{float64(x), float64(y)}
			if isPointInTriangle(p, dstTri) {
				srcX := matrix[0][0]*float64(x) + matrix[0][1]*float64(y) + matrix[0][2]
				srcY := matrix[1][0]*float64(x) + matrix[1][1]*float64(y) + matrix[1][2]

				if srcX >= 0 && srcX < float64(img.Bounds().Dx()) &&
					srcY >= 0 && srcY < float64(img.Bounds().Dy()) {
					newImg.Set(x, y, bilinearInterpolation(img, srcX, srcY))
				} else {
					newImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // 背景黑色
				}
			}
		}
	}
}

func case22() {
	src := getTest2Img()

	rgbaImg := image.NewRGBA(src.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), src, image.Point{}, draw.Src)

	// 定义源三角形和目标三角形
	srcTri := Triangle{
		Point{100, 100}, Point{200, 100}, Point{150, 200},
	}
	dstTri := Triangle{
		Point{150, 100}, Point{250, 150}, Point{200, 250},
	}

	// 创建输出图像
	newImg := image.NewRGBA(rgbaImg.Bounds())

	// 应用三角形变换
	transformTriangle(rgbaImg, srcTri, dstTri, newImg)

	outputFileName := "output_case22.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case23 对图像指定区域进行马赛克处理

func case23() {
	src := getTest2Img()

	// 创建一个可绘制的图像副本
	bounds := src.Bounds()
	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, src, bounds.Min, draw.Src)

	// 指定马赛克区域和块大小
	x, y := 100, 100
	width, height := 800, 400
	blockSize := 64

	// 对指定区域进行马赛克处理
	for i := y; i < y+height; i += blockSize {
		for j := x; j < x+width; j += blockSize {
			var r, g, b, a uint32
			count := 0
			for m := 0; m < blockSize; m++ {
				for n := 0; n < blockSize; n++ {
					if i+m < y+height && j+n < x+width {
						pr, pg, pb, pa := drawImg.At(j+n, i+m).RGBA()
						r += pr
						g += pg
						b += pb
						a += pa
						count++
					}
				}
			}
			if count > 0 {
				r /= uint32(count)
				g /= uint32(count)
				b /= uint32(count)
				a /= uint32(count)
				c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
				for m := 0; m < blockSize; m++ {
					for n := 0; n < blockSize; n++ {
						if i+m < y+height && j+n < x+width {
							drawImg.Set(j+n, i+m, c)
						}
					}
				}
			}
		}
	}

	outputFileName := "output_case23.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, drawImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case24 图像浮雕效果

func case24() {
	src := getTest2Img()

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 0; y < height-1; y++ {
		for x := 0; x < width-1; x++ {
			currentPixel := src.At(x, y)
			nextPixel := src.At(x+1, y+1)

			r1, g1, b1, _ := currentPixel.RGBA()
			r2, g2, b2, _ := nextPixel.RGBA()

			r := int(r1/256) - int(r2/256) + 128
			g := int(g1/256) - int(g2/256) + 128
			b := int(b1/256) - int(b2/256) + 128

			if r < 0 {
				r = 0
			} else if r > 255 {
				r = 255
			}
			if g < 0 {
				g = 0
			} else if g > 255 {
				g = 255
			}
			if b < 0 {
				b = 0
			} else if b > 255 {
				b = 255
			}

			result.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	outputFileName := "output_case24.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case25 彩色图像的平滑处理
// 平滑滤波可以使图像模糊，从而减少图像中的细节，使图像变得柔和。

func case25() {
	src := getTest2Img()

	// 定义平滑处理的核大小
	kernelSize := 3

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	halfKernel := kernelSize / 2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var rSum, gSum, bSum, count int
			for ky := -halfKernel; ky <= halfKernel; ky++ {
				for kx := -halfKernel; kx <= halfKernel; kx++ {
					nx := x + kx
					ny := y + ky
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						r, g, b, _ := src.At(nx, ny).RGBA()
						rSum += int(r / 256)
						gSum += int(g / 256)
						bSum += int(b / 256)
						count++
					}
				}
			}
			if count > 0 {
				r := uint8(rSum / count)
				g := uint8(gSum / count)
				b := uint8(bSum / count)
				result.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}

	outputFileName := "output_case25.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case26 图像的锐化处理
// 锐化的主要目的是突出图像的边缘和细节，使图像变得清晰。

func case26() {
	src := getTest2Img()

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var rSum, gSum, bSum int
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					nx := x + kx - 1
					ny := y + ky - 1
					r, g, b, _ := src.At(nx, ny).RGBA()
					factor := laplacianKernel[ky][kx]
					rSum += int(r/256) * factor
					gSum += int(g/256) * factor
					bSum += int(b/256) * factor
				}
			}
			// 获取当前像素的原始值
			r0, g0, b0, _ := src.At(x, y).RGBA()
			r0 = r0 / 256
			g0 = g0 / 256
			b0 = b0 / 256

			// 计算锐化后的颜色值
			r := int(r0) + rSum
			g := int(g0) + gSum
			b := int(b0) + bSum

			// 确保颜色值在 0 到 255 之间
			if r < 0 {
				r = 0
			} else if r > 255 {
				r = 255
			}
			if g < 0 {
				g = 0
			} else if g > 255 {
				g = 255
			}
			if b < 0 {
				b = 0
			} else if b > 255 {
				b = 255
			}

			result.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	outputFileName := "output_case26.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// 拉普拉斯算子核
var laplacianKernel = [3][3]int{
	{0, -1, 0},
	{-1, 4, -1},
	{0, -1, 0},
}

// ========================================================================

// case27  彩色图像的分割
// 彩色图像分割有多种方法，除了之前提到的基于阈值的分割; 该例子 基于 K - Means 聚类算法的彩色图像分割

func case27() {
	src := getTest2Img()

	// 设定聚类的类别数和最大迭代次数
	k := 3
	maxIterations := 100

	// 进行图像分割
	segmentedImg := Segment(src, k, maxIterations)

	outputFileName := "output_case27.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, segmentedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// Point 表示 RGB 颜色空间中的一个点
type Point27 [3]float64

// distance 计算两个点之间的欧几里得距离
func distance(p1, p2 Point27) float64 {
	return math.Sqrt(math.Pow(p1[0]-p2[0], 2) + math.Pow(p1[1]-p2[1], 2) + math.Pow(p1[2]-p2[2], 2))
}

// kMeans 实现 K - Means 聚类算法
func kMeans(points []Point27, k int, maxIterations int) ([]int, []Point27) {
	rand.Seed(time.Now().UnixNano())
	centers := make([]Point27, k)
	for i := range centers {
		centers[i] = points[rand.Intn(len(points))]
	}

	labels := make([]int, len(points))
	for iter := 0; iter < maxIterations; iter++ {
		// 分配点到最近的中心
		for i, p := range points {
			minDist := math.MaxFloat64
			for j, c := range centers {
				dist := distance(p, c)
				if dist < minDist {
					minDist = dist
					labels[i] = j
				}
			}
		}

		// 更新中心
		newCenters := make([]Point27, k)
		counts := make([]int, k)
		for i, label := range labels {
			newCenters[label][0] += points[i][0]
			newCenters[label][1] += points[i][1]
			newCenters[label][2] += points[i][2]
			counts[label]++
		}
		for i := range newCenters {
			if counts[i] > 0 {
				newCenters[i][0] /= float64(counts[i])
				newCenters[i][1] /= float64(counts[i])
				newCenters[i][2] /= float64(counts[i])
			}
		}

		// 检查是否收敛
		converged := true
		for i := range centers {
			if distance(centers[i], newCenters[i]) > 1e-6 {
				converged = false
				break
			}
		}
		if converged {
			break
		}
		centers = newCenters
	}
	return labels, centers
}

// Segment 函数使用 K - Means 进行图像分割
func Segment(img image.Image, k int, maxIterations int) draw.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	points := make([]Point27, 0, width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			points = append(points, Point27{float64(r / 256), float64(g / 256), float64(b / 256)})
		}
	}

	labels, centers := kMeans(points, k, maxIterations)

	result := image.NewRGBA(bounds)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			label := labels[index]
			center := centers[label]
			r := uint8(center[0])
			g := uint8(center[1])
			b := uint8(center[2])
			result.Set(x, y, color.RGBA{r, g, b, 255})
			index++
		}
	}
	return result
}

// ========================================================================

// case28 彩色图像的边缘提取
// 采用 Sobel 算子进行边缘提取，Sobel 算子是一种常用的边缘检测算子，能分别计算图像在水平和垂直方向上的梯度，从而检测出图像的边缘

func case28() {
	src := getTest2Img()

	// Sobel 算子的水平和垂直核
	var sobelX = [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	var sobelY = [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var gxR, gyR, gxG, gyG, gxB, gyB int
			// 计算水平和垂直方向的梯度
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					nx := x + kx - 1
					ny := y + ky - 1
					r, g, b, _ := src.At(nx, ny).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					gxR += int(r) * sobelX[ky][kx]
					gyR += int(r) * sobelY[ky][kx]
					gxG += int(g) * sobelX[ky][kx]
					gyG += int(g) * sobelY[ky][kx]
					gxB += int(b) * sobelX[ky][kx]
					gyB += int(b) * sobelY[ky][kx]
				}
			}
			// 计算梯度幅值
			gradR := math.Sqrt(float64(gxR*gxR + gyR*gyR))
			gradG := math.Sqrt(float64(gxG*gxG + gyG*gyG))
			gradB := math.Sqrt(float64(gxB*gxB + gyB*gyB))

			// 确保梯度幅值在 0 到 255 之间
			if gradR > 255 {
				gradR = 255
			}
			if gradG > 255 {
				gradG = 255
			}
			if gradB > 255 {
				gradB = 255
			}

			result.Set(x, y, color.RGBA{uint8(gradR), uint8(gradG), uint8(gradB), 255})
		}
	}

	outputFileName := "output_case28.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case29 图像颜色反转

func case29() {
	src := getTest2Img()

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			// 将颜色值转换为 0 - 255 范围
			r = r / 256
			g = g / 256
			b = b / 256
			a = a / 256

			// 反转颜色
			r = 255 - r
			g = 255 - g
			b = 255 - b

			result.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	outputFileName := "output_case29.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case30 图像腐蚀
// 如果输入是彩色图像，可能需要先进行二值化处理

func case30() {
	//src := getTest2Img()

	src := getImg("./output_case19.jpg")

	result := case30_1(src)

	outputFileName := "output_case30.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// 图像腐蚀
func case30_1(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 3x3 结构元素
	structuringElement := [3][3]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			allForeground := true
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					if structuringElement[ky][kx] {
						nx := x + kx - 1
						ny := y + ky - 1
						r, _, _, _ := src.At(nx, ny).RGBA()
						if r/256 < 128 {
							allForeground = false
							break
						}
					}
				}
				if !allForeground {
					break
				}
			}
			if allForeground {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return result
}

// ========================================================================

// case31 图像膨胀
// 如果是彩色图像，建议先进行二值化处理再进行膨胀操作，以达到更好的效果。

func case31() {
	src := getImg("./output_case19.jpg")

	result := case31_1(src)

	outputFileName := "output_case31.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// 图像膨胀
func case31_1(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 3x3 结构元素
	structuringElement := [3][3]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			anyForeground := false
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					if structuringElement[ky][kx] {
						nx := x + kx - 1
						ny := y + ky - 1
						r, _, _, _ := src.At(nx, ny).RGBA()
						if r/256 >= 128 {
							anyForeground = true
							break
						}
					}
				}
				if anyForeground {
					break
				}
			}
			if anyForeground {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return result
}

// ========================================================================

// case32 图像的开运算
// 图像的开运算（Opening）是一种形态学操作，它是先对图像进行腐蚀操作，然后再进行膨胀操作。开运算可以去除图像中的小物体、分离物体以及平滑物体的边界

func case32() {
	src := getImg("./output_case30.jpg")

	result := case31_1(src)

	outputFileName := "output_case32.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case33 图像的闭运算
// 图像的闭运算（Closing）是一种形态学操作，它是先对图像进行膨胀操作，然后再进行腐蚀操作。闭运算常用于填充物体内的小孔、连接邻近的物体等。

func case33() {
	src := getImg("./output_case19.jpg")

	result := case31_1(src)

	result = case30_1(result)

	outputFileName := "output_case33.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case34 图像的内边界提取
// 图像的内边界提取可以通过图像腐蚀操作与原图像做差来实现
// 若输入是彩色图像，可能需要先进行二值化处理

func case34() {
	src := getImg("./output_case19.jpg")

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 先对图像进行腐蚀操作
	erodedImg := case30_1(src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := src.At(x, y).RGBA()
			r2, g2, b2, a2 := erodedImg.At(x, y).RGBA()

			r1 = r1 / 256
			g1 = g1 / 256
			b1 = b1 / 256
			a1 = a1 / 256

			r2 = r2 / 256
			g2 = g2 / 256
			b2 = b2 / 256
			a2 = a2 / 256

			// 通过原图像与腐蚀后的图像做差提取内边界
			if r1 > 128 && r2 < 128 {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	outputFileName := "output_case34.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case35 图像的外边界提取
// 图像的外边界提取可以通过对图像进行膨胀操作，然后将膨胀后的图像与原图像做差来实现
// 若输入是彩色图像，可能需要先进行二值化处理

func case35() {
	src := getImg("./output_case19.jpg")

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 先对图像进行膨胀操作
	dilatedImg := case31_1(src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := src.At(x, y).RGBA()
			r2, g2, b2, a2 := dilatedImg.At(x, y).RGBA()

			r1 = r1 / 256
			g1 = g1 / 256
			b1 = b1 / 256
			a1 = a1 / 256

			r2 = r2 / 256
			g2 = g2 / 256
			b2 = b2 / 256
			a2 = a2 / 256

			// 通过膨胀后的图像与原图像做差提取外边界
			if r2 > 128 && r1 < 128 {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	outputFileName := "output_case35.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case36 利用均值迭代阀值分割法分割图像
// 均值迭代阈值分割法是一种自动确定图像分割阈值的方法，其基本思想是通过迭代计算，不断更新阈值，直到阈值收敛
// 如果输入是彩色图像，需要先将其转换为灰度图像

func case36() {
	src := getImg("./output_case18.jpg")

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 初始化阈值
	var threshold uint8 = 128
	var newThreshold uint8
	var diff float64 = 1.0

	for diff > 0.5 {
		var sum1, sum2 int
		var count1, count2 int

		// 遍历图像像素，根据当前阈值将像素分为两类
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, _, _, _ := src.At(x, y).RGBA()
				gray := uint8(r / 256)

				if gray <= threshold {
					sum1 += int(gray)
					count1++
				} else {
					sum2 += int(gray)
					count2++
				}
			}
		}

		// 计算两类像素的平均灰度值
		var mean1, mean2 float64
		if count1 > 0 {
			mean1 = float64(sum1) / float64(count1)
		}
		if count2 > 0 {
			mean2 = float64(sum2) / float64(count2)
		}

		// 计算新的阈值
		newThreshold = uint8((mean1 + mean2) / 2)

		// 计算阈值的差值
		diff = float64(newThreshold) - float64(threshold)
		if diff < 0 {
			diff = -diff
		}

		// 更新阈值
		threshold = newThreshold
	}

	// 根据最终阈值进行图像分割
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := src.At(x, y).RGBA()
			gray := uint8(r / 256)

			if gray <= threshold {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	outputFileName := "output_case36.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case37 最大类间方差法分割法分割图像
// 最大类间方差法（Otsu 算法）是一种常用的图像阈值分割方法，它通过最大化类间方差来自动确定一个最优的阈值，将图像分为前景和背景两部分
// 如果输入是彩色图像，需要先将其转换为灰度图像

func case37() {
	src := getImg("./output_case18.jpg")

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 统计每个灰度级的像素数量
	histogram := [256]int{}
	totalPixels := width * height

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := src.At(x, y).RGBA()
			gray := uint8(r / 256)
			histogram[gray]++
		}
	}

	// 初始化最大类间方差和最优阈值
	var maxVariance float64
	var optimalThreshold uint8

	// 遍历所有可能的阈值
	for t := 0; t < 256; t++ {
		var w0, w1, u0, u1, sum0, sum1, count0, count1 int

		// 计算前景和背景的像素数量和灰度值总和
		for i := 0; i < t; i++ {
			count0 += histogram[i]
			sum0 += i * histogram[i]
		}
		for i := t; i < 256; i++ {
			count1 += histogram[i]
			sum1 += i * histogram[i]
		}

		// 避免除零错误
		if count0 == 0 || count1 == 0 {
			continue
		}

		// 计算前景和背景的概率
		w0 = count0 * 100 / totalPixels
		w1 = count1 * 100 / totalPixels

		// 计算前景和背景的平均灰度值
		u0 = sum0 / count0
		u1 = sum1 / count1

		// 计算类间方差
		variance := float64(w0*w1) * float64(u0-u1) * float64(u0-u1)

		// 更新最大类间方差和最优阈值
		if variance > maxVariance {
			maxVariance = variance
			optimalThreshold = uint8(t)
		}
	}

	// 根据最优阈值进行图像分割
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := src.At(x, y).RGBA()
			gray := uint8(r / 256)

			if gray <= optimalThreshold {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	outputFileName := "output_case37.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case38 自适应阀值分割法分割图像
// 自适应阈值分割法是一种根据图像局部区域的特性来确定每个像素的阈值，从而实现图像分割的方法
// 若输入是彩色图像，需要先将其转换为灰度图像

func case38() {
	src := getImg("./output_case18.jpg")

	// 设置局部块大小和常数 C
	blockSize := 15
	C := 5

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 遍历图像的每个像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 计算局部区域的边界
			xStart := x - blockSize/2
			if xStart < bounds.Min.X {
				xStart = bounds.Min.X
			}
			xEnd := x + blockSize/2
			if xEnd >= bounds.Max.X {
				xEnd = bounds.Max.X - 1
			}
			yStart := y - blockSize/2
			if yStart < bounds.Min.Y {
				yStart = bounds.Min.Y
			}
			yEnd := y + blockSize/2
			if yEnd >= bounds.Max.Y {
				yEnd = bounds.Max.Y - 1
			}

			// 计算局部区域的像素灰度值总和
			var sum int
			var count int
			for j := yStart; j <= yEnd; j++ {
				for i := xStart; i <= xEnd; i++ {
					r, _, _, _ := src.At(i, j).RGBA()
					gray := uint8(r / 256)
					sum += int(gray)
					count++
				}
			}

			// 计算局部区域的平均灰度值
			localMean := sum / count

			// 获取当前像素的灰度值
			r, _, _, _ := src.At(x, y).RGBA()
			gray := uint8(r / 256)

			// 根据局部平均灰度值和常数 C 确定阈值
			threshold := localMean - C

			// 根据阈值进行分割
			if int(gray) <= threshold {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	outputFileName := "output_case38.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case39 最大熵分割法分割图像
// 最大熵分割法是基于图像灰度直方图，通过最大化前景和背景的熵之和来确定最佳分割阈值的方法
// 若输入是彩色图像，需要先将其转换为灰度图像

func case39() {
	src := getImg("./output_case18.jpg")

	// 计算熵的函数
	entropy := func(histogram []float64) float64 {
		var ent float64
		for _, p := range histogram {
			if p > 0 {
				ent -= p * math.Log(p)
			}
		}
		return ent
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 统计灰度直方图
	histogram := make([]int, 256)
	totalPixels := width * height
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := src.At(x, y).RGBA()
			gray := int(r / 256)
			histogram[gray]++
		}
	}

	// 计算概率分布
	prob := make([]float64, 256)
	for i := 0; i < 256; i++ {
		prob[i] = float64(histogram[i]) / float64(totalPixels)
	}

	// 寻找最大熵对应的阈值
	var maxEntropy float64
	var optimalThreshold int
	for t := 0; t < 256; t++ {
		// 前景和背景的概率分布
		foreProb := prob[:t+1]
		backProb := prob[t+1:]

		// 计算前景和背景的熵
		foreEntropy := entropy(foreProb)
		backEntropy := entropy(backProb)

		// 计算总熵
		totalEntropy := foreEntropy + backEntropy

		if totalEntropy > maxEntropy {
			maxEntropy = totalEntropy
			optimalThreshold = t
		}
	}

	// 根据最优阈值进行图像分割
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := src.At(x, y).RGBA()
			gray := int(r / 256)
			if gray <= optimalThreshold {
				result.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				result.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	outputFileName := "output_case39.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case40 图像调整色相
// 要实现图像色相的调整，可以通过将 RGB 颜色空间转换为 HSV（Hue, Saturation, Value）颜色空间，调整色相（Hue）值后再转换回 RGB 颜色空间

func case40() {
	src := getTest2Img()

	// 定义色相调整值（可以根据需要修改）
	hueAdjustment := 44.0

	// RGBToHSV 将 RGB 颜色转换为 HSV 颜色
	RGBToHSV := func(r, g, b uint8) (float64, float64, float64) {
		rNorm := float64(r) / 255.0
		gNorm := float64(g) / 255.0
		bNorm := float64(b) / 255.0
		maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
		minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
		delta := maxVal - minVal

		var h, s, v float64
		v = maxVal

		if delta == 0 {
			h = 0
		} else {
			s = delta / maxVal
			if maxVal == rNorm {
				h = math.Mod((gNorm-bNorm)/delta, 6)
			} else if maxVal == gNorm {
				h = (bNorm-rNorm)/delta + 2
			} else {
				h = (rNorm-gNorm)/delta + 4
			}
			h *= 60
			if h < 0 {
				h += 360
			}
		}
		return h, s, v
	}

	// HSVToRGB 将 HSV 颜色转换为 RGB 颜色
	HSVToRGB := func(h, s, v float64) (uint8, uint8, uint8) {
		c := v * s
		hPrime := h / 60
		x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
		var r1, g1, b1 float64
		switch {
		case 0 <= hPrime && hPrime < 1:
			r1 = c
			g1 = x
			b1 = 0
		case 1 <= hPrime && hPrime < 2:
			r1 = x
			g1 = c
			b1 = 0
		case 2 <= hPrime && hPrime < 3:
			r1 = 0
			g1 = c
			b1 = x
		case 3 <= hPrime && hPrime < 4:
			r1 = 0
			g1 = x
			b1 = c
		case 4 <= hPrime && hPrime < 5:
			r1 = x
			g1 = 0
			b1 = c
		case 5 <= hPrime && hPrime < 6:
			r1 = c
			g1 = 0
			b1 = x
		}
		m := v - c
		r := uint8((r1 + m) * 255)
		g := uint8((g1 + m) * 255)
		b := uint8((b1 + m) * 255)
		return r, g, b
	}

	// AdjustHue 调整图像的色相

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			h, s, v := RGBToHSV(uint8(r), uint8(g), uint8(b))
			// 调整色相
			h = math.Mod(h+hueAdjustment, 360)
			if h < 0 {
				h += 360
			}
			r1, g1, b1 := HSVToRGB(float64(h), float64(s), float64(v))
			result.Set(x, y, color.RGBA{r1, g1, b1, uint8(a)})
		}
	}

	outputFileName := "output_case40.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case41 图像调整饱和度
// 饱和度调整可以通过将 RGB 颜色空间转换为 HSV 颜色空间，调整饱和度值后再转换回 RGB 颜色空间来完成

func case41() {
	src := getTest2Img()

	// 定义饱和度调整值（可以根据需要修改）
	saturationAdjustment := 0.1

	// RGBToHSV 将 RGB 颜色转换为 HSV 颜色
	RGBToHSV := func(r, g, b uint8) (float64, float64, float64) {
		rNorm := float64(r) / 255.0
		gNorm := float64(g) / 255.0
		bNorm := float64(b) / 255.0
		maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
		minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
		delta := maxVal - minVal

		var h, s, v float64
		v = maxVal

		if delta == 0 {
			h = 0
		} else {
			s = delta / maxVal
			if maxVal == rNorm {
				h = math.Mod((gNorm-bNorm)/delta, 6)
			} else if maxVal == gNorm {
				h = (bNorm-rNorm)/delta + 2
			} else {
				h = (rNorm-gNorm)/delta + 4
			}
			h *= 60
			if h < 0 {
				h += 360
			}
		}
		return h, s, v
	}

	// HSVToRGB 将 HSV 颜色转换为 RGB 颜色
	HSVToRGB := func(h, s, v float64) (uint8, uint8, uint8) {
		c := v * s
		hPrime := h / 60
		x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
		var r1, g1, b1 float64
		switch {
		case 0 <= hPrime && hPrime < 1:
			r1 = c
			g1 = x
			b1 = 0
		case 1 <= hPrime && hPrime < 2:
			r1 = x
			g1 = c
			b1 = 0
		case 2 <= hPrime && hPrime < 3:
			r1 = 0
			g1 = c
			b1 = x
		case 3 <= hPrime && hPrime < 4:
			r1 = 0
			g1 = x
			b1 = c
		case 4 <= hPrime && hPrime < 5:
			r1 = x
			g1 = 0
			b1 = c
		case 5 <= hPrime && hPrime < 6:
			r1 = c
			g1 = 0
			b1 = x
		}
		m := v - c
		r := uint8((r1 + m) * 255)
		g := uint8((g1 + m) * 255)
		b := uint8((b1 + m) * 255)
		return r, g, b
	}

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			h, s, v := RGBToHSV(uint8(r), uint8(g), uint8(b))

			// 调整饱和度
			s += saturationAdjustment
			if s < 0 {
				s = 0
			} else if s > 1 {
				s = 1
			}

			r1, g1, b1 := HSVToRGB(h, s, v)
			result.Set(x, y, color.RGBA{r1, g1, b1, uint8(a)})
		}
	}

	outputFileName := "output_case41.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case42 图像调整明暗度
// 图像明暗度调整可以通过简单地对图像中每个像素的 RGB 值进行线性调整来实现

func case42() {
	src := getTest2Img()

	// 定义明暗度调整值（可以根据需要修改）
	brightnessAdjustment := -50

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整亮度
			newR := int(r) + brightnessAdjustment
			newG := int(g) + brightnessAdjustment
			newB := int(b) + brightnessAdjustment

			// 确保颜色值在 0 到 255 范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case42.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case43 调整色彩平衡
// 色彩平衡调整主要是分别对图像中红、绿、蓝三个通道的值进行调整

func case43() {
	src := getTest2Img()

	// 定义色彩平衡调整值（可以根据需要修改）
	redAdjustment := 20
	greenAdjustment := -10
	blueAdjustment := 30

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整红色通道
			newR := int(r) + redAdjustment
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}

			// 调整绿色通道
			newG := int(g) + greenAdjustment
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}

			// 调整蓝色通道
			newB := int(b) + blueAdjustment
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case43.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case44 调整亮度
// 要实现图像亮度的调整，可通过对图像每个像素的 RGB 值进行线性变换来达成

func case44() {
	src := getTest2Img()

	// 定义亮度调整值（可按需修改）
	brightness := 1.5

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整亮度
			newR := int(float64(r) * brightness)
			newG := int(float64(g) * brightness)
			newB := int(float64(b) * brightness)

			// 确保颜色值在 0 到 255 范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case44.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case45 调整对比度
// 图像对比度调整，核心思路是通过对图像中每个像素的 RGB 值进行线性变换，将其映射到一个新的范围，从而改变图像的对比度

func case45() {

	src := getTest2Img()

	// 定义对比度调整值（可以根据需要修改）
	contrast := 50.0

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 对比度调整因子
	factor := (259 * (contrast + 255)) / (255 * (259 - contrast))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整对比度
			newR := int(factor*(float64(r)-128) + 128)
			newG := int(factor*(float64(g)-128) + 128)
			newB := int(factor*(float64(b)-128) + 128)

			// 确保颜色值在 0 到 255 范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case45.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case46 调整锐度
// 图像锐度调整，通常可以使用拉普拉斯算子进行图像锐化处理。拉普拉斯算子是一种二阶导数算子，它可以增强图像中的边缘和细节，从而达到锐化图像的效果

func case46() {

	src := getTest2Img()

	// 定义锐度调整值（可以根据需要修改）
	sharpness := 0.5

	// 拉普拉斯算子
	var laplacianKernel = [][]int{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var rSum, gSum, bSum int

			// 应用拉普拉斯算子
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					r, g, b, _ := src.At(x+kx, y+ky).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256

					kernelValue := laplacianKernel[ky+1][kx+1]
					rSum += int(r) * kernelValue
					gSum += int(g) * kernelValue
					bSum += int(b) * kernelValue
				}
			}

			// 获取当前像素的原始值
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整锐度
			newR := int(r) + int(float64(rSum)*sharpness)
			newG := int(g) + int(float64(gSum)*sharpness)
			newB := int(b) + int(float64(bSum)*sharpness)

			// 确保颜色值在 0 到 255 范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	// 处理边缘像素，简单复制原始像素值
	for y := 0; y < height; y++ {
		result.Set(0, y, src.At(0, y))
		result.Set(width-1, y, src.At(width-1, y))
	}
	for x := 0; x < width; x++ {
		result.Set(x, 0, src.At(x, 0))
		result.Set(x, height-1, src.At(x, height-1))
	}

	outputFileName := "output_case46.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case47 调整色阶
// 色阶调整是一种通过改变图像中像素的亮度分布来调整图像对比度和颜色的方法。它主要涉及到对图像的暗部、中间调和亮部进行重新映射

func case47() {

	src := getTest2Img()

	// 定义色阶调整值（可以根据需要修改）
	blackPoint := 30.0
	whitePoint := 220.0
	gamma := 1.2

	// adjustChannel 调整单个通道的色阶
	adjustChannel := func(value, blackPoint, whitePoint, gamma float64) float64 {
		// 第一步：将输入值限制在黑点和白点之间
		if value < blackPoint {
			value = 0
		} else if value > whitePoint {
			value = 255
		} else {
			// 线性映射到 0 - 255 范围
			value = (value - blackPoint) / (whitePoint - blackPoint) * 255
		}
		// 第二步：应用伽马校正
		if gamma != 1 {
			value = 255 * math.Pow(value/255, 1/gamma)
		}
		return value
	}

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整红色通道
			newR := adjustChannel(float64(r), blackPoint, whitePoint, gamma)
			// 调整绿色通道
			newG := adjustChannel(float64(g), blackPoint, whitePoint, gamma)
			// 调整蓝色通道
			newB := adjustChannel(float64(b), blackPoint, whitePoint, gamma)

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case47.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case48 调整曝光度
// 图像曝光度调整，可通过对图像每个像素的 RGB 值进行指数变换来实现

func case48() {

	src := getTest2Img()

	// 定义曝光度调整值（可按需修改）
	exposure := 0.2

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整曝光度
			newR := int(math.Min(255, math.Max(0, float64(r)*math.Pow(2, exposure))))
			newG := int(math.Min(255, math.Max(0, float64(g)*math.Pow(2, exposure))))
			newB := int(math.Min(255, math.Max(0, float64(b)*math.Pow(2, exposure))))

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case48.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case49 调整色温
// 图像色温调整，核心思路是通过调整图像中 RGB 颜色通道的比例来模拟不同的色温效果。不同的色温会让图像
// 呈现出偏暖（如红色、黄色调）或偏冷（如蓝色调）的视觉感受

func case49() {

	src := getTest2Img()

	// 定义色温调整值（可以根据需要修改）
	temperature := 5000.0

	// calculateColorGains 根据色温计算 RGB 增益
	calculateColorGains := func(temperature float64) (float64, float64, float64) {
		temperature = math.Max(1000, math.Min(40000, temperature)) / 100
		var r, g, b float64

		// 计算红色增益
		if temperature <= 66 {
			r = 255
		} else {
			r = temperature - 60
			r = 329.698727446 * math.Pow(r, -0.1332047592)
			if r < 0 {
				r = 0
			}
			if r > 255 {
				r = 255
			}
		}

		// 计算绿色增益
		if temperature <= 66 {
			g = temperature
			g = 99.4708025861*math.Log(g) - 161.1195681661
			if g < 0 {
				g = 0
			}
			if g > 255 {
				g = 255
			}
		} else {
			g = temperature - 60
			g = 288.1221695283 * math.Pow(g, -0.0755148492)
			if g < 0 {
				g = 0
			}
			if g > 255 {
				g = 255
			}
		}

		// 计算蓝色增益
		if temperature >= 66 {
			b = 255
		} else {
			if temperature <= 19 {
				b = 0
			} else {
				b = temperature - 10
				b = 138.5177312231*math.Log(b) - 305.0447927307
				if b < 0 {
					b = 0
				}
				if b > 255 {
					b = 255
				}
			}
		}

		// 归一化增益
		maxValue := math.Max(r, math.Max(g, b))
		rGain := r / maxValue
		gGain := g / maxValue
		bGain := b / maxValue

		return rGain, gGain, bGain
	}

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	// 根据色温计算 RGB 增益
	rGain, gGain, bGain := calculateColorGains(temperature)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 应用 RGB 增益
			newR := int(math.Min(255, math.Max(0, float64(r)*rGain)))
			newG := int(math.Min(255, math.Max(0, float64(g)*gGain)))
			newB := int(math.Min(255, math.Max(0, float64(b)*bGain)))

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	outputFileName := "output_case49.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// ========================================================================

// case50 调整色调
// 图像色调调整，通常会先把 RGB 颜色空间转换为 HSV（Hue, Saturation, Value）颜色空间，接着调整色相（Hue）值，最后再转换回 RGB 颜色空间

func case50() {

	src := getTest2Img()

	// 定义色调调整值（可以根据需要修改）
	hueAdjustment := 90.0

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			h, s, v := RGBToHSV(uint8(r), uint8(g), uint8(b))

			// 调整色调
			h = math.Mod(h+hueAdjustment, 360)
			if h < 0 {
				h += 360
			}

			r1, g1, b1 := HSVToRGB(h, s, v)
			result.Set(x, y, color.RGBA{r1, g1, b1, uint8(a)})
		}
	}

	outputFileName := "output_case50.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// RGBToHSV 将 RGB 颜色转换为 HSV 颜色
func RGBToHSV(r, g, b uint8) (float64, float64, float64) {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0
	maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
	minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := maxVal - minVal

	var h, s, v float64
	v = maxVal

	if delta == 0 {
		h = 0
	} else {
		s = delta / maxVal
		if maxVal == rNorm {
			h = math.Mod((gNorm-bNorm)/delta, 6)
		} else if maxVal == gNorm {
			h = (bNorm-rNorm)/delta + 2
		} else {
			h = (rNorm-gNorm)/delta + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, v
}

// HSVToRGB 将 HSV 颜色转换为 RGB 颜色
func HSVToRGB(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	hPrime := h / 60
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
	var r1, g1, b1 float64
	switch {
	case 0 <= hPrime && hPrime < 1:
		r1 = c
		g1 = x
		b1 = 0
	case 1 <= hPrime && hPrime < 2:
		r1 = x
		g1 = c
		b1 = 0
	case 2 <= hPrime && hPrime < 3:
		r1 = 0
		g1 = c
		b1 = x
	case 3 <= hPrime && hPrime < 4:
		r1 = 0
		g1 = x
		b1 = c
	case 4 <= hPrime && hPrime < 5:
		r1 = x
		g1 = 0
		b1 = c
	case 5 <= hPrime && hPrime < 6:
		r1 = c
		g1 = 0
		b1 = x
	}
	m := v - c
	r := uint8((r1 + m) * 255)
	g := uint8((g1 + m) * 255)
	b := uint8((b1 + m) * 255)
	return r, g, b
}

// ========================================================================

// case51 图像降噪
// 图像降噪是图像处理中常见的操作，其中高斯模糊是一种简单且常用的降噪方法。通过调整高斯核的大小（即标准差），可以控制降噪的程度

func case51() {

	src := getTest2Img()

	// 定义降噪程度（可以根据需要修改）
	sigma := 2.0

	// 生成一维高斯核
	generateGaussianKernel := func(sigma float64) []float64 {
		size := int(math.Ceil(sigma * 3))
		kernel := make([]float64, 2*size+1)
		sum := 0.0
		for i := -size; i <= size; i++ {
			kernel[i+size] = math.Exp(-float64(i*i) / (2 * sigma * sigma))
			sum += kernel[i+size]
		}
		// 归一化
		for i := range kernel {
			kernel[i] /= sum
		}
		return kernel
	}

	// 一维高斯模糊
	gaussianBlur1D := func(img image.Image, kernel []float64) draw.Image {
		bounds := img.Bounds()
		//width := bounds.Dx()
		//height := bounds.Dy()
		result := image.NewRGBA(bounds)

		// 水平方向模糊
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var rSum, gSum, bSum, aSum float64
				kernelSize := len(kernel)
				halfKernelSize := kernelSize / 2
				for i := -halfKernelSize; i <= halfKernelSize; i++ {
					newX := x + i
					if newX < bounds.Min.X {
						newX = bounds.Min.X
					} else if newX >= bounds.Max.X {
						newX = bounds.Max.X - 1
					}
					r, g, b, a := img.At(newX, y).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					rSum += float64(r) * kernel[i+halfKernelSize]
					gSum += float64(g) * kernel[i+halfKernelSize]
					bSum += float64(b) * kernel[i+halfKernelSize]
					aSum += float64(a) * kernel[i+halfKernelSize]
				}
				result.Set(x, y, color.RGBA{
					uint8(rSum),
					uint8(gSum),
					uint8(bSum),
					uint8(aSum),
				})
			}
		}

		// 垂直方向模糊
		temp := image.NewRGBA(bounds)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				var rSum, gSum, bSum, aSum float64
				kernelSize := len(kernel)
				halfKernelSize := kernelSize / 2
				for i := -halfKernelSize; i <= halfKernelSize; i++ {
					newY := y + i
					if newY < bounds.Min.Y {
						newY = bounds.Min.Y
					} else if newY >= bounds.Max.Y {
						newY = bounds.Max.Y - 1
					}
					r, g, b, a := result.At(x, newY).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					rSum += float64(r) * kernel[i+halfKernelSize]
					gSum += float64(g) * kernel[i+halfKernelSize]
					bSum += float64(b) * kernel[i+halfKernelSize]
					aSum += float64(a) * kernel[i+halfKernelSize]
				}
				temp.Set(x, y, color.RGBA{
					uint8(rSum),
					uint8(gSum),
					uint8(bSum),
					uint8(aSum),
				})
			}
		}
		return temp
	}

	// 图像降噪
	denoiseImage := func(img image.Image, sigma float64) draw.Image {
		kernel := generateGaussianKernel(sigma)
		return gaussianBlur1D(img, kernel)
	}

	result := denoiseImage(src, sigma)

	outputFileName := "output_case51.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case52 图像模糊
// 图像模糊，常见的方法是使用高斯模糊，它通过对图像的每个像素周围的邻域进行加权平均来实现模糊效果，并且可以通过调整高斯核的标准差来控制模糊程度

func case52() {

	src := getTest2Img()

	// 定义模糊程度（可以根据需要修改）
	sigma := 4.0

	// 生成一维高斯核
	generateGaussianKernel := func(sigma float64) []float64 {
		size := int(math.Ceil(sigma * 3))
		kernel := make([]float64, 2*size+1)
		sum := 0.0
		for i := -size; i <= size; i++ {
			kernel[i+size] = math.Exp(-float64(i*i) / (2 * sigma * sigma))
			sum += kernel[i+size]
		}
		// 归一化
		for i := range kernel {
			kernel[i] /= sum
		}
		return kernel
	}

	// 一维高斯模糊
	gaussianBlur1D := func(img image.Image, kernel []float64) draw.Image {
		bounds := img.Bounds()
		//width := bounds.Dx()
		//height := bounds.Dy()
		result := image.NewRGBA(bounds)

		// 水平方向模糊
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var rSum, gSum, bSum, aSum float64
				kernelSize := len(kernel)
				halfKernelSize := kernelSize / 2
				for i := -halfKernelSize; i <= halfKernelSize; i++ {
					newX := x + i
					if newX < bounds.Min.X {
						newX = bounds.Min.X
					} else if newX >= bounds.Max.X {
						newX = bounds.Max.X - 1
					}
					r, g, b, a := img.At(newX, y).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					rSum += float64(r) * kernel[i+halfKernelSize]
					gSum += float64(g) * kernel[i+halfKernelSize]
					bSum += float64(b) * kernel[i+halfKernelSize]
					aSum += float64(a) * kernel[i+halfKernelSize]
				}
				result.Set(x, y, color.RGBA{
					uint8(rSum),
					uint8(gSum),
					uint8(bSum),
					uint8(aSum),
				})
			}
		}

		// 垂直方向模糊
		temp := image.NewRGBA(bounds)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				var rSum, gSum, bSum, aSum float64
				kernelSize := len(kernel)
				halfKernelSize := kernelSize / 2
				for i := -halfKernelSize; i <= halfKernelSize; i++ {
					newY := y + i
					if newY < bounds.Min.Y {
						newY = bounds.Min.Y
					} else if newY >= bounds.Max.Y {
						newY = bounds.Max.Y - 1
					}
					r, g, b, a := result.At(x, newY).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					rSum += float64(r) * kernel[i+halfKernelSize]
					gSum += float64(g) * kernel[i+halfKernelSize]
					bSum += float64(b) * kernel[i+halfKernelSize]
					aSum += float64(a) * kernel[i+halfKernelSize]
				}
				temp.Set(x, y, color.RGBA{
					uint8(rSum),
					uint8(gSum),
					uint8(bSum),
					uint8(aSum),
				})
			}
		}
		return temp
	}

	// 图像模糊
	blurImage := func(img image.Image, sigma float64) draw.Image {
		kernel := generateGaussianKernel(sigma)
		return gaussianBlur1D(img, kernel)
	}

	// 模糊处理
	result := blurImage(src, sigma)

	outputFileName := "output_case52.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case53 抠图实现
// 实现抠图，有多种方法，这里介绍一种基于简单的颜色阈值的抠图方法，它适用于背景颜色较为单一的图像

func case53() {

	src := getImg("./output_case19.jpg")

	// 定义目标背景颜色和颜色容差
	targetColor := color.RGBA{255, 255, 255, 255} // 白色背景
	tolerance := 30

	bounds := src.Bounds()
	//width := bounds.Dx()
	//height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 计算颜色差值
			dr := int(r) - int(targetColor.R)
			dg := int(g) - int(targetColor.G)
			db := int(b) - int(targetColor.B)
			diff := dr*dr + dg*dg + db*db

			if diff <= tolerance*tolerance {
				// 如果颜色差值在阈值范围内，设置为透明
				result.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				// 否则保留原像素颜色
				result.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}
	}

	outputFileName := "output_case53.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, result, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case54 调整图像透明度
// 若要调整图像的透明度，可以遍历图像的每个像素，然后对其 alpha 通道的值进行修改
// 若输入图像是 JPEG 格式，它本身不支持透明度，即便代码调整了 alpha 通道，保存成 JPEG 时也会丢失透明信息。所以要保证输入图像是支持透明度的格式，如 PNG。
// 同样，若输出格式为 JPEG，也会丢失透明信息。要把输出格式设为支持透明度的 PNG。

func case54() {

	src := getImg("./test3.png")

	// 定义透明度调整因子（可按需修改）
	alphaFactor := 0.4

	bounds := src.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			a = a / 256

			// 调整 alpha 值
			newA := uint8(float64(a) * alphaFactor)
			result.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), newA})
		}
	}

	outputFileName := "output_case54.png"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = png.Encode(outputFile, result)
	if err != nil {
		panic(err)
	}

}

// ========================================================================

// case55 生成纯色图

func case55() {
	// 定义图像的宽度和高度
	width := 800
	height := 600

	// 定义纯色图的颜色，这里使用红色
	solidColor := color.RGBA{255, 0, 0, 255}

	// 生成纯色图
	// 创建一个新的 RGBA 图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 遍历图像的每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 将每个像素设置为指定的颜色
			img.Set(x, y, solidColor)
		}
	}
	// 创建输出文件
	outputFile, err := os.Create("output_case55.png")
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	// 将生成的图像保存为 PNG 格式
	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatalf("无法保存图像: %v", err)
	}

}

// ========================================================================

// case56 计算两幅图像的余弦相似度
// 此方法仅适用于尺寸相同的图像。若图像尺寸不同，需先对图像进行缩放处理
// 基于图像的灰度信息，没有考虑颜色和纹理等更复杂的特征

/*
使用场景

- 图像检索与识别
相似图像搜索：在图像数据库中，根据用户提供的查询图像，通过计算余弦相似度来找到与之相似的图像。例如，在百度图片、谷歌图片等搜索引擎中，
用户上传一张图片，搜索引擎会返回与之相似的图片结果，帮助用户快速找到所需图像。
图像分类与识别：在图像识别系统中，将待识别图像与预定义的各类别图像特征向量进行余弦相似度计算，以确定该图像属于哪个类别。如在人脸识别
系统中，通过计算输入人脸图像与数据库中已知人脸图像的余弦相似度，来判断是否为同一人，进而实现身份识别

- 图像质量评估
判断图像失真程度：通过计算原始图像与处理后图像的余弦相似度，可以评估图像在压缩、滤波、降噪等处理过程中的失真程度。余弦相似度越高，
说明处理后的图像与原始图像越相似，图像质量损失越小。例如，在评估 JPEG 图像压缩算法对图像质量的影响时，可计算压缩前后图像的余弦相似度来衡量压缩效果。
比较不同图像处理算法效果：对于同一张图像，使用不同的图像处理算法（如不同的去噪算法）进行处理后，通过计算处理后图像与原始图像的余弦
相似度，比较不同算法对图像质量的保持能力，从而选择出最适合的算法。

*/

func case56() {

	src1 := getImg("./output_case4.jpg")
	src2 := getImg("./output_case2.jpg")

	// 将图像转换为灰度向量
	vec1 := ImageToGrayVector(src1)
	vec2 := ImageToGrayVector(src2)

	// 计算余弦相似度
	similarity := CosineSimilarity(vec1, vec2)
	println("两幅图像的余弦相似度为:", similarity)

}

// ImageToGrayVector 将图像转换为灰度向量
func ImageToGrayVector(img image.Image) []float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	vector := make([]float64, width*height)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			vector[index] = float64(gray.Y)
			index++
		}
	}
	return vector
}

// CosineSimilarity 计算两个向量的余弦相似度
func CosineSimilarity(vec1, vec2 []float64) float64 {
	dotProduct := 0.0
	normVec1 := 0.0
	normVec2 := 0.0
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		normVec1 += vec1[i] * vec1[i]
		normVec2 += vec2[i] * vec2[i]
	}
	if normVec1 == 0 || normVec2 == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normVec1) * math.Sqrt(normVec2))
}

// ========================================================================

// 图像相似度计算有哪些方法

/*
- 欧氏距离  原理：将图像表示为向量，通常是将图像的像素值按一定顺序排列成一维向量。然后计算两个向量在欧氏空间中的距离，距离越小，图像越相似。

- 余弦相似度  原理：将图像向量归一化后，计算两个向量的夹角余弦值。余弦值越接近 1，说明两个向量的方向越相似，即图像越相似。

- 结构相似性指数（SSIM）   原理：从图像的亮度、对比度和结构三个方面来衡量图像的相似性。它通过比较图像的局部块的均值、方差和协方差来评
估结构信息的相似程度。

- 直方图相似度   原理：统计图像的颜色直方图，即图像中不同颜色值出现的频率。通过比较两张图像的颜色直方图的相似度来衡量图像的相似性。常
用的直方图相似度计算方法有卡方距离、巴氏距离等。

- 基于特征点匹配的方法  原理：首先检测图像中的特征点，如 SIFT（尺度不变特征变换）、SURF（加速稳健特征）等特征点。然后通过匹配这些特征点来
计算图像的相似度。通常根据匹配的特征点数量、匹配的准确性等指标来衡量图像的相似程度。

- 深度学习方法  原理：利用深度神经网络，如卷积神经网络（CNN），对图像进行特征提取和表示。将图像输入到预训练的网络中，得到图像的特征向量，然
后通过计算特征向量之间的距离或相似度来衡量图像的相似性。常见的方法有基于 Siamese 网络、Triplet 网络等的图像相似度计算模型。

*/

// ========================================================================

// 使用感知哈希算法获取图片的“指纹”字符串

/*
感知哈希算法是一类用于生成图像 “指纹”（即哈希值）的算法，这些 “指纹” 字符串能够以一种紧凑的形式概括图像的感知特征。通过比较这些 “指纹”，
可以快速判断两幅图像在视觉上是否相似，而无需像传统方法那样对整个图像像素进行逐点比较
这些算法常用于图像去重、图像检索等领域

均值哈希算法（Average Hash，AHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 8x8 像素），忽略图像的细节和纵横比，这样可以确保所有图像在相同的尺度下进行比较。
灰度化：把缩放后的彩色图像转换为灰度图像，简化处理过程。
计算均值：计算灰度图像中所有像素的平均值。
生成哈希值：将每个像素的灰度值与平均值进行比较，如果像素的灰度值大于平均值，则该位置的哈希位设为 1，否则设为 0。最终得到一个由 0 和
1 组成的二进制字符串，这就是图像的 “指纹”。

感知哈希算法（Perceptual Hash，PHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 32x32 像素），以减少计算量。
灰度化：将彩色图像转换为灰度图像。
离散余弦变换（DCT）：对灰度图像进行 DCT 变换，将图像从空间域转换到频率域。DCT 变换可以将图像的能量集中在低频部分，而高频部分则包含图像的细节信息。
取低频系数：选取 DCT 变换结果的左上角 8x8 区域，这些系数代表了图像的低频特征。
计算均值：计算选取的 8x8 区域的系数平均值。
生成哈希值：将每个系数与平均值进行比较，如果系数大于平均值，则该位置的哈希位设为 1，否则设为 0。最终得到一个 64 位的二进制字符串作为图像的 “指纹”。


差异哈希算法（Difference Hash，DHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 9x8 像素）。
灰度化：将彩色图像转换为灰度图像。
计算差异：比较相邻像素的灰度值，如果右边像素的灰度值大于左边像素的灰度值，则该位置的哈希位设为 1，否则设为 0。这样可以得到一个 64 位的
二进制字符串作为图像的 “指纹”。

*/

// ========================================================================

// case57 均值哈希算法获取图像的指纹字符串

func case57() {

	src1 := getImg("./output_case4.jpg")
	// 计算均值哈希值
	hash1 := AverageHash(src1)
	fmt.Println("图像的均值哈希指纹字符串:", hash1)

	src2 := getImg("./output_case5.jpg")
	// 计算均值哈希值
	hash2 := AverageHash(src2)
	fmt.Println("图像的均值哈希指纹字符串:", hash2)

	src3 := getImg("./output_case3.jpg")
	// 计算均值哈希值
	hash3 := AverageHash(src3)
	fmt.Println("图像的均值哈希指纹字符串:", hash3)

	src4 := getImg("./output_case7.jpg")
	// 计算均值哈希值
	hash4 := AverageHash(src4)
	fmt.Println("图像的均值哈希指纹字符串:", hash4)
}

// AverageHash 计算图像的均值哈希值
func AverageHash(img image.Image) string {
	// 缩放图像到 8x8
	resized := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Src, nil)

	// 灰度化并计算像素总和
	var total int
	grayPixels := make([]uint8, 64)
	index := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// 将颜色值转换为 0 - 255 范围
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := uint8((r + g + b) / 3)
			grayPixels[index] = gray
			total += int(gray)
			index++
		}
	}

	// 计算平均灰度值
	average := total / 64

	// 生成哈希值
	var hashStr string
	for _, pixel := range grayPixels {
		if int(pixel) > average {
			hashStr += "1"
		} else {
			hashStr += "0"
		}
	}

	return hashStr
}

// ========================================================================

// case58 感知哈希算法获取图像的指纹字符串

func case58() {
	src1 := getImg("./output_case4.jpg")
	// 计算均值哈希值
	hash1 := pHash(src1)
	fmt.Println("图像的感知哈希指纹字符串:", hash1)

	src2 := getImg("./output_case5.jpg")
	// 计算均值哈希值
	hash2 := pHash(src2)
	fmt.Println("图像的感知哈希指纹字符串:", hash2)

	src3 := getImg("./output_case3.jpg")
	// 计算均值哈希值
	hash3 := pHash(src3)
	fmt.Println("图像的感知哈希指纹字符串:", hash3)

	src4 := getImg("./output_case7.jpg")
	// 计算均值哈希值
	hash4 := pHash(src4)
	fmt.Println("图像的感知哈希指纹字符串:", hash4)
}

// 二维离散余弦变换
func dct2d(data [][]float64) [][]float64 {
	N := len(data)
	result := make([][]float64, N)
	for i := range result {
		result[i] = make([]float64, N)
	}

	for u := 0; u < N; u++ {
		for v := 0; v < N; v++ {
			var sum float64
			Cu := 1.0
			Cv := 1.0
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					if u == 0 {
						Cu = 1.0 / math.Sqrt(2)
					}
					if v == 0 {
						Cv = 1.0 / math.Sqrt(2)
					}
					sum += data[x][y] * math.Cos((2*float64(x)+1)*float64(u)*math.Pi/(2*float64(N))) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/(2*float64(N)))
				}
			}
			result[u][v] = 2.0 / float64(N) * Cu * Cv * sum
		}
	}
	return result
}

// 感知哈希算法
func pHash(img image.Image) string {
	// 缩放图像到 32x32
	resized := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Src, nil)

	// 灰度化
	grayPixels := make([][]float64, 32)
	for i := range grayPixels {
		grayPixels[i] = make([]float64, 32)
	}
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := float64(r+g+b) / 3
			grayPixels[y][x] = gray
		}
	}

	// 二维离散余弦变换
	dctResult := dct2d(grayPixels)

	// 取左上角 8x8 低频分量
	lowFreq := make([][]float64, 8)
	for i := range lowFreq {
		lowFreq[i] = make([]float64, 8)
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			lowFreq[y][x] = dctResult[y][x]
		}
	}

	// 计算低频分量的平均值
	var sum float64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sum += lowFreq[y][x]
		}
	}
	average := sum / 64

	// 生成哈希值
	var hashStr string
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if lowFreq[y][x] > average {
				hashStr += "1"
			} else {
				hashStr += "0"
			}
		}
	}

	return hashStr
}

// ========================================================================

// case59 差异哈希算法获取图像的指纹字符串

func case59() {
	src1 := getImg("./output_case4.jpg")
	// 计算均值哈希值
	hash1 := DifferenceHash(src1)
	fmt.Println("图像的差异哈希指纹字符串:", hash1)

	src2 := getImg("./output_case5.jpg")
	// 计算均值哈希值
	hash2 := DifferenceHash(src2)
	fmt.Println("图像的差异哈希指纹字符串:", hash2)

	src3 := getImg("./output_case3.jpg")
	// 计算均值哈希值
	hash3 := DifferenceHash(src3)
	fmt.Println("图像的差异哈希指纹字符串:", hash3)

	src4 := getImg("./output_case7.jpg")
	// 计算均值哈希值
	hash4 := DifferenceHash(src4)
	fmt.Println("图像的差异哈希指纹字符串:", hash4)
}

// DifferenceHash 计算图像的差异哈希值
func DifferenceHash(img image.Image) string {
	// 缩放图像到 9x8
	resized := image.NewRGBA(image.Rect(0, 0, 9, 8))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Src, nil)

	// 灰度化
	grayPixels := make([]uint8, 9*8)
	index := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 9; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := uint8((r + g + b) / 3)
			grayPixels[index] = gray
			index++
		}
	}

	// 生成哈希值
	var hashStr string
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			currentIndex := y*9 + x
			nextIndex := currentIndex + 1
			if grayPixels[nextIndex] > grayPixels[currentIndex] {
				hashStr += "1"
			} else {
				hashStr += "0"
			}
		}
	}

	return hashStr
}

// ========================================================================

// case60 图像的高斯模糊
// 高斯模糊是一种常见的图像处理技术，用于减少图像中的噪声和细节，使图像变得模糊平滑

func case60() {
	src := getTest2Img()

	// 生成高斯核
	kernelSize := 44
	sigma := 444.0
	kernel := generateGaussianKernel(kernelSize, sigma)

	// 进行高斯模糊
	blurred := gaussianBlur(src, kernel)

	outputFileName := "output_case60.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, blurred, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}

}

// 生成高斯核
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}
	center := size / 2
	sum := 0.0
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			dx := float64(x - center)
			dy := float64(y - center)
			kernel[x][y] = math.Exp(-(dx*dx+dy*dy)/(2*sigma*sigma)) / (2 * math.Pi * sigma * sigma)
			sum += kernel[x][y]
		}
	}
	// 归一化
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			kernel[x][y] /= sum
		}
	}
	return kernel
}

// 图像高斯模糊
func gaussianBlur(src image.Image, kernel [][]float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum float64
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					nx := x + kx - halfKernel
					ny := y + ky - halfKernel
					if nx < bounds.Min.X {
						nx = bounds.Min.X
					} else if nx >= bounds.Max.X {
						nx = bounds.Max.X - 1
					}
					if ny < bounds.Min.Y {
						ny = bounds.Min.Y
					} else if ny >= bounds.Max.Y {
						ny = bounds.Max.Y - 1
					}
					r, g, b, a := src.At(nx, ny).RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8
					rSum += float64(r) * kernel[kx][ky]
					gSum += float64(g) * kernel[kx][ky]
					bSum += float64(b) * kernel[kx][ky]
					aSum += float64(a) * kernel[kx][ky]
				}
			}
			dst.Set(x, y, color.RGBA{
				uint8(rSum),
				uint8(gSum),
				uint8(bSum),
				uint8(aSum),
			})
		}
	}
	return dst
}

// ========================================================================

// ========================================================================

// todo gif相关  图片格式转换
// todo 加水印，文字水印，图片水印  多种效果
