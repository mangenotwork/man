package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// 为图片添加圆角效果
func roundCorners(img image.Image, radius int) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// 创建一个与原图大小相同的透明画布
	draw.Draw(result, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// 计算四个角的圆心坐标
	corners := []struct{ x, y int }{
		{bounds.Min.X + radius - 1, bounds.Min.Y + radius - 1}, // 左上角
		{bounds.Max.X - radius, bounds.Min.Y + radius - 1},     // 右上角
		{bounds.Min.X + radius - 1, bounds.Max.Y - radius},     // 左下角
		{bounds.Max.X - radius, bounds.Max.Y - radius},         // 右下角
	}

	// 复制原图到结果图，同时处理圆角
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 默认在圆角区域外
			inCornerArea := false

			// 检查是否在四个角的区域内
			if x < bounds.Min.X+radius && y < bounds.Min.Y+radius {
				inCornerArea = true
			} else if x >= bounds.Max.X-radius && y < bounds.Min.Y+radius {
				inCornerArea = true
			} else if x < bounds.Min.X+radius && y >= bounds.Max.Y-radius {
				inCornerArea = true
			} else if x >= bounds.Max.X-radius && y >= bounds.Max.Y-radius {
				inCornerArea = true
			}

			// 如果不在任何角的区域内，直接复制像素
			if !inCornerArea {
				result.Set(x, y, img.At(x, y))
				continue
			}

			// 检查是否在圆角范围内
			isInRound := false
			for _, corner := range corners {
				dx := x - corner.x
				dy := y - corner.y
				if dx*dx+dy*dy <= radius*radius {
					isInRound = true
					break
				}
			}

			// 如果在圆角范围内，复制像素
			if isInRound {
				result.Set(x, y, img.At(x, y))
			}
		}
	}

	return result
}

func roundCorners2(img image.Image, radius int, j1, j2, j3, j4 bool) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// 创建一个与原图大小相同的透明画布
	draw.Draw(result, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// 计算四个角的圆心坐标
	corners := make([]struct{ x, y int }, 0)
	//{bounds.Min.X + radius - 1, bounds.Min.Y + radius - 1}, // 左上角
	//{bounds.Max.X - radius, bounds.Min.Y + radius - 1},     // 右上角
	//{bounds.Min.X + radius - 1, bounds.Max.Y - radius},     // 左下角
	//{bounds.Max.X - radius, bounds.Max.Y - radius}, // 右下角

	if j1 {
		corners = append(corners, struct{ x, y int }{bounds.Min.X + radius - 1, bounds.Min.Y + radius - 1}) // 左上角
	}
	if j2 {
		corners = append(corners, struct{ x, y int }{bounds.Max.X - radius, bounds.Min.Y + radius - 1}) // 右上角
	}
	if j3 {
		corners = append(corners, struct{ x, y int }{bounds.Min.X + radius - 1, bounds.Max.Y - radius}) // 左下角
	}
	if j4 {
		corners = append(corners, struct{ x, y int }{bounds.Max.X - radius, bounds.Max.Y - radius}) // 右下角
	}

	// 复制原图到结果图，同时处理圆角
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 默认在圆角区域外
			inCornerArea := false

			// 检查是否在四个角的区域内

			if j1 && x < bounds.Min.X+radius && y < bounds.Min.Y+radius {
				inCornerArea = true
			}
			if j2 && x >= bounds.Max.X-radius && y < bounds.Min.Y+radius {
				inCornerArea = true
			}
			if j3 && x < bounds.Min.X+radius && y >= bounds.Max.Y-radius {
				inCornerArea = true
			}
			if j4 && x >= bounds.Max.X-radius && y >= bounds.Max.Y-radius {
				inCornerArea = true
			}

			//if x < bounds.Min.X+radius && y < bounds.Min.Y+radius {
			//	inCornerArea = true
			//	//} else if x >= bounds.Max.X-radius && y < bounds.Min.Y+radius {
			//	//inCornerArea = true
			//	//} else if x < bounds.Min.X+radius && y >= bounds.Max.Y-radius {
			//	//	inCornerArea = true
			//} else if x >= bounds.Max.X-radius && y >= bounds.Max.Y-radius {
			//	inCornerArea = true
			//}

			// 如果不在任何角的区域内，直接复制像素
			if !inCornerArea {
				result.Set(x, y, img.At(x, y))
				continue
			}

			// 检查是否在圆角范围内
			isInRound := false
			for _, corner := range corners {
				dx := x - corner.x
				dy := y - corner.y
				if dx*dx+dy*dy <= radius*radius {
					isInRound = true
					break
				}
			}

			// 如果在圆角范围内，复制像素
			if isInRound {
				result.Set(x, y, img.At(x, y))
			}
		}
	}

	return result
}

func main() {
	// 打开原始图片文件
	inFile, err := os.Open("input.png")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// 解码图片
	img, _, err := image.Decode(inFile)
	if err != nil {
		panic(err)
	}

	// 设置圆角半径（这里设为图片宽度的10%）
	radius := img.Bounds().Dx() / 10

	// 确保半径不小于1
	if radius < 1 {
		radius = 1
	}

	// 应用圆角效果
	roundedImg := roundCorners2(img, radius, false, true, true, false) // 左上角 右上角 左下角 右下角

	// 创建输出文件
	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// 保存处理后的图片
	err = png.Encode(outFile, roundedImg)
	if err != nil {
		panic(err)
	}

	println("图片圆角处理完成，已保存为 output.png")
}
