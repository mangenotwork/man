package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/gif" // 使用 _ 导入 image/gif、image/jpeg 和 image/png 包，这会调用这些包的 init 函数，从而注册相应的图像格式
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

func main() {

	// case1()

}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// convert image to NRGBA
func convertToNRGBA(src image.Image) *image.NRGBA {
	srcBounds := src.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

func DecodeImage(filePath string) (img image.Image, err error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err = image.Decode(reader)

	return
}

// 根据文件名打开图片,并编码,返回编码对象和文件类型
func LoadImage(path string) (img image.Image, filetype string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, filetype, err = image.Decode(file)
	return
}

func case1() {
	_, fileType, err := LoadImage("./test2.jpg")
	log.Println("err = ", err)
	log.Println("fileType = ", fileType)
}

// 如果您有未知格式的图像数据，请使用图像。解码功能可以检测格式。公认的格式集是在运行时构建的，不限于标准包库中的格式。
// 图像格式包通常在init函数中注册其格式，主包将“下划线导入”此类包，仅用于格式注册的副作用。
// convertToPNG converts from any recognized format to PNG.
func convertToPNG(w io.Writer, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

//	=========================================================================================================  颜色融合
//
// Color 表示一个RGBA颜色
type Color struct {
	R, G, B, A float64 // 范围从0.0到1.0
}

// NewColor 创建一个新的RGBA颜色
func NewColor(r, g, b, a float64) Color {
	return Color{
		R: clamp(r, 0.0, 1.0),
		G: clamp(g, 0.0, 1.0),
		B: clamp(b, 0.0, 1.0),
		A: clamp(a, 0.0, 1.0),
	}
}

// 辅助函数：将值限制在min和max之间
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Add 颜色加法混合
func (c Color) Add(other Color) Color {
	return Color{
		R: clamp(c.R+other.R, 0.0, 1.0),
		G: clamp(c.G+other.G, 0.0, 1.0),
		B: clamp(c.B+other.B, 0.0, 1.0),
		A: clamp(c.A+other.A, 0.0, 1.0),
	}
}

// Subtract 颜色减法混合
func (c Color) Subtract(other Color) Color {
	return Color{
		R: clamp(c.R-other.R, 0.0, 1.0),
		G: clamp(c.G-other.G, 0.0, 1.0),
		B: clamp(c.B-other.B, 0.0, 1.0),
		A: clamp(c.A-other.A, 0.0, 1.0),
	}
}

// Multiply 颜色乘法混合（常用于变暗）
func (c Color) Multiply(other Color) Color {
	return Color{
		R: clamp(c.R*other.R, 0.0, 1.0),
		G: clamp(c.G*other.G, 0.0, 1.0),
		B: clamp(c.B*other.B, 0.0, 1.0),
		A: clamp(c.A*other.A, 0.0, 1.0),
	}
}

// Blend 颜色混合，根据比例混合两种颜色
func (c Color) Blend(other Color, ratio float64) Color {
	ratio = clamp(ratio, 0.0, 1.0)
	invRatio := 1.0 - ratio
	return Color{
		R: c.R*invRatio + other.R*ratio,
		G: c.G*invRatio + other.G*ratio,
		B: c.B*invRatio + other.B*ratio,
		A: c.A*invRatio + other.A*ratio,
	}
}

// String 返回颜色的字符串表示形式（用于调试）
func (c Color) String() string {
	return fmt.Sprintf("RGBA(%d, %d, %d, %.2f)",
		int(c.R*255), int(c.G*255), int(c.B*255), c.A)
}

// ToHex 返回颜色的十六进制表示形式
func (c Color) ToHex() string {
	return fmt.Sprintf("#%02X%02X%02X",
		int(c.R*255), int(c.G*255), int(c.B*255))
}

// 从HSV颜色空间转换到RGB
func HSVtoRGB(h, s, v float64) Color {
	h = math.Mod(h, 360.0)
	if h < 0 {
		h += 360.0
	}
	s = clamp(s, 0.0, 1.0)
	v = clamp(v, 0.0, 1.0)

	if s == 0 {
		return NewColor(v, v, v, 1.0)
	}

	h /= 60.0
	hi := int(h) % 6
	f := h - float64(int(h))
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	var r, g, b float64
	switch hi {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return NewColor(r, g, b, 1.0)
}

func case2() {
	// 定义几个基本颜色变量
	red := NewColor(1.0, 0.0, 0.0, 1.0)
	green := NewColor(0.0, 1.0, 0.0, 1.0)
	blue := NewColor(0.0, 0.0, 1.0, 1.0)
	yellow := NewColor(1.0, 1.0, 0.0, 1.0)

	// 创建一个从HSV转换的颜色
	purple := HSVtoRGB(270, 0.8, 0.8)

	// 演示颜色混合
	fmt.Println("基本颜色:")
	fmt.Println("红色:", red.ToHex())
	fmt.Println("绿色:", green.ToHex())
	fmt.Println("蓝色:", blue.ToHex())
	fmt.Println("黄色:", yellow.ToHex())
	fmt.Println("紫色:", purple.ToHex())

	fmt.Println("\n颜色混合示例:")
	// 混合红色和绿色
	yellowMix := red.Blend(green, 0.5)
	fmt.Println("红色 + 绿色 =", yellowMix.ToHex(), yellowMix)

	// 混合蓝色和黄色
	greenMix := blue.Blend(yellow, 0.5)
	fmt.Println("蓝色 + 黄色 =", greenMix.ToHex(), greenMix)

	// 混合所有基本颜色
	mixed := red.Blend(green, 0.25).Blend(blue, 0.25).Blend(yellow, 0.25)
	fmt.Println("所有颜色混合 =", mixed.ToHex(), mixed)

	// 颜色减法
	darker := red.Subtract(NewColor(0.3, 0.0, 0.0, 1.0))
	fmt.Println("红色变暗 =", darker.ToHex(), darker)

	// 颜色乘法
	darker2 := red.Multiply(NewColor(0.7, 0.7, 0.7, 1.0))
	fmt.Println("红色乘暗 =", darker2.ToHex(), darker2)
}
