package common

//CompareBytes 比较两个byte数组
import "runtime"

//ConfigRuntime 多核运行
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
}
