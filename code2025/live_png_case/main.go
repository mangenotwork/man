package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	uploadDir   = "./"
	staticDir   = "./"      // 存放前端HTML文件的目录
	maxFileSize = 100 << 20 // 100MB
)

func main() {
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("无法创建上传目录: %v", err)
	}

	// 路由设置
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/images/", serveImage)

	// 启动服务器
	log.Println("服务器运行在 http://localhost:17777")
	log.Fatal(http.ListenAndServe(":17777", nil))
}

// 提供前端页面
func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 读取并返回HTML文件
	content, err := os.ReadFile(filepath.Join(staticDir, "index.html"))
	if err != nil {
		http.Error(w, "无法加载页面", http.StatusInternalServerError)
		log.Printf("读取HTML文件错误: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// 处理文件上传
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 设置最大文件大小
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		http.Error(w, "文件过大", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "未选择文件", http.StatusBadRequest)
		return
	}

	// 保存所有文件
	for _, fileHeader := range files {
		// 打开上传的文件
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "无法打开文件", http.StatusInternalServerError)
			log.Printf("打开文件错误: %v", err)
			return
		}
		defer file.Close()

		// 创建保存路径
		dstPath := filepath.Join(uploadDir, fileHeader.Filename)

		// 创建目标文件
		dstFile, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "无法创建文件", http.StatusInternalServerError)
			log.Printf("创建文件错误: %v", err)
			return
		}
		defer dstFile.Close()

		// 复制文件内容
		if _, err := io.Copy(dstFile, file); err != nil {
			http.Error(w, "文件保存失败", http.StatusInternalServerError)
			log.Printf("复制文件错误: %v", err)
			return
		}

		log.Printf("文件保存成功: %s", fileHeader.Filename)
	}

	// 处理实况图关联信息
	livePairs := r.MultipartForm.Value["live_pairs"]
	for _, pair := range livePairs {
		parts := strings.Split(pair, "|")
		if len(parts) == 2 {
			log.Printf("实况图关联: %s -> %s", parts[0], parts[1])
			// 这里可以添加逻辑将关联信息保存到数据库
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("上传成功"))
}

const liveDir = "./live"

// serveImage 处理图片请求
func serveImage(w http.ResponseWriter, r *http.Request) {
	// 获取URL路径中的图片文件名
	filename := strings.TrimPrefix(r.URL.Path, "/images/")
	if filename == "" {
		http.Error(w, "图片文件名不能为空", http.StatusBadRequest)
		return
	}

	// 构建图片文件的完整路径
	filePath := filepath.Join(liveDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		log.Printf("检查图片文件时出错: %v", err)
		return
	}

	// 设置响应头，根据文件扩展名自动设置Content-Type
	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 打开文件并发送给客户端
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "无法打开图片", http.StatusInternalServerError)
		log.Printf("打开图片文件时出错: %v", err)
		return
	}
	defer file.Close()

	// 将文件内容复制到响应中
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("发送图片时出错: %v", err)
	}
}
