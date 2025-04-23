package code2025

import (
	"strings"
	"time"
)

var timeFormat = map[string]string{
	"yyyy-mm-dd hh:mm:ss": "2006-01-02 15:04:05",
	"yyyy-mm-dd hh:mm":    "2006-01-02 15:04",
	"yyyy-mm-dd hh":       "2006-01-02 15",
	"yyyy-mm-dd":          "2006-01-02",
	"yyyy-mm":             "2006-01",
	"yyyy":                "2006",
	"mm-dd hh:mm:ss":      "01-02 15:04:05",
}

var timeFormat2 = map[string]string{
	"yyyy-mm-dd hh:mm:ss": "2006-01-02 15:04:05",
	"yyyy-mm-dd hh:mm":    "2006-01-02 15:04",
	"yyyy-mm-dd hh":       "2006-01-02 15",
	"yyyy-mm-dd":          "2006-01-02",
	"yyyy-mm":             "2006-01",
	"yyyy":                "2006",
	"mm-dd hh:mm:ss":      "01-02 15:04:05",
}

// timeData 时间值
// timeStr 时间格式  yyyy-mm-dd hh:mm:ss  yyyy/mm/dd hh:mm:ss  ....
func TimeFormat(timeData time.Time, timeStr string) string {

	// 定义替换规则
	replacer := strings.NewReplacer(
		"yyyy", "2006",
		"mm", "01",
		"dd", "02",
		"hh", "15",
		"mm", "04",
		"ss", "05",
	)

	// 批量替换
	result := replacer.Replace(timeStr)

	return timeData.Format(result)
}
