package http_kuai_multipart

import (
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析 multipart 表单，设置最大内存为 10 MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "解析表单失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 获取普通表单字段
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")

	// 输出文本字段内容
	w.Write([]byte("收到的文本数据:\n"))
	w.Write([]byte("name: " + name + "\n"))
	w.Write([]byte("description: " + description + "\n\n"))

	// 获取文件字段
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "获取文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 输出文件信息
	w.Write([]byte("收到的文件数据:\n"))
	w.Write([]byte("文件名: " + handler.Filename + "\n"))
	w.Write([]byte("文件大小: " + string(handler.Size) + " bytes\n"))
	w.Write([]byte("文件类型: " + handler.Header.Get("Content-Type") + "\n\n"))

	// 保存文件到本地
	dst, err := os.Create("upload_" + handler.Filename)
	if err != nil {
		http.Error(w, "创建文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "保存文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("文件上传成功!"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	println("服务器启动在 :8080 端口")
	http.ListenAndServe(":8080", nil)
}
