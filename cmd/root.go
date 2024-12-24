package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
	hook "github.com/robotn/gohook"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

type TemplatesFlag map[string][]string

// Set 实现 cli.Generic 接口
func (m *TemplatesFlag) Set(value string) error {
	if m == nil {
		return errors.New("templates flag is nil")
	}
	tempM := make(TemplatesFlag)
	err := json.Unmarshal([]byte(value), &tempM)
	if err != nil {
		return errors.Wrap(err, "unmarshal error")
	}

	for k, v := range tempM {
		if !strings.HasPrefix(k, "template") {
			return errors.Errorf("invalid template key: %s, must start with template", k)
		}

		var index int64
		index, err = strconv.ParseInt(strings.TrimPrefix(k, "template"), 10, 64)
		if err != nil {
			return errors.Wrapf(err, "invalid template key: %s, must end with number", k)
		}

		if index < 1 || index > 5 {
			return errors.Errorf("invalid template key: %s, must between 1 and 5", k)
		}

		(*m)[strconv.Itoa(int(index))] = v
	}

	return nil
}

func (m *TemplatesFlag) String() string {
	var str []string
	for k, v := range *m {
		str = append(str, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(str, ", ")
}

var (
	templates = make(TemplatesFlag)
	shortcut  string
)

func Run(ctx context.Context) error {
	app := &cli.App{
		Name:    "lolscm",
		Version: "v1.0.0",
		Usage:   "英雄联盟快捷键模拟键盘操作发送游戏对局内的消息",
		Authors: []*cli.Author{
			{Name: "hezebin", Email: "ihezebin@qq.com"},
		},
		Flags: []cli.Flag{
			&cli.GenericFlag{
				Destination: &templates,
				Name:        "templates",
				Aliases:     []string{"t"},
				Value:       nil,
				Usage:       `消息模板，JSON 格式：{"template1": []}`,
			},
			&cli.StringFlag{
				Required:    true,
				Destination: &shortcut,
				Name:        "shortcut",
				Aliases:     []string{"s"},
				Usage:       `消息快捷键，Windows 下通常为 Alt+Shift+模板序号（1-5），MacOS 下通常为 Command+Shift+模板序号（1-5）`,
				Action: func(c *cli.Context, s string) error {
					keys := strings.Split(s, "+")
					if len(keys) > 2 {
						return errors.New("invalid shortcut len")
					}

					allowedKey := []string{"alt", "option", "command", "shift", "1", "2", "3", "4", "5"}
					for _, key := range keys {
						if !slices.Contains(allowedKey, strings.ToLower(key)) {
							return errors.Errorf("invalid shortcut key: %s", key)
						}
					}
					return nil
				},
			},
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			var debounceTimer *time.Timer
			for index, messages := range templates {
				hook.Register(hook.KeyDown, append(strings.Split(strings.ToLower(shortcut), "+"), index), func(e hook.Event) {
					// 如果定时器已存在，取消它
					if debounceTimer != nil {
						debounceTimer.Stop()
					}

					// 重设定时器
					debounceTimer = time.AfterFunc(time.Second, func() {
						// 在延迟时间后触发一次输出
						for _, message := range messages {
							msgs := strings.Split(message, "\n")
							for _, msg := range msgs {
								err := robotgo.KeyTap("enter")
								if err != nil {
									logger.WithError(err).Errorf(ctx, "key tap enter error")
									continue
								}
								robotgo.MilliSleep(50)
								robotgo.TypeStr(msg)
								robotgo.MilliSleep(50)
								err = robotgo.KeyTap("enter")
								if err != nil {
									logger.WithError(err).Errorf(ctx, "key tap enter error")
									continue
								}
								time.Sleep(time.Millisecond * 100)
							}
						}
					})
				})
			}

			// 关闭程序
			hook.Register(hook.KeyDown, append(strings.Split(shortcut, "+"), "0"), func(e hook.Event) {
				hook.End()
			})

			s := hook.Start()
			<-hook.Process(s)

			return nil
		},
	}

	return app.Run(os.Args)
}
