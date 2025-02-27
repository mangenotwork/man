package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	// 创建一个context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	/*
		// 设置Chrome会话上下文和超时时间
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
	*/

	// 创建一个任务列表，用于指定要执行的浏览器操作
	var nodes = ""
	err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Navigate(`https://jobs.51job.com/all/coAmZSPgRuBTlSMwFkVzA.html`), // 导航到目标网址
			chromedp.WaitVisible(`h1`, chromedp.ByQuery),                               // 等待h1元素出现
			chromedp.OuterHTML(`html`, &nodes, chromedp.ByQuery),                       // 获取h1元素的HTML代码
		},
	)
	if err != nil {
		log.Fatalf("获取网页失败: %v", err)
	}

	fmt.Println("页面标题:", nodes) // 打印获取到的标题

	//jg := RegFindAll(`url: '(.*?)'`, nodes)
	//fmt.Println("jg = ", jg)
	//if len(jg) > 0 && len(jg[0]) > 1 {
	//	var nodes2 = ""
	//	err = chromedp.Run(ctx,
	//		chromedp.Tasks{
	//			chromedp.Navigate(jg[0][1]),                         // 导航到目标网址
	//			chromedp.WaitVisible(`h1`, chromedp.ByQuery),        // 等待h1元素出现
	//			chromedp.OuterHTML(`h1`, &nodes2, chromedp.ByQuery), // 获取h1元素的HTML代码
	//		},
	//	)
	//	if err != nil {
	//		log.Fatalf("获取网页失败: %v", err)
	//	}
	//	fmt.Println("页面标题:", nodes2) // 打印获取到的标题
	//}
}

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}

func case1() {

	// 设置Chrome会话上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//var audioURL string

	// 创建一个新的Chrome会话
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//关闭无头模式，方便调试
		chromedp.Flag("headless", false),
		//防止监测webdriver
		chromedp.Flag("enable-automation", false),
		//禁用 blink 特征
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		//忽略浏览器的风险提示（但好像并没什么用）
		chromedp.Flag("ignore-certificate-errors", true),
		//关闭浏览器声音（也没用）
		chromedp.Flag("mute-audio", false),
		//禁用GPU加速（有时可以解决一些渲染问题）
		chromedp.Flag("disable-gpu", true),
		//设置浏览器尺寸
		chromedp.WindowSize(1150, 1000),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()
	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer taskCancel()

	// 启动浏览器并导航
	err := chromedp.Run(taskCtx,
		//打开youtube
		chromedp.Navigate("https://xxxxxxxxx"),
		//等待按钮绘制完成
		chromedp.WaitVisible("a[title='Shorts']", chromedp.ByQuery),
		//点击短视频按钮
		chromedp.Click("a[title='Shorts']", chromedp.ByQuery),
	)

	err = chromedp.Run(taskCtx, chromedp.Sleep(time.Hour*1))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}

/*
将网页截取成图片有两个函数：chromedp.Screenshot和chromedp.FullScreenshot。其中chromedp.Screenshot是按网页中的某个div的元素截取。而chromedp.FullScreenshot是截取整个网页。
*/
func case2() {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`https://pkg.go.dev/`, `img.Homepage-logo`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(`https://brank.as/`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

/*

其他功能
模拟表单提交：可以使用chromedp.Submit函数模拟表单提交。
模拟鼠标滚动：可以使用chromedp.ScrollIntoView函数模拟鼠标滚动。
模拟键盘输入：可以使用chromedp.KeyEvent函数模拟键盘输入。

所以示例在:  https://github.com/chromedp/examples

*/
