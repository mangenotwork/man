/*

FreeType是一个可移植的，高效的字体引擎。

字体在电脑上的显示有两种方式：点阵和矢量。对于一个字，点阵字体保存的是每个点的渲染信息。这个方式的劣势在于保存的数据量非常大，并且对放大缩小等操作支持不好。
因此出现了矢量字体。对于一个字，矢量字体保存的是字的绘制公式。这个绘制公式包括了字体轮廓（outline）和字体精调（hint）。字体轮廓使用贝塞尔曲线来绘制出字
的外部线条。在大分辨率的情况下就需要对字体进行精调了。这个绘制字的公式就叫做字体数据（glyph）。在字体文件中，每个字对应一个glyph。那么字体文件中就存在一
个字符映射表（charmap）。

对于矢量字体，其中用的最为广泛的是TrueType。它的扩展名一般为otf或者ttf。在windows，linux，osx上都得到广泛支持。我们平时看到的.ttf和.ttc的字体文
件就是TrueType字体。其中ttc是多个ttf的集合文件（collection）。


*/

package main

import (
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

const (
	dx       = 100                                    // 图片的大小 宽度
	dy       = 40                                     // 图片的大小 高度
	fontFile = "AlimamaFangYuanTiVF-Thin.3b91e1b.ttf" // 需要使用的字体文件
	fontSize = 20                                     // 字体尺寸
	fontDPI  = 72                                     // 屏幕每英寸的分辨率
)

func main() {

	// 需要保存的文件
	imgcounter := 123
	imgfile, _ := os.Create(fmt.Sprintf("%03d.png", imgcounter))
	defer imgfile.Close()

	// 新建一个 指定大小的 RGBA位图
	img := image.NewNRGBA(image.Rect(0, 0, dx, dy))

	// 画背景
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			// 设置某个点的颜色，依次是 RGBA
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}

	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(fontSize)>>8)) // 字出现的位置

	_, err = c.DrawString("ABCDE", pt)
	if err != nil {
		log.Println(err)
		return
	}

	// 以PNG格式保存文件
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}

}
