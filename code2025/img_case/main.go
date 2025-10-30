package main

import (
	"container/heap"
	"container/list"
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

	// case41()

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

	// case61()

	// case62()

	// case63()

	// case64()

	// case65()

	// case66()

	// case68()

	// case69()

	// case70()

	// case71()

	// case72()

	// case73()

	// case74()

	// case75()

	// case76()

	// case78()

	// case79()

	// case80()

	// case81()

	// case82()

	case83()
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

// case61 计算图片的RGB分量矩阵
// RGB 分量矩阵本质是通过调整红、绿、蓝三通道的权重比例，实现对色彩的精准控制。它是一种线性变换工具，
//核心是用 3x3 矩阵定义每个输出通道（R/G/B）由输入三通道按何种比例混合而成。

// 实用价值
//图像分析：通过协方差矩阵判断通道相关性（如风景图中 G 和 B 通道通常正相关，因绿色植物和蓝色天空常共存）。
//色彩校正：根据通道方差调整对比度（方差小的通道需增强反差）。
//特征提取：作为图像识别的预处理特征，反映色彩分布规律。

func case61() {
	src := getTest2Img()
	// 提取RGB分量
	rgbValues := ExtractRGB(src)
	fmt.Printf("图片尺寸: %dx%d, 总像素: %d\n",
		src.Bounds().Max.X, src.Bounds().Max.Y, len(rgbValues))

	// 计算协方差矩阵
	covMatrix := ComputeCovarianceMatrix(rgbValues)

	// 打印结果（保留2位小数）
	fmt.Println("RGB分量协方差矩阵:")
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Printf("%.2f\t", covMatrix[i][j])
		}
		fmt.Println()
	}
}

// RGBValue 存储单个像素的RGB值（0-255）
type RGBValue struct {
	R, G, B uint8
}

// ExtractRGB 提取图片中所有像素的RGB值
func ExtractRGB(img image.Image) []RGBValue {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgbValues := make([]RGBValue, 0, width*height)

	// 遍历每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素颜色（RGBA格式，含Alpha通道）
			r, g, b, _ := img.At(x, y).RGBA()

			// 转换为0-255范围（RGBA返回的是0-65535，右移8位等价于除以256）
			rgb := RGBValue{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			}
			rgbValues = append(rgbValues, rgb)
		}
	}
	return rgbValues
}

// ComputeCovarianceMatrix 计算RGB三通道的协方差矩阵（3x3）
// 矩阵含义：Cov[i][j]表示第i通道与第j通道的协方差（i,j=0:R,1:G,2:B）
func ComputeCovarianceMatrix(rgbValues []RGBValue) [3][3]float64 {
	n := len(rgbValues)
	if n == 0 {
		return [3][3]float64{}
	}

	// 1. 计算三通道的均值
	var meanR, meanG, meanB float64
	for _, rgb := range rgbValues {
		meanR += float64(rgb.R)
		meanG += float64(rgb.G)
		meanB += float64(rgb.B)
	}
	meanR /= float64(n)
	meanG /= float64(n)
	meanB /= float64(n)

	// 2. 计算协方差矩阵元素
	var cov [3][3]float64
	for _, rgb := range rgbValues {
		r := float64(rgb.R) - meanR
		g := float64(rgb.G) - meanG
		b := float64(rgb.B) - meanB

		cov[0][0] += r * r // Cov(R,R)
		cov[0][1] += r * g // Cov(R,G)
		cov[0][2] += r * b // Cov(R,B)
		cov[1][1] += g * g // Cov(G,G)
		cov[1][2] += g * b // Cov(G,B)
		cov[2][2] += b * b // Cov(B,B)
	}

	// 协方差公式：E[(X-μX)(Y-μY)]，除以n-1（无偏估计）
	divisor := float64(n - 1)
	if divisor == 0 {
		divisor = 1 // 避免除以0（单像素情况）
	}
	cov[0][0] /= divisor
	cov[0][1] /= divisor
	cov[0][2] /= divisor
	cov[1][0] = cov[0][1] // 协方差矩阵对称
	cov[1][1] /= divisor
	cov[1][2] /= divisor
	cov[2][0] = cov[0][2]
	cov[2][1] = cov[1][2]
	cov[2][2] /= divisor

	return cov
}

// ========================================================================

// case62 计算图片的RGB分量矩阵 - 打印矩阵，终端输出二维数值

func case62() {
	src := getTestImg()
	// 提取RGB三通道矩阵
	matrices, err := ExtractRGBMatrices(src)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 打印矩阵
	PrintMatrix("R通道", matrices.R)
	PrintMatrix("G通道", matrices.G)
	PrintMatrix("B通道", matrices.B)
}

// 配置：终端显示的矩阵尺寸（可根据终端宽度调整）
const (
	targetWidth  = 20 // 矩阵宽度（列数）
	targetHeight = 20 // 矩阵高度（行数）
)

// RGBMatrices 存储R、G、B三通道的二维矩阵
type RGBMatrices struct {
	R [][]uint8 // R通道矩阵：[行][列] = 像素值
	G [][]uint8 // G通道矩阵
	B [][]uint8 // B通道矩阵
}

// ExtractRGBMatrices 提取缩放后的RGB三通道矩阵
func ExtractRGBMatrices(img image.Image) (RGBMatrices, error) {
	// 获取原图尺寸
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y
	if origWidth == 0 || origHeight == 0 {
		return RGBMatrices{}, fmt.Errorf("图片尺寸无效")
	}

	// 初始化三通道矩阵（目标尺寸）
	rMat := make([][]uint8, targetHeight)
	gMat := make([][]uint8, targetHeight)
	bMat := make([][]uint8, targetHeight)
	for i := range rMat {
		rMat[i] = make([]uint8, targetWidth)
		gMat[i] = make([]uint8, targetWidth)
		bMat[i] = make([]uint8, targetWidth)
	}

	// 计算缩放步长（原图到目标尺寸的采样间隔）
	stepX := float64(origWidth) / float64(targetWidth)
	stepY := float64(origHeight) / float64(targetHeight)

	// 遍历目标矩阵的每个位置，采样原图像素
	for ty := 0; ty < targetHeight; ty++ {
		for tx := 0; tx < targetWidth; tx++ {
			// 计算原图中对应的采样坐标（取整）
			origX := int(float64(tx)*stepX) + bounds.Min.X
			origY := int(float64(ty)*stepY) + bounds.Min.Y

			// 确保坐标在原图范围内（避免越界）
			if origX >= bounds.Max.X {
				origX = bounds.Max.X - 1
			}
			if origY >= bounds.Max.Y {
				origY = bounds.Max.Y - 1
			}

			// 提取像素的RGBA值（转换为0-255）
			r, g, b, _ := img.At(origX, origY).RGBA()
			r8 := uint8(r >> 8) // 16位转8位
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 存入矩阵
			rMat[ty][tx] = r8
			gMat[ty][tx] = g8
			bMat[ty][tx] = b8
		}
	}

	return RGBMatrices{R: rMat, G: gMat, B: bMat}, nil
}

// PrintMatrix 在终端打印二维矩阵
func PrintMatrix(name string, matrix [][]uint8) {
	fmt.Printf("\n===== %s 矩阵（%dx%d） =====\n", name, len(matrix), len(matrix[0]))
	for _, row := range matrix {
		for _, val := range row {
			fmt.Printf("%3d ", val) // 占3位宽度，右对齐，确保列对齐
		}
		fmt.Println() // 每行结束换行
	}
}

// ========================================================================

// case63 打印图片的灰度矩阵，在终端显示二维数值

func case63() {
	src := getTestImg()
	// 提取灰度矩阵
	grayMatrix, err := ExtractGrayMatrix(src)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 打印灰度矩阵
	PrintGrayMatrix(grayMatrix)
}

// RGBToGray 将RGB值转换为灰度值（0-255）
// 公式：gray = 0.299*R + 0.587*G + 0.114*B（人眼感知加权）
func RGBToGray(r, g, b uint8) uint8 {
	grayFloat := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(grayFloat)
}

// ExtractGrayMatrix 提取缩放后的灰度矩阵
func ExtractGrayMatrix(img image.Image) ([][]uint8, error) {
	// 获取原图尺寸
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y
	if origWidth == 0 || origHeight == 0 {
		return nil, fmt.Errorf("图片尺寸无效")
	}

	// 初始化灰度矩阵（目标尺寸）
	grayMatrix := make([][]uint8, targetHeight)
	for i := range grayMatrix {
		grayMatrix[i] = make([]uint8, targetWidth)
	}

	// 计算缩放步长（原图到目标尺寸的采样间隔）
	stepX := float64(origWidth) / float64(targetWidth)
	stepY := float64(origHeight) / float64(targetHeight)

	// 遍历目标矩阵的每个位置，采样并转换为灰度值
	for ty := 0; ty < targetHeight; ty++ {
		for tx := 0; tx < targetWidth; tx++ {
			// 计算原图中对应的采样坐标（取整）
			origX := int(float64(tx)*stepX) + bounds.Min.X
			origY := int(float64(ty)*stepY) + bounds.Min.Y

			// 确保坐标在原图范围内（避免越界）
			if origX >= bounds.Max.X {
				origX = bounds.Max.X - 1
			}
			if origY >= bounds.Max.Y {
				origY = bounds.Max.Y - 1
			}

			// 提取像素的RGBA值（转换为0-255的RGB）
			r, g, b, _ := img.At(origX, origY).RGBA()
			r8 := uint8(r >> 8) // 16位转8位（RGBA返回0-65535）
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 转换为灰度值并存入矩阵
			grayMatrix[ty][tx] = RGBToGray(r8, g8, b8)
		}
	}

	return grayMatrix, nil
}

// PrintGrayMatrix 在终端打印灰度矩阵
func PrintGrayMatrix(matrix [][]uint8) {
	fmt.Printf("\n===== 灰度矩阵（%dx%d） =====\n", len(matrix), len(matrix[0]))
	for _, row := range matrix {
		for _, val := range row {
			fmt.Printf("%3d ", val) // 占3位宽度，右对齐，确保列对齐
		}
		fmt.Println() // 每行结束换行
	}
}

// ========================================================================

// case64 打印图片的二值矩阵，在终端显示二维数值

func case64() {
	src := getTestImg()
	// 提取二值矩阵
	binaryMatrix, err := ExtractBinaryMatrix(src)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 打印二值矩阵
	PrintBinaryMatrix(binaryMatrix)
}

// 配置：终端显示的二值矩阵尺寸和阈值
const (
	threshold = 128 // 二值化阈值（灰度值 >= threshold 为1，否则为0）
)

// ExtractBinaryMatrix 提取缩放后的二值矩阵
func ExtractBinaryMatrix(img image.Image) ([][]int, error) {
	// 获取原图尺寸
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y
	if origWidth == 0 || origHeight == 0 {
		return nil, fmt.Errorf("图片尺寸无效")
	}

	// 初始化二值矩阵（目标尺寸，元素为0或1）
	binaryMatrix := make([][]int, targetHeight)
	for i := range binaryMatrix {
		binaryMatrix[i] = make([]int, targetWidth)
	}

	// 计算缩放步长（原图到目标尺寸的采样间隔）
	stepX := float64(origWidth) / float64(targetWidth)
	stepY := float64(origHeight) / float64(targetHeight)

	// 遍历目标矩阵的每个位置，采样并转换为二值
	for ty := 0; ty < targetHeight; ty++ {
		for tx := 0; tx < targetWidth; tx++ {
			// 计算原图中对应的采样坐标（取整）
			origX := int(float64(tx)*stepX) + bounds.Min.X
			origY := int(float64(ty)*stepY) + bounds.Min.Y

			// 确保坐标在原图范围内
			if origX >= bounds.Max.X {
				origX = bounds.Max.X - 1
			}
			if origY >= bounds.Max.Y {
				origY = bounds.Max.Y - 1
			}

			// 提取RGB并转为灰度值
			r, g, b, _ := img.At(origX, origY).RGBA()
			gray := RGBToGray(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			// 二值化：灰度 >= 阈值为1，否则为0
			if gray >= threshold {
				binaryMatrix[ty][tx] = 1
			} else {
				binaryMatrix[ty][tx] = 0
			}
		}
	}

	return binaryMatrix, nil
}

// PrintBinaryMatrix 在终端打印二值矩阵
func PrintBinaryMatrix(matrix [][]int) {
	fmt.Printf("\n===== 二值矩阵（阈值=%d，%dx%d） =====\n", threshold, len(matrix), len(matrix[0]))
	for _, row := range matrix {
		for _, val := range row {
			fmt.Printf("%2d ", val) // 占2位宽度，确保列对齐
		}
		fmt.Println() // 每行结束换行
	}
}

// ========================================================================

// case65  图像的线性变换

// 线性变换公式：核心公式newVal = a×oldVal + b中：
//a控制对比度：a > 1增强对比度（明暗差异变大），0 < a < 1降低对比度（明暗差异变小），a < 0反转颜色（亮变暗、暗变亮）。
//b控制亮度：b > 0整体变亮，b < 0整体变暗（例如b=50会让每个像素值增加 50）。

// 示例效果
//增强对比度 + 增亮：a=1.5, b=20 → 图片明暗差异更明显，整体更亮。
//降低对比度 + 变暗：a=0.5, b=-30 → 色彩更柔和，整体偏暗。
//颜色反转：a=-1, b=255 → 公式变为newVal = 255 - oldVal，实现负片效果。

//扩展场景
//图像预处理：调整亮度 / 对比度以提升后续识别（如 OCR、目标检测）的准确性。
//风格化：通过反转颜色（a=-1）生成负片效果，或通过低对比度（a=0.2）生成朦胧感。
//批量处理：结合文件遍历，对文件夹内所有图片应用统一变换（如批量提亮暗图）。

func case65() {
	src := getTest2Img()
	a := 1.5  // 0.5   // 缩放系数（增强对比度）
	b := 20.0 // -30.0 // 偏移量（增加亮度）
	// 应用线性变换
	transformedImg := LinearTransform(src, a, b)
	outputFileName := "output_case65.jpg"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, transformedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		panic(err)
	}
}

// LinearTransform 对图片应用线性变换：newVal = a*oldVal + b
// 参数：
//
//	img: 输入图片
//	a: 缩放系数（控制对比度，a>1增强，0<a<1减弱，a<0反转）
//	b: 偏移量（控制亮度，b>0变亮，b<0变暗）
//
// 返回：处理后的图片
func LinearTransform(img image.Image, a, b float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建输出图片（RGBA格式，支持透明通道）
	outImg := image.NewRGBA(bounds)

	// 遍历每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取原像素的RGBA值（范围0-65535）
			r, g, bVal, aVal := img.At(x, y).RGBA()

			// 转换为0-255范围（右移8位）
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(bVal >> 8)

			// 对R、G、B通道应用线性变换，并截断到0-255
			newR := clampLinear(float64(r8), a, b)
			newG := clampLinear(float64(g8), a, b)
			newB := clampLinear(float64(b8), a, b)

			// 写入新像素（Alpha通道保持不变，转回0-65535）
			outImg.SetRGBA(
				x, y,
				color.RGBA{
					R: newR,
					G: newG,
					B: newB,
					A: uint8(aVal >> 8), // 保留原透明度
				},
			)
		}
	}

	return outImg
}

// clampLinear 应用线性变换并截断到0-255
func clampLinear(val, a, b float64) uint8 {
	result := a*val + b
	// 确保结果在0-255范围内（防止溢出）
	if result < 0 {
		return 0
	}
	if result > 255 {
		return 255
	}
	return uint8(result)
}

// ========================================================================

// case66 图像细化

// 图像细化（Image Thinning）是将二值图像中的线条或区域缩减为单像素宽度骨架的过程，
//核心是保留图像拓扑结构（如连接性、分支）的同时去除冗余像素，广泛用于指纹识别、文字识别、轮廓分析等场景。
//常用的细化算法中，Zhang-Suen 算法因实现简单、效果稳定被广泛采用

func case66() {
	src := getTestImg()
	// 2. 针对文字的精准二值化
	bin := BinarizeForText(src)
	fmt.Println("二值化完成")

	// 3. 应用Zhang-Suen细化
	thinned := ZhangSuenThinningForText(bin)
	fmt.Println("细化完成")
	outputPath := "output_case66.png"
	// 4. 保存结果
	if err := SaveThinned(thinned, outputPath); err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}
	fmt.Printf("细化结果已保存到 %s\n", outputPath)
}

// 【优化1：针对文字的二值化阈值】
// 黑色文字在浅色背景上，将阈值调低至50，确保文字完全转为前景
const textThreshold = 50

// 8邻域索引（顺时针：p2-p9，对应坐标偏移）
var neighbors = []image.Point{
	{0, -1},  // p2: (x, y-1)
	{1, -1},  // p3: (x+1, y-1)
	{1, 0},   // p4: (x+1, y)
	{1, 1},   // p5: (x+1, y+1)
	{0, 1},   // p6: (x, y+1)
	{-1, 1},  // p7: (x-1, y+1)
	{-1, 0},  // p8: (x-1, y)
	{-1, -1}, // p9: (x-1, y-1)
}

// 【优化2：精准二值化】
// 针对黑色文字+浅色背景，强制将文字转为前景（255），背景转为0
func BinarizeForText(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bin := make([][]uint8, height)
	for y := 0; y < height; y++ {
		bin[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// 转为灰度后，低于阈值的文字像素设为255，否则为0
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			if gray <= textThreshold {
				bin[y][x] = 255 // 文字（前景）
			} else {
				bin[y][x] = 0 // 背景
			}
		}
	}
	return bin
}

// countForeground 计算8邻域中前景像素（255）的数量
func countForeground(bin [][]uint8, x, y int) int {
	count := 0
	height, width := len(bin), len(bin[0])
	for _, n := range neighbors {
		nx, ny := x+n.X, y+n.Y
		// 边界外视为背景
		if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
			count++
		}
	}
	return count
}

// countConnections 计算8邻域的连接数（衡量像素的"拐角"程度）
func countConnections(bin [][]uint8, x, y int) int {
	height, width := len(bin), len(bin[0])
	// 取p2-p9的二值（1=前景，0=背景）
	p := make([]int, 8)
	for i, n := range neighbors {
		nx, ny := x+n.X, y+n.Y
		if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
			p[i] = 1
		}
	}
	// 连接数=相邻像素从0→1的次数（p9与p2相邻）
	count := 0
	for i := 0; i < 8; i++ {
		j := (i + 1) % 8
		if p[i] == 0 && p[j] == 1 {
			count++
		}
	}
	return count
}

// 【优化3：增加端点保护】
// 文字的端点（N(p1)=1）不应被删除，避免笔画断裂
func ZhangSuenThinningForText(bin [][]uint8) [][]uint8 {
	height, width := len(bin), len(bin[0])
	// 复制原图避免修改输入
	thinned := make([][]uint8, height)
	for y := range bin {
		thinned[y] = make([]uint8, width)
		copy(thinned[y], bin[y])
	}

	for {
		deleted := false // 标记本轮是否删除像素
		// 第一步迭代：标记符合条件的像素
		toDelete1 := make([][]bool, height)
		for y := 0; y < height; y++ {
			toDelete1[y] = make([]bool, width)
			for x := 0; x < width; x++ {
				if thinned[y][x] != 255 {
					continue // 跳过背景
				}
				n := countForeground(thinned, x, y)
				c := countConnections(thinned, x, y)
				// 【新增】端点保护：N(p1)=1时不删除
				if n == 1 {
					continue
				}
				// 原条件：2 ≤ N(p1) ≤ 6 且 C(p1)=1
				if n < 2 || n > 6 || c != 1 {
					continue
				}
				// 条件4-5：p2*p4*p6=0 且 p4*p6*p8=0
				p2 := getNeighbor(thinned, x, y, 0)
				p4 := getNeighbor(thinned, x, y, 2)
				p6 := getNeighbor(thinned, x, y, 4)
				p8 := getNeighbor(thinned, x, y, 6)
				if p2*p4*p6 == 0 && p4*p6*p8 == 0 {
					toDelete1[y][x] = true
					deleted = true
				}
			}
		}
		// 执行第一步删除
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if toDelete1[y][x] {
					thinned[y][x] = 0
				}
			}
		}

		// 第二步迭代：标记符合条件的像素
		toDelete2 := make([][]bool, height)
		for y := 0; y < height; y++ {
			toDelete2[y] = make([]bool, width)
			for x := 0; x < width; x++ {
				if thinned[y][x] != 255 {
					continue // 跳过背景
				}
				n := countForeground(thinned, x, y)
				c := countConnections(thinned, x, y)
				// 【新增】端点保护：N(p1)=1时不删除
				if n == 1 {
					continue
				}
				// 原条件：2 ≤ N(p1) ≤ 6 且 C(p1)=1
				if n < 2 || n > 6 || c != 1 {
					continue
				}
				// 条件4-5：p2*p4*p8=0 且 p2*p6*p8=0
				p2 := getNeighbor(thinned, x, y, 0)
				p4 := getNeighbor(thinned, x, y, 2)
				p6 := getNeighbor(thinned, x, y, 4)
				p8 := getNeighbor(thinned, x, y, 6)
				if p2*p4*p8 == 0 && p2*p6*p8 == 0 {
					toDelete2[y][x] = true
					deleted = true
				}
			}
		}
		// 执行第二步删除
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if toDelete2[y][x] {
					thinned[y][x] = 0
				}
			}
		}

		// 若本轮无删除，终止迭代
		if !deleted {
			break
		}
	}

	return thinned
}

// getNeighbor 获取指定邻域像素的二值（1=前景，0=背景）
func getNeighbor(bin [][]uint8, x, y, idx int) int {
	n := neighbors[idx]
	nx, ny := x+n.X, y+n.Y
	height, width := len(bin), len(bin[0])
	if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
		return 1
	}
	return 0
}

// SaveThinned 保存细化后的图像
func SaveThinned(thinned [][]uint8, path string) error {
	height, width := len(thinned), len(thinned[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := thinned[y][x]
			img.SetRGBA(x, y, color.RGBA{val, val, val, 255}) // 灰度图（骨架为白）
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// ========================================================================

// case67 图片的孔洞填充
// 适用场景：适用于印刷文字修复（填补文字内部的孔洞）、工业缺陷检测（识别封闭区域的孔洞）、图像分割后处理等场景。

func case67() {
	src := getTestImg()
	// 2. 二值化
	bin, err := Binarize(src)
	if err != nil {
		fmt.Printf("二值化失败: %v\n", err)
		return
	}

	// 3. 识别孔洞
	holes := FindHoles(bin)
	fmt.Printf("识别到 %d 个孔洞\n", len(holes))

	// 4. 填充孔洞
	FillHoles(bin, holes)
	outputPath := "output_case67.png"
	// 5. 保存结果
	if err := SaveImage(bin, outputPath); err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		return
	}

	fmt.Printf("孔洞填充完成，结果已保存到 %s\n", outputPath)
}

// SaveImage 保存处理后的图像
func SaveImage(bin [][]uint8, path string) error {
	height, width := len(bin), len(bin[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := bin[y][x]
			img.SetRGBA(x, y, color.RGBA{val, val, val, 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

const (
	Background = 0   // 背景像素值
	Foreground = 255 // 前景像素值
	Threshold  = 128 // 二值化阈值（可根据图片调整）
)

// 四连通方向（上下左右），若需八连通可添加对角线方向
var dirs = []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

// Binarize 二值化图像，返回二维像素数组
func Binarize(img image.Image) ([][]uint8, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bin := make([][]uint8, height)
	for y := 0; y < height; y++ {
		bin[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			if gray < Threshold {
				bin[y][x] = Foreground // 深色区域设为前景
			} else {
				bin[y][x] = Background // 浅色区域设为背景
			}
		}
	}
	return bin, nil
}

// HoleRegion 表示一个孔洞区域的所有像素坐标
type HoleRegion struct {
	pixels      []image.Point // 孔洞内所有像素的坐标
	touchesEdge bool          // 是否接触图像边缘（true则不是孔洞）
}

// FindHoles 识别所有孔洞区域
func FindHoles(bin [][]uint8) []HoleRegion {
	height, width := len(bin), len(bin[0])
	visited := make([][]bool, height)
	for y := range visited {
		visited[y] = make([]bool, width)
	}
	var holes []HoleRegion

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bin[y][x] == Background && !visited[y][x] {
				// BFS遍历当前背景连通域
				queue := list.New()
				queue.PushBack(image.Point{x, y})
				visited[y][x] = true
				var region []image.Point
				touchesEdge := false

				for queue.Len() > 0 {
					p := queue.Remove(queue.Front()).(image.Point)
					region = append(region, p)

					// 检查是否接触图像边缘
					if p.X == 0 || p.X == width-1 || p.Y == 0 || p.Y == height-1 {
						touchesEdge = true
					}

					// 遍历四连通邻居
					for _, d := range dirs {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == Background && !visited[ny][nx] {
							visited[ny][nx] = true
							queue.PushBack(image.Point{nx, ny})
						}
					}
				}

				// 仅保留不接触边缘的背景连通域（即孔洞）
				if !touchesEdge {
					holes = append(holes, HoleRegion{pixels: region, touchesEdge: touchesEdge})
				}
			}
		}
	}

	return holes
}

// FillHoles 填充孔洞区域
func FillHoles(bin [][]uint8, holes []HoleRegion) {
	for _, hole := range holes {
		for _, p := range hole.pixels {
			bin[p.Y][p.X] = Foreground
		}
	}
}

// ========================================================================

// case68 图像的小波阈值去噪

// 算法原理
//小波阈值去噪的核心是利用小波变换将图像分解为低频近似分量（保留图像主要结构）和高频细节分量（包含噪声和边缘），
//然后对高频分量应用阈值函数（收缩噪声系数），最后通过逆变换重构去噪后的图像。

func case68() {
	src := "./test2.jpg"

	// 1. 读取带噪声的图像
	noisyMatrix, err := ReadImageToMatrix(src)
	if err != nil {
		fmt.Printf("读取图像失败: %v\n", err)
		return
	}

	// 2. 小波分解
	low, highX, highY, highXY := HaarWaveletDecompose(noisyMatrix)

	// 3. 估计噪声水平
	noiseLevel := EstimateNoiseLevel(highX, highY, highXY)
	fmt.Printf("估计噪声水平: %.2f\n", noiseLevel)

	// 4. 高频系数阈值处理
	WaveletThresholding(highX, highY, highXY, noiseLevel)

	// 5. 小波重构
	denoisedMatrix := HaarWaveletReconstruct(low, highX, highY, highXY)

	outputPath := "output_case68.png"

	// 6. 保存去噪结果
	if err := SaveImage(denoisedMatrix, outputPath); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("小波阈值去噪完成，结果已保存到 %s\n", outputPath)
}

const (
	ThresholdRatio = 0.5 // 阈值系数（可根据噪声强度调整，范围0.0~0.5）
)

// 二维图像的像素矩阵（灰度值0-255）
type ImageMatrix [][]uint8

// ReadImage 读取图片并转为灰度矩阵
func ReadImageToMatrix(path string) (ImageMatrix, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	matrix := make(ImageMatrix, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			matrix[y][x] = gray
		}
	}
	return matrix, nil
}

// HaarWaveletDecompose 对图像进行Haar小波分解（一级）
func HaarWaveletDecompose(matrix ImageMatrix) (low, highX, highY, highXY ImageMatrix) {
	height, width := len(matrix), len(matrix[0])
	// 确保尺寸为偶数（Haar小波分解要求长度为2的倍数）
	if height%2 != 0 {
		height--
	}
	if width%2 != 0 {
		width--
	}

	// 初始化四个子带（低频、水平高频、垂直高频、对角高频）
	low = make(ImageMatrix, height/2)
	highX = make(ImageMatrix, height/2)
	highY = make(ImageMatrix, height/2)
	highXY = make(ImageMatrix, height/2)
	for y := 0; y < height/2; y++ {
		low[y] = make([]uint8, width/2)
		highX[y] = make([]uint8, width/2)
		highY[y] = make([]uint8, width/2)
		highXY[y] = make([]uint8, width/2)
	}

	// 对每一行进行一维Haar分解
	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x += 2 {
			// 提取2x2块的四个像素
			a := float64(matrix[y][x])
			b := float64(matrix[y][x+1])
			c := float64(matrix[y+1][x])
			d := float64(matrix[y+1][x+1])

			// 计算低频和高频系数
			avg := (a + b + c + d) / 4.0
			hl := (a + b - c - d) / 4.0
			lh := (a - b + c - d) / 4.0
			hh := (a - b - c + d) / 4.0

			// 存入子带（注意坐标映射）
			lowY := y / 2
			lowX := x / 2
			low[lowY][lowX] = uint8(math.Round(avg))
			highX[lowY][lowX] = uint8(math.Round(hl))
			highY[lowY][lowX] = uint8(math.Round(lh))
			highXY[lowY][lowX] = uint8(math.Round(hh))
		}
	}
	return
}

// HaarWaveletReconstruct 小波逆变换（重构图像）
func HaarWaveletReconstruct(low, highX, highY, highXY ImageMatrix) ImageMatrix {
	lowHeight, lowWidth := len(low), len(low[0])
	height, width := lowHeight*2, lowWidth*2

	matrix := make(ImageMatrix, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]uint8, width)
	}

	// 对每个2x2块进行逆变换
	for y := 0; y < lowHeight; y++ {
		for x := 0; x < lowWidth; x++ {
			avg := float64(low[y][x])
			hl := float64(highX[y][x])
			lh := float64(highY[y][x])
			hh := float64(highXY[y][x])

			// 逆变换公式
			a := avg + hl + lh + hh
			b := avg + hl - lh - hh
			c := avg - hl + lh - hh
			d := avg - hl - lh + hh

			// 映射回原坐标并裁剪到0-255
			matrix[2*y][2*x] = clamp(a)
			matrix[2*y][2*x+1] = clamp(b)
			matrix[2*y+1][2*x] = clamp(c)
			matrix[2*y+1][2*x+1] = clamp(d)
		}
	}
	return matrix
}

// clamp 将数值裁剪到0-255范围
func clamp(val float64) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

// WaveletThresholding 对高频系数应用软阈值
func WaveletThresholding(highX, highY, highXY ImageMatrix, noiseLevel float64) {
	height, width := len(highX), len(highX[0])
	threshold := noiseLevel * ThresholdRatio

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 软阈值：|x| < T → 0；|x| ≥ T → x - sign(x)*T
			processCoeff := func(coeff *uint8) {
				val := float64(*coeff)
				if math.Abs(val) < threshold {
					*coeff = 0
				} else {
					*coeff = uint8(math.Round(val - math.Copysign(threshold, val)))
				}
			}

			processCoeff(&highX[y][x])
			processCoeff(&highY[y][x])
			processCoeff(&highXY[y][x])
		}
	}
}

// EstimateNoiseLevel 估计噪声水平（基于高频系数的标准差）
func EstimateNoiseLevel(highX, highY, highXY ImageMatrix) float64 {
	var sum float64
	count := 0
	height, width := len(highX), len(highX[0])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += math.Pow(float64(highX[y][x]), 2)
			sum += math.Pow(float64(highY[y][x]), 2)
			sum += math.Pow(float64(highXY[y][x]), 2)
			count += 3
		}
	}

	if count == 0 {
		return 0
	}
	variance := sum / float64(count)
	return math.Sqrt(variance)
}

// ========================================================================

// case69 Sobel算子边缘检测

// 算法原理
//Sobel 算子通过两个 3×3 卷积核（分别检测水平梯度和垂直梯度），计算每个像素的梯度幅值，从而识别图像边缘：
//水平梯度核（Gₓ）：[[-1, 0, 1], [-2, 0, 2], [-1, 0, 1]]
//垂直梯度核（Gᵧ）：[[-1, -2, -1], [0, 0, 0], [1, 2, 1]]
//梯度幅值：|Gₓ| + |Gᵧ|（或√(Gₓ² + Gᵧ²)，前者更高效）
//阈值处理：梯度幅值超过阈值的像素视为边缘（设为 255），否则为背景（0）。

func case69() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	// 执行Sobel边缘检测
	edges := SobelEdgeDetection(gray, threshold)

	// 保存结果
	outputPath := "output_case69.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// 转为灰度矩阵
func ToGrayMatrix(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	gray := make([][]uint8, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 加权平均转为灰度（人眼对绿色敏感更高）
			grayVal := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			gray[y][x] = grayVal
		}
	}
	return gray
}

// Sobel边缘检测（返回二值边缘图，255为边缘，0为背景）
func SobelEdgeDetection(gray [][]uint8, threshold uint8) [][]uint8 {
	height := len(gray)
	if height == 0 {
		return nil
	}
	width := len(gray[0])
	edges := make([][]uint8, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]uint8, width)
	}

	// 遍历所有非边界像素（确保3×3邻域有效）
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 计算水平梯度Gx
			gx := -int(gray[y-1][x-1]) - 2*int(gray[y][x-1]) - int(gray[y+1][x-1]) +
				int(gray[y-1][x+1]) + 2*int(gray[y][x+1]) + int(gray[y+1][x+1])

			// 计算垂直梯度Gy
			gy := -int(gray[y-1][x-1]) - 2*int(gray[y-1][x]) - int(gray[y-1][x+1]) +
				int(gray[y+1][x-1]) + 2*int(gray[y+1][x]) + int(gray[y+1][x+1])

			// 梯度幅值（L1范数，|Gx| + |Gy|，高效且效果接近L2）
			magnitude := math.Abs(float64(gx)) + math.Abs(float64(gy))

			// 阈值处理
			if magnitude >= float64(threshold) {
				edges[y][x] = 255
			} else {
				edges[y][x] = 0
			}
		}
	}

	return edges
}

// ========================================================================

// case70 Roberts算子边缘检测

// 算法原理
//Roberts 算子通过两个 2×2 卷积核检测梯度：
//水平梯度核（Gₓ）：[[1, 0], [0, -1]]
//垂直梯度核（Gᵧ）：[[0, 1], [-1, 0]]
//梯度幅值：|Gₓ| + |Gᵧ|（或√(Gₓ² + Gᵧ²)，前者更高效）
//阈值处理：梯度幅值超过阈值的像素视为边缘（设为 255），否则为背景（0）。

func case70() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	edges := RobertsEdgeDetection(gray, threshold)

	// 保存结果
	outputPath := "output_case70.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// Roberts边缘检测（返回二值边缘图，255为边缘，0为背景）
func RobertsEdgeDetection(gray [][]uint8, threshold uint8) [][]uint8 {
	height := len(gray)
	if height < 2 {
		return nil
	}
	width := len(gray[0])
	if width < 2 {
		return nil
	}

	edges := make([][]uint8, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]uint8, width)
	}

	// 遍历所有非边界像素（确保2×2邻域有效）
	for y := 0; y < height-1; y++ {
		for x := 0; x < width-1; x++ {
			// 计算水平梯度Gx = gray[y][x] - gray[y+1][x+1]
			gx := int(gray[y][x]) - int(gray[y+1][x+1])
			// 计算垂直梯度Gy = gray[y][x+1] - gray[y+1][x]
			gy := int(gray[y][x+1]) - int(gray[y+1][x])

			// 梯度幅值（L1范数，|Gx| + |Gy|）
			magnitude := math.Abs(float64(gx)) + math.Abs(float64(gy))

			// 阈值处理
			if magnitude >= float64(threshold) {
				edges[y][x] = 255
			} else {
				edges[y][x] = 0
			}
		}
	}

	return edges
}

// ========================================================================

// case71 Prewitt算子边缘检测

// 算法原理
//Prewitt 算子通过两个 3×3 卷积核分别计算水平和垂直方向的梯度：
//水平梯度核（Gₓ）：检测垂直边缘（左右变化）
//[[-1, 0, 1], [-1, 0, 1], [-1, 0, 1]]
//垂直梯度核（Gᵧ）：检测水平边缘（上下变化）
//[[-1, -1, -1], [0, 0, 0], [1, 1, 1]]

func case71() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	edges := PrewittEdgeDetection(gray, threshold)

	// 保存结果
	outputPath := "output_case71.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// PrewittEdgeDetection 应用Prewitt算子检测边缘
// 参数：gray为灰度矩阵，threshold为边缘阈值（0-255）
// 返回：二值边缘矩阵（255为边缘，0为背景）
func PrewittEdgeDetection(gray [][]uint8, threshold uint8) [][]uint8 {
	height := len(gray)
	if height < 3 { // 至少3行才能应用3x3核
		return nil
	}
	width := len(gray[0])
	if width < 3 { // 至少3列
		return nil
	}

	edges := make([][]uint8, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]uint8, width)
	}

	// 遍历所有非边界像素（确保3x3邻域有效）
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 计算水平梯度Gx（垂直边缘）
			gx := -int(gray[y-1][x-1]) + int(gray[y-1][x+1]) +
				-int(gray[y][x-1]) + int(gray[y][x+1]) +
				-int(gray[y+1][x-1]) + int(gray[y+1][x+1])

			// 计算垂直梯度Gy（水平边缘）
			gy := -int(gray[y-1][x-1]) - int(gray[y-1][x]) - int(gray[y-1][x+1]) +
				int(gray[y+1][x-1]) + int(gray[y+1][x]) + int(gray[y+1][x+1])

			// 梯度幅值（用L1范数，计算高效；L2范数为math.Sqrt(float64(gx*gx + gy*gy))）
			magnitude := math.Abs(float64(gx)) + math.Abs(float64(gy))

			// 阈值判断：超过阈值则为边缘
			if magnitude >= float64(threshold) {
				edges[y][x] = 255
			} else {
				edges[y][x] = 0
			}
		}
	}

	return edges
}

// ========================================================================

// case72 Laplacian算子边缘检测

// 算法原理
//Laplacian 算子通过计算图像的二阶导数来识别边缘，核心是 3×3 卷积核（常用形式）：
//标准核 1（4 邻域）：[[0, 1, 0], [1, -4, 1], [0, 1, 0]]
//标准核 2（8 邻域）：[[1, 1, 1], [1, -8, 1], [1, 1, 1]]

func case72() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	edges := LaplacianEdgeDetection(gray, threshold)

	// 保存结果
	outputPath := "output_case72.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// 选择Laplacian卷积核（8邻域核，对边缘更敏感）
var laplacianKernelSZ = [3][3]int{
	{1, 1, 1},
	{1, -8, 1},
	{1, 1, 1},
}

// LaplacianEdgeDetection 应用Laplacian算子检测边缘
func LaplacianEdgeDetection(gray [][]uint8, threshold uint8) [][]uint8 {
	height := len(gray)
	if height < 3 {
		return nil
	}
	width := len(gray[0])
	if width < 3 {
		return nil
	}

	edges := make([][]uint8, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]uint8, width)
	}

	// 遍历非边界像素（确保3×3邻域有效）
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 与Laplacian核卷积计算二阶导数响应
			var response int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					// 核索引映射：kernel[ky+1][kx+1] 对应核的位置
					response += int(gray[y+ky][x+kx]) * laplacianKernelSZ[ky+1][kx+1]
				}
			}

			// 取绝对值作为边缘强度
			strength := math.Abs(float64(response))

			// 阈值处理
			if strength >= float64(threshold) {
				edges[y][x] = 255
			} else {
				edges[y][x] = 0
			}
		}
	}

	return edges
}

// ========================================================================

// case73 LoG算子边缘检测

// 算法原理
//LoG 算子的核心是 “先平滑，再求二阶导数”：
//高斯滤波：用高斯核卷积图像，抑制噪声（高斯函数的标准差 σ 控制平滑程度，σ 越大平滑越强）。
//Laplacian 变换：对平滑后的图像应用 Laplacian 算子，检测灰度的快速变化（边缘）。
//零交叉检测：LoG 响应值穿过零的位置即为边缘（实际实现中常用阈值简化：响应绝对值超过阈值的区域视为边缘）。

func case73() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	edges := LoGEdgeDetection(gray)

	// 保存结果
	outputPath := "output_case73.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// 配置参数
const (
	sigma73      = 1.0 // 高斯核标准差（控制平滑程度，建议0.5-2.0）
	kernelSize73 = 5   // 高斯核尺寸（需为奇数，3/5/7，与sigma匹配）
	threshold73  = 30  // 边缘检测阈值（建议20-50）
)

// GenerateGaussianKernel 生成高斯核（归一化）
func GenerateGaussianKernel(sigma float64, size int) [][]float64 {
	kernel := make([][]float64, size)
	center := size / 2
	var sum float64 // 用于归一化

	// 计算高斯核值
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			// 高斯函数：G(x,y) = (1/(2πσ²)) * e^(-(x²+y²)/(2σ²))
			dx := float64(x - center)
			dy := float64(y - center)
			expTerm := math.Exp(-(dx*dx + dy*dy) / (2 * sigma * sigma))
			kernel[y][x] = expTerm / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}

	// 归一化（确保核的总和为1，避免图像亮度变化）
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}

	return kernel
}

// GaussianBlur 高斯滤波（卷积操作）
func GaussianBlur(gray [][]uint8, kernel [][]float64) [][]uint8 {
	height := len(gray)
	width := len(gray[0])
	ksize := len(kernel)
	radius := ksize / 2

	blurred := make([][]uint8, height)
	for y := 0; y < height; y++ {
		blurred[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			var sum float64
			// 与高斯核卷积
			for ky := 0; ky < ksize; ky++ {
				for kx := 0; kx < ksize; ky++ {
					// 计算图像坐标（边界外视为0）
					imgY := y + ky - radius
					imgX := x + kx - radius
					if imgY >= 0 && imgY < height && imgX >= 0 && imgX < width {
						sum += float64(gray[imgY][imgX]) * kernel[ky][kx]
					}
				}
			}
			// 裁剪到0-255
			blurred[y][x] = clamp73(sum)
		}
	}
	return blurred
}

// Laplacian 对图像应用Laplacian算子（8邻域核）
func Laplacian(img [][]uint8) [][]int {
	height := len(img)
	width := len(img[0])
	laplacianKernel := [3][3]int{ // 8邻域Laplacian核
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}

	response := make([][]int, height)
	for y := 0; y < height; y++ {
		response[y] = make([]int, width)
		for x := 0; x < width; x++ {
			// 边界像素响应为0
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				response[y][x] = 0
				continue
			}
			// 与Laplacian核卷积
			var val int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					val += int(img[y+ky][x+kx]) * laplacianKernel[ky+1][kx+1]
				}
			}
			response[y][x] = val
		}
	}
	return response
}

// LoGEdgeDetection LoG边缘检测（高斯平滑 + Laplacian + 阈值）
func LoGEdgeDetection(gray [][]uint8) [][]uint8 {
	// 1. 生成高斯核并平滑图像
	gaussianKernel := GenerateGaussianKernel(sigma73, kernelSize73)
	blurred := GaussianBlur(gray, gaussianKernel)

	// 2. 应用Laplacian算子
	laplacianResponse := Laplacian(blurred)

	// 3. 阈值处理（取绝对值，超过阈值为边缘）
	height := len(gray)
	width := len(gray[0])
	edges := make([][]uint8, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			if math.Abs(float64(laplacianResponse[y][x])) >= float64(threshold73) {
				edges[y][x] = 255
			} else {
				edges[y][x] = 0
			}
		}
	}

	return edges
}

// clamp 裁剪到0-255
func clamp73(val float64) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

// ========================================================================

// case74 Canny算子边缘检测

// 算法原理
//Canny 算子分 5 个关键步骤：
//高斯滤波：用高斯核平滑图像，抑制噪声（σ 控制平滑程度）。
//梯度计算：用 Sobel 算子计算水平 / 垂直梯度，得到梯度幅值和方向。
//非极大值抑制：在梯度方向上保留局部最大值，将边缘细化为单像素宽度。
//双阈值检测：用高 / 低阈值（如high=100，low=50）将像素分为 “强边缘”“弱边缘”“非边缘”。
//边缘连接：弱边缘若与强边缘相连则保留，否则抑制，形成完整边缘

func case74() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	edges := CannyEdgeDetection(gray)

	// 保存结果
	outputPath := "output_case74.png"
	if err := SaveImage(edges, outputPath); err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("Sobel边缘检测完成，结果已保存到 %s\n", outputPath)
}

// 配置参数（可根据图像调整）
const (
	sigma74         = 1.0 // 高斯核标准差（建议0.8-1.2）
	kernelSize74    = 5   // 高斯核尺寸（奇数，3/5，与sigma匹配）
	highThreshold74 = 100 // 高阈值（强边缘）
	lowThreshold74  = 50  // 低阈值（弱边缘，通常为高阈值的1/2）
)

// 梯度方向量化（0°、45°、90°、135°，单位：弧度）
const (
	dir0case74   = 0               // 0°（水平）
	dir45case74  = math.Pi / 4     // 45°
	dir90case74  = math.Pi / 2     // 90°（垂直）
	dir135case74 = 3 * math.Pi / 4 // 135°
)

// GenerateGaussianKernel 生成高斯核（归一化）
func GenerateGaussianKernel74(sigma float64, size int) [][]float64 {
	kernel := make([][]float64, size)
	center := size / 2
	var sum float64
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := float64(x - center)
			dy := float64(y - center)
			expTerm := math.Exp(-(dx*dx + dy*dy) / (2 * sigma * sigma))
			kernel[y][x] = expTerm / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}
	// 归一化
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}
	return kernel
}

// GaussianBlur 高斯滤波
func GaussianBlur74(gray [][]uint8, kernel [][]float64) [][]uint8 {
	height := len(gray)
	width := len(gray[0])
	ksize := len(kernel)
	radius := ksize / 2

	blurred := make([][]uint8, height)
	for y := 0; y < height; y++ {
		blurred[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			var sum float64
			for ky := 0; ky < ksize; ky++ {
				for kx := 0; kx < ksize; kx++ {
					imgY := y + ky - radius
					imgX := x + kx - radius
					if imgY >= 0 && imgY < height && imgX >= 0 && imgX < width {
						sum += float64(gray[imgY][imgX]) * kernel[ky][kx]
					}
				}
			}
			blurred[y][x] = clamp74(sum)
		}
	}
	return blurred
}

// CalculateGradients 计算梯度幅值和方向（用Sobel算子）
func CalculateGradients(blurred [][]uint8) (magnitude [][]float64, direction [][]float64) {
	height := len(blurred)
	width := len(blurred[0])
	magnitude = make([][]float64, height)
	direction = make([][]float64, height)
	for y := 0; y < height; y++ {
		magnitude[y] = make([]float64, width)
		direction[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 边界像素梯度为0
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				magnitude[y][x] = 0
				direction[y][x] = 0
				continue
			}
			// Sobel算子计算Gx（水平梯度）和Gy（垂直梯度）
			gx := -int(blurred[y-1][x-1]) - 2*int(blurred[y][x-1]) - int(blurred[y+1][x-1]) +
				int(blurred[y-1][x+1]) + 2*int(blurred[y][x+1]) + int(blurred[y+1][x+1])

			gy := -int(blurred[y-1][x-1]) - 2*int(blurred[y-1][x]) - int(blurred[y-1][x+1]) +
				int(blurred[y+1][x-1]) + 2*int(blurred[y+1][x]) + int(blurred[y+1][x+1])

			// 梯度幅值（L2范数）
			magnitude[y][x] = math.Hypot(float64(gx), float64(gy))
			// 梯度方向（弧度，0~π）
			dir := math.Atan2(float64(gy), float64(gx))
			if dir < 0 {
				dir += math.Pi // 统一到0~π范围
			}
			direction[y][x] = dir
		}
	}
	return
}

// NonMaxSuppression 非极大值抑制（细化边缘为单像素）
func NonMaxSuppression(magnitude, direction [][]float64) [][]float64 {
	height := len(magnitude)
	width := len(magnitude[0])
	suppressed := make([][]float64, height)
	for y := 0; y < height; y++ {
		suppressed[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 边界像素直接抑制
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				suppressed[y][x] = 0
				continue
			}

			dir := direction[y][x]
			mag := magnitude[y][x]
			var neighbor1, neighbor2 float64

			// 根据梯度方向确定相邻像素（量化到4个方向）
			switch {
			// 0°（水平方向，左右相邻）
			case (dir >= dir0case74-dir45case74/4 && dir < dir0case74+dir45case74/4) ||
				(dir >= dir90case74+dir45case74*3/4): // 接近180°合并到0°
				neighbor1 = magnitude[y][x-1] // 左
				neighbor2 = magnitude[y][x+1] // 右
			// 45°（对角线，左上-右下）
			case dir >= dir45case74-dir45case74/4 && dir < dir45case74+dir45case74/4:
				neighbor1 = magnitude[y-1][x+1] // 右上
				neighbor2 = magnitude[y+1][x-1] // 左下
			// 90°（垂直方向，上下相邻）
			case dir >= dir90case74-dir45case74/4 && dir < dir90case74+dir45case74/4:
				neighbor1 = magnitude[y-1][x] // 上
				neighbor2 = magnitude[y+1][x] // 下
			// 135°（对角线，左上-右下）
			case dir >= dir135case74-dir45case74/4 && dir < dir135case74+dir45case74/4:
				neighbor1 = magnitude[y-1][x-1] // 左上
				neighbor2 = magnitude[y+1][x+1] // 右下
			}

			// 仅保留局部最大值
			if mag >= neighbor1 && mag >= neighbor2 {
				suppressed[y][x] = mag
			} else {
				suppressed[y][x] = 0
			}
		}
	}
	return suppressed
}

// DoubleThreshold 双阈值分类（强边缘/弱边缘/非边缘）
func DoubleThreshold(suppressed [][]float64) (edges [][]int) {
	height := len(suppressed)
	width := len(suppressed[0])
	edges = make([][]int, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]int, width)
		for x := 0; x < width; x++ {
			val := suppressed[y][x]
			if val >= highThreshold74 {
				edges[y][x] = 2 // 强边缘
			} else if val >= lowThreshold74 {
				edges[y][x] = 1 // 弱边缘
			} else {
				edges[y][x] = 0 // 非边缘
			}
		}
	}
	return
}

// Hysteresis 边缘连接（弱边缘与强边缘相连则保留）
func Hysteresis(edges [][]int) [][]uint8 {
	height := len(edges)
	width := len(edges[0])
	result := make([][]uint8, height)
	for y := 0; y < height; y++ {
		result[y] = make([]uint8, width)
	}

	// 标记已访问的弱边缘
	visited := make([][]bool, height)
	for y := range visited {
		visited[y] = make([]bool, width)
	}

	// 8邻域方向
	dirs := []image.Point{{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1}}

	// 遍历所有强边缘，用BFS连接弱边缘
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if edges[y][x] == 2 { // 强边缘
				queue := list.New()
				queue.PushBack(image.Point{x, y})
				result[y][x] = 255 // 标记为边缘

				for queue.Len() > 0 {
					p := queue.Remove(queue.Front()).(image.Point)
					// 检查8邻域的弱边缘
					for _, d := range dirs {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							if edges[ny][nx] == 1 && !visited[ny][nx] {
								visited[ny][nx] = true
								result[ny][nx] = 255 // 弱边缘连接到强边缘，保留
								queue.PushBack(image.Point{nx, ny})
							}
						}
					}
				}
			}
		}
	}

	return result
}

// CannyEdgeDetection Canny算子主函数
func CannyEdgeDetection(gray [][]uint8) [][]uint8 {
	// 1. 高斯滤波去噪
	gaussianKernel := GenerateGaussianKernel74(sigma74, kernelSize74)
	blurred := GaussianBlur74(gray, gaussianKernel)

	// 2. 计算梯度幅值和方向
	magnitude, direction := CalculateGradients(blurred)

	// 3. 非极大值抑制
	suppressed := NonMaxSuppression(magnitude, direction)

	// 4. 双阈值检测
	thresholded := DoubleThreshold(suppressed)

	// 5. 边缘连接
	edges := Hysteresis(thresholded)

	return edges
}

// clamp 裁剪到0-255
func clamp74(val float64) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

// ========================================================================

// case75 分水岭分割

// 算法原理
//分水岭分割的核心流程：
//1. 预处理：去噪（形态学开运算）、二值化，得到前景（目标）和背景的初步分离。
//2. 前景 / 背景标记：
//确定 “确定前景”（目标内部，通过腐蚀操作获取）。
//确定 “确定背景”（背景区域，通过膨胀操作获取）。
//计算 “未知区域”（前景与背景之间的过渡带）。
//3. 距离变换：对前景计算距离变换，找到每个目标的 “种子点”（距离变换的局部最大值，即目标中心）。
//4. 分水岭变换：以种子点为起点，模拟水流淹没过程，不同种子点的水流相遇处形成 “分水岭”（分割线），最终将图像分割为多个区域。

func case75() {
	src := getTest2Img()
	gray := ToGrayMatrix(src)

	// 2. 二值化
	bin := Binarize75(gray)

	// 3. 预处理：开运算（去噪）
	eroded1 := Erode(bin)
	opened := Dilate(eroded1) // 开运算 = 腐蚀 + 膨胀

	// 4. 确定前景（进一步腐蚀，确保是目标内部）
	foreground := Erode(opened)

	// 5. 确定背景（膨胀开运算结果）
	background := Dilate(opened)
	background = Dilate(background) // 多次膨胀扩大背景

	// 6. 计算未知区域（背景 - 前景）
	height, width := len(bin), len(bin[0])
	unknown := make([][]uint8, height)
	for y := 0; y < height; y++ {
		unknown[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			if background[y][x] == 255 && foreground[y][x] == 0 {
				unknown[y][x] = 255 // 未知区域
			} else {
				unknown[y][x] = 0
			}
		}
	}

	// 7. 距离变换找种子点
	dist := DistanceTransform(opened)
	seeds := MarkSeeds(dist, foreground)

	// 8. 分水岭分割
	labels := Watershed(gray, seeds, unknown)

	outputPath := "output_case75.png"

	// 9. 可视化并保存结果
	vis := VisualizeLabels(labels)
	if err := SaveImage(vis, outputPath); err != nil {
		fmt.Printf("保存图片错误: %v\n", err)
		return
	}

	fmt.Printf("分水岭分割完成，结果已保存到 %s\n", outputPath)

}

// 配置参数
const (
	threshold76     = 128 // 二值化阈值
	kernelSize76    = 3   // 形态学操作核大小（3x3）
	distThreshold76 = 5.0 // 距离变换种子点阈值（控制种子数量）
)

// 形态学结构元素（3x3矩形核）
var kernel76 = [3][3]int{
	{1, 1, 1},
	{1, 1, 1},
	{1, 1, 1},
}

// 图像像素点结构（用于优先队列）
type Pixel struct {
	val   int // 灰度值（高度）
	x, y  int // 坐标
	label int // 所属区域标记（-1为未标记，0为分水岭）
}

// 优先队列（最小堆，用于模拟水流从低到高淹没）
type PriorityQueue []*Pixel

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].val < pq[j].val }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Pixel)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Binarize 二值化（前景255，背景0）
func Binarize75(gray [][]uint8) [][]uint8 {
	height := len(gray)
	width := len(gray[0])
	bin := make([][]uint8, height)
	for y := 0; y < height; y++ {
		bin[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			if gray[y][x] < threshold76 {
				bin[y][x] = 255 // 前景（假设目标为深色）
			} else {
				bin[y][x] = 0 // 背景
			}
		}
	}
	return bin
}

// Erode 形态学腐蚀（去除边缘像素，缩小前景）
func Erode(bin [][]uint8) [][]uint8 {
	height := len(bin)
	width := len(bin[0])
	eroded := make([][]uint8, height)
	for y := 0; y < height; y++ {
		eroded[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// 边界像素腐蚀为背景
			if x < 1 || x >= width-1 || y < 1 || y >= height-1 {
				eroded[y][x] = 0
				continue
			}
			// 3x3邻域全为前景才保留
			allForeground := true
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					if bin[y+ky][x+kx] == 0 {
						allForeground = false
						break
					}
				}
				if !allForeground {
					break
				}
			}
			if allForeground {
				eroded[y][x] = 255
			} else {
				eroded[y][x] = 0
			}
		}
	}
	return eroded
}

// Dilate 形态学膨胀（扩展边缘像素，扩大前景）
func Dilate(bin [][]uint8) [][]uint8 {
	height := len(bin)
	width := len(bin[0])
	dilated := make([][]uint8, height)
	for y := 0; y < height; y++ {
		dilated[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// 边界像素膨胀为背景
			if x < 1 || x >= width-1 || y < 1 || y >= height-1 {
				dilated[y][x] = 0
				continue
			}
			// 3x3邻域有一个前景则保留
			hasForeground := false
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					if bin[y+ky][x+kx] == 255 {
						hasForeground = true
						break
					}
				}
				if hasForeground {
					break
				}
			}
			if hasForeground {
				dilated[y][x] = 255
			} else {
				dilated[y][x] = 0
			}
		}
	}
	return dilated
}

// DistanceTransform 计算前景到背景的欧氏距离变换
func DistanceTransform(bin [][]uint8) [][]float64 {
	height := len(bin)
	width := len(bin[0])
	dist := make([][]float64, height)
	for y := 0; y < height; y++ {
		dist[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			if bin[y][x] == 0 { // 背景距离为0
				dist[y][x] = 0
				continue
			}
			// 找最近的背景像素计算欧氏距离
			minDist := math.MaxFloat64
			for by := 0; by < height; by++ {
				for bx := 0; bx < width; bx++ {
					if bin[by][bx] == 0 {
						d := math.Hypot(float64(x-bx), float64(y-by))
						if d < minDist {
							minDist = d
						}
					}
				}
			}
			dist[y][x] = minDist
		}
	}
	return dist
}

// MarkSeeds 标记前景种子点（距离变换的局部最大值）
func MarkSeeds(dist [][]float64, eroded [][]uint8) [][]int {
	height := len(dist)
	width := len(dist[0])
	seeds := make([][]int, height)
	for y := 0; y < height; y++ {
		seeds[y] = make([]int, width)
		for x := 0; x < width; x++ {
			seeds[y][x] = -1 // -1表示非种子
		}
	}

	// 仅在确定前景（eroded）中找种子
	label := 1
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if eroded[y][x] == 255 && dist[y][x] > distThreshold76 {
				// 判断是否为局部最大值（3x3邻域）
				isMax := true
				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; ky++ {
						if dist[y+ky][x+kx] > dist[y][x] {
							isMax = false
							break
						}
					}
					if !isMax {
						break
					}
				}
				if isMax {
					seeds[y][x] = label
					label++
				}
			}
		}
	}
	return seeds
}

// Watershed 分水岭变换
func Watershed(gray [][]uint8, seeds [][]int, unknown [][]uint8) [][]int {
	height := len(gray)
	width := len(gray[0])
	labels := make([][]int, height) // 标记结果：0为分水岭，>0为区域
	for y := 0; y < height; y++ {
		labels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			labels[y][x] = -1 // 初始未标记
		}
	}

	// 8邻域方向
	dirs := []image.Point{{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1}}

	// 初始化优先队列（按灰度值从小到大，即从低到高淹没）
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// 加入种子点（前景种子）和确定背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if seeds[y][x] > 0 { // 前景种子
				labels[y][x] = seeds[y][x]
				heap.Push(&pq, &Pixel{val: int(gray[y][x]), x: x, y: y, label: seeds[y][x]})
			} else if unknown[y][x] == 0 { // 确定背景（非未知区域）
				labels[y][x] = 0 // 背景标记为0（非分水岭）
				heap.Push(&pq, &Pixel{val: int(gray[y][x]), x: x, y: y, label: 0})
			}
		}
	}

	// 模拟水流淹没
	for pq.Len() > 0 {
		p := heap.Pop(&pq).(*Pixel)
		x, y := p.x, p.y

		// 处理8邻域
		for _, d := range dirs {
			nx, ny := x+d.X, y+d.Y
			if nx >= 0 && nx < width && ny >= 0 && ny < height && labels[ny][nx] == -1 {
				// 未知区域像素，标记为当前区域
				labels[ny][nx] = p.label
				heap.Push(&pq, &Pixel{val: int(gray[ny][nx]), x: nx, y: ny, label: p.label})
			} else if nx >= 0 && nx < width && ny >= 0 && ny < height && labels[ny][nx] != -1 && labels[ny][nx] != p.label {
				// 不同区域相遇，标记为分水岭（0）
				labels[ny][nx] = 0
			}
		}
	}

	return labels
}

// VisualizeLabels 将分割结果可视化（不同区域不同颜色）
func VisualizeLabels(labels [][]int) [][]uint8 {
	height := len(labels)
	width := len(labels[0])
	// 定义颜色表（8种颜色，可扩展）
	colors := []color.RGBA{
		{255, 0, 0, 255},   // 分水岭（红）
		{0, 255, 0, 255},   // 区域1（绿）
		{0, 0, 255, 255},   // 区域2（蓝）
		{255, 255, 0, 255}, // 区域3（黄）
		{255, 0, 255, 255}, // 区域4（品红）
		{0, 255, 255, 255}, // 区域5（青）
		{128, 0, 0, 255},   // 区域6（深红）
		{0, 128, 0, 255},   // 区域7（深绿）
	}

	result := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			label := labels[y][x]
			if label == 0 { // 分水岭
				result.SetRGBA(x, y, colors[0])
			} else if label > 0 { // 区域（循环使用颜色表）
				colorIdx := label % (len(colors) - 1)
				result.SetRGBA(x, y, colors[colorIdx+1])
			} else { // 未标记区域（黑）
				result.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	// 转为灰度矩阵（便于保存）
	vis := make([][]uint8, height)
	for y := 0; y < height; y++ {
		vis[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			vis[y][x] = uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
		}
	}
	return vis
}

// ========================================================================

// case76 图像提取质心

// 算法原理
//预处理：将图像转为二值化（前景255，背景0），区分目标与背景。
//连通域分析：用 BFS/DFS 标记所有独立的目标区域（连通域）。
//质心计算：对每个连通域，计算所有像素的 x 坐标之和与 y 坐标之和，分别除以像素总数，得到质心坐标(Cx, Cy) (学习了解公式可以问问ai)

func case76() {
	src := getTestImg()
	gray := ToGrayMatrix(src)

	// 2. 二值化（区分目标和背景）
	bin := Binarize76(gray)

	// 3. 找到所有连通域（目标区域）
	components := FindConnectedComponents(bin)
	fmt.Printf("检测到 %d 个目标区域\n", len(components))

	// 4. 计算每个区域的质心
	components = CalculateCentroids(components)

	// 5. 输出质心坐标
	for i, comp := range components {
		fmt.Printf("目标 %d 质心: (%d, %d)\n", i+1, comp.centroid.X, comp.centroid.Y)
	}

	outputPath := "output_case76.png"

	// 6. 在图像上绘制质心并保存
	resultImg := DrawCentroids(gray, components)
	if err := SaveImage76(resultImg, outputPath); err != nil {
		fmt.Printf("保存图像错误: %v\n", err)
		return
	}

	fmt.Printf("质心提取完成，结果已保存到 %s\n", outputPath)
}

// SaveImage 保存图像
func SaveImage76(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()
	return png.Encode(file, img)
}

// 连通域结构：包含像素坐标和质心
type ConnectedComponent struct {
	pixels   []image.Point // 区域内所有像素
	centroid image.Point   // 质心坐标
}

// Binarize 二值化：前景255（目标），背景0（非目标）
func Binarize76(gray [][]uint8) [][]uint8 {
	height := len(gray)
	width := len(gray[0])
	bin := make([][]uint8, height)
	for y := 0; y < height; y++ {
		bin[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			if gray[y][x] < threshold76 { // 假设目标为深色，可根据实际反转
				bin[y][x] = 255 // 前景（目标）
			} else {
				bin[y][x] = 0 // 背景
			}
		}
	}
	return bin
}

// FindConnectedComponents 找到所有连通域（4连通）
func FindConnectedComponents(bin [][]uint8) []ConnectedComponent {
	height := len(bin)
	width := len(bin[0])
	visited := make([][]bool, height) // 标记是否访问过
	for y := range visited {
		visited[y] = make([]bool, width)
	}

	var components []ConnectedComponent
	dirs := []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // 4连通方向

	// 遍历所有像素，寻找未访问的前景像素作为连通域起点
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bin[y][x] == 255 && !visited[y][x] {
				// BFS遍历连通域
				queue := list.New()
				queue.PushBack(image.Point{x, y})
				visited[y][x] = true
				var component ConnectedComponent

				// 收集当前连通域的所有像素
				for queue.Len() > 0 {
					p := queue.Remove(queue.Front()).(image.Point)
					component.pixels = append(component.pixels, p)

					// 检查4邻域
					for _, d := range dirs {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 && !visited[ny][nx] {
							visited[ny][nx] = true
							queue.PushBack(image.Point{nx, ny})
						}
					}
				}

				components = append(components, component)
			}
		}
	}

	return components
}

// CalculateCentroids 计算每个连通域的质心
func CalculateCentroids(components []ConnectedComponent) []ConnectedComponent {
	for i := range components {
		comp := &components[i]
		if len(comp.pixels) == 0 {
			continue // 空区域跳过
		}

		var sumX, sumY int
		for _, p := range comp.pixels {
			sumX += p.X
			sumY += p.Y
		}

		// 质心坐标（四舍五入为整数）
		n := len(comp.pixels)
		comp.centroid.X = sumX / n
		comp.centroid.Y = sumY / n
	}
	return components
}

// DrawCentroids 在图像上绘制质心（十字标记）
func DrawCentroids(gray [][]uint8, components []ConnectedComponent) image.Image {
	height := len(gray)
	width := len(gray[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 先绘制灰度图作为背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := gray[y][x]
			img.SetRGBA(x, y, color.RGBA{val, val, val, 255})
		}
	}

	// 绘制每个质心（红色十字，线宽2像素）
	crossSize := 5 // 十字大小
	for _, comp := range components {
		cx, cy := comp.centroid.X, comp.centroid.Y
		if cx < 0 || cx >= width || cy < 0 || cy >= height {
			continue // 质心超出图像范围
		}

		// 绘制十字横线（左右）
		for dx := -crossSize; dx <= crossSize; dx++ {
			x := cx + dx
			if x >= 0 && x < width {
				img.SetRGBA(x, cy, color.RGBA{255, 0, 0, 255})   // 红线
				img.SetRGBA(x, cy+1, color.RGBA{255, 0, 0, 255}) // 线宽2
			}
		}

		// 绘制十字竖线（上下）
		for dy := -crossSize; dy <= crossSize; dy++ {
			y := cy + dy
			if y >= 0 && y < height {
				img.SetRGBA(cx, y, color.RGBA{255, 0, 0, 255})   // 红线
				img.SetRGBA(cx+1, y, color.RGBA{255, 0, 0, 255}) // 线宽2
			}
		}
	}

	return img
}

// ========================================================================

// case77 基于 BT.601 标准（常用于传统视频和图像领域）实现的 RGB 到 YUV 转换函数。
//YUV 颜色空间中，Y表示亮度（Luminance），U和V表示色度（Chrominance），该实现适用于 8 位（0-255）的 RGB 输入，
//输出符合 BT.601 标准范围的 YUV 值（Y:16-235，U/V:16-240）。

func case77() {
	// 测试示例：常见颜色的RGB转YUV
	testCases := []struct {
		name    string
		r, g, b uint8
	}{
		{"黑色", 0, 0, 0},
		{"白色", 255, 255, 255},
		{"红色", 255, 0, 0},
		{"绿色", 0, 255, 0},
		{"蓝色", 0, 0, 255},
		{"黄色", 255, 255, 0},
	}

	for _, tc := range testCases {
		y, u, v := RGBToYUV(tc.r, tc.g, tc.b)
		fmt.Printf("%s (R:%d, G:%d, B:%d) → Y:%d, U:%d, V:%d\n",
			tc.name, tc.r, tc.g, tc.b, y, u, v)
	}

	// 测试案例：使用前文RGB转YUV的结果反向验证（应还原原始RGB）
	testCases2 := []struct {
		name                string
		y, u, v             uint8 // 输入YUV（来自RGB转YUV的输出）
		wantR, wantG, wantB uint8 // 期望还原的RGB
	}{
		{"黑色", 16, 128, 128, 0, 0, 0},
		{"白色", 235, 128, 128, 255, 255, 255},
		{"红色", 81, 128, 240, 255, 0, 0},   // 红色RGB转YUV的结果
		{"绿色", 145, 44, 16, 0, 255, 0},    // 绿色RGB转YUV的结果
		{"蓝色", 32, 240, 128, 0, 0, 255},   // 蓝色RGB转YUV的结果
		{"黄色", 210, 44, 240, 255, 255, 0}, // 黄色RGB转YUV的结果
	}

	for _, tc := range testCases2 {
		r, g, b := YUVToRGB(tc.y, tc.u, tc.v)
		fmt.Printf("%s (Y:%d, U:%d, V:%d) → RGB(%d,%d,%d) 期望(%d,%d,%d)\n",
			tc.name, tc.y, tc.u, tc.v, r, g, b, tc.wantR, tc.wantG, tc.wantB)
	}
}

// RGBToYUV 将8位RGB值（0-255）转换为BT.601标准的YUV值
// Y范围：16-235（亮度），U和V范围：16-240（色度）
func RGBToYUV(r, g, b uint8) (y, u, v uint8) {
	// 将RGB转换为浮点数以便计算
	rf := float64(r)
	gf := float64(g)
	bf := float64(b)

	// 计算Y分量（亮度）
	yf := 16.0 + (65.738*rf+129.057*gf+25.064*bf)/256.0
	// 计算U分量（蓝色色度）
	uf := 128.0 + (-37.945*rf-74.494*gf+112.439*bf)/256.0
	// 计算V分量（红色色度）
	vf := 128.0 + (112.439*rf-94.154*gf-18.285*bf)/256.0

	// 四舍五入到整数，并裁剪到标准范围
	yInt := int(yf + 0.5)
	uInt := int(uf + 0.5)
	vInt := int(vf + 0.5)

	// 确保Y在16-235范围内
	yInt = clamp77(yInt, 16, 235)
	// 确保U和V在16-240范围内
	uInt = clamp77(uInt, 16, 240)
	vInt = clamp77(vInt, 16, 240)

	return uint8(yInt), uint8(uInt), uint8(vInt)
}

// clamp 辅助函数：将值限制在[min, max]范围内
func clamp77(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// YUVToRGB 将BT.601标准的YUV值转换为8位RGB值（0-255）
// 输入范围：Y∈[16,235]，U∈[16,240]，V∈[16,240]
// 输出范围：R、G、B∈[0,255]
func YUVToRGB(y, u, v uint8) (r, g, b uint8) {
	// 1. 先将YUV钳位到标准范围（避免输入异常值导致转换错误）
	yClamped := clampY(y)
	uClamped := clampUV(u)
	vClamped := clampUV(v)

	// 2. 转换为浮点数并计算偏移校正（消除基准值）
	yf := float64(yClamped - 16)  // Y基准值16，校正后范围0-219
	uf := float64(uClamped - 128) // U基准值128，校正后范围-112-112
	vf := float64(vClamped - 128) // V基准值128，校正后范围-112-112

	// 3. 应用BT.601逆转换公式计算RGB
	rf := 1.164*yf + 1.596*vf
	gf := 1.164*yf - 0.391*uf - 0.813*vf
	bf := 1.164*yf + 2.018*uf

	// 4. 四舍五入并钳位到0-255范围
	rInt := int(rf + 0.5)
	gInt := int(gf + 0.5)
	bInt := int(bf + 0.5)

	rInt = clampRGB(rInt)
	gInt = clampRGB(gInt)
	bInt = clampRGB(bInt)

	return uint8(rInt), uint8(gInt), uint8(bInt)
}

// 辅助函数：将Y钳位到16-235范围
func clampY(y uint8) uint8 {
	if y < 16 {
		return 16
	}
	if y > 235 {
		return 235
	}
	return y
}

// 辅助函数：将U/V钳位到16-240范围
func clampUV(uv uint8) uint8 {
	if uv < 16 {
		return 16
	}
	if uv > 240 {
		return 240
	}
	return uv
}

// 辅助函数：将RGB钳位到0-255范围
func clampRGB(x int) int {
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return x
}

// ========================================================================

// case78 基于 Y 分量的双边滤波磨皮算法
// 该算法通过在 YUV 颜色空间中对亮度分量（Y）进行双边滤波，实现平滑皮肤瑕疵同时保留五官边缘的效果，最后转换回 RGB 空间输出结果。

// 算法原理
//颜色空间转换：将 RGB 图像转为 YUV，分离出亮度分量（Y）和色度分量（U、V）。皮肤瑕疵主要体现在亮度变化上，仅处理 Y 分量可减少对色彩的干扰。
//双边滤波：一种同时考虑空间距离和像素值差异的滤波方法：
//空间权重：距离越近的像素权重越大（高斯分布）。
//范围权重：亮度差异越小的像素权重越大（高斯分布）。
//两者结合可平滑相似亮度区域（皮肤），同时保留亮度突变区域（五官边缘）。
//逆转换：将处理后的 Y 分量与原始 U、V 分量结合，转回 RGB 图像。

// 扩展，人物描边区域标记，只对指定区域进行磨皮

func case78() {
	src := getImg("./test2.jpg")
	outputPath := "output_case78.jpg"

	// 转换为YUV
	Y, U, V := imgToYUV(src)

	// 双边滤波Y分量
	filteredY := bilateralFilterY(Y)

	// 转换回RGB
	resultImg := yuvToImg(filteredY, U, V)

	// 保存图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()
	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
	println("磨皮完成，结果已保存至:", outputPath)
}

// 双边滤波参数（可调整）
const (
	windowRadius = 3    // 滤波窗口半径（7x7窗口，平滑效果更自然）
	sigmaS       = 2.0  // 空间高斯标准差，控制空间影响范围
	sigmaR       = 30.0 // 亮度差异标准差，控制磨皮强度（值越大平滑越强）
)

// RGB转YUV（严格遵循BT.601标准，修复U/V分量计算）
// 输入：R/G/B (0-255)
// 输出：Y(16-235), U(16-240), V(16-240)
func RGBToYUV78(r, g, b uint8) (y, u, v uint8) {
	rf, gf, bf := float64(r), float64(g), float64(b)

	// 标准BT.601 Y分量公式（亮度）
	yf := 0.299*rf + 0.587*gf + 0.114*bf
	// 转换Y到16-235范围（0-255 → 16-235）
	yf = 16 + (yf * 219 / 255)

	// 标准BT.601 U(Cb)分量公式（蓝色色度）
	// 正确公式：Cb = 128 - 0.1687*R - 0.3313*G + 0.5*B
	uf := 128.0 - 0.168736*rf - 0.331264*gf + 0.5*bf

	// 标准BT.601 V(Cr)分量公式（红色色度）
	// 正确公式：Cr = 128 + 0.5*R - 0.4187*G - 0.0813*B
	vf := 128.0 + 0.5*rf - 0.418688*gf - 0.081312*bf

	// 钳位到标准范围（防止溢出导致偏色）
	y = clamp78(yf, 16, 235)
	u = clamp78(uf, 16, 240)
	v = clamp78(vf, 16, 240)
	return
}

// YUV转RGB（修复逆转换公式，确保颜色还原正确）
// 输入：Y(16-235), U(16-240), V(16-240)
// 输出：R/G/B (0-255)
func YUVToRGB78(y, u, v uint8) (r, g, b uint8) {
	// 将YUV从标准范围转换为归一化值
	yNorm := (float64(y) - 16) * (255.0 / 219.0) // Y:16-235 → 0-255
	uNorm := float64(u) - 128.0                  // U:16-240 → -112-112
	vNorm := float64(v) - 128.0                  // V:16-240 → -112-112

	// 标准BT.601逆转换公式
	rf := yNorm + 1.402*vNorm                   // 红色 = 亮度 + 红色色度影响
	gf := yNorm - 0.34414*uNorm - 0.71414*vNorm // 绿色 = 亮度 - 蓝/红色度影响
	bf := yNorm + 1.772*uNorm                   // 蓝色 = 亮度 + 蓝色色度影响

	// 钳位到0-255（防止溢出导致颜色错误）
	r = clamp78(rf, 0, 255)
	g = clamp78(gf, 0, 255)
	b = clamp78(bf, 0, 255)
	return
}

// 通用钳位函数（处理float64到uint8的范围限制）
func clamp78(val float64, min, max float64) uint8 {
	if val < min {
		return uint8(min)
	}
	if val > max {
		return uint8(max)
	}
	return uint8(val + 0.5) // 四舍五入
}

// 图像转YUV分量（分离亮度和色度）
func imgToYUV(img image.Image) (Y, U, V [][]uint8) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	Y = make([][]uint8, height)
	U = make([][]uint8, height)
	V = make([][]uint8, height)
	for i := range Y {
		Y[i] = make([]uint8, width)
		U[i] = make([]uint8, width)
		V[i] = make([]uint8, width)
	}

	// 遍历每个像素转换为YUV
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 从图像中获取RGB值（注意RGBA返回的是0-65535，需转换为0-255）
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			Y[y][x], U[y][x], V[y][x] = RGBToYUV78(r8, g8, b8)
		}
	}
	return Y, U, V
}

// 双边滤波（仅平滑Y分量，保留U/V分量以维持颜色）
func bilateralFilterY(Y [][]uint8) [][]uint8 {
	height := len(Y)
	if height == 0 {
		return nil
	}
	width := len(Y[0])
	result := make([][]uint8, height)
	for i := range result {
		result[i] = make([]uint8, width)
	}

	// 预计算空间权重（高斯分布，只与距离有关）
	spaceWeights := make([][]float64, 2*windowRadius+1)
	for i := range spaceWeights {
		spaceWeights[i] = make([]float64, 2*windowRadius+1)
		for j := range spaceWeights[i] {
			dx := i - windowRadius
			dy := j - windowRadius
			// 空间高斯公式：exp(-(dx²+dy²)/(2σ²))
			spaceWeights[i][j] = math.Exp(-(float64(dx*dx + dy*dy)) / (2 * sigmaS * sigmaS))
		}
	}

	// 遍历每个像素计算滤波结果
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			centerY := float64(Y[y][x])
			sumWeight := 0.0
			sumValue := 0.0

			// 遍历窗口内的邻域像素
			for ky := -windowRadius; ky <= windowRadius; ky++ {
				for kx := -windowRadius; kx <= windowRadius; kx++ {
					nx := x + kx
					ny := y + ky
					// 处理边界（超出图像范围的像素忽略）
					if nx < 0 || nx >= width || ny < 0 || ny >= height {
						continue
					}

					// 空间权重（预计算的高斯值）
					sw := spaceWeights[ky+windowRadius][kx+windowRadius]
					// 范围权重（基于亮度差异的高斯值）
					deltaY := float64(Y[ny][nx]) - centerY
					rw := math.Exp(-(deltaY * deltaY) / (2 * sigmaR * sigmaR))
					// 总权重 = 空间权重 × 范围权重
					totalWeight := sw * rw

					sumWeight += totalWeight
					sumValue += totalWeight * float64(Y[ny][nx])
				}
			}

			// 加权平均得到滤波后的Y值
			if sumWeight > 0 {
				result[y][x] = uint8(sumValue / sumWeight)
			} else {
				result[y][x] = Y[y][x] // 权重为0时保留原始值
			}
		}
	}

	return result
}

// YUV分量转图像（合并处理后的Y和原始U/V）
func yuvToImg(Y, U, V [][]uint8) image.Image {
	height := len(Y)
	if height == 0 {
		return nil
	}
	width := len(Y[0])
	bounds := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(bounds)

	// 遍历每个像素转换回RGB
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b := YUVToRGB78(Y[y][x], U[y][x], V[y][x])
			rgba.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return rgba
}

// ========================================================================

// case79 灰度化

// 灰度化原理说明
//均值灰度化：直接取 RGB 三个分量的算术平均值，简单但未考虑人眼对颜色的敏感度差异。公式：Gray = (R + G + B) / 3
//经典灰度化：基于 BT.601 标准，按人眼对红、绿、蓝的敏感度分配权重（绿色敏感度最高）。公式：Gray = 0.299×R + 0.587×G + 0.114×B
//Photoshop 灰度化：采用 Rec.709 标准（高清视频 / 图像常用），权重更贴近现代显示设备的色域特性。公式：Gray = 0.2126×R + 0.7152×G + 0.0722×B

func case79() {
	inputPath := "./test2.jpg"
	// 分别应用三种灰度化方法并保存
	if err := grayscaleImage(inputPath, "output_case79_mean_gray.jpg", MeanGrayscale); err != nil {
		fmt.Printf("均值灰度化失败: %v\n", err)
	} else {
		fmt.Println("均值灰度化完成: mean_gray.jpg")
	}

	if err := grayscaleImage(inputPath, "output_case79_classic_gray.jpg", ClassicGrayscale); err != nil {
		fmt.Printf("经典灰度化失败: %v\n", err)
	} else {
		fmt.Println("经典灰度化完成: classic_gray.jpg")
	}

	if err := grayscaleImage(inputPath, "output_case79_ps_gray.jpg", PhotoshopGrayscale); err != nil {
		fmt.Printf("Photoshop灰度化失败: %v\n", err)
	} else {
		fmt.Println("Photoshop灰度化完成: ps_gray.jpg")
	}
}

// MeanGrayscale 均值灰度化：(R + G + B) / 3
func MeanGrayscale(r, g, b uint8) uint8 {
	// 先转换为int避免uint8溢出（如255+255+255=765，超出uint8范围）
	sum := int(r) + int(g) + int(b)
	return uint8(sum / 3)
}

// ClassicGrayscale 经典灰度化（BT.601标准）：0.299R + 0.587G + 0.114B
func ClassicGrayscale(r, g, b uint8) uint8 {
	// 转换为float64计算，四舍五入后转回uint8
	gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(gray + 0.5) // +0.5实现四舍五入
}

// PhotoshopGrayscale Photoshop灰度化（Rec.709标准）：0.2126R + 0.7152G + 0.0722B
func PhotoshopGrayscale(r, g, b uint8) uint8 {
	gray := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
	return uint8(gray + 0.5) // 四舍五入
}

// 将图像按指定灰度化函数处理并保存
func grayscaleImage(inputPath, outputPath string, grayscaleFunc func(r, g, b uint8) uint8) error {
	// 读取输入图像
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("解码图片失败: %v", err)
	}

	// 创建灰度图像（使用RGBA格式存储，R=G=B=灰度值）
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	grayImg := image.NewRGBA(bounds)

	// 遍历每个像素应用灰度化
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取RGB值（注意：RGBA返回0-65535，需转换为0-255）
			r, g, b, a := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算灰度值
			gray := grayscaleFunc(r8, g8, b8)

			// 灰度图像的RGB分量均为灰度值，保留原始透明度
			grayImg.SetRGBA(x, y, color.RGBA{gray, gray, gray, uint8(a >> 8)})
		}
	}

	// 保存输出图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()

	// 以高质量保存JPEG
	return jpeg.Encode(outputFile, grayImg, &jpeg.Options{Quality: 95})
}

// ========================================================================

// case80 图像阈值化

//阈值化原理说明
//阈值化是将灰度图像（像素值 0-255）转换为二值图像（或特定范围值）的过程，核心是通过预设阈值T划分像素：
//1. 简单阈值化：灰度值 ≥ T → 255（白色），否则 → 0（黑色）
//2. 反阈值化：灰度值 ≥ T → 0（黑色），否则 → 255（白色）
//3. 截断阈值化：灰度值 ≥ T → T（阈值），否则保持原始灰度值

func case80() {

	inputPath := "./test2.jpg"

	// 生成三种阈值化结果
	if err := thresholdImage(inputPath, "output_case80_simple_threshold.jpg", threshold, SimpleThreshold); err != nil {
		fmt.Printf("简单阈值化失败: %v\n", err)
	} else {
		fmt.Println("简单阈值化完成: simple_threshold.jpg")
	}

	if err := thresholdImage(inputPath, "output_case80_inverse_threshold.jpg", threshold, InverseThreshold); err != nil {
		fmt.Printf("反阈值化失败: %v\n", err)
	} else {
		fmt.Println("反阈值化完成: inverse_threshold.jpg")
	}

	if err := thresholdImage(inputPath, "output_case80_truncate_threshold.jpg", threshold, TruncateThreshold); err != nil {
		fmt.Printf("截断阈值化失败: %v\n", err)
	} else {
		fmt.Println("截断阈值化完成: truncate_threshold.jpg")
	}
}

// 灰度化函数（复用之前的经典灰度化，作为阈值化的前置处理）
// ClassicGrayscale 经典灰度化（BT.601标准）：0.299R + 0.587G + 0.114B
func ClassicGrayscale80(r, g, b uint8) uint8 {
	gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(gray + 0.5) // 四舍五入
}

//  阈值化核心函数

// SimpleThreshold 简单阈值化：gray >= threshold → 255，否则 → 0
func SimpleThreshold(gray, threshold uint8) uint8 {
	if gray >= threshold {
		return 255
	}
	return 0
}

// InverseThreshold 反阈值化：gray >= threshold → 0，否则 → 255
func InverseThreshold(gray, threshold uint8) uint8 {
	if gray >= threshold {
		return 0
	}
	return 255
}

// TruncateThreshold 截断阈值化：gray >= threshold → threshold，否则保持原值
func TruncateThreshold(gray, threshold uint8) uint8 {
	if gray >= threshold {
		return threshold
	}
	return gray
}

// 先将彩色图像转为灰度图，再应用阈值化处理
func thresholdImage(inputPath, outputPath string, threshold uint8, thresholdFunc func(gray, threshold uint8) uint8) error {
	// 1. 读取输入图像
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("解码图片失败: %v", err)
	}

	// 2. 转换为灰度图（阈值化需基于灰度图）
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	grayImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			gray := ClassicGrayscale80(r8, g8, b8)
			grayImg.SetRGBA(x, y, color.RGBA{gray, gray, gray, uint8(a >> 8)})
		}
	}

	// 3. 应用阈值化处理
	resultImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 从灰度图中取灰度值（R=G=B，取R即可）
			gray, _, _, a := grayImg.At(x, y).RGBA()
			gray8 := uint8(gray >> 8)

			// 应用阈值化函数
			thresholded := thresholdFunc(gray8, threshold)

			// 阈值化后的像素：RGB分量相同，保留透明度
			resultImg.SetRGBA(x, y, color.RGBA{thresholded, thresholded, thresholded, uint8(a >> 8)})
		}
	}

	// 4. 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()

	return jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
}

// ========================================================================

// case81 图像均值滤波

// 均值滤波原理
//均值滤波是一种线性平滑滤波，核心逻辑是：
//定义一个滑动窗口（通常为奇数大小，如 3x3、5x5），窗口中心为当前处理的像素。
//计算窗口内所有像素的 RGB 分量平均值（每个通道单独计算）。
//用平均值替换窗口中心像素的 RGB 值，Alpha 通道（透明度）保持不变。

func case81() {

	src := getImg("./test2.jpg")

	outputPath := "output_case81.jpg"

	// 解析窗口半径
	radius := 20 // 无效窗口半径: 1,2,3...

	// 应用均值滤波
	filteredImg := MeanFilter(src, radius)

	// 保存输出图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("无法创建输出文件: %v", err))
	}
	defer outputFile.Close()

	// 以高质量保存JPEG（减少压缩失真）
	jpeg.Encode(outputFile, filteredImg, &jpeg.Options{Quality: 95})
	fmt.Printf("均值滤波完成（窗口半径=%d），结果已保存至: %s\n", radius, outputPath)
}

// MeanFilter 对图像应用均值滤波
// input: 输入图像
// radius: 窗口半径（窗口大小为 2*radius + 1，如radius=1对应3x3窗口）
// 返回处理后的图像
func MeanFilter(input image.Image, radius int) image.Image {
	if radius < 1 {
		return input // 半径为0时返回原图
	}

	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds) // 输出图像（RGBA格式）

	// 遍历每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算窗口范围（处理边界：避免超出图像范围）
			minX := max(0, x-radius)
			maxX := min(width-1, x+radius)
			minY := max(0, y-radius)
			maxY := min(height-1, y+radius)

			// 累加窗口内所有像素的RGB值，并计数像素数量
			var rSum, gSum, bSum int
			pixelCount := 0

			for wy := minY; wy <= maxY; wy++ {
				for wx := minX; wx <= maxX; wx++ {
					// 获取窗口内像素的RGBA值（0-65535），转换为0-255
					r, g, b, _ := input.At(wx, wy).RGBA()
					r8 := int(r >> 8)
					g8 := int(g >> 8)
					b8 := int(b >> 8)

					rSum += r8
					gSum += g8
					bSum += b8
					pixelCount++
				}
			}

			// 计算平均值（四舍五入）
			avgR := uint8(rSum/pixelCount + (rSum%pixelCount)*2/pixelCount) // 简单四舍五入
			avgG := uint8(gSum/pixelCount + (gSum%pixelCount)*2/pixelCount)
			avgB := uint8(bSum/pixelCount + (bSum%pixelCount)*2/pixelCount)

			// 保留原始Alpha通道
			_, _, _, a := input.At(x, y).RGBA()
			a8 := uint8(a >> 8)

			// 设置输出像素
			output.SetRGBA(x, y, color.RGBA{avgR, avgG, avgB, a8})
		}
	}

	return output
}

// ========================================================================

// case82 图像USM锐化

// USM 锐化通过增强图像边缘细节实现锐化效果，原理是用原图减去其模糊版本得到边缘掩码，再将掩码叠加回原图，既能增强边缘又不易引入过多噪点，
//是 Photoshop 等软件中常用的锐化算法。

// USM 锐化原理
//生成模糊图像：对原图进行高斯模糊（模糊程度决定锐化的范围，通常用小窗口高斯滤波）。
//计算边缘掩码：边缘掩码 = 原图 - 模糊图像（掩码保留了原图中比模糊图更锐利的细节）。
//叠加掩码增强锐化：锐化图 = 原图 + 数量 × 边缘掩码（“数量” 控制锐化强度，值越大锐化越明显）。

func case82() {
	src := getImg("./test2.jpg")

	outputPath := "output_case82.jpg" // 输出图片路径
	radius := 1                       // 模糊窗口半径（1→3x3窗口，推荐1-3）
	sigma := 1.0                      // 高斯标准差（推荐0.5-2.0，与半径匹配）
	amount := 1.5                     // 锐化强度（推荐0.5-2.0，值越大锐化越明显）

	// 应用USM锐化
	sharpenedImg := USMSharpen(src, radius, sigma, amount)

	// 保存输出图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, sharpenedImg, &jpeg.Options{Quality: 95})
	println("USM锐化完成！")
	println("输出图片:", outputPath)
	println("使用参数: 半径=", radius, " σ=", sigma, " 强度=", amount)
}

// 生成高斯核（权重矩阵）
func generateGaussianKernel82(radius int, sigma float64) [][]float64 {
	size := 2*radius + 1
	kernel := make([][]float64, size)
	sum := 0.0

	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := x - radius
			dy := y - radius
			exponent := -(float64(dx*dx + dy*dy)) / (2 * sigma * sigma)
			kernel[y][x] = math.Exp(exponent) / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}

	// 归一化权重
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}

	return kernel
}

// 高斯模糊（供USM调用）
func gaussianBlur82(input image.Image, radius int, sigma float64) image.Image {
	if radius < 1 || sigma <= 0 {
		return input
	}

	kernel := generateGaussianKernel82(radius, sigma)
	kernelSize := 2*radius + 1

	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var rSum, gSum, bSum float64

			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					imgX := x + (kx - radius)
					imgY := y + (ky - radius)

					if imgX < 0 || imgX >= width || imgY < 0 || imgY >= height {
						continue
					}

					r, g, b, _ := input.At(imgX, imgY).RGBA()
					r8 := float64(r >> 8)
					g8 := float64(g >> 8)
					b8 := float64(b >> 8)

					weight := kernel[ky][kx]
					rSum += weight * r8
					gSum += weight * g8
					bSum += weight * b8
				}
			}

			r := clamp82(rSum)
			g := clamp82(gSum)
			b := clamp82(bSum)
			_, _, _, a := input.At(x, y).RGBA()
			a8 := uint8(a >> 8)

			output.SetRGBA(x, y, color.RGBA{r, g, b, a8})
		}
	}

	return output
}

// 辅助函数：四舍五入并钳位到0-255
func clamp82(val float64) uint8 {
	intVal := int(val + 0.5)
	if intVal < 0 {
		return 0
	}
	if intVal > 255 {
		return 255
	}
	return uint8(intVal)
}

// USMSharpen 应用USM锐化
func USMSharpen(input image.Image, radius int, sigma, amount float64) image.Image {
	if radius < 1 || sigma <= 0 || amount <= 0 {
		return input
	}

	// 1. 生成模糊图像
	blurredImg := gaussianBlur82(input, radius, sigma)

	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds)

	// 2. 计算锐化结果
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 原图像素（0-255）
			rOrig, gOrig, bOrig, aOrig := input.At(x, y).RGBA()
			rO := float64(rOrig >> 8)
			gO := float64(gOrig >> 8)
			bO := float64(bOrig >> 8)

			// 模糊图像素（0-255）
			rBlur, gBlur, bBlur, _ := blurredImg.At(x, y).RGBA()
			rB := float64(rBlur >> 8)
			gB := float64(gBlur >> 8)
			bB := float64(bBlur >> 8)

			// 3. 计算边缘掩码并叠加
			maskR := rO - rB
			maskG := gO - gB
			maskB := bO - bB

			rSharp := rO + amount*maskR
			gSharp := gO + amount*maskG
			bSharp := bO + amount*maskB

			// 钳位并设置像素
			r := clamp82(rSharp)
			g := clamp82(gSharp)
			b := clamp82(bSharp)
			a := uint8(aOrig >> 8)

			output.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}

	return output
}

// ========================================================================

// case83 LUT（Lookup Table，查找表）滤镜是通过预定义的颜色映射关系转换图像色彩的技术，
//广泛用于照片风格化（如复古、电影感、高对比度等效果）。其核心原理是将图像的 RGB 值通过 “查找表” 映射到新的 RGB 值，实现快速色彩转换。

// 实现思路
//1.LUT 结构：使用 3D LUT（立方体结构），维度为N×N×N（常用 17×17×17，兼顾精度与性能），每个节点存储对应输入 RGB 的输出 RGB 值。
//2.索引映射：将图像像素的 RGB 值（0-255）归一化到 0-1 范围，再映射到 LUT 的索引范围（0 到 N-1）。
//3.三线性插值：由于输入值可能落在 LUT 节点之间，通过三线性插值计算精确的映射结果，避免色彩断层。

func case83() {
	// 配置参数（可修改）
	inputPath := "test2.jpg"          // 输入图像路径
	outputPath := "output_case83.jpg" // 输出图像路径

	// 初始化LUT（复古暖色调效果）
	initLUT()

	// 读取输入图像
	file, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图片: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("无法解码图片: " + err.Error())
	}

	// 应用LUT滤镜
	resultImg := applyLUTFilter(img)

	// 保存输出图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
	println("LUT滤镜应用完成！")
	println("输入图片:", inputPath)
	println("输出图片:", outputPath)
}

// 3D LUT配置（17×17×17，兼顾精度与性能）
const lutSize = 17           // LUT维度（N）
const lutScale = lutSize - 1 // 索引最大值（N-1）

// LUT节点：存储输出RGB值（0-255）
type LUTNode struct {
	R, G, B uint8
}

// 全局LUT表（17×17×17）
var lut [lutSize][lutSize][lutSize]LUTNode

// 初始化示例LUT（模拟“复古暖色调”效果：增强红色/黄色，降低蓝色）
func initLUT() {
	for r := 0; r < lutSize; r++ {
		for g := 0; g < lutSize; g++ {
			for b := 0; b < lutSize; b++ {
				// 将LUT索引转换为0-255的RGB值
				rgbR := uint8(r * 255 / lutScale)
				rgbG := uint8(g * 255 / lutScale)
				rgbB := uint8(b * 255 / lutScale)

				// 复古暖色调算法：增强红色，降低蓝色，微调绿色
				newR := clamp83(rgbR + 20) // 红色+20（暖化）
				newG := clamp83(rgbG + 5)  // 绿色+5（柔和）
				newB := clamp83(rgbB - 30) // 蓝色-30（降低冷色）

				lut[r][g][b] = LUTNode{newR, newG, newB}
			}
		}
	}
}

func clamp83(x uint8) uint8 {
	if x > 255 {
		return 255
	}
	return x
}

// 三线性插值：计算输入RGB在LUT中的映射值
func lutLookup(r, g, b uint8) (uint8, uint8, uint8) {
	// 将RGB（0-255）映射到LUT索引（0.0 - lutScale）
	rNorm := float64(r) / 255.0 * lutScale
	gNorm := float64(g) / 255.0 * lutScale
	bNorm := float64(b) / 255.0 * lutScale

	// 分解为整数索引和小数部分（用于插值）
	r0 := int(rNorm)
	g0 := int(gNorm)
	b0 := int(bNorm)
	rf := rNorm - float64(r0)
	gf := gNorm - float64(g0)
	bf := bNorm - float64(b0)

	// 确保索引不越界（极端值保护）
	r1 := min(r0+1, lutScale)
	g1 := min(g0+1, lutScale)
	b1 := min(b0+1, lutScale)

	// 取8个相邻节点（三线性插值需要的立方体顶点）
	c000 := lut[r0][g0][b0]
	c001 := lut[r0][g0][b1]
	c010 := lut[r0][g1][b0]
	c011 := lut[r0][g1][b1]
	c100 := lut[r1][g0][b0]
	c101 := lut[r1][g0][b1]
	c110 := lut[r1][g1][b0]
	c111 := lut[r1][g1][b1]

	// 三线性插值计算（分步骤在三个维度上插值）
	// 1. R维度插值
	v00 := mixNode(c000, c100, rf)
	v01 := mixNode(c001, c101, rf)
	v10 := mixNode(c010, c110, rf)
	v11 := mixNode(c011, c111, rf)

	// 2. G维度插值
	v0 := mixNode(v00, v10, gf)
	v1 := mixNode(v01, v11, gf)

	// 3. B维度插值
	result := mixNode(v0, v1, bf)

	return result.R, result.G, result.B
}

// 辅助函数：在两个LUT节点之间按比例插值
func mixNode(a, b LUTNode, t float64) LUTNode {
	return LUTNode{
		R: uint8(float64(a.R)*(1-t) + float64(b.R)*t + 0.5),
		G: uint8(float64(a.G)*(1-t) + float64(b.G)*t + 0.5),
		B: uint8(float64(a.B)*(1-t) + float64(b.B)*t + 0.5),
	}
}

// 应用LUT滤镜到图像
func applyLUTFilter(input image.Image) image.Image {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds)

	// 遍历每个像素应用LUT
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素RGB值（0-255）
			r, g, b, a := input.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 通过LUT映射新的RGB值
			newR, newG, newB := lutLookup(r8, g8, b8)

			// 设置输出像素（保留透明度）
			output.SetRGBA(x, y, color.RGBA{newR, newG, newB, a8})
		}
	}

	return output
}

// ========================================================================

// case84

// ========================================================================

// case85

// ========================================================================

// case86

// ========================================================================

// case87

// ========================================================================

// case88

// ========================================================================

// case89

// ========================================================================

// case90

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================
// todo gif相关  图片格式转换
// todo 加水印，文字水印，图片水印  多种效果
