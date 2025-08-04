package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

// 开发框架文档  https://docs.fyne.io/started/goroutines

func main() {
	fmt.Println("程序启动...") // 日志1

	a := app.New()
	fmt.Println("App初始化完成") // 日志2

	w := a.NewWindow("李漫")
	fmt.Println("窗口创建完成") // 日志3

	contentBox := container.NewVBox()
	contentBox.Add(widget.NewLabel("Hello, World!"))
	fmt.Println("内容设置完成") // 日志4

	w.Resize(fyne.NewSize(400, 400))

	// 实时更新时间
	output := canvas.NewText(time.Now().Format(time.TimeOnly), color.NRGBA{G: 0xff, A: 0xff})
	output.TextStyle.Monospace = true
	output.TextSize = 32
	contentBox.Add(output)
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			fyne.Do(func() {
				output.Text = time.Now().Format(time.TimeOnly)
				output.Refresh()
			})
		}
	}()

	w.SetContent(contentBox)
	w.ShowAndRun()
	fmt.Println("事件循环启动") // 日志5（这行可能不会执行，因为ShowAndRun会阻塞）
}
