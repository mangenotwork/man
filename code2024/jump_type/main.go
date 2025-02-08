package main

import (
	"log"
	"strings"
)

func main() {
	link := "https://www.ecosmos.vip/cn/activity/124?"
	log.Println(getJumpType(link))
}

var matchJumpType = map[string]int64{
	"/activity/":  2, // 2:活动详情
	"/metaverse/": 3, //  3:企业元宇宙详情
	// 4:直播详情
}

func getJumpType(link string) (int64, string) {
	var (
		jumpType  int64  = 1
		jumpValue string = ""
	)

	doMain := ".ecosmos.vip"
	if len(doMain) > 0 && string(doMain[0]) == "." {
		doMain = doMain[1:]
	}

	if strings.Contains(link, doMain) {
		for k, v := range matchJumpType {
			if strings.Contains(link, k) {
				jumpType = v

				tempList := strings.Split(link, k)
				if len(tempList) > 1 {
					tempStr := tempList[len(tempList)-1]
					jumpValue = strings.Split(tempStr, "?")[0]
				}

			}
		}
	}

	return jumpType, jumpValue
}
