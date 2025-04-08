package draw

import (
	"testing"
)

func Test_Case(t *testing.T) {
	////BusinessCard(0, 6, 1)
	//inputPath := "../../api/assets/images/girl-10.png"
	//outputPath := "../../api/assets/images/girl-10-min.png"
	////width := 460
	////height := 290
	////ResizeImage(inputPath, outputPath, width, height)
	//scaleFactor := 0.437
	//ScaleDownImage(inputPath, outputPath, scaleFactor)

	//inputPath := "../../api/assets/images/business_card/girlNor-25.png"
	//outputPath := "../../api/assets/images/business_card/girlNor-25-min.png"
	////width := 460
	////height := 290
	//ResizeImage(inputPath, outputPath, width, height)
	//scaleFactor := 0.8
	ScaleDownImage2(inputPath, outputPath, scaleFactor)

	inputPath := "../../api/assets/images/bg/dcbk.png"
	outputPath := "../../api/assets/images/bg/dcbk-min.png"
	scaleFactor := 0.437
	ScaleDownImage2(inputPath, outputPath, scaleFactor)

	//buf, err := BusinessCard(0, 6, 1, "name", "公司试试哈斯哈斯哈斯哈斯", "aaaaaaaaa")
	//if err != nil {
	//	fmt.Println("生成图片失败:", err)
	//	return
	//}
	//
	//file, err := os.Create("./output.jpg") // 根据实际情况更改文件名和格式
	//if err != nil {
	//	fmt.Println("创建图片文件失败:", err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = io.Copy(file, buf)
	//if err != nil {
	//	fmt.Println("将Buffer数据写入图片文件失败:", err)
	//	return
	//}
	//
	//fmt.Println("图片数据已成功存入图片文件")

}
