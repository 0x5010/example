package exec

import (
	"os/exec"
)

// 假如我们要测试函数 call
func call(cmd string) (int, string) {
	bytes, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return 1, reportExecFailed(err.Error())
	}
	output := string(bytes)
	return 0, output
}

func reportExecFailed(msg string) string {
	return "error"
}
