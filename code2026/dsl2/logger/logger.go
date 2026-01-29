package logger

import "log"

var IsDebug = false

func Debug(v ...interface{}) {
	if IsDebug {
		value := make([]interface{}, 0)
		value = append(value, "[DEBUG]")
		value = append(value, v...)
		log.Println(value...)
	}

}
