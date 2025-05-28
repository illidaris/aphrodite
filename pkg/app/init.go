package app

import (
	"fmt"
	"runtime"

	_ "go.uber.org/automaxprocs" // 使用 Uber automaxprocs 正确设置 Go 程序线程数
)

func init() {
	println("aphrodite init")
	fmt.Printf("runtime info: os=%v,arch=%v,cpu=%v,gomaxprocs=%v\n",
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.GOMAXPROCS(0))
}
