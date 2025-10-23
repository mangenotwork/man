//package main
//
//import (
//	"image"
//	"image/color"
//	"image/draw"
//	"image/jpeg"
//	"math"
//	"os"
//)
//
//// 高斯模糊（仅处理局部区域）
//func gaussianBlur(img *image.RGBA, radius int) *image.RGBA {
//	bounds := img.Bounds()
//	width, height := bounds.Max.X, bounds.Max.Y
//	output := image.NewRGBA(bounds)
//	kernelSize := radius*2 + 1
//	sigma := float64(radius) * 0.577
//
//	// 生成高斯核
//	kernel := make([]float64, kernelSize)
//	for i := 0; i < kernelSize; i++ {
//		x := float64(i - radius)
//		kernel[i] = math.Exp(-(x*x)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
//	}
//	// 归一化
//	var sum float64
//	for _, v := range kernel {
//		sum += v
//	}
//	for i := range kernel {
//		kernel[i] /= sum
//	}
//
//	// 横向模糊
//	temp := image.NewRGBA(bounds)
//	for y := 0; y < height; y++ {
//		for x := 0; x < width; x++ {
//			var r, g, b, a float64
//			for k := 0; k < kernelSize; k++ {
//				kx := x + k - radius
//				if kx < 0 {
//					kx = 0
//				} else if kx >= width {
//					kx = width - 1
//				}
//				px := img.RGBAAt(kx, y)
//				weight := kernel[k]
//				r += float64(px.R) * weight
//				g += float64(px.G) * weight
//				b += float64(px.B) * weight
//				a += float64(px.A) * weight
//			}
//			temp.SetRGBA(x, y, color.RGBA{
//				R: uint8(clamp(r, 0, 255)),
//				G: uint8(clamp(g, 0, 255)),
//				B: uint8(clamp(b, 0, 255)),
//				A: uint8(clamp(a, 0, 255)),
//			})
//		}
//	}
//
//	// 纵向模糊
//	for x := 0; x < width; x++ {
//		for y := 0; y < height; y++ {
//			var r, g, b, a float64
//			for k := 0; k < kernelSize; k++ {
//				ky := y + k - radius
//				if ky < 0 {
//					ky = 0
//				} else if ky >= height {
//					ky = height - 1
//				}
//				px := temp.RGBAAt(x, ky)
//				weight := kernel[k]
//				r += float64(px.R) * weight
//				g += float64(px.G) * weight
//				b += float64(px.B) * weight
//				a += float64(px.A) * weight
//			}
//			output.SetRGBA(x, y, color.RGBA{
//				R: uint8(clamp(r, 0, 255)),
//				G: uint8(clamp(g, 0, 255)),
//				B: uint8(clamp(b, 0, 255)),
//				A: uint8(clamp(a, 0, 255)),
//			})
//		}
//	}
//
//	return output
//}
//
//// 手动复制指定区域像素
//func copyRegion(original *image.RGBA, x, y, w, h int) *image.RGBA {
//	dst := image.NewRGBA(image.Rect(0, 0, w, h))
//	for dy := 0; dy < h; dy++ {
//		for dx := 0; dx < w; dx++ {
//			srcX := x + dx
//			srcY := y + dy
//			px := original.RGBAAt(srcX, srcY)
//			dst.SetRGBA(dx, dy, px)
//		}
//	}
//	return dst
//}
//
//// 判断点是否在圆角矩形内（相对坐标）
//// dx, dy: 区域内的相对坐标（0~w-1, 0~h-1）
//// w, h: 矩形宽高
//// radius: 圆角半径
//func isInRoundedRect(dx, dy, w, h, radius int) bool {
//	// 限制圆角半径最大值（不超过宽高的一半）
//	maxRadius := min(w, h) / 2
//	if radius > maxRadius {
//		radius = maxRadius
//	}
//
//	// 四个角的圆角判断
//	// 左上角：x < radius && y < radius
//	if dx < radius && dy < radius {
//		// 到左上角的距离是否小于等于半径
//		dx2 := dx - radius
//		dy2 := dy - radius
//		return dx2*dx2+dy2*dy2 <= radius*radius
//	}
//	// 右上角：x >= w-radius && y < radius
//	if dx >= w-radius && dy < radius {
//		dx2 := dx - (w - radius)
//		dy2 := dy - radius
//		return dx2*dx2+dy2*dy2 <= radius*radius
//	}
//	// 左下角：x < radius && y >= h-radius
//	if dx < radius && dy >= h-radius {
//		dx2 := dx - radius
//		dy2 := dy - (h - radius)
//		return dx2*dx2+dy2*dy2 <= radius*radius
//	}
//	// 右下角：x >= w-radius && y >= h-radius
//	if dx >= w-radius && dy >= h-radius {
//		dx2 := dx - (w - radius)
//		dy2 := dy - (h - radius)
//		return dx2*dx2+dy2*dy2 <= radius*radius
//	}
//
//	// 非角区域：在矩形内部即可
//	return true
//}
//
//func min(a, b int) int {
//	if a < b {
//		return a
//	}
//	return b
//}
//
//// 带圆角的毛玻璃效果
//func drawRoundedFrostedGlass(original *image.RGBA, x, y, w, h, blurRadius, cornerRadius int, bgColor color.RGBA) *image.RGBA {
//	bounds := original.Bounds()
//	imgW, imgH := bounds.Max.X, bounds.Max.Y
//	result := image.NewRGBA(bounds)
//	draw.Draw(result, bounds, original, image.Point{}, draw.Over)
//
//	// 校验区域是否在图像内
//	if x < 0 || y < 0 || x+w > imgW || y+h > imgH {
//		return result
//	}
//
//	// 1. 复制指定区域像素
//	background := copyRegion(original, x, y, w, h)
//
//	// 2. 应用模糊
//	blurredBg := gaussianBlur(background, blurRadius)
//
//	// 3. 半透明背景色混合
//	bgR8, bgG8, bgB8 := bgColor.R, bgColor.G, bgColor.B
//	bgAlpha := float64(bgColor.A) / 255.0
//
//	for dy := 0; dy < h; dy++ {
//		for dx := 0; dx < w; dx++ {
//			px := blurredBg.RGBAAt(dx, dy)
//			mixedR := uint8(float64(px.R)*(1-bgAlpha) + float64(bgR8)*bgAlpha)
//			mixedG := uint8(float64(px.G)*(1-bgAlpha) + float64(bgG8)*bgAlpha)
//			mixedB := uint8(float64(px.B)*(1-bgAlpha) + float64(bgB8)*bgAlpha)
//			blurredBg.SetRGBA(dx, dy, color.RGBA{mixedR, mixedG, mixedB, px.A})
//		}
//	}
//
//	// 4. 带圆角绘制回原图（逐像素判断）
//	for dy := 0; dy < h; dy++ {
//		for dx := 0; dx < w; dx++ {
//			// 判断是否在圆角矩形内
//			if isInRoundedRect(dx, dy, w, h, cornerRadius) {
//				// 绘制到原图的绝对坐标 (x+dx, y+dy)
//				px := blurredBg.RGBAAt(dx, dy)
//				result.SetRGBA(x+dx, y+dy, px)
//			}
//		}
//	}
//
//	return result
//}
//
//func clamp(v, min, max float64) float64 {
//	if v < min {
//		return min
//	}
//	if v > max {
//		return max
//	}
//	return v
//}
//
//func main() {
//	input, err := os.Open("input.jpg")
//	if err != nil {
//		panic(err)
//	}
//	defer input.Close()
//
//	img, _, err := image.Decode(input)
//	if err != nil {
//		panic(err)
//	}
//	rgbaImg := image.NewRGBA(img.Bounds())
//	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)
//
//	// 绘制带圆角的毛玻璃
//	frosted := drawRoundedFrostedGlass(
//		rgbaImg,
//		300, 150, // 位置(x,y)
//		600, 600, // 宽高(w,h)
//		20,                            // 模糊半径
//		36,                            // 圆角半径（建议5-30）
//		color.RGBA{255, 255, 255, 64}, // 半透明背景色
//	)
//
//	output, err := os.Create("frosted_rounded3.jpg")
//	if err != nil {
//		panic(err)
//	}
//	defer output.Close()
//	jpeg.Encode(output, frosted, &jpeg.Options{Quality: 95})
//}

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

// 高斯模糊（复用原有逻辑，支持对任意RGBA图像模糊）
func gaussianBlur(img *image.RGBA, radius int) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	output := image.NewRGBA(bounds)
	kernelSize := radius*2 + 1
	sigma := float64(radius) * 0.577

	kernel := make([]float64, kernelSize)
	for i := 0; i < kernelSize; i++ {
		x := float64(i - radius)
		kernel[i] = math.Exp(-(x*x)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
	}
	var sum float64
	for _, v := range kernel {
		sum += v
	}
	for i := range kernel {
		kernel[i] /= sum
	}

	temp := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a float64
			for k := 0; k < kernelSize; k++ {
				kx := x + k - radius
				if kx < 0 {
					kx = 0
				} else if kx >= width {
					kx = width - 1
				}
				px := img.RGBAAt(kx, y)
				weight := kernel[k]
				r += float64(px.R) * weight
				g += float64(px.G) * weight
				b += float64(px.B) * weight
				a += float64(px.A) * weight
			}
			temp.SetRGBA(x, y, color.RGBA{
				R: uint8(clamp(r, 0, 255)),
				G: uint8(clamp(g, 0, 255)),
				B: uint8(clamp(b, 0, 255)),
				A: uint8(clamp(a, 0, 255)),
			})
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var r, g, b, a float64
			for k := 0; k < kernelSize; k++ {
				ky := y + k - radius
				if ky < 0 {
					ky = 0
				} else if ky >= height {
					ky = height - 1
				}
				px := temp.RGBAAt(x, ky)
				weight := kernel[k]
				r += float64(px.R) * weight
				g += float64(px.G) * weight
				b += float64(px.B) * weight
				a += float64(px.A) * weight
			}
			output.SetRGBA(x, y, color.RGBA{
				R: uint8(clamp(r, 0, 255)),
				G: uint8(clamp(g, 0, 255)),
				B: uint8(clamp(b, 0, 255)),
				A: uint8(clamp(a, 0, 255)),
			})
		}
	}

	return output
}

// 复制原图指定区域像素
func copyRegion(original *image.RGBA, x, y, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			srcX := x + dx
			srcY := y + dy
			px := original.RGBAAt(srcX, srcY)
			dst.SetRGBA(dx, dy, px)
		}
	}
	return dst
}

// 判断点是否在圆角矩形内（相对坐标）
func isInRoundedRect(dx, dy, w, h, radius int) bool {
	maxRadius := min(w, h) / 2
	if radius > maxRadius {
		radius = maxRadius
	}

	if dx < radius && dy < radius {
		dx2 := dx - radius
		dy2 := dy - radius
		return dx2*dx2+dy2*dy2 <= radius*radius
	}
	if dx >= w-radius && dy < radius {
		dx2 := dx - (w - radius)
		dy2 := dy - radius
		return dx2*dx2+dy2*dy2 <= radius*radius
	}
	if dx < radius && dy >= h-radius {
		dx2 := dx - radius
		dy2 := dy - (h - radius)
		return dx2*dx2+dy2*dy2 <= radius*radius
	}
	if dx >= w-radius && dy >= h-radius {
		dx2 := dx - (w - radius)
		dy2 := dy - (h - radius)
		return dx2*dx2+dy2*dy2 <= radius*radius
	}

	return true
}

// 绘制阴影（独立函数）
// x,y: 毛玻璃区域左上角坐标
// w,h: 毛玻璃宽高
// cornerRadius: 圆角半径（需与毛玻璃一致）
// shadowDx, shadowDy: 阴影偏移量（右、下方向）
// shadowBlur: 阴影模糊半径
// shadowColor: 阴影颜色（建议半透明黑色）
func drawShadow(original *image.RGBA, x, y, w, h, cornerRadius, shadowDx, shadowDy, shadowBlur int, shadowColor color.RGBA) *image.RGBA {
	bounds := original.Bounds()
	//imgW, imgH := bounds.Max.X, bounds.Max.Y
	// 创建阴影临时图像（比毛玻璃大，预留模糊空间）
	shadowW := w + shadowBlur*2
	shadowH := h + shadowBlur*2
	shadowImg := image.NewRGBA(image.Rect(0, 0, shadowW, shadowH))

	// 1. 绘制阴影形状（与毛玻璃同圆角，居中绘制在shadowImg中）
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			if isInRoundedRect(dx, dy, w, h, cornerRadius) {
				// 阴影形状在shadowImg中的位置（居中，预留模糊空间）
				shadowX := dx + shadowBlur
				shadowY := dy + shadowBlur
				shadowImg.SetRGBA(shadowX, shadowY, shadowColor)
			}
		}
	}

	// 2. 对阴影形状应用模糊
	blurredShadow := gaussianBlur(shadowImg, shadowBlur)

	// 3. 绘制阴影到原图（阴影位置 = 毛玻璃位置 + 偏移量 - 模糊预留空间）
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, original, image.Point{}, draw.Over)
	shadowDrawX := x + shadowDx - shadowBlur
	shadowDrawY := y + shadowDy - shadowBlur
	shadowDrawRect := image.Rect(shadowDrawX, shadowDrawY, shadowDrawX+shadowW, shadowDrawY+shadowH)
	// 确保阴影不超出原图边界
	shadowDrawRect = shadowDrawRect.Intersect(bounds)
	draw.Draw(result, shadowDrawRect, blurredShadow, image.Point{shadowBlur - shadowDx, shadowBlur - shadowDy}, draw.Over)

	return result
}

// 带圆角和阴影的毛玻璃效果
func drawRoundedFrostedGlassWithShadow(original *image.RGBA, x, y, w, h, blurRadius, cornerRadius int, bgColor color.RGBA, shadowDx, shadowDy, shadowBlur int, shadowColor color.RGBA) *image.RGBA {
	bounds := original.Bounds()
	imgW, imgH := bounds.Max.X, bounds.Max.Y

	// 1. 先绘制阴影（底层）
	result := drawShadow(original, x, y, w, h, cornerRadius, shadowDx, shadowDy, shadowBlur, shadowColor)

	// 2. 校验毛玻璃区域是否在图像内
	if x < 0 || y < 0 || x+w > imgW || y+h > imgH {
		return result
	}

	// 3. 复制毛玻璃区域像素并模糊
	background := copyRegion(original, x, y, w, h)
	blurredBg := gaussianBlur(background, blurRadius)

	// 4. 混合半透明背景色
	bgR8, bgG8, bgB8 := bgColor.R, bgColor.G, bgColor.B
	bgAlpha := float64(bgColor.A) / 255.0
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px := blurredBg.RGBAAt(dx, dy)
			mixedR := uint8(float64(px.R)*(1-bgAlpha) + float64(bgR8)*bgAlpha)
			mixedG := uint8(float64(px.G)*(1-bgAlpha) + float64(bgG8)*bgAlpha)
			mixedB := uint8(float64(px.B)*(1-bgAlpha) + float64(bgB8)*bgAlpha)
			blurredBg.SetRGBA(dx, dy, color.RGBA{mixedR, mixedG, mixedB, px.A})
		}
	}

	// 5. 绘制带圆角的毛玻璃（上层）
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			if isInRoundedRect(dx, dy, w, h, cornerRadius) {
				px := blurredBg.RGBAAt(dx, dy)
				result.SetRGBA(x+dx, y+dy, px)
			}
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func main() {
	input, err := os.Open("input.jpg")
	if err != nil {
		panic(err)
	}
	defer input.Close()

	img, _, err := image.Decode(input)
	if err != nil {
		panic(err)
	}
	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)

	// 绘制带圆角和阴影的毛玻璃
	frosted := drawRoundedFrostedGlassWithShadow(
		rgbaImg,
		200, 150, // 毛玻璃位置(x,y)
		600, 600, // 毛玻璃宽高(w,h)
		20,                            // 毛玻璃模糊半径
		36,                            // 圆角半径
		color.RGBA{255, 255, 255, 64}, // 毛玻璃背景色（半透明白）
		10, 10,                        // 阴影偏移量（向右5px，向下5px）
		8,                       // 阴影模糊半径（建议5-15）
		color.RGBA{0, 0, 0, 40}, // 阴影颜色（半透明黑）
	)

	output, err := os.Create("frosted_rounded_shadow_2.jpg")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	jpeg.Encode(output, frosted, &jpeg.Options{Quality: 95})
}
