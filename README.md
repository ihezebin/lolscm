# lolscm

NAME:
lolscm - 英雄联盟快捷键模拟键盘操作发送游戏对局内的消息

USAGE:
lolscm [global options] command [command options]

VERSION:
v1.0.0

AUTHOR:
hezebin <ihezebin@qq.com>

COMMANDS:
help, h Shows a list of commands or help for one command

GLOBAL OPTIONS:
--templates value, -t value 消息模板，JSON 格式：{"template1": []}
--shortcut value, -s value 消息快捷键，Windows 下通常为 Alt+Shift+模板序号（1-5），MacOS 下通常为 Command+Shift+模板序号（1-5）
--help, -h show help
--version, -v print the version

# MacOS 编译

```bash
go build -o lolscm main.go
```

# MacOS 下交叉编译 Windows

需要先安装 MinGW32

```bash
brew install mingw-w64

x86_64-w64-mingw32-gcc --version
```

```bash
GOARCH=amd64 GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o lolscm.exe main.go
```