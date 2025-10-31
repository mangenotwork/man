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

	//case1()

	//case2()

	//case3()

	//case4()

	//case5()

	case6()
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

// ========================================================================

// case8 BEEPS滤波算法

// ========================================================================

// case9 DCT降噪滤波算法

// ========================================================================

// case10 非局部均值滤波算法

// ========================================================================

// case11 加权中值滤波算法

// ========================================================================

// case12 基于颜色空间的皮肤检测算法

// ========================================================================

// case13 基于高斯模型的肤色概率计算方法

// ========================================================================

// case14 皮肤美白算法 LUT调色法

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
