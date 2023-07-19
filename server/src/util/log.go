package util

import "log"

func LogError(err error, errMsg string) {
	if errMsg != "" {
		log.Println(errMsg)
		return
	}
	if err != nil {
		log.Println(err.Error())
	}
}
