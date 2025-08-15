package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// 自定义数据结构：表示目录树节点
type TreeNode struct {
	ID       string      // 节点唯一标识
	Name     string      // 节点显示名称
	Children []*TreeNode // 子节点
	IsFile   bool        // 区分文件(叶子)和文件夹(分支)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Tree组件数据绑定示例")

	// 1. 准备数据源（模拟目录结构）
	rootNodes := []*TreeNode{
		{
			ID:   "docs",
			Name: "文档",
			Children: []*TreeNode{
				{ID: "doc1", Name: "报告.docx", IsFile: true},
				{ID: "doc2", Name: "计划.xlsx", IsFile: true},
			},
			IsFile: false,
		},
		{
			ID:   "images",
			Name: "图片",
			Children: []*TreeNode{
				{ID: "img1", Name: "截图1.png", IsFile: true},
				{ID: "img2", Name: "截图2.png", IsFile: true},
			},
			IsFile: false,
		},
		{ID: "readme", Name: "说明.txt", IsFile: true},
	}

	// 2. 构建节点ID与数据的映射
	nodeMap := make(map[string]*TreeNode)
	var buildNodeMap func(parentID string, nodes []*TreeNode)
	buildNodeMap = func(parentID string, nodes []*TreeNode) {
		for _, node := range nodes {
			fullID := parentID + "/" + node.ID
			if parentID == "" {
				fullID = node.ID // 根节点ID
			}
			nodeMap[fullID] = node

			if len(node.Children) > 0 {
				buildNodeMap(fullID, node.Children)
			}
		}
	}
	buildNodeMap("", rootNodes)

	// 3. 创建Tree组件并绑定数据
	tree := widget.NewTree(
		// 回调1：获取子节点ID列表
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" { // 根节点
				ids := make([]widget.TreeNodeID, len(rootNodes))
				for i, node := range rootNodes {
					ids[i] = node.ID
				}
				return ids
			}

			if node, ok := nodeMap[id]; ok {
				ids := make([]widget.TreeNodeID, len(node.Children))
				for i, child := range node.Children {
					ids[i] = id + "/" + child.ID
				}
				return ids
			}
			return nil
		},

		// 回调2：判断是否为叶子节点
		func(id widget.TreeNodeID) bool {
			if node, ok := nodeMap[id]; ok {
				return node.IsFile
			}
			return false
		},

		// 回调3：创建节点UI组件（关键修改：提前设置颜色）
		func(branch bool) fyne.CanvasObject {
			if branch {
				// 分支节点：创建带样式的标签（加粗+黑色）
				label := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				// 用CanvasObject包装以支持颜色设置
				return label
			} else {
				// 叶子节点：使用canvas.Text直接设置颜色
				text := canvas.NewText("", color.Gray{80})
				text.TextSize = 14 // 与Label保持一致的字体大小
				return text
			}
		},

		// 回调4：更新节点显示内容
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			if node, ok := nodeMap[id]; ok {
				if branch {
					// 分支节点：更新Label文本
					o.(*widget.Label).SetText(node.Name)
				} else {
					// 叶子节点：更新canvas.Text文本
					o.(*canvas.Text).Text = node.Name
					o.(*canvas.Text).Refresh() // 刷新显示
				}
			}
		},
	)

	// 4. 节点选择事件
	tree.OnSelected = func(id widget.TreeNodeID) {
		if node, ok := nodeMap[id]; ok {
			myWindow.SetTitle("选中: " + node.Name)
		}
	}

	// 布局设置
	leftContainer := container.New(layout.NewStackLayout(), tree)

	middle := canvas.NewText("选择左侧文件/文件夹查看详情", color.Black)
	mainContainer := container.NewHSplit(leftContainer, container.NewVBox(middle))
	mainContainer.SetOffset(0.25)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.ShowAndRun()
}
