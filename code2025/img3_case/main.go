package main

import (
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

	case14()

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

// ========================================================================

// case16 人像磨皮 通到磨皮算法

// ========================================================================

// case17 人像磨皮 高反差磨皮算法

// ========================================================================

// case18 人像磨皮 细节叠加磨皮算法

// ========================================================================

// case19 图像放射变换

// ========================================================================

// case20 图像透视变换

// ========================================================================

// case21 图像反距离加权(IDW)插值变形算法

// ========================================================================

// case22 图像特征线变换算法

// ========================================================================

// case23 图像MLS变形算法

// ========================================================================

// case24 图像MRLS算法

// ========================================================================

// case25 图像三角剖分变形算法

// ========================================================================

// case26 人像分割

// ========================================================================

// case27 背景虚化

// ========================================================================

// case28 L-G算子处理

// ========================================================================
