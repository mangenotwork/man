package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// 图片信息结构体
type ImageInfo struct {
	FileName       string            `json:"fileName"`
	FileSize       int64             `json:"fileSize"`
	MimeType       string            `json:"mimeType"`
	Width          int               `json:"width"`
	Height         int               `json:"height"`
	Orientation    int               `json:"orientation"`
	CameraMake     string            `json:"cameraMake"`
	CameraModel    string            `json:"cameraModel"`
	DateTime       string            `json:"dateTime"`
	Error          string            `json:"error"`
	GPSCoordinates map[string]string `json:"gpsCoordinates"`
}

func main() {
	filePath := "./test4.jpg"

	info, err := GetImageInfo(filePath)
	if err != nil {
		log.Fatalf("获取图片信息失败: %v", err)
	}

	if true {
		outputJSON(info)
	} else {
		outputText(info)
	}
}

// 获取图片信息
func GetImageInfo(filePath string) (*ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	info := &ImageInfo{
		FileName: filepath.Base(filePath),
		FileSize: fileInfo.Size(),
	}

	// 读取文件前512字节来确定MIME类型
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	info.MimeType = http.DetectContentType(buf[:n])

	// 将文件指针重置到开始位置
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("重置文件指针失败: %w", err)
	}

	// 确定图片格式并获取尺寸
	switch info.MimeType {
	case "image/jpeg":
		if err := extractJPEGInfo(file, info); err != nil {
			info.Error = err.Error()
		}
	case "image/png":
		if err := extractPNGInfo(file, info); err != nil {
			info.Error = err.Error()
		}
	default:
		info.Error = fmt.Sprintf("不支持的图片格式: %s", info.MimeType)
	}

	return info, nil
}

// 提取JPEG图片信息
func extractJPEGInfo(file *os.File, info *ImageInfo) error {
	// 重置文件指针
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("重置文件指针失败: %w", err)
	}

	// 手动解析JPEG文件结构
	width, height, err := parseJPEGDimensions(file)
	if err != nil {
		return fmt.Errorf("解析JPEG尺寸失败: %w", err)
	}
	info.Width = width
	info.Height = height

	// 重置文件指针以读取EXIF
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("重置文件指针失败: %w", err)
	}

	// 手动解析JPEG文件以提取EXIF信息
	if err := extractJPEGEXIFInfo(file, info); err != nil {
		return fmt.Errorf("提取EXIF信息失败: %w", err)
	}

	return nil
}

// 解析JPEG尺寸
func parseJPEGDimensions(file *os.File) (int, int, error) {
	// 读取JPEG文件头
	buf := make([]byte, 2)
	_, err := io.ReadFull(file, buf)
	if err != nil {
		return 0, 0, fmt.Errorf("读取JPEG文件头失败: %w", err)
	}

	// 检查JPEG文件签名
	if buf[0] != 0xFF || buf[1] != 0xD8 {
		return 0, 0, fmt.Errorf("不是有效的JPEG文件")
	}

	// 查找SOF标记（Start of Frame）
	for {
		// 读取标记
		_, err := io.ReadFull(file, buf)
		if err != nil {
			return 0, 0, fmt.Errorf("读取JPEG标记失败: %w", err)
		}

		// 检查标记是否为SOI（文件开始）
		if buf[0] != 0xFF {
			return 0, 0, fmt.Errorf("无效的JPEG标记")
		}

		// 检查是否为SOF标记（0xC0-0xC3, 0xC5-0xC7, 0xC9-0xCB, 0xCD-0xCF）
		if (buf[1] >= 0xC0 && buf[1] <= 0xC3) ||
			(buf[1] >= 0xC5 && buf[1] <= 0xC7) ||
			(buf[1] >= 0xC9 && buf[1] <= 0xCB) ||
			(buf[1] >= 0xCD && buf[1] <= 0xCF) {

			// 读取长度（2字节）
			lengthBytes := make([]byte, 2)
			_, err := io.ReadFull(file, lengthBytes)
			if err != nil {
				return 0, 0, fmt.Errorf("读取SOF长度失败: %w", err)
			}

			// 跳过数据精度字节
			_, err = file.Seek(1, io.SeekCurrent)
			if err != nil {
				return 0, 0, fmt.Errorf("跳过数据精度字节失败: %w", err)
			}

			// 读取高度和宽度
			heightBytes := make([]byte, 2)
			widthBytes := make([]byte, 2)

			_, err = io.ReadFull(file, heightBytes)
			if err != nil {
				return 0, 0, fmt.Errorf("读取高度失败: %w", err)
			}

			_, err = io.ReadFull(file, widthBytes)
			if err != nil {
				return 0, 0, fmt.Errorf("读取宽度失败: %w", err)
			}

			height := int(binary.BigEndian.Uint16(heightBytes))
			width := int(binary.BigEndian.Uint16(widthBytes))

			return width, height, nil
		} else {
			// 不是SOF标记，跳过此段
			lengthBytes := make([]byte, 2)
			_, err := io.ReadFull(file, lengthBytes)
			if err != nil {
				return 0, 0, fmt.Errorf("读取标记长度失败: %w", err)
			}
			segmentLength := int(binary.BigEndian.Uint16(lengthBytes)) - 2
			_, err = file.Seek(int64(segmentLength), io.SeekCurrent)
			if err != nil {
				return 0, 0, fmt.Errorf("跳过段失败: %w", err)
			}
		}
	}
}

// 提取JPEG中的EXIF信息（不使用第三方库）
func extractJPEGEXIFInfo(file *os.File, info *ImageInfo) error {
	// 读取JPEG文件头
	buf := make([]byte, 2)
	_, err := io.ReadFull(file, buf)
	if err != nil {
		return fmt.Errorf("读取JPEG文件头失败: %w", err)
	}

	// 检查JPEG文件签名
	if buf[0] != 0xFF || buf[1] != 0xD8 {
		return fmt.Errorf("不是有效的JPEG文件")
	}

	// 查找EXIF数据
	for {
		// 读取标记
		_, err := io.ReadFull(file, buf)
		if err != nil {
			if err == io.EOF {
				return nil // 没有找到EXIF
			}
			return fmt.Errorf("读取JPEG标记失败: %w", err)
		}

		// 检查标记是否为SOI（文件开始）
		if buf[0] != 0xFF {
			return fmt.Errorf("无效的JPEG标记")
		}

		// 检查是否为APP1标记（包含EXIF数据）
		if buf[1] == 0xE1 {
			// 读取长度（2字节）
			lengthBytes := make([]byte, 2)
			_, err := io.ReadFull(file, lengthBytes)
			if err != nil {
				return fmt.Errorf("读取APP1长度失败: %w", err)
			}
			segmentLength := int(binary.BigEndian.Uint16(lengthBytes)) - 2

			// 读取APP1段内容
			app1Data := make([]byte, segmentLength)
			_, err = io.ReadFull(file, app1Data)
			if err != nil {
				return fmt.Errorf("读取APP1段失败: %w", err)
			}

			// 检查是否为EXIF数据
			if len(app1Data) >= 6 && string(app1Data[:6]) == "Exif\x00\x00" {
				// 跳过"Exif\x00\x00"
				exifData := app1Data[6:]
				fmt.Printf("找到EXIF数据，长度: %d 字节\n", len(exifData)) // 调试输出
				if err := parseEXIF(exifData, info); err != nil {
					return fmt.Errorf("解析EXIF数据失败: %w", err)
				}
				return nil
			}
		} else {
			// 不是APP1标记，跳过此段
			lengthBytes := make([]byte, 2)
			_, err := io.ReadFull(file, lengthBytes)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("读取标记长度失败: %w", err)
			}
			segmentLength := int(binary.BigEndian.Uint16(lengthBytes)) - 2
			_, err = file.Seek(int64(segmentLength), io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("跳过段失败: %w", err)
			}
		}
	}
}

// 解析EXIF数据
func parseEXIF(data []byte, info *ImageInfo) error {
	// 检查字节序
	if len(data) < 8 {
		return fmt.Errorf("EXIF数据太短")
	}

	var byteOrder binary.ByteOrder
	if string(data[0:2]) == "II" {
		byteOrder = binary.LittleEndian
		fmt.Println("使用小端字节序") // 调试输出
	} else if string(data[0:2]) == "MM" {
		byteOrder = binary.BigEndian
		fmt.Println("使用大端字节序") // 调试输出
	} else {
		return fmt.Errorf("未知的字节序: %q", string(data[0:2]))
	}

	// 检查TIFF标记
	if byteOrder.Uint16(data[2:4]) != 0x002A {
		return fmt.Errorf("无效的TIFF标记: 0x%04X", byteOrder.Uint16(data[2:4]))
	}

	// 获取第一个IFD的偏移量
	ifdOffset := int(byteOrder.Uint32(data[4:8]))
	if ifdOffset >= len(data) {
		return fmt.Errorf("IFD偏移量超出数据范围: %d >= %d", ifdOffset, len(data))
	}

	fmt.Printf("主IFD偏移量: 0x%X\n", ifdOffset) // 调试输出

	// 解析主IFD
	err := parseIFD(data, ifdOffset, byteOrder, info, true)
	if err != nil {
		return err
	}

	// 查找EXIF子IFD
	exifIfdOffset := findEXIFSubIFD(data, ifdOffset, byteOrder)
	if exifIfdOffset > 0 && exifIfdOffset < len(data) {
		fmt.Printf("找到EXIF子IFD，偏移量: 0x%X\n", exifIfdOffset) // 调试输出
		err := parseEXIFSubIFD(data, exifIfdOffset, byteOrder, info)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("未找到EXIF子IFD") // 调试输出
	}

	return nil
}

// 解析IFD
func parseIFD(data []byte, ifdOffset int, byteOrder binary.ByteOrder, info *ImageInfo, isMain bool) error {
	// 读取条目数量
	entryCount := int(byteOrder.Uint16(data[ifdOffset : ifdOffset+2]))
	fmt.Printf("IFD条目数量: %d\n", entryCount) // 调试输出

	// 遍历每个条目
	for i := 0; i < entryCount; i++ {
		entryOffset := ifdOffset + 2 + i*12
		if entryOffset+12 > len(data) {
			fmt.Printf("条目 %d 超出数据范围，跳过\n", i) // 调试输出
			break
		}

		tag := byteOrder.Uint16(data[entryOffset : entryOffset+2])
		typ := byteOrder.Uint16(data[entryOffset+2 : entryOffset+4])
		count := byteOrder.Uint32(data[entryOffset+4 : entryOffset+8])

		// 计算值的字节数
		var typeSize int
		switch typ {
		case 1: // BYTE
			typeSize = 1
		case 2: // ASCII
			typeSize = 1
		case 3: // SHORT
			typeSize = 2
		case 4: // LONG
			typeSize = 4
		case 5: // RATIONAL
			typeSize = 8
		default:
			continue // 跳过不支持的类型
		}

		valueBytes := int(count) * typeSize
		valueOffset := entryOffset + 8

		// 如果值超过4字节，则从偏移量处读取
		if valueBytes > 4 {
			valueOffset = int(byteOrder.Uint32(data[entryOffset+8 : entryOffset+12]))
		}

		// 调试输出重要标签
		if tag == 0x010F || tag == 0x0110 || tag == 0x0112 || tag == 0x8769 || (isMain && tag == 0x9003) {
			tagName := "Unknown"
			switch tag {
			case 0x010F:
				tagName = "Make"
			case 0x0110:
				tagName = "Model"
			case 0x0112:
				tagName = "Orientation"
			case 0x8769:
				tagName = "ExifSubIFD"
			case 0x9003:
				tagName = "DateTimeOriginal"
			}
			fmt.Printf("找到标签: %s (0x%04X), 类型: %d, 计数: %d, 值偏移: 0x%X\n",
				tagName, tag, typ, count, valueOffset)
		}

		// 处理感兴趣的标签
		switch tag {
		case 0x010F: // Make
			if typ == 2 { // ASCII
				makeStr := readASCIIString(data, valueOffset, int(count))
				info.CameraMake = makeStr
				fmt.Printf("相机制造商: %q\n", makeStr) // 调试输出
			}
		case 0x0110: // Model
			if typ == 2 { // ASCII
				modelStr := readASCIIString(data, valueOffset, int(count))
				info.CameraModel = modelStr
				fmt.Printf("相机型号: %q\n", modelStr) // 调试输出
			}
		case 0x0112: // Orientation
			if typ == 3 { // SHORT
				orientation := int(byteOrder.Uint16(data[valueOffset : valueOffset+2]))
				info.Orientation = orientation
				fmt.Printf("图片方向: %d\n", orientation) // 调试输出
			}
		case 0x9003: // DateTimeOriginal (在主IFD中查找)
			if isMain && typ == 2 { // ASCII
				dateTime := readASCIIString(data, valueOffset, int(count))
				info.DateTime = dateTime
				fmt.Printf("主IFD中找到拍摄时间: %q\n", dateTime) // 调试输出
			}
		}
	}

	// 获取下一个IFD的偏移量
	nextIFDOffset := int(byteOrder.Uint32(data[ifdOffset+2+entryCount*12 : ifdOffset+2+entryCount*12+4]))
	if nextIFDOffset > 0 && nextIFDOffset < len(data) {
		fmt.Printf("下一个IFD偏移量: 0x%X\n", nextIFDOffset) // 调试输出
		return parseIFD(data, nextIFDOffset, byteOrder, info, false)
	}

	return nil
}

// 查找EXIF子IFD
func findEXIFSubIFD(data []byte, ifdOffset int, byteOrder binary.ByteOrder) int {
	// 读取条目数量
	entryCount := int(byteOrder.Uint16(data[ifdOffset : ifdOffset+2]))

	// 遍历每个条目查找EXIF子IFD
	for i := 0; i < entryCount; i++ {
		entryOffset := ifdOffset + 2 + i*12
		if entryOffset+12 > len(data) {
			continue
		}

		tag := byteOrder.Uint16(data[entryOffset : entryOffset+2])
		if tag == 0x8769 { // ExifSubIFD
			return int(byteOrder.Uint32(data[entryOffset+8 : entryOffset+12]))
		}
	}

	return 0
}

// 解析EXIF子IFD
func parseEXIFSubIFD(data []byte, ifdOffset int, byteOrder binary.ByteOrder, info *ImageInfo) error {
	// 读取条目数量
	entryCount := int(byteOrder.Uint16(data[ifdOffset : ifdOffset+2]))
	fmt.Printf("EXIF子IFD条目数量: %d\n", entryCount) // 调试输出

	// 遍历每个条目
	for i := 0; i < entryCount; i++ {
		entryOffset := ifdOffset + 2 + i*12
		if entryOffset+12 > len(data) {
			fmt.Printf("EXIF子IFD条目 %d 超出数据范围，跳过\n", i) // 调试输出
			continue                                   // 修复：使用continue而不是break，继续处理后续条目
		}

		tag := byteOrder.Uint16(data[entryOffset : entryOffset+2])
		typ := byteOrder.Uint16(data[entryOffset+2 : entryOffset+4])
		count := byteOrder.Uint32(data[entryOffset+4 : entryOffset+8])

		log.Printf("tag %x", tag)

		// 计算值的字节数
		var typeSize int
		switch typ {
		case 1: // BYTE
			typeSize = 1
		case 2: // ASCII
			typeSize = 1
		case 3: // SHORT
			typeSize = 2
		case 4: // LONG
			typeSize = 4
		case 5: // RATIONAL
			typeSize = 8
		default:
			continue // 跳过不支持的类型
		}

		valueBytes := int(count) * typeSize
		valueOffset := entryOffset + 8

		// 如果值超过4字节，则从偏移量处读取
		if valueBytes > 4 {
			valueOffset = int(byteOrder.Uint32(data[entryOffset+8 : entryOffset+12]))
		}

		// 调试输出重要标签
		if tag == 0x9003 {
			fmt.Printf("找到标签: DateTimeOriginal (0x%04X), 类型: %d, 计数: %d, 值偏移: 0x%X\n",
				tag, typ, count, valueOffset)
		}

		// 处理感兴趣的标签
		switch tag {
		case 0x9003: // DateTimeOriginal
			if typ == 2 { // ASCII
				dateTime := readASCIIString(data, valueOffset, int(count))
				if dateTime != "" { // 只在有值时更新
					info.DateTime = dateTime
					fmt.Printf("EXIF子IFD中找到拍摄时间: %q\n", dateTime) // 调试输出
				}
			}
		}
	}

	return nil
}

// 读取ASCII字符串
func readASCIIString(data []byte, offset, length int) string {
	if offset < 0 || offset >= len(data) {
		fmt.Printf("警告: 字符串偏移超出数据范围: %d >= %d\n", offset, len(data))
		return ""
	}

	if offset+length > len(data) {
		length = len(data) - offset
		fmt.Printf("警告: 字符串长度截断为 %d\n", length)
	}

	// 查找字符串结束符
	end := offset + length
	for i := offset; i < end; i++ {
		if data[i] == 0 {
			return string(data[offset:i])
		}
	}

	// 如果没有找到结束符，使用整个长度
	return string(data[offset:end])
}

// 提取PNG图片信息
func extractPNGInfo(file *os.File, info *ImageInfo) error {
	// 重置文件指针
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("重置文件指针失败: %w", err)
	}

	// 读取PNG文件头
	buf := make([]byte, 8)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return fmt.Errorf("读取PNG文件头失败: %w", err)
	}

	// 检查PNG文件签名
	if !isValidPNGSignature(buf) {
		return fmt.Errorf("不是有效的PNG文件")
	}

	// 查找IHDR块（包含尺寸信息）
	for {
		// 读取块长度（4字节）
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(file, lengthBytes)
		if err != nil {
			return fmt.Errorf("读取PNG块长度失败: %w", err)
		}
		blockLength := int(binary.BigEndian.Uint32(lengthBytes))

		// 读取块类型（4字节）
		blockType := make([]byte, 4)
		_, err = io.ReadFull(file, blockType)
		if err != nil {
			return fmt.Errorf("读取PNG块类型失败: %w", err)
		}

		// 检查是否为IHDR块
		if string(blockType) == "IHDR" {
			// 读取宽度和高度
			dimensions := make([]byte, 8)
			_, err := io.ReadFull(file, dimensions)
			if err != nil {
				return fmt.Errorf("读取PNG尺寸失败: %w", err)
			}

			width := int(binary.BigEndian.Uint32(dimensions[:4]))
			height := int(binary.BigEndian.Uint32(dimensions[4:]))

			info.Width = width
			info.Height = height

			// 跳过IHDR块的剩余部分
			_, err = file.Seek(int64(blockLength-8), io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("跳过IHDR块剩余部分失败: %w", err)
			}

			// 读取CRC校验（4字节）
			_, err = file.Seek(4, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("跳过CRC校验失败: %w", err)
			}

			break
		} else {
			// 不是IHDR块，跳过整个块
			_, err = file.Seek(int64(blockLength+4), io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("跳过PNG块失败: %w", err)
			}
		}
	}

	// 注意：PNG格式通常不包含像相机型号、拍摄时间这样的EXIF元数据
	// 它可能包含一些文本元数据（tEXt、iTXt或zTXt块），但解析这些需要更多工作

	return nil
}

// 检查PNG文件签名
func isValidPNGSignature(buf []byte) bool {
	expected := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i := 0; i < 8; i++ {
		if buf[i] != expected[i] {
			return false
		}
	}
	return true
}

// 以JSON格式输出
func outputJSON(info *ImageInfo) {
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatalf("转换为JSON失败: %v", err)
	}
	fmt.Println(string(jsonData))
}

// 以文本格式输出
func outputText(info *ImageInfo) {
	fmt.Printf("文件名: %s\n", info.FileName)
	fmt.Printf("文件大小: %.2f KB\n", float64(info.FileSize)/1024)
	fmt.Printf("MIME类型: %s\n", info.MimeType)
	fmt.Println()

	if info.Width > 0 && info.Height > 0 {
		fmt.Printf("图片尺寸: %dx%d 像素\n", info.Width, info.Height)
	}

	if info.Orientation > 0 {
		fmt.Printf("图片方向: %d\n", info.Orientation)
	}

	if info.CameraMake != "" {
		fmt.Printf("相机制造商: %s\n", info.CameraMake)
	}

	if info.CameraModel != "" {
		fmt.Printf("相机型号: %s\n", info.CameraModel)
	}

	if info.DateTime != "" {
		fmt.Printf("拍摄时间: %s\n", info.DateTime)
	}

	if info.Error != "" {
		fmt.Printf("\n警告: %s\n", info.Error)
	}
}
