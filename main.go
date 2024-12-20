package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var debounceTimer *time.Timer

func getCurrentOS() string {
	return runtime.GOOS
}

func main() {
	os := getCurrentOS()
	// Alt+Shift+1
	k := []string{"alt", "shift", "1"}
	killK := []string{"command", "shift", "0"}
	if os == "darwin" {
		// Command+Shift+1
		k = []string{"command", "shift", "1"}
		killK = []string{"command", "shift", "0"}
	}

	time.After(time.Millisecond)

	hook.Register(hook.KeyDown, k, func(e hook.Event) {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		// 如果定时器已存在，取消它
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		// 重设定时器
		debounceTimer = time.AfterFunc(1000*time.Millisecond, func() {
			// 在延迟时间后触发一次输出
			sendMessage()
		})
	})

	hook.Register(hook.KeyDown, killK, func(e hook.Event) {
		hook.End()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func sendMessage() {
	fmt.Println("发送消息")
	robotgo.KeyTap("enter")
	robotgo.MilliSleep(50)
	robotgo.TypeStr("你好，欢迎使用 HLOLT!")
	robotgo.MilliSleep(50)
	robotgo.KeyTap("enter")
}
