package util

import (
	"fmt"
	"github.com/liuyp5181/base/log"
	"runtime"
)

func Go(f func()) {
	go func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			errStr := fmt.Sprintf("panic, err = %v\n%v", err, string(buf))
			log.Error(errStr)
		}
		f()
	}()
}
