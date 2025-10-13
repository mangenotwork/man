package http_kuai_TransferEncoding

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 发送请求到分块传输服务器
	resp, err := http.Get("http://localhost:8080/chunked")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应是否使用分块传输
	if resp.TransferEncoding != nil && len(resp.TransferEncoding) > 0 && resp.TransferEncoding[0] == "chunked" {
		fmt.Println("服务器使用分块传输编码")
	} else {
		fmt.Println("服务器未使用分块传输编码")
	}

	// 读取分块数据
	reader := bufio.NewReader(resp.Body)
	for {
		// 读取每个块
		chunk, err := io.ReadAll(reader)
		if err != nil {
			if err == io.EOF {
				// 所有块读取完成
				break
			}
			fmt.Printf("读取数据块失败: %v\n", err)
			return
		}

		if len(chunk) > 0 {
			fmt.Printf("收到数据块: %s", chunk)
		}
	}

	fmt.Println("所有数据块接收完成")
}
