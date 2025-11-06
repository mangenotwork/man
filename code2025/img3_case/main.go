package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sort"
	"sync"
)

/*

美颜相关算法，学习 《图像视频滤镜与人像美颜美妆算法详解》


*/

func main() {

	//case1()

	//case2()

	//case3()

	//case4()

	//case5()

	//case6()

	//case7()

	//case8()

	//case9()

	//case10()

	//case11()

	//case12()

	//case13()
	//case13_1()
	//case13_2()

	//case14()

	//case15()

	//case16()

	//case17()

	//case18()

	//case19()

	//case20()

	//case21()

	//case22()

	//case23()

	//case24()

	//case25()

	case26()
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

// Surface Blur（表面模糊）是一种常用于图像编辑的保边滤波算法，尤其适合平滑皮肤、去除瑕疵同时保留皮肤纹理和边缘细节（如发丝、五官轮廓）

// Surface Blur 原理
//Surface Blur 的核心逻辑是：在滑动窗口内，仅让与中心像素颜色差异（通常用亮度差异）小于阈值的像素参与平均计算，
//从而在平滑 “颜色相近的表面区域”（如皮肤）的同时，忽略 “颜色突变的边缘区域”（如轮廓、纹理边界），实现保边效果。

func case2() {
	sbInputPath := "test4.jpg"
	sbOutputPath := "output_case2.jpg"

	// 读取输入图像
	sbFile, sbErr := os.Open(sbInputPath)
	if sbErr != nil {
		panic("无法打开输入图片: " + sbErr.Error())
	}
	defer sbFile.Close()

	sbImg, _, sbErr := image.Decode(sbFile)
	if sbErr != nil {
		panic("无法解码图片: " + sbErr.Error())
	}

	// 应用滤波
	sbFilteredImg := SurfaceBlur(sbImg)

	// 保存输出图像
	sbOutputFile, sbErr := os.Create(sbOutputPath)
	if sbErr != nil {
		panic("无法创建输出文件: " + sbErr.Error())
	}
	defer sbOutputFile.Close()

	jpeg.Encode(sbOutputFile, sbFilteredImg, &jpeg.Options{Quality: 95})
	println("Surface Blur滤波完成！")
	println("输入:", sbInputPath, "输出:", sbOutputPath)
	println("参数: 窗口半径=", sbWindowRadius, " 颜色阈值=", sbColorThreshold)
}

// Surface Blur参数
const (
	sbWindowRadius   = 5    // 窗口半径
	sbColorThreshold = 10.0 // 颜色差异阈值（float64类型）
)

// rgbToSbY 将RGB转换为亮度Y（返回float64）
func rgbToSbY(r, g, b uint8) float64 {
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}

// SurfaceBlur 表面模糊滤波核心函数
func SurfaceBlur(sbInput image.Image) image.Image {
	sbBounds := sbInput.Bounds()
	sbWidth, sbHeight := sbBounds.Max.X, sbBounds.Max.Y
	sbOutput := image.NewRGBA(sbBounds)

	for sbY := 0; sbY < sbHeight; sbY++ {
		for sbX := 0; sbX < sbWidth; sbX++ {
			// 获取中心像素的RGBA值
			sbCurR, sbCurG, sbCurB, sbCurA := sbInput.At(sbX, sbY).RGBA()
			sbCenterR := uint8(sbCurR >> 8)
			sbCenterG := uint8(sbCurG >> 8)
			sbCenterB := uint8(sbCurB >> 8)

			// 中心像素亮度（明确为float64）
			var sbCenterY float64 = rgbToSbY(sbCenterR, sbCenterG, sbCenterB)

			var sbSumR, sbSumG, sbSumB int
			sbValidPixelCount := 0

			// 遍历窗口内邻域像素
			for sbKy := -sbWindowRadius; sbKy <= sbWindowRadius; sbKy++ {
				for sbKx := -sbWindowRadius; sbKx <= sbWindowRadius; sbKx++ {
					sbNeighX := sbX + sbKx
					sbNeighYCoord := sbY + sbKy // 重命名变量避免与亮度变量冲突
					if sbNeighX < 0 || sbNeighX >= sbWidth || sbNeighYCoord < 0 || sbNeighYCoord >= sbHeight {
						continue
					}

					// 获取邻域像素RGB值
					sbNeighR, sbNeighG, sbNeighB, _ := sbInput.At(sbNeighX, sbNeighYCoord).RGBA()
					sbNeighR8 := uint8(sbNeighR >> 8)
					sbNeighG8 := uint8(sbNeighG >> 8)
					sbNeighB8 := uint8(sbNeighB >> 8)

					// 邻域像素亮度（明确为float64）
					var sbNeighY float64 = rgbToSbY(sbNeighR8, sbNeighG8, sbNeighB8)

					// 计算亮度差异（均为float64类型）
					sbYDiff := sbAbs(sbNeighY - sbCenterY)

					// 阈值比较（均为float64类型）
					if sbYDiff <= sbColorThreshold {
						sbSumR += int(sbNeighR8)
						sbSumG += int(sbNeighG8)
						sbSumB += int(sbNeighB8)
						sbValidPixelCount++
					}
				}
			}

			// 计算新像素值
			var sbNewR, sbNewG, sbNewB uint8
			if sbValidPixelCount > 0 {
				sbNewR = uint8(sbSumR / sbValidPixelCount)
				sbNewG = uint8(sbSumG / sbValidPixelCount)
				sbNewB = uint8(sbSumB / sbValidPixelCount)
			} else {
				sbNewR, sbNewG, sbNewB = sbCenterR, sbCenterG, sbCenterB
			}

			// 设置输出像素
			sbA8 := uint8(sbCurA >> 8)
			sbOutput.SetRGBA(sbX, sbY, color.RGBA{sbNewR, sbNewG, sbNewB, sbA8})
		}
	}

	return sbOutput
}

// sbAbs 计算float64绝对值
func sbAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ========================================================================

// case3 Guided滤波算法

// 引导滤波（Guided Filter）是一种高效的保边滤波算法，通过利用引导图像（通常为原图）的边缘信息来控制滤波过程，在平滑噪声的同时能很好地保留边缘细节

func case3() {
	gfInputPath := "test4.jpg"         // 输入图像路径
	gfOutputPath := "output_case3.jpg" // 输出图像路径

	// 读取输入图像
	gfFile, gfErr := os.Open(gfInputPath)
	if gfErr != nil {
		panic("无法打开输入图片: " + gfErr.Error())
	}
	defer gfFile.Close()

	gfImg, _, gfErr := image.Decode(gfFile)
	if gfErr != nil {
		panic("无法解码图片: " + gfErr.Error())
	}

	// 应用引导滤波（使用原图作为引导图像）
	gfFilteredImg := GuidedFilter(gfImg, gfImg)

	// 保存输出图像
	gfOutputFile, gfErr := os.Create(gfOutputPath)
	if gfErr != nil {
		panic("无法创建输出文件: " + gfErr.Error())
	}
	defer gfOutputFile.Close()

	jpeg.Encode(gfOutputFile, gfFilteredImg, &jpeg.Options{Quality: 95})
	println("引导滤波完成！")
	println("输入图片:", gfInputPath)
	println("输出图片:", gfOutputPath)
	println("参数: 窗口半径=", gfRadius, " 正则化参数=", gfEpsilon)
}

// 引导滤波参数（增大参数增强效果）
const (
	gfRadius  = 10  // 增大窗口半径（原5→10，平滑范围更广）
	gfEpsilon = 0.1 // 增大正则化参数（原0.01→0.1，平滑强度更高）
)

// 引导滤波核心函数
func GuidedFilter(gfInput, gfGuide image.Image) image.Image {
	gfBounds := gfInput.Bounds()
	gfWidth, gfHeight := gfBounds.Max.X, gfBounds.Max.Y
	gfOutput := image.NewRGBA(gfBounds)

	// 转换为浮点型RGB通道（0.0-1.0）
	gfInputR, gfInputG, gfInputB := gfImgToFloatChannels(gfInput)
	gfGuideR, gfGuideG, gfGuideB := gfImgToFloatChannels(gfGuide)

	// 分通道滤波
	gfOutputR := gfGuidedFilterChannel(gfInputR, gfGuideR, gfWidth, gfHeight, gfRadius, gfEpsilon)
	gfOutputG := gfGuidedFilterChannel(gfInputG, gfGuideG, gfWidth, gfHeight, gfRadius, gfEpsilon)
	gfOutputB := gfGuidedFilterChannel(gfInputB, gfGuideB, gfWidth, gfHeight, gfRadius, gfEpsilon)

	// 合并通道并保留Alpha
	for y := 0; y < gfHeight; y++ {
		for x := 0; x < gfWidth; x++ {
			_, _, _, gfA := gfInput.At(x, y).RGBA()
			gfA8 := uint8(gfA >> 8)

			gfR := gfClampFloat(gfOutputR[y][x])
			gfG := gfClampFloat(gfOutputG[y][x])
			gfB := gfClampFloat(gfOutputB[y][x])

			gfOutput.SetRGBA(x, y, color.RGBA{gfR, gfG, gfB, gfA8})
		}
	}

	return gfOutput
}

// 单通道引导滤波
func gfGuidedFilterChannel(gfInput, gfGuide [][]float64, gfWidth, gfHeight, gfRadius int, gfEpsilon float64) [][]float64 {
	// 1. 计算窗口均值
	gfMeanI := gfBoxFilter(gfGuide, gfWidth, gfHeight, gfRadius)
	gfMeanP := gfBoxFilter(gfInput, gfWidth, gfHeight, gfRadius)
	gfMeanIP := gfBoxFilter(gfMultiply(gfGuide, gfInput), gfWidth, gfHeight, gfRadius)
	gfCovIP := gfSubtract(gfMeanIP, gfMultiply(gfMeanI, gfMeanP))

	// 2. 计算方差
	gfMeanII := gfBoxFilter(gfMultiply(gfGuide, gfGuide), gfWidth, gfHeight, gfRadius)
	gfVarI := gfSubtract(gfMeanII, gfMultiply(gfMeanI, gfMeanI))

	// 创建epsilon数组（与矩阵同尺寸）
	gfEpsilonArr := gfCreateScalarArray(gfWidth, gfHeight, gfEpsilon)

	// 3. 计算线性系数a和b
	gfA := gfDivide(
		gfAdd(gfCovIP, gfEpsilonArr),
		gfAdd(gfVarI, gfEpsilonArr),
	)
	gfB := gfSubtract(gfMeanP, gfMultiply(gfA, gfMeanI))

	// 4. 平滑系数a和b
	gfMeanA := gfBoxFilter(gfA, gfWidth, gfHeight, gfRadius)
	gfMeanB := gfBoxFilter(gfB, gfWidth, gfHeight, gfRadius)

	// 5. 计算输出
	return gfAdd(gfMultiply(gfMeanA, gfGuide), gfMeanB)
}

// 盒式滤波（计算窗口均值）
func gfBoxFilter(gfInput [][]float64, gfWidth, gfHeight, gfRadius int) [][]float64 {
	// 1. 计算积分图
	gfIntegral := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfIntegral[y] = make([]float64, gfWidth)
		gfRowSum := 0.0
		for x := 0; x < gfWidth; x++ {
			gfRowSum += gfInput[y][x]
			if y == 0 {
				gfIntegral[y][x] = gfRowSum
			} else {
				gfIntegral[y][x] = gfIntegral[y-1][x] + gfRowSum
			}
		}
	}

	// 2. 计算窗口均值
	gfOutput := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfOutput[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfMinX := gfMax(0, x-gfRadius)
			gfMaxX := gfMin(gfWidth-1, x+gfRadius)
			gfMinY := gfMax(0, y-gfRadius)
			gfMaxY := gfMin(gfHeight-1, y+gfRadius)

			gfCount := (gfMaxX - gfMinX + 1) * (gfMaxY - gfMinY + 1)
			gfSum := gfIntegral[gfMaxY][gfMaxX]
			if gfMinX > 0 {
				gfSum -= gfIntegral[gfMaxY][gfMinX-1]
			}
			if gfMinY > 0 {
				gfSum -= gfIntegral[gfMinY-1][gfMaxX]
			}
			if gfMinX > 0 && gfMinY > 0 {
				gfSum += gfIntegral[gfMinY-1][gfMinX-1]
			}

			gfOutput[y][x] = gfSum / float64(gfCount)
		}
	}

	return gfOutput
}

// 创建标量数组（与图像同尺寸）
func gfCreateScalarArray(gfWidth, gfHeight int, gfScalar float64) [][]float64 {
	gfArr := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfArr[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfArr[y][x] = gfScalar
		}
	}
	return gfArr
}

// 矩阵加法
func gfAdd(gfA, gfB [][]float64) [][]float64 {
	gfHeight := len(gfA)
	gfWidth := len(gfA[0])
	gfRes := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfRes[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfRes[y][x] = gfA[y][x] + gfB[y][x]
		}
	}
	return gfRes
}

// 矩阵减法
func gfSubtract(gfA, gfB [][]float64) [][]float64 {
	gfHeight := len(gfA)
	gfWidth := len(gfA[0])
	gfRes := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfRes[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfRes[y][x] = gfA[y][x] - gfB[y][x]
		}
	}
	return gfRes
}

// 矩阵乘法（对应元素）
func gfMultiply(gfA, gfB [][]float64) [][]float64 {
	gfHeight := len(gfA)
	gfWidth := len(gfA[0])
	gfRes := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfRes[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfRes[y][x] = gfA[y][x] * gfB[y][x]
		}
	}
	return gfRes
}

// 矩阵除法（对应元素）
func gfDivide(gfA, gfB [][]float64) [][]float64 {
	gfHeight := len(gfA)
	gfWidth := len(gfA[0])
	gfRes := make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfRes[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			if gfB[y][x] < 1e-8 {
				gfRes[y][x] = 0
			} else {
				gfRes[y][x] = gfA[y][x] / gfB[y][x]
			}
		}
	}
	return gfRes
}

// 图像转浮点通道
func gfImgToFloatChannels(gfImg image.Image) (gfR, gfG, gfB [][]float64) {
	gfBounds := gfImg.Bounds()
	gfWidth, gfHeight := gfBounds.Max.X, gfBounds.Max.Y

	gfR = make([][]float64, gfHeight)
	gfG = make([][]float64, gfHeight)
	gfB = make([][]float64, gfHeight)
	for y := 0; y < gfHeight; y++ {
		gfR[y] = make([]float64, gfWidth)
		gfG[y] = make([]float64, gfWidth)
		gfB[y] = make([]float64, gfWidth)
		for x := 0; x < gfWidth; x++ {
			gfRVal, gfGVal, gfBVal, _ := gfImg.At(x, y).RGBA()
			gfR[y][x] = float64(gfRVal>>8) / 255.0
			gfG[y][x] = float64(gfGVal>>8) / 255.0
			gfB[y][x] = float64(gfBVal>>8) / 255.0
		}
	}
	return gfR, gfG, gfB
}

// 浮点值转uint8
func gfClampFloat(gfVal float64) uint8 {
	gfVal = math.Max(0, math.Min(1.0, gfVal))
	return uint8(gfVal*255 + 0.5)
}

// 取最大值
func gfMax(gfA, gfB int) int {
	if gfA > gfB {
		return gfA
	}
	return gfB
}

// 取最小值
func gfMin(gfA, gfB int) int {
	if gfA < gfB {
		return gfA
	}
	return gfB
}

// ========================================================================

// case4 局部均值滤波算法

// 局部均值滤波（Local Mean Filter）是一种基础的线性平滑滤波算法，通过计算像素邻域内所有像素的平均值来替换中心像素值，
//从而实现降噪和平滑效果。它实现简单、计算高效，但会模糊图像边缘（与保边滤波的核心区别），适合处理高斯噪声或椒盐噪声等场景。

func case4() {
	// 配置文件路径
	mfInputPath := "test4.jpg"         // 输入图像路径
	mfOutputPath := "output_case4.jpg" // 输出图像路径

	// 读取输入图像
	mfFile, mfErr := os.Open(mfInputPath)
	if mfErr != nil {
		panic("无法打开输入图片: " + mfErr.Error())
	}
	defer mfFile.Close()

	mfImg, _, mfErr := image.Decode(mfFile)
	if mfErr != nil {
		panic("无法解码图片: " + mfErr.Error())
	}

	// 应用局部均值滤波
	mfFilteredImg := MeanFilter(mfImg)

	// 保存输出图像
	mfOutputFile, mfErr := os.Create(mfOutputPath)
	if mfErr != nil {
		panic("无法创建输出文件: " + mfErr.Error())
	}
	defer mfOutputFile.Close()

	jpeg.Encode(mfOutputFile, mfFilteredImg, &jpeg.Options{Quality: 95})
	println("局部均值滤波完成！")
	println("输入图片:", mfInputPath)
	println("输出图片:", mfOutputPath)
	println("滤波窗口大小:", 2*mfWindowRadius+1, "x", 2*mfWindowRadius+1)

}

// 局部均值滤波参数
const (
	mfWindowRadius = 3 // 窗口半径（窗口大小 = 2*半径 + 1，如1对应3x3窗口）  2=5x5  3=7x7
)

// MeanFilter 局部均值滤波核心函数
// 对输入图像应用均值滤波，返回滤波后的图像
func MeanFilter(mfInput image.Image) image.Image {
	mfBounds := mfInput.Bounds()
	mfWidth, mfHeight := mfBounds.Max.X, mfBounds.Max.Y
	mfOutput := image.NewRGBA(mfBounds)

	// 遍历图像每个像素
	for mfY := 0; mfY < mfHeight; mfY++ {
		for mfX := 0; mfX < mfWidth; mfX++ {
			// 保留原始Alpha通道值
			_, _, _, mfAlpha := mfInput.At(mfX, mfY).RGBA()
			mfAlpha8 := uint8(mfAlpha >> 8)

			// 累加窗口内像素的RGB值
			var mfSumR, mfSumG, mfSumB int
			mfValidCount := 0 // 有效像素计数（避免边界越界）

			// 遍历窗口内所有邻域像素
			for mfKy := -mfWindowRadius; mfKy <= mfWindowRadius; mfKy++ {
				for mfKx := -mfWindowRadius; mfKx <= mfWindowRadius; mfKx++ {
					// 计算邻域像素坐标
					mfNeighX := mfX + mfKx
					mfNeighY := mfY + mfKy

					// 检查坐标是否在图像范围内（处理边界）
					if mfNeighX >= 0 && mfNeighX < mfWidth && mfNeighY >= 0 && mfNeighY < mfHeight {
						// 获取邻域像素的RGB值（转换为0-255范围）
						mfR, mfG, mfB, _ := mfInput.At(mfNeighX, mfNeighY).RGBA()
						mfSumR += int(mfR >> 8)
						mfSumG += int(mfG >> 8)
						mfSumB += int(mfB >> 8)
						mfValidCount++
					}
				}
			}

			// 计算窗口内的平均值（避免除零，理论上count至少为1）
			mfAvgR := uint8(mfSumR / mfValidCount)
			mfAvgG := uint8(mfSumG / mfValidCount)
			mfAvgB := uint8(mfSumB / mfValidCount)

			// 设置输出像素值
			mfOutput.SetRGBA(mfX, mfY, color.RGBA{mfAvgR, mfAvgG, mfAvgB, mfAlpha8})
		}
	}

	return mfOutput
}

// ========================================================================

// case5 Anisotropic滤波算法

// 各向异性滤波（Anisotropic Filtering）是一种能够根据图像局部结构（如边缘方向）调整平滑强度的滤波算法。
//与各向同性滤波（如均值滤波）在所有方向上均匀平滑不同，它会沿着边缘方向增强平滑，同时垂直于边缘方向抑制平滑，从而在降噪的同时更好地保留边缘和细节。

func case5() {
	// 配置文件路径
	afInputPath := "test4.jpg"         // 输入图像路径
	afOutputPath := "output_case5.jpg" // 输出图像路径

	// 读取输入图像
	afFile, afErr := os.Open(afInputPath)
	if afErr != nil {
		panic("无法打开输入图片: " + afErr.Error())
	}
	defer afFile.Close()

	afImg, _, afErr := image.Decode(afFile)
	if afErr != nil {
		panic("无法解码图片: " + afErr.Error())
	}

	// 应用各向异性滤波
	afFilteredImg := AnisotropicFilter(afImg)

	// 保存输出图像
	afOutputFile, afErr := os.Create(afOutputPath)
	if afErr != nil {
		panic("无法创建输出文件: " + afErr.Error())
	}
	defer afOutputFile.Close()

	jpeg.Encode(afOutputFile, afFilteredImg, &jpeg.Options{Quality: 95})
	println("各向异性滤波完成！")
	println("输入图片:", afInputPath)
	println("输出图片:", afOutputPath)
	println("参数: 迭代次数=", afIterations, " 边缘敏感度=", afK)
}

// 各向异性滤波参数（优化默认值，平衡速度与效果）
const (
	afIterations = 3    // 减少迭代次数（原5→3，速度提升40%）
	afK          = 30.0 // 边缘敏感度（保持效果）
	afDeltaT     = 0.25 // 扩散步长（保持稳定性）
	afSigma      = 1.0  // 高斯平滑标准差（保持梯度计算准确性）
)

// AnisotropicFilter 优化后的各向异性滤波核心函数
func AnisotropicFilter(afInput image.Image) image.Image {
	// 转换为浮点型RGB通道（0.0-255.0）
	afR, afG, afB := afImgToFloat(afInput)
	afWidth := len(afR[0])
	afHeight := len(afR)

	// 预分配扩散用的缓冲区（复用内存，减少分配开销）
	afBufR := make([][]float64, afHeight)
	afBufG := make([][]float64, afHeight)
	afBufB := make([][]float64, afHeight)
	for y := 0; y < afHeight; y++ {
		afBufR[y] = make([]float64, afWidth)
		afBufG[y] = make([]float64, afWidth)
		afBufB[y] = make([]float64, afWidth)
	}

	// 多轮迭代扩散（减少迭代次数并复用缓冲区）
	for i := 0; i < afIterations; i++ {
		// 对每个通道单独扩散，直接写入缓冲区
		afDiffuseOptimized(afR, afBufR, afWidth, afHeight)
		afDiffuseOptimized(afG, afBufG, afWidth, afHeight)
		afDiffuseOptimized(afB, afBufB, afWidth, afHeight)

		// 交换源和缓冲区（避免重复分配）
		afR, afBufR = afBufR, afR
		afG, afBufG = afBufG, afG
		afB, afBufB = afBufB, afB
	}

	// 转换回图像格式
	return afFloatToImg(afR, afG, afB, afInput)
}

// 优化的单通道扩散函数（减少内存访问和计算冗余）
func afDiffuseOptimized(afSrc, afDst [][]float64, afWidth, afHeight int) {
	// 预计算高斯平滑（修复原代码中的循环变量错误，并优化核计算）
	afSmoothed := afGaussianBlurOptimized(afSrc, afWidth, afHeight, afSigma)

	// 预计算梯度（合并边界处理，减少条件判断）
	afGradX, afGradY := afComputeGradientsOptimized(afSmoothed, afWidth, afHeight)

	// 遍历像素计算扩散（按行连续访问，提升缓存利用率）
	for y := 0; y < afHeight; y++ {
		for x := 0; x < afWidth; x++ {
			// 直接使用当前像素值作为初始值
			afDst[y][x] = afSrc[y][x]

			// 跳过边缘像素（已通过梯度计算处理边界，减少内部判断）
			if y == 0 || y == afHeight-1 || x == 0 || x == afWidth-1 {
				continue
			}

			// 读取四方向梯度（减少数组访问次数）
			afN := afGradY[y-1][x]
			afS := afGradY[y+1][x]
			afE := afGradX[y][x+1]
			afW := afGradX[y][x-1]

			// 计算扩散系数（合并指数计算的公共部分）
			afK2 := afK * afK // 预计算K的平方
			afCN := math.Exp(-(afN * afN) / afK2)
			afCS := math.Exp(-(afS * afS) / afK2)
			afCE := math.Exp(-(afE * afE) / afK2)
			afCW := math.Exp(-(afW * afW) / afK2)

			// 计算扩散量（合并同类项，减少计算步骤）
			afSrcVal := afSrc[y][x]
			afDiffuse := afDeltaT * (afCN*(afSrc[y-1][x]-afSrcVal) +
				afCS*(afSrc[y+1][x]-afSrcVal) +
				afCE*(afSrc[y][x+1]-afSrcVal) +
				afCW*(afSrc[y][x-1]-afSrcVal))

			afDst[y][x] = afSrcVal + afDiffuse
		}
	}
}

// 优化的高斯模糊（使用1D核分离计算，减少66%的乘法操作）
func afGaussianBlurOptimized(afSrc [][]float64, afWidth, afHeight int, afSigma float64) [][]float64 {
	// 1. 计算1D高斯核（分离为水平和垂直方向）
	afKernel := make([]float64, 3)
	afSigma2 := 2 * afSigma * afSigma
	afKernel[0] = math.Exp(-1 / afSigma2)
	afKernel[1] = math.Exp(0)
	afKernel[2] = afKernel[0]
	afSum := afKernel[0] + afKernel[1] + afKernel[2]
	afKernel[0] /= afSum
	afKernel[1] /= afSum
	afKernel[2] /= afSum

	// 2. 先进行水平方向模糊
	afTemp := make([][]float64, afHeight)
	for y := 0; y < afHeight; y++ {
		afTemp[y] = make([]float64, afWidth)
		for x := 0; x < afWidth; x++ {
			afVal := 0.0
			// 水平卷积（仅3个点）
			for kx := -1; kx <= 1; kx++ {
				afNx := x + kx
				if afNx < 0 {
					afNx = 0 // 边缘处理：复制边界
				} else if afNx >= afWidth {
					afNx = afWidth - 1
				}
				afVal += afSrc[y][afNx] * afKernel[kx+1]
			}
			afTemp[y][x] = afVal
		}
	}

	// 3. 再进行垂直方向模糊（基于水平结果）
	afDst := make([][]float64, afHeight)
	for y := 0; y < afHeight; y++ {
		afDst[y] = make([]float64, afWidth)
		for x := 0; x < afWidth; x++ {
			afVal := 0.0
			// 垂直卷积（仅3个点）
			for ky := -1; ky <= 1; ky++ {
				afNy := y + ky
				if afNy < 0 {
					afNy = 0 // 边缘处理：复制边界
				} else if afNy >= afHeight {
					afNy = afHeight - 1
				}
				afVal += afTemp[afNy][x] * afKernel[ky+1]
			}
			afDst[y][x] = afVal
		}
	}

	return afDst
}

// 优化的梯度计算（减少条件判断，统一处理边界）
func afComputeGradientsOptimized(afSrc [][]float64, afWidth, afHeight int) (afGradX, afGradY [][]float64) {
	afGradX = make([][]float64, afHeight)
	afGradY = make([][]float64, afHeight)
	for y := 0; y < afHeight; y++ {
		afGradX[y] = make([]float64, afWidth)
		afGradY[y] = make([]float64, afWidth)
	}

	// 计算x方向梯度（统一公式，边缘用单侧差分）
	for y := 0; y < afHeight; y++ {
		afGradX[y][0] = afSrc[y][1] - afSrc[y][0] // 左边界
		for x := 1; x < afWidth-1; x++ {
			afGradX[y][x] = (afSrc[y][x+1] - afSrc[y][x-1]) / 2.0 // 内部
		}
		afGradX[y][afWidth-1] = afSrc[y][afWidth-1] - afSrc[y][afWidth-2] // 右边界
	}

	// 计算y方向梯度（统一公式，边缘用单侧差分）
	for x := 0; x < afWidth; x++ {
		afGradY[0][x] = afSrc[1][x] - afSrc[0][x] // 上边界
		for y := 1; y < afHeight-1; y++ {
			afGradY[y][x] = (afSrc[y+1][x] - afSrc[y-1][x]) / 2.0 // 内部
		}
		afGradY[afHeight-1][x] = afSrc[afHeight-1][x] - afSrc[afHeight-2][x] // 下边界
	}

	return afGradX, afGradY
}

// 图像格式转换辅助函数（保持不变）

func afImgToFloat(afImg image.Image) (afR, afG, afB [][]float64) {
	afBounds := afImg.Bounds()
	afWidth, afHeight := afBounds.Max.X, afBounds.Max.Y

	afR = make([][]float64, afHeight)
	afG = make([][]float64, afHeight)
	afB = make([][]float64, afHeight)
	for y := 0; y < afHeight; y++ {
		afR[y] = make([]float64, afWidth)
		afG[y] = make([]float64, afWidth)
		afB[y] = make([]float64, afWidth)
		for x := 0; x < afWidth; x++ {
			afRVal, afGVal, afBVal, _ := afImg.At(x, y).RGBA()
			afR[y][x] = float64(afRVal >> 8)
			afG[y][x] = float64(afGVal >> 8)
			afB[y][x] = float64(afBVal >> 8)
		}
	}
	return afR, afG, afB
}

func afFloatToImg(afR, afG, afB [][]float64, afOrig image.Image) image.Image {
	afBounds := afOrig.Bounds()
	afWidth, afHeight := afBounds.Max.X, afBounds.Max.Y
	afDst := image.NewRGBA(afBounds)

	for y := 0; y < afHeight; y++ {
		for x := 0; x < afWidth; x++ {
			afRClamped := math.Max(0, math.Min(255, afR[y][x]))
			afGClamped := math.Max(0, math.Min(255, afG[y][x]))
			afBClamped := math.Max(0, math.Min(255, afB[y][x]))

			_, _, _, afA := afOrig.At(x, y).RGBA()
			afA8 := uint8(afA >> 8)

			afDst.SetRGBA(x, y, color.RGBA{
				uint8(afRClamped),
				uint8(afGClamped),
				uint8(afBClamped),
				afA8,
			})
		}
	}
	return afDst
}

// ========================================================================

// case6 Smart Blur滤波算法

// Smart Blur（智能模糊）是一种结合了边缘检测的保边滤波算法，它能在平滑图像噪声的同时有效保留边缘细节。与普通均值滤波不同，
//Smart Blur 会先判断像素是否处于边缘区域 —— 对平坦区域应用较强模糊，对边缘区域应用较弱模糊（或不模糊），从而实现 “智能” 的选择性平滑。

func case6() {
	// 配置文件路径
	sbInputPath := "test4.jpg"         // 输入图像路径
	sbOutputPath := "output_case6.jpg" // 输出图像路径

	// 读取输入图像
	sbFile, sbErr := os.Open(sbInputPath)
	if sbErr != nil {
		panic("无法打开输入图片: " + sbErr.Error())
	}
	defer sbFile.Close()

	sbImg, _, sbErr := image.Decode(sbFile)
	if sbErr != nil {
		panic("无法解码图片: " + sbErr.Error())
	}

	// 应用Smart Blur滤波
	sbFilteredImg := SmartBlur(sbImg)

	// 保存输出图像
	sbOutputFile, sbErr := os.Create(sbOutputPath)
	if sbErr != nil {
		panic("无法创建输出文件: " + sbErr.Error())
	}
	defer sbOutputFile.Close()

	jpeg.Encode(sbOutputFile, sbFilteredImg, &jpeg.Options{Quality: 95})
	println("Smart Blur滤波完成！")
	println("输入图片:", sbInputPath)
	println("输出图片:", sbOutputPath)
	println("参数: 窗口半径=", sbWindowRadius, " 亮度阈值=", sbLuminanceThresh)
}

// Smart Blur参数
const (
	sbuWindowRadius   = 5    // 模糊窗口半径（窗口大小=2*半径+1，推荐3-10）
	sbLuminanceThresh = 10.0 // 亮度差异阈值（控制边缘敏感度，越小越敏感）
)

// SmartBlur 智能模糊核心函数
func SmartBlur(sbInput image.Image) image.Image {
	sbBounds := sbInput.Bounds()
	sbWidth, sbHeight := sbBounds.Max.X, sbBounds.Max.Y
	sbOutput := image.NewRGBA(sbBounds)

	// 预计算所有像素的亮度值（减少重复计算）
	sbLuminance := sbComputeLuminance(sbInput, sbWidth, sbHeight)

	// 遍历每个像素
	for sbY := 0; sbY < sbHeight; sbY++ {
		for sbX := 0; sbX < sbWidth; sbX++ {
			// 保留原始Alpha通道
			_, _, _, sbAlpha := sbInput.At(sbX, sbY).RGBA()
			sbAlpha8 := uint8(sbAlpha >> 8)

			// 中心像素亮度
			sbCenterLum := sbLuminance[sbY][sbX]

			// 累加符合条件的邻域像素（亮度差异小于阈值）
			var sbSumR, sbSumG, sbSumB int
			sbValidCount := 0

			// 遍历窗口内的邻域像素
			for sbKy := -sbuWindowRadius; sbKy <= sbuWindowRadius; sbKy++ {
				for sbKx := -sbuWindowRadius; sbKx <= sbuWindowRadius; sbKx++ {
					// 计算邻域坐标（处理边界）
					sbNeighX := sbX + sbKx
					sbNeighY := sbY + sbKy
					if sbNeighX < 0 || sbNeighX >= sbWidth || sbNeighY < 0 || sbNeighY >= sbHeight {
						continue
					}

					// 检查亮度差异（小于阈值才纳入计算，即非边缘区域）
					sbNeighLum := sbLuminance[sbNeighY][sbNeighX]
					if math.Abs(sbNeighLum-sbCenterLum) <= sbLuminanceThresh {
						// 累加RGB值
						sbR, sbG, sbB, _ := sbInput.At(sbNeighX, sbNeighY).RGBA()
						sbSumR += int(sbR >> 8)
						sbSumG += int(sbG >> 8)
						sbSumB += int(sbB >> 8)
						sbValidCount++
					}
				}
			}

			// 计算平均值（无有效像素时使用原像素值）
			var sbNewR, sbNewG, sbNewB uint8
			if sbValidCount > 0 {
				sbNewR = uint8(sbSumR / sbValidCount)
				sbNewG = uint8(sbSumG / sbValidCount)
				sbNewB = uint8(sbSumB / sbValidCount)
			} else {
				// 边缘区域：直接使用原像素值
				sbR, sbG, sbB, _ := sbInput.At(sbX, sbY).RGBA()
				sbNewR = uint8(sbR >> 8)
				sbNewG = uint8(sbG >> 8)
				sbNewB = uint8(sbB >> 8)
			}

			// 设置输出像素
			sbOutput.SetRGBA(sbX, sbY, color.RGBA{sbNewR, sbNewG, sbNewB, sbAlpha8})
		}
	}

	return sbOutput
}

// 计算图像的亮度矩阵（基于BT.709标准）
func sbComputeLuminance(sbImg image.Image, sbWidth, sbHeight int) [][]float64 {
	sbLum := make([][]float64, sbHeight)
	for y := 0; y < sbHeight; y++ {
		sbLum[y] = make([]float64, sbWidth)
		for x := 0; x < sbWidth; x++ {
			sbR, sbG, sbB, _ := sbImg.At(x, y).RGBA()
			// 转换为0-255范围并计算亮度：Y = 0.2126R + 0.7152G + 0.0722B
			sbR8 := float64(sbR >> 8)
			sbG8 := float64(sbG >> 8)
			sbB8 := float64(sbB >> 8)
			sbLum[y][x] = 0.2126*sbR8 + 0.7152*sbG8 + 0.0722*sbB8
		}
	}
	return sbLum
}

// ========================================================================

// case7 MeanShift滤波算法

// MeanShift（均值漂移）滤波是一种基于密度估计的保边滤波算法，其核心思想是通过迭代寻找像素邻域内的密度峰值（模式） 来更新像素值。
//与其他滤波不同，它同时考虑像素的空间位置和色彩相似性，能在平滑噪声的同时很好地保留边缘和纹理细节，广泛用于图像分割、目标跟踪和纹理平滑等场景。

func case7() {
	// 配置文件路径
	msInputPath := "test4.jpg"         // 输入图像路径
	msOutputPath := "output_case7.jpg" // 输出图像路径

	// 读取输入图像
	msFile, msErr := os.Open(msInputPath)
	if msErr != nil {
		panic("无法打开输入图片: " + msErr.Error())
	}
	defer msFile.Close()

	msImg, _, msErr := image.Decode(msFile)
	if msErr != nil {
		panic("无法解码图片: " + msErr.Error())
	}

	// 应用MeanShift滤波
	msFilteredImg := MeanShiftFilter(msImg)

	// 保存输出图像
	msOutputFile, msErr := os.Create(msOutputPath)
	if msErr != nil {
		panic("无法创建输出文件: " + msErr.Error())
	}
	defer msOutputFile.Close()

	jpeg.Encode(msOutputFile, msFilteredImg, &jpeg.Options{Quality: 95})
	println("MeanShift滤波完成！")
	println("输入图片:", msInputPath)
	println("输出图片:", msOutputPath)
	println("参数: 空间半径=", msSpatialRadius, " 色彩带宽=", msColorBandwidth, " 最大迭代次数=", msMaxIterations)
}

// MeanShift参数（平衡效果与效率）
const (
	msSpatialRadius  = 5    // 空间窗口半径（控制空间范围，推荐3-10）
	msColorBandwidth = 20.0 // 色彩带宽（控制色彩相似度，推荐10-30）
	msMaxIterations  = 5    // 最大迭代次数（避免无限循环）
	msEpsilon        = 1e-3 // 收敛阈值（迭代停止条件）
)

// 像素点结构（包含空间坐标和色彩值）
type msPixel struct {
	x, y    int     // 空间坐标
	r, g, b float64 // 色彩值（0-255）
}

// MeanShiftFilter MeanShift滤波核心函数
func MeanShiftFilter(msInput image.Image) image.Image {
	msBounds := msInput.Bounds()
	msWidth, msHeight := msBounds.Max.X, msBounds.Max.Y
	msOutput := image.NewRGBA(msBounds)

	// 预计算所有像素的RGB值（减少重复访问）
	msPixels := msGetPixelArray(msInput, msWidth, msHeight)

	// 遍历每个像素应用MeanShift
	for y := 0; y < msHeight; y++ {
		for x := 0; x < msWidth; x++ {
			// 对当前像素执行MeanShift迭代
			msResult := msMeanShiftIterate(x, y, msPixels, msWidth, msHeight)

			// 保留原始Alpha通道
			_, _, _, msAlpha := msInput.At(x, y).RGBA()
			msAlpha8 := uint8(msAlpha >> 8)

			// 设置输出像素
			msOutput.SetRGBA(x, y, color.RGBA{
				uint8(msClamp(msResult.r, 0, 255)),
				uint8(msClamp(msResult.g, 0, 255)),
				uint8(msClamp(msResult.b, 0, 255)),
				msAlpha8,
			})
		}
	}

	return msOutput
}

// 对单个像素执行MeanShift迭代
func msMeanShiftIterate(msX, msY int, msPixels [][]msPixel, msWidth, msHeight int) msPixel {
	// 初始化当前点为输入像素
	msCurrent := msPixels[msY][msX]

	// 迭代寻找均值漂移的收敛点
	for iter := 0; iter < msMaxIterations; iter++ {
		var msSumX, msSumY, msSumR, msSumG, msSumB float64
		var msTotalWeight float64

		// 遍历空间窗口内的所有像素
		msMinY := msMax(0, msY-msSpatialRadius)
		msMaxY := msMin(msHeight-1, msY+msSpatialRadius)
		msMinX := msMax(0, msX-msSpatialRadius)
		msMaxX := msMin(msWidth-1, msX+msSpatialRadius)

		for y := msMinY; y <= msMaxY; y++ {
			for x := msMinX; x <= msMaxX; x++ {
				// 计算空间距离（仅考虑窗口内像素）
				msSpatialDist := math.Hypot(float64(x-msX), float64(y-msY))
				if msSpatialDist > float64(msSpatialRadius) {
					continue // 超出空间窗口，跳过
				}

				// 获取邻域像素色彩
				msNeigh := msPixels[y][x]

				// 计算色彩距离（RGB三维空间距离）
				msColorDist := math.Sqrt(
					(msNeigh.r-msCurrent.r)*(msNeigh.r-msCurrent.r) +
						(msNeigh.g-msCurrent.g)*(msNeigh.g-msCurrent.g) +
						(msNeigh.b-msCurrent.b)*(msNeigh.b-msCurrent.b),
				)

				// 计算权重（色彩相似度越高，权重越大）
				if msColorDist > msColorBandwidth {
					continue // 超出色彩带宽，权重为0
				}
				msWeight := math.Exp(-(msColorDist * msColorDist) / (msColorBandwidth * msColorBandwidth))

				// 累加加权和
				msSumX += float64(x) * msWeight
				msSumY += float64(y) * msWeight
				msSumR += msNeigh.r * msWeight
				msSumG += msNeigh.g * msWeight
				msSumB += msNeigh.b * msWeight
				msTotalWeight += msWeight
			}
		}

		// 计算新的均值位置（避免除零）
		if msTotalWeight < 1e-6 {
			break
		}
		msNewX := msSumX / msTotalWeight
		msNewY := msSumY / msTotalWeight
		msNewR := msSumR / msTotalWeight
		msNewG := msSumG / msTotalWeight
		msNewB := msSumB / msTotalWeight

		// 检查收敛（位置和色彩变化均小于阈值）
		msDeltaX := math.Abs(msNewX - float64(msX))
		msDeltaY := math.Abs(msNewY - float64(msY))
		msDeltaColor := math.Sqrt(
			(msNewR-msCurrent.r)*(msNewR-msCurrent.r) +
				(msNewG-msCurrent.g)*(msNewG-msCurrent.g) +
				(msNewB-msCurrent.b)*(msNewB-msCurrent.b),
		)

		if msDeltaX < msEpsilon && msDeltaY < msEpsilon && msDeltaColor < msEpsilon {
			break // 已收敛，停止迭代
		}

		// 更新当前点和坐标
		msX = int(msNewX + 0.5) // 四舍五入到整数坐标
		msY = int(msNewY + 0.5)
		msCurrent.r = msNewR
		msCurrent.g = msNewG
		msCurrent.b = msNewB
	}

	return msCurrent
}

// 辅助函数

// 将图像转换为像素数组（包含空间坐标和色彩值）
func msGetPixelArray(msImg image.Image, msWidth, msHeight int) [][]msPixel {
	msPixels := make([][]msPixel, msHeight)
	for y := 0; y < msHeight; y++ {
		msPixels[y] = make([]msPixel, msWidth)
		for x := 0; x < msWidth; x++ {
			msR, msG, msB, _ := msImg.At(x, y).RGBA()
			msPixels[y][x] = msPixel{
				x: x,
				y: y,
				r: float64(msR >> 8),
				g: float64(msG >> 8),
				b: float64(msB >> 8),
			}
		}
	}
	return msPixels
}

// 限制值在[min, max]范围内
func msClamp(msVal, msMin, msMax float64) float64 {
	if msVal < msMin {
		return msMin
	}
	if msVal > msMax {
		return msMax
	}
	return msVal
}

// 取最大值
func msMax(msA, msB int) int {
	if msA > msB {
		return msA
	}
	return msB
}

// 取最小值
func msMin(msA, msB int) int {
	if msA < msB {
		return msA
	}
	return msB
}

// ========================================================================

// case8 BEEPS滤波算法

// BEEPS（Bilateral Edge-Enhancing Preserving Smoothing）滤波算法是一种基于双边滤波思想的改进型保边滤波算法，
//其核心特点是在平滑噪声的同时增强边缘对比度，通过动态调整滤波权重实现 “平滑平坦区域、增强边缘细节” 的双重效果。
//与传统双边滤波相比，BEEPS 对边缘的响应更敏感，能在去噪的同时避免边缘模糊，甚至轻微增强边缘清晰度，适合用于图像预处理、细节增强等场景。

func case8() {
	// 配置文件路径
	beepsInputPath := "test4.jpg"         // 输入图像路径
	beepsOutputPath := "output_case8.jpg" // 输出图像路径

	// 读取输入图像
	beepsFile, beepsErr := os.Open(beepsInputPath)
	if beepsErr != nil {
		panic("无法打开输入图片: " + beepsErr.Error())
	}
	defer beepsFile.Close()

	beepsImg, _, beepsErr := image.Decode(beepsFile)
	if beepsErr != nil {
		panic("无法解码图片: " + beepsErr.Error())
	}

	// 应用BEEPS滤波
	beepsFilteredImg := BEEPSFilter(beepsImg)

	// 保存输出图像
	beepsOutputFile, beepsErr := os.Create(beepsOutputPath)
	if beepsErr != nil {
		panic("无法创建输出文件: " + beepsErr.Error())
	}
	defer beepsOutputFile.Close()

	jpeg.Encode(beepsOutputFile, beepsFilteredImg, &jpeg.Options{Quality: 95})
	println("BEEPS滤波完成！")
	println("输入图片:", beepsInputPath)
	println("输出图片:", beepsOutputPath)
	println("参数: 窗口半径=", beepsWindowRadius, " 边缘增强系数=", beepsEdgeBoost)
}

// BEEPS滤波参数
const (
	beepsWindowRadius = 3    // 窗口半径（推荐3-5，控制平滑范围）
	beepsSpatialSigma = 2.0  // 空间高斯标准差（控制空间权重衰减速度）
	beepsEdgeSigma    = 10.0 // 边缘高斯标准差（控制边缘敏感度，值越小越敏感）
	beepsEdgeBoost    = 1.2  // 边缘增强系数（>1增强边缘，<1抑制边缘）
)

// BEEPSFilter 核心滤波函数
func BEEPSFilter(beepsInput image.Image) image.Image {
	beepsBounds := beepsInput.Bounds()
	beepsWidth, beepsHeight := beepsBounds.Max.X, beepsBounds.Max.Y
	beepsOutput := image.NewRGBA(beepsBounds)

	// 预计算梯度（边缘强度）矩阵
	beepsGrad := beepsComputeGradient(beepsInput, beepsWidth, beepsHeight)

	// 遍历每个像素
	for y := 0; y < beepsHeight; y++ {
		for x := 0; x < beepsWidth; x++ {
			// 保留原始Alpha通道
			_, _, _, beepsAlpha := beepsInput.At(x, y).RGBA()
			beepsAlpha8 := uint8(beepsAlpha >> 8)

			// 中心像素的RGB值
			beepsCR, beepsCG, beepsCB, _ := beepsInput.At(x, y).RGBA()
			beepsCtrR := float64(beepsCR >> 8)
			beepsCtrG := float64(beepsCG >> 8)
			beepsCtrB := float64(beepsCB >> 8)

			var beepsSumR, beepsSumG, beepsSumB, beepsTotalWeight float64

			// 遍历邻域窗口
			beepsMinY := beepsMax(0, y-beepsWindowRadius)
			beepsMaxY := beepsMin(beepsHeight-1, y+beepsWindowRadius)
			beepsMinX := beepsMax(0, x-beepsWindowRadius)
			beepsMaxX := beepsMin(beepsWidth-1, x+beepsWindowRadius)

			for ny := beepsMinY; ny <= beepsMaxY; ny++ {
				for nx := beepsMinX; nx <= beepsMaxX; nx++ {
					// 计算空间距离权重（高斯衰减）
					beepsDx := float64(nx - x)
					beepsDy := float64(ny - y)
					beepsSpatialDist := beepsDx*beepsDx + beepsDy*beepsDy
					beepsSpatialW := math.Exp(-beepsSpatialDist / (2 * beepsSpatialSigma * beepsSpatialSigma))

					// 获取邻域像素RGB值
					beepsNR, beepsNG, beepsNB, _ := beepsInput.At(nx, ny).RGBA()
					beepsNeighR := float64(beepsNR >> 8)
					beepsNeighG := float64(beepsNG >> 8)
					beepsNeighB := float64(beepsNB >> 8)

					// 计算色彩差异（用于边缘判断）
					beepsColorDiff := math.Sqrt(
						(beepsNeighR-beepsCtrR)*(beepsNeighR-beepsCtrR) +
							(beepsNeighG-beepsCtrG)*(beepsNeighG-beepsCtrG) +
							(beepsNeighB-beepsCtrB)*(beepsNeighB-beepsCtrB),
					)

					// 结合梯度计算边缘权重（低梯度区域权重高）
					beepsEdgeStr := beepsGrad[ny][nx] // 边缘强度（梯度值）
					beepsEdgeW := math.Exp(-(beepsColorDiff*beepsColorDiff + beepsEdgeStr*beepsEdgeStr) / (2 * beepsEdgeSigma * beepsEdgeSigma))

					// 边缘增强调整：对边缘两侧像素施加差异化权重
					if beepsEdgeStr > 5.0 { // 阈值判断是否为显著边缘
						// 边缘方向的简单判断（基于色彩差异符号）
						if (beepsNeighR > beepsCtrR && beepsColorDiff > 3.0) ||
							(beepsNeighG > beepsCtrG && beepsColorDiff > 3.0) ||
							(beepsNeighB > beepsCtrB && beepsColorDiff > 3.0) {
							beepsEdgeW *= beepsEdgeBoost // 边缘亮侧增强
						} else {
							beepsEdgeW /= beepsEdgeBoost // 边缘暗侧减弱（增强对比度）
						}
					}

					// 总权重 = 空间权重 × 边缘权重
					beepsWeight := beepsSpatialW * beepsEdgeW

					// 累加加权像素值
					beepsSumR += beepsNeighR * beepsWeight
					beepsSumG += beepsNeighG * beepsWeight
					beepsSumB += beepsNeighB * beepsWeight
					beepsTotalWeight += beepsWeight
				}
			}

			// 计算输出像素值（避免除零）
			var beepsOutR, beepsOutG, beepsOutB float64
			if beepsTotalWeight > 1e-6 {
				beepsOutR = beepsSumR / beepsTotalWeight
				beepsOutG = beepsSumG / beepsTotalWeight
				beepsOutB = beepsSumB / beepsTotalWeight
			} else {
				beepsOutR = beepsCtrR
				beepsOutG = beepsCtrG
				beepsOutB = beepsCtrB
			}

			// 钳位并设置输出像素
			beepsOutput.SetRGBA(x, y, color.RGBA{
				uint8(beepsClamp(beepsOutR, 0, 255)),
				uint8(beepsClamp(beepsOutG, 0, 255)),
				uint8(beepsClamp(beepsOutB, 0, 255)),
				beepsAlpha8,
			})
		}
	}

	return beepsOutput
}

// 计算图像梯度（边缘强度）矩阵
func beepsComputeGradient(beepsImg image.Image, beepsWidth, beepsHeight int) [][]float64 {
	beepsGrad := make([][]float64, beepsHeight)
	for y := 0; y < beepsHeight; y++ {
		beepsGrad[y] = make([]float64, beepsWidth)
		for x := 0; x < beepsWidth; x++ {
			// 取邻域像素计算梯度（简单 Sobel 近似）
			beepsR, beepsG, beepsB, _ := beepsImg.At(x, y).RGBA()
			beepsCtr := float64(beepsR>>8+beepsG>>8+beepsB>>8) / 3.0 // 灰度值

			// 右邻域
			beepsRight := beepsCtr
			if x+1 < beepsWidth {
				beepsR, beepsG, beepsB, _ := beepsImg.At(x+1, y).RGBA()
				beepsRight = float64(beepsR>>8+beepsG>>8+beepsB>>8) / 3.0
			}

			// 下邻域
			beepsDown := beepsCtr
			if y+1 < beepsHeight {
				beepsR, beepsG, beepsB, _ := beepsImg.At(x, y+1).RGBA()
				beepsDown = float64(beepsR>>8+beepsG>>8+beepsB>>8) / 3.0
			}

			// 梯度强度（x和y方向差异的平方和开方）
			beepsGrad[y][x] = math.Sqrt((beepsRight-beepsCtr)*(beepsRight-beepsCtr) + (beepsDown-beepsCtr)*(beepsDown-beepsCtr))
		}
	}
	return beepsGrad
}

// 辅助函数：值钳位
func beepsClamp(beepsVal, beepsMin, beepsMax float64) float64 {
	if beepsVal < beepsMin {
		return beepsMin
	}
	if beepsVal > beepsMax {
		return beepsMax
	}
	return beepsVal
}

// 辅助函数：取最大值
func beepsMax(beepsA, beepsB int) int {
	if beepsA > beepsB {
		return beepsA
	}
	return beepsB
}

// 辅助函数：取最小值
func beepsMin(beepsA, beepsB int) int {
	if beepsA < beepsB {
		return beepsA
	}
	return beepsB
}

// ========================================================================

// case9 DCT降噪滤波算法

// DCT（Discrete Cosine Transform，离散余弦变换）降噪滤波算法的核心思想是利用图像信号与噪声在频率域的分布差异实现降噪：
//图像的主要信息（如轮廓、平滑区域）集中在低频分量，而噪声主要分布在高频分量。通过对图像分块进行 DCT 变换，对高频系数施加阈值处理（抑制噪声），
//再通过逆 DCT 变换恢复图像，从而实现降噪效果。

func case9() {
	// 配置文件路径
	dctInputPath := "test4.jpg"         // 输入图像路径
	dctOutputPath := "output_case9.jpg" // 输出图像路径

	// 读取输入图像
	dctFile, dctErr := os.Open(dctInputPath)
	if dctErr != nil {
		panic("无法打开输入图片: " + dctErr.Error())
	}
	defer dctFile.Close()

	dctImg, _, dctErr := image.Decode(dctFile)
	if dctErr != nil {
		panic("无法解码图片: " + dctErr.Error())
	}

	// 应用DCT降噪
	dctFilteredImg := DCTDenoise(dctImg)

	// 保存输出图像
	dctOutputFile, dctErr := os.Create(dctOutputPath)
	if dctErr != nil {
		panic("无法创建输出文件: " + dctErr.Error())
	}
	defer dctOutputFile.Close()

	jpeg.Encode(dctOutputFile, dctFilteredImg, &jpeg.Options{Quality: 95})
	println("DCT降噪完成！")
	println("输入图片:", dctInputPath)
	println("输出图片:", dctOutputPath)
	println("参数: 块大小=", dctBlockSize, "x", dctBlockSize, " 阈值=", dctThreshold, " 软阈值=", dctIsSoft)
}

// DCT降噪参数
const (
	dctBlockSize = 8    // 分块大小（8x8，经典选择）
	dctThreshold = 20.0 // 高频系数阈值（控制降噪强度，推荐10-50）
	dctIsSoft    = true // 阈值类型：true=软阈值（更平滑），false=硬阈值（保留更多细节）
)

// DCTDenoise 核心DCT降噪函数
func DCTDenoise(dctInput image.Image) image.Image {
	dctBounds := dctInput.Bounds()
	dctWidth, dctHeight := dctBounds.Max.X, dctBounds.Max.Y
	dctOutput := image.NewRGBA(dctBounds)

	// 分通道处理（RGB分别降噪）
	dctR, dctG, dctB := dctImgToFloatChannels(dctInput)

	// 对每个通道应用DCT降噪
	dctDenoisedR := dctProcessChannel(dctR, dctWidth, dctHeight)
	dctDenoisedG := dctProcessChannel(dctG, dctWidth, dctHeight)
	dctDenoisedB := dctProcessChannel(dctB, dctWidth, dctHeight)

	// 合并通道并保留Alpha值
	for y := 0; y < dctHeight; y++ {
		for x := 0; x < dctWidth; x++ {
			_, _, _, dctA := dctInput.At(x, y).RGBA()
			dctA8 := uint8(dctA >> 8)

			dctR := dctClamp(dctDenoisedR[y][x], 0, 255)
			dctG := dctClamp(dctDenoisedG[y][x], 0, 255)
			dctB := dctClamp(dctDenoisedB[y][x], 0, 255)

			dctOutput.SetRGBA(x, y, color.RGBA{uint8(dctR), uint8(dctG), uint8(dctB), dctA8})
		}
	}

	return dctOutput
}

// 单通道DCT降噪处理
func dctProcessChannel(dctChan [][]float64, dctWidth, dctHeight int) [][]float64 {
	// 初始化输出通道
	dctOutput := make([][]float64, dctHeight)
	for y := 0; y < dctHeight; y++ {
		dctOutput[y] = make([]float64, dctWidth)
	}

	// 预计算8x8 DCT变换矩阵（提高效率）
	dctMat := dctComputeMatrix()

	// 遍历所有块（按块大小步进）
	for y := 0; y < dctHeight; y += dctBlockSize {
		for x := 0; x < dctWidth; x += dctBlockSize {
			// 提取当前块（处理边缘块可能不足8x8的情况）
			dctBlock := dctExtractBlock(dctChan, x, y, dctWidth, dctHeight)

			// 1. 对块进行DCT变换
			dctTransformed := dctApply(dctBlock, dctMat)

			// 2. 对高频系数应用阈值滤波（抑制噪声）
			dctFiltered := dctApplyThreshold(dctTransformed)

			// 3. 逆DCT变换还原块
			dctInverse := dctInverseApply(dctFiltered, dctMat)

			// 4. 将处理后的块写入输出
			dctWriteBlock(dctOutput, dctInverse, x, y, dctWidth, dctHeight)
		}
	}

	return dctOutput
}

// 提取8x8块（边缘块用0填充）
func dctExtractBlock(dctChan [][]float64, dctX, dctY, dctWidth, dctHeight int) [][]float64 {
	dctBlock := make([][]float64, dctBlockSize)
	for i := 0; i < dctBlockSize; i++ {
		dctBlock[i] = make([]float64, dctBlockSize)
		for j := 0; j < dctBlockSize; j++ {
			// 计算原始图像坐标
			dctImgX := dctX + j
			dctImgY := dctY + i
			// 边缘检查：超出图像范围则用0填充
			if dctImgX < dctWidth && dctImgY < dctHeight {
				dctBlock[i][j] = dctChan[dctImgY][dctImgX]
			} else {
				dctBlock[i][j] = 0
			}
		}
	}
	return dctBlock
}

// 将处理后的块写入输出图像
func dctWriteBlock(dctOutput [][]float64, dctBlock [][]float64, dctX, dctY, dctWidth, dctHeight int) {
	for i := 0; i < dctBlockSize; i++ {
		for j := 0; j < dctBlockSize; j++ {
			dctImgX := dctX + j
			dctImgY := dctY + i
			if dctImgX < dctWidth && dctImgY < dctHeight {
				dctOutput[dctImgY][dctImgX] = dctBlock[i][j]
			}
		}
	}
}

// 计算8x8 DCT变换矩阵
func dctComputeMatrix() [][]float64 {
	dctMat := make([][]float64, dctBlockSize)
	for i := 0; i < dctBlockSize; i++ {
		dctMat[i] = make([]float64, dctBlockSize)
		for j := 0; j < dctBlockSize; j++ {
			var dctAlpha float64
			if i == 0 {
				dctAlpha = 1.0 / math.Sqrt(float64(dctBlockSize)) // 第一行系数特殊处理
			} else {
				dctAlpha = math.Sqrt(2.0 / float64(dctBlockSize))
			}
			// DCT变换公式
			dctMat[i][j] = dctAlpha * math.Cos(
				math.Pi*float64(2*j+1)*float64(i)/(2*float64(dctBlockSize)),
			)
		}
	}
	return dctMat
}

// 应用DCT变换
func dctApply(dctBlock, dctMat [][]float64) [][]float64 {
	// 结果矩阵初始化
	dctResult := make([][]float64, dctBlockSize)
	for i := 0; i < dctBlockSize; i++ {
		dctResult[i] = make([]float64, dctBlockSize)
	}

	// DCT变换：result = mat * block * mat^T（矩阵乘法）
	for i := 0; i < dctBlockSize; i++ {
		for j := 0; j < dctBlockSize; j++ {
			var dctSum float64
			for k := 0; k < dctBlockSize; k++ {
				for l := 0; l < dctBlockSize; l++ {
					dctSum += dctMat[i][k] * dctBlock[k][l] * dctMat[j][l]
				}
			}
			dctResult[i][j] = dctSum
		}
	}
	return dctResult
}

// 应用逆DCT变换
func dctInverseApply(dctBlock, dctMat [][]float64) [][]float64 {
	// 结果矩阵初始化
	dctResult := make([][]float64, dctBlockSize)
	for i := 0; i < dctBlockSize; i++ {
		dctResult[i] = make([]float64, dctBlockSize)
	}

	// 逆DCT变换：result = mat^T * block * mat（矩阵乘法）
	for i := 0; i < dctBlockSize; i++ {
		for j := 0; j < dctBlockSize; j++ {
			var dctSum float64
			for k := 0; k < dctBlockSize; k++ {
				for l := 0; l < dctBlockSize; l++ {
					dctSum += dctMat[k][i] * dctBlock[k][l] * dctMat[l][j]
				}
			}
			dctResult[i][j] = dctSum
		}
	}
	return dctResult
}

// 对DCT系数应用阈值滤波
func dctApplyThreshold(dctBlock [][]float64) [][]float64 {
	dctResult := make([][]float64, dctBlockSize)
	for i := 0; i < dctBlockSize; i++ {
		dctResult[i] = make([]float64, dctBlockSize)
		for j := 0; j < dctBlockSize; j++ {
			dctVal := dctBlock[i][j]
			// 低频分量（左上角）通常保留，高频分量（右下角）施加阈值
			if math.Abs(dctVal) < dctThreshold {
				if dctIsSoft {
					// 软阈值：保留符号，减去阈值
					if dctVal > 0 {
						dctResult[i][j] = dctVal - dctThreshold
					} else {
						dctResult[i][j] = dctVal + dctThreshold
					}
				} else {
					// 硬阈值：直接置0
					dctResult[i][j] = 0
				}
			} else {
				// 大于阈值的系数保留
				dctResult[i][j] = dctVal
			}
		}
	}
	return dctResult
}

// 将图像转换为浮点型通道（0-255）
func dctImgToFloatChannels(dctImg image.Image) (dctR, dctG, dctB [][]float64) {
	dctBounds := dctImg.Bounds()
	dctWidth, dctHeight := dctBounds.Max.X, dctBounds.Max.Y

	dctR = make([][]float64, dctHeight)
	dctG = make([][]float64, dctHeight)
	dctB = make([][]float64, dctHeight)
	for y := 0; y < dctHeight; y++ {
		dctR[y] = make([]float64, dctWidth)
		dctG[y] = make([]float64, dctWidth)
		dctB[y] = make([]float64, dctWidth)
		for x := 0; x < dctWidth; x++ {
			dctRVal, dctGVal, dctBVal, _ := dctImg.At(x, y).RGBA()
			dctR[y][x] = float64(dctRVal >> 8)
			dctG[y][x] = float64(dctGVal >> 8)
			dctB[y][x] = float64(dctBVal >> 8)
		}
	}
	return dctR, dctG, dctB
}

// 钳位值到[min, max]
func dctClamp(dctVal, dctMin, dctMax float64) float64 {
	if dctVal < dctMin {
		return dctMin
	}
	if dctVal > dctMax {
		return dctMax
	}
	return dctVal
}

// ========================================================================

// case10 非局部均值滤波算法

// 非局部均值滤波（Non-Local Means, NLM）是一种基于全局相似性的降噪算法，其核心思想是：图像中存在大量重复或相似的局部结构（如纹理、图案），
//噪声会破坏这种相似性，通过寻找全局范围内的相似块并加权平均，可以有效还原真实信号。与局部均值滤波仅利用局部邻域不同，NLM 能利用图像的全局冗余信息，
//在强降噪的同时更好地保留细节，尤其适合处理高斯噪声。

func case10() {
	inputPath := "test4.jpg"
	outputPath := "output_case10.jpg"

	// 读取图像
	file, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图片: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("无法解码图片: " + err.Error())
	}

	// 应用非局部均值滤波
	filteredImg := NonLocalMeans(img)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, filteredImg, &jpeg.Options{Quality: 95})
	println("非局部均值滤波完成（无init函数）！")
	println("输入图片:", inputPath)
	println("输出图片:", outputPath)
	println("参数: 块大小=", nlmPatchSize, "x", nlmPatchSize,
		" 搜索窗口=", nlmSearchSize, "x", nlmSearchSize,
		" 衰减参数=", nlmH)
}

// 非局部均值滤波参数
const (
	nlmPatchSize  = 7    // 块大小（7x7，推荐5-9，需为奇数）
	nlmSearchSize = 21   // 搜索窗口大小（21x21，推荐15-25，需为奇数）
	nlmH          = 10.0 // 相似度衰减参数（控制权重，推荐5-20）
	nlmSigma      = 1.0  // 块内高斯加权标准差
)

// NonLocalMeans 非局部均值滤波核心函数（无init版本）
func NonLocalMeans(nlmInput image.Image) image.Image {
	nlmBounds := nlmInput.Bounds()
	nlmWidth, nlmHeight := nlmBounds.Max.X, nlmBounds.Max.Y
	nlmOutput := image.NewRGBA(nlmBounds)

	// 1. 预计算块内高斯权重（替代init函数的逻辑）
	gaussianWeights := computeGaussianWeights()

	// 2. 转换为浮点型RGB通道（0-255）
	nlmR, nlmG, nlmB := nlmImgToFloat(nlmInput)

	// 3. 分通道处理（传递高斯权重）
	nlmDenoisedR := nlmProcessChannel(nlmR, nlmWidth, nlmHeight, gaussianWeights)
	nlmDenoisedG := nlmProcessChannel(nlmG, nlmWidth, nlmHeight, gaussianWeights)
	nlmDenoisedB := nlmProcessChannel(nlmB, nlmWidth, nlmHeight, gaussianWeights)

	// 4. 合并通道并保留Alpha
	for y := 0; y < nlmHeight; y++ {
		for x := 0; x < nlmWidth; x++ {
			_, _, _, nlmA := nlmInput.At(x, y).RGBA()
			nlmA8 := uint8(nlmA >> 8)

			r := nlmClamp(nlmDenoisedR[y][x], 0, 255)
			g := nlmClamp(nlmDenoisedG[y][x], 0, 255)
			b := nlmClamp(nlmDenoisedB[y][x], 0, 255)

			nlmOutput.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), nlmA8})
		}
	}

	return nlmOutput
}

// 计算块内高斯权重（替代init函数）
func computeGaussianWeights() [][]float64 {
	half := nlmPatchSize / 2
	weights := make([][]float64, nlmPatchSize)
	for i := 0; i < nlmPatchSize; i++ {
		weights[i] = make([]float64, nlmPatchSize)
		for j := 0; j < nlmPatchSize; j++ {
			dx := float64(j - half)
			dy := float64(i - half)
			weights[i][j] = math.Exp(-(dx*dx + dy*dy) / (2 * nlmSigma * nlmSigma))
		}
	}
	return weights
}

// 单通道非局部均值处理（接收高斯权重参数）
func nlmProcessChannel(nlmChan [][]float64, nlmWidth, nlmHeight int, gaussianWeights [][]float64) [][]float64 {
	nlmOutput := make([][]float64, nlmHeight)
	for y := 0; y < nlmHeight; y++ {
		nlmOutput[y] = make([]float64, nlmWidth)
	}

	halfPatch := nlmPatchSize / 2
	halfSearch := nlmSearchSize / 2

	// 遍历每个像素
	for y := 0; y < nlmHeight; y++ {
		for x := 0; x < nlmWidth; x++ {
			// 确定搜索窗口范围（避免越界）
			startY := nlmMax(0, y-halfSearch)
			endY := nlmMin(nlmHeight-1, y+halfSearch)
			startX := nlmMax(0, x-halfSearch)
			endX := nlmMin(nlmWidth-1, x+halfSearch)

			var totalWeight float64
			var weightedSum float64

			// 在搜索窗口内寻找相似块
			for ny := startY; ny <= endY; ny++ {
				for nx := startX; nx <= endX; nx++ {
					// 计算当前块与中心块的相似度（传递高斯权重）
					distance := nlmComputeDistance(nlmChan, x, y, nx, ny, nlmWidth, nlmHeight, halfPatch, gaussianWeights)

					// 距离越小，权重越大（指数衰减）
					weight := math.Exp(-distance / (nlmH * nlmH))
					totalWeight += weight

					// 累加相似块中心像素的加权值
					weightedSum += weight * nlmChan[ny][nx]
				}
			}

			// 归一化权重并计算输出值
			if totalWeight > 1e-6 {
				nlmOutput[y][x] = weightedSum / totalWeight
			} else {
				nlmOutput[y][x] = nlmChan[y][x] // 无相似块时保留原值
			}
		}
	}

	return nlmOutput
}

// 计算两个块之间的加权欧氏距离（接收高斯权重参数）
func nlmComputeDistance(nlmChan [][]float64, x1, y1, x2, y2, width, height, halfPatch int, gaussianWeights [][]float64) float64 {
	var distance float64

	// 遍历块内所有像素
	for i := -halfPatch; i <= halfPatch; i++ {
		for j := -halfPatch; j <= halfPatch; j++ {
			// 计算块内像素坐标（处理边界）
			p1x := x1 + j
			p1y := y1 + i
			p2x := x2 + j
			p2y := y2 + i

			// 边界检查（超出范围则用边缘像素）
			p1x = nlmClampInt(p1x, 0, width-1)
			p1y = nlmClampInt(p1y, 0, height-1)
			p2x = nlmClampInt(p2x, 0, width-1)
			p2y = nlmClampInt(p2y, 0, height-1)

			// 像素值差异（应用高斯权重）
			diff := nlmChan[p1y][p1x] - nlmChan[p2y][p2x]
			distance += gaussianWeights[i+halfPatch][j+halfPatch] * diff * diff
		}
	}

	return distance
}

// 图像转浮点通道
func nlmImgToFloat(nlmImg image.Image) (r, g, b [][]float64) {
	bounds := nlmImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	r = make([][]float64, height)
	g = make([][]float64, height)
	b = make([][]float64, height)
	for y := 0; y < height; y++ {
		r[y] = make([]float64, width)
		g[y] = make([]float64, width)
		b[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			rVal, gVal, bVal, _ := nlmImg.At(x, y).RGBA()
			r[y][x] = float64(rVal >> 8)
			g[y][x] = float64(gVal >> 8)
			b[y][x] = float64(bVal >> 8)
		}
	}
	return r, g, b
}

// 浮点值钳位
func nlmClamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// 整数钳位
func nlmClampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// 取最大值
func nlmMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 取最小值
func nlmMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ========================================================================

// case11 加权中值滤波算法

// 加权中值滤波（Weighted Median Filter）是中值滤波的改进版本，通过为窗口内不同位置的像素分配差异化权重，
//实现 “重要像素（如中心附近或边缘区域）影响更大” 的选择性滤波。相比普通中值滤波（所有像素权重均等），
//它能在去除脉冲噪声（如椒盐噪声）的同时更好地保留边缘和细节，平衡降噪效果与细节损失。

func case11() {
	inputPath := "test4.jpg"
	outputPath := "output_case11.jpg"

	// 读取图像
	file, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图片: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("无法解码图片: " + err.Error())
	}

	// 应用加权中值滤波
	filteredImg := WeightedMedianFilter(img)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, filteredImg, &jpeg.Options{Quality: 95})
	println("加权中值滤波完成！")
	println("输入图片:", inputPath)
	println("输出图片:", outputPath)
	println("参数: 窗口大小=", wmWindowSize, "x", wmWindowSize)
}

// 加权中值滤波参数
const (
	wmWindowSize = 3 // 窗口大小（3x3，推荐3-7，需为奇数）
)

// 生成权重矩阵（中心权重最高，向外递减，示例为3x3）
// 可根据需求调整权重分布（需保证权重为非负整数）
func generateWeightMatrix() [][]int {
	switch wmWindowSize {
	case 3:
		return [][]int{
			{1, 2, 1},
			{2, 4, 2},
			{1, 2, 1},
		}
	case 5:
		return [][]int{
			{1, 1, 2, 1, 1},
			{1, 2, 3, 2, 1},
			{2, 3, 5, 3, 2},
			{1, 2, 3, 2, 1},
			{1, 1, 2, 1, 1},
		}
	default:
		panic("不支持的窗口大小（仅支持3或5）")
	}
}

// WeightedMedianFilter 加权中值滤波核心函数
func WeightedMedianFilter(wmInput image.Image) image.Image {
	wmBounds := wmInput.Bounds()
	wmWidth, wmHeight := wmBounds.Max.X, wmBounds.Max.Y
	wmOutput := image.NewRGBA(wmBounds)

	// 生成权重矩阵
	weightMatrix := generateWeightMatrix()
	halfWin := wmWindowSize / 2 // 窗口半宽（用于计算边界）

	// 分通道处理（RGB）
	wmR, wmG, wmB := wmImgToChannels(wmInput)
	wmDenoisedR := wmProcessChannel(wmR, wmWidth, wmHeight, weightMatrix, halfWin)
	wmDenoisedG := wmProcessChannel(wmG, wmWidth, wmHeight, weightMatrix, halfWin)
	wmDenoisedB := wmProcessChannel(wmB, wmWidth, wmHeight, weightMatrix, halfWin)

	// 合并通道并保留Alpha
	for y := 0; y < wmHeight; y++ {
		for x := 0; x < wmWidth; x++ {
			_, _, _, wmA := wmInput.At(x, y).RGBA()
			wmA8 := uint8(wmA >> 8)

			r := wmClamp(wmDenoisedR[y][x], 0, 255)
			g := wmClamp(wmDenoisedG[y][x], 0, 255)
			b := wmClamp(wmDenoisedB[y][x], 0, 255)

			wmOutput.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), wmA8})
		}
	}

	return wmOutput
}

// 单通道加权中值处理
func wmProcessChannel(wmChan [][]uint8, wmWidth, wmHeight int, weights [][]int, halfWin int) [][]uint8 {
	output := make([][]uint8, wmHeight)
	for y := 0; y < wmHeight; y++ {
		output[y] = make([]uint8, wmWidth)
		for x := 0; x < wmWidth; x++ {
			// 收集窗口内像素值，并按权重扩展
			weightedValues := wmCollectWeightedValues(wmChan, x, y, wmWidth, wmHeight, weights, halfWin)
			// 计算中值
			median := wmComputeMedian(weightedValues)
			output[y][x] = median
		}
	}
	return output
}

// 收集窗口内像素值并按权重扩展（权重w对应的值会被添加w次）
func wmCollectWeightedValues(wmChan [][]uint8, x, y, width, height int, weights [][]int, halfWin int) []uint8 {
	var values []uint8

	// 遍历窗口内所有位置
	for ky := 0; ky < wmWindowSize; ky++ {
		for kx := 0; kx < wmWindowSize; kx++ {
			// 计算原始图像坐标（相对于中心像素的偏移）
			imgX := x + (kx - halfWin)
			imgY := y + (ky - halfWin)

			// 边界处理：超出范围则使用最近的边缘像素
			imgX = wmClampInt(imgX, 0, width-1)
			imgY = wmClampInt(imgY, 0, height-1)

			// 获取当前位置的权重
			w := weights[ky][kx]
			if w <= 0 {
				continue // 权重为0则跳过
			}

			// 按权重复制像素值（添加w次）
			val := wmChan[imgY][imgX]
			for i := 0; i < w; i++ {
				values = append(values, val)
			}
		}
	}

	return values
}

// 计算序列的中值
func wmComputeMedian(values []uint8) uint8 {
	if len(values) == 0 {
		return 0 // 理论上不会发生（权重总和至少为1）
	}

	// 排序序列
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	// 取中值（奇数长度取中间，偶数长度取中间左值）
	mid := len(values) / 2
	return values[mid]
}

// 图像转换为RGB通道（uint8类型，0-255）
func wmImgToChannels(wmImg image.Image) (r, g, b [][]uint8) {
	bounds := wmImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	r = make([][]uint8, height)
	g = make([][]uint8, height)
	b = make([][]uint8, height)
	for y := 0; y < height; y++ {
		r[y] = make([]uint8, width)
		g[y] = make([]uint8, width)
		b[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			rVal, gVal, bVal, _ := wmImg.At(x, y).RGBA()
			r[y][x] = uint8(rVal >> 8)
			g[y][x] = uint8(gVal >> 8)
			b[y][x] = uint8(bVal >> 8)
		}
	}
	return r, g, b
}

// 整数钳位
func wmClampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// 字节钳位（0-255）
func wmClamp(val uint8, min, max uint8) uint8 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// ========================================================================

// case12 基于颜色空间的皮肤检测算法

// 基于颜色空间的皮肤检测算法利用了皮肤颜色在特定颜色空间中具有稳定聚类特性的特点，通过定义肤色在该空间中的阈值范围，
//实现对皮肤区域的快速识别。常用的颜色空间包括 YCrCb（亮度 - 色度分离，抗光照变化能力强）、HSV（色调 - 饱和度 - 明度，直观反映颜色特征）等，
//其中 YCrCb 因能有效分离亮度和色度分量，在皮肤检测中应用最广泛。

// 算法原理
//以 YCrCb 颜色空间为例，核心步骤如下：
//颜色空间转换：将输入图像从 RGB 转换为 YCrCb（Y 为亮度，Cr 和 Cb 为色度分量）。
//肤色阈值定义：皮肤在 Cr（红色分量）和 Cb（蓝色分量）通道中分布在特定范围（例如：Cr∈[133, 173]，Cb∈[77, 127]，不同人种略有差异）。
//像素分类：遍历图像像素，判断其 Cr 和 Cb 值是否在肤色阈值范围内，是则标记为皮肤，否则为非皮肤。
//后处理（可选）：通过形态学操作（如腐蚀、膨胀）去除噪声，使皮肤区域更连贯。

func case12() {
	inputPath := "test4.jpg"                 // 输入图像路径
	resultPath := "output_case12_result.jpg" // 皮肤高亮结果
	maskPath := "output_case12_mask.jpg"     // 皮肤掩码结果

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

	// 执行皮肤检测
	resultImg, maskImg := SkinDetection(img)

	// 保存皮肤高亮结果
	outputFile, err := os.Create(resultPath)
	if err != nil {
		panic("无法创建结果文件: " + err.Error())
	}
	defer outputFile.Close()
	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})

	// 保存皮肤掩码结果
	maskFile, err := os.Create(maskPath)
	if err != nil {
		panic("无法创建掩码文件: " + err.Error())
	}
	defer maskFile.Close()
	jpeg.Encode(maskFile, maskImg, &jpeg.Options{Quality: 95})

	println("皮肤检测完成！")
	println("结果图:", resultPath)
	println("掩码图:", maskPath)
	println("YCrCb阈值: Cr=[", crMin, ",", crMax, "], Cb=[", cbMin, ",", cbMax, "]")
}

// YCrCb颜色空间中的肤色阈值（适用于多数人种，可根据需求调整）
const (
	crMin = 133 // 红色色度最小值
	crMax = 173 // 红色色度最大值
	cbMin = 77  // 蓝色色度最小值
	cbMax = 127 // 蓝色色度最大值
)

// SkinDetection 基于YCrCb的皮肤检测核心函数
func SkinDetection(input image.Image) (image.Image, image.Image) {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 生成皮肤掩码（二值图：皮肤=白色，非皮肤=黑色）
	mask := image.NewRGBA(bounds)
	// 生成带皮肤高亮的结果图（皮肤保留原色，非皮肤灰色）
	result := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取RGB值（0-255）
			r, g, b, a := input.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 1. RGB转换为YCrCb
			yVal, cr, cb := rgbToYCrCb(r8, g8, b8)

			// 2. 判断是否为皮肤（Cr和Cb在阈值范围内）
			isSkin := cr >= crMin && cr <= crMax && cb >= cbMin && cb <= cbMax

			// 3. 更新掩码和结果图
			if isSkin {
				// 掩码：皮肤区域为白色
				mask.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
				// 结果图：皮肤保留原色
				result.SetRGBA(x, y, color.RGBA{r8, g8, b8, a8})
			} else {
				// 掩码：非皮肤区域为黑色
				mask.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
				// 结果图：非皮肤区域为灰色（取亮度值）
				gray := uint8(yVal) // 用Y亮度值作为灰色
				result.SetRGBA(x, y, color.RGBA{gray, gray, gray, a8})
			}
		}
	}

	// 4. 形态学后处理（去除小噪声，可选）
	processedMask := morphCleanup(mask, width, height)

	return result, processedMask
}

// RGB转YCrCb颜色空间（ITU-R BT.601标准）
func rgbToYCrCb(r, g, b uint8) (y, cr, cb uint8) {
	// 转换公式（RGB值范围0-255）
	y = uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) + 0.5)
	cr = uint8(128 + 0.713*(float64(r)-float64(y)) + 0.5)
	cb = uint8(128 + 0.564*(float64(b)-float64(y)) + 0.5)
	return y, cr, cb
}

// 形态学后处理：腐蚀+膨胀去除小噪声（简单实现）
func morphCleanup(mask image.Image, width, height int) image.Image {
	// 腐蚀：去除孤立的小亮点（3x3窗口）
	eroded := erode(mask, width, height)
	// 膨胀：恢复皮肤区域的连通性
	dilated := dilate(eroded, width, height)
	return dilated
}

// 腐蚀操作
func erode(mask image.Image, width, height int) image.Image {
	result := image.NewRGBA(mask.Bounds())
	half := 1 // 3x3窗口半宽

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 3x3窗口内是否全为皮肤（白色）
			allSkin := true
			for ky := -half; ky <= half; ky++ {
				for kx := -half; kx <= half; kx++ {
					nx := x + kx
					ny := y + ky
					if nx < 0 || nx >= width || ny < 0 || ny >= height {
						allSkin = false
						break
					}
					r, _, _, _ := mask.At(nx, ny).RGBA()
					if r>>8 < 255 { // 非白色（非皮肤）
						allSkin = false
						break
					}
				}
				if !allSkin {
					break
				}
			}
			if allSkin {
				result.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return result
}

// 膨胀操作
func dilate(mask image.Image, width, height int) image.Image {
	result := image.NewRGBA(mask.Bounds())
	half := 1 // 3x3窗口半宽

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 3x3窗口内是否有皮肤（白色）
			hasSkin := false
			for ky := -half; ky <= half; ky++ {
				for kx := -half; kx <= half; kx++ {
					nx := x + kx
					ny := y + ky
					if nx < 0 || nx >= width || ny < 0 || ny >= height {
						continue
					}
					r, _, _, _ := mask.At(nx, ny).RGBA()
					if r>>8 == 255 { // 白色（皮肤）
						hasSkin = true
						break
					}
				}
				if hasSkin {
					break
				}
			}
			if hasSkin {
				result.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				result.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return result
}

// ========================================================================

// case13 基于高斯模型的肤色概率计算方法

// 基于高斯模型的肤色概率计算方法通过统计肤色在特定颜色空间中的分布特性，建立概率模型来量化每个像素属于肤色的可能性。
//与传统阈值法的 “非黑即白” 判断不同，该方法能输出连续的概率值（0-1），更适合复杂场景（如光照变化、肤色差异），是肤色检测中更鲁棒的方案。

// 算法原理
//颜色空间选择：通常使用 YCrCb 颜色空间，重点关注 Cr（红色色度）和 Cb（蓝色色度）分量（二者对肤色的区分度高，且受亮度影响小）。
//高斯模型假设：假设肤色的 Cr 和 Cb 值服从二维正态分布（高斯分布），其概率密度函数（PDF）由均值向量（μ）和协方差矩阵（Σ）决定。
//参数估计：通过大量肤色样本（训练集）计算 μ 和 Σ：
//均值向量 μ = [μ_Cr, μ_Cb]（Cr 和 Cb 的样本均值）
//协方差矩阵 Σ（描述 Cr 和 Cb 的相关性及离散程度）
//概率计算：对输入图像的每个像素，计算其 Cr、Cb 值在该高斯模型下的概率密度，作为肤色概率（值越高，越可能是肤色）。

func case13() {
	inputPath := "test4.jpg"               // 输入图像
	resultPath := "output_case13_skin.jpg" // 概率可视化结果
	maskPath := "output_case13_mask.jpg"   // 二值掩码（阈值0.5）

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

	// 计算肤色概率并生成结果
	probMap, resultImg := GaussianSkinProbability(img)

	// 生成二值掩码（阈值0.5，可根据需求调整）
	maskImg := probMapToMask(probMap, img.Bounds().Max.X, img.Bounds().Max.Y, 0.5)

	// 保存结果
	outputFile, err := os.Create(resultPath)
	if err != nil {
		panic("无法创建结果文件: " + err.Error())
	}
	defer outputFile.Close()
	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})

	maskFile, err := os.Create(maskPath)
	if err != nil {
		panic("无法创建掩码文件: " + err.Error())
	}
	defer maskFile.Close()
	jpeg.Encode(maskFile, maskImg, &jpeg.Options{Quality: 95})

	println("高斯模型肤色概率计算完成！")
	println("概率可视化结果:", resultPath)
	println("二值掩码（阈值0.5）:", maskPath)
}

// 高斯模型参数
type GaussianParams struct {
	meanCr float64       // Cr分量均值
	meanCb float64       // Cb分量均值
	cov    [2][2]float64 // 协方差矩阵
	covInv [2][2]float64 // 协方差矩阵的逆
	covDet float64       // 协方差矩阵的行列式
}

// 初始化高斯模型参数
func initGaussianParams() GaussianParams {
	meanCr := 152.0
	meanCb := 107.0

	cov := [2][2]float64{
		{160.0, 20.0},
		{20.0, 140.0},
	}

	covDet := cov[0][0]*cov[1][1] - cov[0][1]*cov[1][0]
	invDet := 1.0 / covDet
	covInv := [2][2]float64{
		{cov[1][1] * invDet, -cov[0][1] * invDet},
		{-cov[1][0] * invDet, cov[0][0] * invDet},
	}

	return GaussianParams{
		meanCr: meanCr,
		meanCb: meanCb,
		cov:    cov,
		covInv: covInv,
		covDet: covDet,
	}
}

// 计算二维高斯分布的概率密度
func gaussianPDF(cr, cb float64, params GaussianParams) float64 {
	dCr := cr - params.meanCr
	dCb := cb - params.meanCb

	mahalanobis := dCr*(dCr*params.covInv[0][0]+dCb*params.covInv[0][1]) +
		dCb*(dCr*params.covInv[1][0]+dCb*params.covInv[1][1])

	normalizer := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	pdf := normalizer * math.Exp(-0.5*mahalanobis)

	peakPDF := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	return pdf / peakPDF
}

// 基于高斯模型的肤色概率计算主函数（修正变量类型）
func GaussianSkinProbability(input image.Image) (probMap [][]float64, resultImg *image.RGBA) {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 关键修正：将resultImg声明为具体类型*image.RGBA（而非image.Image接口）
	resultImg = image.NewRGBA(bounds)

	params := initGaussianParams()
	probMap = make([][]float64, height)
	for y := 0; y < height; y++ {
		probMap[y] = make([]float64, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			_, cr, cb := rgbToYCrCb(r8, g8, b8)

			prob := gaussianPDF(float64(cr), float64(cb), params)
			probMap[y][x] = prob

			// 现在可以正常调用SetRGBA，因为resultImg是*image.RGBA类型
			if prob > 0.5 {
				alpha := prob
				resultR := uint8(float64(r8)*alpha + 255*(1-alpha)*0.3)
				resultG := uint8(float64(g8)*alpha + 255*(1-alpha)*0.3)
				resultB := uint8(float64(b8)*alpha + 255*(1-alpha)*0.3)
				resultImg.SetRGBA(x, y, color.RGBA{resultR, resultG, resultB, a8})
			} else {
				gray := uint8(50 + 155*(1-prob))
				resultImg.SetRGBA(x, y, color.RGBA{gray, gray, gray, a8})
			}
		}
	}

	return probMap, resultImg
}

// 概率图转二值掩码（修正mask类型）
func probMapToMask(probMap [][]float64, width, height int, threshold float64) *image.RGBA {
	// 修正：返回具体类型*image.RGBA
	mask := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if probMap[y][x] >= threshold {
				mask.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				mask.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return mask
}

// ========================================================================

// case13_1  应用，将检测的皮肤跟换肤色，比如换成黑褐色变成黑人

func case13_1() {
	inputPath := "test4.jpg"                             // 输入图像路径
	outputPath := "output_case13_1_darkbrown_result.jpg" // 转换结果路径

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

	// 执行皮肤转黑褐色处理
	resultImg := SkinToDarkBrown(img)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
	println("皮肤转黑褐色处理完成！")
	println("输出图片:", outputPath)
}

// 初始化高斯模型参数（通用肤色分布）
func initGaussianParams1() GaussianParams {
	return GaussianParams{
		meanCr: 152.0,
		meanCb: 107.0,
		cov:    [2][2]float64{{160.0, 20.0}, {20.0, 140.0}},
		covInv: [2][2]float64{{0.0063, -0.0009}, {-0.0009, 0.0071}}, // 预计算的逆矩阵
		covDet: 160.0*140.0 - 20.0*20.0,                             // 行列式=22000
	}
}

// 目标黑褐色（可根据需求调整深浅）
var (
	targetR = uint8(40) // 黑褐色红色分量
	targetG = uint8(25) // 黑褐色绿色分量
	targetB = uint8(15) // 黑褐色蓝色分量
)

// 计算肤色概率（基于高斯模型）
func gaussianPDF1(cr, cb float64, params GaussianParams) float64 {
	dCr := cr - params.meanCr
	dCb := cb - params.meanCb

	mahalanobis := dCr*(dCr*params.covInv[0][0]+dCb*params.covInv[0][1]) +
		dCb*(dCr*params.covInv[1][0]+dCb*params.covInv[1][1])

	normalizer := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	pdf := normalizer * math.Exp(-0.5*mahalanobis)
	peakPDF := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	return pdf / peakPDF // 归一化到[0,1]
}

// 核心功能：检测皮肤并转换为黑褐色
func SkinToDarkBrown(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	resultImg := image.NewRGBA(bounds)
	params := initGaussianParams1()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 读取原始像素值
			r, g, b, a := input.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 转换为YCrCb并计算肤色概率
			_, cr, cb := rgbToYCrCb(r8, g8, b8)
			prob := gaussianPDF1(float64(cr), float64(cb), params)

			// 根据肤色概率调整颜色（概率越高，越接近黑褐色）
			var finalR, finalG, finalB uint8
			if prob > 0.1 { // 低概率阈值，确保弱皮肤区域也被处理
				// 混合比例：皮肤概率越高，目标色占比越大（0.7-1.0）
				mixRatio := math.Min(1.0, prob*1.2) // 增强转换强度

				// 线性混合：原肤色与黑褐色过渡
				finalR = uint8(float64(r8)*(1-mixRatio) + float64(targetR)*mixRatio)
				finalG = uint8(float64(g8)*(1-mixRatio) + float64(targetG)*mixRatio)
				finalB = uint8(float64(b8)*(1-mixRatio) + float64(targetB)*mixRatio)
			} else {
				// 非皮肤区域保留原色
				finalR, finalG, finalB = r8, g8, b8
			}

			// 设置结果像素（保留Alpha通道）
			resultImg.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, a8})
		}
	}

	return resultImg
}

// ========================================================================

// case13_2 应用，将检测的皮肤跟换肤色，比如换成黑褐色变成黑人 - 脸部阴影部分更加真实

// 要让脸部阴影部分更真实，需要结合光照特性和肤色自然过渡规律：阴影区域不仅亮度更低，还会保留肤色的基础色调（不会完全失去色彩信息），
//且阴影与高光区域的边界应平滑过渡（避免生硬断层）

func case13_2() {
	inputPath := "test4.jpg"                           // 输入图像
	outputPath := "output_case13_2_realistic_skin.jpg" // 带真实阴影的结果

	file, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图片: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("无法解码图片: " + err.Error())
	}

	resultImg := FastSkinConversion(img)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
	println("带真实阴影的肤色转换完成！")
	println("输出图片:", outputPath)
}

// 高斯肤色模型参数（保持原特性）
type GaussianParams3 struct {
	meanCr float64       // Cr均值
	meanCb float64       // Cb均值
	covInv [2][2]float64 // 协方差逆矩阵
	covDet float64       // 协方差行列式
}

// 初始化肤色模型
func initGaussianParams3() GaussianParams3 {
	return GaussianParams3{
		meanCr: 145.0,
		meanCb: 95.0,
		covInv: [2][2]float64{{0.007, -0.001}, {-0.001, 0.008}},
		covDet: 18000.0,
	}
}

// 黑人肤色基准值
var (
	baseR = uint8(40)
	baseG = uint8(25)
	baseB = uint8(15)
)

// 核心参数（精简且高效）
const (
	shadowBrightnessThresh = 85
	shadowDeepenRatio      = 0.25
	smoothFactor           = 0.6
	colorBlendPower        = 0.8
	gaussianKernelSize     = 3 // 保持3x3核，平衡速度与效果
	parallelBlocks         = 4 // 并行处理的图像块数量
)

// 预计算3x3高斯核（全局复用，避免重复生成）
var gaussianKernel = [3][3]float64{
	{0.0751, 0.1238, 0.0751},
	{0.1238, 0.2042, 0.1238},
	{0.0751, 0.1238, 0.0751},
}

// 核心函数：高性能肤色转换
func FastSkinConversion(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	resultImg := image.NewRGBA(bounds)
	params := initGaussianParams3()

	// 1. 合并计算：一次遍历同时获取肤色概率图和平滑亮度图（减少1次完整图像遍历）
	probMap, smoothY := precomputeProbAndSmoothY(input, width, height, params)

	// 2. 并行处理初步肤色转换（按行分块，利用多核加速）
	var wg sync.WaitGroup
	rowsPerBlock := (height + parallelBlocks - 1) / parallelBlocks // 每块处理的行数
	for b := 0; b < parallelBlocks; b++ {
		startY := b * rowsPerBlock
		endY := startY + rowsPerBlock
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(bStart, bEnd int) {
			defer wg.Done()
			// 块内处理：初步肤色转换+阴影处理
			for y := bStart; y < bEnd; y++ {
				for x := 0; x < width; x++ {
					r, g, b, a := input.At(x, y).RGBA()
					r8 := uint8(r >> 8)
					g8 := uint8(g >> 8)
					b8 := uint8(b >> 8)
					a8 := uint8(a >> 8)

					prob := probMap[y][x]
					smoothYVal := smoothY[y][x]

					// 非线性混合（复用预计算的probMap，避免重复计算）
					mixRatio := 1 - math.Pow(1-prob, colorBlendPower)
					mixRatio = math.Min(1.0, mixRatio*1.05)

					// 基础肤色混合
					finalR := uint8(float64(r8)*(1-mixRatio) + float64(baseR)*mixRatio)
					finalG := uint8(float64(g8)*(1-mixRatio) + float64(baseG)*mixRatio)
					finalB := uint8(float64(b8)*(1-mixRatio) + float64(baseB)*mixRatio)

					// 阴影处理（精简计算逻辑）
					if prob > 0.2 && smoothYVal < shadowBrightnessThresh {
						shadowDepth := 1 - math.Pow(float64(smoothYVal)/shadowBrightnessThresh, 1.2)
						deepen := shadowDepth * shadowDeepenRatio
						finalR = uint8(float64(finalR) * (1 - deepen*0.7))
						finalG = uint8(float64(finalG) * (1 - deepen*0.9))
						finalB = uint8(float64(finalB) * (1 - deepen*0.85))
					}

					resultImg.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, a8})
				}
			}
		}(startY, endY)
	}
	wg.Wait() // 等待并行处理完成

	// 3. 高效高斯平滑（仅对皮肤区域，且减少边界判断）
	smoothed := applyFastSkinSmooth(resultImg, probMap, width, height)

	return smoothed
}

// 合并计算：一次遍历同时生成肤色概率图和平滑亮度图（减少图像访问次数）
func precomputeProbAndSmoothY(img image.Image, width, height int, params GaussianParams3) ([][]float64, [][]float64) {
	probMap := make([][]float64, height)
	smoothY := make([][]float64, height)
	halfKernel := gaussianKernelSize / 2

	for y := 0; y < height; y++ {
		probMap[y] = make([]float64, width)
		smoothY[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 一次获取当前像素的RGB值，同时用于概率计算和亮度计算
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			_, cr, cb := rgbToYCrCb(r8, g8, b8)

			// 计算肤色概率（复用当前像素的Cr/Cb）
			probMap[y][x] = gaussianPDF3(float64(cr), float64(cb), params)

			// 计算高斯平滑亮度（预计算边界范围，减少循环内判断）
			minX, maxX := max(0, x-halfKernel), min(width-1, x+halfKernel)
			minY, maxY := max(0, y-halfKernel), min(height-1, y+halfKernel)
			var weightedSum, kernelSum float64

			for ky := minY; ky <= maxY; ky++ {
				for kx := minX; kx <= maxX; kx++ {
					// 计算核索引（避免负数，直接映射）
					ki := ky - (y - halfKernel)
					kj := kx - (x - halfKernel)
					weight := gaussianKernel[ki][kj]

					// 复用邻域像素的Y值计算（避免重复调用At和转换）
					nr, ng, nb, _ := img.At(kx, ky).RGBA()
					nyVal, _, _ := rgbToYCrCb(uint8(nr>>8), uint8(ng>>8), uint8(nb>>8))
					weightedSum += float64(nyVal) * weight
					kernelSum += weight
				}
			}
			smoothY[y][x] = weightedSum / kernelSum
		}
	}
	return probMap, smoothY
}

// 高效皮肤区域平滑（减少边界判断，复用概率图）
func applyFastSkinSmooth(img *image.RGBA, probMap [][]float64, width, height int) *image.RGBA {
	smoothed := image.NewRGBA(img.Bounds())
	halfKernel := gaussianKernelSize / 2

	// 并行处理平滑（再次分块加速）
	var wg sync.WaitGroup
	rowsPerBlock := (height + parallelBlocks - 1) / parallelBlocks
	for b := 0; b < parallelBlocks; b++ {
		startY := b * rowsPerBlock
		endY := startY + rowsPerBlock
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(bStart, bEnd int) {
			defer wg.Done()
			for y := bStart; y < bEnd; y++ {
				for x := 0; x < width; x++ {
					prob := probMap[y][x]
					origR, origG, origB, origA := img.At(x, y).RGBA()
					origR8 := uint8(origR >> 8)
					origG8 := uint8(origG >> 8)
					origB8 := uint8(origB >> 8)
					origA8 := uint8(origA >> 8)

					if prob < 0.1 { // 非皮肤区域直接复制
						smoothed.SetRGBA(x, y, color.RGBA{origR8, origG8, origB8, origA8})
						continue
					}

					// 预计算邻域范围，减少循环内条件判断
					minX, maxX := max(0, x-halfKernel), min(width-1, x+halfKernel)
					minY, maxY := max(0, y-halfKernel), min(height-1, y+halfKernel)
					var rSum, gSum, bSum, aSum, weightSum float64

					for ky := minY; ky <= maxY; ky++ {
						for kx := minX; kx <= maxX; kx++ {
							nprob := probMap[ky][kx]
							ki := ky - (y - halfKernel)
							kj := kx - (x - halfKernel)
							weight := gaussianKernel[ki][kj] * math.Min(prob, nprob)
							if weight < 0.01 {
								continue
							}

							nr, ng, nb, na := img.At(kx, ky).RGBA()
							rSum += float64(nr>>8) * weight
							gSum += float64(ng>>8) * weight
							bSum += float64(nb>>8) * weight
							aSum += float64(na>>8) * weight
							weightSum += weight
						}
					}

					// 混合平滑结果与原始像素（使用smoothFactor）
					if weightSum > 0 {
						smoothR := uint8(rSum / weightSum)
						smoothG := uint8(gSum / weightSum)
						smoothB := uint8(bSum / weightSum)
						smoothA := uint8(aSum / weightSum)

						finalR := uint8(float64(origR8)*(1-smoothFactor) + float64(smoothR)*smoothFactor)
						finalG := uint8(float64(origG8)*(1-smoothFactor) + float64(smoothG)*smoothFactor)
						finalB := uint8(float64(origB8)*(1-smoothFactor) + float64(smoothB)*smoothFactor)
						finalA := uint8(float64(origA8)*(1-smoothFactor) + float64(smoothA)*smoothFactor)

						smoothed.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, finalA})
					} else {
						smoothed.SetRGBA(x, y, color.RGBA{origR8, origG8, origB8, origA8})
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return smoothed
}

// 辅助函数：高斯概率计算（保持不变）
func gaussianPDF3(cr, cb float64, params GaussianParams3) float64 {
	dCr := cr - params.meanCr
	dCb := cb - params.meanCb
	mahalanobis := dCr*(dCr*params.covInv[0][0]+dCb*params.covInv[0][1]) +
		dCb*(dCr*params.covInv[1][0]+dCb*params.covInv[1][1])
	normalizer := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	pdf := normalizer * math.Exp(-0.5*mahalanobis)
	peakPDF := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	return math.Min(1.0, pdf/peakPDF)
}

// ========================================================================

// case14 皮肤美白算法 LUT调色法

// 基于 LUT（Lookup Table，查找表）调色法的皮肤美白算法，通过预定义颜色映射关系实现高效的肤色调整：
//先设计针对皮肤区域的亮度 / 色度映射曲线（LUT），再通过查找表快速将输入颜色映射到美白后的颜色，避免实时复杂计算，兼顾效率与自然度

// 算法原理
//颜色空间选择：使用 YCrCb 颜色空间，重点调整 Y（亮度）通道（提亮肤色），同时微调 Cr（红色色度）和 Cb（蓝色色度）通道（避免美白后肤色偏黄 / 偏红）。
//LUT 设计：
//Y 通道 LUT：对皮肤区域的 Y 值（亮度）设计 “低亮度大幅提亮，高亮度小幅提亮” 的非线性曲线（避免过曝）。
//Cr/Cb 通道 LUT：轻微压低 Cr 值（减少红色感）、微调 Cb 值（保持肤色通透）。
//非皮肤区域：LUT 设为 “输入 = 输出”，不改变原始颜色。
//肤色掩码融合：结合肤色检测得到的概率图，按皮肤概率动态混合原始颜色与 LUT 映射颜色（概率越高，美白效果越强）

func case14() {
	inputPath := "test4.jpg"          // 输入图像路径
	outputPath := "output_case14.jpg" // 美白结果路径

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 执行LUT皮肤美白
	whitenedImg := LUTSkinWhitening(inputImg)

	// 保存美白结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	// 以高质量保存JPEG
	jpeg.Encode(outputFile, whitenedImg, &jpeg.Options{Quality: 95})
	println("LUT皮肤美白处理完成！")
	println("美白结果已保存至:", outputPath)
}

// 皮肤检测的高斯模型参数（用于区分皮肤区域）
type SkinGaussianParams struct {
	meanCr float64       // 皮肤Cr分量均值
	meanCb float64       // 皮肤Cb分量均值
	covInv [2][2]float64 // 协方差矩阵的逆
	covDet float64       // 协方差矩阵的行列式
}

// 初始化皮肤检测的高斯模型参数（适配黄/白种人皮肤）
func initSkinGaussianParams() SkinGaussianParams {
	return SkinGaussianParams{
		meanCr: 155.0, // 黄/白种人皮肤Cr值偏高（偏红）
		meanCb: 100.0, // 皮肤Cb值中等
		covInv: [2][2]float64{{0.008, -0.001}, {-0.001, 0.009}},
		covDet: 20000.0, // 皮肤Cr/Cb分布的行列式
	}
}

// LUT美白核心参数（控制美白效果）
const (
	skinWhitenIntensity = 0.7  // 美白强度（0-1，值越高美白越明显）
	maxSkinY            = 240  // 皮肤最大亮度（避免过曝）
	crReduceRatio       = 0.1  // Cr通道降低比例（减少红色感）
	cbEnhanceRatio      = 0.05 // Cb通道增强比例（提升通透感）
)

// 预生成皮肤美白的LUT（查找表）：启动时计算，加速实时处理
var (
	skinYWhiteningLUT [256]uint8 // Y通道（亮度）美白映射表
	skinCrAdjustLUT   [256]uint8 // Cr通道（红色度）调整表
	skinCbAdjustLUT   [256]uint8 // Cb通道（蓝色度）调整表
)

// 初始化LUT表（程序启动时执行，预计算颜色映射关系）
func init() {
	// 1. Y通道LUT：低亮度皮肤大幅提亮，高亮度皮肤小幅提亮（非线性美白）
	for y := 0; y < 256; y++ {
		// 非线性因子：暗部（y小）提亮多，亮部（y大）提亮少，避免高光过曝
		nonlinearFactor := 1.0 - math.Pow(float64(y)/255.0, 2)
		// 计算美白后的亮度：原始亮度 + 可提升空间 * 强度 * 非线性因子
		whitenedY := float64(y) + (maxSkinY-float64(y))*skinWhitenIntensity*nonlinearFactor
		if whitenedY > maxSkinY {
			whitenedY = maxSkinY // 限制最大亮度，防止过曝
		}
		skinYWhiteningLUT[y] = uint8(whitenedY)
	}

	// 2. Cr通道LUT：降低红色度（避免美白后皮肤偏红）
	for cr := 0; cr < 256; cr++ {
		adjustedCr := float64(cr) * (1 - crReduceRatio*skinWhitenIntensity)
		skinCrAdjustLUT[cr] = clampColor(adjustedCr)
	}

	// 3. Cb通道LUT：轻微提升蓝色度（增加皮肤通透感）
	for cb := 0; cb < 256; cb++ {
		adjustedCb := float64(cb) * (1 + cbEnhanceRatio*skinWhitenIntensity)
		skinCbAdjustLUT[cb] = clampColor(adjustedCb)
	}
}

// 核心函数：基于LUT的皮肤美白主流程
func LUTSkinWhitening(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	resultImg := image.NewRGBA(bounds)
	skinGaussian := initSkinGaussianParams()

	// 1. 预计算皮肤概率图（每个像素属于皮肤的概率，0-1）
	skinProbMap := precomputeSkinProbMap(input, width, height, skinGaussian)

	// 2. 并行处理图像（按行分块，利用多核加速）
	var wg sync.WaitGroup
	rowsPerBlock := (height + 3) / 4 // 分成4块并行处理
	for block := 0; block < 4; block++ {
		startY := block * rowsPerBlock
		endY := startY + rowsPerBlock
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(blockStartY, blockEndY int) {
			defer wg.Done()
			processImageBlock(input, resultImg, skinProbMap, blockStartY, blockEndY, width)
		}(startY, endY)
	}
	wg.Wait() // 等待所有块处理完成

	return resultImg
}

// 预计算皮肤概率图（区分皮肤和非皮肤区域）
func precomputeSkinProbMap(img image.Image, width, height int, gaussian SkinGaussianParams) [][]float64 {
	skinProbMap := make([][]float64, height)
	for y := 0; y < height; y++ {
		skinProbMap[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 转换为YCrCb，提取Cr/Cb用于皮肤概率计算
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			_, cr, cb := rgbToYCrCbForSkin(r8, g8, b8)
			// 计算当前像素的皮肤概率
			skinProbMap[y][x] = skinGaussianProbability(float64(cr), float64(cb), gaussian)
		}
	}
	return skinProbMap
}

// 处理单个图像块（应用LUT美白）
func processImageBlock(input image.Image, result *image.RGBA, skinProbMap [][]float64, startY, endY, width int) {
	for y := startY; y < endY; y++ {
		for x := 0; x < width; x++ {
			// 读取原始像素并转换为YCrCb
			r, g, b, a := input.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			a8 := uint8(a >> 8)
			yVal, cr, cb := rgbToYCrCbForSkin(r8, g8, b8)

			// 获取当前像素的皮肤概率（0-1）
			skinProb := skinProbMap[y][x]

			// 按皮肤概率混合：原始值与LUT美白值
			// Y通道（亮度）：美白核心，权重随皮肤概率变化
			whitenedY := skinYWhiteningLUT[yVal]
			finalY := blendBySkinProb(yVal, whitenedY, skinProb)

			// Cr通道（红色度）：轻微调整，权重低于亮度
			adjustedCr := skinCrAdjustLUT[cr]
			finalCr := blendBySkinProb(cr, adjustedCr, skinProb*0.8)

			// Cb通道（蓝色度）：微调通透感，权重最低
			adjustedCb := skinCbAdjustLUT[cb]
			finalCb := blendBySkinProb(cb, adjustedCb, skinProb*0.5)

			// 转换回RGB并写入结果
			finalR, finalG, finalB := yCrCbToRGBForSkin(finalY, finalCr, finalCb)
			result.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, a8})
		}
	}
}

// 计算皮肤的高斯概率（判断像素是否为皮肤）
func skinGaussianProbability(cr, cb float64, params SkinGaussianParams) float64 {
	dCr := cr - params.meanCr // Cr与均值的偏差
	dCb := cb - params.meanCb // Cb与均值的偏差
	// 计算马氏距离的平方（衡量与皮肤分布的偏离程度）
	mahalanobisDistSq := dCr*(dCr*params.covInv[0][0]+dCb*params.covInv[0][1]) +
		dCb*(dCr*params.covInv[1][0]+dCb*params.covInv[1][1])
	// 高斯概率密度公式
	normalizer := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	probDensity := normalizer * math.Exp(-0.5*mahalanobisDistSq)
	// 归一化到[0,1]范围
	peakDensity := 1.0 / (2 * math.Pi * math.Sqrt(params.covDet))
	return math.Min(1.0, probDensity/peakDensity)
}

// 按皮肤概率混合两个颜色值（原始值与调整值）
func blendBySkinProb(original, adjusted uint8, skinProb float64) uint8 {
	return uint8(float64(original)*(1-skinProb) + float64(adjusted)*skinProb)
}

// 将颜色值限制在[0,255]范围内
func clampColor(value float64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value + 0.5) // 四舍五入
}

// RGB转YCrCb（用于皮肤处理的颜色空间转换）
func rgbToYCrCbForSkin(r, g, b uint8) (y, cr, cb uint8) {
	y = uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) + 0.5)
	cr = uint8(128 + 0.713*(float64(r)-float64(y)) + 0.5)
	cb = uint8(128 + 0.564*(float64(b)-float64(y)) + 0.5)
	return y, cr, cb
}

// YCrCb转RGB（皮肤处理后转回显示用的RGB）
func yCrCbToRGBForSkin(y, cr, cb uint8) (r, g, b uint8) {
	yFloat := float64(y)
	crFloat := float64(cr) - 128.0 // 还原Cr偏移
	cbFloat := float64(cb) - 128.0 // 还原Cb偏移

	// 转换公式（确保在有效范围）
	r = clampColor(yFloat + 1.403*crFloat)
	g = clampColor(yFloat - 0.344*cbFloat - 0.714*crFloat)
	b = clampColor(yFloat + 1.773*cbFloat)
	return r, g, b
}

// ========================================================================

// case15 人像磨皮 通用磨皮算法

// 通用人像磨皮算法的核心是在平滑皮肤纹理（如毛孔、痘印、细纹）的同时保留五官边缘（如眉毛、眼线、唇线、发丝），避免 “模糊成一团” 的不自然效果。
//常用的边缘保留滤波（如双边滤波、表面模糊）是实现该效果的基础，结合皮肤区域检测可进一步提升针对性（只磨皮不影响背景）

//算法原理
//皮肤区域检测：基于 YCrCb 颜色空间的高斯模型，生成皮肤概率图（0-1），区分皮肤与非皮肤区域（头发、背景、衣物等）。
//边缘保留滤波：使用双边滤波对图像进行处理 —— 对颜色差异小的区域（皮肤内部）进行模糊，对颜色差异大的区域（边缘）保留细节，避免边缘模糊。
//动态混合：根据皮肤概率和边缘强度，混合原始图像与滤波结果：
//高皮肤概率 + 低边缘强度（如脸颊）：更多保留滤波结果（磨皮强）。
//低皮肤概率 + 高边缘强度（如眉毛、发丝）：更多保留原始图像（不磨皮）

func case15() {
	inputPath := "test4.jpg"          // 输入人像路径
	outputPath := "output_case15.jpg" // 肤质柔化结果路径

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 执行人像肤质柔化
	softenedImg := PortraitSkinSoftening(inputImg)

	// 保存柔化结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, softenedImg, &jpeg.Options{Quality: 95})
	println("人像肤质柔化处理完成！")
	println("结果保存至:", outputPath)
}

// 人像肤质检测的高斯模型参数（区分肤质区域与非肤质区域）
type PortraitSkinGaussian struct {
	meanLum  float64       // 肤质区域亮度均值
	meanRed  float64       // 肤质Cr分量均值（偏红特征）
	meanBlue float64       // 肤质Cb分量均值（偏蓝特征）
	covInv   [3][3]float64 // 3D协方差逆矩阵（亮度+Cr+Cb）
	covDet   float64       // 协方差行列式
}

// 初始化人像肤质检测模型（适配多肤色）
func initPortraitSkinGaussian() PortraitSkinGaussian {
	return PortraitSkinGaussian{
		meanLum:  150.0, // 典型肤质亮度均值
		meanRed:  155.0, // 肤质红色度均值
		meanBlue: 100.0, // 肤质蓝色度均值
		covInv: [3][3]float64{ // 预计算的3D协方差逆矩阵
			{0.0005, 0.0001, -0.0001},
			{0.0001, 0.008, -0.001},
			{-0.0001, -0.001, 0.009},
		},
		covDet: 350000.0, // 3D协方差行列式值
	}
}

// 肤质柔化核心参数（可调节柔化效果）
const (
	skinSoftIntensity   = 0.8  // 肤质柔化强度（0-1，越高柔化越明显）
	spatialBlurSigma    = 2.0  // 空间模糊范围（控制柔化区域大小）
	colorSensitivity    = 30.0 // 颜色敏感度（控制边缘保留严格度）
	edgeIntensityThresh = 20.0 // 边缘强度阈值（高于此值保留原始细节）
	parallelChunkCount  = 4    // 并行处理的图像块数量
)

// 核心函数：人像肤质柔化主流程
func PortraitSkinSoftening(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	srcImg := image.NewRGBA(bounds) // 转换为RGBA格式便于像素操作

	// 复制输入图像到RGBA（统一处理格式）
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			srcImg.SetRGBA(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8),
			})
		}
	}

	// 1. 计算人像肤质概率图（区分肤质区域与非肤质区域）
	portraitSkinProb := calcPortraitSkinProb(srcImg, width, height)

	// 2. 对原图应用边缘保留模糊（柔化肤质但不模糊边缘）
	blurredImg := applyEdgePreservingBlur(srcImg, width, height)

	// 3. 计算图像边缘强度图（用于保留五官、发丝等关键边缘）
	edgeIntensityMap := calcEdgeIntensity(srcImg, width, height)

	// 4. 动态融合：结合肤质概率和边缘强度，生成最终柔化结果
	resultImg := fuseSoftAndEdge(srcImg, blurredImg, portraitSkinProb, edgeIntensityMap, width, height)

	return resultImg
}

// 计算人像肤质概率图（0-1，1表示高概率肤质区域）
func calcPortraitSkinProb(img *image.RGBA, width, height int) [][]float64 {
	skinProb := make([][]float64, height)
	gaussianModel := initPortraitSkinGaussian()

	for y := 0; y < height; y++ {
		skinProb[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			lum, cr, cb := rgbToYCrCb15(r8, g8, b8) // 亮度+色度转换

			// 计算当前像素与肤质模型的匹配概率
			prob := portraitSkinGaussianProb(
				float64(lum), float64(cr), float64(cb),
				gaussianModel,
			)
			skinProb[y][x] = prob
		}
	}
	return skinProb
}

// 应用边缘保留模糊（双边滤波变种，柔化肤质同时保留边缘）
func applyEdgePreservingBlur(src *image.RGBA, width, height int) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	halfKernel := int(spatialBlurSigma * 2.0) // 模糊核大小（约为2*sigma*2）

	var wg sync.WaitGroup
	rowsPerChunk := (height + parallelChunkCount - 1) / parallelChunkCount

	for chunk := 0; chunk < parallelChunkCount; chunk++ {
		startY := chunk * rowsPerChunk
		endY := startY + rowsPerChunk
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(cStart, cEnd int) {
			defer wg.Done()
			for y := cStart; y < cEnd; y++ {
				for x := 0; x < width; x++ {
					// 双边滤波：空间权重（距离）+ 颜色权重（相似度）
					var rSum, gSum, bSum, weightTotal float64
					for ky := -halfKernel; ky <= halfKernel; ky++ {
						for kx := -halfKernel; kx <= halfKernel; ky++ {
							nx, ny := x+kx, y+ky
							if nx < 0 || nx >= width || ny < 0 || ny >= height {
								continue
							}

							// 空间权重（高斯距离，控制模糊范围）
							spatialDist := float64(kx*kx + ky*ky)
							spatialW := math.Exp(-spatialDist / (2 * spatialBlurSigma * spatialBlurSigma))

							// 颜色权重（RGB差异，控制边缘敏感度）
							srcR, srcG, srcB, _ := src.At(x, y).RGBA()
							nR, nG, nB, _ := src.At(nx, ny).RGBA()
							colorDist := float64(
								(srcR-nR)*(srcR-nR) +
									(srcG-nG)*(srcG-nG) +
									(srcB-nB)*(srcB-nB),
							)
							colorW := math.Exp(-colorDist / (2 * colorSensitivity * colorSensitivity * 65536)) // 65536为RGBA最大值平方

							// 总权重与加权求和
							totalW := spatialW * colorW
							weightTotal += totalW
							rSum += float64(nR>>8) * totalW
							gSum += float64(nG>>8) * totalW
							bSum += float64(nB>>8) * totalW
						}
					}

					// 归一化并设置像素值
					if weightTotal > 0 {
						r := uint8(rSum / weightTotal)
						g := uint8(gSum / weightTotal)
						b := uint8(bSum / weightTotal)
						_, _, _, a := src.At(x, y).RGBA()
						dst.SetRGBA(x, y, color.RGBA{r, g, b, uint8(a >> 8)})
					} else {
						dst.SetRGBA(x, y, src.At(x, y).(color.RGBA))
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return dst
}

// 计算图像边缘强度（基于Sobel算子，用于识别五官边缘）
func calcEdgeIntensity(img *image.RGBA, width, height int) [][]float64 {
	edgeMap := make([][]float64, height)
	// Sobel算子（水平/垂直梯度模板）
	sobelX := [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}} // 水平梯度
	sobelY := [3][3]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}} // 垂直梯度

	for y := 1; y < height-1; y++ {
		edgeMap[y] = make([]float64, width)
		for x := 1; x < width-1; x++ {
			var gradX, gradY int
			// 3x3邻域计算梯度
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					nx, ny := x+kx, y+ky
					r, g, b, _ := img.At(nx, ny).RGBA()
					gray := int(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)) // 转灰度值
					gradX += gray * sobelX[ky+1][kx+1]
					gradY += gray * sobelY[ky+1][kx+1]
				}
			}
			// 梯度强度（边缘强度）
			edgeStrength := math.Sqrt(float64(gradX*gradX + gradY*gradY))
			edgeMap[y][x] = edgeStrength
		}
	}
	// 边界像素边缘强度设为0（简化处理）
	for x := 0; x < width; x++ {
		edgeMap[0][x] = 0
		edgeMap[height-1][x] = 0
	}
	for y := 0; y < height; y++ {
		edgeMap[y][0] = 0
		edgeMap[y][width-1] = 0
	}
	return edgeMap
}

// 融合原始图像与柔化结果（基于肤质概率和边缘强度）
func fuseSoftAndEdge(
	src, blurred *image.RGBA,
	skinProb [][]float64,
	edgeIntensity [][]float64,
	width, height int,
) *image.RGBA {
	result := image.NewRGBA(src.Bounds())
	var wg sync.WaitGroup
	rowsPerChunk := (height + parallelChunkCount - 1) / parallelChunkCount

	for chunk := 0; chunk < parallelChunkCount; chunk++ {
		startY := chunk * rowsPerChunk
		endY := startY + rowsPerChunk
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(cStart, cEnd int) {
			defer wg.Done()
			for y := cStart; y < cEnd; y++ {
				for x := 0; x < width; x++ {
					// 获取当前像素的肤质概率和边缘强度
					prob := skinProb[y][x]
					edge := edgeIntensity[y][x]

					// 计算柔化权重：肤质概率越高、边缘越弱，柔化结果占比越高
					edgeFactor := math.Max(0, 1-edge/edgeIntensityThresh) // 边缘越强，该值越小
					softWeight := prob * edgeFactor * skinSoftIntensity
					softWeight = math.Min(1.0, softWeight) // 权重限制在[0,1]

					// 融合原始像素与柔化像素
					srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
					blurR, blurG, blurB, _ := blurred.At(x, y).RGBA()

					finalR := uint8(float64(srcR>>8)*(1-softWeight) + float64(blurR>>8)*softWeight)
					finalG := uint8(float64(srcG>>8)*(1-softWeight) + float64(blurG>>8)*softWeight)
					finalB := uint8(float64(srcB>>8)*(1-softWeight) + float64(blurB>>8)*softWeight)

					result.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, uint8(srcA >> 8)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return result
}

// 计算肤质区域的高斯概率（3D：亮度+Cr+Cb）
func portraitSkinGaussianProb(lum, cr, cb float64, model PortraitSkinGaussian) float64 {
	dLum := lum - model.meanLum
	dCr := cr - model.meanRed
	dCb := cb - model.meanBlue

	// 马氏距离平方（衡量与肤质分布的偏离程度）
	mahalanobis := dLum*(dLum*model.covInv[0][0]+dCr*model.covInv[0][1]+dCb*model.covInv[0][2]) +
		dCr*(dLum*model.covInv[1][0]+dCr*model.covInv[1][1]+dCb*model.covInv[1][2]) +
		dCb*(dLum*model.covInv[2][0]+dCr*model.covInv[2][1]+dCb*model.covInv[2][2])

	// 高斯概率密度计算
	normalizer := 1.0 / (math.Pow(2*math.Pi, 1.5) * math.Sqrt(model.covDet))
	pdf := normalizer * math.Exp(-0.5*mahalanobis)
	peakPDF := 1.0 / (math.Pow(2*math.Pi, 1.5) * math.Sqrt(model.covDet))
	return math.Min(1.0, pdf/peakPDF)
}

// RGB转YCrCb（用于肤质检测的颜色空间转换）
func rgbToYCrCb15(r, g, b uint8) (lum, cr, cb uint8) {
	lum = uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) + 0.5)
	cr = uint8(128 + 0.713*(float64(r)-float64(lum)) + 0.5)
	cb = uint8(128 + 0.564*(float64(b)-float64(lum)) + 0.5)
	return lum, cr, cb
}

// ========================================================================

// case16 人像磨皮 通到磨皮算法

//通道磨皮的核心原理是利用图像通道提取皮肤纹理（如毛孔、痘印），通过弱化纹理实现磨皮，同时保留五官边缘，
//常见于专业图像编辑（如 Photoshop 的通道磨皮技巧）。该算法通过分离 “肤色底色” 与 “皮肤纹理”，仅对纹理进行平滑，效果更自然。

//算法原理
//通道分离：提取图像的 RGB 通道（或 Lab 通道的明度通道），聚焦于亮度变化剧烈的通道（通常是红色通道，因皮肤纹理在红色通道更明显）。
//高反差保留：对通道进行高斯模糊，再用原始通道减去模糊结果，得到 “高反差通道”—— 该通道保留了皮肤纹理（高频信息）和边缘，弱化了均匀肤色（低频信息）。
//纹理蒙版生成：将高反差通道转换为灰度蒙版，用于区分 “需要磨皮的纹理区域” 和 “需要保留的边缘区域”。
//混合磨皮：用原始图像与平滑图像（如高斯模糊结果）基于蒙版混合 —— 纹理区域（蒙版亮部）保留更多平滑效果，边缘区域（蒙版暗部）保留更多原始细节

func case16() {
	inputPath := "test4.jpg"          // 输入人像路径
	outputPath := "output_case16.jpg" // 通道磨皮结果路径

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 执行通道磨皮
	smoothedImg := ChannelSkinSmoothing(inputImg)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, smoothedImg, &jpeg.Options{Quality: 95})
	println("通道磨皮处理完成！")
	println("结果保存至:", outputPath)
}

// 通道磨皮核心参数（可调节效果）
const (
	blurRadius         = 3.0 // 高斯模糊半径（控制平滑程度，越大磨皮越强）
	textureSensitivity = 0.3 // 纹理敏感度（0-1，越高保留纹理越多）
	channelToUse       = 0   // 用于磨皮的通道（0=R, 1=G, 2=B，通常选R通道）
	parallelSegments   = 4   // 并行处理的图像段数
)

// 核心函数：通道磨皮主流程
func ChannelSkinSmoothing(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	srcImg := image.NewRGBA(bounds)

	// 转换输入图像为RGBA格式
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			srcImg.SetRGBA(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8),
			})
		}
	}

	// 1. 提取目标通道（如红色通道）
	targetChannel := extractChannel(srcImg, width, height, channelToUse)

	// 2. 对目标通道进行高斯模糊（获取平滑底色）
	blurredChannel := gaussianBlurChannel(targetChannel, width, height, blurRadius)

	// 3. 计算高反差通道（原始通道 - 模糊通道，保留纹理和边缘）
	highPassChannel := calcHighPassChannel(targetChannel, blurredChannel, width, height)

	// 4. 生成纹理蒙版（基于高反差通道，区分纹理与边缘）
	textureMask := generateTextureMask(highPassChannel, width, height)

	// 5. 对原图进行整体平滑（用于混合的平滑图像）
	smoothedImg := gaussianBlurImage(srcImg, width, height, blurRadius)

	// 6. 基于纹理蒙版混合原图与平滑图，实现通道磨皮
	resultImg := blendByTextureMask(srcImg, smoothedImg, textureMask, width, height)

	return resultImg
}

// 提取指定通道（0=R, 1=G, 2=B）
func extractChannel(img *image.RGBA, width, height, channel int) [][]uint8 {
	channelData := make([][]uint8, height)
	for y := 0; y < height; y++ {
		channelData[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			// 根据通道索引选择数据
			switch channel {
			case 0:
				channelData[y][x] = r8
			case 1:
				channelData[y][x] = g8
			case 2:
				channelData[y][x] = b8
			}
		}
	}
	return channelData
}

// 对单通道数据进行高斯模糊
func gaussianBlurChannel(channel [][]uint8, width, height int, radius float64) [][]uint8 {
	blurred := make([][]uint8, height)
	kernel := generateGaussianKernel(radius)
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2

	var wg sync.WaitGroup
	rowsPerSeg := (height + parallelSegments - 1) / parallelSegments

	for seg := 0; seg < parallelSegments; seg++ {
		startY := seg * rowsPerSeg
		endY := startY + rowsPerSeg
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(sStart, sEnd int) {
			defer wg.Done()
			for y := sStart; y < sEnd; y++ {
				blurred[y] = make([]uint8, width)
				for x := 0; x < width; x++ {
					var sum float64
					var weightTotal float64
					// 应用高斯核
					for ky := -halfKernel; ky <= halfKernel; ky++ {
						for kx := -halfKernel; kx <= halfKernel; kx++ {
							nx, ny := x+kx, y+ky
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								weight := kernel[ky+halfKernel][kx+halfKernel]
								sum += float64(channel[ny][nx]) * weight
								weightTotal += weight
							}
						}
					}
					if weightTotal > 0 {
						blurred[y][x] = uint8(sum / weightTotal)
					} else {
						blurred[y][x] = channel[y][x]
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return blurred
}

// 生成高斯核（用于模糊计算）
func generateGaussianKernel(radius float64) [][]float64 {
	size := int(radius*2) + 1
	kernel := make([][]float64, size)
	halfSize := size / 2
	sigma := radius / 2.0 // sigma通常为半径的一半
	sigmaSq := sigma * sigma

	// 计算高斯权重
	var total float64
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := x - halfSize
			dy := y - halfSize
			kernel[y][x] = math.Exp(-(float64(dx*dx+dy*dy) / (2 * sigmaSq)))
			total += kernel[y][x]
		}
	}

	// 归一化核权重
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= total
		}
	}
	return kernel
}

// 计算高反差通道（原始通道 - 模糊通道）
func calcHighPassChannel(original, blurred [][]uint8, width, height int) [][]int {
	highPass := make([][]int, height)
	for y := 0; y < height; y++ {
		highPass[y] = make([]int, width)
		for x := 0; x < width; x++ {
			// 高反差 = 原始值 - 模糊值（保留纹理）
			highPass[y][x] = int(original[y][x]) - int(blurred[y][x])
		}
	}
	return highPass
}

// 生成纹理蒙版（高反差通道转换为0-255的灰度蒙版）
func generateTextureMask(highPass [][]int, width, height int) [][]uint8 {
	mask := make([][]uint8, height)
	// 先找到高反差通道的最大/最小值，用于归一化
	minVal, maxVal := 0, 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if highPass[y][x] < minVal {
				minVal = highPass[y][x]
			}
			if highPass[y][x] > maxVal {
				maxVal = highPass[y][x]
			}
		}
	}
	// 归一化到0-255（蒙版越亮表示纹理越强）
	rangeVal := maxVal - minVal
	if rangeVal == 0 {
		rangeVal = 1 // 避免除零
	}
	for y := 0; y < height; y++ {
		mask[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			normalized := (highPass[y][x] - minVal) * 255 / rangeVal
			// 反转蒙版并调整敏感度（降低纹理保留强度）
			mask[y][x] = uint8(255 - normalized)
			mask[y][x] = uint8(float64(mask[y][x]) * (1 - textureSensitivity))
		}
	}
	return mask
}

// 对整幅图像进行高斯模糊（用于混合的平滑图）
func gaussianBlurImage(src *image.RGBA, width, height int, radius float64) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	kernel := generateGaussianKernel(radius)
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2

	var wg sync.WaitGroup
	rowsPerSeg := (height + parallelSegments - 1) / parallelSegments

	for seg := 0; seg < parallelSegments; seg++ {
		startY := seg * rowsPerSeg
		endY := startY + rowsPerSeg
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(sStart, sEnd int) {
			defer wg.Done()
			for y := sStart; y < sEnd; y++ {
				for x := 0; x < width; x++ {
					var rSum, gSum, bSum, weightTotal float64
					// 应用高斯核
					for ky := -halfKernel; ky <= halfKernel; ky++ {
						for kx := -halfKernel; kx <= halfKernel; kx++ {
							nx, ny := x+kx, y+ky
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								weight := kernel[ky+halfKernel][kx+halfKernel]
								r, g, b, _ := src.At(nx, ny).RGBA()
								rSum += float64(r>>8) * weight
								gSum += float64(g>>8) * weight
								bSum += float64(b>>8) * weight
								weightTotal += weight
							}
						}
					}
					if weightTotal > 0 {
						r := uint8(rSum / weightTotal)
						g := uint8(gSum / weightTotal)
						b := uint8(bSum / weightTotal)
						_, _, _, a := src.At(x, y).RGBA()
						dst.SetRGBA(x, y, color.RGBA{r, g, b, uint8(a >> 8)})
					} else {
						dst.SetRGBA(x, y, src.At(x, y).(color.RGBA))
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return dst
}

// 基于纹理蒙版混合原图与平滑图
func blendByTextureMask(src, smoothed *image.RGBA, mask [][]uint8, width, height int) *image.RGBA {
	result := image.NewRGBA(src.Bounds())
	var wg sync.WaitGroup
	rowsPerSeg := (height + parallelSegments - 1) / parallelSegments

	for seg := 0; seg < parallelSegments; seg++ {
		startY := seg * rowsPerSeg
		endY := startY + rowsPerSeg
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(sStart, sEnd int) {
			defer wg.Done()
			for y := sStart; y < sEnd; y++ {
				for x := 0; x < width; x++ {
					// 蒙版值（0-255）→ 混合权重（0-1）：值越高，平滑图占比越高
					maskVal := float64(mask[y][x]) / 255.0
					// 原图与平滑图混合
					srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
					smooR, smooG, smooB, _ := smoothed.At(x, y).RGBA()

					finalR := uint8(float64(srcR>>8)*(1-maskVal) + float64(smooR>>8)*maskVal)
					finalG := uint8(float64(srcG>>8)*(1-maskVal) + float64(smooG>>8)*maskVal)
					finalB := uint8(float64(srcB>>8)*(1-maskVal) + float64(smooB>>8)*maskVal)

					result.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, uint8(srcA >> 8)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return result
}

// ========================================================================

// case17 人像磨皮 高反差磨皮算法

// 高反差磨皮算法（High Pass Sharpening for Skin Smoothing）的核心逻辑是通过 “高反差保留” 技术分离皮肤的 “底色” 与 “纹理 / 瑕疵”，仅对平坦的皮肤底色进行平滑，
//同时保留五官边缘和必要的皮肤质感，实现自然磨皮效果。该算法在专业修图中应用广泛，尤其适合需要保留皮肤细节（如毛孔）同时去除痘印、色斑等瑕疵的场景。

//算法原理
//高反差保留（High Pass Filter）：
//对图像进行高斯模糊，得到 “模糊图像”（保留底色，去除高频纹理）；用 “原始图像” 减去 “模糊图像”，得到 “高反差图像”——
//该图像仅保留亮度变化剧烈的区域（如痘印、皱纹、边缘），平坦区域（如脸颊）接近黑色。
//生成磨皮蒙版：
//将高反差图像转换为灰度蒙版，通过阈值处理强化 “需要磨皮的区域”（蒙版暗部，对应平坦皮肤）和 “需要保留的区域”（蒙版亮部，对应边缘 / 纹理）。
//混合磨皮：
//用 “模糊图像”（平滑底色）与 “原始图像” 基于蒙版混合 —— 蒙版暗部区域采用模糊图像（磨皮），蒙版亮部区域保留原始图像（保细节）。

func case17() {
	inputPath := "test4.jpg"          // 输入人像路径
	outputPath := "output_case17.jpg" // 高反差磨皮结果路径

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 执行高反差磨皮
	smoothedImg := HighPassSkinSmoothing(inputImg)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, smoothedImg, &jpeg.Options{Quality: 95})
	println("高反差磨皮处理完成！")
	println("结果保存至:", outputPath)
}

// 高反差磨皮核心参数（控制磨皮强度与细节保留）
const (
	highPassRadius  = 3.0 // 高反差保留半径（控制保留细节的大小，1-5）
	blurStrength    = 1.5 // 基础模糊强度（高反差的辅助参数）
	maskThreshold   = 30  // 蒙版阈值（0-255，值越高保留细节越多）
	skinSmoothRatio = 0.8 // 磨皮混合比例（0-1，越高磨皮越明显）
	parallelChunks  = 4   // 并行处理的图像块数
)

// 核心函数：高反差磨皮主流程
func HighPassSkinSmoothing(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	srcImg := image.NewRGBA(bounds)

	// 转换输入图像为RGBA格式（统一处理）
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			srcImg.SetRGBA(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8),
			})
		}
	}

	// 1. 对原图进行高斯模糊（获取皮肤底色）
	blurredImg := gaussianBlur(srcImg, width, height, highPassRadius*blurStrength)

	// 2. 计算高反差图像（原始图 - 模糊图，保留纹理和边缘）
	highPassImg := calcHighPassImage(srcImg, blurredImg, width, height)

	// 3. 生成磨皮蒙版（基于高反差图像，区分磨皮区域与保留区域）
	smoothMask := generateSmoothMask(highPassImg, width, height)

	// 4. 基于蒙版混合原图与模糊图，实现高反差磨皮
	resultImg := blendWithMask(srcImg, blurredImg, smoothMask, width, height)

	return resultImg
}

// 高斯模糊（修正ky/kx循环范围，避免索引越界）
func gaussianBlur(src *image.RGBA, width, height int, radius float64) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	kernel := createGaussianKernel(radius)
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2 // 此时kernelSize为奇数，halfKernel是正确的中心偏移

	var wg sync.WaitGroup
	rowsPerChunk := (height + parallelChunks - 1) / parallelChunks

	for chunk := 0; chunk < parallelChunks; chunk++ {
		startY := chunk * rowsPerChunk
		endY := startY + rowsPerChunk
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(cStart, cEnd int) {
			defer wg.Done()
			for y := cStart; y < cEnd; y++ {
				for x := 0; x < width; x++ {
					var rSum, gSum, bSum, weightTotal float64
					// 严格限制ky和kx在[-halfKernel, halfKernel]范围内（闭区间）
					for ky := -halfKernel; ky <= halfKernel; ky++ {
						for kx := -halfKernel; kx <= halfKernel; kx++ {
							nx, ny := x+kx, y+ky
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								// 计算核索引：ky+halfKernel范围是0~kernelSize-1（无越界）
								kernelY := ky + halfKernel
								kernelX := kx + halfKernel
								weight := kernel[kernelY][kernelX]
								r, g, b, _ := src.At(nx, ny).RGBA()
								rSum += float64(r>>8) * weight
								gSum += float64(g>>8) * weight
								bSum += float64(b>>8) * weight
								weightTotal += weight
							}
						}
					}
					if weightTotal > 0 {
						r := uint8(rSum / weightTotal)
						g := uint8(gSum / weightTotal)
						b := uint8(bSum / weightTotal)
						_, _, _, a := src.At(x, y).RGBA()
						dst.SetRGBA(x, y, color.RGBA{r, g, b, uint8(a >> 8)})
					} else {
						dst.SetRGBA(x, y, src.At(x, y).(color.RGBA))
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return dst
}

// 创建高斯核（确保尺寸为奇数，避免索引越界）
func createGaussianKernel(radius float64) [][]float64 {
	// 强制核尺寸为奇数：2*ceil(radius) + 1（确保至少3x3）
	radiusCeil := math.Ceil(radius)
	size := int(2*radiusCeil) + 1 // 保证是奇数（如radius=3.0→size=7，radius=4.5→size=11）
	if size < 3 {
		size = 3 // 最小核尺寸为3x3，避免过小导致计算异常
	}
	kernel := make([][]float64, size)
	halfKernel := size / 2 // 整数除法，如size=7→halfKernel=3（索引0~6）
	sigma := radius / 1.5  // sigma控制高斯分布的陡峭程度

	var total float64
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := x - halfKernel
			dy := y - halfKernel
			kernel[y][x] = math.Exp(-(float64(dx*dx+dy*dy) / (2 * sigma * sigma)))
			total += kernel[y][x]
		}
	}

	// 归一化核权重
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= total
		}
	}
	return kernel
}

// 计算高反差图像（原始图 - 模糊图）
func calcHighPassImage(src, blurred *image.RGBA, width, height int) *image.Gray {
	highPass := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 原始像素转灰度
			srcR, srcG, srcB, _ := src.At(x, y).RGBA()
			srcGray := uint8(0.299*float64(srcR>>8) + 0.587*float64(srcG>>8) + 0.114*float64(srcB>>8))
			// 模糊像素转灰度
			blurR, blurG, blurB, _ := blurred.At(x, y).RGBA()
			blurGray := uint8(0.299*float64(blurR>>8) + 0.587*float64(blurG>>8) + 0.114*float64(blurB>>8))
			// 高反差 = 原始灰度 - 模糊灰度（取绝对值，确保非负）
			diff := int(srcGray) - int(blurGray)
			if diff < 0 {
				diff = -diff
			}
			highPass.SetGray(x, y, color.Gray{uint8(diff)})
		}
	}
	return highPass
}

// 生成磨皮蒙版（高反差图像→阈值处理→反转）
func generateSmoothMask(highPass *image.Gray, width, height int) [][]float64 {
	mask := make([][]float64, height)
	for y := 0; y < height; y++ {
		mask[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 高反差图像的灰度值（0-255）
			hpVal := highPass.GrayAt(x, y).Y
			// 阈值处理：高于阈值的区域（边缘/纹理）保留原图，低于阈值的区域（平坦皮肤）磨皮
			if hpVal > maskThreshold {
				mask[y][x] = 0.0 // 0表示完全保留原图
			} else {
				// 低于阈值的区域：值越小，磨皮权重越高（越平坦磨皮越强）
				mask[y][x] = skinSmoothRatio * (1.0 - float64(hpVal)/float64(maskThreshold))
			}
		}
	}
	return mask
}

// 基于蒙版混合原图与模糊图（实现磨皮）
func blendWithMask(src, blurred *image.RGBA, mask [][]float64, width, height int) *image.RGBA {
	result := image.NewRGBA(src.Bounds())
	var wg sync.WaitGroup
	rowsPerChunk := (height + parallelChunks - 1) / parallelChunks

	for chunk := 0; chunk < parallelChunks; chunk++ {
		startY := chunk * rowsPerChunk
		endY := startY + rowsPerChunk
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(cStart, cEnd int) {
			defer wg.Done()
			for y := cStart; y < cEnd; y++ {
				for x := 0; x < width; x++ {
					// 磨皮权重（0-1）：值越高，模糊图占比越高
					weight := mask[y][x]
					// 混合像素值
					srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
					blurR, blurG, blurB, _ := blurred.At(x, y).RGBA()

					finalR := uint8(float64(srcR>>8)*(1-weight) + float64(blurR>>8)*weight)
					finalG := uint8(float64(srcG>>8)*(1-weight) + float64(blurG>>8)*weight)
					finalB := uint8(float64(srcB>>8)*(1-weight) + float64(blurB>>8)*weight)

					result.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, uint8(srcA >> 8)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return result
}

// ========================================================================

// case18 人像磨皮 细节叠加磨皮算法

// 细节叠加磨皮算法（Detail-Preserving Overlay Smoothing）的核心是在平滑皮肤瑕疵的基础上，
//将原始图像的高频细节（如细腻毛孔、皮肤纹理）叠加回处理结果，解决传统磨皮 “过度模糊导致塑料感” 的问题，实现 “平滑且有质感” 的自然效果

//算法原理
//基础平滑（去瑕疵）：对图像进行重度高斯模糊，得到 “平滑图”—— 该图去除了皮肤表面的痘印、色斑等瑕疵，但也丢失了正常纹理（如毛孔）。
//细节提取：对原图进行轻度高斯模糊，用 “原始图像” 减去 “轻度模糊图”，得到 “细节图”—— 该图仅保留皮肤的高频细节（如细腻纹理、毛孔），过滤掉低频的均匀色块。
//皮肤区域蒙版：通过肤色检测生成皮肤概率图（0-1），确保仅在皮肤区域进行磨皮和细节叠加，避免影响头发、背景等非皮肤区域。
//细节叠加：将 “平滑图” 与 “细节图” 按比例混合（平滑图占主导，细节图叠加补充质感），并结合皮肤蒙版控制范围，最终得到 “平滑且保留细节” 的结果。

func case18() {
	inputPath := "test4.jpg"          // 输入人像路径
	outputPath := "output_case18.jpg" // 细节叠加磨皮结果路径

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 执行细节叠加磨皮
	smoothedImg := DetailOverlaySkinSmoothing(inputImg)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, smoothedImg, &jpeg.Options{Quality: 95})
	println("细节叠加磨皮处理完成！")
	println("结果保存至:", outputPath)
}

// 皮肤检测的高斯模型参数（用于区分皮肤区域）
type SkinGaussianModel struct {
	meanY  float64       // 皮肤亮度均值
	meanCr float64       // 皮肤Cr分量均值
	meanCb float64       // 皮肤Cb分量均值
	covInv [3][3]float64 // 3D协方差逆矩阵（Y, Cr, Cb）
	covDet float64       // 协方差行列式
}

// 初始化皮肤检测模型（适配多肤色）
func initSkinGaussianModel() SkinGaussianModel {
	return SkinGaussianModel{
		meanY:  145.0, // 皮肤平均亮度
		meanCr: 153.0, // 皮肤Cr均值（偏红特征）
		meanCb: 102.0, // 皮肤Cb均值
		covInv: [3][3]float64{ // 预计算的3D协方差逆矩阵
			{0.0006, 0.0001, -0.0001},
			{0.0001, 0.007, -0.0009},
			{-0.0001, -0.0009, 0.008},
		},
		covDet: 320000.0, // 3D协方差行列式
	}
}

// 细节叠加磨皮核心参数（控制效果与质感）
const (
	heavyBlurRadius    = 4.0 // 重度模糊半径（去瑕疵，1-6，值越大磨皮越强）
	lightBlurRadius    = 1.0 // 轻度模糊半径（提细节，0.5-2，值越小保留细节越细）
	detailOverlayRatio = 0.4 // 细节叠加比例（0-1，值越高质感越强）
	skinMaskThreshold  = 0.3 // 皮肤蒙版阈值（0-1，低于此值不磨皮）
)

// 核心函数：细节叠加磨皮主流程
func DetailOverlaySkinSmoothing(input image.Image) *image.RGBA {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	srcImg := image.NewRGBA(bounds)

	// 转换输入图像为RGBA格式
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			srcImg.SetRGBA(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8),
			})
		}
	}

	// 1. 生成皮肤区域蒙版（0-1，1表示纯皮肤）
	skinMask := generateSkinMask(srcImg, width, height)

	// 2. 对原图进行重度模糊（去除瑕疵，得到平滑图）
	smoothedImg := gaussianBlur18(srcImg, width, height, heavyBlurRadius)

	// 3. 提取细节图（原始图 - 轻度模糊图，保留高频纹理）
	detailImg := extractDetailMap(srcImg, width, height, lightBlurRadius)

	// 4. 细节叠加：平滑图 + 细节图 * 比例，结合皮肤蒙版
	resultImg := overlayDetail(smoothedImg, srcImg, detailImg, skinMask, width, height)

	return resultImg
}

// 生成皮肤区域蒙版（基于3D高斯模型）
func generateSkinMask(img *image.RGBA, width, height int) [][]float64 {
	skinMask := make([][]float64, height)
	model := initSkinGaussianModel()

	for y := 0; y < height; y++ {
		skinMask[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			yVal, cr, cb := rgbToYCrCb18(r8, g8, b8)

			// 计算皮肤概率（0-1）
			prob := skinGaussianProbability18(
				float64(yVal), float64(cr), float64(cb),
				model,
			)
			// 低于阈值的区域设为0（非皮肤）
			if prob < skinMaskThreshold {
				skinMask[y][x] = 0.0
			} else {
				skinMask[y][x] = prob // 皮肤区域保留概率值
			}
		}
	}
	return skinMask
}

// 高斯模糊（支持不同半径，用于生成平滑图和轻度模糊图）
func gaussianBlur18(src *image.RGBA, width, height int, radius float64) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	kernel := createGaussianKernel18(radius)
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2 // 核中心偏移量

	var wg sync.WaitGroup
	rowsPerBlock := (height + parallelBlocks - 1) / parallelBlocks

	for block := 0; block < parallelBlocks; block++ {
		startY := block * rowsPerBlock
		endY := startY + rowsPerBlock
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(bStart, bEnd int) {
			defer wg.Done()
			for y := bStart; y < bEnd; y++ {
				for x := 0; x < width; x++ {
					var rSum, gSum, bSum, weightTotal float64
					// 遍历高斯核（确保索引不越界）
					for ky := -halfKernel; ky <= halfKernel; ky++ {
						for kx := -halfKernel; kx <= halfKernel; kx++ {
							nx, ny := x+kx, y+ky
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								// 核索引计算（确保在0~kernelSize-1范围内）
								kernelY := ky + halfKernel
								kernelX := kx + halfKernel
								weight := kernel[kernelY][kernelX]

								r, g, b, _ := src.At(nx, ny).RGBA()
								rSum += float64(r>>8) * weight
								gSum += float64(g>>8) * weight
								bSum += float64(b>>8) * weight
								weightTotal += weight
							}
						}
					}
					if weightTotal > 0 {
						r := uint8(rSum / weightTotal)
						g := uint8(gSum / weightTotal)
						b := uint8(bSum / weightTotal)
						_, _, _, a := src.At(x, y).RGBA()
						dst.SetRGBA(x, y, color.RGBA{r, g, b, uint8(a >> 8)})
					} else {
						dst.SetRGBA(x, y, src.At(x, y).(color.RGBA))
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return dst
}

// 创建高斯核（确保尺寸为奇数，避免索引越界）
func createGaussianKernel18(radius float64) [][]float64 {
	radiusCeil := math.Ceil(radius)
	size := int(2*radiusCeil) + 1 // 强制奇数尺寸（如radius=4→size=9）
	if size < 3 {
		size = 3 // 最小3x3核
	}
	kernel := make([][]float64, size)
	halfSize := size / 2
	sigma := radius / 1.2 // 控制高斯分布的扩散程度

	var total float64
	for y := 0; y < size; y++ {
		kernel[y] = make([]float64, size)
		for x := 0; x < size; x++ {
			dx := x - halfSize
			dy := y - halfSize
			kernel[y][x] = math.Exp(-(float64(dx*dx+dy*dy) / (2 * sigma * sigma)))
			total += kernel[y][x]
		}
	}

	// 归一化核权重
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= total
		}
	}
	return kernel
}

// 提取细节图（原始图 - 轻度模糊图，保留高频纹理）
func extractDetailMap(src *image.RGBA, width, height int, lightRadius float64) *image.RGBA {
	// 生成轻度模糊图（仅模糊微小瑕疵，保留大部分纹理）
	lightBlurImg := gaussianBlur(src, width, height, lightRadius)
	detailImg := image.NewRGBA(src.Bounds())

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 原始像素
			srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
			srcR8, srcG8, srcB8 := uint8(srcR>>8), uint8(srcG>>8), uint8(srcB>>8)
			// 轻度模糊像素
			blurR, blurG, blurB, _ := lightBlurImg.At(x, y).RGBA()
			blurR8, blurG8, blurB8 := uint8(blurR>>8), uint8(blurG>>8), uint8(blurB>>8)

			// 细节 = 原始 - 轻度模糊（取差值，保留高频信息）
			detailR := clampColor18(int(srcR8) - int(blurR8))
			detailG := clampColor18(int(srcG8) - int(blurG8))
			detailB := clampColor18(int(srcB8) - int(blurB8))

			detailImg.SetRGBA(x, y, color.RGBA{detailR, detailG, detailB, uint8(srcA >> 8)})
		}
	}
	return detailImg
}

// 细节叠加：平滑图 + 细节图 * 比例，结合皮肤蒙版
func overlayDetail(smoothed, src, detail *image.RGBA, skinMask [][]float64, width, height int) *image.RGBA {
	result := image.NewRGBA(src.Bounds())
	var wg sync.WaitGroup
	rowsPerBlock := (height + parallelBlocks - 1) / parallelBlocks

	for block := 0; block < parallelBlocks; block++ {
		startY := block * rowsPerBlock
		endY := startY + rowsPerBlock
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go func(bStart, bEnd int) {
			defer wg.Done()
			for y := bStart; y < bEnd; y++ {
				for x := 0; x < width; x++ {
					// 皮肤蒙版权重（0-1，非皮肤区域为0）
					maskWeight := skinMask[y][x]
					if maskWeight < 0.01 {
						// 非皮肤区域直接用原图
						result.SetRGBA(x, y, src.At(x, y).(color.RGBA))
						continue
					}

					// 平滑图像素
					smoothR, smoothG, smoothB, _ := smoothed.At(x, y).RGBA()
					smoothR8, smoothG8, smoothB8 := uint8(smoothR>>8), uint8(smoothG>>8), uint8(smoothB>>8)
					// 细节图像素
					dtlR, dtlG, dtlB, _ := detail.At(x, y).RGBA()
					dtlR8, dtlG8, dtlB8 := uint8(dtlR>>8), uint8(dtlG>>8), uint8(dtlB>>8)
					// 原图Alpha通道
					_, _, _, srcA := src.At(x, y).RGBA()
					srcA8 := uint8(srcA >> 8)

					// 细节叠加公式：平滑图 + 细节图 * 叠加比例 * 皮肤权重
					finalR := clampColor18(int(smoothR8) + int(float64(dtlR8)*detailOverlayRatio*maskWeight))
					finalG := clampColor18(int(smoothG8) + int(float64(dtlG8)*detailOverlayRatio*maskWeight))
					finalB := clampColor18(int(smoothB8) + int(float64(dtlB8)*detailOverlayRatio*maskWeight))

					result.SetRGBA(x, y, color.RGBA{finalR, finalG, finalB, srcA8})
				}
			}
		}(startY, endY)
	}
	wg.Wait()

	return result
}

// 计算皮肤的高斯概率（3D：Y, Cr, Cb）
func skinGaussianProbability18(y, cr, cb float64, model SkinGaussianModel) float64 {
	dY := y - model.meanY
	dCr := cr - model.meanCr
	dCb := cb - model.meanCb

	// 马氏距离平方（衡量与皮肤分布的偏离）
	mahalanobis := dY*(dY*model.covInv[0][0]+dCr*model.covInv[0][1]+dCb*model.covInv[0][2]) +
		dCr*(dY*model.covInv[1][0]+dCr*model.covInv[1][1]+dCb*model.covInv[1][2]) +
		dCb*(dY*model.covInv[2][0]+dCr*model.covInv[2][1]+dCb*model.covInv[2][2])

	// 高斯概率密度
	normalizer := 1.0 / (math.Pow(2*math.Pi, 1.5) * math.Sqrt(model.covDet))
	pdf := normalizer * math.Exp(-0.5*mahalanobis)
	peakPDF := 1.0 / (math.Pow(2*math.Pi, 1.5) * math.Sqrt(model.covDet))
	return math.Min(1.0, pdf/peakPDF)
}

// RGB转YCrCb（用于皮肤检测）
func rgbToYCrCb18(r, g, b uint8) (y, cr, cb uint8) {
	y = uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) + 0.5)
	cr = uint8(128 + 0.713*(float64(r)-float64(y)) + 0.5)
	cb = uint8(128 + 0.564*(float64(b)-float64(y)) + 0.5)
	return y, cr, cb
}

// 颜色值限制在[0,255]
func clampColor18(val int) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

// ========================================================================

// case19 图像放射变换

// 图像放射变换（Affine Transformation）是一种保持直线和平行性的几何变换，能实现平移、旋转、缩放、剪切等操作，广泛应用于图像对齐、
//校正、视角变换等场景。其核心是通过线性变换（矩阵乘法）+ 平移（向量加法）的组合，将原始图像的坐标映射到新位置。

func case19() {
	// 读取输入图像
	inputPath := "test4.jpg"
	outputPath := "output_case19.jpg"
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义放射变换参数（示例：缩放0.8倍 + 旋转30度 + 平移50像素 + 轻微剪切）
	params := AffineParams{
		ScaleX:     0.8, // X方向缩放0.8倍
		ScaleY:     0.8, // Y方向缩放0.8倍
		RotateDeg:  30,  // 顺时针旋转30度
		TranslateX: 50,  // X方向右移50像素
		TranslateY: 50,  // Y方向下移50像素
		SkewX:      0.2, // 轻微X方向剪切
	}

	// 应用放射变换
	resultImg := ApplyAffineTransform(srcImg, params)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, resultImg, &jpeg.Options{Quality: 95})
	println("放射变换完成，结果保存至:", outputPath)
}

// 放射变换参数（支持组合变换：先缩放→旋转→平移→剪切）
type AffineParams struct {
	ScaleX     float64 // X方向缩放因子（1.0为不变）
	ScaleY     float64 // Y方向缩放因子
	RotateDeg  float64 // 旋转角度（度）
	TranslateX float64 // X方向平移量（像素）
	TranslateY float64 // Y方向平移量（像素）
	SkewX      float64 // X方向剪切因子（0为不变）
}

// 构建放射变换矩阵（按：缩放→旋转→剪切→平移的顺序组合）
func buildAffineMatrix(params AffineParams) [3][3]float64 {
	// 1. 缩放矩阵
	scale := [3][3]float64{
		{params.ScaleX, 0, 0},
		{0, params.ScaleY, 0},
		{0, 0, 1},
	}

	// 2. 旋转矩阵（角度转弧度）
	theta := params.RotateDeg * math.Pi / 180.0
	cosθ, sinθ := math.Cos(theta), math.Sin(theta)
	rotate := [3][3]float64{
		{cosθ, -sinθ, 0},
		{sinθ, cosθ, 0},
		{0, 0, 1},
	}

	// 3. 剪切矩阵（仅X方向）
	skew := [3][3]float64{
		{1, params.SkewX, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	// 4. 平移矩阵
	translate := [3][3]float64{
		{1, 0, params.TranslateX},
		{0, 1, params.TranslateY},
		{0, 0, 1},
	}

	// 组合矩阵：先缩放→旋转→剪切→平移（矩阵乘法顺序：右乘）
	// 最终矩阵 = 平移 × 剪切 × 旋转 × 缩放
	mat := multiplyMatrices(translate, skew)
	mat = multiplyMatrices(mat, rotate)
	mat = multiplyMatrices(mat, scale)

	return mat
}

// 矩阵乘法（3x3矩阵）
func multiplyMatrices(a, b [3][3]float64) [3][3]float64 {
	var res [3][3]float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			res[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j]
		}
	}
	return res
}

// 计算变换矩阵的逆矩阵（用于逆映射）
func invertMatrix(mat [3][3]float64) (inv [3][3]float64, ok bool) {
	// 计算行列式
	det := mat[0][0]*(mat[1][1]*mat[2][2]-mat[1][2]*mat[2][1]) -
		mat[0][1]*(mat[1][0]*mat[2][2]-mat[1][2]*mat[2][0]) +
		mat[0][2]*(mat[1][0]*mat[2][1]-mat[1][1]*mat[2][0])

	if det == 0 {
		return inv, false // 矩阵不可逆
	}

	// 伴随矩阵 / 行列式
	inv[0][0] = (mat[1][1]*mat[2][2] - mat[1][2]*mat[2][1]) / det
	inv[0][1] = (mat[0][2]*mat[2][1] - mat[0][1]*mat[2][2]) / det
	inv[0][2] = (mat[0][1]*mat[1][2] - mat[0][2]*mat[1][1]) / det
	inv[1][0] = (mat[1][2]*mat[2][0] - mat[1][0]*mat[2][2]) / det
	inv[1][1] = (mat[0][0]*mat[2][2] - mat[0][2]*mat[2][0]) / det
	inv[1][2] = (mat[0][2]*mat[1][0] - mat[0][0]*mat[1][2]) / det
	inv[2][0] = (mat[1][0]*mat[2][1] - mat[1][1]*mat[2][0]) / det
	inv[2][1] = (mat[0][1]*mat[2][0] - mat[0][0]*mat[2][1]) / det
	inv[2][2] = (mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]) / det

	return inv, true
}

// 应用放射变换到图像
func ApplyAffineTransform(src image.Image, params AffineParams) image.Image {
	srcBounds := src.Bounds()
	srcW, srcH := srcBounds.Max.X, srcBounds.Max.Y

	// 1. 构建变换矩阵及逆矩阵
	affineMat := buildAffineMatrix(params)
	invMat, ok := invertMatrix(affineMat)
	if !ok {
		panic("变换矩阵不可逆，无法执行放射变换")
	}

	// 2. 计算输出图像的边界（避免裁剪）
	minX, minY, maxX, maxY := float64(srcW), float64(srcH), 0.0, 0.0
	corners := [][2]int{{0, 0}, {srcW, 0}, {0, srcH}, {srcW, srcH}} // 原图四个角点
	for _, corner := range corners {
		x, y := float64(corner[0]), float64(corner[1])
		// 计算变换后的坐标
		xPrime := affineMat[0][0]*x + affineMat[0][1]*y + affineMat[0][2]
		yPrime := affineMat[1][0]*x + affineMat[1][1]*y + affineMat[1][2]
		// 更新边界
		if xPrime < minX {
			minX = xPrime
		}
		if xPrime > maxX {
			maxX = xPrime
		}
		if yPrime < minY {
			minY = yPrime
		}
		if yPrime > maxY {
			maxY = yPrime
		}
	}
	outW := int(maxX - minX + 1)
	outH := int(maxY - minY + 1)
	outBounds := image.Rect(0, 0, outW, outH)
	dst := image.NewRGBA(outBounds)

	// 3. 对输出图像每个像素执行逆变换+双线性插值
	for yPrime := 0; yPrime < outH; yPrime++ {
		for xPrime := 0; xPrime < outW; xPrime++ {
			// 输出坐标（x', y'）转换为原始坐标系（减去偏移）
			x := float64(xPrime) + minX
			y := float64(yPrime) + minY

			// 逆变换：计算原始图像中的对应坐标（x, y）
			srcX := invMat[0][0]*x + invMat[0][1]*y + invMat[0][2]
			srcY := invMat[1][0]*x + invMat[1][1]*y + invMat[1][2]

			// 双线性插值获取像素值
			r, g, b, a := bilinearInterpolation(src, srcX, srcY, srcW, srcH)
			dst.SetRGBA(xPrime, yPrime, color.RGBA{r, g, b, a})
		}
	}

	return dst
}

// 双线性插值（处理浮点数坐标）
func bilinearInterpolation(src image.Image, x, y float64, w, h int) (r, g, b, a uint8) {
	// 取整数部分（左上角坐标）
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	// 取小数部分（插值权重）
	u := x - float64(x0)
	v := y - float64(y0)

	// 边界检查（超出范围用边缘像素）
	x1 := clamp19(x0+1, 0, w-1)
	y1 := clamp19(y0+1, 0, h-1)
	x0 = clamp19(x0, 0, w-1)
	y0 = clamp19(y0, 0, h-1)

	// 获取四个邻域像素
	p00 := getPixelRGBA(src, x0, y0)
	p10 := getPixelRGBA(src, x1, y0)
	p01 := getPixelRGBA(src, x0, y1)
	p11 := getPixelRGBA(src, x1, y1)

	// 双线性插值公式
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// 获取像素的RGBA值
func getPixelRGBA(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// 限制值在[min, max]范围内
func clamp19(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ========================================================================

// case20 图像透视变换

// 图像透视变换（Perspective Transformation）是一种能够模拟 3D 视角变化的几何变换，可将倾斜、变形的平面（如斜拍的文档、建筑）校正为正面平行视角。
//与放射变换（保持平行线）不同，透视变换允许平行线相交（如道路向远方汇聚），更贴近人眼的透视效果，广泛应用于文档校正、图像对齐、AR 视角模拟等场景。

func case20() {
	inputPath := "test4.jpg"          // 输入图像（如倾斜的文档）
	outputPath := "output_case20.jpg" // 透视校正结果

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}
	srcBounds := srcImg.Bounds()
	srcW, srcH := srcBounds.Max.X, srcBounds.Max.Y

	// 定义4对对应坐标（原始倾斜坐标→目标校正坐标）
	// 原始坐标：倾斜文档的四个角（顺时针）
	origCoords := [4]PerspCoord{
		{X: 120, Y: 90},                                // 左上角
		{X: float64(srcW - 60), Y: 130},                // 右上角
		{X: float64(srcW - 90), Y: float64(srcH - 70)}, // 右下角
		{X: 70, Y: float64(srcH - 50)},                 // 左下角
	}

	// 目标坐标：校正后文档的矩形角（宽700，高900）
	tgtCoords := [4]PerspCoord{
		{X: 0, Y: 0},     // 左上角
		{X: 700, Y: 0},   // 右上角
		{X: 700, Y: 900}, // 右下角
		{X: 0, Y: 900},   // 左下角
	}

	// 执行透视视角校正
	correctedImg := ExecutePerspViewAdjust(srcImg, origCoords, tgtCoords)

	// 保存校正结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, correctedImg, &jpeg.Options{Quality: 95})
	println("透视视角校正完成！")
	println("结果保存至:", outputPath)
}

// 坐标点结构（替代通用Point，强调透视变换中的坐标特性）
type PerspCoord struct {
	X, Y float64 // 二维坐标值
}

// 推导透视映射矩阵（通过4对对应坐标，区别于get/calc）
// origCoords: 原始图像中的4个坐标点，tgtCoords: 目标视角中的4个对应点
func derivePerspMapMatrix(origCoords [4]PerspCoord, tgtCoords [4]PerspCoord) (mapMat [3][3]float64, valid bool) {
	// 构建线性方程组（8个方程求解8个矩阵参数）
	eqMatrix := make([][]float64, 8)
	eqValues := make([]float64, 8)

	for i := 0; i < 4; i++ {
		x, y := origCoords[i].X, origCoords[i].Y
		xPrime, yPrime := tgtCoords[i].X, tgtCoords[i].Y

		// 第2i行方程：x' = (a*x + b*y + c)/(g*x + h*y + 1)
		eqMatrix[2*i] = []float64{x, y, 1, 0, 0, 0, -x * xPrime, -y * xPrime}
		eqValues[2*i] = xPrime

		// 第2i+1行方程：y' = (d*x + e*y + f)/(g*x + h*y + 1)
		eqMatrix[2*i+1] = []float64{0, 0, 0, x, y, 1, -x * yPrime, -y * yPrime}
		eqValues[2*i+1] = yPrime
	}

	// 高斯消元求解方程组（区别于gaussianElimination）
	params, solved := gaussSolveEquations(eqMatrix, eqValues)
	if !solved {
		return mapMat, false
	}

	// 组装透视映射矩阵（固定最后一个元素为1）
	mapMat[0][0] = params[0]
	mapMat[0][1] = params[1]
	mapMat[0][2] = params[2]
	mapMat[1][0] = params[3]
	mapMat[1][1] = params[4]
	mapMat[1][2] = params[5]
	mapMat[2][0] = params[6]
	mapMat[2][1] = params[7]
	mapMat[2][2] = 1.0

	return mapMat, true
}

// 高斯消元法求解线性方程组（命名更强调“求解方程”）
func gaussSolveEquations(eqMat [][]float64, eqVals []float64) (solution []float64, success bool) {
	eqCount := len(eqMat)
	for i := 0; i < eqCount; i++ {
		// 寻找主元（最大系数所在行）
		pivotRow := i
		for j := i; j < eqCount; j++ {
			if math.Abs(eqMat[j][i]) > math.Abs(eqMat[pivotRow][i]) {
				pivotRow = j
			}
		}
		// 主元接近0，矩阵奇异（无解）
		if math.Abs(eqMat[pivotRow][i]) < 1e-8 {
			return nil, false
		}
		// 交换当前行与主元行
		eqMat[i], eqMat[pivotRow] = eqMat[pivotRow], eqMat[i]
		eqVals[i], eqVals[pivotRow] = eqVals[pivotRow], eqVals[i]

		// 消元操作（化为上三角矩阵）
		for j := i + 1; j < eqCount; j++ {
			factor := eqMat[j][i] / eqMat[i][i]
			eqVals[j] -= factor * eqVals[i]
			for k := i; k < eqCount; k++ {
				eqMat[j][k] -= factor * eqMat[i][k]
			}
		}
	}

	// 回代求解参数
	solution = make([]float64, eqCount)
	for i := eqCount - 1; i >= 0; i-- {
		sum := eqVals[i]
		for j := i + 1; j < eqCount; j++ {
			sum -= eqMat[i][j] * solution[j]
		}
		solution[i] = sum / eqMat[i][i]
	}
	return solution, true
}

// 求透视矩阵的逆矩阵（用于反向坐标映射，区别于invert）
func reversePerspMatrix(mat [3][3]float64) (invMat [3][3]float64, valid bool) {
	// 计算矩阵行列式
	det := mat[0][0]*(mat[1][1]*mat[2][2]-mat[1][2]*mat[2][1]) -
		mat[0][1]*(mat[1][0]*mat[2][2]-mat[1][2]*mat[2][0]) +
		mat[0][2]*(mat[1][0]*mat[2][1]-mat[1][1]*mat[2][0])

	if math.Abs(det) < 1e-8 {
		return invMat, false // 行列式为0，矩阵不可逆
	}

	// 计算伴随矩阵并除以行列式（得到逆矩阵）
	invMat[0][0] = (mat[1][1]*mat[2][2] - mat[1][2]*mat[2][1]) / det
	invMat[0][1] = (mat[0][2]*mat[2][1] - mat[0][1]*mat[2][2]) / det
	invMat[0][2] = (mat[0][1]*mat[1][2] - mat[0][2]*mat[1][1]) / det
	invMat[1][0] = (mat[1][2]*mat[2][0] - mat[1][0]*mat[2][2]) / det
	invMat[1][1] = (mat[0][0]*mat[2][2] - mat[0][2]*mat[2][0]) / det
	invMat[1][2] = (mat[0][2]*mat[1][0] - mat[0][0]*mat[1][2]) / det
	invMat[2][0] = (mat[1][0]*mat[2][1] - mat[1][1]*mat[2][0]) / det
	invMat[2][1] = (mat[0][1]*mat[2][0] - mat[0][0]*mat[2][1]) / det
	invMat[2][2] = (mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]) / det

	return invMat, true
}

// 执行透视视角校正（核心函数，区别于ApplyPerspectiveTransform）
func ExecutePerspViewAdjust(srcImg image.Image, origCoords, tgtCoords [4]PerspCoord) image.Image {
	srcBounds := srcImg.Bounds()
	srcWidth, srcHeight := srcBounds.Max.X, srcBounds.Max.Y

	// 1. 推导透视映射矩阵
	perspMapMat, valid := derivePerspMapMatrix(origCoords, tgtCoords)
	if !valid {
		panic("透视映射矩阵推导失败，可能因坐标共线或输入错误")
	}

	// 2. 计算逆矩阵用于反向映射
	invPerspMat, valid := reversePerspMatrix(perspMapMat)
	if !valid {
		panic("透视矩阵不可逆，无法执行视角校正")
	}

	// 3. 确定输出图像边界（包含所有目标坐标）
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, coord := range tgtCoords {
		if coord.X < minX {
			minX = coord.X
		}
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.Y < minY {
			minY = coord.Y
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
	}
	outWidth := int(maxX - minX + 1)
	outHeight := int(maxY - minY + 1)
	outBounds := image.Rect(0, 0, outWidth, outHeight)
	dstImg := image.NewRGBA(outBounds)

	// 4. 对输出图像每个像素执行反向映射+双线性采样
	for yPrime := 0; yPrime < outHeight; yPrime++ {
		for xPrime := 0; xPrime < outWidth; xPrime++ {
			// 输出坐标映射到目标坐标系（抵消偏移）
			mapX := float64(xPrime) + minX
			mapY := float64(yPrime) + minY

			// 反向透视映射：计算原始图像中的对应坐标
			w := invPerspMat[2][0]*mapX + invPerspMat[2][1]*mapY + invPerspMat[2][2]
			if math.Abs(w) < 1e-8 {
				continue // 避免除零错误
			}
			origX := (invPerspMat[0][0]*mapX + invPerspMat[0][1]*mapY + invPerspMat[0][2]) / w
			origY := (invPerspMat[1][0]*mapX + invPerspMat[1][1]*mapY + invPerspMat[1][2]) / w

			// 双线性采样获取像素值（区别于interpolation）
			r, g, b, a := biLinearSample(srcImg, origX, origY, srcWidth, srcHeight)
			dstImg.SetRGBA(xPrime, yPrime, color.RGBA{r, g, b, a})
		}
	}

	return dstImg
}

// 双线性采样（处理浮点数坐标，区别于interpolation）
func biLinearSample(srcImg image.Image, x, y float64, imgWidth, imgHeight int) (r, g, b, a uint8) {
	x0 := int(math.Floor(x)) // 左上角x坐标
	y0 := int(math.Floor(y)) // 左上角y坐标
	u := x - float64(x0)     // x方向小数权重
	v := y - float64(y0)     // y方向小数权重

	// 边界检查（超出范围取边缘像素）
	x1 := limitRange(x0+1, 0, imgWidth-1)
	y1 := limitRange(y0+1, 0, imgHeight-1)
	x0 = limitRange(x0, 0, imgWidth-1)
	y0 = limitRange(y0, 0, imgHeight-1)

	// 获取四个邻域像素的RGBA值
	p00 := fetchPixelRGBA(srcImg, x0, y0)
	p10 := fetchPixelRGBA(srcImg, x1, y0)
	p01 := fetchPixelRGBA(srcImg, x0, y1)
	p11 := fetchPixelRGBA(srcImg, x1, y1)

	// 双线性加权计算最终像素值
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// 获取像素的RGBA值（区别于getPixelRGBA）
func fetchPixelRGBA(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// 限制数值在指定范围内（区别于clamp）
func limitRange(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// ========================================================================

// case21 图像反距离加权(IDW)插值变形算法

// IDW 插值变形通过控制点的位移带动周围像素平滑变形：
//每个像素的变形位移由所有控制点的位移加权求和得到，权重与像素到控制点的距离成反比（距离越近影响越大）。
//采用逆映射（从变形后坐标反推原始坐标）避免像素空洞，通过迭代逼近提高精度。
//双线性采样保证变形后图像的平滑性，避免锯齿。

func case21() {
	// 输入输出路径
	inputPath := "test4.jpg"          // 原始图像
	outputPath := "output_case21.jpg" // IDW变形结果

	// 读取原始图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义变形控制点（示例：将图像右上角区域向右侧拉伸）
	ctrlPoints := []DeformCtrlPoint{
		// 左上角点不变
		{OrigX: 50, OrigY: 50, TargetX: 50, TargetY: 50},
		// 右上角点向右移动100像素
		{OrigX: 300, OrigY: 50, TargetX: 400, TargetY: 50},
		// 右下角点向右移动50像素
		{OrigX: 300, OrigY: 200, TargetX: 350, TargetY: 200},
		// 左下角点不变
		{OrigX: 50, OrigY: 200, TargetX: 50, TargetY: 200},
	}

	// 设置IDW变形参数
	idwParams := IDWDeformParams{
		PowerFactor: 2.0,  // 权重与距离平方成反比（常用值）
		MinDist:     1e-3, // 最小距离阈值（避免除零）
	}

	// 执行IDW插值变形
	deformedImg := performIDWDeformation(srcImg, ctrlPoints, idwParams)

	// 保存变形结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, deformedImg, &jpeg.Options{Quality: 95})
	println("IDW插值图像变形完成！")
	println("结果保存至:", outputPath)
}

// 变形控制点结构（包含原始位置和目标变形位置）
type DeformCtrlPoint struct {
	OrigX, OrigY float64 // 原始图像中的控制点坐标
	TargetX      float64 // 变形后该点的目标X坐标
	TargetY      float64 // 变形后该点的目标Y坐标
}

// IDW变形核心参数（控制权重衰减和插值效果）
type IDWDeformParams struct {
	PowerFactor float64 // 距离权重的幂指数（通常取2，值越大权重衰减越快）
	MinDist     float64 // 最小距离阈值（避免除零，通常设为1e-3）
}

// 计算单个像素的IDW变形位移（核心函数：根据控制点计算偏移）
func calcPixelDisplacement(
	pixelX, pixelY float64,
	ctrlPoints []DeformCtrlPoint,
	params IDWDeformParams,
) (deltaX, deltaY float64) {
	var totalWeight float64    // 权重总和
	var weightedDeltaX float64 // X方向加权位移总和
	var weightedDeltaY float64 // Y方向加权位移总和

	for _, ctrl := range ctrlPoints {
		// 计算当前像素到控制点的欧氏距离
		dx := pixelX - ctrl.OrigX
		dy := pixelY - ctrl.OrigY
		dist := math.Hypot(dx, dy)

		// 距离过小则直接使用该控制点的位移（避免除零）
		if dist < params.MinDist {
			return ctrl.TargetX - ctrl.OrigX, ctrl.TargetY - ctrl.OrigY
		}

		// 反距离加权：权重 = 1 / (距离^幂指数)
		weight := 1.0 / math.Pow(dist, params.PowerFactor)
		totalWeight += weight

		// 累加该控制点对当前像素的加权位移
		deltaXCtrl := ctrl.TargetX - ctrl.OrigX // 控制点的X方向总位移
		deltaYCtrl := ctrl.TargetY - ctrl.OrigY // 控制点的Y方向总位移
		weightedDeltaX += weight * deltaXCtrl
		weightedDeltaY += weight * deltaYCtrl
	}

	// 归一化位移（除以总权重）
	if totalWeight > 1e-8 {
		deltaX = weightedDeltaX / totalWeight
		deltaY = weightedDeltaY / totalWeight
	}
	return deltaX, deltaY
}

// 执行IDW插值图像变形（主函数：遍历像素计算变形后图像）
func performIDWDeformation(
	srcImg image.Image,
	ctrlPoints []DeformCtrlPoint,
	params IDWDeformParams,
) image.Image {
	srcBounds := srcImg.Bounds()
	srcWidth, srcHeight := srcBounds.Max.X, srcBounds.Max.Y

	// 计算变形后图像的边界（避免裁剪控制点的目标位置）
	minDeformX, minDeformY := math.Inf(1), math.Inf(1)
	maxDeformX, maxDeformY := math.Inf(-1), math.Inf(-1)

	// 包含原始图像边界和所有控制点目标位置
	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			dx, dy := calcPixelDisplacement(float64(x), float64(y), ctrlPoints, params)
			deformX := float64(x) + dx
			deformY := float64(y) + dy
			minDeformX = math.Min(minDeformX, deformX)
			maxDeformX = math.Max(maxDeformX, deformX)
			minDeformY = math.Min(minDeformY, deformY)
			maxDeformY = math.Max(maxDeformY, deformY)
		}
	}
	for _, ctrl := range ctrlPoints {
		minDeformX = math.Min(minDeformX, ctrl.TargetX)
		maxDeformX = math.Max(maxDeformX, ctrl.TargetX)
		minDeformY = math.Min(minDeformY, ctrl.TargetY)
		maxDeformY = math.Max(maxDeformY, ctrl.TargetY)
	}

	// 确定输出图像尺寸
	outWidth := int(maxDeformX - minDeformX + 1)
	outHeight := int(maxDeformY - minDeformY + 1)
	outBounds := image.Rect(0, 0, outWidth, outHeight)
	dstImg := image.NewRGBA(outBounds)

	// 对输出图像每个像素执行逆映射（从变形后坐标反推原始坐标）
	for deformY := 0; deformY < outHeight; deformY++ {
		for deformX := 0; deformX < outWidth; deformX++ {
			// 将输出坐标映射到变形坐标系（抵消偏移）
			origDeformX := float64(deformX) + minDeformX
			origDeformY := float64(deformY) + minDeformY

			// 逆推：寻找原始图像中对应此变形坐标的像素（迭代逼近）
			// 原理：变形后坐标 = 原始坐标 + 位移 → 原始坐标 ≈ 变形后坐标 - 位移
			rawX, rawY := origDeformX, origDeformY // 初始猜测值
			// 迭代优化（通常3-5次即可收敛）
			for iter := 0; iter < 5; iter++ {
				dx, dy := calcPixelDisplacement(rawX, rawY, ctrlPoints, params)
				rawX = origDeformX - dx
				rawY = origDeformY - dy
			}

			// 双线性采样获取原始图像像素值
			r, g, b, a := sampleBilinear(
				srcImg, rawX, rawY,
				float64(srcWidth), float64(srcHeight),
			)
			dstImg.SetRGBA(deformX, deformY, color.RGBA{r, g, b, a})
		}
	}

	return dstImg
}

// 双线性采样（为IDW变形定制，区别于通用插值函数）
func sampleBilinear(
	img image.Image,
	x, y float64,
	imgWidth, imgHeight float64,
) (r, g, b, a uint8) {
	// 边界裁剪（超出图像范围的像素用边缘值）
	x = math.Max(0, math.Min(x, imgWidth-1))
	y = math.Max(0, math.Min(y, imgHeight-1))

	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	u := x - float64(x0) // X方向小数部分（权重）
	v := y - float64(y0) // Y方向小数部分（权重）

	x1 := x0 + 1
	y1 := y0 + 1
	// 边界检查（防止越界）
	if x1 >= int(imgWidth) {
		x1 = x0
	}
	if y1 >= int(imgHeight) {
		y1 = y0
	}

	// 获取四个邻域像素的RGBA值
	p00 := getDeformPixel(img, x0, y0)
	p10 := getDeformPixel(img, x1, y0)
	p01 := getDeformPixel(img, x0, y1)
	p11 := getDeformPixel(img, x1, y1)

	// 双线性加权计算最终像素值
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// 获取变形计算中的像素RGBA值（专用函数，区别于通用获取函数）
func getDeformPixel(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

// ========================================================================

// case22 图像特征线变换算法

//图像特征线变换算法（Feature Line Warping）是一种基于用户定义的特征线（如物体轮廓、边缘、结构线） 进行图像变形的技术，
//核心是通过约束特征线的变换（如移动、拉伸、弯曲），使特征线之间的区域按平滑过渡的方式跟随变形，同时保持图像的结构连贯性。
//该算法广泛应用于图像编辑（如人脸变形、物体形状调整）、动画生成等场景。

//算法原理
//特征线定义：用户指定多条特征线（每条线由起点和终点组成），包含原始特征线（OriginalSeg）和目标特征线（TargetSeg）—— 目标线是特征线变换后的期望位置。
//投影与距离计算：对图像中每个像素，计算其到每条特征线的垂直距离（影响权重）和投影点（在线段上的位置比例）。
//变形位移计算：根据像素在原始特征线上的投影点，映射到目标特征线的对应投影点，结合距离权重（距离特征线越近，受其影响越大），计算像素的最终位移。
//逆映射与采样：通过逆变换确定输出图像像素在原始图像中的对应位置，用双线性采样填充像素值，保证变形平滑。

func case22() {
	inputPath := "test4.jpg"          // 输入图像（如人脸）
	outputPath := "output_case22.jpg" // 特征线变换结果

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义特征线（示例：调整人脸嘴部轮廓）
	// 原始特征线：嘴部上边缘和下边缘
	featLines := []FeatureLine{
		{
			OrigSeg: LineSegment{ // 原始上唇线
				StartX: 150, StartY: 220,
				EndX: 250, EndY: 220,
			},
			TargetSeg: LineSegment{ // 目标上唇线（上移10像素）
				StartX: 150, StartY: 210,
				EndX: 250, EndY: 210,
			},
		},
		{
			OrigSeg: LineSegment{ // 原始下唇线
				StartX: 150, StartY: 240,
				EndX: 250, EndY: 240,
			},
			TargetSeg: LineSegment{ // 目标下唇线（下移10像素）
				StartX: 150, StartY: 250,
				EndX: 250, EndY: 250,
			},
		},
	}

	// 特征线变换参数
	warpParams := LineWarpParams{
		WeightPower: 2.0,  // 权重与距离平方成反比
		MinLineDist: 1e-3, // 最小距离阈值
		BlendRadius: 50.0, // 特征线间混合半径（控制过渡范围）
	}

	// 执行特征线变换
	warpedImg := performFeatureWarp(srcImg, featLines, warpParams)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, warpedImg, &jpeg.Options{Quality: 95})
	println("特征线变换完成！")
	println("结果保存至:", outputPath)
}

// 线段结构（特征线的基本组成单位）
type LineSegment struct {
	StartX, StartY float64 // 线段起点坐标
	EndX, EndY     float64 // 线段终点坐标
}

// 特征线结构（包含原始线段和目标变换线段）
type FeatureLine struct {
	OrigSeg   LineSegment // 原始图像中的特征线
	TargetSeg LineSegment // 变换后目标位置的特征线
}

// 特征线变换参数（控制变形平滑度和权重衰减）
type LineWarpParams struct {
	WeightPower float64 // 距离权重的幂指数（值越大，近线区域受影响越显著）
	MinLineDist float64 // 最小线距离阈值（避免除零，通常设为1e-3）
	BlendRadius float64 // 特征线间混合半径（控制多线影响的过渡范围）
}

// 计算像素到线段的垂直距离和投影参数
func calcLineProj(px, py float64, seg LineSegment) (dist float64, t float64) {
	// 线段向量
	segDX := seg.EndX - seg.StartX
	segDY := seg.EndY - seg.StartY
	// 像素到线段起点的向量
	pxDX := px - seg.StartX
	pxDY := py - seg.StartY

	// 线段长度的平方
	segLenSq := segDX*segDX + segDY*segDY
	if segLenSq < 1e-8 { // 线段长度为0（点）
		return math.Hypot(pxDX, pxDY), 0.0
	}

	// 投影参数t（0~1表示投影在 segment 上，<0在起点外，>1在终点外）
	t = (pxDX*segDX + pxDY*segDY) / segLenSq
	t = math.Max(0, math.Min(t, 1)) // 限制在 segment 上

	// 投影点坐标
	projX := seg.StartX + t*segDX
	projY := seg.StartY + t*segDY

	// 垂直距离
	dist = math.Hypot(px-projX, py-projY)
	return dist, t
}

// 计算单个像素的特征线变换位移
func getLineWarpDelta(
	px, py float64,
	featLines []FeatureLine,
	params LineWarpParams,
) (deltaX, deltaY float64) {
	var totalWeight float64
	var weightedDeltaX, weightedDeltaY float64

	for _, featLine := range featLines {
		// 计算像素到原始特征线的距离和投影参数t
		dist, t := calcLineProj(px, py, featLine.OrigSeg)

		// 距离过近时直接使用特征线的位移（避免权重异常）
		if dist < params.MinLineDist {
			// 原始特征线在t处的点
			origX := featLine.OrigSeg.StartX + t*(featLine.OrigSeg.EndX-featLine.OrigSeg.StartX)
			origY := featLine.OrigSeg.StartY + t*(featLine.OrigSeg.EndY-featLine.OrigSeg.StartY)
			// 目标特征线在t处的点
			targetX := featLine.TargetSeg.StartX + t*(featLine.TargetSeg.EndX-featLine.TargetSeg.StartX)
			targetY := featLine.TargetSeg.StartY + t*(featLine.TargetSeg.EndY-featLine.TargetSeg.StartY)
			return targetX - origX, targetY - origY
		}

		// 距离权重：1/(距离^幂指数)，并加入混合半径平滑
		weight := 1.0 / math.Pow(dist, params.WeightPower)
		// 特征线长度归一化（避免长线过度影响）
		segLen := math.Hypot(
			featLine.OrigSeg.EndX-featLine.OrigSeg.StartX,
			featLine.OrigSeg.EndY-featLine.OrigSeg.StartY,
		)
		weight /= (1 + segLen/params.BlendRadius)

		totalWeight += weight

		// 计算该特征线在投影点t处的位移
		origX := featLine.OrigSeg.StartX + t*(featLine.OrigSeg.EndX-featLine.OrigSeg.StartX)
		origY := featLine.OrigSeg.StartY + t*(featLine.OrigSeg.EndY-featLine.OrigSeg.StartY)
		targetX := featLine.TargetSeg.StartX + t*(featLine.TargetSeg.EndX-featLine.TargetSeg.StartX)
		targetY := featLine.TargetSeg.StartY + t*(featLine.TargetSeg.EndY-featLine.TargetSeg.StartY)
		lineDeltaX := targetX - origX
		lineDeltaY := targetY - origY

		// 累加加权位移
		weightedDeltaX += weight * lineDeltaX
		weightedDeltaY += weight * lineDeltaY
	}

	// 归一化总位移
	if totalWeight > 1e-8 {
		deltaX = weightedDeltaX / totalWeight
		deltaY = weightedDeltaY / totalWeight
	}
	return deltaX, deltaY
}

// 执行特征线变换（主函数）
func performFeatureWarp(
	srcImg image.Image,
	featLines []FeatureLine,
	params LineWarpParams,
) image.Image {
	srcBounds := srcImg.Bounds()
	srcW, srcH := float64(srcBounds.Max.X), float64(srcBounds.Max.Y)

	// 计算变换后图像的边界（包含所有特征线目标位置和变形区域）
	minWarpX, minWarpY := math.Inf(1), math.Inf(1)
	maxWarpX, maxWarpY := math.Inf(-1), math.Inf(-1)

	// 遍历原始图像边界点计算变形后位置
	for _, corner := range []struct{ x, y float64 }{ // 用变量corner接收结构体实例
		{0, 0}, {srcW, 0}, {0, srcH}, {srcW, srcH},
	} {
		// 通过结构体字段访问x和y
		dx, dy := getLineWarpDelta(corner.x, corner.y, featLines, params)
		warpX := corner.x + dx
		warpY := corner.y + dy
		minWarpX = math.Min(minWarpX, warpX)
		maxWarpX = math.Max(maxWarpX, warpX)
		minWarpY = math.Min(minWarpY, warpY)
		maxWarpY = math.Max(maxWarpY, warpY)
	}

	// 包含所有目标特征线的端点
	for _, fl := range featLines {
		// 分两次调用math.Min：先比较当前最小值与起点，再比较结果与终点
		minWarpX = math.Min(math.Min(minWarpX, fl.TargetSeg.StartX), fl.TargetSeg.EndX)
		// 同理处理最大值
		maxWarpX = math.Max(math.Max(maxWarpX, fl.TargetSeg.StartX), fl.TargetSeg.EndX)

		// Y方向同样处理
		minWarpY = math.Min(math.Min(minWarpY, fl.TargetSeg.StartY), fl.TargetSeg.EndY)
		maxWarpY = math.Max(math.Max(maxWarpY, fl.TargetSeg.StartY), fl.TargetSeg.EndY)
	}

	// 输出图像尺寸
	outW := int(maxWarpX - minWarpX + 1)
	outH := int(maxWarpY - minWarpY + 1)
	outBounds := image.Rect(0, 0, outW, outH)
	dstImg := image.NewRGBA(outBounds)

	// 逆映射：从变换后坐标反推原始坐标（迭代优化）
	for warpY := 0; warpY < outH; warpY++ {
		for warpX := 0; warpX < outW; warpX++ {
			// 变换后坐标映射到全局坐标系
			globalX := float64(warpX) + minWarpX
			globalY := float64(warpY) + minWarpY

			// 迭代逼近原始坐标（变形后坐标 = 原始坐标 + 位移 → 原始坐标 ≈ 变形后坐标 - 位移）
			rawX, rawY := globalX, globalY
			for iter := 0; iter < 5; iter++ { // 5次迭代足够收敛
				dx, dy := getLineWarpDelta(rawX, rawY, featLines, params)
				rawX = globalX - dx
				rawY = globalY - dy
			}

			// 双线性采样获取像素值
			r, g, b, a := sampleLineWarpPixel(srcImg, rawX, rawY, srcW, srcH)
			dstImg.SetRGBA(warpX, warpY, color.RGBA{r, g, b, a})
		}
	}

	return dstImg
}

// 特征线变换专用双线性采样
func sampleLineWarpPixel(
	img image.Image,
	x, y, imgW, imgH float64,
) (r, g, b, a uint8) {
	// 边界裁剪（超出图像范围用边缘像素）
	x = math.Max(0, math.Min(x, imgW-1))
	y = math.Max(0, math.Min(y, imgH-1))

	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	u := x - float64(x0)
	v := y - float64(y0)

	x1 := x0 + 1
	y1 := y0 + 1
	if x1 >= int(imgW) {
		x1 = x0
	}
	if y1 >= int(imgH) {
		y1 = y0
	}

	// 获取邻域像素
	p00 := getWarpPixel(img, x0, y0)
	p10 := getWarpPixel(img, x1, y0)
	p01 := getWarpPixel(img, x0, y1)
	p11 := getWarpPixel(img, x1, y1)

	// 双线性加权
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// 特征线变换专用像素获取
func getWarpPixel(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

// ========================================================================

// case23 图像MLS变形算法

// 图像 MLS（Moving Least Squares，移动最小二乘）变形算法是一种基于局部加权拟合的高精度图像变形技术。
//它通过对控制点的移动建立局部最小二乘模型，计算每个像素的变形位移，能在保持图像局部几何特征（如形状相似性）的同时实现平滑变形，
//广泛应用于人脸编辑、物体形状调整、动画生成等场景。

//算法原理
//控制点定义：用户指定一组控制点（MLSControlPoint），包含原始位置（OrigX, OrigY）和目标位置（TargetX, TargetY）。
//局部权重计算：对每个像素，根据与控制点的距离计算权重（距离越近，权重越大，影响越显著）。
//最小二乘拟合：基于权重构建局部坐标变换模型（支持刚性、相似、仿射等变换模式），通过最小二乘求解最优变换矩阵。
//位移计算：利用拟合的变换矩阵，将像素从原始坐标映射到目标坐标，得到变形位移。
//逆映射与采样：通过逆变换确定输出图像像素在原始图像中的对应位置，双线性采样填充像素值，保证变形平滑。

func case23() {
	inputPath := "test4.jpg"          // 输入图像
	outputPath := "output_case23.jpg" // MLS变形结果

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义MLS控制点（示例：将图像右侧区域向右拉伸）
	ctrlPoints := []MLSControlPoint{
		{OrigX: 100, OrigY: 100, TargetX: 100, TargetY: 100}, // 左上角不变
		{OrigX: 300, OrigY: 100, TargetX: 400, TargetY: 100}, // 右上角右移100
		{OrigX: 300, OrigY: 300, TargetX: 400, TargetY: 300}, // 右下角右移100
		{OrigX: 100, OrigY: 300, TargetX: 100, TargetY: 300}, // 左下角不变
	}

	// MLS变形参数（相似变换模式，保持形状）
	mlsParams := MLSWarpParams{
		Mode:        MLSSimilar, // 相似变换（推荐用于大多数场景）
		WeightPower: 2.0,        // 权重与距离平方成反比
		Epsilon:     1e-4,       // 最小距离阈值
	}

	// 执行MLS变形
	warpedImg := performMLSWarp(srcImg, ctrlPoints, mlsParams)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, warpedImg, &jpeg.Options{Quality: 95})
	println("MLS变形完成！")
	println("结果保存至:", outputPath)
}

// MLS控制点结构（原始位置与目标变形位置）
type MLSControlPoint struct {
	OrigX, OrigY     float64 // 原始图像中的控制点坐标
	TargetX, TargetY float64 // 变形后目标位置坐标
}

// MLS变形模式（控制局部变换的刚性程度）
type MLSDeformMode int

const (
	MLSRigid   MLSDeformMode = iota // 刚性变换（保持旋转和缩放）
	MLSSimilar                      // 相似变换（保持形状相似性）
	MLSAffine                       // 仿射变换（允许剪切，灵活性最高）
)

// MLS变形参数
type MLSWarpParams struct {
	Mode        MLSDeformMode // 变形模式
	WeightPower float64       // 距离权重幂指数（值越大，局部影响越集中）
	Epsilon     float64       // 最小距离阈值（避免除零，通常1e-4）
}

// 计算单个像素的MLS变形位移
func calcMLSDelta(
	px, py float64,
	ctrlPoints []MLSControlPoint,
	params MLSWarpParams,
) (deltaX, deltaY float64) {
	// 1. 计算权重总和与加权中心（原始坐标）
	var sumWeights float64
	var sumOrigX, sumOrigY float64     // 原始控制点的加权中心
	var sumTargetX, sumTargetY float64 // 目标控制点的加权中心

	for _, ctrl := range ctrlPoints {
		// 像素到控制点的距离
		dx := px - ctrl.OrigX
		dy := py - ctrl.OrigY
		dist := math.Hypot(dx, dy)
		if dist < params.Epsilon {
			// 距离过近，直接返回控制点的位移
			return ctrl.TargetX - ctrl.OrigX, ctrl.TargetY - ctrl.OrigY
		}

		// 权重 = 1 / (距离^幂指数)
		weight := 1.0 / math.Pow(dist, params.WeightPower)
		sumWeights += weight

		// 累加加权坐标
		sumOrigX += weight * ctrl.OrigX
		sumOrigY += weight * ctrl.OrigY
		sumTargetX += weight * ctrl.TargetX
		sumTargetY += weight * ctrl.TargetY
	}

	if sumWeights < 1e-8 {
		return 0, 0 // 无有效控制点影响
	}

	// 加权中心坐标（原始与目标）
	avgOrigX := sumOrigX / sumWeights
	avgOrigY := sumOrigY / sumWeights
	avgTargetX := sumTargetX / sumWeights
	avgTargetY := sumTargetY / sumWeights

	// 2. 计算矩阵M和V（用于最小二乘拟合）
	var m00, m01, m11 float64 // 矩阵M的元素（对称矩阵）
	var v0, v1 float64        // 向量V的元素

	for _, ctrl := range ctrlPoints {
		dx := px - ctrl.OrigX
		dy := py - ctrl.OrigY
		dist := math.Hypot(dx, dy)
		if dist < params.Epsilon {
			continue // 已处理过近距离点
		}
		weight := 1.0 / math.Pow(dist, params.WeightPower)

		// 相对于原始加权中心的偏移
		qx := ctrl.OrigX - avgOrigX
		qy := ctrl.OrigY - avgOrigY
		// 相对于目标加权中心的偏移
		pqX := ctrl.TargetX - avgTargetX
		pqY := ctrl.TargetY - avgTargetY

		// 累加矩阵M和向量V的加权值
		m00 += weight * (qx*qx + qy*qy)
		m01 += weight * qx * qy
		m11 += weight * (qx*qx + qy*qy) // 刚性/相似模式下M的特殊结构
		v0 += weight * (pqX*qx + pqY*qy)
		v1 += weight * (-pqX*qy + pqY*qx)

		// 仿射模式下修正矩阵M（非对称）
		if params.Mode == MLSAffine {
			m01 += weight * qx * qy
			m11 += weight * qy * qy
			v1 += weight * (pqX*qx + pqY*qy) // 仿射模式下V1的计算不同
		}
	}

	// 3. 求解变换矩阵参数（根据变形模式）
	var a, b, c, d float64 // 变换矩阵 [[a, b], [c, d]]
	detM := m00*m11 - m01*m01
	if math.Abs(detM) < 1e-8 {
		return 0, 0 // 矩阵奇异，无法求解
	}

	switch params.Mode {
	case MLSRigid:
		// 刚性变换：a=d, b=-c（保持旋转和缩放）
		a = (m00*v0 + m01*v1) / detM
		b = (m01*v0 + m11*v1) / detM
		c = -b
		d = a
	case MLSSimilar:
		// 相似变换：a*d - b*c > 0（保持形状相似）
		a = (m00*v0 + m01*v1) / detM
		b = (m01*v0 + m11*v1) / detM
		c = (-m01*v0 - m11*v1) / detM
		d = (m00*v0 + m01*v1) / detM
	case MLSAffine:
		// 仿射变换：无约束（最灵活）
		a = (m11*v0 - m01*v1) / detM
		b = (m00*v1 - m01*v0) / detM
		c = (m11*v1 - m01*v0) / detM // 注意仿射模式下v1定义不同
		d = (m00*v0 - m01*v1) / detM
	}

	// 4. 计算像素的目标坐标
	pxRelOrigX := px - avgOrigX
	pxRelOrigY := py - avgOrigY
	targetX := avgTargetX + a*pxRelOrigX + b*pxRelOrigY
	targetY := avgTargetY + c*pxRelOrigX + d*pxRelOrigY

	return targetX - px, targetY - py // 位移 = 目标坐标 - 原始坐标
}

// 执行MLS变形（主函数）
func performMLSWarp(
	srcImg image.Image,
	ctrlPoints []MLSControlPoint,
	params MLSWarpParams,
) image.Image {
	srcBounds := srcImg.Bounds()
	srcW, srcH := float64(srcBounds.Max.X), float64(srcBounds.Max.Y)

	// 计算变形后图像的边界（包含所有变形区域）
	minWarpX, minWarpY := math.Inf(1), math.Inf(1)
	maxWarpX, maxWarpY := math.Inf(-1), math.Inf(-1)

	// 1. 包含原始图像四个角点的变形位置
	corners := []struct{ x, y float64 }{
		{0, 0}, {srcW, 0}, {0, srcH}, {srcW, srcH},
	}
	for _, corner := range corners {
		dx, dy := calcMLSDelta(corner.x, corner.y, ctrlPoints, params)
		warpX := corner.x + dx
		warpY := corner.y + dy
		minWarpX = math.Min(minWarpX, warpX)
		maxWarpX = math.Max(maxWarpX, warpX)
		minWarpY = math.Min(minWarpY, warpY)
		maxWarpY = math.Max(maxWarpY, warpY)
	}

	// 2. 包含所有控制点的目标位置
	for _, ctrl := range ctrlPoints {
		minWarpX = math.Min(minWarpX, ctrl.TargetX)
		maxWarpX = math.Max(maxWarpX, ctrl.TargetX)
		minWarpY = math.Min(minWarpY, ctrl.TargetY)
		maxWarpY = math.Max(maxWarpY, ctrl.TargetY)
	}

	// 输出图像尺寸
	outW := int(maxWarpX - minWarpX + 1)
	outH := int(maxWarpY - minWarpY + 1)
	outBounds := image.Rect(0, 0, outW, outH)
	dstImg := image.NewRGBA(outBounds)

	// 逆映射：从变形后坐标反推原始坐标（迭代优化）
	for warpY := 0; warpY < outH; warpY++ {
		for warpX := 0; warpX < outW; warpX++ {
			// 变形后坐标映射到全局坐标系
			globalX := float64(warpX) + minWarpX
			globalY := float64(warpY) + minWarpY

			// 迭代逼近原始坐标（变形后坐标 = 原始坐标 + 位移 → 原始坐标 ≈ 变形后坐标 - 位移）
			rawX, rawY := globalX, globalY
			for iter := 0; iter < 5; iter++ { // 5次迭代足够收敛
				dx, dy := calcMLSDelta(rawX, rawY, ctrlPoints, params)
				rawX = globalX - dx
				rawY = globalY - dy
			}

			// 双线性采样获取像素值
			r, g, b, a := sampleMLSPixel(srcImg, rawX, rawY, srcW, srcH)
			dstImg.SetRGBA(warpX, warpY, color.RGBA{r, g, b, a})
		}
	}

	return dstImg
}

// MLS变形专用双线性采样
func sampleMLSPixel(
	img image.Image,
	x, y, imgW, imgH float64,
) (r, g, b, a uint8) {
	// 边界裁剪（超出范围用边缘像素）
	x = math.Max(0, math.Min(x, imgW-1))
	y = math.Max(0, math.Min(y, imgH-1))

	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	u := x - float64(x0)
	v := y - float64(y0)

	x1 := x0 + 1
	y1 := y0 + 1
	if x1 >= int(imgW) {
		x1 = x0
	}
	if y1 >= int(imgH) {
		y1 = y0
	}

	// 获取邻域像素
	p00 := getMLSPixel(img, x0, y0)
	p10 := getMLSPixel(img, x1, y0)
	p01 := getMLSPixel(img, x0, y1)
	p11 := getMLSPixel(img, x1, y1)

	// 双线性加权计算
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// MLS变形专用像素获取
func getMLSPixel(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

// ========================================================================

// case24 图像MRLS算法  - 效果没成功，示例失败????

// 图像 MRLS（Moving Robust Least Squares，移动鲁棒最小二乘）算法是 MLS（移动最小二乘）算法的抗干扰改进版本。
//它通过引入鲁棒损失函数（如 Huber、Tukey 损失）降低异常值（如噪声控制点、误标定点）对变形结果的影响，在保持局部几何特征的同时，
//显著提升变形的稳定性和抗噪性，适用于控制点存在误差或噪声的场景（如手动标记的控制点、动态捕捉的噪声数据）。

//算法原理
//MRLS 在 MLS 基础上的核心改进是对异常控制点进行加权抑制，具体流程：
//控制点定义：同 MLS，包含原始位置（MRLSOrigX, MRLSOrigY）和目标位置（MRLSTargetX, MRLSTargetY），但允许存在少量异常点。
//双重权重计算：
//距离权重：与 MLS 一致，基于像素到控制点的距离（近点权重高）。
//鲁棒权重：通过损失函数计算，降低异常控制点的权重（偏离拟合模型越远的点，权重越低）。
//鲁棒最小二乘拟合：结合双重权重构建局部变换模型（刚性、相似、仿射），求解时抑制异常点的干扰。
//位移计算与逆映射：同 MLS，通过拟合模型计算像素位移，逆映射采样保证图像平滑。

func case24() {
	inputPath := "test4.jpg"          // 含噪声控制点的输入图像
	outputPath := "output_case24.jpg" // MRLS变形结果

	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义含噪声的MRLS控制点（第3个为异常点）
	ctrlPoints := []MRLSControlPoint{
		{OrigX: 100, OrigY: 100, TargetX: 100, TargetY: 100}, // 正常点
		{OrigX: 300, OrigY: 100, TargetX: 400, TargetY: 100}, // 正常点
		{OrigX: 300, OrigY: 300, TargetX: 600, TargetY: 300}, // 异常点（过度偏移）
		{OrigX: 100, OrigY: 300, TargetX: 100, TargetY: 300}, // 正常点
	}

	// MRLS参数（Tukey损失抑制异常点）
	mrlsParams := MRLSWarpParams{
		Mode:        MRLSSimilar, // 相似变换
		LossType:    MRLSTukey,   // 采用Tukey损失（抗噪性强）
		WeightPower: 2.0,         // 距离权重幂指数
		Epsilon:     1e-4,        // 最小距离阈值
		LossScale:   50.0,        // 损失阈值（大于此值视为异常点）
	}

	// 执行MRLS变形
	warpedImg := performMRLSWarp(srcImg, ctrlPoints, mrlsParams)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, warpedImg, &jpeg.Options{Quality: 95})
	println("MRLS变形完成！")
	println("结果保存至:", outputPath)
}

// MRLS控制点结构（含原始与目标位置）
type MRLSControlPoint struct {
	OrigX, OrigY     float64 // 原始坐标
	TargetX, TargetY float64 // 目标坐标
}

// MRLS鲁棒损失函数类型
type MRLSLossType int

const (
	MRLSHuber MRLSLossType = iota // Huber损失（对小误差线性，大误差饱和）
	MRLSTukey                     // Tukey损失（大误差直接归零，抗噪性最强）
	MRLSL1                        // L1损失（对异常值敏感度低于L2）
)

// MRLS变形模式（控制变换刚性）
type MRLSDeformMode int

const (
	MRLSRigid   MRLSDeformMode = iota // 刚性变换
	MRLSSimilar                       // 相似变换
	MRLSAffine                        // 仿射变换
)

// MRLS算法参数
type MRLSWarpParams struct {
	Mode        MRLSDeformMode // 变形模式
	LossType    MRLSLossType   // 鲁棒损失类型
	WeightPower float64        // 距离权重幂指数（>0）
	Epsilon     float64        // 最小距离阈值（避免除零，1e-4）
	LossScale   float64        // 损失函数缩放因子（控制异常值判断阈值）
}

// 计算鲁棒权重（基于损失函数）
func calcRobustWeight(residual float64, lossType MRLSLossType, scale float64) float64 {
	if scale < 1e-8 {
		return 1.0 // 缩放因子无效时退化到普通权重
	}
	r := residual / scale // 归一化残差

	switch lossType {
	case MRLSHuber:
		// Huber损失：|r|≤1时权重1，|r|>1时权重1/|r|
		if math.Abs(r) <= 1.0 {
			return 1.0
		}
		return 1.0 / math.Abs(r)
	case MRLSTukey:
		// Tukey损失：|r|≤1时权重(1-r²)²，|r|>1时权重0（完全忽略）
		if math.Abs(r) > 1.0 {
			return 0.0
		}
		return math.Pow(1-r*r, 2)
	case MRLSL1:
		// L1损失：权重1/|r|（对异常值敏感度低于L2）
		if math.Abs(r) < 1e-8 {
			return 1.0
		}
		return 1.0 / math.Abs(r)
	default:
		return 1.0
	}
}

// 计算单个像素的MRLS变形位移
func calcMRLSDelta(
	px, py float64,
	ctrlPoints []MRLSControlPoint,
	params MRLSWarpParams,
) (deltaX, deltaY float64) {
	// 若无控制点，返回0位移
	if len(ctrlPoints) == 0 {
		return 0, 0
	}

	// 初始化：首次拟合（无鲁棒权重，类似MLS）
	var initSumW float64
	var initSumOrigX, initSumOrigY float64
	var initSumTargetX, initSumTargetY float64

	for _, ctrl := range ctrlPoints {
		dx := px - ctrl.OrigX
		dy := py - ctrl.OrigY
		dist := math.Hypot(dx, dy)
		if dist < params.Epsilon {
			return ctrl.TargetX - ctrl.OrigX, ctrl.TargetY - ctrl.OrigY
		}
		wDist := 1.0 / math.Pow(dist, params.WeightPower) // 距离权重
		initSumW += wDist
		initSumOrigX += wDist * ctrl.OrigX
		initSumOrigY += wDist * ctrl.OrigY
		initSumTargetX += wDist * ctrl.TargetX
		initSumTargetY += wDist * ctrl.TargetY
	}

	if initSumW < 1e-8 {
		return 0, 0
	}

	// 首次拟合的加权中心
	avgOInitX := initSumOrigX / initSumW
	avgOInitY := initSumOrigY / initSumW
	avgTInitX := initSumTargetX / initSumW
	avgTInitY := initSumTargetY / initSumW

	// 计算首次拟合的残差（用于鲁棒权重）
	robustWeights := make([]float64, len(ctrlPoints))
	for i, ctrl := range ctrlPoints {
		qx := ctrl.OrigX - avgOInitX
		qy := ctrl.OrigY - avgOInitY
		pqX := ctrl.TargetX - avgTInitX
		pqY := ctrl.TargetY - avgTInitY

		// 残差：目标偏移与初始拟合偏移的差异（简化计算）
		residual := math.Hypot(pqX-qx, pqY-qy) // 假设初始为刚性变换
		robustWeights[i] = calcRobustWeight(residual, params.LossType, params.LossScale)
	}

	// 二次拟合：结合距离权重和鲁棒权重
	var sumW float64
	var sumOrigX, sumOrigY float64
	var sumTargetX, sumTargetY float64
	var m00, m01, m11 float64 // 拟合矩阵M
	var v0, v1 float64        // 拟合向量V

	for i, ctrl := range ctrlPoints {
		dx := px - ctrl.OrigX
		dy := py - ctrl.OrigY
		dist := math.Hypot(dx, dy)
		if dist < params.Epsilon {
			continue
		}
		wDist := 1.0 / math.Pow(dist, params.WeightPower)
		wTotal := wDist * robustWeights[i] // 总权重 = 距离权重 × 鲁棒权重
		sumW += wTotal

		sumOrigX += wTotal * ctrl.OrigX
		sumOrigY += wTotal * ctrl.OrigY
		sumTargetX += wTotal * ctrl.TargetX
		sumTargetY += wTotal * ctrl.TargetY

		// 相对于加权中心的偏移（实时计算当前平均，避免累计误差）
		currentAvgO_X := sumOrigX / sumW
		currentAvgO_Y := sumOrigY / sumW
		currentAvgT_X := sumTargetX / sumW
		currentAvgT_Y := sumTargetY / sumW

		qx := ctrl.OrigX - currentAvgO_X
		qy := ctrl.OrigY - currentAvgO_Y
		pqX := ctrl.TargetX - currentAvgT_X
		pqY := ctrl.TargetY - currentAvgT_Y

		// 累加矩阵M和向量V（根据变形模式）
		m00 += wTotal * (qx*qx + qy*qy)
		m01 += wTotal * qx * qy
		if params.Mode == MRLSAffine {
			m11 += wTotal * qy * qy // 仿射模式下M非对称
			v1 += wTotal * (pqX*qx + pqY*qy)
		} else {
			m11 += wTotal * (qx*qx + qy*qy) // 刚性/相似模式下M对称
			v1 += wTotal * (-pqX*qy + pqY*qx)
		}
		v0 += wTotal * (pqX*qx + pqY*qy)
	}

	if sumW < 1e-8 {
		return 0, 0
	}

	// 最终加权中心
	avgOrigX := sumOrigX / sumW
	avgOrigY := sumOrigY / sumW
	avgTargetX := sumTargetX / sumW
	avgTargetY := sumTargetY / sumW

	// 求解变换矩阵参数
	detM := m00*m11 - m01*m01
	if math.Abs(detM) < 1e-8 {
		return 0, 0
	}

	var a, b, c, d float64
	switch params.Mode {
	case MRLSRigid:
		a = (m00*v0 + m01*v1) / detM
		b = (m01*v0 + m11*v1) / detM
		c = -b
		d = a
	case MRLSSimilar:
		a = (m00*v0 + m01*v1) / detM
		b = (m01*v0 + m11*v1) / detM
		c = (-m01*v0 - m11*v1) / detM
		d = (m00*v0 + m01*v1) / detM
	case MRLSAffine:
		a = (m11*v0 - m01*v1) / detM
		b = (m00*v1 - m01*v0) / detM
		c = (m11*v1 - m01*v0) / detM
		d = (m00*v0 - m01*v1) / detM
	}

	// 计算目标坐标与位移
	pxRelX := px - avgOrigX
	pxRelY := py - avgOrigY
	targetX := avgTargetX + a*pxRelX + b*pxRelY
	targetY := avgTargetY + c*pxRelX + d*pxRelY

	return targetX - px, targetY - py
}

// 执行MRLS变形（主函数，修复边界计算错误）
func performMRLSWarp(
	srcImg image.Image,
	ctrlPoints []MRLSControlPoint,
	params MRLSWarpParams,
) image.Image {
	// 若无控制点，直接返回原图（避免后续计算错误）
	if len(ctrlPoints) == 0 {
		return srcImg
	}

	srcBounds := srcImg.Bounds()
	srcW, srcH := float64(srcBounds.Max.X), float64(srcBounds.Max.Y)

	// 初始化边界（使用图像实际范围作为初始值，避免无穷大导致的计算异常）
	minWarpX, minWarpY := srcW, srcH
	maxWarpX, maxWarpY := 0.0, 0.0

	// 1. 包含原始图像四角的变形位置
	corners := []struct{ x, y float64 }{
		{0, 0}, {srcW, 0}, {0, srcH}, {srcW, srcH},
	}
	for _, corner := range corners {
		dx, dy := calcMRLSDelta(corner.x, corner.y, ctrlPoints, params)
		warpX := corner.x + dx
		warpY := corner.y + dy
		// 更新边界时增加浮点数精度容错
		minWarpX = math.Min(minWarpX, math.Nextafter(warpX, math.Inf(-1)))
		maxWarpX = math.Max(maxWarpX, math.Nextafter(warpX, math.Inf(1)))
		minWarpY = math.Min(minWarpY, math.Nextafter(warpY, math.Inf(-1)))
		maxWarpY = math.Max(maxWarpY, math.Nextafter(warpY, math.Inf(1)))
	}

	// 2. 包含所有控制点目标位置
	for _, ctrl := range ctrlPoints {
		minWarpX = math.Min(minWarpX, ctrl.TargetX)
		maxWarpX = math.Max(maxWarpX, ctrl.TargetX)
		minWarpY = math.Min(minWarpY, ctrl.TargetY)
		maxWarpY = math.Max(maxWarpY, ctrl.TargetY)
	}

	// 修复核心：计算输出图像尺寸（处理负数、精度和超大值问题）
	deltaX := maxWarpX - minWarpX
	deltaY := maxWarpY - minWarpY

	// 处理负数尺寸（因计算误差或异常控制点导致）
	if deltaX < 0 {
		deltaX = 0
	}
	if deltaY < 0 {
		deltaY = 0
	}

	// 处理浮点数精度问题（避免极小值导致0尺寸）
	minValidSize := 1.0 // 最小有效尺寸（至少1像素）
	if deltaX < minValidSize {
		deltaX = minValidSize
	}
	if deltaY < minValidSize {
		deltaY = minValidSize
	}

	// 限制最大尺寸（避免超出系统内存限制，可根据需求调整）
	maxValidSize := 10000.0 // 最大10000像素
	if deltaX > maxValidSize {
		deltaX = maxValidSize
	}
	if deltaY > maxValidSize {
		deltaY = maxValidSize
	}

	// 转换为整数尺寸（+1确保包含边界像素）
	outW := int(deltaX + 1)
	outH := int(deltaY + 1)

	// 最终兜底检查（确保尺寸为正）
	if outW <= 0 || outH <= 0 {
		outW = 1
		outH = 1
	}

	// 创建输出图像（此时尺寸一定合法）
	outBounds := image.Rect(0, 0, outW, outH)
	dstImg := image.NewRGBA(outBounds)

	// 逆映射：迭代反推原始坐标
	for warpY := 0; warpY < outH; warpY++ {
		for warpX := 0; warpX < outW; warpX++ {
			globalX := float64(warpX) + minWarpX
			globalY := float64(warpY) + minWarpY

			// 迭代优化原始坐标
			rawX, rawY := globalX, globalY
			for iter := 0; iter < 5; iter++ {
				dx, dy := calcMRLSDelta(rawX, rawY, ctrlPoints, params)
				rawX = globalX - dx
				rawY = globalY - dy
			}

			// 双线性采样
			r, g, b, a := sampleMRLSPixel(srcImg, rawX, rawY, srcW, srcH)
			dstImg.SetRGBA(warpX, warpY, color.RGBA{r, g, b, a})
		}
	}

	return dstImg
}

// MRLS专用双线性采样
func sampleMRLSPixel(
	img image.Image,
	x, y, imgW, imgH float64,
) (r, g, b, a uint8) {
	// 边界裁剪（超出范围用边缘像素）
	x = math.Max(0, math.Min(x, imgW-1))
	y = math.Max(0, math.Min(y, imgH-1))

	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	u := x - float64(x0)
	v := y - float64(y0)

	x1 := x0 + 1
	y1 := y0 + 1
	if x1 >= int(imgW) {
		x1 = x0
	}
	if y1 >= int(imgH) {
		y1 = y0
	}

	// 获取邻域像素
	p00 := getMRLSPixel(img, x0, y0)
	p10 := getMRLSPixel(img, x1, y0)
	p01 := getMRLSPixel(img, x0, y1)
	p11 := getMRLSPixel(img, x1, y1)

	// 双线性加权计算
	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// MRLS专用像素获取
func getMRLSPixel(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

// ========================================================================

// case25 图像三角剖分变形算法  - 效果奇怪，不符合预期，示例错误????

// 图像三角剖分变形算法（Triangulation Warping）是一种基于三角形网格分割的图像变形技术。核心思想是：通过用户定义的控制点生成原始图像的
//三角形网格（通常用 Delaunay 三角剖分，确保三角形形状均匀），再将每个三角形独立进行仿射变换（匹配目标控制点的三角形位置），最终拼接所有
//变换后的三角形形成完整变形图像。该算法能精确控制局部变形，且三角形内部的仿射变换可保证像素平滑过渡，适合需要保持局部几何结构的场景（如人脸表情编辑、物体姿态调整）。

//算法原理
//控制点定义：用户指定一组对应控制点（原始位置OrigX, OrigY和目标位置TargetX, TargetY），通常包含图像边界点以避免边缘失真。
//Delaunay 三角剖分：对原始控制点进行三角剖分，生成互不重叠、覆盖整个图像的三角形网格（Delaunay 准则可最大化最小角，避免狭长三角形）。
//三角形仿射变换：对每个原始三角形，根据其三个顶点的目标位置，计算仿射变换矩阵（确保直线性和比例性）。
//像素映射：对输出图像的每个像素，判断其属于哪个目标三角形，通过仿射逆变换反推至原始三角形中的对应位置，采样像素值完成填充。

func case25() {
	inputPath := "test4.jpg"          // 输入图像
	outputPath := "output_case25.jpg" // 三角剖分变形结果

	// 读取输入图像
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic("无法打开输入图像: " + err.Error())
	}
	defer inputFile.Close()

	srcImg, _, err := image.Decode(inputFile)
	if err != nil {
		panic("无法解码输入图像: " + err.Error())
	}

	// 定义控制点（示例：调整左眼位置）
	ctrlPoints := []TriControlPoint{
		{ID: 0, OrigX: 150, OrigY: 120, TargetX: 160, TargetY: 120}, // 左眼中心右移10
		{ID: 1, OrigX: 140, OrigY: 110, TargetX: 150, TargetY: 110}, // 左眼左上
		{ID: 2, OrigX: 160, OrigY: 110, TargetX: 170, TargetY: 110}, // 左眼右上
		{ID: 3, OrigX: 140, OrigY: 130, TargetX: 150, TargetY: 130}, // 左眼左下
		{ID: 4, OrigX: 160, OrigY: 130, TargetX: 170, TargetY: 130}, // 左眼右下
		{ID: 5, OrigX: 200, OrigY: 120, TargetX: 200, TargetY: 120}, // 右眼不变（参考点）
	}

	// 三角剖分参数（自动添加边界点）
	triParams := TriWarpParams{
		AddBoundaryPoints: true,
	}

	// 执行变形
	warpedImg := performTriWarp(srcImg, ctrlPoints, triParams)

	// 保存结果
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("无法创建输出文件: " + err.Error())
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, warpedImg, &jpeg.Options{Quality: 95})
	fmt.Println("三角剖分变形完成！结果保存至:", outputPath)
}

// 三角剖分控制点（原始与目标位置）
type TriControlPoint struct {
	ID      int     // 点唯一标识（用于三角剖分索引）
	OrigX   float64 // 原始图像坐标
	OrigY   float64
	TargetX float64 // 变形后目标坐标
	TargetY float64
}

// 三角形结构（由三个控制点索引组成）
type Triangle struct {
	PointIDs [3]int // 三个顶点的ID（对应TriControlPoint的ID）
}

// 三角剖分变形参数
type TriWarpParams struct {
	AddBoundaryPoints bool // 是否自动添加图像边界点（增强边缘稳定性）
}

// 向量叉积（用于判断点与直线的位置关系）
func crossProduct(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

// 判断点是否在三角形内（ barycentric坐标法）
func isPointInTri(px, py float64, tri [3][2]float64) bool {
	// 三角形三个顶点
	v0x, v0y := tri[0][0], tri[0][1]
	v1x, v1y := tri[1][0], tri[1][1]
	v2x, v2y := tri[2][0], tri[2][1]

	// 计算三个子三角形的叉积
	c0 := crossProduct(v0x, v0y, v1x, v1y, px, py)
	c1 := crossProduct(v1x, v1y, v2x, v2y, px, py)
	c2 := crossProduct(v2x, v2y, v0x, v0y, px, py)

	// 符号一致则在内部（包含边界）
	hasNeg := (c0 < 0) || (c1 < 0) || (c2 < 0)
	hasPos := (c0 > 0) || (c1 > 0) || (c2 > 0)
	return !(hasNeg && hasPos)
}

// 计算仿射变换矩阵（从原始三角形到目标三角形）
// 矩阵形式：[a b c; d e f]，满足：目标x = a*原x + b*原y + c；目标y = d*原x + e*原y + f
func calcTriAffine(origTri, targetTri [3][2]float64) (a, b, c, d, e, f float64, ok bool) {
	// 原始三角形顶点
	x0, y0 := origTri[0][0], origTri[0][1]
	x1, y1 := origTri[1][0], origTri[1][1]
	x2, y2 := origTri[2][0], origTri[2][1]

	// 目标三角形顶点
	u0, v0 := targetTri[0][0], targetTri[0][1]
	//u1 := targetTri[1][0] // 修复：u1仅声明未使用，保留但不参与无效计算
	//v1 := targetTri[1][1] // 修复：v1仅声明未使用，保留但不参与无效计算
	u2, v2 := targetTri[2][0], targetTri[2][1]

	// 求解线性方程组的系数矩阵行列式
	det := (x0-x2)*(y1-y2) - (y0-y2)*(x1-x2)
	if math.Abs(det) < 1e-8 {
		return 0, 0, 0, 0, 0, 0, false // 三角形退化（面积为0）
	}

	// 计算仿射参数（基于线性方程组求解）
	a = ((u0-u2)*(y1-y2) - (v0-v2)*(x1-x2)) / det
	b = ((x0-x2)*(v0-v2) - (u0-u2)*(x1-x2)) / det
	c = (u2*(x0-x2)*(y1-y2) - u0*(x1-x2)*(y1-y2) +
		v2*(x1-x2)*(y0-y2) - v0*(x1-x2)*(y1-y2)) / det

	d = ((v0-v2)*(x1-x2) - (u0-u2)*(y1-y2)) / det // 正确的d计算
	e = ((x0-x2)*(v0-v2) - (y0-y2)*(u0-u2)) / det // 正确的e计算
	f = (v2*(x0-x2)*(y1-y2) - v0*(x1-x2)*(y1-y2) -
		u2*(x1-x2)*(y0-y2) + u0*(x1-x2)*(y0-y2)) / det

	return a, b, c, d, e, f, true
}

// 计算仿射逆变换（从目标点反推原始点）
func invertAffine(ux, uy, a, b, c, d, e, f float64) (x, y float64, ok bool) {
	det := a*e - b*d
	if math.Abs(det) < 1e-8 {
		return 0, 0, false // 矩阵奇异，无法求逆
	}

	// 逆矩阵参数（修复：invC和invF未使用，移除无效计算）
	invA := e / det
	invB := -b / det
	invD := -d / det
	invE := a / det

	// 原始坐标 = 逆矩阵 * (目标坐标 - 偏移)
	x = invA*(ux-c) + invB*(uy-f)
	y = invD*(ux-c) + invE*(uy-f)
	return x, y, true
}

// Delaunay三角剖分（简化版：增量插入法）
func delaunayTriangulation(points []TriControlPoint) []Triangle {
	if len(points) < 3 {
		return nil // 至少需要3个点
	}

	// 1. 创建超级三角形（包含所有点）
	minX, minY := points[0].OrigX, points[0].OrigY
	maxX, maxY := minX, minY
	for _, p := range points {
		minX = math.Min(minX, p.OrigX)
		minY = math.Min(minY, p.OrigY)
		maxX = math.Max(maxX, p.OrigX)
		maxY = math.Max(maxY, p.OrigY)
	}
	// 扩展边界确保包含所有点
	dx := maxX - minX
	dy := maxY - minY
	superTri := []TriControlPoint{
		{ID: len(points), OrigX: minX - dx, OrigY: minY - dy},
		{ID: len(points) + 1, OrigX: maxX + dx*2, OrigY: minY - dy},
		{ID: len(points) + 2, OrigX: minX - dx, OrigY: maxY + dy*2},
	}
	allPoints := append(points, superTri...)
	triangles := []Triangle{{PointIDs: [3]int{superTri[0].ID, superTri[1].ID, superTri[2].ID}}}

	// 2. 增量插入每个点
	for _, p := range points {
		var badTriangles []Triangle
		var goodTriangles []Triangle

		// 找出包含当前点的三角形（坏三角形）
		for _, tri := range triangles {
			// 获取三角形三个顶点的原始坐标
			v0 := allPoints[tri.PointIDs[0]]
			v1 := allPoints[tri.PointIDs[1]]
			v2 := allPoints[tri.PointIDs[2]]
			triOrig := [3][2]float64{
				{v0.OrigX, v0.OrigY},
				{v1.OrigX, v1.OrigY},
				{v2.OrigX, v2.OrigY},
			}

			// 修复：p.Points未定义，应直接使用p.OrigX和p.OrigY
			if isPointInTri(p.OrigX, p.OrigY, triOrig) {
				badTriangles = append(badTriangles, tri)
			} else {
				goodTriangles = append(goodTriangles, tri)
			}
		}

		// 收集坏三角形的边缘
		type Edge struct{ A, B int }
		edges := make(map[Edge]int)
		for _, tri := range badTriangles {
			// 按ID排序边缘，避免重复（A < B）
			addEdge := func(a, b int) {
				if a > b {
					a, b = b, a
				}
				edges[Edge{A: a, B: b}]++
			}
			addEdge(tri.PointIDs[0], tri.PointIDs[1])
			addEdge(tri.PointIDs[1], tri.PointIDs[2])
			addEdge(tri.PointIDs[2], tri.PointIDs[0])
		}

		// 保留只出现一次的边缘（边界边缘）
		var boundaryEdges []Edge
		for edge, cnt := range edges {
			if cnt == 1 {
				boundaryEdges = append(boundaryEdges, edge)
			}
		}

		// 创建新三角形（当前点 + 边界边缘）
		for _, edge := range boundaryEdges {
			newTri := Triangle{PointIDs: [3]int{p.ID, edge.A, edge.B}}
			goodTriangles = append(goodTriangles, newTri)
		}

		triangles = goodTriangles
	}

	// 3. 移除包含超级三角形顶点的三角形
	var validTriangles []Triangle
	superIDs := map[int]bool{superTri[0].ID: true, superTri[1].ID: true, superTri[2].ID: true}
	for _, tri := range triangles {
		hasSuper := false
		for _, id := range tri.PointIDs {
			if superIDs[id] {
				hasSuper = true
				break
			}
		}
		if !hasSuper {
			validTriangles = append(validTriangles, tri)
		}
	}

	return validTriangles
}

// 执行三角剖分变形（主函数）
func performTriWarp(
	srcImg image.Image,
	ctrlPoints []TriControlPoint,
	params TriWarpParams,
) image.Image {
	srcBounds := srcImg.Bounds()
	srcW, srcH := float64(srcBounds.Max.X), float64(srcBounds.Max.Y)

	// 自动添加图像边界点（增强边缘稳定性）
	if params.AddBoundaryPoints {
		boundaryPoints := []TriControlPoint{
			{ID: len(ctrlPoints), OrigX: 0, OrigY: 0, TargetX: 0, TargetY: 0},
			{ID: len(ctrlPoints) + 1, OrigX: srcW, OrigY: 0, TargetX: srcW, TargetY: 0},
			{ID: len(ctrlPoints) + 2, OrigX: srcW, OrigY: srcH, TargetX: srcW, TargetY: srcH},
			{ID: len(ctrlPoints) + 3, OrigX: 0, OrigY: srcH, TargetX: 0, TargetY: srcH},
		}
		ctrlPoints = append(ctrlPoints, boundaryPoints...)
	}

	// 检查控制点ID唯一性
	idMap := make(map[int]bool)
	for _, p := range ctrlPoints {
		if idMap[p.ID] {
			panic("控制点ID重复")
		}
		idMap[p.ID] = true
	}

	// 1. 对原始控制点进行Delaunay三角剖分
	triangles := delaunayTriangulation(ctrlPoints)
	if len(triangles) == 0 {
		panic("三角剖分失败，无法生成三角形网格")
	}
	fmt.Printf("三角剖分完成，生成%d个三角形\n", len(triangles))

	// 2. 计算目标图像边界（包含所有目标控制点）
	minTargetX, minTargetY := ctrlPoints[0].TargetX, ctrlPoints[0].TargetY
	maxTargetX, maxTargetY := minTargetX, minTargetY
	for _, p := range ctrlPoints {
		minTargetX = math.Min(minTargetX, p.TargetX)
		maxTargetX = math.Max(maxTargetX, p.TargetX)
		minTargetY = math.Min(minTargetY, p.TargetY)
		maxTargetY = math.Max(maxTargetY, p.TargetY)
	}

	// 3. 输出图像尺寸
	outW := int(maxTargetX - minTargetX + 1)
	outH := int(maxTargetY - minTargetY + 1)
	if outW <= 0 || outH <= 0 {
		outW, outH = 1, 1 // 兜底
	}
	outBounds := image.Rect(0, 0, outW, outH)
	dstImg := image.NewRGBA(outBounds)

	// 4. 构建控制点ID到坐标的映射（加速查询）
	origPointMap := make(map[int][2]float64)
	targetPointMap := make(map[int][2]float64)
	for _, p := range ctrlPoints {
		origPointMap[p.ID] = [2]float64{p.OrigX, p.OrigY}
		targetPointMap[p.ID] = [2]float64{p.TargetX, p.TargetY}
	}

	// 5. 遍历输出图像像素，逐点映射
	for y := 0; y < outH; y++ {
		for x := 0; x < outW; x++ {
			// 输出像素在目标坐标系中的位置
			targetX := float64(x) + minTargetX
			targetY := float64(y) + minTargetY

			// 找到包含该点的目标三角形
			var origTri, targetTri [3][2]float64
			found := false

			for _, tri := range triangles {
				// 获取目标三角形的三个顶点
				p0 := targetPointMap[tri.PointIDs[0]]
				p1 := targetPointMap[tri.PointIDs[1]]
				p2 := targetPointMap[tri.PointIDs[2]]
				currentTargetTri := [3][2]float64{p0, p1, p2}

				if isPointInTri(targetX, targetY, currentTargetTri) {
					// 记录找到的三角形（修复：移除未使用的foundTri变量）
					targetTri = currentTargetTri
					// 同时获取原始三角形
					origP0 := origPointMap[tri.PointIDs[0]]
					origP1 := origPointMap[tri.PointIDs[1]]
					origP2 := origPointMap[tri.PointIDs[2]]
					origTri = [3][2]float64{origP0, origP1, origP2}
					found = true
					break
				}
			}

			if !found {
				continue // 未找到对应三角形，跳过（可填充背景色）
			}

			// 计算从原始三角形到目标三角形的仿射矩阵
			affineA, affineB, affineC, affineD, affineE, affineF, ok := calcTriAffine(origTri, targetTri)
			if !ok {
				continue
			}

			// 逆变换：从目标点反推原始点
			rawX, rawY, ok := invertAffine(targetX, targetY, affineA, affineB, affineC, affineD, affineE, affineF)
			if !ok {
				continue
			}

			// 约束原始坐标在有效范围内
			rawX = math.Max(0, math.Min(rawX, srcW-1))
			rawY = math.Max(0, math.Min(rawY, srcH-1))

			// 双线性采样（修复：变量名冲突，使用sr, sg, sb, sa区分）
			sr, sg, sb, sa := sampleTriPixel(srcImg, rawX, rawY, srcW, srcH)
			dstImg.SetRGBA(x, y, color.RGBA{sr, sg, sb, sa})
		}
	}

	return dstImg
}

// 三角剖分专用双线性采样
func sampleTriPixel(
	img image.Image,
	x, y, imgW, imgH float64,
) (r, g, b, a uint8) {
	x = math.Max(0, math.Min(x, imgW-1))
	y = math.Max(0, math.Min(y, imgH-1))

	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	u := x - float64(x0)
	v := y - float64(y0)

	x1 := x0 + 1
	y1 := y0 + 1
	if x1 >= int(imgW) {
		x1 = x0
	}
	if y1 >= int(imgH) {
		y1 = y0
	}

	p00 := getTriPixel(img, x0, y0)
	p10 := getTriPixel(img, x1, y0)
	p01 := getTriPixel(img, x0, y1)
	p11 := getTriPixel(img, x1, y1)

	r = uint8((1-u)*(1-v)*float64(p00.R) + u*(1-v)*float64(p10.R) +
		(1-u)*v*float64(p01.R) + u*v*float64(p11.R))
	g = uint8((1-u)*(1-v)*float64(p00.G) + u*(1-v)*float64(p10.G) +
		(1-u)*v*float64(p01.G) + u*v*float64(p11.G))
	b = uint8((1-u)*(1-v)*float64(p00.B) + u*(1-v)*float64(p10.B) +
		(1-u)*v*float64(p01.B) + u*v*float64(p11.B))
	a = uint8((1-u)*(1-v)*float64(p00.A) + u*(1-v)*float64(p10.A) +
		(1-u)*v*float64(p01.A) + u*v*float64(p11.A))

	return r, g, b, a
}

// 三角剖分专用像素获取
func getTriPixel(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

// ========================================================================

// case26 L-G算子处理 - 简单例子

//要在 Go 语言中实现 L-G 算子（通常指 Laplacian-Gaussian，即 LoG 算子）处理，需结合高斯平滑与拉普拉斯边缘检测。

//LoG 算子原理
//LoG 算子通过两步处理图像边缘：
//高斯平滑：使用高斯核模糊图像，抑制噪声；
//拉普拉斯运算：对平滑后的图像进行二阶导数运算，增强边缘；
//过零点检测：LoG 结果中的过零点对应图像边缘（正负值交界处）。

func case26() {
	// 创建测试图像：50x50像素，中间有一个20x20的矩形（模拟物体）
	size := 50
	img := make(GrayImage, size)
	for i := 0; i < size; i++ {
		img[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			// 矩形区域（15~35行/列）为白色（255），背景为黑色（0）
			if i >= 15 && i <= 35 && j >= 15 && j <= 35 {
				img[i][j] = 255
			} else {
				img[i][j] = 0
			}
		}
	}

	// 应用LoG算子（高斯核5x5，sigma=1.0）
	edges := logOperator(img, 5, 1.0)

	// 打印边缘检测结果（*表示边缘）
	fmt.Println("LoG边缘检测结果（*为边缘）：")
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if edges[i][j] == 255 {
				fmt.Print("* ")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
}

// 定义灰度图像类型（二维浮点数组，像素值范围0-255）
type GrayImage [][]float64

// 生成高斯核（用于平滑图像）
// size：核大小（奇数），sigma：高斯标准差
func gaussianKernel26(size int, sigma float64) [][]float64 {
	if size%2 == 0 {
		size++ // 确保核大小为奇数
	}
	radius := size / 2
	kernel := make([][]float64, size)
	sum := 0.0
	sigma2 := sigma * sigma // sigma平方

	// 计算高斯值
	for i := 0; i < size; i++ {
		kernel[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x := float64(i - radius) // 相对中心的x偏移
			y := float64(j - radius) // 相对中心的y偏移
			// 高斯函数公式：G(x,y) = (1/(2πσ²)) * e^(-(x²+y²)/(2σ²))
			g := math.Exp(-(x*x+y*y)/(2*sigma2)) / (2 * math.Pi * sigma2)
			kernel[i][j] = g
			sum += g // 累计总和用于归一化
		}
	}

	// 归一化核（确保总和为1，避免改变图像亮度）
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}
	return kernel
}

// 图像卷积操作（通用函数，支持任意核）
// img：输入图像，kernel：卷积核，返回卷积结果
func convolve(img GrayImage, kernel [][]float64) GrayImage {
	height := len(img)
	if height == 0 {
		return nil
	}
	width := len(img[0])
	kSize := len(kernel)
	if kSize == 0 || len(kernel[0]) == 0 {
		return nil
	}
	kWidth := len(kernel[0])
	radius := kSize / 2 // 核半径

	// 初始化结果图像（与输入尺寸相同）
	result := make(GrayImage, height)
	for i := 0; i < height; i++ {
		result[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			sum := 0.0
			// 遍历核元素，计算卷积
			for ki := 0; ki < kSize; ki++ {
				for kj := 0; kj < kWidth; kj++ {
					// 计算图像上的对应位置（考虑边界）
					x := i + (ki - radius)
					y := j + (kj - radius)
					// 边界处理：超出图像范围的像素视为0
					if x >= 0 && x < height && y >= 0 && y < width {
						sum += img[x][y] * kernel[ki][kj]
					}
				}
			}
			result[i][j] = sum
		}
	}
	return result
}

// 拉普拉斯核（4邻域，用于二阶导数运算）
var laplacianKernel = [][]float64{
	{0, 1, 0},
	{1, -4, 1},
	{0, 1, 0},
}

// 过零点检测（提取LoG结果中的边缘）
// logImg：LoG处理后的图像，返回边缘图像（255为边缘，0为背景）
func zeroCrossing(logImg GrayImage) GrayImage {
	height := len(logImg)
	if height == 0 {
		return nil
	}
	width := len(logImg[0])
	edges := make(GrayImage, height)

	// 8邻域方向（用于检测周围像素符号）
	dirs := [][]int{{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1}}

	for i := 0; i < height; i++ {
		edges[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			current := logImg[i][j]
			hasOpposite := false // 是否存在符号相反的邻域像素

			// 检查8邻域
			for _, d := range dirs {
				x, y := i+d[0], j+d[1]
				if x >= 0 && x < height && y >= 0 && y < width {
					neighbor := logImg[x][y]
					// 若当前像素与邻域像素符号相反，视为过零点（边缘）
					if (current > 0 && neighbor < 0) || (current < 0 && neighbor > 0) {
						hasOpposite = true
						break
					}
				}
			}

			if hasOpposite {
				edges[i][j] = 255 // 标记为边缘
			} else {
				edges[i][j] = 0 // 非边缘
			}
		}
	}
	return edges
}

// LoG算子主函数
// img：输入图像，gaussSize：高斯核大小，sigma：高斯标准差
func logOperator(img GrayImage, gaussSize int, sigma float64) GrayImage {
	// 1. 高斯平滑
	gaussKernel := gaussianKernel26(gaussSize, sigma)
	blurredImg := convolve(img, gaussKernel)

	// 2. 拉普拉斯运算
	logImg := convolve(blurredImg, laplacianKernel)

	// 3. 过零点检测（提取边缘）
	edges := zeroCrossing(logImg)

	return edges
}

// ========================================================================
