package http_kuai_TransferEncoding

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 注册分块传输处理函数
	http.HandleFunc("/chunked", chunkedHandler)

	// 启动服务器
	fmt.Println("服务器启动在 :8080")
	http.ListenAndServe(":8080", nil)
}

// chunkedHandler 实现分块传输
func chunkedHandler(w http.ResponseWriter, r *http.Request) {
	// 确保响应使用分块传输编码
	// 当不设置 Content-Length 时，Go 会自动使用分块传输
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 分5次发送数据块
	for i := 1; i <= 5; i++ {
		// 每个块的格式: [长度]\r\n[数据]\r\n
		// Go 的 ResponseWriter 会自动处理块格式，我们只需要写入数据
		chunk := fmt.Sprintf("这是第 %d 个数据块\n", i)

		// 写入数据块
		_, err := fmt.Fprint(w, chunk)
		if err != nil {
			fmt.Printf("写入数据块失败: %v\n", err)
			return
		}

		// 刷新缓冲区，确保数据立即发送
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		} else {
			fmt.Println("响应不支持刷新，可能无法实时看到分块数据")
		}

		// 等待1秒，模拟数据生成延迟
		time.Sleep(1 * time.Second)
	}

	// 分块传输结束时，Go 会自动发送终止块 (0\r\n\r\n)
	fmt.Println("所有数据块发送完成")
}
