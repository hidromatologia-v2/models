package logs

import (
	"log"
	"strings"
)

func LogOnError(err error) {
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		log.Print(err)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
