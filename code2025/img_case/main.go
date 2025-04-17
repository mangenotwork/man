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
	"os"
	"strings"
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

	case19()

}

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

// todo 调整色相 调整饱和度 调整明暗度 调整色彩平衡 调整亮度 调整对比度 调整锐度 调整色阶 调整曝光度 调整色温 调整色调 锐化  降噪  模糊
