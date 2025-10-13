package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 1. 翻译映射表
var translations = map[string]map[string]string{
	"en": {
		"appTitle":  "Auto Multi-language",
		"welcome":   "Welcome to the example",
		"inputHint": "Enter text here",
		"switchBtn": "Switch to Chinese",
		"submitBtn": "Submit",
	},
	"zh-CN": {
		"appTitle":  "自动多语言示例",
		"welcome":   "欢迎使用本示例",
		"inputHint": "在此输入文本",
		"switchBtn": "切换到英文",
		"submitBtn": "提交",
	},
}

// 2. 全局状态：当前语言+组件映射表（自动关联翻译键和组件）
var (
	currentLang = "en"
	// 存储需要翻译的组件：key是翻译键，value是组件和更新方法
	translatableComponents = struct {
		sync.Mutex
		list map[string][]interface{} // 支持多种组件类型（Label/Button/Entry等）
	}{
		list: make(map[string][]interface{}),
	}
)

// 3. 核心函数：注册组件与翻译键（一行代码完成关联）
// 支持 Label/Button/Entry.PlaceHolder 等常见组件
func registerTranslatable(key string, component interface{}) {
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	translatableComponents.list[key] = append(translatableComponents.list[key], component)
}

// 4. 翻译函数：获取当前语言的文本
func t(key string) string {
	if langMap, ok := translations[currentLang]; ok {
		if text, ok := langMap[key]; ok {
			return text
		}
	}
	return translations["en"][key] // 默认英文
}

// 5. 切换语言并自动更新所有组件（核心：批量刷新）
func switchLang() {
	// 切换当前语言
	if currentLang == "en" {
		currentLang = "zh-CN"
	} else {
		currentLang = "en"
	}

	// 遍历所有注册的组件，自动更新文本（无需手动写每个组件的刷新逻辑）
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	for key, components := range translatableComponents.list {
		text := t(key)
		for _, comp := range components {
			switch c := comp.(type) {
			case *widget.Label:
				c.SetText(text)
			case *widget.Button:
				c.SetText(text)
			case *widget.Entry:
				c.SetPlaceHolder(text) // Entry 适配占位文本
			}
		}
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow(t("appTitle"))

	// 6. 创建原生组件并关联翻译键（仅需一行注册代码）
	welcomeLabel := widget.NewLabel(t("welcome"))
	registerTranslatable("welcome", welcomeLabel) // 关联标签
	welcomeLabel2 := widget.NewLabel(t("welcome"))
	registerTranslatable("welcome", welcomeLabel2) // 关联标签

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder(t("inputHint"))
	registerTranslatable("inputHint", inputEntry) // 关联输入框占位符

	submitBtn := widget.NewButton(t("submitBtn"), func() {})
	registerTranslatable("submitBtn", submitBtn) // 关联按钮

	switchBtn := widget.NewButton(t("switchBtn"), switchLang)
	registerTranslatable("switchBtn", switchBtn) // 关联切换按钮

	// 布局
	content := container.NewVBox(
		welcomeLabel,
		welcomeLabel2,
		inputEntry,
		submitBtn,
		switchBtn,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 250))
	myWindow.ShowAndRun()
}
