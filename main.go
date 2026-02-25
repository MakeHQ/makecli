/**
 * [INPUT]: 依赖 cmd 包的 Execute、internal/build 的 Version/Date
 * [OUTPUT]: 可执行二进制 makecli
 * [POS]: 程序入口，将构建元数据传入命令层
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package main

import (
	"os"

	"github.com/MakeHQ/makecli/cmd"
	"github.com/MakeHQ/makecli/internal/build"
)

func main() {
	if err := cmd.Execute(build.Version, build.Date); err != nil {
		os.Exit(1)
	}
}
