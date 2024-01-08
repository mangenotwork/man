package main

import (
	"crypto/md5"
	"fmt"
	"time"
)

func main() {
	agentID := "abcd_RMB"
	agentKey := "be8e08f95357d215921f91c6a533f74d3194de52"
	//utc4, _ := time.LoadLocation(GetLinuxZone("GMT-04:00"))
	//tm := LocChange(time.Now(), utc4)
	keyG := fmt.Sprintf("%x", md5.Sum([]byte("19011"+agentID+agentKey)))
	fmt.Println(keyG) // 9424ffdd6de016a5f90a97a55d99c717
	querystring := "Token=Test006&GameId=1&Lang=zh-CN&AgentId=abcd_RMB"
	md5string := fmt.Sprintf("%x", md5.Sum([]byte(querystring+keyG)))
	randomText1 := "123456"
	randomText2 := "abcdef"
	key := randomText1 + md5string + randomText2
	fmt.Println(key) // 1234563e4c3b321eac16f20633f683be08d237abcdef
}

// 获取 linux 24 时区
func GetLinuxZone(tFmt string) string {
	var zoneMap = map[string]string{
		"GMT-11:00": "Pacific/Apia",
		"GMT-10:00": "Pacific/Honolulu",
		"GMT-09:30": "Pacific/Marquesas",
		"GMT-09:00": "Pacific/Gambier",
		"GMT-08:00": "Pacific/Pitcairn",
		"GMT-07:00": "America/Phoenix",
		"GMT-06:00": "America/Costa_Rica",
		"GMT-05:00": "America/Bogota",
		"GMT-04:00": "America/Guadeloupe",
		"GMT-03:00": "America/Argentina/Buenos_Aires",
		"GMT-02:00": "America/Noronha",
		"GMT-01:00": "Atlantic/Cape_Verde",
		"GMT+00:00": "Africa/Dakar",
		"GMT+01:00": "Africa/Kinshasa",
		"GMT+02:00": "Africa/Harare",
		"GMT+03:00": "Africa/Mogadishu",
		"GMT+04:00": "Indian/Mauritius",
		"GMT+04:30": "Asia/Kabul",
		"GMT+05:00": "Asia/Samarkand",
		"GMT+05:30": "Asia/Colombo", //    mx300:Asia/Calcutta
		"GMT+06:00": "Asia/Dhaka",
		"GMT+06:30": "Asia/Rangoon",
		"GMT+07:00": "Asia/Vientiane",
		"GMT+08:00": "Asia/Shanghai",
		"GMT+09:00": "Asia/Tokyo",
		"GMT+10:00": "Pacific/Guam",
		"GMT+11:00": "Pacific/Noumea",
		"GMT+11:30": "Pacific/Norfolk",
		"GMT+12:00": "Pacific/Nauru",
		"GMT+13:00": "Pacific/Tongatapu",
	}

	s, ok := zoneMap[tFmt]
	if !ok {
		return "没有找到对应的格式"
	}

	return s
}

var (
	locationUTC = time.UTC
)

// 指定时间变换时区
func LocChange(t time.Time, loc *time.Location) time.Time {
	if loc == nil { // 不设定,则使用标准时,对标Parse将时间解释为UTC时间
		loc = locationUTC
	}

	return t.In(loc)
}
