package main

import "runtime"

func isMobile() bool {
	return runtime.GOOS == "android" || runtime.GOOS == "ios"
}
