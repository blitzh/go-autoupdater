package util

import "runtime"

func GOOS() string   { return runtime.GOOS }
func GOARCH() string { return runtime.GOARCH }
