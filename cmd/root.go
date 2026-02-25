/**
 * [INPUT]: 依赖 github.com/spf13/cobra
 * [OUTPUT]: 对外提供 Execute 函数、rootCmd 根命令
 * [POS]: cmd 模块的入口，挂载所有子命令
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "makecli",
	Short: "makecli — make your workflow faster",
}

// Execute 是程序入口，由 main.go 调用
func Execute(version, buildDate string) error {
	rootCmd.AddCommand(newVersionCmd(version, buildDate))
	return rootCmd.Execute()
}
