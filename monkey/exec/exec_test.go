package exec

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
)

func TestExecSussess(t *testing.T) {
	// 恢复 patch 修改
	defer monkey.UnpatchAll()
	// mock 掉 exec.Command 返回的 *exec.Cmd 的 CombinedOutput 方法
	monkey.PatchInstanceMethod(
		reflect.TypeOf((*exec.Cmd)(nil)),
		"CombinedOutput", func(_ *exec.Cmd) ([]byte, error) {
			return []byte("results"), nil
		},
	)

	rc, output := call("any")
	if rc != 0 {
		t.Fail()
	}
	if output != "results" {
		t.Fail()
	}
}

func TestExecFailed(t *testing.T) {
	defer monkey.UnpatchAll()
	// 上次 mock 的是执行成功的情况，这一次轮到执行失败
	monkey.PatchInstanceMethod(
		reflect.TypeOf((*exec.Cmd)(nil)),
		"CombinedOutput", func(_ *exec.Cmd) ([]byte, error) {
			return []byte(""), fmt.Errorf("sth bad happened")
		},
	)
	// mock 掉 reportExecFailed 函数
	monkey.Patch(reportExecFailed, func(msg string) string {
		return msg
	})

	rc, output := call("any")
	if rc != 1 {
		t.Fail()
	}
	if output != "sth bad happened" {
		t.Fail()
	}
}
