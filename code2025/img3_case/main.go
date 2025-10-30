package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

/*

美颜相关算法，学习 《图像视频滤镜与人像美颜美妆算法详解》


*/

func main() {

	case1()

}

// ========================================================================

// case1 双边滤波（Bilateral Filter）是一种经典的保边滤波算法，通过同时考虑空间距离和像素值差异来计算权重，
//在平滑噪声的同时能有效保留图像边缘（传统高斯滤波仅考虑空间距离，会模糊边缘）

// 双边滤波原理
//双边滤波的核心是为每个像素的邻域像素计算 “双重权重”：
//1. 空间权重（Spatial Weight）：基于像素间的欧氏距离，距离越近权重越大（高斯分布）。
//2. 范围权重（Range Weight）：基于像素值的差异，差异越小权重越大（高斯分布）。

func case1() {
	// 配置参数（可修改）
	inputPath := "test4.jpg"         // 输入图像路径
	outputPath := "output_case1.jpg" // 输出图像路径

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

	// 应用双边滤波
	filteredImg := bilateralFilter(img)

	// 保存输出图像
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, filteredImg, &jpeg.Options{Quality: 95})
	println("双边滤波完成！")
	println("输入图片:", inputPath)
	println("输出图片:", outputPath)
	println("使用参数: 窗口半径=", windowRadius, " σ_s=", sigmaS, " σ_r=", sigmaR)
}

// 双边滤波参数（可调整）
const (
	windowRadius = 2    // 滤波窗口半径（窗口大小为2*radius+1，如2→5x5窗口）
	sigmaS       = 1.5  // 空间高斯标准差（控制空间权重衰减速度）
	sigmaR       = 30.0 // 范围高斯标准差（控制像素差异权重衰减，值越小边缘保留越严格）
)

// 双边滤波核心函数
// input: 输入图像
// 返回滤波后的图像
func bilateralFilter(input image.Image) image.Image {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds) // 输出图像（RGBA格式）

	// 预计算空间权重（仅与距离有关，可复用）
	spaceWeights := precomputeSpaceWeights(windowRadius, sigmaS)

	// 遍历每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取中心像素的RGB值（0-255）
			curR, curG, curB, curA := input.At(x, y).RGBA()
			centerR := float64(curR >> 8)
			centerG := float64(curG >> 8)
			centerB := float64(curB >> 8)

			var sumR, sumG, sumB, totalWeight float64

			// 遍历窗口内的邻域像素
			for ky := -windowRadius; ky <= windowRadius; ky++ {
				for kx := -windowRadius; kx <= windowRadius; kx++ {
					// 计算邻域像素坐标（处理边界，避免越界）
					neighX := x + kx
					neighY := y + ky
					if neighX < 0 || neighX >= width || neighY < 0 || neighY >= height {
						continue // 忽略图像外的像素
					}

					// 获取邻域像素的RGB值（0-255）
					neighR, neighG, neighB, _ := input.At(neighX, neighY).RGBA()
					r := float64(neighR >> 8)
					g := float64(neighG >> 8)
					b := float64(neighB >> 8)

					// 1. 获取预计算的空间权重
					sw := spaceWeights[ky+windowRadius][kx+windowRadius]

					// 2. 计算范围权重（基于RGB三通道的差异，取平均差异更稳定）
					dr := r - centerR
					dg := g - centerG
					db := b - centerB
					delta := (dr*dr + dg*dg + db*db) / 3.0 // 三通道平均差异
					rw := math.Exp(-delta / (2 * sigmaR * sigmaR))

					// 总权重 = 空间权重 × 范围权重
					weight := sw * rw

					// 累加加权像素值和总权重
					sumR += weight * r
					sumG += weight * g
					sumB += weight * b
					totalWeight += weight
				}
			}

			// 计算加权平均（避免除零）
			if totalWeight > 0 {
				sumR /= totalWeight
				sumG /= totalWeight
				sumB /= totalWeight
			}

			// 转换为uint8并钳位到0-255
			r := clamp(sumR)
			g := clamp(sumG)
			b := clamp(sumB)
			a := uint8(curA >> 8) // 保留原始透明度

			output.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}

	return output
}

// 预计算空间权重（基于窗口半径和sigmaS）
func precomputeSpaceWeights(radius int, sigmaS float64) [][]float64 {
	size := 2*radius + 1
	weights := make([][]float64, size)
	for i := 0; i < size; i++ {
		weights[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			// 计算相对于窗口中心的偏移量
			dx := i - radius
			dy := j - radius
			// 空间高斯函数：exp(-(dx²+dy²)/(2σ_s²))
			weights[i][j] = math.Exp(-(float64(dx*dx + dy*dy)) / (2 * sigmaS * sigmaS))
		}
	}
	return weights
}

// 辅助函数：将浮点数转换为uint8并钳位到0-255
func clamp(val float64) uint8 {
	intVal := int(val + 0.5) // 四舍五入
	if intVal < 0 {
		return 0
	}
	if intVal > 255 {
		return 255
	}
	return uint8(intVal)
}

// ========================================================================

// case2 Surface Blur滤波算法

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================

// ========================================================================
