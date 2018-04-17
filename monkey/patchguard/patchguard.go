package main

import (
	"fmt"
	"strings"

	"github.com/bouk/monkey"
)

func main() {
	var guard *monkey.PatchGuard
	guard = monkey.Patch(fmt.Println, func(a ...interface{}) (int, error) {
		s := make([]interface{}, len(a))
		for i, v := range a {
			s[i] = strings.Replace(fmt.Sprint(v), "hell", "*bleep*", -1)
		}

		// 取消patch
		guard.Unpatch()
		defer guard.Restore()
		// 使用默认的fmt.Println
		return fmt.Println(s...)
	})
	fmt.Println("what the hell?") // what the *bleep*?
	fmt.Println("what the hell?") // what the *bleep*?
}
