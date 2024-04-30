package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

/*

snotify 可以用于各种场景，例如：

实时构建系统：当源代码发生变化时，自动重新编译和部署应用。
文件同步工具：监视文件或目录的变化，并将更改实时同步到其他位置。
数据库存储跟踪：当数据库日志文件发生变动时，触发相应的处理逻辑。
编辑器功能扩展：为文本编辑器添加实时保存、语法检查等功能。

Create：文件或目录被创建
Write：文件内容被写入
Remove：文件或目录被删除
Rename：文件或目录被重命名
Chmod：文件或目录权限发生变化

*/

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("/tmp")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
