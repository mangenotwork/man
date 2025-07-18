package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

/*
// case1 并没能实现
import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

// LiveMetadata 存储实况图的动态信息
type LiveMetadata struct {
	Version    string   `json:"version"`    // 元数据版本
	Duration   float64  `json:"duration"`   // 动态时长(秒)
	FPS        int      `json:"fps"`        // 帧率
	FrameCount int      `json:"frameCount"` // 总帧数
	FrameURL   string   `json:"frameUrl"`   // 帧序列URL模板
	Trigger    string   `json:"trigger"`    // 触发方式(如"long_press")
	Thumbnails []string `json:"thumbnails"` // 缩略图URL
}

// 创建可正常打开的实况图（严格遵循JPEG格式规范）
func createLiveJPEG(inputImagePath, outputPath string, meta LiveMetadata) error {
	// 1. 读取并解码原始图片（获取image对象）
	inputFile, err := os.Open(inputImagePath)
	if err != nil {
		return fmt.Errorf("无法打开输入图片: %v", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("输入图片解码失败（请确保是 valid JPEG）: %v", err)
	}

	// 2. 生成动态元数据（JSON格式）
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("元数据序列化失败: %v", err)
	}

	// 3. 构建JPEG的APP1段（存储元数据）
	metaPrefix := []byte("LiveMeta:") // 自定义标记，用于后续提取
	app1Data := append(metaPrefix, metaJSON...)
	app1Length := len(app1Data) + 2 // +2是因为长度字段本身占2字节
	if app1Length > 65535 {
		return fmt.Errorf("元数据过大（超过JPEG APP1段最大限制）")
	}

	// 4. 生成完整的JPEG文件内容
	var jpegBuffer bytes.Buffer

	// 4.1 先写入JPEG起始标记（SOI: 0xFFD8）
	jpegBuffer.Write([]byte{0xFF, 0xD8})

	// 4.2 写入APP1段（包含我们的元数据）
	jpegBuffer.Write([]byte{0xFF, 0xE1})                                     // APP1标记
	jpegBuffer.Write([]byte{byte(app1Length >> 8), byte(app1Length & 0xFF)}) // 长度（大端序）
	jpegBuffer.Write(app1Data)

	// 4.3 写入图片数据（使用临时缓冲区处理编码结果）
	var imgBuffer bytes.Buffer
	if err := jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 90}); err != nil {
		return fmt.Errorf("图片编码失败: %v", err)
	}

	// 处理编码结果，移除自动添加的SOI/EOI标记
	imgData := imgBuffer.Bytes()
	if len(imgData) >= 2 && imgData[0] == 0xFF && imgData[1] == 0xD8 {
		imgData = imgData[2:] // 移除SOI
	}
	if len(imgData) >= 2 && imgData[len(imgData)-2] == 0xFF && imgData[len(imgData)-1] == 0xD9 {
		imgData = imgData[:len(imgData)-2] // 移除EOI
	}

	// 将处理后的图片数据写入最终缓冲区
	if _, err := jpegBuffer.Write(imgData); err != nil {
		return fmt.Errorf("写入图片数据失败: %v", err)
	}

	// 4.4 写入JPEG结束标记（EOI）
	jpegBuffer.Write([]byte{0xFF, 0xD9})

	// 5. 验证生成的JPEG是否有效
	_, format, err := image.Decode(bytes.NewReader(jpegBuffer.Bytes()))
	if err != nil || format != "jpeg" {
		return fmt.Errorf("生成的文件不是有效的JPEG格式: %v", err)
	}

	// 6. 写入输出文件
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("无法创建输出目录: %v", err)
	}

	if err := os.WriteFile(outputPath, jpegBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// extractLiveMetadata 从JPEG中提取动态元数据
func extractLiveMetadata(jpegPath string) (LiveMetadata, error) {
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return LiveMetadata{}, fmt.Errorf("无法读取文件: %v", err)
	}

	metaPrefix := []byte("LiveMeta:")
	metaStart := bytes.Index(data, metaPrefix)
	if metaStart == -1 {
		return LiveMetadata{}, fmt.Errorf("未找到元数据")
	}

	// 从APP1段中提取元数据（跳过前缀）
	metaStart += len(metaPrefix)
	metaEnd := metaStart
	for metaEnd < len(data) && !(data[metaEnd] == 0xFF && metaEnd+1 < len(data) && data[metaEnd+1] != 0x00) {
		metaEnd++
	}

	var meta LiveMetadata
	if err := json.Unmarshal(data[metaStart:metaEnd], &meta); err != nil {
		return LiveMetadata{}, fmt.Errorf("解析元数据失败: %v", err)
	}

	return meta, nil
}

func main() {
	outputPath := "./out.jpg"
	inputPath := "./1.jpg"
	frameURL := "http://10.0.40.203:17777/images/{frame}.jpg"

	// 构建元数据
	meta := LiveMetadata{
		Version:    "1.0",
		Duration:   3,  // 动态时长(秒)
		FPS:        30, // 帧率
		FrameCount: 90, // 总帧数
		FrameURL:   frameURL,
		Trigger:    "long_press",
		Thumbnails: []string{
			"http://10.0.40.203:17777/images/2.jpg",
		},
	}

	// 创建实况图
	if err := createLiveJPEG(inputPath, outputPath, meta); err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}

	fmt.Printf("成功创建实况图: %s\n", outputPath)
	fmt.Println("可以使用 -extract -input [文件路径] 提取元数据进行验证")

	// 提取元数据模式
	meta2, err := extractLiveMetadata(outputPath)
	if err != nil {
		fmt.Printf("提取失败: %v\n", err)
		return
	}

	fmt.Println("提取到的实况图元数据:")
	jsonMeta, _ := json.MarshalIndent(meta2, "", "  ")
	fmt.Println(string(jsonMeta))
	return

}
*/

// case2 也没有实现效果

// 苹果Live Photo所需的关键元数据标记
const (
	// EXIF APP1段标记
	app1Marker = 0xFFE1
	// 苹果私有元数据标签
	appleLivePhotoVersionTag = 0xA005 // 版本标记
	livePhotoUUIDTag         = 0xA006 // UUID标记
)

// 生成符合Live Photo规范的JPEG
func createLivePhotoJPEG(inputPath, outputPath, uuidStr string) error {
	// 1. 读取原始图片
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("无法解码图片: %v", err)
	}

	// 2. 处理UUID（转换为16字节）
	uuidBytes, err := parseUUID(uuidStr)
	if err != nil {
		return fmt.Errorf("UUID格式错误: %v", err)
	}

	// 3. 构建苹果私有EXIF数据
	exifData := buildAppleExifData(uuidBytes)

	// 4. 构建完整的JPEG数据
	jpegData, err := buildJPEGWithExif(img, exifData)
	if err != nil {
		return fmt.Errorf("构建JPEG失败: %v", err)
	}

	// 5. 写入输出文件
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	if err := os.WriteFile(outputPath, jpegData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 解析UUID字符串为16字节
func parseUUID(uuidStr string) ([]byte, error) {
	// 移除UUID中的短横线
	cleaned := make([]rune, 0, 32)
	for _, c := range uuidStr {
		if c != '-' {
			cleaned = append(cleaned, c)
		}
	}

	if len(cleaned) != 32 {
		return nil, fmt.Errorf("UUID格式应为标准8-4-4-4-12格式，如: 123e4567-e89b-12d3-a456-426614174000")
	}

	// 转换为16字节
	uuidBytes := make([]byte, 16)
	for i := 0; i < 16; i++ {
		hexStr := string(cleaned[i*2 : (i+1)*2])
		val, err := hex.DecodeString(hexStr)
		if err != nil || len(val) != 1 {
			return nil, fmt.Errorf("UUID格式错误，位置 %d", i)
		}
		uuidBytes[i] = val[0]
	}

	return uuidBytes, nil
}

// 构建包含苹果私有元数据的EXIF数据
func buildAppleExifData(uuid []byte) []byte {
	var exifBuffer bytes.Buffer

	// TIFF头部 (小端序)
	exifBuffer.Write([]byte{0x49, 0x49, 0x2A, 0x00})          // II* 标记，小端序
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(8)) // 偏移量

	// IFD0 结构
	// 条目数量 (2个: 版本和UUID)
	binary.Write(&exifBuffer, binary.LittleEndian, uint16(2))

	// 第一个条目: AppleLivePhotoVersion (0xA005)
	binary.Write(&exifBuffer, binary.LittleEndian, uint16(appleLivePhotoVersionTag)) // 标签
	binary.Write(&exifBuffer, binary.LittleEndian, uint16(4))                        // 类型: 无符号长整型
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(1))                        // 数量
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(1))                        // 值: 1 (版本号)

	// 第二个条目: LivePhotoUUID (0xA006)
	binary.Write(&exifBuffer, binary.LittleEndian, uint16(livePhotoUUIDTag)) // 标签
	binary.Write(&exifBuffer, binary.LittleEndian, uint16(7))                // 类型: 未定义
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(16))               // 数量: 16字节
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(0x1C))             // 值偏移量

	// IFD0 结束标记
	binary.Write(&exifBuffer, binary.LittleEndian, uint32(0))

	// UUID值 (16字节)
	exifBuffer.Write(uuid)

	return exifBuffer.Bytes()
}

// 构建包含EXIF数据的完整JPEG
func buildJPEGWithExif(img image.Image, exifData []byte) ([]byte, error) {
	var jpegBuffer bytes.Buffer

	// 1. 写入JPEG起始标记 (SOI)
	jpegBuffer.Write([]byte{0xFF, 0xD8})

	// 2. 写入EXIF APP1段
	exifApp1Data := append([]byte("Exif\000\000"), exifData...) // EXIF前缀
	app1Length := len(exifApp1Data) + 2                         // 长度包含自身2字节

	// 写入APP1标记和长度
	jpegBuffer.Write([]byte{byte(app1Marker >> 8), byte(app1Marker & 0xFF)})
	jpegBuffer.Write([]byte{byte(app1Length >> 8), byte(app1Length & 0xFF)})
	jpegBuffer.Write(exifApp1Data)

	// 3. 写入图片数据
	var imgBuffer bytes.Buffer
	if err := jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}

	// 移除图片数据中可能包含的SOI标记
	imgBytes := imgBuffer.Bytes()
	if len(imgBytes) >= 2 && imgBytes[0] == 0xFF && imgBytes[1] == 0xD8 {
		imgBytes = imgBytes[2:]
	}

	jpegBuffer.Write(imgBytes)

	return jpegBuffer.Bytes(), nil
}

func main() {
	//inputPath := flag.String("input", "", "输入JPEG图片路径")
	//outputPath := flag.String("output", "live_photo.jpg", "输出JPEG路径")
	//uuid := flag.String("uuid", "", "标准UUID (如: 123e4567-e89b-12d3-a456-426614174000)")
	flag.Parse()
	inputPath := "./1.jpg"
	outputPath := "./123e4567-e89b-12d3-a456-426614174000.jpg"
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	//if *inputPath == "" || *uuid == "" {
	//	fmt.Println("用法: go run live_photo.go -input 原始图片.jpg -uuid 标准UUID")
	//	fmt.Println("示例UUID: 123e4567-e89b-12d3-a456-426614174000")
	//	return
	//}

	if err := createLivePhotoJPEG(inputPath, outputPath, uuid); err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}

	fmt.Printf("成功生成符合Live Photo规范的JPEG: %s\n", outputPath)
	fmt.Println("使用说明:")
	fmt.Println("1. 需创建同名UUID的MOV视频文件(如UUID为xxx, 视频命名为xxx.mov)")
	fmt.Println("2. 将JPEG和MOV一起导入iOS设备相册")
	fmt.Println("3. 长按图片即可触发实况效果")
}
