package util

import "log"

func Assert(condition bool, msg string){
	if condition {log.Println(msg + " PASS")}else{log.Println(msg + " FAILED")}
}
