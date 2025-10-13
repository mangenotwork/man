package http_kuai_multipart

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// 创建 multipart 表单数据
func createMultipartFormData(fieldName, fileName string, fileContent []byte) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 添加普通表单字段
	if err := w.WriteField("name", "测试文件"); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("description", "这是一个用于测试的文件"); err != nil {
		return nil, "", err
	}

	// 创建文件字段
	fileField, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, "", err
	}

	// 将文件内容写入字段
	if _, err := fileField.Write(fileContent); err != nil {
		return nil, "", err
	}

	// 关闭 multipart writer
	if err := w.Close(); err != nil {
		return nil, "", err
	}

	// 返回缓冲区、boundary
	return &b, w.Boundary(), nil
}

func main() {
	// 读取要上传的文件
	fileContent, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 创建 multipart 表单数据
	buffer, boundary, err := createMultipartFormData("file", "test.txt", fileContent)
	if err != nil {
		fmt.Printf("创建表单数据失败: %v\n", err)
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", "http://localhost:8080/upload", buffer)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Content-Type，包含 boundary
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	fmt.Printf("服务器响应: %s\n", respBody)
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
}
