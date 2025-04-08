package draw

import (
	"bytes"
	"ecosmos-api/common/logger"
	drawx "golang.org/x/image/draw"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var manCharacterMap = map[int]string{
	0:  "./assets/images/business_card/manNor-1-min.png",
	1:  "./assets/images/business_card/manNor-2-min.png",
	2:  "./assets/images/business_card/manNor-3-min.png",
	3:  "./assets/images/business_card/manNor-4-min.png",
	4:  "./assets/images/business_card/manNor-5-min.png",
	5:  "./assets/images/business_card/manNor-6-min.png",
	6:  "./assets/images/business_card/manNor-7-min.png",
	7:  "./assets/images/business_card/manNor-8-min.png",
	8:  "./assets/images/business_card/manNor-9-min.png",
	9:  "./assets/images/business_card/manNor-10-min.png",
	10: "./assets/images/business_card/manNor-11-min.png",
	11: "./assets/images/business_card/manNor-12-min.png",
	12: "./assets/images/business_card/manNor-13-min.png",
	13: "./assets/images/business_card/manNor-14-min.png",
	14: "./assets/images/business_card/manNor-15-min.png",
	15: "./assets/images/business_card/manNor-16-min.png",
	16: "./assets/images/business_card/manNor-17-min.png",
	17: "./assets/images/business_card/manNor-18-min.png",
	18: "./assets/images/business_card/manNor-19-min.png",
	19: "./assets/images/business_card/manNor-20-min.png",
	20: "./assets/images/business_card/manNor-21-min.png",
	21: "./assets/images/business_card/manNor-22-min.png",
	22: "./assets/images/business_card/manNor-23-min.png",
	23: "./assets/images/business_card/manNor-24-min.png",
	24: "./assets/images/business_card/manNor-25-min.png",
	25: "./assets/images/business_card/manNor-26-min.png",
}

var girlCharacterMap = map[int]string{
	0:  "./assets/images/business_card/girlNor-1-min.png",
	1:  "./assets/images/business_card/girlNor-2-min.png",
	2:  "./assets/images/business_card/girlNor-3-min.png",
	3:  "./assets/images/business_card/girlNor-4-min.png",
	4:  "./assets/images/business_card/girlNor-5-min.png",
	5:  "./assets/images/business_card/girlNor-6-min.png",
	6:  "./assets/images/business_card/girlNor-7-min.png",
	7:  "./assets/images/business_card/girlNor-8-min.png",
	8:  "./assets/images/business_card/girlNor-9-min.png",
	9:  "./assets/images/business_card/girlNor-10-min.png",
	10: "./assets/images/business_card/girlNor-11-min.png",
	11: "./assets/images/business_card/girlNor-12-min.png",
	12: "./assets/images/business_card/girlNor-13-min.png",
	13: "./assets/images/business_card/girlNor-14-min.png",
	14: "./assets/images/business_card/girlNor-15-min.png",
	15: "./assets/images/business_card/girlNor-16-min.png",
	16: "./assets/images/business_card/girlNor-17-min.png",
	17: "./assets/images/business_card/girlNor-18-min.png",
	18: "./assets/images/business_card/girlNor-19-min.png",
	19: "./assets/images/business_card/girlNor-20-min.png",
	20: "./assets/images/business_card/girlNor-21-min.png",
	21: "./assets/images/business_card/girlNor-22-min.png",
	22: "./assets/images/business_card/girlNor-23-min.png",
	23: "./assets/images/business_card/girlNor-24-min.png",
	24: "./assets/images/business_card/girlNor-25-min.png",
}

var manSampleMap = map[int]string{
	0:      "./assets/images/bg/bg_card1-min.png",
	1:      "./assets/images/bg/bg_card1-min.png",
	2:      "./assets/images/bg/bg_card2-min.png",
	3:      "./assets/images/bg/bg_card3-min.png",
	4:      "./assets/images/bg/bg_card4-min.png",
	5:      "./assets/images/bg/bg_card5-min.png",
	6:      "./assets/images/bg/bg_card6-min.png",
	7:      "./assets/images/bg/bg_card7-min.png",
	8:      "./assets/images/bg/bg_card8-min.png",
	9:      "./assets/images/bg/bg_card9-min.png",
	10:     "./assets/images/bg/bg_card10-min.png",
	100001: "./assets/images/bg/dcbk-min.png",
	11:     "./assets/images/bg/dcbk-min.png",
}

var girlSampleMap = map[int]string{
	0:      "./assets/images/bg/bg_card1-min.png",
	1:      "./assets/images/bg/bg_card1-min.png",
	2:      "./assets/images/bg/bg_card2-min.png",
	3:      "./assets/images/bg/bg_card3-min.png",
	4:      "./assets/images/bg/bg_card4-min.png",
	5:      "./assets/images/bg/bg_card5-min.png",
	6:      "./assets/images/bg/bg_card6-min.png",
	7:      "./assets/images/bg/bg_card7-min.png",
	8:      "./assets/images/bg/bg_card8-min.png",
	9:      "./assets/images/bg/bg_card9-min.png",
	10:     "./assets/images/bg/bg_card10-min.png",
	100001: "./assets/images/bg/dcbk-min.png",
	11:     "./assets/images/bg/dcbk-min.png",
}

func BusinessCard(sex, avatar, sample int, name, company, position string) (*bytes.Buffer, error) {

	const width = 500
	const height = 400

	img1File, err := os.Open("./assets/images/business_card/bg_share_behind.png")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer func() {
		_ = img1File.Close()
	}()
	img1, err := png.Decode(img1File)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	bounds1 := img1.Bounds()

	samplePath := manSampleMap[sample]
	if sex == 1 {
		samplePath = girlSampleMap[sample]
	}

	sampleFile, err := os.Open(samplePath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer func() {
		_ = sampleFile.Close()
	}()
	imgSample, err := png.Decode(sampleFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	boundsSample := imgSample.Bounds()

	bk2File, err := os.Open("./assets/images/business_card/bg_share_front.png")
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		_ = bk2File.Close()
	}()
	bk2, err := png.Decode(bk2File)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	boundsBk2 := bk2.Bounds()

	path := manCharacterMap[avatar]
	if sex == 1 {
		path = girlCharacterMap[avatar]
	}

	img2File, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer func() {
		_ = img2File.Close()
	}()
	img2, err := png.Decode(img2File)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	bounds2 := img2.Bounds()

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(newImg, bounds1, img1, image.Point{}, draw.Src)

	draw.Draw(
		newImg,
		image.Rect(30, 62, 30+boundsSample.Dx(), 62+boundsSample.Dy()),
		imgSample,
		image.Point{},
		draw.Over,
	)

	centerX := -22 // 600/2 - bounds2.Dx()/2
	centerY := 10  // 220 - bounds2.Dy()/2

	draw.Draw(
		newImg,
		image.Rect(centerX, centerY, centerX+bounds2.Dx(), centerY+bounds2.Dy()),
		img2,
		image.Point{},
		draw.Over,
	)

	draw.Draw(
		newImg,
		image.Rect(0, 0, boundsBk2.Dx(), boundsBk2.Dy()),
		bk2,
		image.Point{},
		draw.Over,
	)

	touMinFile, err := os.Open("./assets/images/toumin-min.png")
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		_ = touMinFile.Close()
	}()
	touMin, err := png.Decode(touMinFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	boundsTouMin := touMin.Bounds()

	col := color.RGBA{0, 0, 0, 255}

	draw.Draw(
		newImg,
		image.Rect(284, 210, 284+boundsTouMin.Dx(), 210+boundsTouMin.Dy()),
		touMin,
		image.Point{},
		draw.Over,
	)

	fontFile, err := os.Open("./assets/SourceHanSansCN-Regular.otf")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer func() {
		_ = fontFile.Close()
	}()

	fontData, err := io.ReadAll(fontFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	f, err := opentype.Parse(fontData)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    28,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	faceTextWidth := font.MeasureString(face, name).Ceil()

	faceStartX := newImg.Bounds().Dx() - faceTextWidth - 56

	point := fixed.Point26_6{fixed.Int26_6(faceStartX * 64), fixed.Int26_6(150 * 64)}
	d := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(name)

	gs, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    21, // 设置字体大小为 24
		DPI:     52,
		Hinting: font.HintingFull,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	gsTextWidth := font.MeasureString(gs, company).Ceil()

	gsStartX := newImg.Bounds().Dx() - gsTextWidth - 56

	pointGs := fixed.Point26_6{fixed.Int26_6(gsStartX * 64), fixed.Int26_6(206 * 64)}
	dGs := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(col),
		Face: gs,
		Dot:  pointGs,
	}
	dGs.DrawString(company)

	fzr, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    21,
		DPI:     52,
		Hinting: font.HintingFull,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	fzrTextWidth := font.MeasureString(fzr, position).Ceil()

	fzrStartX := newImg.Bounds().Dx() - fzrTextWidth - 56

	pointFzr := fixed.Point26_6{fixed.Int26_6(fzrStartX * 64), fixed.Int26_6(180 * 64)}
	dFzr := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(col),
		Face: fzr,
		Dot:  pointFzr,
	}
	dFzr.DrawString(position)

	//yx, err := opentype.NewFace(f, &opentype.FaceOptions{
	//	Size:    14,
	//	DPI:     36,
	//	Hinting: font.HintingFull,
	//})
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}
	//yxTextWidth := font.MeasureString(yx, "123456789@ecosmos.vip").Ceil()
	//
	//yxStartX := newImg.Bounds().Dx() - yxTextWidth - 80
	//
	//pointYx := fixed.Point26_6{fixed.Int26_6(yxStartX * 64), fixed.Int26_6(196 * 64)}
	//dYx := &font.Drawer{
	//	Dst:  newImg,
	//	Src:  image.NewUniform(col),
	//	Face: yx,
	//	Dot:  pointYx,
	//}
	//dYx.DrawString("123456789@ecomos.vip")

	//sj, err := opentype.NewFace(f, &opentype.FaceOptions{
	//	Size:    14,
	//	DPI:     36,
	//	Hinting: font.HintingFull,
	//})
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}
	//sjTextWidth := font.MeasureString(sj, "12345678911").Ceil()
	//
	//sjStartX := newImg.Bounds().Dx() - sjTextWidth - 80
	//
	//pointSj := fixed.Point26_6{fixed.Int26_6(sjStartX * 64), fixed.Int26_6(182 * 64)}
	//dSj := &font.Drawer{
	//	Dst:  newImg,
	//	Src:  image.NewUniform(col),
	//	Face: sj,
	//	Dot:  pointSj,
	//}
	//dSj.DrawString("12345678911")

	newByte := new(bytes.Buffer)

	err = png.Encode(newByte, newImg)
	if err != nil {
		logger.Error(err)
	}

	return newByte, err
}

func ResizeImage(inputPath string, outputPath string, width int, height int) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		logger.Error(err)
		return
	}

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Src)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		_ = outputFile.Close()
	}()

	err = png.Encode(outputFile, newImg)
	if err != nil {
		logger.Error(err)
		return
	}
}

func ScaleDownImage(inputPath, outputPath string, scaleFactor float64) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		_ = inputFile.Close()
	}()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		logger.Error(err)
		return
	}

	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * scaleFactor)
	newHeight := int(float64(bounds.Dy()) * scaleFactor)

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / scaleFactor)
			srcY := int(float64(y) / scaleFactor)
			r, g, b, a := img.At(srcX, srcY).RGBA()
			newImg.Set(x, y,
				color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)},
			)
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, newImg)
	if err != nil {
		logger.Error(err)
		return
	}
}

func ScaleDownImage2(inputPath, outputPath string, scaleFactor float64) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		logger.Error(err)
		return
	}

	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * scaleFactor)
	newHeight := int(float64(bounds.Dy()) * scaleFactor)

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	drawx.ApproxBiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		_ = outputFile.Close()
	}()

	err = png.Encode(outputFile, newImg)
	if err != nil {
		logger.Error(err)
		return
	}
}
