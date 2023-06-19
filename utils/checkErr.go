package utils

import "log"

func CheckErr(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}
